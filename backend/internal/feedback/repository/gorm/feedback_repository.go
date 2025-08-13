package gorm

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	feedbackmodel "kyooar/internal/feedback/model"
	sharedModels "kyooar/internal/shared/models"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type feedbackRepository struct {
	*sharedRepos.BaseRepository[feedbackmodel.Feedback]
}

func NewFeedbackRepository(db *gorm.DB) *feedbackRepository {
	return &feedbackRepository{
		BaseRepository: sharedRepos.NewBaseRepository[feedbackmodel.Feedback](db),
	}
}

func (r *feedbackRepository) Create(ctx context.Context, feedback *feedbackmodel.Feedback) error {
	return r.BaseRepository.Create(ctx, feedback)
}

func (r *feedbackRepository) FindByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Feedback, error) {
	return r.BaseRepository.FindByID(ctx, id)
}

func (r *feedbackRepository) Update(ctx context.Context, feedback *feedbackmodel.Feedback) error {
	return r.BaseRepository.Update(ctx, feedback)
}

func (r *feedbackRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.BaseRepository.Delete(ctx, id)
}

func (r *feedbackRepository) FindByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int) (*sharedModels.PageResponse[feedbackmodel.Feedback], error) {
	var feedbacks []feedbackmodel.Feedback
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).Where("organization_id = ?", organizationID).Count(&total)

	query := r.DB.WithContext(ctx).Preload("Product").Preload("QRCode").
		Where("organization_id = ?", organizationID).
		Limit(limit).
		Offset((page - 1) * limit).
		Order("created_at DESC")

	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	if err := r.populateQuestionDataBatch(ctx, feedbacks); err != nil {
		fmt.Printf("Error populating question data in batch: %v\n", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &sharedModels.PageResponse[feedbackmodel.Feedback]{
		Data:       feedbacks,
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) FindByOrganizationIDWithFilters(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, page, limit int, filters feedbackmodel.FeedbackFilter) (*sharedModels.PageResponse[feedbackmodel.Feedback], error) {
	var feedbacks []feedbackmodel.Feedback
	var total int64

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	baseQuery := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).Where("organization_id = ?", organizationID)

	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		baseQuery = baseQuery.Where("(LOWER(customer_name) LIKE ? OR LOWER(customer_email) LIKE ? OR EXISTS (SELECT 1 FROM json_array_elements(responses) AS response WHERE LOWER(response->>'answer') LIKE ?))",
			searchTerm, searchTerm, searchTerm)
	}

	if filters.RatingMin != nil {
		baseQuery = baseQuery.Where("overall_rating >= ?", *filters.RatingMin)
	}

	if filters.RatingMax != nil {
		baseQuery = baseQuery.Where("overall_rating <= ?", *filters.RatingMax)
	}

	if filters.DateFrom != nil {
		baseQuery = baseQuery.Where("created_at >= ?", *filters.DateFrom)
	}

	if filters.DateTo != nil {
		endDate := filters.DateTo.Add(24 * time.Hour).Add(-time.Second)
		baseQuery = baseQuery.Where("created_at <= ?", endDate)
	}

	if filters.ProductID != nil {
		baseQuery = baseQuery.Where("product_id = ?", *filters.ProductID)
	}

	if filters.IsComplete != nil {
		baseQuery = baseQuery.Where("is_complete = ?", *filters.IsComplete)
	}

	baseQuery.Count(&total)

	query := baseQuery.
		Preload("Product").
		Preload("QRCode").
		Limit(limit).
		Offset((page - 1) * limit).
		Order("created_at DESC")

	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	if err := r.populateQuestionDataBatch(ctx, feedbacks); err != nil {
		fmt.Printf("Error populating question data in batch: %v\n", err)
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	return &sharedModels.PageResponse[feedbackmodel.Feedback]{
		Data:       feedbacks,
		Page:       page,
		Limit:      limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) GetStatsByOrganization(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) (*feedbackmodel.FeedbackStats, error) {
	stats := &feedbackmodel.FeedbackStats{}

	var totalFeedbacks int64
	if err := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).
		Where("organization_id = ?", organizationID).
		Count(&totalFeedbacks).Error; err != nil {
		return nil, err
	}
	stats.TotalFeedbacks = totalFeedbacks

	var avgRating sql.NullFloat64
	if err := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).
		Select("AVG(overall_rating)").
		Where("organization_id = ? AND overall_rating > 0", organizationID).
		Scan(&avgRating).Error; err != nil {
		return nil, err
	}
	if avgRating.Valid {
		stats.AverageRating = avgRating.Float64
	}

	today := time.Now().Truncate(24 * time.Hour)
	var todayFeedbacks int64
	if err := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).
		Where("organization_id = ? AND created_at >= ?", organizationID, today).
		Count(&todayFeedbacks).Error; err != nil {
		return nil, err
	}
	stats.FeedbacksToday = todayFeedbacks

	startOfWeek := time.Now().AddDate(0, 0, -int(time.Now().Weekday())).Truncate(24 * time.Hour)
	var thisWeekFeedbacks int64
	if err := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).
		Where("organization_id = ? AND created_at >= ?", organizationID, startOfWeek).
		Count(&thisWeekFeedbacks).Error; err != nil {
		return nil, err
	}
	stats.FeedbacksThisWeek = thisWeekFeedbacks

	return stats, nil
}

