package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"kyooar/internal/analytics/models"
	"kyooar/internal/shared/logger"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type TimeSeriesRepository interface {
	Create(ctx context.Context, metric *models.TimeSeriesMetric) error
	CreateBatch(ctx context.Context, metrics []models.TimeSeriesMetric) error
	GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) ([]models.TimeSeriesMetric, error)
	GetComparison(ctx context.Context, request models.ComparisonRequest) ([]models.TimeSeriesMetric, []models.TimeSeriesMetric, error)
	DeleteOldMetrics(ctx context.Context, before time.Time) error
	DeleteOrganizationMetrics(ctx context.Context, organizationID uuid.UUID) error
	GetLatestMetric(ctx context.Context, organizationID uuid.UUID, metricType string, productID *uuid.UUID) (*models.TimeSeriesMetric, error)
}

type timeSeriesRepository struct {
	db *gorm.DB
}

func NewTimeSeriesRepository(i *do.Injector) (TimeSeriesRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &timeSeriesRepository{db: db}, nil
}

func (r *timeSeriesRepository) Create(ctx context.Context, metric *models.TimeSeriesMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

func (r *timeSeriesRepository) CreateBatch(ctx context.Context, metrics []models.TimeSeriesMetric) error {
	if len(metrics) == 0 {
		return nil
	}
	
	// Temporarily silence database logs for bulk operations to reduce noise
	dbWithSilentLogger := r.db.Session(&gorm.Session{Logger: gormLogger.Default.LogMode(gormLogger.Silent)})
	
	// Use larger batch size for better performance (500 instead of 100)
	batchSize := 500
	for i := 0; i < len(metrics); i += batchSize {
		end := i + batchSize
		if end > len(metrics) {
			end = len(metrics)
		}
		
		// Use CreateInBatches for more efficient bulk inserts
		if err := dbWithSilentLogger.WithContext(ctx).CreateInBatches(metrics[i:end], batchSize).Error; err != nil {
			logger.Error("Failed to create time series batch", err, logrus.Fields{
				"batch_start": i,
				"batch_end":   end,
				"total_metrics": len(metrics),
			})
			return err
		}
	}
	
	return nil
}

func (r *timeSeriesRepository) GetTimeSeries(ctx context.Context, request models.TimeSeriesRequest) ([]models.TimeSeriesMetric, error) {
	var metrics []models.TimeSeriesMetric
	
	query := r.db.WithContext(ctx).Model(&models.TimeSeriesMetric{}).
		Where("organization_id = ?", request.OrganizationID).
		Where("timestamp >= ? AND timestamp <= ?", request.StartDate, request.EndDate).
		Where("granularity = ?", request.Granularity)
	
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
		})
		return nil, err
	}
	
	
	return metrics, nil
}

func (r *timeSeriesRepository) GetComparison(ctx context.Context, request models.ComparisonRequest) ([]models.TimeSeriesMetric, []models.TimeSeriesMetric, error) {
	// Get metrics for period 1
	period1Request := models.TimeSeriesRequest{
		OrganizationID: request.OrganizationID,
		ProductID:      request.ProductID,
		QuestionID:     request.QuestionID,
		MetricTypes:    request.MetricTypes,
		StartDate:      request.Period1Start,
		EndDate:        request.Period1End,
		Granularity:    "daily", // Default to daily for comparisons
	}
	
	period1Metrics, err := r.GetTimeSeries(ctx, period1Request)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get period 1 metrics: %w", err)
	}
	
	// Get metrics for period 2
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

func (r *timeSeriesRepository) DeleteOldMetrics(ctx context.Context, before time.Time) error {
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

func (r *timeSeriesRepository) DeleteOrganizationMetrics(ctx context.Context, organizationID uuid.UUID) error {
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

func (r *timeSeriesRepository) GetLatestMetric(ctx context.Context, organizationID uuid.UUID, metricType string, productID *uuid.UUID) (*models.TimeSeriesMetric, error) {
	var metric models.TimeSeriesMetric
	
	query := r.db.WithContext(ctx).Model(&models.TimeSeriesMetric{}).
		Where("organization_id = ?", organizationID).
		Where("metric_type = ?", metricType)
	
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	} else {
		query = query.Where("product_id IS NULL")
	}
	
	err := query.Order("timestamp DESC").First(&metric).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	
	return &metric, nil
}