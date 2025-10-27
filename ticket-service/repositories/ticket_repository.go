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

// TicketRepository handles database operations for tickets
type TicketRepository struct {
	db     *sqlx.DB
	logger *zap.Logger
}

// NewTicketRepository creates a new ticket repository
func NewTicketRepository(db *sqlx.DB, logger *zap.Logger) *TicketRepository {
	return &TicketRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new ticket
func (r *TicketRepository) Create(ctx context.Context, ticket *models.Ticket) error {
	query := `
		INSERT INTO tickets (
			id, event_id, seat_id, zone_id, user_id, booking_session_id,
			ticket_number, ticket_type, pricing_category, base_price, final_price,
			currency, discount_amount, discount_reason, status, payment_status,
			payment_method, payment_reference, qr_code, barcode, valid_from,
			valid_until, metadata, created_by, updated_by
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16,
			$17, $18, $19, $20, $21, $22, $23, $24, $25
		)
	`

	// Generate UUID if not provided
	if ticket.ID == "" {
		ticket.ID = uuid.New().String()
	}

	_, err := r.db.ExecContext(ctx, query,
		ticket.ID, ticket.EventID, ticket.SeatID, ticket.ZoneID, ticket.UserID,
		ticket.BookingSessionID, ticket.TicketNumber, ticket.TicketType,
		ticket.PricingCategory, ticket.BasePrice, ticket.FinalPrice,
		ticket.Currency, ticket.DiscountAmount, ticket.DiscountReason,
		ticket.Status, ticket.PaymentStatus, ticket.PaymentMethod,
		ticket.PaymentReference, ticket.QRCode, ticket.Barcode,
		ticket.ValidFrom, ticket.ValidUntil, ticket.Metadata,
		ticket.CreatedBy, ticket.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to create ticket",
			zap.String("ticket_id", ticket.ID),
			zap.String("event_id", ticket.EventID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create ticket: %w", err)
	}

	r.logger.Info("Ticket created successfully",
		zap.String("ticket_id", ticket.ID),
		zap.String("event_id", ticket.EventID),
	)

	return nil
}

// GetByID retrieves a ticket by ID
func (r *TicketRepository) GetByID(ctx context.Context, id string) (*models.Ticket, error) {
	query := `SELECT * FROM tickets WHERE id = $1`

	var ticket models.Ticket
	err := r.db.GetContext(ctx, &ticket, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket not found: %s", id)
		}
		r.logger.Error("Failed to get ticket by ID",
			zap.String("ticket_id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return &ticket, nil
}

// GetByTicketNumber retrieves a ticket by ticket number
func (r *TicketRepository) GetByTicketNumber(ctx context.Context, ticketNumber string) (*models.Ticket, error) {
	query := `SELECT * FROM tickets WHERE ticket_number = $1`

	var ticket models.Ticket
	err := r.db.GetContext(ctx, &ticket, query, ticketNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("ticket not found: %s", ticketNumber)
		}
		r.logger.Error("Failed to get ticket by number",
			zap.String("ticket_number", ticketNumber),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	return &ticket, nil
}

// GetByUserID retrieves tickets for a user
func (r *TicketRepository) GetByUserID(ctx context.Context, userID string, page, limit int) ([]*models.Ticket, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM tickets WHERE user_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, userID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	// Get tickets with pagination
	query := `
		SELECT * FROM tickets 
		WHERE user_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var tickets []*models.Ticket
	err = r.db.SelectContext(ctx, &tickets, query, userID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get tickets by user ID",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get tickets: %w", err)
	}

	return tickets, total, nil
}

// GetByEventID retrieves tickets for an event
func (r *TicketRepository) GetByEventID(ctx context.Context, eventID string, page, limit int) ([]*models.Ticket, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM tickets WHERE event_id = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, eventID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	// Get tickets with pagination
	query := `
		SELECT * FROM tickets 
		WHERE event_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var tickets []*models.Ticket
	err = r.db.SelectContext(ctx, &tickets, query, eventID, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get tickets by event ID",
			zap.String("event_id", eventID),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get tickets: %w", err)
	}

	return tickets, total, nil
}

// GetByBookingSessionID retrieves tickets for a booking session
func (r *TicketRepository) GetByBookingSessionID(ctx context.Context, bookingSessionID string) ([]*models.Ticket, error) {
	query := `SELECT * FROM tickets WHERE booking_session_id = $1 ORDER BY created_at ASC`

	var tickets []*models.Ticket
	err := r.db.SelectContext(ctx, &tickets, query, bookingSessionID)
	if err != nil {
		r.logger.Error("Failed to get tickets by booking session ID",
			zap.String("booking_session_id", bookingSessionID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get tickets: %w", err)
	}

	return tickets, nil
}

// Update updates a ticket
func (r *TicketRepository) Update(ctx context.Context, ticket *models.Ticket) error {
	query := `
		UPDATE tickets SET
			status = $2, payment_status = $3, payment_method = $4,
			payment_reference = $5, qr_code = $6, barcode = $7,
			used_at = $8, cancelled_at = $9, cancelled_reason = $10,
			refunded_at = $11, refunded_amount = $12, metadata = $13,
			updated_by = $14, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		ticket.ID, ticket.Status, ticket.PaymentStatus, ticket.PaymentMethod,
		ticket.PaymentReference, ticket.QRCode, ticket.Barcode,
		ticket.UsedAt, ticket.CancelledAt, ticket.CancelledReason,
		ticket.RefundedAt, ticket.RefundedAmount, ticket.Metadata,
		ticket.UpdatedBy,
	)

	if err != nil {
		r.logger.Error("Failed to update ticket",
			zap.String("ticket_id", ticket.ID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found: %s", ticket.ID)
	}

	r.logger.Info("Ticket updated successfully",
		zap.String("ticket_id", ticket.ID),
	)

	return nil
}

// UpdateStatus updates ticket status
func (r *TicketRepository) UpdateStatus(ctx context.Context, id, status, updatedBy string) error {
	query := `UPDATE tickets SET status = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id, status, updatedBy)
	if err != nil {
		r.logger.Error("Failed to update ticket status",
			zap.String("ticket_id", id),
			zap.String("status", status),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update ticket status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found: %s", id)
	}

	return nil
}

// UpdatePaymentStatus updates ticket payment status
func (r *TicketRepository) UpdatePaymentStatus(ctx context.Context, id, paymentStatus, paymentMethod, paymentReference, updatedBy string) error {
	query := `
		UPDATE tickets SET 
			payment_status = $2, payment_method = $3, payment_reference = $4,
			updated_by = $5, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, paymentStatus, paymentMethod, paymentReference, updatedBy)
	if err != nil {
		r.logger.Error("Failed to update ticket payment status",
			zap.String("ticket_id", id),
			zap.String("payment_status", paymentStatus),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update ticket payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found: %s", id)
	}

	return nil
}

// Delete deletes a ticket (soft delete by updating status)
func (r *TicketRepository) Delete(ctx context.Context, id, reason, deletedBy string) error {
	query := `
		UPDATE tickets SET 
			status = 'cancelled', cancelled_at = CURRENT_TIMESTAMP, 
			cancelled_reason = $2, updated_by = $3, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, reason, deletedBy)
	if err != nil {
		r.logger.Error("Failed to delete ticket",
			zap.String("ticket_id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete ticket: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("ticket not found: %s", id)
	}

	r.logger.Info("Ticket deleted successfully",
		zap.String("ticket_id", id),
	)

	return nil
}

// Search searches tickets with filters
func (r *TicketRepository) Search(ctx context.Context, filters map[string]interface{}, page, limit int) ([]*models.Ticket, int, error) {
	whereClause := "WHERE 1=1"
	args := []interface{}{}
	argIndex := 1

	// Build where clause dynamically
	if eventID, ok := filters["event_id"].(string); ok && eventID != "" {
		whereClause += fmt.Sprintf(" AND event_id = $%d", argIndex)
		args = append(args, eventID)
		argIndex++
	}

	if userID, ok := filters["user_id"].(string); ok && userID != "" {
		whereClause += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, userID)
		argIndex++
	}

	if status, ok := filters["status"].(string); ok && status != "" {
		whereClause += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	if paymentStatus, ok := filters["payment_status"].(string); ok && paymentStatus != "" {
		whereClause += fmt.Sprintf(" AND payment_status = $%d", argIndex)
		args = append(args, paymentStatus)
		argIndex++
	}

	if ticketType, ok := filters["ticket_type"].(string); ok && ticketType != "" {
		whereClause += fmt.Sprintf(" AND ticket_type = $%d", argIndex)
		args = append(args, ticketType)
		argIndex++
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM tickets %s", whereClause)
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	// Get tickets with pagination
	query := fmt.Sprintf(`
		SELECT * FROM tickets %s 
		ORDER BY created_at DESC 
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	offset := (page - 1) * limit
	args = append(args, limit, offset)

	var tickets []*models.Ticket
	err = r.db.SelectContext(ctx, &tickets, query, args...)
	if err != nil {
		r.logger.Error("Failed to search tickets",
			zap.Any("filters", filters),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to search tickets: %w", err)
	}

	return tickets, total, nil
}

// GetExpiredTickets retrieves expired tickets
func (r *TicketRepository) GetExpiredTickets(ctx context.Context, before time.Time) ([]*models.Ticket, error) {
	query := `
		SELECT * FROM tickets 
		WHERE valid_until < $1 AND status IN ('pending', 'confirmed')
		ORDER BY valid_until ASC
	`

	var tickets []*models.Ticket
	err := r.db.SelectContext(ctx, &tickets, query, before)
	if err != nil {
		r.logger.Error("Failed to get expired tickets",
			zap.Time("before", before),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get expired tickets: %w", err)
	}

	return tickets, nil
}

// GetTicketsByStatus retrieves tickets by status
func (r *TicketRepository) GetTicketsByStatus(ctx context.Context, status string, page, limit int) ([]*models.Ticket, int, error) {
	// Count total
	countQuery := `SELECT COUNT(*) FROM tickets WHERE status = $1`
	var total int
	err := r.db.GetContext(ctx, &total, countQuery, status)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tickets: %w", err)
	}

	// Get tickets with pagination
	query := `
		SELECT * FROM tickets 
		WHERE status = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3
	`

	offset := (page - 1) * limit
	var tickets []*models.Ticket
	err = r.db.SelectContext(ctx, &tickets, query, status, limit, offset)
	if err != nil {
		r.logger.Error("Failed to get tickets by status",
			zap.String("status", status),
			zap.Error(err),
		)
		return nil, 0, fmt.Errorf("failed to get tickets: %w", err)
	}

	return tickets, total, nil
}
