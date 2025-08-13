package product

import (
	"github.com/samber/do"

	productHandlers "kyooar/internal/product/handlers"
	productRepos "kyooar/internal/product/repositories"
	productServices "kyooar/internal/product/services"
)

func RegisterNewModule(container *do.Injector) {
	do.Provide(container, productRepos.NewProductRepository)
	do.Provide(container, productServices.NewProductService)
	do.Provide(container, productHandlers.NewProductHandler)
	do.Provide(container, productHandlers.NewProductPublicHandler)
}