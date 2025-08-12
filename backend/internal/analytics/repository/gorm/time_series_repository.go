package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	models "kyooar/internal/analytics/model"
	"kyooar/internal/shared/logger"
	"github.com/sirupsen/logrus"
)

type TimeSeriesRepository struct {
	db *gorm.DB
}

func NewTimeSeriesRepository(db *gorm.DB) *TimeSeriesRepository {
	return &TimeSeriesRepository{db: db}
}

func (r *TimeSeriesRepository) Create(ctx context.Context, metric *models.TimeSeriesMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

func (r *TimeSeriesRepository) FindByFilters(ctx context.Context, filters models.TimeSeriesFilters) ([]*models.TimeSeriesMetric, error) {
	var metrics []*models.TimeSeriesMetric
	
	query := r.db.WithContext(ctx).Model(&models.TimeSeriesMetric{})
	
	if filters.OrganizationID != nil {
		query = query.Where("organization_id = ?", *filters.OrganizationID)
	}
	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}
	if filters.QuestionID != nil {
		query = query.Where("question_id = ?", *filters.QuestionID)
	}
	if filters.MetricType != "" {
		query = query.Where("metric_type = ?", filters.MetricType)
	}
	if !filters.StartDate.IsZero() {
		query = query.Where("timestamp >= ?", filters.StartDate)
	}
	if !filters.EndDate.IsZero() {
		query = query.Where("timestamp <= ?", filters.EndDate)
	}
	if filters.Granularity != "" {
		query = query.Where("granularity = ?", filters.Granularity)
	}
	
	err := query.Order("timestamp ASC").Find(&metrics).Error
	return metrics, err
}

func (r *TimeSeriesRepository) DeleteOlderThan(ctx context.Context, organizationID uuid.UUID, cutoffTime time.Time) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ? AND timestamp < ?", organizationID, cutoffTime).
		Delete(&models.TimeSeriesMetric{})
	
	if result.Error != nil {
		logger.Error("Failed to delete old metrics", result.Error, logrus.Fields{
			"organization_id": organizationID,
			"cutoff_time":     cutoffTime,
		})
		return result.Error
	}
	
	logger.Info("Deleted old time series metrics", logrus.Fields{
		"organization_id": organizationID,
		"count":           result.RowsAffected,
		"cutoff_time":     cutoffTime,
	})
	
	return nil
}

func (r *TimeSeriesRepository) GetAggregatedData(ctx context.Context, organizationID uuid.UUID, metricType string, granularity string, startDate, endDate time.Time, productID *uuid.UUID, questionID *uuid.UUID) ([]*models.TimeSeriesDataPoint, error) {
	var dataPoints []*models.TimeSeriesDataPoint
	
	var dateTrunc string
	switch granularity {
	case models.GranularityHourly:
		dateTrunc = "hour"
	case models.GranularityDaily:
		dateTrunc = "day"
	case models.GranularityWeekly:
		dateTrunc = "week"
	case models.GranularityMonthly:
		dateTrunc = "month"
	default:
		dateTrunc = "day"
	}
	
	selectClause := fmt.Sprintf(`
		DATE_TRUNC('%s', timestamp) as timestamp,
		COALESCE(AVG(value), 0) as value,
		COALESCE(SUM(count), 0) as count
	`, dateTrunc)
	
	groupClause := fmt.Sprintf(`DATE_TRUNC('%s', timestamp)`, dateTrunc)
	
	query := r.db.WithContext(ctx).
		Select(selectClause).
		Table("time_series_metrics").
		Where("organization_id = ?", organizationID).
		Where("metric_type = ?", metricType).
		Where("timestamp >= ? AND timestamp <= ?", startDate, endDate).
		Group(groupClause)
	
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}
	
	if questionID != nil {
		query = query.Where("question_id = ?", *questionID)
	}
	
	err := query.Order("timestamp ASC").Scan(&dataPoints).Error
	if err != nil {
		logger.Error("Failed to get aggregated time series data", err, logrus.Fields{
			"organization_id": organizationID,
			"metric_type":     metricType,
			"granularity":     granularity,
			"start_date":      startDate,
			"end_date":        endDate,
		})
		return nil, err
	}
	
	return dataPoints, nil
}

func (r *TimeSeriesRepository) BatchCreate(ctx context.Context, metrics []*models.TimeSeriesMetric) error {
	if len(metrics) == 0 {
		return nil
	}
	
	dbWithSilentLogger := r.db.Session(&gorm.Session{Logger: gormLogger.Default.LogMode(gormLogger.Silent)})
	
	batchSize := 500
	for i := 0; i < len(metrics); i += batchSize {
		end := i + batchSize
		if end > len(metrics) {
			end = len(metrics)
		}
		
		if err := dbWithSilentLogger.WithContext(ctx).CreateInBatches(metrics[i:end], batchSize).Error; err != nil {
			logger.Error("Failed to create time series batch", err, logrus.Fields{
				"batch_start":   i,
				"batch_end":     end,
				"total_metrics": len(metrics),
			})
			return err
		}
	}
	
	return nil
}

func (r *TimeSeriesRepository) DeleteOrganizationMetrics(ctx context.Context, organizationID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Delete(&models.TimeSeriesMetric{})
	
	if result.Error != nil {
		logger.Error("Failed to delete organization metrics", result.Error, logrus.Fields{
			"organization_id": organizationID,
		})
		return result.Error
	}
	
	return nil
}

