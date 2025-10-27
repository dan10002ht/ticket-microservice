package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

// Migration represents a database migration
type Migration struct {
	Version int
	Name    string
	UpSQL   string
	DownSQL string
}

// Migrator handles database migrations
type Migrator struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, logger *zap.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// RunMigrations runs all pending migrations
func (m *Migrator) RunMigrations(migrationsDir string) error {
	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Load migration files
	migrations, err := m.loadMigrations(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range migrations {
		if _, applied := appliedMigrations[migration.Version]; !applied {
			if err := m.runMigration(migration); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", migration.Version, err)
			}
		}
	}

	m.logger.Info("All migrations completed successfully")
	return nil
}

// createMigrationsTable creates the migrations tracking table
func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	return err
}

// getAppliedMigrations returns a map of applied migration versions
func (m *Migrator) getAppliedMigrations() (map[int]bool, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		applied[version] = true
	}

	return applied, rows.Err()
}

// loadMigrations loads migration files from directory
func (m *Migrator) loadMigrations(migrationsDir string) ([]Migration, error) {
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, err
	}

	var migrations []Migration
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".up.sql") {
			version, name, err := parseMigrationFilename(file.Name())
			if err != nil {
				m.logger.Warn("Skipping invalid migration file", zap.String("file", file.Name()))
				continue
			}

			upSQL, err := m.readSQLFile(filepath.Join(migrationsDir, file.Name()))
			if err != nil {
				return nil, fmt.Errorf("failed to read up migration %s: %w", file.Name(), err)
			}

			downFile := strings.Replace(file.Name(), ".up.sql", ".down.sql", 1)
			downSQL, err := m.readSQLFile(filepath.Join(migrationsDir, downFile))
			if err != nil {
				m.logger.Warn("No down migration found", zap.String("file", downFile))
				downSQL = ""
			}

			migrations = append(migrations, Migration{
				Version: version,
				Name:    name,
				UpSQL:   upSQL,
				DownSQL: downSQL,
			})
		}
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// runMigration executes a single migration
func (m *Migrator) runMigration(migration Migration) error {
	m.logger.Info("Running migration",
		zap.Int("version", migration.Version),
		zap.String("name", migration.Name),
	)

	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Execute migration SQL
	if _, err := tx.Exec(migration.UpSQL); err != nil {
		return fmt.Errorf("failed to execute migration SQL: %w", err)
	}

	// Record migration as applied
	query := "INSERT INTO schema_migrations (version, name) VALUES ($1, $2)"
	if _, err := tx.Exec(query, migration.Version, migration.Name); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	m.logger.Info("Migration completed successfully",
		zap.Int("version", migration.Version),
		zap.String("name", migration.Name),
	)

	return nil
}

// parseMigrationFilename parses migration filename to extract version and name
func parseMigrationFilename(filename string) (int, string, error) {
	// Expected format: 001_create_tickets_table.up.sql
	parts := strings.Split(filename, "_")
	if len(parts) < 2 {
		return 0, "", fmt.Errorf("invalid migration filename format")
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid version number: %w", err)
	}

	name := strings.Join(parts[1:], "_")
	name = strings.TrimSuffix(name, ".up.sql")

	return version, name, nil
}

// readSQLFile reads SQL content from file
func (m *Migrator) readSQLFile(filepath string) (string, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
