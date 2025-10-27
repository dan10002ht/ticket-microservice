package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"

	"ticket-service/models"
)

// BookingSessionRepository handles database operations for booking sessions
type BookingSessionRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewBookingSessionRepository creates a new booking session repository
func NewBookingSessionRepository(db *sqlx.DB, logger *zap.Logger) *BookingSessionRepository {
	return &BookingSessionRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new booking session
func (r *BookingSessionRepository) Create(ctx context.Context, session *models.BookingSession) error {
	query := `
		INSERT INTO booking_sessions (
			id, user_id, event_id, session_token, status, seat_count,
			total_amount, currency, expires_at, ip_address, user_agent,
			metadata, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`

	// Generate UUID if not provided
	if session.ID == "" {
		session.ID = uuid.New().String()
	}

	_, err := r.db.ExecContext(ctx, query,
		session.ID, session.UserID, session.EventID, session.SessionToken,
		session.Status, session.SeatCount, session.TotalAmount, session.Currency,
		session.ExpiresAt, session.IPAddress, session.UserAgent,
		session.Metadata, session.CreatedBy, session.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to create booking session",
			zap.String("session_id", session.ID),
			zap.String("user_id", session.UserID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create booking session: %w", err)
	}

	r.logger.Info("Booking session created successfully",
		zap.String("session_id", session.ID),
		zap.String("user_id", session.UserID),
	)

	return nil
}

// GetByID retrieves a booking session by ID
func (r *BookingSessionRepository) GetByID(ctx context.Context, id string) (*models.BookingSession, error) {
	query := `SELECT * FROM booking_sessions WHERE id = $1`

	var session models.BookingSession
	err := r.db.GetContext(ctx, &session, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking session not found: %s", id)
		}
		r.logger.Error("Failed to get booking session by ID",
			zap.String("session_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get booking session: %w", err)
	}

	return &session, nil
}

// GetBySessionToken retrieves a booking session by session token
func (r *BookingSessionRepository) GetBySessionToken(ctx context.Context, sessionToken string) (*models.BookingSession, error) {
	query := `SELECT * FROM booking_sessions WHERE session_token = $1`

	var session models.BookingSession
	err := r.db.GetContext(ctx, &session, query, sessionToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking session not found: %s", sessionToken)
		}
		r.logger.Error("Failed to get booking session by token",
			zap.String("session_token", sessionToken),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get booking session: %w", err)
	}

	return &session, nil
}

// GetByUserID retrieves booking sessions for a user
func (r *BookingSessionRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*models.BookingSession, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM booking_sessions WHERE user_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking sessions: %w", err)
	}

	// Get sessions with pagination
	query := `
		SELECT * FROM booking_sessions 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var sessions []*models.BookingSession
	err = r.db.SelectContext(ctx, &sessions, query, userID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get booking sessions by user ID",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get booking sessions: %w", err)
	}

	return sessions, total, nil
}

// GetByEventID retrieves booking sessions for an event
func (r *BookingSessionRepository) GetByEventID(ctx context.Context, eventID string, page, limit int) ([]*models.BookingSession, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM booking_sessions WHERE event_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking sessions: %w", err)
	}

	// Get sessions with pagination
	query := `
		SELECT * FROM booking_sessions 
		WHERE event_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var sessions []*models.BookingSession
	err = r.db.SelectContext(ctx, &sessions, query, eventID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get booking sessions by event ID",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get booking sessions: %w", err)
	}

	return sessions, total, nil
}

