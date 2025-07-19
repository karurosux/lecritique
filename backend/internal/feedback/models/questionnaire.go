package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	menuModels "lecritique/internal/menu/models"
	restaurantModels "lecritique/internal/restaurant/models"
	sharedModels "lecritique/internal/shared/models"
)

type Questionnaire struct {
	sharedModels.BaseModel
	RestaurantID uuid.UUID          `gorm:"not null" json:"restaurant_id"`
	Restaurant   restaurantModels.Restaurant         `json:"restaurant,omitempty"`
	DishID       *uuid.UUID         `json:"dish_id"`
	Dish         *menuModels.Dish              `json:"dish,omitempty"`
	Name         string             `gorm:"not null" json:"name"`
	Description  string             `json:"description"`
	IsDefault    bool               `gorm:"default:false" json:"is_default"`
	IsActive     bool               `gorm:"default:true" json:"is_active"`
	Questions    []Question         `json:"questions,omitempty"`
}

type Question struct {
	sharedModels.BaseModel
	DishID       uuid.UUID        `gorm:"not null;index" json:"dish_id"`
	Dish         *menuModels.Dish `json:"dish,omitempty"`
	Text         string           `gorm:"not null" json:"text"`
	Type         QuestionType     `gorm:"not null" json:"type"`
	IsRequired   bool             `gorm:"default:true" json:"is_required"`
	DisplayOrder int              `gorm:"default:0" json:"display_order"`
	Options      pq.StringArray   `gorm:"type:text[]" json:"options" swaggertype:"array,string"`
	MinValue     *int             `json:"min_value"`
	MaxValue     *int             `json:"max_value"`
	MinLabel     string           `json:"min_label"`
	MaxLabel     string           `json:"max_label"`
}

type QuestionType string

const (
	QuestionTypeRating       QuestionType = "rating"       // 1-5 stars
	QuestionTypeScale        QuestionType = "scale"        // 1-10 scale
	QuestionTypeMultiChoice  QuestionType = "multi_choice" // Multiple choice
	QuestionTypeSingleChoice QuestionType = "single_choice" // Single choice
	QuestionTypeText         QuestionType = "text"         // Free text
	QuestionTypeYesNo        QuestionType = "yes_no"       // Yes/No
)

type QuestionTemplate struct {
	sharedModels.BaseModel
	Category     string        `gorm:"not null" json:"category"`
	Name         string        `gorm:"not null" json:"name"`
	Description  string        `json:"description"`
	Text         string        `gorm:"not null" json:"text"`
	Type         QuestionType  `gorm:"not null" json:"type"`
	Options      pq.StringArray `gorm:"type:text[]" json:"options" swaggertype:"array,string"`
	MinValue     *int          `json:"min_value"`
	MaxValue     *int          `json:"max_value"`
	MinLabel     string        `json:"min_label"`
	MaxLabel     string        `json:"max_label"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags" swaggertype:"array,string"`
	IsActive     bool          `gorm:"default:true" json:"is_active"`
}

// Request/Response models for API
type CreateQuestionnaireRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	DishID      *uuid.UUID `json:"dish_id"`
	IsDefault   bool       `json:"is_default"`
}

type UpdateQuestionnaireRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
	IsActive    bool   `json:"is_active"`
}

type GenerateQuestionnaireRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	IsDefault   bool   `json:"is_default"`
}

type GeneratedQuestion struct {
	Text     string         `json:"text"`
	Type     QuestionType   `json:"type"`
	Options  pq.StringArray `json:"options,omitempty" swaggertype:"array,string"`
	MinValue *int           `json:"min_value,omitempty"`
	MaxValue *int           `json:"max_value,omitempty"`
	MinLabel string         `json:"min_label,omitempty"`
	MaxLabel string         `json:"max_label,omitempty"`
}

// Request/Response models for Questions
type CreateQuestionRequest struct {
	Text         string         `json:"text" binding:"required"`
	Type         QuestionType   `json:"type" binding:"required"`
	IsRequired   bool           `json:"is_required"`
	Options      pq.StringArray `json:"options,omitempty" swaggertype:"array,string"`
	MinValue     *int           `json:"min_value,omitempty"`
	MaxValue     *int           `json:"max_value,omitempty"`
	MinLabel     string         `json:"min_label,omitempty"`
	MaxLabel     string         `json:"max_label,omitempty"`
}

type UpdateQuestionRequest struct {
	Text         string         `json:"text"`
	Type         QuestionType   `json:"type"`
	IsRequired   bool           `json:"is_required"`
	Options      pq.StringArray `json:"options,omitempty" swaggertype:"array,string"`
	MinValue     *int           `json:"min_value,omitempty"`
	MaxValue     *int           `json:"max_value,omitempty"`
	MinLabel     string         `json:"min_label,omitempty"`
	MaxLabel     string         `json:"max_label,omitempty"`
}
