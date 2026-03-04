package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"go.uber.org/zap"
)

type Migration struct {
	Version int
	Name    string
	UpSQL   string
}

type Migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewMigrator(db *sql.DB, logger *zap.Logger) *Migrator {
	return &Migrator{db: db, logger: logger}
}

func (m *Migrator) RunMigrations(migrationsDir string) error {
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	applied, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("get applied migrations: %w", err)
	}

	migrations, err := m.loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("load migrations: %w", err)
	}

	for _, mig := range migrations {
		if applied[mig.Version] {
			continue
		}
		if err := m.runMigration(mig); err != nil {
			return fmt.Errorf("run migration %d: %w", mig.Version, err)
		}
	}

	m.logger.Info("Migrations completed")
	return nil
}

func (m *Migrator) createMigrationsTable() error {
	_, err := m.db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version    INTEGER PRIMARY KEY,
		name       VARCHAR(255) NOT NULL,
		applied_at TIMESTAMPTZ DEFAULT NOW()
	)`)
	return err
}

func (m *Migrator) getAppliedMigrations() (map[int]bool, error) {
	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		applied[v] = true
	}
	return applied, rows.Err()
}

func (m *Migrator) loadMigrations(dir string) ([]Migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			m.logger.Warn("Migrations directory not found, skipping", zap.String("dir", dir))
			return nil, nil
		}
		return nil, err
	}

	var migrations []Migration
	for _, entry := range entries {
		if entry.Type() == fs.ModeDir || !strings.HasSuffix(entry.Name(), ".up.sql") {
			continue
		}
		version, name, err := parseMigrationFilename(entry.Name())
		if err != nil {
			m.logger.Warn("Skipping invalid migration file", zap.String("file", entry.Name()))
			continue
		}
		sql, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", entry.Name(), err)
		}
		migrations = append(migrations, Migration{Version: version, Name: name, UpSQL: string(sql)})
	}

	sort.Slice(migrations, func(i, j int) bool { return migrations[i].Version < migrations[j].Version })
	return migrations, nil
}

func (m *Migrator) runMigration(mig Migration) error {
	m.logger.Info("Applying migration", zap.Int("version", mig.Version), zap.String("name", mig.Name))
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(mig.UpSQL); err != nil {
		return fmt.Errorf("execute SQL: %w", err)
	}
	if _, err := tx.Exec("INSERT INTO schema_migrations (version, name) VALUES ($1, $2)", mig.Version, mig.Name); err != nil {
		return fmt.Errorf("record migration: %w", err)
	}
	return tx.Commit()
}

func parseMigrationFilename(filename string) (int, string, error) {
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) < 2 {
		return 0, "", fmt.Errorf("invalid migration filename: %s", filename)
	}
	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid version in %s: %w", filename, err)
	}
	name := strings.TrimSuffix(parts[1], ".up.sql")
	return version, name, nil
}
