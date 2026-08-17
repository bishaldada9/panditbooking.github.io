package dto

import "github.com/google/uuid"

type CreateRitualCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type CreateRitualRequest struct {
	CategoryID       uuid.UUID `json:"category_id" validate:"required"`
	Name             string    `json:"name" validate:"required"`
	Description      string    `json:"description"`
	Duration         string    `json:"duration" validate:"required"`
	BasePrice        float64   `json:"base_price" validate:"required,min=0"`
	RequiredItems    string    `json:"required_items"`
	Procedure        string    `json:"procedure"`
	PanditCommission float64   `json:"pandit_commission"`
}

type RitualCategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	IsActive    bool      `json:"is_active"`
}

type RitualResponse struct {
	ID               uuid.UUID `json:"id"`
	CategoryID       uuid.UUID `json:"category_id"`
	CategoryName     string    `json:"category_name"`
	Name             string    `json:"name"`
	Slug             string    `json:"slug"`
	Description      string    `json:"description"`
	Duration         string    `json:"duration"`
	BasePrice        float64   `json:"base_price"`
	RequiredItems    string    `json:"required_items"`
	Procedure        string    `json:"procedure"`
	PanditCommission float64   `json:"pandit_commission"`
}
