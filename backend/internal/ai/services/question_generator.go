package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/lecritique/api/internal/feedback/models"
	menuModels "github.com/lecritique/api/internal/menu/models"
	"github.com/lecritique/api/internal/shared/config"
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
	Type         models.QuestionType  `json:"type"`
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

func (qg *QuestionGenerator) GenerateQuestionsForDish(ctx context.Context, dish *menuModels.Dish) ([]GeneratedQuestion, error) {
	prompt := qg.buildPromptForDish(dish)
	return qg.provider.GenerateQuestions(ctx, prompt)
}

func (qg *QuestionGenerator) buildPromptForDish(dish *menuModels.Dish) string {
	tagsStr := ""
	if len(dish.Tags) > 0 {
		tagsStr = fmt.Sprintf("\nTags: %s", strings.Join(dish.Tags, ", "))
	}

	return fmt.Sprintf(`Generate 5-7 specific feedback questions for the following dish in a restaurant. The questions should help the restaurant gather actionable feedback to improve this specific dish.

Dish Name: %s
Description: %s
Category: %s
Price: %.2f %s%s

Generate questions in the following JSON format:
[
  {
    "text": "Question text here",
    "type": "rating|scale|multi_choice|single_choice|text|yes_no",
    "options": ["option1", "option2"], // only for choice types
    "min_value": 1, // only for scale
    "max_value": 10, // only for scale
    "min_label": "label", // only for scale
    "max_label": "label" // only for scale
  }
]

Guidelines:
1. Make questions specific to the dish (not generic)
2. Include a mix of question types
3. Focus on actionable feedback (taste, texture, presentation, portion, temperature, etc.)
4. For rating questions, use 1-5 scale
5. For scale questions, use 1-10 scale with descriptive labels
6. Keep questions concise and clear
7. Avoid yes/no questions unless they're very specific

Return ONLY the JSON array, no additional text.`, 
		dish.Name,
		dish.Description,
		dish.Category,
		dish.Price,
		dish.Currency,
		tagsStr,
	)
}