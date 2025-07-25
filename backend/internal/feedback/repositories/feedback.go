package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/feedback/models"
	sharedModels "kyooar/internal/shared/models"
	sharedRepos "kyooar/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type FeedbackFilter struct {
	Search      string     `json:"search,omitempty"`
	RatingMin   *int       `json:"rating_min,omitempty"`
	RatingMax   *int       `json:"rating_max,omitempty"`
	DateFrom    *time.Time `json:"date_from,omitempty"`
	DateTo      *time.Time `json:"date_to,omitempty"`
	ProductID      *uuid.UUID `json:"product_id,omitempty"`
	IsComplete  *bool      `json:"is_complete,omitempty"`
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Feedback, error)
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error)
	FindByOrganizationIDWithFilters(ctx context.Context, organizationID uuid.UUID, req sharedModels.PageRequest, filters FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error)
	FindByProductID(ctx context.Context, productID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error)
	CountByOrganizationID(ctx context.Context, organizationID uuid.UUID, since time.Time) (int64, error)
	CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error)
	CountByQRCodeID(ctx context.Context, qrCodeID uuid.UUID) (int64, error)
	CountByQRCodeIDs(ctx context.Context, qrCodeIDs []uuid.UUID) (map[uuid.UUID]int64, error)
	GetAverageRating(ctx context.Context, organizationID uuid.UUID, productID *uuid.UUID) (float64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// Analytics methods
	FindByOrganizationIDForAnalytics(ctx context.Context, organizationID uuid.UUID, limit int) ([]models.Feedback, error)
	FindByProductIDForAnalytics(ctx context.Context, productID uuid.UUID, limit int) ([]models.Feedback, error)
	GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]models.Question, error)
}

type feedbackRepository struct {
	*sharedRepos.BaseRepository[models.Feedback]
}

func NewFeedbackRepository(i *do.Injector) (FeedbackRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &feedbackRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Feedback](db),
	}, nil
}

func (r *feedbackRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error) {
	var feedbacks []models.Feedback
	var total int64
	
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	
	// Count total
	r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("organization_id = ?", organizationID).Count(&total)
	
	// Get data
	query := r.DB.WithContext(ctx).Preload("Product").Preload("QRCode").
		Where("organization_id = ?", organizationID).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("created_at DESC")
	
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	
	// Populate question text for all feedback in batch
	if err := r.populateQuestionDataBatch(ctx, feedbacks); err != nil {
		// Log error but don't fail the entire request
		fmt.Printf("Error populating question data in batch: %v\n", err)
	}
	
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}
	
	return &sharedModels.PageResponse[models.Feedback]{
		Data:       feedbacks,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) FindByOrganizationIDWithFilters(ctx context.Context, organizationID uuid.UUID, req sharedModels.PageRequest, filters FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error) {
	var feedbacks []models.Feedback
	var total int64
	
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	
	// Build base query
	baseQuery := r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("organization_id = ?", organizationID)
	
	// Apply filters
	baseQuery = r.applyFilters(baseQuery, filters)
	
	// Count total
	baseQuery.Count(&total)
	
	// Get data with preloads
	query := r.DB.WithContext(ctx).Preload("Product").Preload("QRCode").
		Where("organization_id = ?", organizationID)
	
	// Apply the same filters to the data query
	query = r.applyFilters(query, filters)
	
	query = query.Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("created_at DESC")
	
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	
	// Populate question text for all feedback in batch
	if err := r.populateQuestionDataBatch(ctx, feedbacks); err != nil {
		// Log error but don't fail the entire request
		fmt.Printf("Error populating question data in batch: %v\n", err)
	}
	
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}
	
	return &sharedModels.PageResponse[models.Feedback]{
		Data:       feedbacks,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) applyFilters(query *gorm.DB, filters FeedbackFilter) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where(
			"LOWER(COALESCE(customer_name, '')) LIKE ? OR LOWER(COALESCE(customer_email, '')) LIKE ?",
			searchTerm, searchTerm,
		)
	}
	
	// Rating filters
	if filters.RatingMin != nil {
		query = query.Where("overall_rating >= ?", *filters.RatingMin)
	}
	if filters.RatingMax != nil {
		query = query.Where("overall_rating <= ?", *filters.RatingMax)
	}
	
	// Date filters
	if filters.DateFrom != nil {
		query = query.Where("DATE(created_at) >= DATE(?)", *filters.DateFrom)
	}
	if filters.DateTo != nil {
		query = query.Where("DATE(created_at) <= DATE(?)", *filters.DateTo)
	}
	
	// Product filter
	if filters.ProductID != nil {
		query = query.Where("product_id = ?", *filters.ProductID)
	}
	
	// Completion filter
	if filters.IsComplete != nil {
		query = query.Where("is_complete = ?", *filters.IsComplete)
	}
	
	return query
}

