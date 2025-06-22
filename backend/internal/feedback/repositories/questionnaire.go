package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/feedback/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	Create(ctx context.Context, questionnaire *models.Questionnaire) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Questionnaire, error)
	FindByDishID(ctx context.Context, dishID uuid.UUID) (*models.Questionnaire, error)
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Questionnaire, error)
	Update(ctx context.Context, questionnaire *models.Questionnaire) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type questionnaireRepository struct {
	*sharedRepos.BaseRepository[models.Questionnaire]
}

func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	return &questionnaireRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Questionnaire](db),
	}
}

func (r *questionnaireRepository) FindByDishID(ctx context.Context, dishID uuid.UUID) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	err := r.DB.WithContext(ctx).Preload("Questions").
		Where("dish_id = ? AND is_active = ?", dishID, true).
		First(&questionnaire).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Questionnaire, error) {
	var questionnaires []models.Questionnaire
	err := r.DB.WithContext(ctx).Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&questionnaires).Error
	return questionnaires, err
}

type QuestionTemplateRepository interface {
	FindAll(ctx context.Context) ([]models.QuestionTemplate, error)
	FindByCategory(ctx context.Context, category string) ([]models.QuestionTemplate, error)
}

type questionTemplateRepository struct {
	*sharedRepos.BaseRepository[models.QuestionTemplate]
}

func NewQuestionTemplateRepository(db *gorm.DB) QuestionTemplateRepository {
	return &questionTemplateRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.QuestionTemplate](db),
	}
}

func (r *questionTemplateRepository) FindAll(ctx context.Context) ([]models.QuestionTemplate, error) {
	var templates []models.QuestionTemplate
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).
		Order("category, name").
		Find(&templates).Error
	return templates, err
}

func (r *questionTemplateRepository) FindByCategory(ctx context.Context, category string) ([]models.QuestionTemplate, error) {
	var templates []models.QuestionTemplate
	err := r.DB.WithContext(ctx).Where("category = ? AND is_active = ?", category, true).
		Order("name").
		Find(&templates).Error
	return templates, err
}
