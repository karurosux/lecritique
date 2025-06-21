package repositories

import (
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/shared/models"
	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	Create(questionnaire *models.Questionnaire) error
	FindByID(id uuid.UUID, preloads ...string) (*models.Questionnaire, error)
	FindByDishID(dishID uuid.UUID) (*models.Questionnaire, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.Questionnaire, error)
	Update(questionnaire *models.Questionnaire) error
	Delete(id uuid.UUID) error
}

type questionnaireRepository struct {
	*repositories.BaseRepository[models.Questionnaire]
}

func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	return &questionnaireRepository{
		repositories.BaseRepository: Newrepositories.BaseRepository[models.Questionnaire](db),
	}
}

func (r *questionnaireRepository) FindByDishID(dishID uuid.UUID) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	err := r.db.Preload("Questions").
		Where("dish_id = ? AND is_active = ?", dishID, true).
		First(&questionnaire).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.Questionnaire, error) {
	var questionnaires []models.Questionnaire
	err := r.db.Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&questionnaires).Error
	return questionnaires, err
}

type QuestionTemplateRepository interface {
	FindAll() ([]models.QuestionTemplate, error)
	FindByCategory(category string) ([]models.QuestionTemplate, error)
}

type questionTemplateRepository struct {
	*repositories.BaseRepository[models.QuestionTemplate]
}

func NewQuestionTemplateRepository(db *gorm.DB) QuestionTemplateRepository {
	return &questionTemplateRepository{
		repositories.BaseRepository: Newrepositories.BaseRepository[models.QuestionTemplate](db),
	}
}

func (r *questionTemplateRepository) FindAll() ([]models.QuestionTemplate, error) {
	var templates []models.QuestionTemplate
	err := r.db.Where("is_active = ?", true).
		Order("category, name").
		Find(&templates).Error
	return templates, err
}

func (r *questionTemplateRepository) FindByCategory(category string) ([]models.QuestionTemplate, error) {
	var templates []models.QuestionTemplate
	err := r.db.Where("category = ? AND is_active = ?", category, true).
		Order("name").
		Find(&templates).Error
	return templates, err
}