func (r *TimeSeriesRepository) CreateBatch(ctx context.Context, metrics []models.TimeSeriesMetric) error {
	if len(metrics) == 0 {
		return nil
	}
	
	dbWithSilentLogger := r.db.Session(&gorm.Session{Logger: gormLogger.Default.LogMode(gormLogger.Silent)})
	
	batchSize := 500
	for i := 0; i < len(metrics); i += batchSize {
		end := i + batchSize
		if end > len(metrics) {
			end = len(metrics)
		}
		
		if err := dbWithSilentLogger.WithContext(ctx).CreateInBatches(metrics[i:end], batchSize).Error; err != nil {
			logger.Error("Failed to create time series batch", err, logrus.Fields{
				"batch_start":   i,
				"batch_end":     end,
				"total_metrics": len(metrics),
			})
			return err
		}
	}
	
	return nil
}

func (r *TimeSeriesRepository) GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) ([]models.TimeSeriesMetric, error) {
	var metrics []models.TimeSeriesMetric
	
	var dateTrunc string
	switch request.Granularity {
	case models.GranularityHourly:
		dateTrunc = "hour"
	case models.GranularityDaily:
		dateTrunc = "day"
	case models.GranularityWeekly:
		dateTrunc = "week"
	case models.GranularityMonthly:
		dateTrunc = "month"
	default:
		dateTrunc = "day"
	}
	
	selectClause := fmt.Sprintf(`
		gen_random_uuid() as id,
		account_id,
		organization_id,
		product_id,
		question_id,
		metric_type,
		metric_name,
		DATE_TRUNC('%s', timestamp) as timestamp,
		'%s' as granularity,
		COALESCE(AVG(value), 0) as value,
		COALESCE(SUM(count), 0) as count,
		metadata,
		NOW() as created_at,
		NOW() as updated_at
	`, dateTrunc, request.Granularity)
	
	groupClause := fmt.Sprintf(`
		account_id,
		organization_id,
		product_id,
		question_id,
		metric_type,
		metric_name,
		DATE_TRUNC('%s', timestamp),
		metadata
	`, dateTrunc)
	
	query := r.db.WithContext(ctx).
		Select(selectClause).
		Table("time_series_metrics").
		Where("organization_id = ?", request.OrganizationID).
		Where("timestamp >= ? AND timestamp <= ?", request.StartDate, request.EndDate).
		Where("granularity = ?", models.GranularityDaily).
		Group(groupClause)
	
	if request.ProductID != nil {
		query = query.Where("product_id = ?", *request.ProductID)
	}
	
	if request.QuestionID != nil {
		query = query.Where("question_id = ?", *request.QuestionID)
	}
	
	if len(request.MetricTypes) > 0 {
		query = query.Where("metric_type IN ?", request.MetricTypes)
	}
	
	err := query.Order("timestamp ASC").Find(&metrics).Error
	if err != nil {
		logger.Error("Failed to get time series data", err, logrus.Fields{
			"organization_id": request.OrganizationID,
			"start_date":      request.StartDate,
			"end_date":        request.EndDate,
			"granularity":     request.Granularity,
		})
		return nil, err
	}
	
	return metrics, nil
}

func (r *TimeSeriesRepository) GetComparison(ctx context.Context, request models.ComparisonRequest) ([]models.TimeSeriesMetric, []models.TimeSeriesMetric, error) {
	period1Request := models.TimeSeriesRequest{
		OrganizationID: request.OrganizationID,
		ProductID:      request.ProductID,
		QuestionID:     request.QuestionID,
		MetricTypes:    request.MetricTypes,
		StartDate:      request.Period1Start,
		EndDate:        request.Period1End,
		Granularity:    "daily",
	}
	
	period1Metrics, err := r.GetTimeSeries(ctx, period1Request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get period 1 metrics: %w", err)
	}
	
	period2Request := models.TimeSeriesRequest{
		OrganizationID: request.OrganizationID,
		ProductID:      request.ProductID,
		QuestionID:     request.QuestionID,
		MetricTypes:    request.MetricTypes,
		StartDate:      request.Period2Start,
		EndDate:        request.Period2End,
		Granularity:    "daily",
	}
	
	period2Metrics, err := r.GetTimeSeries(ctx, period2Request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get period 2 metrics: %w", err)
	}
	
	return period1Metrics, period2Metrics, nil
}

func (r *TimeSeriesRepository) DeleteOldMetrics(ctx context.Context, before time.Time) error {
	result := r.db.WithContext(ctx).
		Where("timestamp < ?", before).
		Delete(&models.TimeSeriesMetric{})
	
	if result.Error != nil {
		logger.Error("Failed to delete old metrics", result.Error, logrus.Fields{
			"before": before,
		})
		return result.Error
	}
	
	logger.Info("Deleted old time series metrics", logrus.Fields{
		"count":  result.RowsAffected,
		"before": before,
	})
	
	return nil
}

func (r *TimeSeriesRepository) HasMetricsWithPattern(ctx context.Context, pattern string) bool {
	var count int64
	
	err := r.db.WithContext(ctx).Model(&models.TimeSeriesMetric{}).
		Where("metric_type LIKE ?", pattern+"%").
		Count(&count).Error
		
	if err != nil {
		logger.Error("Failed to check metrics with pattern", err, logrus.Fields{
			"pattern": pattern,
		})
		return false
	}
	
	return count > 0
}

func (r *TimeSeriesRepository) GetMetricTypesByPattern(ctx context.Context, pattern string) []string {
	var metricTypes []string
	
	err := r.db.WithContext(ctx).Model(&models.TimeSeriesMetric{}).
		Select("DISTINCT metric_type").
		Where("metric_type LIKE ?", pattern+"%").
		Pluck("metric_type", &metricTypes).Error
		
	if err != nil {
		logger.Error("Failed to get metric types by pattern", err, logrus.Fields{
			"pattern": pattern,
		})
		return []string{}
	}
	
	return metricTypes
}