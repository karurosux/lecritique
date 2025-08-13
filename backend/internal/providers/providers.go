package providers

import (
	
	organization "kyooar/internal/organization"
	
	
	
	
	
	
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/middleware"
	sharedServices "kyooar/internal/shared/services"
	
	"github.com/samber/do"
	"gorm.io/gorm"
)

func RegisterAll(i *do.Injector, cfg *config.Config, db *gorm.DB) {
	do.ProvideValue(i, cfg)
	do.ProvideValue(i, db)
	
	do.Provide(i, sharedServices.NewEmailService)
	do.Provide(i, middleware.NewMiddlewareProvider)
	
	organization.RegisterNewModule(i)
}
