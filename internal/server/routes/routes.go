package routes

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/shared/config"
	"github.com/lecritique/api/internal/handlers"
	"github.com/lecritique/api/shared/middleware"
	"github.com/lecritique/api/internal/repositories"
	"github.com/lecritique/api/internal/services"
	"gorm.io/gorm"
)

func Setup(e *echo.Echo, db *gorm.DB, cfg *config.Config) {
	// Initialize repositories
	accountRepo := repositories.NewAccountRepository(db)
	restaurantRepo := repositories.NewRestaurantRepository(db)
	dishRepo := repositories.NewDishRepository(db)
	qrCodeRepo := repositories.NewQRCodeRepository(db)
	feedbackRepo := repositories.NewFeedbackRepository(db)
	subscriptionRepo := repositories.NewSubscriptionRepository(db)
	questionnaireRepo := repositories.NewQuestionnaireRepository(db)

	// Initialize services
	authService := services.NewAuthService(accountRepo, cfg)
	restaurantService := services.NewRestaurantService(restaurantRepo, subscriptionRepo)
	dishService := services.NewDishService(dishRepo, restaurantRepo)
	qrCodeService := services.NewQRCodeService(qrCodeRepo, restaurantRepo)
	feedbackService := services.NewFeedbackService(feedbackRepo, restaurantRepo, qrCodeRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	restaurantHandler := handlers.NewRestaurantHandler(restaurantService)
	dishHandler := handlers.NewDishHandler(dishService)
	publicHandler := handlers.NewPublicHandler(qrCodeService, feedbackService, dishRepo, questionnaireRepo)

	// Health check
	e.GET("/api/health", healthCheck)

	// Rate limiter
	rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	// API v1 routes
	v1 := e.Group("/api/v1")
	v1.Use(rateLimiter.Middleware())

	// Public routes (no auth required)
	public := v1.Group("/public")
	setupPublicRoutes(public, publicHandler)

	// Auth routes
	auth := v1.Group("/auth")
	setupAuthRoutes(auth, authHandler)

	// Protected routes (auth required)
	protected := v1.Group("")
	protected.Use(middleware.JWTAuth(authService))
	setupProtectedRoutes(protected, restaurantHandler, dishHandler)
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "lecritique-api",
		"time":    time.Now(),
	})
}

func setupPublicRoutes(g *echo.Group, handler *handlers.PublicHandler) {
	g.GET("/qr/:code", handler.ValidateQRCode)
	g.GET("/restaurant/:id/menu", handler.GetRestaurantMenu)
	g.GET("/questionnaire/:restaurantId/:dishId", handler.GetQuestionnaire)
	g.POST("/feedback", handler.SubmitFeedback)
}

func setupAuthRoutes(g *echo.Group, handler *handlers.AuthHandler) {
	g.POST("/register", handler.Register)
	g.POST("/login", handler.Login)
	g.POST("/refresh", handler.RefreshToken)
}

func setupProtectedRoutes(g *echo.Group, restaurantHandler *handlers.RestaurantHandler, dishHandler *handlers.DishHandler) {
	// Restaurant routes
	restaurants := g.Group("/restaurants")
	restaurants.POST("", restaurantHandler.Create)
	restaurants.GET("", restaurantHandler.GetAll)
	restaurants.GET("/:id", restaurantHandler.GetByID)
	restaurants.PUT("/:id", restaurantHandler.Update)
	restaurants.DELETE("/:id", restaurantHandler.Delete)
	restaurants.GET("/:restaurantId/dishes", dishHandler.GetByRestaurant)

	// Dish routes
	dishes := g.Group("/dishes")
	dishes.POST("", dishHandler.Create)
	dishes.GET("/:id", dishHandler.GetByID)
	dishes.PUT("/:id", dishHandler.Update)
	dishes.DELETE("/:id", dishHandler.Delete)
}
