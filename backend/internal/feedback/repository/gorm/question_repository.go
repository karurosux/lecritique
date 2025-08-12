package gorm

import (
	"context"

	"github.com/google/uuid"
	feedbackmodel "kyooar/internal/feedback/model"
	"gorm.io/gorm"
)

type questionRepository struct {
	db *gorm.DB
}

func NewQuestionRepository(db *gorm.DB) *questionRepository {
	return &questionRepository{db: db}
}

func (r *questionRepository) Create(ctx context.Context, question *feedbackmodel.Question) error {
	return r.db.WithContext(ctx).Create(question).Error
}

func (r *questionRepository) FindByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Question, error) {
	var question feedbackmodel.Question
	err := r.db.WithContext(ctx).
		Preload("Product").
		Where("id = ?", id).
		First(&question).Error
	if err != nil {
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) Update(ctx context.Context, question *feedbackmodel.Question) error {
	return r.db.WithContext(ctx).Save(question).Error
}

func (r *questionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&feedbackmodel.Question{}, "id = ?", id).Error
}

func (r *questionRepository) FindByProductID(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error) {
	var questions []*feedbackmodel.Question
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) FindByProductIDs(ctx context.Context, productIDs []uuid.UUID) ([]*feedbackmodel.Question, error) {
	if len(productIDs) == 0 {
		return []*feedbackmodel.Question{}, nil
	}

	var questions []*feedbackmodel.Question
	err := r.db.WithContext(ctx).
		Where("product_id IN ?", productIDs).
		Order("product_id ASC, display_order ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) FindByProductIDsOptimized(ctx context.Context, productIDs []uuid.UUID) ([]*feedbackmodel.BatchQuestionResponse, error) {
	if len(productIDs) == 0 {
		return []*feedbackmodel.BatchQuestionResponse{}, nil
	}

	var results []*feedbackmodel.BatchQuestionResponse
	err := r.db.WithContext(ctx).
		Model(&feedbackmodel.Question{}).
		Select("id, product_id, text, type").
		Where("product_id IN ?", productIDs).
		Order("product_id ASC, display_order ASC").
		Find(&results).Error
	return results, err
}

func (r *questionRepository) ReorderQuestions(ctx context.Context, productID uuid.UUID, questionIDs []uuid.UUID) error {
	tx := r.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for i, questionID := range questionIDs {
		if err := tx.Model(&feedbackmodel.Question{}).
			Where("id = ? AND product_id = ?", questionID, productID).
			Update("display_order", i+1).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (r *questionRepository) GetProductsWithQuestions(ctx context.Context, organizationID uuid.UUID) ([]uuid.UUID, error) {
	var productIDs []uuid.UUID
	
	err := r.db.WithContext(ctx).
		Model(&feedbackmodel.Question{}).
		Select("DISTINCT questions.product_id").
		Joins("JOIN products ON products.id = questions.product_id").
		Where("products.organization_id = ?", organizationID).
		Pluck("questions.product_id", &productIDs).Error
	
	return productIDs, err
}

func (r *questionRepository) GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error) {
	var questions []*feedbackmodel.Question
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetQuestionsByProductIDForAnalytics(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error) {
	var questions []*feedbackmodel.Question
	err := r.db.WithContext(ctx).
		Where("product_id = ?", productID).
		Order("display_order ASC").
		Find(&questions).Error
	return questions, err
}

func (r *questionRepository) GetMaxDisplayOrder(ctx context.Context, productID uuid.UUID) (int, error) {
	var maxOrder int
	err := r.db.WithContext(ctx).
		Model(&feedbackmodel.Question{}).
		Where("product_id = ?", productID).
		Select("COALESCE(MAX(display_order), 0)").
		Scan(&maxOrder).Error
	return maxOrder, err
}