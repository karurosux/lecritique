package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/feedback/models"
	"github.com/lecritique/api/internal/feedback/repositories"
	menuRepos "github.com/lecritique/api/internal/menu/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, accountID, dishID uuid.UUID, request *models.CreateQuestionRequest) (*models.Question, error)
	GetQuestionsByDish(ctx context.Context, accountID, dishID uuid.UUID) ([]*models.Question, error)
	GetQuestion(ctx context.Context, accountID, questionID uuid.UUID) (*models.Question, error)
	UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, request *models.UpdateQuestionRequest) (*models.Question, error)
	DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error
	ReorderQuestions(ctx context.Context, accountID, dishID uuid.UUID, questionIDs []uuid.UUID) error
	GetDishesWithQuestions(ctx context.Context, accountID, restaurantID uuid.UUID) ([]uuid.UUID, error)
}

type questionService struct {
	questionRepo   repositories.QuestionRepository
	dishRepo       menuRepos.DishRepository
	restaurantRepo restaurantRepos.RestaurantRepository
}

func NewQuestionService(questionRepo repositories.QuestionRepository, dishRepo menuRepos.DishRepository, restaurantRepo restaurantRepos.RestaurantRepository) QuestionService {
	return &questionService{
		questionRepo:   questionRepo,
		dishRepo:       dishRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, accountID, dishID uuid.UUID, request *models.CreateQuestionRequest) (*models.Question, error) {
	// Verify dish exists and belongs to account
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}

	if restaurant.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	// Get next display order
	maxOrder, err := s.questionRepo.GetMaxDisplayOrder(ctx, dishID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to get display order")
	}

	// Create question
	question := &models.Question{
		DishID:       dishID,
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

func (s *questionService) GetQuestionsByDish(ctx context.Context, accountID, dishID uuid.UUID) ([]*models.Question, error) {
	// Verify dish exists and belongs to account
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}

	if restaurant.AccountID != accountID {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	questions, err := s.questionRepo.GetQuestionsByDishID(ctx, dishID)
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

	// Verify the question's dish belongs to the account
	dish, err := s.dishRepo.FindByID(ctx, question.DishID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Associated dish not found")
	}

	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}

	if restaurant.AccountID != accountID {
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

func (s *questionService) ReorderQuestions(ctx context.Context, accountID, dishID uuid.UUID, questionIDs []uuid.UUID) error {
	// Verify dish belongs to account
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}

	if restaurant.AccountID != accountID {
		return echo.NewHTTPError(http.StatusNotFound, "Dish not found")
	}

	// Verify all questions belong to this dish
	for _, questionID := range questionIDs {
		question, err := s.questionRepo.GetQuestionByID(ctx, questionID)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Question not found")
		}
		if question.DishID != dishID {
			return echo.NewHTTPError(http.StatusBadRequest, "Question does not belong to this dish")
		}
	}

	if err := s.questionRepo.ReorderQuestions(ctx, dishID, questionIDs); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to reorder questions")
	}

	return nil
}

func (s *questionService) GetDishesWithQuestions(ctx context.Context, accountID, restaurantID uuid.UUID) ([]uuid.UUID, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		// Log for debugging
		fmt.Printf("DEBUG: Restaurant not found. RestaurantID: %s, Error: %v\n", restaurantID, err)
		return nil, echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}

	if restaurant.AccountID != accountID {
		// Log for debugging
		fmt.Printf("DEBUG: Account mismatch. RestaurantAccountID: %s, RequestAccountID: %s\n", restaurant.AccountID, accountID)
		return nil, echo.NewHTTPError(http.StatusNotFound, "Restaurant not found")
	}

	return s.questionRepo.GetDishesWithQuestions(ctx, restaurantID)
}