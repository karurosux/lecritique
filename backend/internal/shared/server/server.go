package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/logger"
	sharedMiddleware "github.com/lecritique/api/internal/shared/middleware"
	"github.com/lecritique/api/internal/shared/cron"
	
	// Domain handlers for route registration
	authHandlers "github.com/lecritique/api/internal/auth/handlers"
	authServices "github.com/lecritique/api/internal/auth/services"
	feedbackHandlers "github.com/lecritique/api/internal/feedback/handlers"
	menuHandlers "github.com/lecritique/api/internal/menu/handlers"
	restaurantHandlers "github.com/lecritique/api/internal/restaurant/handlers"
	qrcodeHandlers "github.com/lecritique/api/internal/qrcode/handlers"
	analyticsHandlers "github.com/lecritique/api/internal/analytics/handlers"
	subscriptionHandlers "github.com/lecritique/api/internal/subscription/handlers"
	
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	cronLib "github.com/robfig/cron/v3"
	
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *gorm.DB
	cron   *cronLib.Cron
	authService authServices.AuthService
}

func New(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	
	// Configure Echo
	e.HideBanner = false  // Show Echo banner for startup info
	e.HidePort = false    // Show port info
	
	// Set logger level based on environment
	if cfg.App.Env == "development" {
		logger.SetLevel("debug")
	}
	
	// Custom error handler
	e.HTTPErrorHandler = customErrorHandler
	
	// Setup middleware
	setupMiddleware(e, cfg)
	
	s := &Server{
		echo:   e,
		config: cfg,
		db:     db,
	}
	
	// Setup routes
	s.setupRoutes()
	
	// Setup cron jobs
	s.setupCronJobs()
	
	return s
}

func (s *Server) setupRoutes() {
	// Health check route
	s.echo.GET("/api/health", s.healthCheck)

	// Swagger documentation with relaxed CSP
	swaggerGroup := s.echo.Group("/swagger")
	swaggerGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Relax CSP for Swagger UI
			c.Response().Header().Set("Content-Security-Policy", "default-src 'self' 'unsafe-inline' 'unsafe-eval'; img-src 'self' data:; font-src 'self' data:")
			return next(c)
		}
	})
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)

	// Rate limiter
	rateLimiter := sharedMiddleware.NewRateLimiter(100, time.Minute)

	// API v1 routes
	v1 := s.echo.Group("/api/v1")
	v1.Use(rateLimiter.Middleware())

	// Register domain routes
	authService := authHandlers.RegisterRoutes(v1, s.db, s.config)
	
	// Store authService for cron jobs
	s.authService = authService
	
	// Public and protected route groups
	public := v1.Group("/public")
	protected := v1.Group("")
	
	// Register all routes
	feedbackHandlers.RegisterRoutes(public, protected, s.db, authService)
	restaurantHandlers.RegisterRoutes(protected, s.db, authService)
	menuHandlers.RegisterRoutes(protected, s.db, authService)
	qrcodeHandlers.RegisterRoutes(protected, s.db, authService)
	analyticsHandlers.RegisterRoutes(protected, s.db, authService)
	subscriptionHandlers.RegisterRoutes(v1, s.db, authService)
}

func (s *Server) setupCronJobs() {
	// Setup deactivation cron job
	if s.authService != nil {
		s.cron = cron.SetupDeactivationCron(s.authService)
		logger.Info("Cron jobs initialized", logrus.Fields{
			"job": "account_deactivation",
		})
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags system
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/health [get]
func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "lecritique-api",
		"time":    time.Now(),
	})
}

func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	// Request ID middleware (first, so all logs include request ID)
	e.Use(middleware.RequestID())
	
	// Custom structured logging middleware
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip health check logs to reduce noise
			if c.Path() == "/api/health" {
				return next(c)
			}
			
			start := time.Now()
			err := next(c)
			stop := time.Now()
			
			req := c.Request()
			res := c.Response()
			
			// Log request details
			fields := logrus.Fields{
				"request_id":    c.Response().Header().Get(echo.HeaderXRequestID),
				"method":        req.Method,
				"path":          req.URL.Path,
				"query":         req.URL.RawQuery,
				"status":        res.Status,
				"latency_ms":    stop.Sub(start).Milliseconds(),
				"bytes_in":      req.Header.Get(echo.HeaderContentLength),
				"bytes_out":     res.Size,
				"user_agent":    req.UserAgent(),
				"remote_ip":     c.RealIP(),
				"referer":       req.Referer(),
			}
			
			msg := fmt.Sprintf("%s %s", req.Method, req.URL.Path)
			
			if err != nil {
				fields["error"] = err.Error()
				logger.Error(msg, err, fields)
			} else if res.Status >= 500 {
				logger.Error(msg, fmt.Errorf("server error"), fields)
			} else if res.Status >= 400 {
				logger.Warn(msg, fields)
			} else {
				logger.Info(msg, fields)
			}
			
			return err
		}
	})
	
	// Recover middleware with structured logging
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Error("panic recovered", err, logrus.Fields{
				"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
				"method":     c.Request().Method,
				"path":       c.Request().URL.Path,
				"stack":      string(stack),
			})
			return err
		},
	}))
	
	// CORS middleware
	corsConfig := middleware.DefaultCORSConfig
	if cfg.App.Env == "production" {
		corsConfig.AllowOrigins = []string{cfg.App.URL}
	} else {
		corsConfig.AllowOrigins = []string{"*"}
	}
	e.Use(middleware.CORSWithConfig(corsConfig))
	
	// Secure middleware
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         "1; mode=block",
		ContentTypeNosniff:    "nosniff",
		XFrameOptions:         "SAMEORIGIN",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	}))
	
	// Body limit
	e.Use(middleware.BodyLimit("2M"))
	
	// Gzip compression
	e.Use(middleware.Gzip())
}

func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"
	
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	}
	
	// Create error response
	errorResponse := map[string]interface{}{
		"success": false,
		"error": map[string]interface{}{
			"code":    fmt.Sprintf("HTTP_%d", code),
			"message": message,
		},
		"timestamp": time.Now().Format(time.RFC3339),
		"request_id": c.Response().Header().Get(echo.HeaderXRequestID),
	}
	
	// Send response with proper JSON formatting
	if !c.Response().Committed {
		// Pretty print JSON in development
		if c.Echo().Debug {
			c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
			encoder := json.NewEncoder(c.Response())
			encoder.SetIndent("", "  ")
			c.Response().WriteHeader(code)
			encoder.Encode(errorResponse)
		} else {
			c.JSON(code, errorResponse)
		}
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.App.Port)
	
	// Enable debug mode in development for pretty JSON
	if s.config.App.Env == "development" {
		s.echo.Debug = true
	}
	
	// Start cron jobs
	if s.cron != nil {
		s.cron.Start()
		logger.Info("Cron scheduler started", nil)
	}
	
	logger.Info("Server starting", logrus.Fields{
		"name":        s.config.App.Name,
		"address":     addr,
		"environment": s.config.App.Env,
		"debug":       s.echo.Debug,
	})
	
	return s.echo.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	// Stop cron jobs
	if s.cron != nil {
		s.cron.Stop()
		logger.Info("Cron scheduler stopped", nil)
	}
	
	return s.echo.Shutdown(ctx)
}
