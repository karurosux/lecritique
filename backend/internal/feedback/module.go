package feedback

import (
	"github.com/labstack/echo/v4"
	"lecritique/internal/feedback/handlers"
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
	v1.GET("/organization/:organizationId/products/:productId/questions", publicHandler.GetProductQuestions)
	v1.GET("/organization/:organizationId/questions/products-with-questions", publicHandler.GetProductesWithQuestions)
	v1.POST("/feedback", publicHandler.SubmitFeedback)
	
	// Protected feedback routes (moved to organization module)
}
