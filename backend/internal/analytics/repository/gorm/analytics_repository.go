package gorm

import (
	"context"
	"time"

	"github.com/google/uuid"
	models "kyooar/internal/analytics/model"
	feedbackModels "kyooar/internal/feedback/models"
	qrcodeModels "kyooar/internal/qrcode/models"
	"gorm.io/gorm"
)

type AnalyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *AnalyticsRepository {
	return &AnalyticsRepository{
		db: db,
	}
}

func (r *AnalyticsRepository) GetQuestionChartData(ctx context.Context, questionID uuid.UUID, filters map[string]interface{}) (*models.ChartData, error) {
	var result models.ChartData
	
	query := r.db.WithContext(ctx).
		Table("feedback_responses fr").
		Joins("JOIN feedbacks f ON fr.feedback_id = f.id").
		Joins("LEFT JOIN questions q ON fr.question_id = q.id").
		Where("fr.question_id = ?", questionID)
	
	if dateFrom, ok := filters["date_from"]; ok {
		query = query.Where("DATE(f.created_at) >= DATE(?)", dateFrom)
	}
	if dateTo, ok := filters["date_to"]; ok {
		query = query.Where("DATE(f.created_at) <= DATE(?)", dateTo)
	}
	if productID, ok := filters["product_id"]; ok {
		query = query.Where("f.product_id = ?", productID)
	}
	
	var aggregatedData []struct {
		ResponseValue interface{} `gorm:"column:response_value"`
		ResponseCount int64       `gorm:"column:response_count"`
	}
	
	err := query.
		Select(`
			fr.response_value,
			COUNT(*) as response_count
		`).
		Group("fr.response_value").
		Scan(&aggregatedData).Error
		
	if err != nil {
		return nil, err
	}
	
	result.QuestionID = questionID
	result.Data = make(map[string]interface{})
	for _, data := range aggregatedData {
		if data.ResponseValue != nil {
			result.Data[data.ResponseValue.(string)] = data.ResponseCount
		}
	}
	
	return &result, nil
}

func (r *AnalyticsRepository) GetOrganizationChartDataBatch(ctx context.Context, organizationID uuid.UUID, questionIDs []uuid.UUID, filters map[string]interface{}) (map[uuid.UUID]*models.ChartData, error) {
	resultMap := make(map[uuid.UUID]*models.ChartData)
	
	for _, questionID := range questionIDs {
		chartData, err := r.GetQuestionChartData(ctx, questionID, filters)
		if err != nil {
			return nil, err
		}
		resultMap[questionID] = chartData
	}
	
	return resultMap, nil
}

func (r *AnalyticsRepository) GetFeedbackCounts(ctx context.Context, organizationID uuid.UUID) (*models.FeedbackCounts, error) {
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

func (r *AnalyticsRepository) GetQRCodeMetrics(ctx context.Context, organizationID uuid.UUID) (*models.QRCodeMetrics, error) {
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

func (r *AnalyticsRepository) GetProductRatingsAndCounts(ctx context.Context, organizationID uuid.UUID, productIDs []uuid.UUID) (map[uuid.UUID]models.ProductMetrics, error) {
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