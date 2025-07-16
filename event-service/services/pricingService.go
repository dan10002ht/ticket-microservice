package services

import (
	"context"
	"event-service/models"
	"event-service/repositories"
	"fmt"

	"github.com/google/uuid"
)

type PricingService struct {
	repo *repositories.EventPricingRepository
}

func NewPricingService(repo *repositories.EventPricingRepository) *PricingService {
	return &PricingService{repo: repo}
}

func (s *PricingService) CreatePricing(ctx context.Context, pricing *models.EventPricing) error {
	if err := s.ValidatePricing(pricing); err != nil {
		return err
	}
	return s.repo.Create(ctx, pricing)
}

func (s *PricingService) GetPricing(ctx context.Context, publicID string) (*models.EventPricing, error) {
	return s.repo.GetByPublicID(ctx, uuid.MustParse(publicID))
}

func (s *PricingService) UpdatePricing(ctx context.Context, pricing *models.EventPricing) error {
	if err := s.ValidatePricing(pricing); err != nil {
		return err
	}
	return s.repo.Update(ctx, pricing)
}

func (s *PricingService) DeletePricing(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, uuid.MustParse(publicID))
}

func (s *PricingService) ListByEventID(ctx context.Context, eventID int64) ([]*models.EventPricing, error) {
	return s.repo.ListByEventID(ctx, eventID)
}

func (s *PricingService) ValidatePricing(pricing *models.EventPricing) error {
	if pricing.EventID == 0 || pricing.ZoneID == 0 || pricing.Price < 0 {
		return fmt.Errorf("invalid pricing data")
	}
	return nil
}

func (s *PricingService) ApplyDiscount(ctx context.Context, eventID int64, basePrice float64, discountCode string, userID string) (finalPrice float64, discountAmount float64, reason string, err error) {
	// Giả lập: lấy rule từ DB, kiểm tra code, tính discount
	// Thực tế: lấy rule từ pricingRepo, kiểm tra code hợp lệ, tính toán
	if discountCode == "" {
		return basePrice, 0, "No discount", nil
	}
	// Giả lập rule: code "PROMO10" giảm 10%, "PROMO50" giảm 50k
	if discountCode == "PROMO10" {
		discountAmount = basePrice * 0.1
		finalPrice = basePrice - discountAmount
		reason = "10% promotion"
		return
	}
	if discountCode == "PROMO50" {
		discountAmount = 50000
		finalPrice = basePrice - discountAmount
		if finalPrice < 0 { finalPrice = 0 }
		reason = "50k promotion"
		return
	}
	// Nếu không hợp lệ
	return basePrice, 0, "Invalid discount code", nil
}

func (s *PricingService) CalculatePrice(ctx context.Context, eventID int64, zoneID int64, quantity int, discountCode string, userID string) (basePrice, finalPrice, discountAmount float64, reason string, err error) {
	// Lấy base price từ pricingRepo
	pricings, err := s.repo.ListByEventID(ctx, eventID)
	if err != nil || len(pricings) == 0 {
		return 0, 0, 0, "No pricing found", err
	}
	var price float64
	for _, p := range pricings {
		if p.ZoneID == zoneID {
			price = p.Price
			break
		}
	}
	if price == 0 {
		return 0, 0, 0, "No price for zone", nil
	}
	basePrice = price * float64(quantity)
	finalPrice, discountAmount, reason, err = s.ApplyDiscount(ctx, eventID, basePrice, discountCode, userID)
	return
} 