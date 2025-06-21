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
	
	// Domain handlers for route registration
	authHandlers "github.com/lecritique/api/internal/auth/handlers"
	// feedbackHandlers "github.com/lecritique/api/internal/feedback/handlers"   // TODO: Fix domain models first
	// menuHandlers "github.com/lecritique/api/internal/menu/handlers"           // TODO: Fix domain models first
	// restaurantHandlers "github.com/lecritique/api/internal/restaurant/handlers" // TODO: Fix domain models first
	
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Server struct {
	echo   *echo.Echo
	config *config.Config
	db     *gorm.DB
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
	
	return s
}

func (s *Server) setupRoutes() {
	// Health check route
	s.echo.GET("/api/health", s.healthCheck)

	// Rate limiter
	rateLimiter := sharedMiddleware.NewRateLimiter(100, time.Minute)

	// API v1 routes
	v1 := s.echo.Group("/api/v1")
	v1.Use(rateLimiter.Middleware())

	// Register domain routes
	authHandlers.RegisterRoutes(v1, s.db, s.config)
	
	// TODO: Uncomment when domains are ready
	// public := v1.Group("/public")
	// protected := v1.Group("")
	// feedbackHandlers.RegisterRoutes(public, s.db)                   // TODO: Fix domain models first
	// restaurantHandlers.RegisterRoutes(protected, s.db, nil)        // TODO: Fix domain models first  
	// menuHandlers.RegisterRoutes(protected, s.db, nil)              // TODO: Fix domain models first
}

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
	
	logger.Info("Server starting", logrus.Fields{
		"name":        s.config.App.Name,
		"address":     addr,
		"environment": s.config.App.Env,
		"debug":       s.echo.Debug,
	})
	
	return s.echo.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.echo.Shutdown(ctx)
}
