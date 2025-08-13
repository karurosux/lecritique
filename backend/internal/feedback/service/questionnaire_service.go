package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	aiServices "kyooar/internal/ai/services"
	feedbackinterface "kyooar/internal/feedback/interface"
	feedbackmodel "kyooar/internal/feedback/model"
	productModels "kyooar/internal/product/models"
)

type questionnaireService struct {
	repo              feedbackinterface.QuestionnaireRepository
	questionGenerator *aiServices.QuestionGenerator
}

func NewQuestionnaireService(
	repo feedbackinterface.QuestionnaireRepository,
	generator *aiServices.QuestionGenerator,
) feedbackinterface.QuestionnaireService {
	return &questionnaireService{
		repo:              repo,
		questionGenerator: generator,
	}
}

func (s *questionnaireService) Create(ctx context.Context, accountID, organizationID uuid.UUID, input *feedbackmodel.Questionnaire) (*feedbackmodel.Questionnaire, error) {
	questionnaire := &feedbackmodel.Questionnaire{
		OrganizationID: organizationID,
		ProductID:      input.ProductID,
		Name:           input.Name,
		Description:    input.Description,
		IsDefault:      input.IsDefault,
		IsActive:       true,
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

func (s *questionnaireService) GetByID(ctx context.Context, accountID, questionnaireID uuid.UUID) (*feedbackmodel.Questionnaire, error) {
	return s.repo.FindByID(ctx, questionnaireID)
}

func (s *questionnaireService) ListByOrganization(ctx context.Context, accountID, organizationID uuid.UUID) ([]*feedbackmodel.Questionnaire, error) {
	return s.repo.FindByOrganizationID(ctx, organizationID)
}

func (s *questionnaireService) Update(ctx context.Context, accountID, questionnaireID uuid.UUID, input *feedbackmodel.Questionnaire) (*feedbackmodel.Questionnaire, error) {
	questionnaire, err := s.repo.FindByID(ctx, questionnaireID)
	if err != nil {
		return nil, err
	}

	questionnaire.Name = input.Name
	questionnaire.Description = input.Description
	questionnaire.IsDefault = input.IsDefault
	questionnaire.IsActive = input.IsActive

	if questionnaire.IsDefault {
		if err := s.repo.DeactivateDefaultQuestionnaires(ctx, questionnaire.OrganizationID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.Update(ctx, questionnaire); err != nil {
		return nil, err
	}

	return questionnaire, nil
}

func (s *questionnaireService) Delete(ctx context.Context, accountID, questionnaireID uuid.UUID) error {
	return s.repo.Delete(ctx, questionnaireID)
}

func (s *questionnaireService) AddQuestion(ctx context.Context, accountID, questionnaireID uuid.UUID, question *feedbackmodel.Question) (*feedbackmodel.Question, error) {
	if err := s.repo.CreateQuestion(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *questionnaireService) UpdateQuestion(ctx context.Context, accountID, questionID uuid.UUID, question *feedbackmodel.Question) (*feedbackmodel.Question, error) {
	if err := s.repo.UpdateQuestion(ctx, question); err != nil {
		return nil, err
	}
	return question, nil
}

func (s *questionnaireService) DeleteQuestion(ctx context.Context, accountID, questionID uuid.UUID) error {
	return s.repo.DeleteQuestion(ctx, questionID)
}

func (s *questionnaireService) ReorderQuestions(ctx context.Context, accountID, questionnaireID uuid.UUID, questionIDs []uuid.UUID) error {
	return s.repo.ReorderQuestions(ctx, questionnaireID, questionIDs)
}

func (s *questionnaireService) GenerateQuestionsForProduct(ctx context.Context, accountID uuid.UUID, product *productModels.Product) ([]*feedbackmodel.GeneratedQuestion, error) {
	if s.questionGenerator == nil {
		return nil, fmt.Errorf("question generator not available")
	}

	return s.questionGenerator.GenerateQuestionsForProduct(ctx, product)
}

func (s *questionnaireService) GenerateAndSaveQuestionnaireForProduct(ctx context.Context, accountID uuid.UUID, product *productModels.Product, name, description string, isDefault bool) (*feedbackmodel.Questionnaire, error) {
	if s.questionGenerator == nil {
		return nil, fmt.Errorf("question generator not available")
	}

	generatedQuestions, err := s.questionGenerator.GenerateQuestionsForProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	questionnaire := &feedbackmodel.Questionnaire{
		OrganizationID: product.OrganizationID,
		ProductID:      &product.ID,
		Name:           name,
		Description:    description,
		IsDefault:      isDefault,
		IsActive:       true,
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
		question := &feedbackmodel.Question{
			ProductID:    product.ID,
			Text:         genQ.Text,
			Type:         genQ.Type,
			IsRequired:   true,
			DisplayOrder: i + 1,
			Options:      genQ.Options,
			MinValue:     genQ.MinValue,
			MaxValue:     genQ.MaxValue,
			MinLabel:     genQ.MinLabel,
			MaxLabel:     genQ.MaxLabel,
		}

		if err := s.repo.CreateQuestion(ctx, question); err != nil {
			return nil, err
		}
		questionnaire.Questions = append(questionnaire.Questions, *question)
	}

	return questionnaire, nil
}