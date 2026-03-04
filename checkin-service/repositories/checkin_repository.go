package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"checkin-service/database"
	"checkin-service/models"
)

type CheckinRepository struct {
	db     *database.DB
	logger *zap.Logger
}

func NewCheckinRepository(db *database.DB) *CheckinRepository {
	return &CheckinRepository{db: db}
}

// ExistsByTicketID returns true if a check-in record already exists for this ticket.
func (r *CheckinRepository) ExistsByTicketID(ticketID string) (bool, error) {
	var count int
	err := r.db.QueryRowx(
		"SELECT COUNT(1) FROM checkins WHERE ticket_id = $1", ticketID,
	).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("ExistsByTicketID: %w", err)
	}
	return count > 0, nil
}

// Create inserts a new check-in record and returns it with the generated ID.
func (r *CheckinRepository) Create(c *models.CheckIn) (*models.CheckIn, error) {
	if err := c.Validate(); err != nil {
		return nil, err
	}

	query := `
		INSERT INTO checkins (ticket_id, event_id, user_id, staff_id, qr_code, status,
		                      check_in_time, device_id, gate, notes)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), $7, $8, $9)
		RETURNING id, check_in_time, created_at`

	err := r.db.QueryRowx(query,
		c.TicketID, c.EventID, c.UserID, c.StaffID, c.QRCode, c.Status,
		c.DeviceID, c.Gate, c.Notes,
	).Scan(&c.ID, &c.CheckInTime, &c.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("Create checkin: %w", err)
	}
	return c, nil
}

// GetByID retrieves a check-in record by its UUID.
func (r *CheckinRepository) GetByID(id string) (*models.CheckIn, error) {
	var c models.CheckIn
	err := r.db.QueryRowx("SELECT * FROM checkins WHERE id = $1", id).StructScan(&c)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("GetByID: %w", err)
	}
	return &c, nil
}

// ListByEvent retrieves paginated check-ins for an event with optional filters.
func (r *CheckinRepository) ListByEvent(eventID, userID, gate string, page, limit int) ([]*models.CheckIn, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Build dynamic WHERE clause
	where := "WHERE event_id = $1"
	args := []interface{}{eventID}
	idx := 2

	if userID != "" {
		where += fmt.Sprintf(" AND user_id = $%d", idx)
		args = append(args, userID)
		idx++
	}
	if gate != "" {
		where += fmt.Sprintf(" AND gate = $%d", idx)
		args = append(args, gate)
		idx++
	}

	// Total count
	var total int
	if err := r.db.QueryRowx(
		fmt.Sprintf("SELECT COUNT(1) FROM checkins %s", where), args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("ListByEvent count: %w", err)
	}

	// Paginated rows
	args = append(args, limit, offset)
	rows, err := r.db.Queryx(
		fmt.Sprintf("SELECT * FROM checkins %s ORDER BY check_in_time DESC LIMIT $%d OFFSET $%d",
			where, idx, idx+1),
		args...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("ListByEvent query: %w", err)
	}
	defer rows.Close()

	var checkins []*models.CheckIn
	for rows.Next() {
		var c models.CheckIn
		if err := rows.StructScan(&c); err != nil {
			return nil, 0, err
		}
		checkins = append(checkins, &c)
	}
	return checkins, total, rows.Err()
}

// GetEventStats returns aggregated check-in counts for an event.
func (r *CheckinRepository) GetEventStats(eventID string) (*models.EventStats, error) {
	stats := &models.EventStats{
		EventID: eventID,
		ByGate:  make(map[string]int),
	}

	row := r.db.QueryRowx(`
		SELECT
			COUNT(*)                AS total,
			COUNT(DISTINCT ticket_id) AS unique_tickets,
			MAX(check_in_time)      AS last_checkin
		FROM checkins
		WHERE event_id = $1
	`, eventID)

	var lastCheckin *string // nullable timestamp
	if err := row.Scan(&stats.TotalCheckins, &stats.UniqueTickets, &lastCheckin); err != nil {
		return nil, fmt.Errorf("GetEventStats aggregate: %w", err)
	}

	// Per-gate breakdown
	gateRows, err := r.db.Queryx(`
		SELECT COALESCE(gate, 'unknown'), COUNT(*) AS cnt
		FROM checkins
		WHERE event_id = $1
		GROUP BY gate
	`, eventID)
	if err != nil {
		return nil, fmt.Errorf("GetEventStats by gate: %w", err)
	}
	defer gateRows.Close()

	for gateRows.Next() {
		var gate string
		var cnt int
		if err := gateRows.Scan(&gate, &cnt); err != nil {
			return nil, err
		}
		stats.ByGate[gate] = cnt
	}
	return stats, gateRows.Err()
}
