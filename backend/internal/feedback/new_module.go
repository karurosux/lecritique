package feedback

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"

	feedbackcontroller "kyooar/internal/feedback/controller"
	feedbackinterface "kyooar/internal/feedback/interface"
	feedbackmiddleware "kyooar/internal/feedback/middleware"
	gormrepo "kyooar/internal/feedback/repository/gorm"
	feedbackservice "kyooar/internal/feedback/service"
	menuRepos "kyooar/internal/menu/repositories"
	menuServices "kyooar/internal/menu/services"
	organizationRepos "kyooar/internal/organization/repositories"
	qrcodeRepos "kyooar/internal/qrcode/repositories"
	"kyooar/internal/shared/config"
)

func ProvideFeedbackRepository(i *do.Injector) (feedbackinterface.FeedbackRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewFeedbackRepository(db), nil
}

func ProvideQuestionRepository(i *do.Injector) (feedbackinterface.QuestionRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewQuestionRepository(db), nil
}

func ProvideQuestionnaireRepository(i *do.Injector) (feedbackinterface.QuestionnaireRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormrepo.NewQuestionnaireRepository(db), nil
}

func ProvideFeedbackService(i *do.Injector) (feedbackinterface.FeedbackService, error) {
	feedbackRepo := do.MustInvoke[feedbackinterface.FeedbackRepository](i)
	organizationRepo := do.MustInvoke[organizationRepos.OrganizationRepository](i)
	qrCodeRepo := do.MustInvoke[qrcodeRepos.QRCodeRepository](i)

	return feedbackservice.NewFeedbackService(
		feedbackRepo,
		organizationRepo,
		qrCodeRepo,
	), nil
}

func ProvideQuestionService(i *do.Injector) (feedbackinterface.QuestionService, error) {
	questionRepo := do.MustInvoke[feedbackinterface.QuestionRepository](i)
	productRepo := do.MustInvoke[menuRepos.ProductRepository](i)
	organizationRepo := do.MustInvoke[organizationRepos.OrganizationRepository](i)

	return feedbackservice.NewQuestionService(
		questionRepo,
		productRepo,
		organizationRepo,
	), nil
}

func ProvideQuestionnaireService(i *do.Injector) (feedbackinterface.QuestionnaireService, error) {
	questionnaireRepo := do.MustInvoke[feedbackinterface.QuestionnaireRepository](i)
	config := do.MustInvoke[*config.Config](i)

	return feedbackservice.NewQuestionnaireService(
		questionnaireRepo,
		config,
	), nil
}

func ProvideFeedbackMiddleware(i *do.Injector) (*feedbackmiddleware.FeedbackMiddleware, error) {
	feedbackService := do.MustInvoke[feedbackinterface.FeedbackService](i)
	return feedbackmiddleware.NewFeedbackMiddleware(feedbackService), nil
}

func ProvideFeedbackController(i *do.Injector) (*feedbackcontroller.FeedbackController, error) {
	feedbackService := do.MustInvoke[feedbackinterface.FeedbackService](i)
	return feedbackcontroller.NewFeedbackController(feedbackService), nil
}

func ProvideQuestionController(i *do.Injector) (*feedbackcontroller.QuestionController, error) {
	questionService := do.MustInvoke[feedbackinterface.QuestionService](i)
	return feedbackcontroller.NewQuestionController(questionService), nil
}

func ProvideQuestionnaireController(i *do.Injector) (*feedbackcontroller.QuestionnaireController, error) {
	questionnaireService := do.MustInvoke[feedbackinterface.QuestionnaireService](i)
	productService := do.MustInvoke[menuServices.ProductService](i)
	return feedbackcontroller.NewQuestionnaireController(questionnaireService, productService), nil
}

func ProvidePublicController(i *do.Injector) (*feedbackcontroller.PublicController, error) {
	feedbackService := do.MustInvoke[feedbackinterface.FeedbackService](i)
	productRepo := do.MustInvoke[menuRepos.ProductRepository](i)
	questionnaireRepo := do.MustInvoke[feedbackinterface.QuestionnaireRepository](i)
	questionRepo := do.MustInvoke[feedbackinterface.QuestionRepository](i)
	return feedbackcontroller.NewPublicController(feedbackService, productRepo, questionnaireRepo, questionRepo), nil
}

