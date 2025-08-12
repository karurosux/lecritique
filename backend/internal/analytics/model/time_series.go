package analyticsmodel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (tsm *TimeSeriesMetric) BeforeCreate(tx *gorm.DB) error {
	if tsm.ID == uuid.Nil {
		tsm.ID = uuid.New()
	}
	return nil
}

type TimeSeriesComparison struct {
	MetricType    string             `json:"metric_type"`
	MetricName    string             `json:"metric_name"`
	Period1       TimePeriodMetrics  `json:"period1"`
	Period2       TimePeriodMetrics  `json:"period2"`
	Change        float64            `json:"change"`
	ChangePercent float64            `json:"change_percent"`
	Trend         string             `json:"trend"`
	Metadata      *string            `json:"metadata,omitempty"`
}

type TimePeriodMetrics struct {
	StartDate        time.Time           `json:"start_date"`
	EndDate          time.Time           `json:"end_date"`
	Value            float64             `json:"value"`
	Count            int64               `json:"count"`
	Average          float64             `json:"average"`
	Min              float64             `json:"min"`
	Max              float64             `json:"max"`
	DataPoints       []TimeSeriesPoint   `json:"data_points"`
	ChoiceDistribution map[string]int64   `json:"choice_distribution,omitempty"`
	MostPopularChoice  *ChoiceInfo        `json:"most_popular_choice,omitempty"`
	TopChoices         []ChoiceInfo       `json:"top_choices,omitempty"`
}

type ChoiceInfo struct {
	Choice string `json:"choice"`
	Count  int64  `json:"count"`
}

type TimeSeriesPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
	Count     int64     `json:"count"`
}

type TimeSeriesFilters struct {
	OrganizationID *uuid.UUID `json:"organization_id,omitempty"`
	ProductID      *uuid.UUID `json:"product_id,omitempty"`
	QuestionID     *uuid.UUID `json:"question_id,omitempty"`
	MetricType     string     `json:"metric_type,omitempty"`
	StartDate      time.Time  `json:"start_date"`
	EndDate        time.Time  `json:"end_date"`
	Granularity    string     `json:"granularity,omitempty"`
}

type TimeSeriesResponse struct {
	Request    TimeSeriesRequest        `json:"request"`
	Series     []TimeSeriesData         `json:"series"`
	Summary    TimeSeriesSummary        `json:"summary"`
}

type TimeSeriesData struct {
	MetricType   string             `json:"metric_type"`
	MetricName   string             `json:"metric_name"`
	ProductID    *uuid.UUID         `json:"product_id,omitempty"`
	ProductName  string             `json:"product_name,omitempty"`
	Points       []TimeSeriesPoint  `json:"points"`
	Statistics   TimeSeriesStats    `json:"statistics"`
	Metadata     map[string]any     `json:"metadata,omitempty"`
	ChoiceSeries []ChoiceSeriesData `json:"choice_series,omitempty"`
}

type ChoiceSeriesData struct {
	Choice     string            `json:"choice"`
	Points     []TimeSeriesPoint `json:"points"`
	Statistics TimeSeriesStats   `json:"statistics"`
}

type TimeSeriesStats struct {
	Total        float64 `json:"total"`
	Average      float64 `json:"average"`
	Min          float64 `json:"min"`
	Max          float64 `json:"max"`
	Count        int64   `json:"count"`
	TrendDirection string `json:"trend_direction"`
	TrendStrength  float64 `json:"trend_strength"`
}

type TimeSeriesSummary struct {
	TotalDataPoints int                    `json:"total_data_points"`
	DateRange       DateRange              `json:"date_range"`
	Granularity     string                 `json:"granularity"`
	MetricsSummary  map[string]interface{} `json:"metrics_summary"`
}

type DateRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
type ComparisonResponse struct {
	Request     ComparisonRequest      `json:"request"`
	Comparisons []TimeSeriesComparison `json:"comparisons"`
	Insights    []ComparisonInsight    `json:"insights"`
}

type ComparisonInsight struct {
	Type        string  `json:"type"`
	Severity    string  `json:"severity"`
	Message     string  `json:"message"`
	MetricType  string  `json:"metric_type"`
	Change      float64 `json:"change"`
	Recommendation string `json:"recommendation,omitempty"`
}

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

const (
	GranularityHourly  = "hourly"
	GranularityDaily   = "daily"
	GranularityWeekly  = "weekly"
	GranularityMonthly = "monthly"
)

const (
	TrendImproving = "improving"
	TrendDeclining = "declining"
	TrendStable    = "stable"
	TrendVolatile  = "volatile"
)