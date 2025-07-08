package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/feedback/models"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *models.Question) error
	GetQuestionsByDishID(ctx context.Context, dishID uuid.UUID) ([]*models.Question, error)
	GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, questionID uuid.UUID) error
	ReorderQuestions(ctx context.Context, dishID uuid.UUID, questionIDs []uuid.UUID) error
	GetMaxDisplayOrder(ctx context.Context, dishID uuid.UUID) (int, error)
	GetDishesWithQuestions(ctx context.Context, restaurantID uuid.UUID) ([]uuid.UUID, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) QuestionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepository) GetQuestionsByDishID(ctx context.Context, dishID uuid.UUID) ([]*models.Question, error) {
	var questions []*models.Question
	err := r.db.WithContext(ctx).
		Where("dish_id = ?", dishID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := r.db.WithContext(ctx).
		Preload("Dish").
		Where("id = ?", questionID).
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

func (r *questionRepository) DeleteQuestion(ctx context.Context, questionID uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Question{}, "id = ?", questionID).Error
}

func (r *questionRepository) GetMaxDisplayOrder(ctx context.Context, dishID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.WithContext(ctx).
		Model(&models.Question{}).
		Where("dish_id = ?", dishID).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *questionRepository) ReorderQuestions(ctx context.Context, dishID uuid.UUID, questionIDs []uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update display order for each question
	for i, questionID := range questionIDs {
		if err := tx.Model(&models.Question{}).
			Where("id = ? AND dish_id = ?", questionID, dishID).
			Update("display_order", i+1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *questionRepository) GetDishesWithQuestions(ctx context.Context, restaurantID uuid.UUID) ([]uuid.UUID, error) {
	var dishIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Table("questions").
		Select("DISTINCT questions.dish_id").
		Joins("JOIN dishes ON dishes.id = questions.dish_id").
		Where("dishes.restaurant_id = ?", restaurantID).
		Scan(&dishIDs).Error
	return dishIDs, err
}