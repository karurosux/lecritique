package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/lecritique/api/shared/models"
)

type Questionnaire struct {
	models.BaseModel
	RestaurantID uuid.UUID  `gorm:"not null" json:"restaurant_id"`
	DishID       *uuid.UUID `json:"dish_id"`
	Name         string     `gorm:"not null" json:"name"`
	Description  string     `json:"description"`
	IsDefault    bool       `gorm:"default:false" json:"is_default"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`
	Questions    []Question `json:"questions,omitempty"`
}

type Question struct {
	models.BaseModel
	QuestionnaireID uuid.UUID     `gorm:"not null" json:"questionnaire_id"`
	Questionnaire   Questionnaire `json:"questionnaire,omitempty"`
	Text            string        `gorm:"not null" json:"text"`
	Type            QuestionType  `gorm:"not null" json:"type"`
	IsRequired      bool          `gorm:"default:true" json:"is_required"`
	DisplayOrder    int           `gorm:"default:0" json:"display_order"`
	Options         pq.StringArray `gorm:"type:text[]" json:"options"`
	MinValue        *int          `json:"min_value"`
	MaxValue        *int          `json:"max_value"`
	MinLabel        string        `json:"min_label"`
	MaxLabel        string        `json:"max_label"`
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
	models.BaseModel
	Category     string        `gorm:"not null" json:"category"`
	Name         string        `gorm:"not null" json:"name"`
	Description  string        `json:"description"`
	Text         string        `gorm:"not null" json:"text"`
	Type         QuestionType  `gorm:"not null" json:"type"`
	Options      pq.StringArray `gorm:"type:text[]" json:"options"`
	MinValue     *int          `json:"min_value"`
	MaxValue     *int          `json:"max_value"`
	MinLabel     string        `json:"min_label"`
	MaxLabel     string        `json:"max_label"`
	Tags         pq.StringArray `gorm:"type:text[]" json:"tags"`
	IsActive     bool          `gorm:"default:true" json:"is_active"`
}