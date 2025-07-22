package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	aiServices "kyooar/internal/ai/services"
	"kyooar/internal/feedback/models"
	"kyooar/internal/feedback/repositories"
	menuModels "kyooar/internal/menu/models"
	"kyooar/internal/shared/config"
	"github.com/samber/do"
)

type QuestionnaireService struct {
	repo              repositories.QuestionnaireRepository
	questionGenerator *aiServices.QuestionGenerator
}

func NewQuestionnaireService(i *do.Injector) (*QuestionnaireService, error) {
	cfg := do.MustInvoke[*config.Config](i)
	repo := do.MustInvoke[repositories.QuestionnaireRepository](i)
	
	generator, err := aiServices.NewQuestionGenerator(cfg)
	if err != nil {
		// AI is optional, so we continue without it
		generator = nil
	}

	return &QuestionnaireService{
		repo:              repo,
		questionGenerator: generator,
	}, nil
}

func (s *QuestionnaireService) Create(ctx context.Context, accountID, organizationID uuid.UUID, input *models.Questionnaire) (*models.Questionnaire, error) {
	// Validate organization belongs to account
	// TODO: Add validation

	questionnaire := &models.Questionnaire{
		OrganizationID: organizationID,
		ProductID:       input.ProductID,
		Name:         input.Name,
		Description:  input.Description,
		IsDefault:    input.IsDefault,
		IsActive:     true,
	}

	if questionnaire.IsDefault {
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, organizationID); err != nil {
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
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, questionnaire.OrganizationID); err != nil {
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

func (s *QuestionnaireService) ListByOrganization(ctx context.Context, accountID, organizationID uuid.UUID) ([]models.Questionnaire, error) {
	// TODO: Validate organization belongs to account

	return s.repo.FindByOrganizationID(ctx, organizationID)
}

func (s *QuestionnaireService) AddQuestion(ctx context.Context, accountID, questionnaireID uuid.UUID, question *models.Question) (*models.Question, error) {
	_, err := s.repo.FindByID(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	// TODO: Validate questionnaire belongs to account

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

func (s *QuestionnaireService) GenerateQuestionsForProduct(ctx context.Context, accountID uuid.UUID, product *menuModels.Product) ([]aiServices.GeneratedQuestion, error) {
	if s.questionGenerator == nil {
		return nil, fmt.Errorf("AI question generation is not configured")
	}

	// TODO: Validate product belongs to account

	questions, err := s.questionGenerator.GenerateQuestionsForProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	return questions, nil
}

func (s *QuestionnaireService) GenerateAndSaveQuestionnaireForProduct(ctx context.Context, accountID uuid.UUID, product *menuModels.Product, name, description string, isDefault bool) (*models.Questionnaire, error) {
	if s.questionGenerator == nil {
		return nil, fmt.Errorf("AI question generation is not configured")
	}

	generatedQuestions, err := s.questionGenerator.GenerateQuestionsForProduct(ctx, product)
	if err != nil {
		return nil, fmt.Errorf("failed to generate questions: %w", err)
	}

	questionnaire := &models.Questionnaire{
		OrganizationID: product.OrganizationID,
		ProductID:       &product.ID,
		Name:         name,
		Description:  description,
		IsDefault:    isDefault,
		IsActive:     true,
	}

	if questionnaire.IsDefault {
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, product.OrganizationID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Create(ctx, questionnaire); err != nil {
		return nil, err
	}

	for i, genQ := range generatedQuestions {
		question := &models.Question{
			// TODO: Update this to use new Question structure with ProductID
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

	return s.repo.FindByIDWithQuestions(ctx, questionnaire.ID)
}
