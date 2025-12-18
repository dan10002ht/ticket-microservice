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

func (s *PricingService) CreatePricing(ctx context.Context, eventID, zoneID, pricingCategory string, basePrice float64, currency, pricingRules, discountRules, validFrom, validUntil, createdBy string) (*models.EventPricing, error) {
	pricing := &models.EventPricing{
		PublicID:        uuid.New().String(),
		EventID:         eventID,
		ZoneID:          zoneID,
		PricingCategory: pricingCategory,
		BasePrice:       basePrice,
		Currency:        currency,
		PricingRules:    pricingRules,
		DiscountRules:   discountRules,
		IsActive:        true,
		ValidFrom:       validFrom,
		ValidUntil:      validUntil,
		CreatedBy:       createdBy,
	}

	if err := s.ValidatePricing(pricing); err != nil {
		return nil, err
	}

	err := s.repo.Create(ctx, pricing)
	if err != nil {
		return nil, err
	}

	return pricing, nil
}

func (s *PricingService) GetPricing(ctx context.Context, publicID string) (*models.EventPricing, error) {
	return s.repo.GetByPublicID(ctx, publicID)
}

func (s *PricingService) UpdatePricing(ctx context.Context, pricingID string, basePrice float64, currency, pricingRules, discountRules string, isActive bool, validFrom, validUntil, updatedBy string) (*models.EventPricing, error) {
	pricing, err := s.repo.GetByPublicID(ctx, pricingID)
	if err != nil {
		return nil, err
	}

	pricing.BasePrice = basePrice
	pricing.Currency = currency
	pricing.PricingRules = pricingRules
	pricing.DiscountRules = discountRules
	pricing.IsActive = isActive
	pricing.ValidFrom = validFrom
	pricing.ValidUntil = validUntil
	pricing.UpdatedBy = updatedBy

	if err := s.ValidatePricing(pricing); err != nil {
		return nil, err
	}

	err = s.repo.Update(ctx, pricing)
	if err != nil {
		return nil, err
	}

	return pricing, nil
}

func (s *PricingService) DeletePricing(ctx context.Context, publicID string) error {
	return s.repo.Delete(ctx, publicID)
}

func (s *PricingService) ListPricing(ctx context.Context, eventID string, isActive bool, page, limit int32) ([]*models.EventPricing, int, error) {
	return s.repo.ListPricing(ctx, eventID, isActive, page, limit)
}

func (s *PricingService) ValidatePricing(pricing *models.EventPricing) error {
	if pricing.EventID == "" || pricing.ZoneID == "" || pricing.BasePrice < 0 {
		return fmt.Errorf("invalid pricing data")
	}
	return nil
}

func (s *PricingService) ApplyDiscount(ctx context.Context, eventID string, originalPrice float64, discountCode string, userID string) (discountAmount float64, finalPrice float64, reason string, isValid bool, err error) {
	if discountCode == "" {
		return 0, originalPrice, "No discount applied", false, nil
	}

	// Simulated discount logic - in production, fetch from DB
	switch discountCode {
	case "PROMO10":
		discountAmount = originalPrice * 0.1
		finalPrice = originalPrice - discountAmount
		reason = "10% promotion discount"
		isValid = true
	case "PROMO50":
		discountAmount = 50000
		finalPrice = originalPrice - discountAmount
		if finalPrice < 0 {
			finalPrice = 0
		}
		reason = "50k fixed discount"
		isValid = true
	default:
		return 0, originalPrice, "Invalid discount code", false, nil
	}

	return discountAmount, finalPrice, reason, isValid, nil
}

func (s *PricingService) CalculatePrice(ctx context.Context, eventID, zoneID, pricingCategory string, quantity int32, discountCode, userID string) (basePrice, finalPrice, discountAmount float64, discountReason, currency, pricingDetails string, err error) {
	// Get pricing for the zone
	pricings, err := s.repo.GetPricingByZone(ctx, eventID, zoneID)
	if err != nil {
		return 0, 0, 0, "", "", "", err
	}

	if len(pricings) == 0 {
		return 0, 0, 0, "No pricing found for zone", "", "", fmt.Errorf("no pricing found for zone")
	}

	// Find matching pricing category or use first active pricing
	var selectedPricing *models.EventPricing
	for _, p := range pricings {
		if p.IsActive {
			if pricingCategory == "" || p.PricingCategory == pricingCategory {
				selectedPricing = p
				break
			}
		}
	}

	if selectedPricing == nil {
		return 0, 0, 0, "No active pricing found", "", "", fmt.Errorf("no active pricing found")
	}

	basePrice = selectedPricing.BasePrice * float64(quantity)
	currency = selectedPricing.Currency

	// Apply discount if provided
	if discountCode != "" {
		var isValid bool
		discountAmount, finalPrice, discountReason, isValid, err = s.ApplyDiscount(ctx, eventID, basePrice, discountCode, userID)
		if err != nil {
			return 0, 0, 0, "", "", "", err
		}
		if !isValid {
			finalPrice = basePrice
			discountAmount = 0
		}
	} else {
		finalPrice = basePrice
		discountReason = "No discount applied"
	}

	pricingDetails = fmt.Sprintf(`{"pricing_id":"%s","category":"%s","unit_price":%f,"quantity":%d}`,
		selectedPricing.PublicID, selectedPricing.PricingCategory, selectedPricing.BasePrice, quantity)

	return basePrice, finalPrice, discountAmount, discountReason, currency, pricingDetails, nil
}

func (s *PricingService) GetPricingByEvent(ctx context.Context, eventID string, isActive bool) ([]*models.EventPricing, error) {
	return s.repo.GetPricingByEvent(ctx, eventID, isActive)
}

func (s *PricingService) GetPricingByZone(ctx context.Context, eventID, zoneID string, isActive bool) ([]*models.EventPricing, error) {
	if isActive {
		return s.repo.GetActivePricingByZone(ctx, eventID, zoneID)
	}
	return s.repo.GetPricingByZone(ctx, eventID, zoneID)
}
