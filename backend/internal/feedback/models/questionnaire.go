package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	menuModels "kyooar/internal/menu/models"
	organizationModels "kyooar/internal/organization/models"
	sharedModels "kyooar/internal/shared/models"
)

type Questionnaire struct {
	sharedModels.BaseModel
	OrganizationID uuid.UUID          `gorm:"not null" json:"organization_id"`
	Organization   organizationModels.Organization         `json:"organization,omitempty"`
	ProductID       *uuid.UUID         `json:"product_id"`
	Product         *menuModels.Product              `json:"product,omitempty"`
	Name         string             `gorm:"not null" json:"name"`
	Description  string             `json:"description"`
	IsDefault    bool               `gorm:"default:false" json:"is_default"`
	IsActive     bool               `gorm:"default:true" json:"is_active"`
	Questions    []Question         `json:"questions,omitempty"`
}

type Question struct {
	sharedModels.BaseModel
	ProductID       uuid.UUID        `gorm:"not null;index" json:"product_id"`
	Product         *menuModels.Product `json:"product,omitempty"`
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
	QuestionTypeRating       QuestionType = "rating"      	QuestionTypeScale        QuestionType = "scale"       	QuestionTypeMultiChoice  QuestionType = "multi_choice"	QuestionTypeSingleChoice QuestionType = "single_choice"	QuestionTypeText         QuestionType = "text"        	QuestionTypeYesNo        QuestionType = "yes_no"      )

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

type CreateQuestionnaireRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	ProductID      *uuid.UUID `json:"product_id"`
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

type BatchQuestionsRequest struct {
	ProductIDs []uuid.UUID `json:"product_ids" binding:"required"`
}

type BatchQuestionResponse struct {
	ID        uuid.UUID    `json:"id"`
	ProductID uuid.UUID    `json:"product_id"`
	Text      string       `json:"text"`
	Type      QuestionType `json:"type"`
}
