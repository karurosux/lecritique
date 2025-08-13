package ai

import (
	"github.com/samber/do"
	"kyooar/internal/ai/services"
	"kyooar/internal/shared/config"
)

func ProvideQuestionGenerator(i *do.Injector) (*services.QuestionGenerator, error) {
	cfg := do.MustInvoke[*config.Config](i)
	return services.NewQuestionGenerator(cfg)
}

func RegisterNewModule(container *do.Injector) error {
	do.Provide(container, ProvideQuestionGenerator)
	return nil
}