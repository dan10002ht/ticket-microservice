package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
	"fmt"

	"github.com/google/uuid"
)

type EventSeatService struct {
	repo *repositories.EventSeatRepository
}

func NewEventSeatService(repo *repositories.EventSeatRepository) *EventSeatService {
	return &EventSeatService{repo: repo}
}

func (s *EventSeatService) DeleteSeat(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, publicID)
}

func (s *EventSeatService) GetSeat(ctx context.Context, publicID string) (*models.EventSeat, error) {
	return s.repo.GetByPublicID(ctx, publicID)
}

func (s *EventSeatService) CreateSeat(ctx context.Context, eventID, zoneID, seatNumber, rowNumber, coordinates, pricingCategory string, basePrice, finalPrice float64, currency string) (*models.EventSeat, error) {
	seat := &models.EventSeat{
		PublicID:        uuid.New().String(),
		EventID:         eventID,
		ZoneID:          zoneID,
		SeatNumber:      seatNumber,
		RowNumber:       rowNumber,
		Coordinates:     coordinates,
		Status:          "available",
		PricingCategory: pricingCategory,
		BasePrice:       basePrice,
		FinalPrice:      finalPrice,
		Currency:        currency,
		Version:         1,
	}

	if err := s.ValidateSeat(seat); err != nil {
		return nil, err
	}

	err := s.repo.Create(ctx, seat)
	if err != nil {
		return nil, err
	}

	return seat, nil
}

func (s *EventSeatService) UpdateSeat(ctx context.Context, seatID, seatNumber, rowNumber, coordinates, pricingCategory string, basePrice, finalPrice float64, currency string) (*models.EventSeat, error) {
	seat, err := s.repo.GetByPublicID(ctx, seatID)
	if err != nil {
		return nil, err
	}

	seat.SeatNumber = seatNumber
	seat.RowNumber = rowNumber
	seat.Coordinates = coordinates
	seat.PricingCategory = pricingCategory
	seat.BasePrice = basePrice
	seat.FinalPrice = finalPrice
	seat.Currency = currency

	if err := s.ValidateSeat(seat); err != nil {
		return nil, err
	}

	err = s.repo.Update(ctx, seat)
	if err != nil {
		return nil, err
	}

	return seat, nil
}

func (s *EventSeatService) ListSeatsByEvent(ctx context.Context, eventID, zoneID, status string, page, limit int32) ([]*models.EventSeat, int, error) {
	return s.repo.ListSeatsByEvent(ctx, eventID, zoneID, status, page, limit)
}

func (s *EventSeatService) ValidateSeat(seat *models.EventSeat) error {
	if seat.EventID == "" || seat.ZoneID == "" || seat.SeatNumber == "" {
		return fmt.Errorf("invalid seat data")
	}
	return nil
}
