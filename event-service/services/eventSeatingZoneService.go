package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
	"fmt"

	"github.com/google/uuid"
)

type EventSeatingZoneService struct {
	repo *repositories.EventSeatingZoneRepository
}

func NewEventSeatingZoneService(repo *repositories.EventSeatingZoneRepository) *EventSeatingZoneService {
	return &EventSeatingZoneService{repo: repo}
}

func (s *EventSeatingZoneService) CreateZone(ctx context.Context, eventID, name, zoneType, coordinates, color string) (*models.EventSeatingZone, error) {
	zone := &models.EventSeatingZone{
		PublicID:    uuid.New().String(),
		EventID:     eventID,
		Name:        name,
		ZoneType:    zoneType,
		Coordinates: coordinates,
		Color:       color,
		SeatCount:   0, // Will be updated when seats are added
	}

	if err := s.ValidateZone(zone); err != nil {
		return nil, err
	}

	err := s.repo.Create(ctx, zone)
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *EventSeatingZoneService) GetZone(ctx context.Context, zoneID string) (*models.EventSeatingZone, error) {
	return s.repo.GetByPublicID(ctx, zoneID)
}

func (s *EventSeatingZoneService) UpdateZone(ctx context.Context, zoneID, name, zoneType, coordinates, color string) (*models.EventSeatingZone, error) {
	zone, err := s.repo.GetByPublicID(ctx, zoneID)
	if err != nil {
		return nil, err
	}

	zone.Name = name
	zone.ZoneType = zoneType
	zone.Coordinates = coordinates
	zone.Color = color

	if err := s.ValidateZone(zone); err != nil {
		return nil, err
	}

	err = s.repo.Update(ctx, zone)
	if err != nil {
		return nil, err
	}

	return zone, nil
}

func (s *EventSeatingZoneService) DeleteZone(ctx context.Context, zoneID string) error {
	return s.repo.Delete(ctx, zoneID)
}

func (s *EventSeatingZoneService) ListZonesByEvent(ctx context.Context, eventID string) ([]*models.EventSeatingZone, error) {
	return s.repo.ListByEventID(ctx, eventID)
}

func (s *EventSeatingZoneService) ValidateZone(zone *models.EventSeatingZone) error {
	if zone.EventID == "" || zone.Name == "" {
		return fmt.Errorf("invalid zone data")
	}
	return nil
}
