package providers

import (
	"github.com/samber/do"
	"gorm.io/gorm"
)

// ProvideDatabase provides the database connection as a singleton
func ProvideDatabase(i *do.Injector, db *gorm.DB) {
	do.ProvideValue(i, db)
}