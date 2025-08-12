package analyticsmodel

import (
	"github.com/google/uuid"
)

type ProductAnalytics struct {
	ProductID     uuid.UUID `json:"product_id"`
	ProductName   string    `json:"product_name"`
	AverageRating float64   `json:"average_rating"`
	TotalFeedback int64     `json:"total_feedback"`
}

type OrganizationAnalytics struct {
	OrganizationID      uuid.UUID          `json:"organization_id"`
	OrganizationName    string             `json:"organization_name"`
	TotalFeedback       int64              `json:"total_feedback"`
	AverageRating       float64            `json:"average_rating"`
	FeedbackToday       int64              `json:"feedback_today"`
	FeedbackThisWeek    int64              `json:"feedback_this_week"`
	FeedbackThisMonth   int64              `json:"feedback_this_month"`
	TopRatedProducts    []ProductAnalytics `json:"top_rated_products"`
	LowestRatedProducts []ProductAnalytics `json:"lowest_rated_products"`
}

type DashboardMetricsResponse struct {
	Success bool                    `json:"success"`
	Data    interface{}             `json:"data"`
}

type ProductAnalyticsResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ProductInsightsResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type TimeSeriesAPIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ComparisonAPIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type ChartDataResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

type MetricsCollectionResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}