package feedback

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	feedbackcontroller "kyooar/internal/feedback/controller"
)

type Module struct {
	injector *do.Injector
}

func NewModule(i *do.Injector) *Module {
	return &Module{injector: i}
}

func (m *Module) RegisterRoutes(v1 *echo.Group) {
	publicController := do.MustInvoke[*feedbackcontroller.PublicController](m.injector)
	v1.GET("/questionnaire/:organizationId/:productId", publicController.GetQuestionnaire)
	v1.GET("/public/organization/:organizationId/products/:productId/questions", publicController.GetProductQuestions)
	v1.GET("/public/organization/:organizationId/questions/products-with-questions", publicController.GetProductsWithQuestions)
	v1.POST("/public/feedback", publicController.SubmitFeedback)
}
