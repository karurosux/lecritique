package analyticsmodel

import (
	"time"

	"github.com/google/uuid"
)

type TimeSeriesRequest struct {
	OrganizationID uuid.UUID    `json:"organization_id" validate:"required,uuid"`
	MetricTypes    []string     `json:"metric_types" validate:"required,min=1"`
	StartDate      time.Time    `json:"start_date" validate:"required"`
	EndDate        time.Time    `json:"end_date" validate:"required"`
	Granularity    string       `json:"granularity" validate:"required,oneof=hourly daily weekly monthly"`
	ProductID      *uuid.UUID   `json:"product_id,omitempty"`
	QuestionID     *uuid.UUID   `json:"question_id,omitempty"`
}

type ComparisonRequest struct {
	OrganizationID uuid.UUID    `json:"organization_id" validate:"required,uuid"`
	MetricTypes    []string     `json:"metric_types" validate:"required,min=1"`
	Period1Start   time.Time    `json:"period1_start" validate:"required"`
	Period1End     time.Time    `json:"period1_end" validate:"required"`
	Period2Start   time.Time    `json:"period2_start" validate:"required"`
	Period2End     time.Time    `json:"period2_end" validate:"required"`
	ProductID      *uuid.UUID   `json:"product_id,omitempty"`
	QuestionID     *uuid.UUID   `json:"question_id,omitempty"`
}

type ChartFiltersRequest struct {
	DateFrom    string     `json:"date_from,omitempty"`
	DateTo      string     `json:"date_to,omitempty"`
	ProductID   *uuid.UUID `json:"product_id,omitempty"`
	QuestionID  *uuid.UUID `json:"question_id,omitempty"`
}