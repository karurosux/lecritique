package service

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	feedbackinterface "kyooar/internal/feedback/interface"
	feedbackmodel "kyooar/internal/feedback/model"
	menuRepos "kyooar/internal/product/repositories"
	organizationinterface "kyooar/internal/organization/interface"
)

type questionService struct {
	questionRepo     feedbackinterface.QuestionRepository
	productRepo      menuRepos.ProductRepository
	organizationRepo organizationinterface.OrganizationRepository
}

func NewQuestionService(
	questionRepo feedbackinterface.QuestionRepository,
	productRepo menuRepos.ProductRepository,
	organizationRepo organizationinterface.OrganizationRepository,
) feedbackinterface.QuestionService {
	return &questionService{
		questionRepo:     questionRepo,
		productRepo:      productRepo,
		organizationRepo: organizationRepo,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, accountID, productID uuid.UUID, request *feedbackmodel.CreateQuestionRequest) (*feedbackmodel.Question, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	question := &feedbackmodel.Question{
		ProductID:    productID,
		Text:         request.Text,
		Type:         request.Type,
		IsRequired:   request.IsRequired,
		Options:      request.Options,
		MinValue:     request.MinValue,
		MaxValue:     request.MaxValue,
		MinLabel:     request.MinLabel,
		MaxLabel:     request.MaxLabel,
		DisplayOrder: 0,
	}

	if err := s.questionRepo.Create(ctx, question); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to create question")
	}

	return question, nil
}

func (s *questionService) GetQuestionsByProduct(ctx context.Context, accountID, productID uuid.UUID) ([]*feedbackmodel.Question, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return s.questionRepo.FindByProductID(ctx, productID)
}

func (s *questionService) GetQuestionsByProducts(ctx context.Context, accountID, organizationID uuid.UUID, productIDs []uuid.UUID) ([]*feedbackmodel.Question, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return s.questionRepo.FindByProductIDs(ctx, productIDs)
}

func (s *questionService) GetQuestionsByProductsOptimized(ctx context.Context, accountID, organizationID uuid.UUID, productIDs []uuid.UUID) ([]*feedbackmodel.BatchQuestionResponse, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return s.questionRepo.FindByProductIDsOptimized(ctx, productIDs)
}

func (s *questionService) GetQuestion(ctx context.Context, accountID, questionID uuid.UUID) (*feedbackmodel.Question, error) {
	question, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Question not found")
	}

	product, err := s.productRepo.FindByID(ctx, question.ProductID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return question, nil
}

func (s *questionService) UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, request *feedbackmodel.UpdateQuestionRequest) (*feedbackmodel.Question, error) {
	question, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Question not found")
	}

	product, err := s.productRepo.FindByID(ctx, question.ProductID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	question.Text = request.Text
	question.Type = request.Type
	question.IsRequired = request.IsRequired
	question.Options = request.Options
	question.MinValue = request.MinValue
	question.MaxValue = request.MaxValue
	question.MinLabel = request.MinLabel
	question.MaxLabel = request.MaxLabel

	if err := s.questionRepo.Update(ctx, question); err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to update question")
	}

	return question, nil
}

func (s *questionService) DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error {
	question, err := s.questionRepo.FindByID(ctx, questionID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Question not found")
	}

	product, err := s.productRepo.FindByID(ctx, question.ProductID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return s.questionRepo.Delete(ctx, questionID)
}

func (s *questionService) ReorderQuestions(ctx context.Context, accountID, productID uuid.UUID, questionIDs []uuid.UUID) error {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Product not found")
	}

	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return s.questionRepo.ReorderQuestions(ctx, productID, questionIDs)
}

func (s *questionService) GetProductsWithQuestions(ctx context.Context, accountID, organizationID uuid.UUID) ([]uuid.UUID, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Organization not found")
	}

	if organization.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusForbidden, "Access denied")
	}

	return s.questionRepo.GetProductsWithQuestions(ctx, organizationID)
}

func (s *questionService) GetQuestionsByProductIDForAnalytics(ctx context.Context, productID uuid.UUID) ([]*feedbackmodel.Question, error) {
	return s.questionRepo.GetQuestionsByProductIDForAnalytics(ctx, productID)
}