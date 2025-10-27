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

// SeatReservationRepository handles database operations for seat reservations
type SeatReservationRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewSeatReservationRepository creates a new seat reservation repository
func NewSeatReservationRepository(db *sqlx.DB, logger *zap.Logger) *SeatReservationRepository {
	return &SeatReservationRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new seat reservation
func (r *SeatReservationRepository) Create(ctx context.Context, reservation *models.SeatReservation) error {
	query := `
		INSERT INTO seat_reservations (
			id, booking_session_id, event_id, seat_id, zone_id,
			reservation_token, status, reserved_at, expires_at,
			pricing_category, base_price, final_price, currency,
			metadata, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16
		)
	`

	// Generate UUID if not provided
	if reservation.ID == "" {
		reservation.ID = uuid.New().String()
	}

	_, err := r.db.ExecContext(ctx, query,
		reservation.ID, reservation.BookingSessionID, reservation.EventID,
		reservation.SeatID, reservation.ZoneID, reservation.ReservationToken,
		reservation.Status, reservation.ReservedAt, reservation.ExpiresAt,
		reservation.PricingCategory, reservation.BasePrice, reservation.FinalPrice,
		reservation.Currency, reservation.Metadata, reservation.CreatedBy,
		reservation.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to create seat reservation",
			zap.String("reservation_id", reservation.ID),
			zap.String("seat_id", reservation.SeatID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create seat reservation: %w", err)
	}

	r.logger.Info("Seat reservation created successfully",
		zap.String("reservation_id", reservation.ID),
		zap.String("seat_id", reservation.SeatID),
	)

	return nil
}

// GetByID retrieves a seat reservation by ID
func (r *SeatReservationRepository) GetByID(ctx context.Context, id string) (*models.SeatReservation, error) {
	query := `SELECT * FROM seat_reservations WHERE id = $1`

	var reservation models.SeatReservation
	err := r.db.GetContext(ctx, &reservation, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("seat reservation not found: %s", id)
		}
		r.logger.Error("Failed to get seat reservation by ID",
			zap.String("reservation_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get seat reservation: %w", err)
	}

	return &reservation, nil
}

// GetByReservationToken retrieves a seat reservation by reservation token
func (r *SeatReservationRepository) GetByReservationToken(ctx context.Context, token string) (*models.SeatReservation, error) {
	query := `SELECT * FROM seat_reservations WHERE reservation_token = $1`

	var reservation models.SeatReservation
	err := r.db.GetContext(ctx, &reservation, query, token)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("seat reservation not found: %s", token)
		}
		r.logger.Error("Failed to get seat reservation by token",
			zap.String("reservation_token", token),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get seat reservation: %w", err)
	}

	return &reservation, nil
}

// GetByBookingSessionID retrieves seat reservations for a booking session
func (r *SeatReservationRepository) GetByBookingSessionID(ctx context.Context, bookingSessionID string) ([]*models.SeatReservation, error) {
	query := `SELECT * FROM seat_reservations WHERE booking_session_id = $1 ORDER BY created_at ASC`

	var reservations []*models.SeatReservation
	err := r.db.SelectContext(ctx, &reservations, query, bookingSessionID)
	if err != nil {
		r.logger.Error("Failed to get seat reservations by booking session ID",
			zap.String("booking_session_id", bookingSessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	return reservations, nil
}

// GetByEventID retrieves seat reservations for an event
func (r *SeatReservationRepository) GetByEventID(ctx context.Context, eventID string, page, limit int) ([]*models.SeatReservation, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM seat_reservations WHERE event_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count seat reservations: %w", err)
	}

	// Get reservations with pagination
	query := `
		SELECT * FROM seat_reservations 
		WHERE event_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var reservations []*models.SeatReservation
	err = r.db.SelectContext(ctx, &reservations, query, eventID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get seat reservations by event ID",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	return reservations, total, nil
}

// GetBySeatID retrieves seat reservations for a specific seat
func (r *SeatReservationRepository) GetBySeatID(ctx context.Context, seatID string) ([]*models.SeatReservation, error) {
	query := `
		SELECT * FROM seat_reservations 
		WHERE seat_id = $1 
		ORDER BY created_at DESC
	`

	var reservations []*models.SeatReservation
	err := r.db.SelectContext(ctx, &reservations, query, seatID)
	if err != nil {
		r.logger.Error("Failed to get seat reservations by seat ID",
			zap.String("seat_id", seatID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	return reservations, nil
}

// GetActiveReservations retrieves active seat reservations
func (r *SeatReservationRepository) GetActiveReservations(ctx context.Context, eventID string) ([]*models.SeatReservation, error) {
	query := `
		SELECT * FROM seat_reservations 
		WHERE event_id = $1 AND status = 'reserved' AND expires_at > CURRENT_TIMESTAMP
		ORDER BY reserved_at ASC
	`

	var reservations []*models.SeatReservation
	err := r.db.SelectContext(ctx, &reservations, query, eventID)
	if err != nil {
		r.logger.Error("Failed to get active seat reservations",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get active seat reservations: %w", err)
	}

	return reservations, nil
}

// GetReservationsBySeatIDs retrieves reservations for multiple seats
func (r *SeatReservationRepository) GetReservationsBySeatIDs(ctx context.Context, seatIDs []string) ([]*models.SeatReservation, error) {
	if len(seatIDs) == 0 {
		return []*models.SeatReservation{}, nil
	}

	// Build dynamic query
	query := `SELECT * FROM seat_reservations WHERE seat_id = ANY($1) ORDER BY created_at DESC`

	var reservations []*models.SeatReservation
	err := r.db.SelectContext(ctx, &reservations, query, seatIDs)
	if err != nil {
		r.logger.Error("Failed to get seat reservations by seat IDs",
			zap.Strings("seat_ids", seatIDs),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	return reservations, nil
}

// Update updates a seat reservation
func (r *SeatReservationRepository) Update(ctx context.Context, reservation *models.SeatReservation) error {
	query := `
		UPDATE seat_reservations SET
			status = $2, released_at = $3, released_reason = $4,
			metadata = $5, updated_by = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		reservation.ID, reservation.Status, reservation.ReleasedAt,
		reservation.ReleasedReason, reservation.Metadata, reservation.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to update seat reservation",
			zap.String("reservation_id", reservation.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update seat reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("seat reservation not found: %s", reservation.ID)
	}

	r.logger.Info("Seat reservation updated successfully",
		zap.String("reservation_id", reservation.ID),
	)

	return nil
}

// UpdateStatus updates seat reservation status
func (r *SeatReservationRepository) UpdateStatus(ctx context.Context, id, status, updatedBy string) error {
	var query string
	var args []interface{}

	switch status {
	case models.ReservationStatusReleased:
		query = `
			UPDATE seat_reservations SET 
				status = $2, released_at = CURRENT_TIMESTAMP, 
				updated_by = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{id, status, updatedBy}
	case models.ReservationStatusConfirmed:
		query = `
			UPDATE seat_reservations SET 
				status = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{id, status, updatedBy}
	default:
		query = `
			UPDATE seat_reservations SET 
				status = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
			WHERE id = $1
		`
		args = []interface{}{id, status, updatedBy}
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to update seat reservation status",
			zap.String("reservation_id", id),
			zap.String("status", status),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update seat reservation status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("seat reservation not found: %s", id)
	}

	return nil
}

// Release releases a seat reservation
func (r *SeatReservationRepository) Release(ctx context.Context, id, reason, releasedBy string) error {
	query := `
		UPDATE seat_reservations SET 
			status = 'released', released_at = CURRENT_TIMESTAMP, 
			released_reason = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, reason, releasedBy)
	if err != nil {
		r.logger.Error("Failed to release seat reservation",
			zap.String("reservation_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to release seat reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("seat reservation not found: %s", id)
	}

	r.logger.Info("Seat reservation released successfully",
		zap.String("reservation_id", id),
	)

	return nil
}

// Confirm confirms a seat reservation
func (r *SeatReservationRepository) Confirm(ctx context.Context, id, confirmedBy string) error {
	query := `
		UPDATE seat_reservations SET 
			status = 'confirmed', updated_by = $2, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, confirmedBy)
	if err != nil {
		r.logger.Error("Failed to confirm seat reservation",
			zap.String("reservation_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to confirm seat reservation: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("seat reservation not found: %s", id)
	}

	r.logger.Info("Seat reservation confirmed successfully",
		zap.String("reservation_id", id),
	)

	return nil
}

// ReleaseByBookingSession releases all reservations for a booking session
func (r *SeatReservationRepository) ReleaseByBookingSession(ctx context.Context, bookingSessionID, reason, releasedBy string) error {
	query := `
		UPDATE seat_reservations SET 
			status = 'released', released_at = CURRENT_TIMESTAMP, 
			released_reason = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE booking_session_id = $1 AND status = 'reserved'
	`

	result, err := r.db.ExecContext(ctx, query, bookingSessionID, reason, releasedBy)
	if err != nil {
		r.logger.Error("Failed to release seat reservations by booking session",
			zap.String("booking_session_id", bookingSessionID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to release seat reservations: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.logger.Info("Seat reservations released by booking session",
		zap.String("booking_session_id", bookingSessionID),
		zap.Int64("rows_affected", rowsAffected),
	)

	return nil
}

// GetExpiredReservations retrieves expired seat reservations
func (r *SeatReservationRepository) GetExpiredReservations(ctx context.Context, before time.Time) ([]*models.SeatReservation, error) {
	query := `
		SELECT * FROM seat_reservations 
		WHERE expires_at < $1 AND status = 'reserved'
		ORDER BY expires_at ASC
	`

	var reservations []*models.SeatReservation
	err := r.db.SelectContext(ctx, &reservations, query, before)
	if err != nil {
		r.logger.Error("Failed to get expired seat reservations",
			zap.Time("before", before),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get expired seat reservations: %w", err)
	}

	return reservations, nil
}

// GetReservationsByStatus retrieves seat reservations by status
func (r *SeatReservationRepository) GetReservationsByStatus(ctx context.Context, status string, page, limit int) ([]*models.SeatReservation, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM seat_reservations WHERE status = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count seat reservations: %w", err)
	}

	// Get reservations with pagination
	query := `
		SELECT * FROM seat_reservations 
		WHERE status = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var reservations []*models.SeatReservation
	err = r.db.SelectContext(ctx, &reservations, query, status, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get seat reservations by status",
			zap.String("status", status),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get seat reservations: %w", err)
	}

	return reservations, total, nil
}

// GetReservationStats retrieves seat reservation statistics
func (r *SeatReservationRepository) GetReservationStats(ctx context.Context, eventID string) (map[string]int, error) {
	query := `
		SELECT 
			status,
			COUNT(*) as count
		FROM seat_reservations 
		WHERE event_id = $1
		GROUP BY status
	`

	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		r.logger.Error("Failed to get reservation stats",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get reservation stats: %w", err)
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, fmt.Errorf("failed to scan reservation stats: %w", err)
		}
		stats[status] = count
	}

	return stats, nil
}

// CheckSeatAvailability checks if a seat is available for reservation
func (r *SeatReservationRepository) CheckSeatAvailability(ctx context.Context, seatID string) (bool, error) {
	query := `
		SELECT COUNT(*) FROM seat_reservations 
		WHERE seat_id = $1 AND status = 'reserved' AND expires_at > CURRENT_TIMESTAMP
	`

	var count int
	err := r.db.GetContext(ctx, &count, query, seatID)
	if err != nil {
		r.logger.Error("Failed to check seat availability",
			zap.String("seat_id", seatID),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check seat availability: %w", err)
	}

	return count == 0, nil
}
