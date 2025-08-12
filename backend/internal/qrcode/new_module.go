package qrcode

import (
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
	"gorm.io/gorm"

	organizationinterface "kyooar/internal/organization/interface"
	qrcodecontroller "kyooar/internal/qrcode/controller"
	qrcodeinterface "kyooar/internal/qrcode/interface"
	gormqrcode "kyooar/internal/qrcode/repository/gorm"
	qrcodeservice "kyooar/internal/qrcode/service"
	sharedMiddleware "kyooar/internal/shared/middleware"
)

func ProvideQRCodeRepository(i *do.Injector) (qrcodeinterface.QRCodeRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return gormqrcode.NewQRCodeRepository(db), nil
}

func ProvideQRCodeService(i *do.Injector) (qrcodeinterface.QRCodeService, error) {
	qrCodeRepo := do.MustInvoke[qrcodeinterface.QRCodeRepository](i)
	organizationRepo := do.MustInvoke[organizationinterface.OrganizationRepository](i)

	return qrcodeservice.NewQRCodeService(
		qrCodeRepo,
		organizationRepo,
	), nil
}

func ProvideQRCodeController(i *do.Injector) (*qrcodecontroller.QRCodeController, error) {
	qrCodeService := do.MustInvoke[qrcodeinterface.QRCodeService](i)
	return qrcodecontroller.NewQRCodeController(qrCodeService), nil
}

func ProvidePublicController(i *do.Injector) (*qrcodecontroller.PublicController, error) {
	qrCodeService := do.MustInvoke[qrcodeinterface.QRCodeService](i)
	return qrcodecontroller.NewPublicController(qrCodeService), nil
}

type QRCodeModule struct {
	injector *do.Injector
}

func NewQRCodeModule(i *do.Injector) *QRCodeModule {
	return &QRCodeModule{injector: i}
}

func (m *QRCodeModule) RegisterRoutes(v1 *echo.Group) {
	qrCodeController := do.MustInvoke[*qrcodecontroller.QRCodeController](m.injector)
	publicController := do.MustInvoke[*qrcodecontroller.PublicController](m.injector)
	
	middlewareProvider := do.MustInvoke[*sharedMiddleware.MiddlewareProvider](m.injector)
	
	// Public routes
	v1.GET("/public/qr/:code", publicController.ValidateQRCode)
	
	// QR code CRUD routes
	qrCodes := v1.Group("/qr-codes")
	qrCodes.Use(middlewareProvider.AuthMiddleware())
	qrCodes.Use(middlewareProvider.TeamAwareMiddleware())
	qrCodes.PATCH("/:id", qrCodeController.Update)
	qrCodes.DELETE("/:id", qrCodeController.Delete)
}

func RegisterNewModule(container *do.Injector) error {
	do.Provide(container, ProvideQRCodeRepository)
	do.Provide(container, ProvideQRCodeService)
	do.Provide(container, ProvideQRCodeController)
	do.Provide(container, ProvidePublicController)

	return nil
}