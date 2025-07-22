package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/feedback/models"
	"lecritique/internal/feedback/repositories"
	menuRepos "lecritique/internal/menu/repositories"
	organizationRepos "lecritique/internal/organization/repositories"
	"github.com/samber/do"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, accountID, productID uuid.UUID, request *models.CreateQuestionRequest) (*models.Question, error)
	GetQuestionsByProduct(ctx context.Context, accountID, productID uuid.UUID) ([]*models.Question, error)
	GetQuestion(ctx context.Context, accountID, questionID uuid.UUID) (*models.Question, error)
	UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, request *models.UpdateQuestionRequest) (*models.Question, error)
	DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error
	ReorderQuestions(ctx context.Context, accountID, productID uuid.UUID, questionIDs []uuid.UUID) error
	GetProductesWithQuestions(ctx context.Context, accountID, organizationID uuid.UUID) ([]uuid.UUID, error)
}

type questionService struct {
	questionRepo   repositories.QuestionRepository
	productRepo       menuRepos.ProductRepository
	organizationRepo organizationRepos.OrganizationRepository
}

func NewQuestionService(i *do.Injector) (QuestionService, error) {
	return &questionService{
		questionRepo:   do.MustInvoke[repositories.QuestionRepository](i),
		productRepo:       do.MustInvoke[menuRepos.ProductRepository](i),
		organizationRepo: do.MustInvoke[organizationRepos.OrganizationRepository](i),
	}, nil
}

func (s *questionService) CreateQuestion(ctx context.Context, accountID, productID uuid.UUID, request *models.CreateQuestionRequest) (*models.Question, error) {
	// Verify product exists and belongs to account
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Get next display order
	maxOrder, err := s.questionRepo.GetMaxDisplayOrder(ctx, productID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get display order")
	}

	// Create question
	question := &models.Question{
		ProductID:       productID,
		Text:         request.Text,
		Type:         request.Type,
		IsRequired:   request.IsRequired,
		DisplayOrder: maxOrder + 1,
		Options:      request.Options,
		MinValue:     request.MinValue,
		MaxValue:     request.MaxValue,
		MinLabel:     request.MinLabel,
		MaxLabel:     request.MaxLabel,
	}

	if err := s.questionRepo.CreateQuestion(ctx, question); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to create question")
	}

	return question, nil
}

func (s *questionService) GetQuestionsByProduct(ctx context.Context, accountID, productID uuid.UUID) ([]*models.Question, error) {
	// Verify product exists and belongs to account
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	questions, err := s.questionRepo.GetQuestionsByProductID(ctx, productID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get questions")
	}

	return questions, nil
}

func (s *questionService) GetQuestion(ctx context.Context, accountID, questionID uuid.UUID) (*models.Question, error) {
	question, err := s.questionRepo.GetQuestionByID(ctx, questionID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Question not found")
	}

	// Verify the question's product belongs to the account
	product, err := s.productRepo.FindByID(ctx, question.ProductID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Associated product not found")
	}

	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Question not found")
	}

	return question, nil
}

func (s *questionService) UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, request *models.UpdateQuestionRequest) (*models.Question, error) {
	// Get existing question and verify access
	question, err := s.GetQuestion(ctx, accountID, questionID)
	if err != nil {
		return nil, err
	}

	// Update fields
	question.Text = request.Text
	question.Type = request.Type
	question.IsRequired = request.IsRequired
	question.Options = request.Options
	question.MinValue = request.MinValue
	question.MaxValue = request.MaxValue
	question.MinLabel = request.MinLabel
	question.MaxLabel = request.MaxLabel

	if err := s.questionRepo.UpdateQuestion(ctx, question); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to update question")
	}

	return question, nil
}

func (s *questionService) DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error {
	// Verify access first
	_, err := s.GetQuestion(ctx, accountID, questionID)
	if err != nil {
		return err
	}

	if err := s.questionRepo.DeleteQuestion(ctx, questionID); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete question")
	}

	return nil
}

func (s *questionService) ReorderQuestions(ctx context.Context, accountID, productID uuid.UUID, questionIDs []uuid.UUID) error {
	// Verify product belongs to account
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	// Verify all questions belong to this product
	for _, questionID := range questionIDs {
		question, err := s.questionRepo.GetQuestionByID(ctx, questionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Question not found")
		}
		if question.ProductID != productID {
			return echo.NewHTTPError(http.StatusBadRequest, "Question does not belong to this product")
		}
	}

	if err := s.questionRepo.ReorderQuestions(ctx, productID, questionIDs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reorder questions")
	}

	return nil
}

func (s *questionService) GetProductesWithQuestions(ctx context.Context, accountID, organizationID uuid.UUID) ([]uuid.UUID, error) {
	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		// Log for debugging
		fmt.Printf("DEBUG: Organization not found. OrganizationID: %s, Error: %v\n", organizationID, err)
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		// Log for debugging
		fmt.Printf("DEBUG: Account mismatch. OrganizationAccountID: %s, RequestAccountID: %s\n", organization.AccountID, accountID)
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	return s.questionRepo.GetProductesWithQuestions(ctx, organizationID)
}
