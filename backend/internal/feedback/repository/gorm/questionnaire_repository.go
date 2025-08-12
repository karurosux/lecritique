package gorm

import (
	"context"
	"github.com/google/uuid"
	feedbackmodel "kyooar/internal/feedback/model"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type questionnaireRepository struct {
	*sharedRepos.BaseRepository[feedbackmodel.Questionnaire]
}

func NewQuestionnaireRepository(db *gorm.DB) *questionnaireRepository {
	return &questionnaireRepository{
		BaseRepository: sharedRepos.NewBaseRepository[feedbackmodel.Questionnaire](db),
	}
}

func (r *questionnaireRepository) Create(ctx context.Context, questionnaire *feedbackmodel.Questionnaire) error {
	return r.DB.WithContext(ctx).Create(questionnaire).Error
}

func (r *questionnaireRepository) FindByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Questionnaire, error) {
	var questionnaire feedbackmodel.Questionnaire
	err := r.DB.WithContext(ctx).First(&questionnaire, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) Update(ctx context.Context, questionnaire *feedbackmodel.Questionnaire) error {
	return r.DB.WithContext(ctx).Save(questionnaire).Error
}

func (r *questionnaireRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&feedbackmodel.Questionnaire{}, "id = ?", id).Error
}

func (r *questionnaireRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*feedbackmodel.Questionnaire, error) {
	var questionnaires []*feedbackmodel.Questionnaire
	err := r.DB.WithContext(ctx).Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Find(&questionnaires).Error
	return questionnaires, err
}

func (r *questionnaireRepository) DeactivateDefaultQuestionnaires(ctx context.Context, organizationID uuid.UUID) error {
	return r.DB.WithContext(ctx).Model(&feedbackmodel.Questionnaire{}).
		Where("organization_id = ? AND is_default = ?", organizationID, true).
		Update("is_default", false).Error
}

func (r *questionnaireRepository) FindByIDWithQuestions(ctx context.Context, id uuid.UUID) (*feedbackmodel.Questionnaire, error) {
	var questionnaire feedbackmodel.Questionnaire
	err := r.DB.WithContext(ctx).Preload("Questions", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).First(&questionnaire, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) FindByProductID(ctx context.Context, productID uuid.UUID) (*feedbackmodel.Questionnaire, error) {
	var questionnaire feedbackmodel.Questionnaire
	err := r.DB.WithContext(ctx).Preload("Questions").
		Where("product_id = ? AND is_active = ?", productID, true).
		First(&questionnaire).Error
	if err != nil {
		return nil, err
	}
	return &questionnaire, nil
}

func (r *questionnaireRepository) CreateQuestion(ctx context.Context, question *feedbackmodel.Question) error {
	return r.DB.WithContext(ctx).Create(question).Error
}

func (r *questionnaireRepository) FindQuestionByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Question, error) {
	var question feedbackmodel.Question
	err := r.DB.WithContext(ctx).First(&question, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionnaireRepository) UpdateQuestion(ctx context.Context, question *feedbackmodel.Question) error {
	return r.DB.WithContext(ctx).Save(question).Error
}

func (r *questionnaireRepository) DeleteQuestion(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&feedbackmodel.Question{}, "id = ?", id).Error
}

func (r *questionnaireRepository) GetMaxQuestionOrder(ctx context.Context, questionnaireID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.DB.WithContext(ctx).
		Model(&feedbackmodel.Question{}).
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
		if err := tx.Model(&feedbackmodel.Question{}).
			Where("id = ? AND questionnaire_id = ?", questionID, questionnaireID).
			Update("display_order", i+1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}