type FeedbackModule struct {
	injector *do.Injector
}

func NewFeedbackModule(i *do.Injector) *FeedbackModule {
	return &FeedbackModule{injector: i}
}

func (m *FeedbackModule) RegisterRoutes(v1 *echo.Group) {
	publicController := do.MustInvoke[*feedbackcontroller.PublicController](m.injector)
	feedbackController := do.MustInvoke[*feedbackcontroller.FeedbackController](m.injector)
	questionController := do.MustInvoke[*feedbackcontroller.QuestionController](m.injector)
	questionnaireController := do.MustInvoke[*feedbackcontroller.QuestionnaireController](m.injector)

	v1.GET("/questionnaire/:organizationId/:productId", publicController.GetQuestionnaire)
	v1.GET("/public/organization/:organizationId/products/:productId/questions", publicController.GetProductQuestions)
	v1.GET("/public/organization/:organizationId/questions/products-with-questions", publicController.GetProductsWithQuestions)
	v1.POST("/public/feedback", publicController.SubmitFeedback)

	orgGroup := v1.Group("/organizations/:organizationId")
	orgGroup.GET("/feedback", feedbackController.GetByOrganization)
	orgGroup.GET("/analytics", feedbackController.GetStats)

	orgGroup.POST("/products/:productId/questions", questionController.CreateQuestion)
	orgGroup.GET("/products/:productId/questions", questionController.GetQuestionsByProduct)
	orgGroup.GET("/questions/batch", questionController.GetQuestionsByProducts)
	orgGroup.GET("/products/:productId/questions/:questionId", questionController.GetQuestion)
	orgGroup.PUT("/products/:productId/questions/:questionId", questionController.UpdateQuestion)
	orgGroup.DELETE("/products/:productId/questions/:questionId", questionController.DeleteQuestion)
	orgGroup.POST("/products/:productId/questions/reorder", questionController.ReorderQuestions)
	orgGroup.GET("/questions/products-with-questions", questionController.GetProductsWithQuestions)

	orgGroup.POST("/questionnaires", questionnaireController.CreateQuestionnaire)
	orgGroup.GET("/questionnaires", questionnaireController.ListQuestionnaires)
	orgGroup.GET("/questionnaires/:id", questionnaireController.GetQuestionnaire)
	orgGroup.PUT("/questionnaires/:id", questionnaireController.UpdateQuestionnaire)
	orgGroup.DELETE("/questionnaires/:id", questionnaireController.DeleteQuestionnaire)
	orgGroup.POST("/questionnaires/:id/questions", questionnaireController.AddQuestion)
	orgGroup.PUT("/questionnaires/:id/questions/:questionId", questionnaireController.UpdateQuestion)
	orgGroup.DELETE("/questionnaires/:id/questions/:questionId", questionnaireController.DeleteQuestion)
	orgGroup.POST("/questionnaires/:id/reorder", questionnaireController.ReorderQuestions)

	aiGroup := v1.Group("/ai")
	aiGroup.POST("/generate-questions/:productId", questionnaireController.GenerateQuestions)
	aiGroup.POST("/generate-questionnaire/:productId", questionnaireController.GenerateAndSaveQuestionnaire)
}

func RegisterNewModule(container *do.Injector) error {
	do.Provide(container, ProvideFeedbackRepository)
	do.Provide(container, ProvideQuestionRepository)
	do.Provide(container, ProvideQuestionnaireRepository)
	do.Provide(container, ProvideFeedbackService)
	do.Provide(container, ProvideQuestionService)
	do.Provide(container, ProvideQuestionnaireService)
	do.Provide(container, ProvideFeedbackMiddleware)
	do.Provide(container, ProvideFeedbackController)
	do.Provide(container, ProvideQuestionController)
	do.Provide(container, ProvideQuestionnaireController)
	do.Provide(container, ProvidePublicController)

	return nil
}