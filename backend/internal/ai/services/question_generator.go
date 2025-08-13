package services

import (
	"context"
	"fmt"
	"strings"

	feedbackmodel "kyooar/internal/feedback/model"
	menuModels "kyooar/internal/product/models"
	"kyooar/internal/shared/config"
)

type QuestionGenerator struct {
	config   *config.Config
	provider AIProvider
}

type AIProvider interface {
	GenerateQuestions(ctx context.Context, prompt string) ([]GeneratedQuestion, error)
}

type GeneratedQuestion struct {
	Text         string               `json:"text"`
	Type         feedbackmodel.QuestionType  `json:"type"`
	Options      []string             `json:"options,omitempty"`
	MinValue     *int                 `json:"min_value,omitempty"`
	MaxValue     *int                 `json:"max_value,omitempty"`
	MinLabel     string               `json:"min_label,omitempty"`
	MaxLabel     string               `json:"max_label,omitempty"`
}

func NewQuestionGenerator(cfg *config.Config) (*QuestionGenerator, error) {
	var provider AIProvider
	
	switch cfg.AI.Provider {
	case "anthropic":
		provider = NewAnthropicProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "openai":
		provider = NewOpenAIProvider(cfg.AI.APIKey, cfg.AI.Model)
	case "gemini":
		provider = NewGeminiProvider(cfg.AI.APIKey, cfg.AI.Model)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", cfg.AI.Provider)
	}

	return &QuestionGenerator{
		config:   cfg,
		provider: provider,
	}, nil
}

func (qg *QuestionGenerator) GenerateQuestionsForProduct(ctx context.Context, product *menuModels.Product) ([]*feedbackmodel.GeneratedQuestion, error) {
	prompt := qg.buildPromptForProduct(product)
	aiQuestions, err := qg.provider.GenerateQuestions(ctx, prompt)
	if err != nil {
		return nil, err
	}

	// Convert from AI service GeneratedQuestion to feedback model GeneratedQuestion
	result := make([]*feedbackmodel.GeneratedQuestion, len(aiQuestions))
	for i, q := range aiQuestions {
		result[i] = &feedbackmodel.GeneratedQuestion{
			Text:     q.Text,
			Type:     q.Type,
			Options:  q.Options,
			MinValue: q.MinValue,
			MaxValue: q.MaxValue,
			MinLabel: q.MinLabel,
			MaxLabel: q.MaxLabel,
		}
	}
	
	return result, nil
}

func (qg *QuestionGenerator) buildPromptForProduct(product *menuModels.Product) string {
	tagsStr := ""
	if len(product.Tags) > 0 {
		tagsStr = fmt.Sprintf("\nTags: %s", strings.Join(product.Tags, ", "))
	}

	return fmt.Sprintf(`Generate 5-7 specific feedback questions for the following product in a organization. The questions should help the organization gather actionable feedback to improve this specific product.

Product Name: %s
Description: %s
Category: %s
Price: %.2f %s%s

Generate questions in the following JSON format:
[
  {
    "text": "Question text here",
    "type": "rating|scale|multi_choice|single_choice|text|yes_no",
    "options": ["option1", "option2"],
    "min_value": 1,
    "max_value": 10,
    "min_label": "label",
    "max_label": "label"
  }
]

Guidelines:
1. Make questions specific to the product (not generic)
2. Include a mix of question types
3. Focus on actionable feedback (taste, texture, presentation, portion, temperature, etc.)
4. For rating questions, use 1-5 scale
5. For scale questions, use 1-10 scale with descriptive labels
6. Keep questions concise and clear
7. Avoid yes/no questions unless they're very specific

Return ONLY the JSON array, no additional text.`, 
		product.Name,
		product.Description,
		product.Category,
		product.Price,
		product.Currency,
		tagsStr,
	)
}
