package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/shared/models"
	"gorm.io/gorm"
)

type FeedbackRepository interface {
	Create(feedback *models.Feedback) error
	FindByID(id uuid.UUID, preloads ...string) (*models.Feedback, error)
	FindByRestaurantID(restaurantID uuid.UUID, req models.PageRequest) (*models.PageResponse, error)
	FindByDishID(dishID uuid.UUID, req models.PageRequest) (*models.PageResponse, error)
	CountByRestaurantID(restaurantID uuid.UUID, since time.Time) (int64, error)
	GetAverageRating(restaurantID uuid.UUID, dishID *uuid.UUID) (float64, error)
	Delete(id uuid.UUID) error
}

type feedbackRepository struct {
	*repositories.BaseRepository[models.Feedback]
}

func NewFeedbackRepository(db *gorm.DB) FeedbackRepository {
	return &feedbackRepository{
		repositories.BaseRepository: Newrepositories.BaseRepository[models.Feedback](db),
	}
}

func (r *feedbackRepository) FindByRestaurantID(restaurantID uuid.UUID, req PageRequest) (*PageResponse[models.Feedback], error) {
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
	r.db.Model(&models.Feedback{}).Where("restaurant_id = ?", restaurantID).Count(&total)
	
	// Get data
	query := r.db.Preload("Dish").Preload("QRCode").
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
	
	return &PageResponse[models.Feedback]{
		Data:       feedbacks,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) FindByDishID(dishID uuid.UUID, req PageRequest) (*PageResponse[models.Feedback], error) {
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
	r.db.Model(&models.Feedback{}).Where("dish_id = ?", dishID).Count(&total)
	
	// Get data
	query := r.db.Preload("QRCode").
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
	
	return &PageResponse[models.Feedback]{
		Data:       feedbacks,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: totalPages,
	}, nil
}

func (r *feedbackRepository) CountByRestaurantID(restaurantID uuid.UUID, since time.Time) (int64, error) {
	var count int64
	err := r.db.Model(&models.Feedback{}).
		Where("restaurant_id = ? AND created_at >= ?", restaurantID, since).
		Count(&count).Error
	return count, err
}

func (r *feedbackRepository) GetAverageRating(restaurantID uuid.UUID, dishID *uuid.UUID) (float64, error) {
	var avg float64
	query := r.db.Model(&models.Feedback{}).Where("restaurant_id = ?", restaurantID)
	if dishID != nil {
		query = query.Where("dish_id = ?", *dishID)
	}
	err := query.Select("AVG(overall_rating)").Row().Scan(&avg)
	return avg, err
}
