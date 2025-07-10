package models

import (
	"time"

	"github.com/google/uuid"
)

// QuestionMetric represents analytics for a single question
type QuestionMetric struct {
	QuestionID   uuid.UUID              `json:"question_id"`
	QuestionText string                 `json:"question_text"`
	QuestionType string                 `json:"question_type"`
	ResponseCount int64                 `json:"response_count"`
	
	// For numeric types (rating, scale)
	AverageScore  *float64              `json:"average_score,omitempty"`
	MinScore      *float64              `json:"min_score,omitempty"`
	MaxScore      *float64              `json:"max_score,omitempty"`
	
	// For choice types (multi_choice, single_choice, yes_no)
	OptionDistribution map[string]int64  `json:"option_distribution,omitempty"`
	
	// For text responses
	TextResponses []string             `json:"text_responses,omitempty"`
	CommonThemes  []string             `json:"common_themes,omitempty"`
	
	// Sentiment for all types (derived)
	PositiveRate  float64              `json:"positive_rate"`
	NeutralRate   float64              `json:"neutral_rate"`
	NegativeRate  float64              `json:"negative_rate"`
}

// DishInsights represents comprehensive analytics for a dish
type DishInsights struct {
	DishID        uuid.UUID           `json:"dish_id"`
	DishName      string              `json:"dish_name"`
	TotalFeedback int64               `json:"total_feedback"`
	CompletionRate float64            `json:"completion_rate"`
	
	// Overall metrics (calculated from all numeric questions)
	OverallScore  float64             `json:"overall_score"`
	ScoreTrend    string              `json:"score_trend"` // "improving", "declining", "stable"
	
	// Question-level metrics
	Questions     []QuestionMetric    `json:"questions"`
	
	// Key insights
	BestAspects   []string           `json:"best_aspects"`   // Top performing questions
	NeedsAttention []string          `json:"needs_attention"` // Low performing questions
	
	// Time-based
	LastFeedback  time.Time          `json:"last_feedback"`
	WeeklyChange  float64            `json:"weekly_change"`
}

// RestaurantInsights represents restaurant-wide analytics
type RestaurantInsights struct {
	RestaurantID      uuid.UUID         `json:"restaurant_id"`
	RestaurantName    string            `json:"restaurant_name"`
	Period            string            `json:"period"` // "today", "week", "month", "all"
	
	// Overall metrics
	TotalFeedback     int64             `json:"total_feedback"`
	ActiveDishes      int               `json:"active_dishes"`
	AverageSatisfaction float64         `json:"average_satisfaction"`
	RecommendationRate float64          `json:"recommendation_rate"`
	SentimentScore    float64           `json:"sentiment_score"`
	
	// Trends
	FeedbackTrend     []TrendPoint      `json:"feedback_trend"`
	SatisfactionTrend []TrendPoint      `json:"satisfaction_trend"`
	
	// Top/Bottom performers
	TopDishes         []DishSummary     `json:"top_dishes"`
	BottomDishes      []DishSummary     `json:"bottom_dishes"`
	
	// Issues requiring attention
	CriticalIssues    []Issue           `json:"critical_issues"`
}

// DashboardMetrics represents operational metrics for the dashboard
type DashboardMetrics struct {
	// Core operational metrics
	TotalFeedbacks    int64             `json:"total_feedbacks"`
	TodaysFeedback    int64             `json:"todays_feedback"`
	TrendVsYesterday  float64           `json:"trend_vs_yesterday"` // percentage
	
	// QR Code operational metrics
	ActiveQRCodes     int64             `json:"active_qr_codes"`
	TotalQRScans      int64             `json:"total_qr_scans"`
	ScansToday        int64             `json:"scans_today"`
	CompletionRate    float64           `json:"completion_rate"` // feedback submissions / QR scans
	
	// Performance metrics
	AverageResponseTime float64         `json:"average_response_time"` // minutes from scan to submission
	PeakHours         []int             `json:"peak_hours"` // top 3 hours of day for responses
	
	// Device/Platform analytics
	DeviceBreakdown   map[string]int64  `json:"device_breakdown"` // platform/browser usage
	
	// QR Code performance analytics
	QRPerformance     []QRCodePerformance `json:"qr_performance"` // individual QR code stats
}

// QRCodePerformance represents performance data for individual QR codes
type QRCodePerformance struct {
	ID             uuid.UUID `json:"id"`
	Label          string    `json:"label"`
	Location       string    `json:"location,omitempty"`
	RestaurantID   uuid.UUID `json:"restaurant_id"`
	RestaurantName string    `json:"restaurant_name"`
	ScansCount     int64     `json:"scans_count"`
	FeedbackCount  int64     `json:"feedback_count"`
	ConversionRate float64   `json:"conversion_rate"` // feedback / scans * 100
}

// Supporting types
type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type DishSummary struct {
	DishID        uuid.UUID `json:"dish_id"`
	DishName      string    `json:"dish_name"`
	Score         float64   `json:"score"`
	FeedbackCount int64     `json:"feedback_count"`
	Trend         string    `json:"trend"` // "up", "down", "stable"
}

type Issue struct {
	DishID       uuid.UUID `json:"dish_id"`
	DishName     string    `json:"dish_name"`
	QuestionText string    `json:"question_text"`
	IssueType    string    `json:"issue_type"` // "low_score", "negative_trend", "complaints"
	Severity     string    `json:"severity"`   // "critical", "warning", "info"
	Description  string    `json:"description"`
}

type QuickIssue struct {
	Title       string `json:"title"`
	Count       int    `json:"count"`
	Severity    string `json:"severity"`
	ActionLink  string `json:"action_link"`
}

type FeedbackSummary struct {
	FeedbackID   uuid.UUID `json:"feedback_id"`
	DishName     string    `json:"dish_name"`
	CustomerName string    `json:"customer_name"`
	Score        float64   `json:"score"`
	Sentiment    string    `json:"sentiment"` // "positive", "neutral", "negative"
	Highlight    string    `json:"highlight"` // Key comment or issue
	CreatedAt    time.Time `json:"created_at"`
}