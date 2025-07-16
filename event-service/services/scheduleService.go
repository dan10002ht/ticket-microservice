package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
	"fmt"

	"github.com/google/uuid"
)

type ScheduleService struct {
	repo *repositories.EventScheduleRepository
}

func NewScheduleService(repo *repositories.EventScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, sched *models.EventSchedule) error {
	if err := s.ValidateSchedule(sched); err != nil {
		return err
	}
	return s.repo.Create(ctx, sched)
}

func (s *ScheduleService) GetSchedule(ctx context.Context, publicID string) (*models.EventSchedule, error) {
	return s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
}

func (s *ScheduleService) UpdateSchedule(ctx context.Context, sched *models.EventSchedule) error {
	if err := s.ValidateSchedule(sched); err != nil {
		return err
	}
	return s.repo.Update(ctx, sched)
}

func (s *ScheduleService) DeleteSchedule(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, uuid.MustParse(publicID))
}

func (s *ScheduleService) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventSchedule, error) {
	return s.repo.ListByEventID(ctx, eventID)
}

func (s *ScheduleService) ValidateSchedule(sched *models.EventSchedule) error {
	if sched.EventID == 0 || sched.ScheduleType == "" {
		return fmt.Errorf("invalid schedule data")
	}
	return nil
} 