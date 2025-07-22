package repositories

import (
	"context"
	"github.com/google/uuid"
	"lecritique/internal/feedback/models"
	sharedRepos "lecritique/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type QuestionnaireRepository interface {
	Create(ctx context.Context, questionnaire *models.Questionnaire) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.Questionnaire, error)
	FindByIDWithQuestions(ctx context.Context, id uuid.UUID) (*models.Questionnaire, error)
	FindByProductID(ctx context.Context, productID uuid.UUID) (*models.Questionnaire, error)
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.Questionnaire, error)
	Update(ctx context.Context, questionnaire *models.Questionnaire) error
	Delete(ctx context.Context, id uuid.UUID) error
	DeactivateDefaultQuestionnaires(ctx context.Context, organizationID uuid.UUID) error
	
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

func NewQuestionnaireRepository(i *do.Injector) (QuestionnaireRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &questionnaireRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Questionnaire](db),
	}, nil
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

func (r *questionnaireRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*models.Questionnaire, error) {
	var questionnaire models.Questionnaire
	err := r.DB.WithContext(ctx).Preload("Questions").
		Where("product_id = ? AND is_active = ?", productID, true).
		First(&questionnaire).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.Questionnaire, error) {
	var questionnaires []models.Questionnaire
	err := r.DB.WithContext(ctx).Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Find(&questionnaires).Error
	return questionnaires, err
}

func (r *questionnaireRepository) DeactivateDefaultQuestionnaires(ctx context.Context, organizationID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&models.Questionnaire{}).
		Where("organization_id = ? AND is_default = ?", organizationID, true).
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

func NewQuestionTemplateRepository(i *do.Injector) (QuestionTemplateRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &questionTemplateRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.QuestionTemplate](db),
	}, nil
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
