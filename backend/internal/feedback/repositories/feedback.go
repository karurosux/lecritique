package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/feedback/models"
	sharedModels "github.com/lecritique/api/internal/shared/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type FeedbackFilter struct {
	Search      string     `json:"search,omitempty"`
	RatingMin   *int       `json:"rating_min,omitempty"`
	RatingMax   *int       `json:"rating_max,omitempty"`
	DateFrom    *time.Time `json:"date_from,omitempty"`
	DateTo      *time.Time `json:"date_to,omitempty"`
	DishID      *uuid.UUID `json:"dish_id,omitempty"`
	IsComplete  *bool      `json:"is_complete,omitempty"`
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Feedback, error)
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error)
	FindByRestaurantIDWithFilters(ctx context.Context, restaurantID uuid.UUID, req sharedModels.PageRequest, filters FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error)
	FindByDishID(ctx context.Context, dishID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error)
	CountByRestaurantID(ctx context.Context, restaurantID uuid.UUID, since time.Time) (int64, error)
	CountByDishID(ctx context.Context, dishID uuid.UUID) (int64, error)
	CountByQRCodeID(ctx context.Context, qrCodeID uuid.UUID) (int64, error)
	GetAverageRating(ctx context.Context, restaurantID uuid.UUID, dishID *uuid.UUID) (float64, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// Analytics methods
	FindByRestaurantIDForAnalytics(ctx context.Context, restaurantID uuid.UUID, limit int) ([]models.Feedback, error)
	FindByDishIDForAnalytics(ctx context.Context, dishID uuid.UUID, limit int) ([]models.Feedback, error)
	GetQuestionsByDishID(ctx context.Context, dishID uuid.UUID) ([]models.Question, error)
}

type feedbackRepository struct {
	*sharedRepos.BaseRepository[models.Feedback]
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Feedback](db),
	}
}

func (r *feedbackRepository) FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error) {
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
	r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("restaurant_id = ?", restaurantID).Count(&total)
	
	// Get data
	query := r.DB.WithContext(ctx).Preload("Dish").Preload("QRCode").
		Where("restaurant_id = ?", restaurantID).
		Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("created_at DESC")
	
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	
	// Populate question text for each feedback
	for i := range feedbacks {
		if err := r.populateQuestionData(ctx, &feedbacks[i]); err != nil {
			// Log error but don't fail the entire request
			fmt.Printf("Error populating question data for feedback %s: %v\n", feedbacks[i].ID, err)
		}
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

func (r *feedbackRepository) FindByRestaurantIDWithFilters(ctx context.Context, restaurantID uuid.UUID, req sharedModels.PageRequest, filters FeedbackFilter) (*sharedModels.PageResponse[models.Feedback], error) {
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
	baseQuery := r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("restaurant_id = ?", restaurantID)
	
	// Apply filters
	baseQuery = r.applyFilters(baseQuery, filters)
	
	// Count total
	baseQuery.Count(&total)
	
	// Get data with preloads
	query := r.DB.WithContext(ctx).Preload("Dish").Preload("QRCode").
		Where("restaurant_id = ?", restaurantID)
	
	// Apply the same filters to the data query
	query = r.applyFilters(query, filters)
	
	query = query.Limit(req.Limit).
		Offset((req.Page - 1) * req.Limit).
		Order("created_at DESC")
	
	if err := query.Find(&feedbacks).Error; err != nil {
		return nil, err
	}
	
	// Populate question text for each feedback
	for i := range feedbacks {
		if err := r.populateQuestionData(ctx, &feedbacks[i]); err != nil {
			// Log error but don't fail the entire request
			fmt.Printf("Error populating question data for feedback %s: %v\n", feedbacks[i].ID, err)
		}
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
	
	// Dish filter
	if filters.DishID != nil {
		query = query.Where("dish_id = ?", *filters.DishID)
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

func (r *feedbackRepository) FindByDishID(ctx context.Context, dishID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error) {
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
	r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("dish_id = ?", dishID).Count(&total)
	
	// Get data
	query := r.DB.WithContext(ctx).Preload("QRCode").
		Where("dish_id = ?", dishID).
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

func (r *feedbackRepository) CountByRestaurantID(ctx context.Context, restaurantID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Feedback{}).
		Where("restaurant_id = ? AND created_at >= ?", restaurantID, since).
		Count(&count).Error
	return count, err
}

func (r *feedbackRepository) CountByDishID(ctx context.Context, dishID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Feedback{}).
		Where("dish_id = ?", dishID).
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

func (r *feedbackRepository) GetAverageRating(ctx context.Context, restaurantID uuid.UUID, dishID *uuid.UUID) (float64, error) {
	var avg float64
	query := r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("restaurant_id = ?", restaurantID)
	
	if dishID != nil {
		query = query.Where("dish_id = ?", *dishID)
		fmt.Printf("🔍 GetAverageRating called WITH dish filter: %s\n", *dishID)
	} else {
		fmt.Printf("🔍 GetAverageRating called WITHOUT dish filter (restaurant-wide)\n")
		// Debug: Check individual ratings for restaurant-wide query
		var allRatings []int
		r.DB.WithContext(ctx).Model(&models.Feedback{}).Where("restaurant_id = ?", restaurantID).Pluck("overall_rating", &allRatings)
		fmt.Printf("🔍 All ratings for restaurant %s: %v\n", restaurantID, allRatings)
	}
	
	err := query.Select("COALESCE(AVG(overall_rating), 0)").Row().Scan(&avg)
	fmt.Printf("🔍 AVG query result: %v, error: %v, dishID: %v\n", avg, err, dishID)
	
	// If we got 0 and it's a restaurant-wide query, something's wrong
	if avg == 0 && dishID == nil {
		fmt.Printf("🚨 WARNING: Restaurant-wide average is 0, this seems wrong!\n")
	}
	
	return avg, err
}

func (r *feedbackRepository) FindByRestaurantIDForAnalytics(ctx context.Context, restaurantID uuid.UUID, limit int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.DB.WithContext(ctx).
		Preload("Dish").
		Where("restaurant_id = ?", restaurantID).
		Limit(limit).
		Order("created_at DESC").
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) FindByDishIDForAnalytics(ctx context.Context, dishID uuid.UUID, limit int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	err := r.DB.WithContext(ctx).
		Preload("Dish").
		Where("dish_id = ?", dishID).
		Limit(limit).
		Order("created_at DESC").
		Find(&feedbacks).Error
	return feedbacks, err
}

func (r *feedbackRepository) GetQuestionsByDishID(ctx context.Context, dishID uuid.UUID) ([]models.Question, error) {
	var questions []models.Question
	err := r.DB.WithContext(ctx).
		Where("dish_id = ?", dishID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}
