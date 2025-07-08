package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	aiServices "github.com/lecritique/api/internal/ai/services"
	"github.com/lecritique/api/internal/feedback/models"
	"github.com/lecritique/api/internal/feedback/repositories"
	menuModels "github.com/lecritique/api/internal/menu/models"
	"github.com/lecritique/api/internal/shared/config"
	"gorm.io/gorm"
)

type QuestionnaireService struct {
	repo              repositories.QuestionnaireRepository
	questionGenerator *aiServices.QuestionGenerator
}

func NewQuestionnaireService(db *gorm.DB, cfg *config.Config) (*QuestionnaireService, error) {
	generator, err := aiServices.NewQuestionGenerator(cfg)
	if err != nil {
		// AI is optional, so we continue without it
		generator = nil
	}

	return &QuestionnaireService{
		repo:              repositories.NewQuestionnaireRepository(db),
		questionGenerator: generator,
	}, nil
}

func (s *QuestionnaireService) Create(ctx context.Context, accountID, restaurantID uuid.UUID, input *models.Questionnaire) (*models.Questionnaire, error) {
	// Validate restaurant belongs to account
	// TODO: Add validation

	questionnaire := &models.Questionnaire{
		RestaurantID: restaurantID,
		DishID:       input.DishID,
		Name:         input.Name,
		Description:  input.Description,
		IsDefault:    input.IsDefault,
		IsActive:     true,
	}

	// If it's a default questionnaire, deactivate other defaults
	if questionnaire.IsDefault {
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, restaurantID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, questionnaire); err != nil {
		return nil, err
	}

	return questionnaire, nil
}

func (s *QuestionnaireService) Update(ctx context.Context, accountID, questionnaireID uuid.UUID, input *models.Questionnaire) (*models.Questionnaire, error) {
	questionnaire, err := s.repo.FindByID(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate questionnaire belongs to account

	questionnaire.Name = input.Name
	questionnaire.Description = input.Description
	questionnaire.IsActive = input.IsActive

	if input.IsDefault && !questionnaire.IsDefault {
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, questionnaire.RestaurantID); err != nil {
			return nil, err
		}
		questionnaire.IsDefault = true
	}

	if err := s.repo.Update(ctx, questionnaire); err != nil {
		return nil, err
	}

	return questionnaire, nil
}

func (s *QuestionnaireService) Delete(ctx context.Context, accountID, questionnaireID uuid.UUID) error {
	questionnaire, err := s.repo.FindByID(ctx, questionnaireID)
	if err != nil {
		return err
	}

	// TODO: Validate questionnaire belongs to account

	if questionnaire.IsDefault {
		return fmt.Errorf("cannot delete default questionnaire")
	}

	return s.repo.Delete(ctx, questionnaireID)
}

func (s *QuestionnaireService) GetByID(ctx context.Context, accountID, questionnaireID uuid.UUID) (*models.Questionnaire, error) {
	questionnaire, err := s.repo.FindByIDWithQuestions(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate questionnaire belongs to account

	return questionnaire, nil
}

func (s *QuestionnaireService) ListByRestaurant(ctx context.Context, accountID, restaurantID uuid.UUID) ([]models.Questionnaire, error) {
	// TODO: Validate restaurant belongs to account

	return s.repo.FindByRestaurantID(ctx, restaurantID)
}

func (s *QuestionnaireService) AddQuestion(ctx context.Context, accountID, questionnaireID uuid.UUID, question *models.Question) (*models.Question, error) {
	_, err := s.repo.FindByID(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate questionnaire belongs to account

	// Get the next display order
	maxOrder, err := s.repo.GetMaxQuestionOrder(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	// TODO: Update this to use new Question structure
	// question.QuestionnaireID = questionnaireID
	question.DisplayOrder = maxOrder + 1

	if err := s.repo.CreateQuestion(ctx, question); err != nil {
		return nil, err
	}

	return question, nil
}

func (s *QuestionnaireService) UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, input *models.Question) (*models.Question, error) {
	question, err := s.repo.FindQuestionByID(ctx, questionID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate question's questionnaire belongs to account

	question.Text = input.Text
	question.Type = input.Type
	question.IsRequired = input.IsRequired
	question.Options = input.Options
	question.MinValue = input.MinValue
	question.MaxValue = input.MaxValue
	question.MinLabel = input.MinLabel
	question.MaxLabel = input.MaxLabel

	if err := s.repo.UpdateQuestion(ctx, question); err != nil {
		return nil, err
	}

	return question, nil
}

func (s *QuestionnaireService) DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error {
	_, err := s.repo.FindQuestionByID(ctx, questionID)
	if err != nil {
		return err
	}

	// TODO: Validate question's questionnaire belongs to account

	return s.repo.DeleteQuestion(ctx, questionID)
}

func (s *QuestionnaireService) ReorderQuestions(ctx context.Context, accountID, questionnaireID uuid.UUID, questionIDs []uuid.UUID) error {
	// TODO: Validate questionnaire belongs to account

	return s.repo.ReorderQuestions(ctx, questionnaireID, questionIDs)
}

// GenerateQuestionsForDish uses AI to generate questions based on dish details
func (s *QuestionnaireService) GenerateQuestionsForDish(ctx context.Context, accountID uuid.UUID, dish *menuModels.Dish) ([]aiServices.GeneratedQuestion, error) {
	if s.questionGenerator == nil {
		return nil, fmt.Errorf("AI question generation is not configured")
	}

	// TODO: Validate dish belongs to account

	questions, err := s.questionGenerator.GenerateQuestionsForDish(ctx, dish)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	return questions, nil
}

// GenerateAndSaveQuestionnaireForDish generates AI questions and creates a complete questionnaire for a dish
func (s *QuestionnaireService) GenerateAndSaveQuestionnaireForDish(ctx context.Context, accountID uuid.UUID, dish *menuModels.Dish, name, description string, isDefault bool) (*models.Questionnaire, error) {
	if s.questionGenerator == nil {
		return nil, fmt.Errorf("AI question generation is not configured")
	}

	// Generate questions using AI
	generatedQuestions, err := s.questionGenerator.GenerateQuestionsForDish(ctx, dish)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	// Create the questionnaire
	questionnaire := &models.Questionnaire{
		RestaurantID: dish.RestaurantID,
		DishID:       &dish.ID,
		Name:         name,
		Description:  description,
		IsDefault:    isDefault,
		IsActive:     true,
	}

	// If it's a default questionnaire, deactivate other defaults
	if questionnaire.IsDefault {
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, dish.RestaurantID); err != nil {
			return nil, err
		}
	}

	// Create the questionnaire
	if err := s.repo.Create(ctx, questionnaire); err != nil {
		return nil, err
	}

	// Convert generated questions to database questions and add them
	for i, genQ := range generatedQuestions {
		question := &models.Question{
			// TODO: Update this to use new Question structure with DishID
			// QuestionnaireID: questionnaire.ID,
			Text:         genQ.Text,
			Type:         genQ.Type,
			IsRequired:   true, // Default to required
			DisplayOrder: i + 1,
			Options:      genQ.Options,
			MinValue:     genQ.MinValue,
			MaxValue:     genQ.MaxValue,
			MinLabel:     genQ.MinLabel,
			MaxLabel:     genQ.MaxLabel,
		}

		if err := s.repo.CreateQuestion(ctx, question); err != nil {
			return nil, fmt.Errorf("failed to create question %d: %w", i+1, err)
		}
	}

	// Reload the questionnaire with questions
	return s.repo.FindByIDWithQuestions(ctx, questionnaire.ID)
}