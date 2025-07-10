package repositories

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/feedback/models"
	sharedModels "github.com/lecritique/api/internal/shared/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *models.Feedback) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Feedback, error)
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID, req sharedModels.PageRequest) (*sharedModels.PageResponse[models.Feedback], error)
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
