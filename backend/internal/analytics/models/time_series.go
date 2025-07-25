package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TimeSeriesMetric represents a single data point in a time series
type TimeSeriesMetric struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	AccountID      uuid.UUID  `gorm:"not null;index" json:"account_id"`
	OrganizationID uuid.UUID  `gorm:"not null;index" json:"organization_id"`
	ProductID      *uuid.UUID `gorm:"index" json:"product_id,omitempty"`
	QuestionID     *uuid.UUID `gorm:"index" json:"question_id,omitempty"`
	MetricType     string     `gorm:"not null;index" json:"metric_type"`
	MetricName     string     `gorm:"not null;index" json:"metric_name"`
	Value          float64    `gorm:"not null" json:"value"`
	Count          int64      `gorm:"not null;default:0" json:"count"`
	Timestamp      time.Time  `gorm:"not null;index" json:"timestamp"`
	Granularity    string     `gorm:"not null;index" json:"granularity"`
	Metadata       *string    `gorm:"type:jsonb" json:"metadata,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// BeforeCreate generates a UUID for the TimeSeriesMetric if not set
func (tsm *TimeSeriesMetric) BeforeCreate(tx *gorm.DB) error {
	if tsm.ID == uuid.Nil {
		tsm.ID = uuid.New()
	}
	return nil
}

// TimeSeriesComparison represents a comparison between two time periods
type TimeSeriesComparison struct {
	MetricType    string             `json:"metric_type"`
	MetricName    string             `json:"metric_name"`
	Period1       TimePeriodMetrics  `json:"period1"`
	Period2       TimePeriodMetrics  `json:"period2"`
	Change        float64            `json:"change"`
	ChangePercent float64            `json:"change_percent"`
	Trend         string             `json:"trend"`
}

// TimePeriodMetrics represents metrics for a specific time period
type TimePeriodMetrics struct {
	StartDate     time.Time          `json:"start_date"`
	EndDate       time.Time          `json:"end_date"`
	Value         float64            `json:"value"`
	Count         int64              `json:"count"`
	Average       float64            `json:"average"`
	Min           float64            `json:"min"`
	Max           float64            `json:"max"`
	DataPoints    []TimeSeriesPoint  `json:"data_points"`
}

// TimeSeriesPoint represents a single point in a time series chart
type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Count     int64     `json:"count"`
}

// TimeSeriesRequest represents a request for time series data
type TimeSeriesRequest struct {
	OrganizationID uuid.UUID  `json:"organization_id"`
	ProductID      *uuid.UUID `json:"product_id,omitempty"`
	QuestionID     *uuid.UUID `json:"question_id,omitempty"`
	MetricTypes    []string   `json:"metric_types"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	Granularity    string     `json:"granularity"`
	GroupBy        []string   `json:"group_by,omitempty"`
}

// TimeSeriesResponse represents the response containing time series data
type TimeSeriesResponse struct {
	Request    TimeSeriesRequest        `json:"request"`
	Series     []TimeSeriesData         `json:"series"`
	Summary    TimeSeriesSummary        `json:"summary"`
}

// TimeSeriesData represents a complete time series dataset
type TimeSeriesData struct {
	MetricType  string             `json:"metric_type"`
	MetricName  string             `json:"metric_name"`
	ProductID   *uuid.UUID         `json:"product_id,omitempty"`
	ProductName string             `json:"product_name,omitempty"`
	Points      []TimeSeriesPoint  `json:"points"`
	Statistics  TimeSeriesStats    `json:"statistics"`
	Metadata    *string            `json:"metadata,omitempty"`
}

// TimeSeriesStats represents statistical information about a time series
type TimeSeriesStats struct {
	Total        float64 `json:"total"`
	Average      float64 `json:"average"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Count        int64   `json:"count"`
	TrendDirection string `json:"trend_direction"`
	TrendStrength  float64 `json:"trend_strength"`
}

// TimeSeriesSummary represents a summary of all time series in the response
type TimeSeriesSummary struct {
	TotalDataPoints int                    `json:"total_data_points"`
	DateRange       DateRange              `json:"date_range"`
	Granularity     string                 `json:"granularity"`
	MetricsSummary  map[string]interface{} `json:"metrics_summary"`
}

// DateRange represents a date range
type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// ComparisonRequest represents a request to compare two time periods
type ComparisonRequest struct {
	OrganizationID uuid.UUID  `json:"organization_id"`
	ProductID      *uuid.UUID `json:"product_id,omitempty"`
	QuestionID     *uuid.UUID `json:"question_id,omitempty"`
	MetricTypes    []string   `json:"metric_types"`
	Period1Start   time.Time  `json:"period1_start"`
	Period1End     time.Time  `json:"period1_end"`
	Period2Start   time.Time  `json:"period2_start"`
	Period2End     time.Time  `json:"period2_end"`
}

// ComparisonResponse represents the response for time period comparison
type ComparisonResponse struct {
	Request     ComparisonRequest      `json:"request"`
	Comparisons []TimeSeriesComparison `json:"comparisons"`
	Insights    []ComparisonInsight    `json:"insights"`
}

// ComparisonInsight represents an insight from the comparison
type ComparisonInsight struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Message     string  `json:"message"`
	MetricType  string  `json:"metric_type"`
	Change      float64 `json:"change"`
	Recommendation string `json:"recommendation,omitempty"`
}

// MetricType constants
const (
	MetricTypeFeedbackCount     = "feedback_count"
	MetricTypeAverageRating     = "average_rating"
	MetricTypeResponseRate      = "response_rate"
	MetricTypeCompletionRate    = "completion_rate"
	MetricTypeSentimentScore    = "sentiment_score"
	MetricTypeQuestionScore     = "question_score"
	MetricTypeQRScanCount       = "qr_scan_count"
	MetricTypeConversionRate    = "conversion_rate"
	MetricTypeResponseTime      = "response_time"
	MetricTypeCustomerSatisfaction = "customer_satisfaction"
)

// Granularity constants
const (
	GranularityHourly  = "hourly"
	GranularityDaily   = "daily"
	GranularityWeekly  = "weekly"
	GranularityMonthly = "monthly"
)

// Trend constants
const (
	TrendImproving = "improving"
	TrendDeclining = "declining"
	TrendStable    = "stable"
	TrendVolatile  = "volatile"
)