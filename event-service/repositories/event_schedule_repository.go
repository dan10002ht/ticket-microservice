package repositories

import (
	"context"
	"event-service/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type EventScheduleRepository struct {
	db *sqlx.DB
}

func NewEventScheduleRepository(db *sqlx.DB) *EventScheduleRepository {
	return &EventScheduleRepository{db: db}
}

func (r *EventScheduleRepository) Create(ctx context.Context, sched *models.EventSchedule) error {
	query := `INSERT INTO event_schedules (public_id, event_id, schedule_type, start_date, end_date, recurrence_rule, is_active, created_at, updated_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9) RETURNING id`
	return r.db.QueryRowContext(ctx, query, sched.PublicID, sched.EventID, sched.ScheduleType, sched.StartDate, sched.EndDate, sched.RecurrenceRule, sched.IsActive, sched.CreatedAt, sched.UpdatedAt).Scan(&sched.ID)
}

func (r *EventScheduleRepository) GetByPublicID(ctx context.Context, publicID uuid.UUID) (*models.EventSchedule, error) {
	query := `SELECT id, public_id, event_id, schedule_type, start_date, end_date, recurrence_rule, is_active, created_at, updated_at FROM event_schedules WHERE public_id = $1`
	var sched models.EventSchedule
	err := r.db.QueryRowContext(ctx, query, publicID).Scan(&sched.ID, &sched.PublicID, &sched.EventID, &sched.ScheduleType, &sched.StartDate, &sched.EndDate, &sched.RecurrenceRule, &sched.IsActive, &sched.CreatedAt, &sched.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &sched, nil
}

func (r *EventScheduleRepository) Update(ctx context.Context, sched *models.EventSchedule) error {
	query := `UPDATE event_schedules SET event_id=$1, schedule_type=$2, start_date=$3, end_date=$4, recurrence_rule=$5, is_active=$6, updated_at=$7 WHERE public_id=$8`
	_, err := r.db.ExecContext(ctx, query, sched.EventID, sched.ScheduleType, sched.StartDate, sched.EndDate, sched.RecurrenceRule, sched.IsActive, sched.UpdatedAt, sched.PublicID)
	return err
}

func (r *EventScheduleRepository) Delete(ctx context.Context, publicID uuid.UUID) error {
	query := `DELETE FROM event_schedules WHERE public_id = $1`
	_, err := r.db.ExecContext(ctx, query, publicID)
	return err
}

func (r *EventScheduleRepository) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSchedule, error) {
	query := `SELECT id, public_id, event_id, schedule_type, start_date, end_date, recurrence_rule, is_active, created_at, updated_at FROM event_schedules WHERE event_id = $1`
	rows, err := r.db.QueryContext(ctx, query, eventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var scheds []*models.EventSchedule
	for rows.Next() {
		var sched models.EventSchedule
		err := rows.Scan(&sched.ID, &sched.PublicID, &sched.EventID, &sched.ScheduleType, &sched.StartDate, &sched.EndDate, &sched.RecurrenceRule, &sched.IsActive, &sched.CreatedAt, &sched.UpdatedAt)
		if err != nil {
			return nil, err
		}
		scheds = append(scheds, &sched)
	}
	return scheds, nil
} 