// populateQuestionData populates the question text and type for each response in the feedback
func (r *feedbackRepository) populateQuestionData(ctx context.Context, feedback *models.Feedback) error {
	if len(feedback.Responses) == 0 {
		return nil
	}
	
	// Extract question IDs from responses
	var questionIDs []uuid.UUID
	for _, response := range feedback.Responses {
		questionIDs = append(questionIDs, response.QuestionID)
	}
	
	// Query questions table to get question data and types
	var questions []struct {
		ID   uuid.UUID `gorm:"column:id"`
		Text string    `gorm:"column:text"`
		Type string    `gorm:"column:type"`
	}
	
	if err := r.DB.WithContext(ctx).Table("questions").
		Select("id, text, type").
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		return err
	}
	
	// Create maps for quick lookup
	questionTextMap := make(map[uuid.UUID]string)
	questionTypeMap := make(map[uuid.UUID]string)
	for _, question := range questions {
		questionTextMap[question.ID] = question.Text
		questionTypeMap[question.ID] = question.Type
	}
	
	// Update responses with question text and type
	for i := range feedback.Responses {
		if text, exists := questionTextMap[feedback.Responses[i].QuestionID]; exists {
			feedback.Responses[i].QuestionText = text
		}
		if qType, exists := questionTypeMap[feedback.Responses[i].QuestionID]; exists {
			feedback.Responses[i].QuestionType = models.QuestionType(qType)
		}
	}
	
	return nil
}

// populateQuestionDataBatch populates question data for multiple feedback items in a single query
func (r *feedbackRepository) populateQuestionDataBatch(ctx context.Context, feedbacks []models.Feedback) error {
	if len(feedbacks) == 0 {
		return nil
	}
	
	// Collect all unique question IDs from all feedback items
	questionIDSet := make(map[uuid.UUID]bool)
	for _, feedback := range feedbacks {
		for _, response := range feedback.Responses {
			if response.QuestionID != uuid.Nil {
				questionIDSet[response.QuestionID] = true
			}
		}
	}
	
	if len(questionIDSet) == 0 {
		return nil
	}
	
	// Convert set to slice
	var questionIDs []uuid.UUID
	for id := range questionIDSet {
		questionIDs = append(questionIDs, id)
	}
	
	// Single query to get all question data including scale labels
	var questions []struct {
		ID       uuid.UUID `gorm:"column:id"`
		Text     string    `gorm:"column:text"`
		Type     string    `gorm:"column:type"`
		MinLabel string    `gorm:"column:min_label"`
		MaxLabel string    `gorm:"column:max_label"`
		MinValue *int      `gorm:"column:min_value"`
		MaxValue *int      `gorm:"column:max_value"`
	}
	
	if err := r.DB.WithContext(ctx).Table("questions").
		Select("id, text, type, min_label, max_label, min_value, max_value").
		Where("id IN ?", questionIDs).
		Find(&questions).Error; err != nil {
		return err
	}
	
	// Create lookup maps
	questionTextMap := make(map[uuid.UUID]string)
	questionTypeMap := make(map[uuid.UUID]string)
	questionDataMap := make(map[uuid.UUID]struct {
		MinLabel string
		MaxLabel string
		MinValue *int
		MaxValue *int
	})
	for _, question := range questions {
		questionTextMap[question.ID] = question.Text
		questionTypeMap[question.ID] = question.Type
		questionDataMap[question.ID] = struct {
			MinLabel string
			MaxLabel string
			MinValue *int
			MaxValue *int
		}{
			MinLabel: question.MinLabel,
			MaxLabel: question.MaxLabel,
			MinValue: question.MinValue,
			MaxValue: question.MaxValue,
		}
	}
	
	// Update all feedback responses
	for i := range feedbacks {
		for j := range feedbacks[i].Responses {
			questionID := feedbacks[i].Responses[j].QuestionID
			if text, exists := questionTextMap[questionID]; exists {
				feedbacks[i].Responses[j].QuestionText = text
			}
			if qType, exists := questionTypeMap[questionID]; exists {
				feedbacks[i].Responses[j].QuestionType = models.QuestionType(qType)
			}
		}
	}
	
	return nil
}

