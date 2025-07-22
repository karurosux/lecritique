package providers

import (
	"kyooar/internal/shared/config"
	"github.com/samber/do"
	"gorm.io/gorm"
)

// ProvideCore provides core dependencies like config and database
func ProvideCore(i *do.Injector, cfg *config.Config, db *gorm.DB) {
	do.ProvideValue(i, cfg)
	do.ProvideValue(i, db)
}