func (r *feedbackRepository) FindByOrganizationIDForAnalytics(ctx context.Context, organizationID uuid.UUID, limit int) ([]feedbackmodel.Feedback, error) {
	var feedbacks []feedbackmodel.Feedback
	
	query := r.DB.WithContext(ctx).
		Preload("Product").
		Where("organization_id = ?", organizationID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *feedbackRepository) FindByQuestionInPeriod(ctx context.Context, questionID uuid.UUID, startDate, endDate time.Time) ([]feedbackmodel.Feedback, error) {
	var feedbacks []feedbackmodel.Feedback

	if err := r.DB.WithContext(ctx).
		Where("created_at BETWEEN ? AND ?", startDate, endDate).
		Where("EXISTS (SELECT 1 FROM json_array_elements(responses) AS response WHERE (response->>'question_id')::uuid = ?)", questionID).
		Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *feedbackRepository) populateQuestionDataBatch(ctx context.Context, feedbacks []feedbackmodel.Feedback) error {
	questionIDs := make(map[uuid.UUID]bool)
	for _, feedback := range feedbacks {
		for _, response := range feedback.Responses {
			if response.QuestionID != uuid.Nil {
				questionIDs[response.QuestionID] = true
			}
		}
	}

	if len(questionIDs) == 0 {
		return nil
	}

	ids := make([]uuid.UUID, 0, len(questionIDs))
	for id := range questionIDs {
		ids = append(ids, id)
	}

	var questions []feedbackmodel.Question
	if err := r.DB.WithContext(ctx).Where("id IN ?", ids).Find(&questions).Error; err != nil {
		return err
	}

	questionMap := make(map[uuid.UUID]feedbackmodel.Question)
	for _, q := range questions {
		questionMap[q.ID] = q
	}

	for i := range feedbacks {
		for j := range feedbacks[i].Responses {
			if q, exists := questionMap[feedbacks[i].Responses[j].QuestionID]; exists {
				feedbacks[i].Responses[j].QuestionText = q.Text
				feedbacks[i].Responses[j].QuestionType = q.Type
			}
		}
	}

	return nil
}

func (r *feedbackRepository) CountByOrganizationID(ctx context.Context, organizationID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	query := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).Where("organization_id = ?", organizationID)
	if !since.IsZero() {
		query = query.Where("created_at >= ?", since)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *feedbackRepository) CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).Where("product_id = ?", productID).Count(&count).Error
	return count, err
}

func (r *feedbackRepository) CountByQRCodeID(ctx context.Context, qrCodeID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).Where("qr_code_id = ?", qrCodeID).Count(&count).Error
	return count, err
}

func (r *feedbackRepository) CountByQRCodeIDs(ctx context.Context, qrCodeIDs []uuid.UUID) (map[uuid.UUID]int64, error) {
	if len(qrCodeIDs) == 0 {
		return make(map[uuid.UUID]int64), nil
	}

	type result struct {
		QRCodeID uuid.UUID `gorm:"column:qr_code_id"`
		Count    int64     `gorm:"column:count"`
	}

	var results []result
	err := r.DB.WithContext(ctx).
		Model(&feedbackmodel.Feedback{}).
		Select("qr_code_id, COUNT(*) as count").
		Where("qr_code_id IN ?", qrCodeIDs).
		Group("qr_code_id").
		Scan(&results).Error

	if err != nil {
		return nil, err
	}

	countMap := make(map[uuid.UUID]int64)
	for _, res := range results {
		countMap[res.QRCodeID] = res.Count
	}

	for _, id := range qrCodeIDs {
		if _, exists := countMap[id]; !exists {
			countMap[id] = 0
		}
	}

	return countMap, nil
}

func (r *feedbackRepository) GetAverageRating(ctx context.Context, organizationID uuid.UUID, productID *uuid.UUID) (float64, error) {
	var avgRating sql.NullFloat64
	query := r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).
		Select("AVG(overall_rating)").
		Where("organization_id = ? AND overall_rating > 0", organizationID)

	if productID != nil {
		query = query.Where("product_id = ?", *productID)
	}

	err := query.Scan(&avgRating).Error
	if err != nil {
		return 0, err
	}

	if avgRating.Valid {
		return avgRating.Float64, nil
	}
	return 0, nil
}

func (r *feedbackRepository) FindByProductID(ctx context.Context, productID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[feedbackmodel.Feedback], error) {
	var feedbacks []feedbackmodel.Feedback
	var total int64

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}

	r.DB.WithContext(ctx).Model(&feedbackmodel.Feedback{}).Where("product_id = ?", productID).Count(&total)

	query := r.DB.WithContext(ctx).Preload("Product").Preload("QRCode").
		Where("product_id = ?", productID).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("created_at DESC")

	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	if err := r.populateQuestionDataBatch(ctx, feedbacks); err != nil {
		fmt.Printf("Error populating question data in batch: %v\n", err)
	}

	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}

	return &sharedModels.PageResponse[feedbackmodel.Feedback]{
		Data:       feedbacks,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) FindByProductIDForAnalytics(ctx context.Context, productID uuid.UUID, limit int) ([]feedbackmodel.Feedback, error) {
	var feedbacks []feedbackmodel.Feedback

	query := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *feedbackRepository) GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]feedbackmodel.Question, error) {
	var questions []feedbackmodel.Question
	err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}