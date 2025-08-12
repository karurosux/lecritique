package models

import (
	"time"

	"github.com/google/uuid"
)

type QuestionMetric struct {
	QuestionID   uuid.UUID              `json:"question_id"`
	QuestionText string                 `json:"question_text"`
	QuestionType string                 `json:"question_type"`
	ResponseCount int64                 `json:"response_count"`
	
	AverageScore  *float64              `json:"average_score,omitempty"`
	MinScore      *float64              `json:"min_score,omitempty"`
	MaxScore      *float64              `json:"max_score,omitempty"`
	
	OptionDistribution map[string]int64  `json:"option_distribution,omitempty"`
	
	TextResponses []string             `json:"text_responses,omitempty"`
	CommonThemes  []string             `json:"common_themes,omitempty"`
	
	PositiveRate  float64              `json:"positive_rate"`
	NeutralRate   float64              `json:"neutral_rate"`
	NegativeRate  float64              `json:"negative_rate"`
}

type ProductInsights struct {
	ProductID        uuid.UUID           `json:"product_id"`
	ProductName      string              `json:"product_name"`
	TotalFeedback int64               `json:"total_feedback"`
	CompletionRate float64            `json:"completion_rate"`
	
	OverallScore  float64             `json:"overall_score"`
	ScoreTrend    string              `json:"score_trend"`
	
	Questions     []QuestionMetric    `json:"questions"`
	
	BestAspects   []string           `json:"best_aspects"`
	NeedsAttention []string          `json:"needs_attention"`
	
	LastFeedback  time.Time          `json:"last_feedback"`
	WeeklyChange  float64            `json:"weekly_change"`
}

type OrganizationInsights struct {
	OrganizationID      uuid.UUID         `json:"organization_id"`
	OrganizationName    string            `json:"organization_name"`
	Period            string            `json:"period"`
	
	TotalFeedback     int64             `json:"total_feedback"`
	ActiveProducts      int               `json:"active_products"`
	AverageSatisfaction float64         `json:"average_satisfaction"`
	RecommendationRate float64          `json:"recommendation_rate"`
	SentimentScore    float64           `json:"sentiment_score"`
	
	FeedbackTrend     []TrendPoint      `json:"feedback_trend"`
	SatisfactionTrend []TrendPoint      `json:"satisfaction_trend"`
	
	TopProducts         []ProductSummary     `json:"top_products"`
	BottomProducts      []ProductSummary     `json:"bottom_products"`
	
	CriticalIssues    []Issue           `json:"critical_issues"`
}

type DashboardMetrics struct {
	TotalFeedbacks    int64             `json:"total_feedbacks"`
	TodaysFeedback    int64             `json:"todays_feedback"`
	TrendVsYesterday  float64           `json:"trend_vs_yesterday"`
	
	ActiveQRCodes     int64             `json:"active_qr_codes"`
	TotalQRScans      int64             `json:"total_qr_scans"`
	ScansToday        int64             `json:"scans_today"`
	CompletionRate    float64           `json:"completion_rate"`
	
	AverageResponseTime float64         `json:"average_response_time"`
	PeakHours         []int             `json:"peak_hours"`
	
	DeviceBreakdown   map[string]int64  `json:"device_breakdown"`
	
	QRPerformance     []QRCodePerformance `json:"qr_performance"`
}

type QRCodePerformance struct {
	ID             uuid.UUID  `json:"id"`
	Label          string     `json:"label"`
	Location       string     `json:"location,omitempty"`
	OrganizationID   uuid.UUID  `json:"organization_id"`
	OrganizationName string     `json:"organization_name"`
	ScansCount     int64      `json:"scans_count"`
	FeedbackCount  int64      `json:"feedback_count"`
	ConversionRate float64    `json:"conversion_rate"`
	LastScan       *time.Time `json:"last_scan,omitempty"`
	IsActive       bool       `json:"is_active"`
}

type TrendPoint struct {
	Date  time.Time `json:"date"`
	Value float64   `json:"value"`
}

type ProductSummary struct {
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`
	Score         float64   `json:"score"`
	FeedbackCount int64     `json:"feedback_count"`
	Trend         string    `json:"trend"`
}

type Issue struct {
	ProductID       uuid.UUID `json:"product_id"`
	ProductName     string    `json:"product_name"`
	QuestionText string    `json:"question_text"`
	IssueType    string    `json:"issue_type"`
	Severity     string    `json:"severity"`
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
	ProductName     string    `json:"product_name"`
	CustomerName string    `json:"customer_name"`
	Score        float64   `json:"score"`
	Sentiment    string    `json:"sentiment"`
	Highlight    string    `json:"highlight"`
	CreatedAt    time.Time `json:"created_at"`
}


type ChartData struct {
	QuestionID   uuid.UUID              `json:"question_id"`
	QuestionText string                 `json:"question_text"`
	QuestionType string                 `json:"question_type"`
	ChartType    string                 `json:"chart_type"`
	ProductID       *uuid.UUID             `json:"product_id,omitempty"`
	ProductName     string                 `json:"product_name,omitempty"`
	Data         map[string]interface{} `json:"data"`
}

type RatingDistribution struct {
	Scale         int                    `json:"scale"`
	Distribution  map[string]int64       `json:"distribution"`
	Average       float64                `json:"average"`
	Total         int64                  `json:"total"`
	Percentages   map[string]float64     `json:"percentages"`
}

  
type ChoiceDistribution struct {
	Options       map[string]int64       `json:"options"`
	Total         int64                  `json:"total"`
	Percentages   map[string]float64     `json:"percentages"`
	IsMultiChoice bool                   `json:"is_multi_choice"`
	Combinations  []ChoiceCombination    `json:"combinations,omitempty"`
}

type ChoiceCombination struct {
	Options []string `json:"options"`
	Count   int64    `json:"count"`
	Percentage float64 `json:"percentage"`
}

type TextSentiment struct {
	Positive    int64    `json:"positive"`
	Neutral     int64    `json:"neutral"`
	Negative    int64    `json:"negative"`
	Total       int64    `json:"total"`
	Samples     []string `json:"samples"`
	Keywords    []string `json:"keywords"`
}

type OrganizationChartData struct {
	OrganizationID uuid.UUID   `json:"organization_id"`
	Charts       []ChartData `json:"charts"`
	Summary      struct {
		TotalResponses int64     `json:"total_responses"`
		DateRange      struct {
			Start time.Time `json:"start"`
			End   time.Time `json:"end"`
		} `json:"date_range"`
		FiltersApplied map[string]interface{} `json:"filters_applied"`
	} `json:"summary"`
}

type ProductAnalytics struct {
	ProductID        uuid.UUID `json:"product_id"`
	ProductName      string    `json:"product_name"`
	AverageRating    float64   `json:"average_rating"`
	TotalFeedback    int64     `json:"total_feedback"`
}
