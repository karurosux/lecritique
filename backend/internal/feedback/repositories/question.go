package repositories

import (
	"context"

	"github.com/google/uuid"
	"lecritique/internal/feedback/models"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type QuestionRepository interface {
	CreateQuestion(ctx context.Context, question *models.Question) error
	GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]*models.Question, error)
	GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*models.Question, error)
	UpdateQuestion(ctx context.Context, question *models.Question) error
	DeleteQuestion(ctx context.Context, questionID uuid.UUID) error
	ReorderQuestions(ctx context.Context, productID uuid.UUID, questionIDs []uuid.UUID) error
	GetMaxDisplayOrder(ctx context.Context, productID uuid.UUID) (int, error)
	GetProductesWithQuestions(ctx context.Context, organizationID uuid.UUID) ([]uuid.UUID, error)
}

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(i *do.Injector) (QuestionRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &questionRepository{db: db}, nil
}

func (r *questionRepository) CreateQuestion(ctx context.Context, question *models.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepository) GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]*models.Question, error) {
	var questions []*models.Question
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetQuestionByID(ctx context.Context, questionID uuid.UUID) (*models.Question, error) {
	var question models.Question
	err := r.db.WithContext(ctx).
		Preload("Product").
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

func (r *questionRepository) GetMaxDisplayOrder(ctx context.Context, productID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.WithContext(ctx).
		Model(&models.Question{}).
		Where("product_id = ?", productID).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}

func (r *questionRepository) ReorderQuestions(ctx context.Context, productID uuid.UUID, questionIDs []uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update display order for each question
	for i, questionID := range questionIDs {
		if err := tx.Model(&models.Question{}).
			Where("id = ? AND product_id = ?", questionID, productID).
			Update("display_order", i+1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *questionRepository) GetProductesWithQuestions(ctx context.Context, organizationID uuid.UUID) ([]uuid.UUID, error) {
	var productIDs []uuid.UUID
	err := r.db.WithContext(ctx).
		Table("questions").
		Select("DISTINCT questions.product_id").
		Joins("JOIN products ON products.id = questions.product_id").
		Where("products.organization_id = ?", organizationID).
		Scan(&productIDs).Error
	return productIDs, err
}
