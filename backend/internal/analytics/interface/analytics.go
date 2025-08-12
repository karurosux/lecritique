package analyticsinterface

import (
	"context"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/analytics/models"
)

type AnalyticsRepository interface {
	GetQuestionChartData(ctx context.Context, questionID uuid.UUID, filters map[string]interface{}) (*models.ChartData, error)
	GetOrganizationChartDataBatch(ctx context.Context, organizationID uuid.UUID, questionIDs []uuid.UUID, filters map[string]interface{}) (map[uuid.UUID]*models.ChartData, error)
	GetFeedbackCounts(ctx context.Context, organizationID uuid.UUID) (*models.FeedbackCounts, error)
	GetQRCodeMetrics(ctx context.Context, organizationID uuid.UUID) (*models.QRCodeMetrics, error)
	GetProductRatingsAndCounts(ctx context.Context, organizationID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]models.ProductMetrics, error)
}

type TimeSeriesRepository interface {
	Create(ctx context.Context, metric *models.TimeSeriesMetric) error
	FindByFilters(ctx context.Context, filters models.TimeSeriesFilters) ([]*models.TimeSeriesMetric, error)
	DeleteOlderThan(ctx context.Context, organizationID uuid.UUID, cutoffTime time.Time) error
	GetAggregatedData(ctx context.Context, organizationID uuid.UUID, metricType string, granularity string, startDate, endDate time.Time, productID *uuid.UUID, questionID *uuid.UUID) ([]*models.TimeSeriesDataPoint, error)
	BatchCreate(ctx context.Context, metrics []*models.TimeSeriesMetric) error
	DeleteOrganizationMetrics(ctx context.Context, organizationID uuid.UUID) error
	CreateBatch(ctx context.Context, metrics []models.TimeSeriesMetric) error
	GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) ([]models.TimeSeriesMetric, error)
	GetComparison(ctx context.Context, request models.ComparisonRequest) ([]models.TimeSeriesMetric, []models.TimeSeriesMetric, error)
	DeleteOldMetrics(ctx context.Context, before time.Time) error
	HasMetricsWithPattern(ctx context.Context, pattern string) bool
	GetMetricTypesByPattern(ctx context.Context, pattern string) []string
}

type AnalyticsService interface {
	GetDashboardMetrics(ctx context.Context, organizationID uuid.UUID) (*models.DashboardMetrics, error)
	GetProductInsights(ctx context.Context, productID uuid.UUID) (*models.ProductInsights, error)
	GetOrganizationInsights(ctx context.Context, organizationID uuid.UUID, period string) (*models.OrganizationInsights, error)
	GetOrganizationChartData(ctx context.Context, organizationID uuid.UUID, filters map[string]interface{}) (*models.OrganizationChartData, error)
	GetQuestionChartData(ctx context.Context, questionID uuid.UUID, filters map[string]interface{}) (*models.ChartData, error)
	GetProductAnalyticsBatch(ctx context.Context, organizationID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]models.ProductAnalytics, error)
}

type TimeSeriesService interface {
	CollectMetrics(ctx context.Context, organizationID uuid.UUID) error
	GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) (*models.TimeSeriesResponse, error)
	GetComparison(ctx context.Context, request models.ComparisonRequest) (*models.ComparisonResponse, error)
	CleanupOldMetrics(ctx context.Context, retentionDays int) error
}