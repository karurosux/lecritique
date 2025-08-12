package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/logger"
	sharedMiddleware "kyooar/internal/shared/middleware"
	"kyooar/internal/shared/cron"
	"kyooar/internal/providers"
	
	authModule "kyooar/internal/auth"
	organizationModule "kyooar/internal/organization"
	menuModule "kyooar/internal/menu"
	feedbackModule "kyooar/internal/feedback"
	qrcodeModule "kyooar/internal/qrcode"
	analyticsModule "kyooar/internal/analytics"
	subscriptionModule "kyooar/internal/subscription"
	
	authinterface "kyooar/internal/auth/interface"
	
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	cronLib "github.com/robfig/cron/v3"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type Server struct {
	echo     *echo.Echo
	config   *config.Config
	db       *gorm.DB
	cron     *cronLib.Cron
	injector *do.Injector
}

func NewWithDI(cfg *config.Config, db *gorm.DB) *Server {
	e := echo.New()
	
	e.HideBanner = false
	e.HidePort = false
	
	if cfg.App.Env == "development" {
		logger.SetLevel("debug")
	}
	
	e.HTTPErrorHandler = customErrorHandler
	
	setupMiddleware(e, cfg)
	
	injector := do.New()
	providers.RegisterAll(injector, cfg, db)
	
	s := &Server{
		echo:     e,
		config:   cfg,
		db:       db,
		injector: injector,
	}
	
	s.setupRoutes()
	
	s.setupCronJobs()
	
	return s
}

func (s *Server) setupRoutes() {
	s.echo.GET("/api/health", s.healthCheck)

	swaggerGroup := s.echo.Group("/swagger")
	swaggerGroup.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Content-Security-Policy", 
				"default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:;")
			return next(c)
		}
	})
	swaggerGroup.GET("/*", echoSwagger.WrapHandler)

	rateLimiter := sharedMiddleware.NewRateLimiter(100, time.Minute)

	v1 := s.echo.Group("/api/v1")
	v1.Use(rateLimiter.Middleware())

	authMod := authModule.NewAuthModule(s.injector)
	authMod.RegisterRoutes(v1)
	
	// Register QRCode module before Organization module since Organization depends on it
	qrcodeModule.RegisterNewModule(s.injector)
	qrcodeMod := qrcodeModule.NewQRCodeModule(s.injector)
	qrcodeMod.RegisterRoutes(v1)
	
	organizationMod := organizationModule.NewOrganizationModule(s.injector)
	organizationMod.RegisterRoutes(v1)
	
	menuMod := menuModule.NewModule(s.injector)
	menuMod.RegisterRoutes(v1)
	
	feedbackMod := feedbackModule.NewModule(s.injector)
	feedbackMod.RegisterRoutes(v1)
	
	analyticsMod := analyticsModule.NewModule(s.injector)
	analyticsMod.RegisterRoutes(v1)
	
	subscriptionMod := subscriptionModule.NewModule(s.injector)
	subscriptionMod.RegisterRoutes(v1)
}

func (s *Server) setupCronJobs() {
	authService := do.MustInvoke[authinterface.AuthService](s.injector)
	s.cron = cron.SetupDeactivationCron(authService)
	logger.Info("Cron jobs initialized", logrus.Fields{
		"job": "account_deactivation",
	})
}

func (s *Server) healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "kyooar-api",
		"time":    time.Now(),
	})
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.config.App.Port)
	logger.Info("Starting server", logrus.Fields{
		"address": addr,
		"env":     s.config.App.Env,
	})
	
	if s.cron != nil {
		s.cron.Start()
	}
	
	return s.echo.Start(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	if s.cron != nil {
		s.cron.Stop()
	}
	return s.echo.Shutdown(ctx)
}

func setupMiddleware(e *echo.Echo, cfg *config.Config) {
	e.Use(middleware.Recover())
	
	if cfg.App.Env == "development" {
		e.Use(middleware.Logger())
	}
	
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS, echo.PATCH},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))
}

func customErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal server error"
	
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		message = fmt.Sprintf("%v", he.Message)
	}
	
	if !c.Response().Committed {
		if c.Request().Method == echo.HEAD {
			c.NoContent(code)
		} else {
			c.JSON(code, map[string]interface{}{
				"error": message,
			})
		}
	}
}
