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
	v1.GET("/questionnaire/:restaurantId/:dishId", publicHandler.GetQuestionnaire)
	v1.GET("/restaurant/:restaurantId/dishes/:dishId/questions", publicHandler.GetDishQuestions)
	v1.GET("/restaurant/:restaurantId/questions/dishes-with-questions", publicHandler.GetDishesWithQuestions)
	v1.POST("/feedback", publicHandler.SubmitFeedback)
	
	// Protected feedback routes (moved to restaurant module)
}