func (r *feedbackRepository) FindByProductID(ctx context.Context, productID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error) {
	var feedbacks []models.Feedback
	var total int64
	
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	
	// Count total
	r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("product_id = ?", productID).Count(&total)
	
	// Get data
	query := r.DB.WithContext(ctx).Preload("QRCode").
		Where("product_id = ?", productID).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("created_at DESC")
	
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	
	totalPages := int(total) / req.Limit
	if int(total)%req.Limit > 0 {
		totalPages++
	}
	
	return &sharedModels.PageResponse[models.Feedback]{
		Data:       feedbacks,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) CountByOrganizationID(ctx context.Context, organizationID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Feedback{}).
		Where("organization_id = ? AND created_at >= ?", organizationID, since).
		Count(&count).Error
	return count, err
}

func (r *feedbackRepository) CountByProductID(ctx context.Context, productID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Feedback{}).
		Where("product_id = ?", productID).
		Count(&count).Error
	return count, err
}

func (r *feedbackRepository) CountByQRCodeID(ctx context.Context, qrCodeID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Feedback{}).
		Where("qr_code_id = ?", qrCodeID).
		Count(&count).Error
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
		Model(&models.Feedback{}).
		Select("qr_code_id, COUNT(*) as count").
		Where("qr_code_id IN ?", qrCodeIDs).
		Group("qr_code_id").
		Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	// Convert to map
	countMap := make(map[uuid.UUID]int64)
	for _, r := range results {
		countMap[r.QRCodeID] = r.Count
	}
	
	// Ensure all QR codes have an entry (even if 0)
	for _, id := range qrCodeIDs {
		if _, exists := countMap[id]; !exists {
			countMap[id] = 0
		}
	}
	
	return countMap, nil
}

func (r *feedbackRepository) GetAverageRating(ctx context.Context, organizationID uuid.UUID, productID *uuid.UUID) (float64, error) {
	var avg float64
	query := r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("organization_id = ?", organizationID)
	
	if productID != nil {
		query = query.Where("product_id = ?", *productID)
		fmt.Printf("🔍 GetAverageRating called WITH product filter: %s\n", *productID)
	} else {
		fmt.Printf("🔍 GetAverageRating called WITHOUT product filter (organization-wide)\n")
		// Debug: Check individual ratings for organization-wide query
		var allRatings []int
		r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("organization_id = ?", organizationID).Pluck("overall_rating", &allRatings)
		fmt.Printf("🔍 All ratings for organization %s: %v\n", organizationID, allRatings)
	}
	
	err := query.Select("COALESCE(AVG(overall_rating), 0)").Row().Scan(&avg)
	fmt.Printf("🔍 AVG query result: %v, error: %v, productID: %v\n", avg, err, productID)
	
	// If we got 0 and it's a organization-wide query, something's wrong
	if avg == 0 && productID == nil {
		fmt.Printf("🚨 WARNING: Organization-wide average is 0, this seems wrong!\n")
	}
	
	return avg, err
}

func (r *feedbackRepository) FindByOrganizationIDForAnalytics(ctx context.Context, organizationID uuid.UUID, limit int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.DB.WithContext(ctx).
		Preload("Product").
		Where("organization_id = ?", organizationID).
		Limit(limit).
		Order("created_at DESC").
		Find(&feedbacks).Error
	
	if err != nil {
		return nil, err
	}
	
	// Populate question data to avoid N+1 queries
	if err := r.populateQuestionDataBatch(ctx, feedbacks); err != nil {
		// Log error but don't fail the entire request since analytics can work without question details
		return feedbacks, nil
	}
	
	return feedbacks, nil
}

func (r *feedbackRepository) FindByProductIDForAnalytics(ctx context.Context, productID uuid.UUID, limit int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.DB.WithContext(ctx).
		Preload("Product").
		Where("product_id = ?", productID).
		Limit(limit).
		Order("created_at DESC").
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]models.Question, error) {
	var questions []models.Question
	err := r.DB.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}