// GetActiveSessions retrieves active booking sessions
func (r *BookingSessionRepository) GetActiveSessions(ctx context.Context, eventID string) ([]*models.BookingSession, error) {
	query := `
		SELECT * FROM booking_sessions 
		WHERE event_id = $1 AND status = 'active' AND expires_at > CURRENT_TIMESTAMP
		ORDER BY created_at ASC
	`

	var sessions []*models.BookingSession
	err := r.db.SelectContext(ctx, &sessions, query, eventID)
	if err != nil {
		r.logger.Error("Failed to get active booking sessions",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get active booking sessions: %w", err)
	}

	return sessions, nil
}

// Update updates a booking session
func (r *BookingSessionRepository) Update(ctx context.Context, session *models.BookingSession) error {
	query := `
		UPDATE booking_sessions SET
			status = $2, seat_count = $3, total_amount = $4,
			completed_at = $5, cancelled_at = $6, cancelled_reason = $7,
			metadata = $8, updated_by = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		session.ID, session.Status, session.SeatCount, session.TotalAmount,
		session.CompletedAt, session.CancelledAt, session.CancelledReason,
		session.Metadata, session.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to update booking session",
			zap.String("session_id", session.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update booking session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking session not found: %s", session.ID)
	}

	r.logger.Info("Booking session updated successfully",
		zap.String("session_id", session.ID),
	)

	return nil
}

// UpdateStatus updates booking session status
func (r *BookingSessionRepository) UpdateStatus(ctx context.Context, id, status, updatedBy string) error {
	var query string
	var args []interface{}

	switch status {
	case models.BookingSessionStatusCompleted:
		query = `
			UPDATE booking_sessions SET 
				status = $2, completed_at = CURRENT_TIMESTAMP, 
				updated_by = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{id, status, updatedBy}
	case models.BookingSessionStatusCancelled:
		query = `
			UPDATE booking_sessions SET 
				status = $2, cancelled_at = CURRENT_TIMESTAMP, 
				updated_by = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{id, status, updatedBy}
	default:
		query = `
			UPDATE booking_sessions SET 
				status = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{id, status, updatedBy}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to update booking session status",
			zap.String("session_id", id),
			zap.String("status", status),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update booking session status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking session not found: %s", id)
	}

	return nil
}

// Cancel cancels a booking session
func (r *BookingSessionRepository) Cancel(ctx context.Context, id, reason, cancelledBy string) error {
	query := `
		UPDATE booking_sessions SET 
			status = 'cancelled', cancelled_at = CURRENT_TIMESTAMP, 
			cancelled_reason = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, reason, cancelledBy)
	if err != nil {
		r.logger.Error("Failed to cancel booking session",
			zap.String("session_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to cancel booking session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking session not found: %s", id)
	}

	r.logger.Info("Booking session cancelled successfully",
		zap.String("session_id", id),
	)

	return nil
}

// Complete completes a booking session
func (r *BookingSessionRepository) Complete(ctx context.Context, id, completedBy string) error {
	query := `
		UPDATE booking_sessions SET 
			status = 'completed', completed_at = CURRENT_TIMESTAMP, 
			updated_by = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, completedBy)
	if err != nil {
		r.logger.Error("Failed to complete booking session",
			zap.String("session_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to complete booking session: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("booking session not found: %s", id)
	}

	r.logger.Info("Booking session completed successfully",
		zap.String("session_id", id),
	)

	return nil
}

// GetExpiredSessions retrieves expired booking sessions
func (r *BookingSessionRepository) GetExpiredSessions(ctx context.Context, before time.Time) ([]*models.BookingSession, error) {
	query := `
		SELECT * FROM booking_sessions 
		WHERE expires_at < $1 AND status = 'active'
		ORDER BY expires_at ASC
	`

	var sessions []*models.BookingSession
	err := r.db.SelectContext(ctx, &sessions, query, before)
	if err != nil {
		r.logger.Error("Failed to get expired booking sessions",
			zap.Time("before", before),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get expired booking sessions: %w", err)
	}

	return sessions, nil
}

// GetSessionsByStatus retrieves booking sessions by status
func (r *BookingSessionRepository) GetSessionsByStatus(ctx context.Context, status string, page, limit int) ([]*models.BookingSession, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM booking_sessions WHERE status = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count booking sessions: %w", err)
	}

	// Get sessions with pagination
	query := `
		SELECT * FROM booking_sessions 
		WHERE status = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var sessions []*models.BookingSession
	err = r.db.SelectContext(ctx, &sessions, query, status, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get booking sessions by status",
			zap.String("status", status),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get booking sessions: %w", err)
	}

	return sessions, total, nil
}

// GetSessionStats retrieves booking session statistics
func (r *BookingSessionRepository) GetSessionStats(ctx context.Context, eventID string) (map[string]int, error) {
	query := `
		SELECT 
			status,
			COUNT(*) as count
		FROM booking_sessions 
		WHERE event_id = $1
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		r.logger.Error("Failed to get session stats",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get session stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan session stats: %w", err)
		}
		stats[status] = count
	}

	return stats, nil
}
