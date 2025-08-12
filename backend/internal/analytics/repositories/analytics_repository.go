package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/analytics/models"
	feedbackModels "kyooar/internal/feedback/models"
	qrcodeModels "kyooar/internal/qrcode/models"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	GetFeedbackCounts(ctx context.Context, organizationID uuid.UUID) (*models.FeedbackCounts, error)
	GetQRCodeMetrics(ctx context.Context, organizationID uuid.UUID) (*models.QRCodeMetrics, error)
	GetFeedbackWithQRCodes(ctx context.Context, organizationID uuid.UUID, limit int) ([]models.FeedbackWithQRCode, error)
	GetQRCodePerformanceMetrics(ctx context.Context, organizationID uuid.UUID) ([]models.QRCodePerformanceData, error)
	GetOrganizationChartData(ctx context.Context, organizationID uuid.UUID, filters models.ChartFilters) (*models.ChartDataResult, error)
	GetTimeSeriesData(ctx context.Context, organizationID uuid.UUID, startDate, endDate time.Time) ([]models.TimeSeriesDataPoint, error)
	GetProductRatingsAndCounts(ctx context.Context, organizationID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]models.ProductMetrics, error)
}

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(i *do.Injector) (AnalyticsRepository, error) {
	return &analyticsRepository{
		db: do.MustInvoke[*gorm.DB](i),
	}, nil
}

func (r *analyticsRepository) GetFeedbackCounts(ctx context.Context, organizationID uuid.UUID) (*models.FeedbackCounts, error) {
	var result models.FeedbackCounts
	
	todayStart := time.Now().Truncate(24 * time.Hour)
	yesterdayStart := todayStart.AddDate(0, 0, -1)
	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	
	err := r.db.WithContext(ctx).
		Model(&feedbackModels.Feedback{}).
		Select(`
			COUNT(*) as total,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as today,
			COUNT(CASE WHEN created_at >= ? AND created_at < ? THEN 1 END) as yesterday,
			COUNT(CASE WHEN created_at >= ? THEN 1 END) as recent_30_days
		`, todayStart, yesterdayStart, todayStart, thirtyDaysAgo).
		Where("organization_id = ?", organizationID).
		Scan(&result).Error
		
	return &result, err
}

func (r *analyticsRepository) GetQRCodeMetrics(ctx context.Context, organizationID uuid.UUID) (*models.QRCodeMetrics, error) {
	var result models.QRCodeMetrics
	todayStart := time.Now().Truncate(24 * time.Hour)
	
	err := r.db.WithContext(ctx).
		Model(&qrcodeModels.QRCode{}).
		Select(`
			COUNT(*) as total_qr_codes,
			COUNT(CASE WHEN is_active = true THEN 1 END) as active_count,
			SUM(scans_count) as total_scans,
			COUNT(CASE WHEN last_scanned_at >= ? THEN 1 END) as scans_today
		`, todayStart).
		Where("organization_id = ?", organizationID).
		Scan(&result).Error
		
	return &result, err
}

func (r *analyticsRepository) GetFeedbackWithQRCodes(ctx context.Context, organizationID uuid.UUID, limit int) ([]models.FeedbackWithQRCode, error) {
	var results []models.FeedbackWithQRCode
	
	err := r.db.WithContext(ctx).
		Table("feedbacks f").
		Select(`
			f.id as feedback_id,
			f.created_at as feedback_created_at,
			f.qr_code_id,
			f.device_info,
			q.last_scanned_at as qr_last_scanned_at
		`).
		Joins("LEFT JOIN qr_codes q ON f.qr_code_id = q.id").
		Where("f.organization_id = ?", organizationID).
		Order("f.created_at DESC").
		Limit(limit).
		Scan(&results).Error
		
	return results, err
}

func (r *analyticsRepository) GetQRCodePerformanceMetrics(ctx context.Context, organizationID uuid.UUID) ([]models.QRCodePerformanceData, error) {
	var results []models.QRCodePerformanceData
	
	err := r.db.WithContext(ctx).
		Table("qr_codes q").
		Select(`
			q.id,
			q.label,
			q.organization_id,
			q.scans_count,
			q.last_scanned_at,
			q.is_active,
			q.location,
			COUNT(f.id) as feedback_count
		`).
		Joins("LEFT JOIN feedbacks f ON q.id = f.qr_code_id").
		Where("q.organization_id = ?", organizationID).
		Group("q.id, q.label, q.organization_id, q.scans_count, q.last_scanned_at, q.is_active, q.location").
		Scan(&results).Error
		
	return results, err
}

