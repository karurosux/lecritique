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
	FindByID(ctx context.Context, id uuid.UUID) (*models.Questionnaire, error)
	FindByIDWithQuestions(ctx context.Context, id uuid.UUID) (*models.Questionnaire, error)
	FindByDishID(ctx context.Context, dishID uuid.UUID) (*models.Questionnaire, error)
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Questionnaire, error)
	Update(ctx context.Context, questionnaire *models.Questionnaire) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeactivateDefaultQuestionnaires(ctx context.Context, restaurantID uuid.UUID) error
	
	// Question methods
	CreateQuestion(ctx context.Context, question *models.Question) error
	FindQuestionByID(ctx context.Context, id uuid.UUID) (*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error
	GetMaxQuestionOrder(ctx context.Context, questionnaireID uuid.UUID) (int, error)
	ReorderQuestions(ctx context.Context, questionnaireID uuid.UUID, questionIDs []uuid.UUID) error
}

type questionnaireRepository struct {
	*sharedRepos.BaseRepository[models.Questionnaire]
}

func NewQuestionnaireRepository(db *gorm.DB) QuestionnaireRepository {
	return &questionnaireRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Questionnaire](db),
	}
}

func (r *questionnaireRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	err := r.DB.WithContext(ctx).First(&questionnaire, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) Create(ctx context.Context, questionnaire *models.Questionnaire) error {
	return r.DB.WithContext(ctx).Create(questionnaire).Error
}

func (r *questionnaireRepository) Update(ctx context.Context, questionnaire *models.Questionnaire) error {
	return r.DB.WithContext(ctx).Save(questionnaire).Error
}

func (r *questionnaireRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.Questionnaire{}, "id = ?", id).Error
}

func (r *questionnaireRepository) FindByIDWithQuestions(ctx context.Context, id uuid.UUID) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	err := r.DB.WithContext(ctx).Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).First(&questionnaire, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
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

func (r *questionnaireRepository) DeactivateDefaultQuestionnaires(ctx context.Context, restaurantID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Questionnaire{}).
		Where("restaurant_id = ? AND is_default = ?", restaurantID, true).
		Update("is_default", false).Error
}

func (r *questionnaireRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	return r.DB.WithContext(ctx).Create(question).Error
}

func (r *questionnaireRepository) FindQuestionByID(ctx context.Context, id uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := r.DB.WithContext(ctx).First(&question, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionnaireRepository) UpdateQuestion(ctx context.Context, question *models.Question) error {
	return r.DB.WithContext(ctx).Save(question).Error
}

func (r *questionnaireRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&models.Question{}, "id = ?", id).Error
}

func (r *questionnaireRepository) GetMaxQuestionOrder(ctx context.Context, questionnaireID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.DB.WithContext(ctx).Model(&models.Question{}).
		Where("questionnaire_id = ?", questionnaireID).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *questionnaireRepository) ReorderQuestions(ctx context.Context, questionnaireID uuid.UUID, questionIDs []uuid.UUID) error {
	tx := r.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	for i, questionID := range questionIDs {
		if err := tx.Model(&models.Question{}).
			Where("id = ? AND questionnaire_id = ?", questionID, questionnaireID).
			Update("display_order", i+1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	
	return tx.Commit().Error
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
