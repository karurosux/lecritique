package feedbackinterface

import (
	"context"

	"github.com/google/uuid"
	feedbackmodel "kyooar/internal/feedback/model"
)

type QuestionRepository interface {
	Create(ctx context.Context, question *feedbackmodel.Question) error
	FindByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Question, error)
	Update(ctx context.Context, question *feedbackmodel.Question) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByProductID(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error)
	FindByProductIDs(ctx context.Context, productIDs []uuid.UUID) ([]*feedbackmodel.Question, error)
	FindByProductIDsOptimized(ctx context.Context, productIDs []uuid.UUID) ([]*feedbackmodel.BatchQuestionResponse, error)
	ReorderQuestions(ctx context.Context, productID uuid.UUID, questionIDs []uuid.UUID) error
	GetProductsWithQuestions(ctx context.Context, organizationID uuid.UUID) ([]uuid.UUID, error)
	GetQuestionsByProductID(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error)
	GetQuestionsByProductIDForAnalytics(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error)
}

type QuestionService interface {
	CreateQuestion(ctx context.Context, accountID, productID uuid.UUID, request *feedbackmodel.CreateQuestionRequest) (*feedbackmodel.Question, error)
	GetQuestionsByProduct(ctx context.Context, accountID, productID uuid.UUID) ([]*feedbackmodel.Question, error)
	GetQuestionsByProducts(ctx context.Context, accountID, organizationID uuid.UUID, productIDs []uuid.UUID) ([]*feedbackmodel.Question, error)
	GetQuestionsByProductsOptimized(ctx context.Context, accountID, organizationID uuid.UUID, productIDs []uuid.UUID) ([]*feedbackmodel.BatchQuestionResponse, error)
	GetQuestion(ctx context.Context, accountID, questionID uuid.UUID) (*feedbackmodel.Question, error)
	UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, request *feedbackmodel.UpdateQuestionRequest) (*feedbackmodel.Question, error)
	DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error
	ReorderQuestions(ctx context.Context, accountID, productID uuid.UUID, questionIDs []uuid.UUID) error
	GetProductsWithQuestions(ctx context.Context, accountID, organizationID uuid.UUID) ([]uuid.UUID, error)
	GetQuestionsByProductIDForAnalytics(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error)
}