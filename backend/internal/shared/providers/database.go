package providers

import (
	"github.com/samber/do"
	"gorm.io/gorm"
)

func ProvideDatabase(i *do.Injector, db *gorm.DB) {
	do.ProvideValue(i, db)
}