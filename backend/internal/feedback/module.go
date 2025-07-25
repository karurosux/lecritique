package feedback

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/feedback/handlers"
	"github.com/samber/do"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	// Get handlers from injector
	publicHandler := do.MustInvoke[*handlers.FeedbackPublicHandler](m.injector)
	
	// Public feedback routes (no auth required)
	v1.GET("/questionnaire/:organizationId/:productId", publicHandler.GetQuestionnaire)
	v1.GET("/public/organization/:organizationId/products/:productId/questions", publicHandler.GetProductQuestions)
	v1.GET("/public/organization/:organizationId/questions/products-with-questions", publicHandler.GetProductsWithQuestions)
	v1.POST("/public/feedback", publicHandler.SubmitFeedback)
	
	// Protected feedback routes (moved to organization module)
}