func (r *analyticsRepository) GetOrganizationChartData(ctx context.Context, organizationID uuid.UUID, filters models.ChartFilters) (*models.ChartDataResult, error) {
	result := &models.ChartDataResult{
		OrganizationID: organizationID,
	}
	
	query := r.db.WithContext(ctx).
		Table("feedbacks f").
		Joins("JOIN feedback_responses fr ON f.id = fr.feedback_id").
		Joins("LEFT JOIN questions q ON fr.question_id = q.id").
		Where("f.organization_id = ?", organizationID)
	
	if filters.DateFrom != nil {
		query = query.Where("DATE(f.created_at) >= DATE(?)", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("DATE(f.created_at) <= DATE(?)", *filters.DateTo)
	}
	if filters.ProductID != nil {
		query = query.Where("f.product_id = ?", *filters.ProductID)
	}
	
	var aggregatedData []struct {
		QuestionID   uuid.UUID `gorm:"column:question_id"`
		QuestionText string    `gorm:"column:question_text"`
		QuestionType string    `gorm:"column:question_type"`
		ProductID    uuid.UUID `gorm:"column:product_id"`
		ResponseCount int64    `gorm:"column:response_count"`
	}
	
	err := query.
		Select(`
			fr.question_id,
			COALESCE(q.text, fr.question_text) as question_text,
			COALESCE(q.type, fr.question_type) as question_type,
			f.product_id,
			COUNT(*) as response_count
		`).
		Group("fr.question_id, question_text, question_type, f.product_id").
		Scan(&aggregatedData).Error
		
	if err != nil {
		return nil, err
	}
	
	var totalCount int64
	r.db.WithContext(ctx).
		Model(&feedbackModels.Feedback{}).
		Where("organization_id = ?", organizationID).
		Count(&totalCount)
		
	result.TotalResponses = totalCount
	result.QuestionData = aggregatedData
	
	return result, nil
}

func (r *analyticsRepository) GetTimeSeriesData(ctx context.Context, organizationID uuid.UUID, startDate, endDate time.Time) ([]models.TimeSeriesDataPoint, error) {
	var results []models.TimeSeriesDataPoint
	
	err := r.db.WithContext(ctx).
		Table("feedbacks").
		Select(`
			DATE(created_at) as date,
			COUNT(*) as count,
			AVG(overall_rating) as average_rating
		`).
		Where("organization_id = ? AND created_at BETWEEN ? AND ?", organizationID, startDate, endDate).
		Group("DATE(created_at)").
		Order("date ASC").
		Scan(&results).Error
		
	return results, err
}

func (r *analyticsRepository) GetProductRatingsAndCounts(ctx context.Context, organizationID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]models.ProductMetrics, error) {
	if len(productIDs) == 0 {
		return make(map[uuid.UUID]models.ProductMetrics), nil
	}
	
	type result struct {
		ProductID     uuid.UUID `gorm:"column:product_id"`
		AverageRating float64   `gorm:"column:average_rating"`
		FeedbackCount int64     `gorm:"column:feedback_count"`
	}
	
	var results []result
	err := r.db.WithContext(ctx).
		Model(&feedbackModels.Feedback{}).
		Select(`
			product_id,
			COALESCE(AVG(overall_rating), 0) as average_rating,
			COUNT(*) as feedback_count
		`).
		Where("organization_id = ? AND product_id IN ?", organizationID, productIDs).
		Group("product_id").
		Scan(&results).Error
		
	if err != nil {
		return nil, err
	}
	
	metricsMap := make(map[uuid.UUID]models.ProductMetrics)
	for _, r := range results {
		metricsMap[r.ProductID] = models.ProductMetrics{
			AverageRating: r.AverageRating,
			FeedbackCount: r.FeedbackCount,
		}
	}
	
	for _, productID := range productIDs {
		if _, exists := metricsMap[productID]; !exists {
			metricsMap[productID] = models.ProductMetrics{
				AverageRating: 0,
				FeedbackCount: 0,
			}
		}
	}
	
	return metricsMap, nil
}