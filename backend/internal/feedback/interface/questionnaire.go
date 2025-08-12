package feedbackinterface

import (
	"context"

	"github.com/google/uuid"
	feedbackmodel "kyooar/internal/feedback/model"
	menuModels "kyooar/internal/menu/models"
)

type QuestionnaireRepository interface {
	Create(ctx context.Context, questionnaire *feedbackmodel.Questionnaire) error
	FindByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Questionnaire, error)
	Update(ctx context.Context, questionnaire *feedbackmodel.Questionnaire) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]*feedbackmodel.Questionnaire, error)
	DeactivateDefaultQuestionnaires(ctx context.Context, organizationID uuid.UUID) error
	CreateQuestion(ctx context.Context, question *feedbackmodel.Question) error
	FindQuestionByID(ctx context.Context, id uuid.UUID) (*feedbackmodel.Question, error)
	UpdateQuestion(ctx context.Context, question *feedbackmodel.Question) error
	DeleteQuestion(ctx context.Context, id uuid.UUID) error
	ReorderQuestions(ctx context.Context, questionnaireID uuid.UUID, questionIDs []uuid.UUID) error
}

type QuestionnaireService interface {
	Create(ctx context.Context, accountID, organizationID uuid.UUID, input *feedbackmodel.Questionnaire) (*feedbackmodel.Questionnaire, error)
	GetByID(ctx context.Context, accountID, questionnaireID uuid.UUID) (*feedbackmodel.Questionnaire, error)
	ListByOrganization(ctx context.Context, accountID, organizationID uuid.UUID) ([]*feedbackmodel.Questionnaire, error)
	Update(ctx context.Context, accountID, questionnaireID uuid.UUID, input *feedbackmodel.Questionnaire) (*feedbackmodel.Questionnaire, error)
	Delete(ctx context.Context, accountID, questionnaireID uuid.UUID) error
	AddQuestion(ctx context.Context, accountID, questionnaireID uuid.UUID, question *feedbackmodel.Question) (*feedbackmodel.Question, error)
	UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, question *feedbackmodel.Question) (*feedbackmodel.Question, error)
	DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error
	ReorderQuestions(ctx context.Context, accountID, questionnaireID uuid.UUID, questionIDs []uuid.UUID) error
	GenerateQuestionsForProduct(ctx context.Context, accountID uuid.UUID, product *menuModels.Product) ([]*feedbackmodel.GeneratedQuestion, error)
	GenerateAndSaveQuestionnaireForProduct(ctx context.Context, accountID uuid.UUID, product *menuModels.Product, name, description string, isDefault bool) (*feedbackmodel.Questionnaire, error)
}