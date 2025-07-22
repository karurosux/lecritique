package qrcode

import (
	"github.com/labstack/echo/v4"
	"kyooar/internal/qrcode/handlers"
	sharedMiddleware "kyooar/internal/shared/middleware"
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
	qrCodeHandler := do.MustInvoke[*handlers.QRCodeHandler](m.injector)
	publicHandler := do.MustInvoke[*handlers.QRCodePublicHandler](m.injector)
	
	// Get middleware provider
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public QR code routes (no auth required)
	v1.GET("/public/qr/:code", publicHandler.ValidateQRCode)
	
	// QR Code routes under organizations (moved to organization module)
	
	// Direct QR code routes
	qrCodes := v1.Group("/qr-codes")
	qrCodes.Use(middlewareProvider.AuthMiddleware())
	qrCodes.Use(middlewareProvider.TeamAwareMiddleware())
	qrCodes.PATCH("/:id", qrCodeHandler.Update)
	qrCodes.DELETE("/:id", qrCodeHandler.Delete)
}
