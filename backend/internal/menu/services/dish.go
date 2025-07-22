package services

import (
	"context"
	"github.com/google/uuid"
	"lecritique/internal/menu/models"
	menuRepos "lecritique/internal/menu/repositories"
	organizationRepos "lecritique/internal/organization/repositories"
	sharedRepos "lecritique/internal/shared/repositories"
	"github.com/samber/do"
)

type ProductService interface {
	Create(ctx context.Context, accountID uuid.UUID, product *models.Product) error
	Update(ctx context.Context, accountID uuid.UUID, productID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error
	GetByID(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) (*models.Product, error)
	GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) ([]models.Product, error)
}

type productService struct {
	productRepo       menuRepos.ProductRepository
	organizationRepo organizationRepos.OrganizationRepository
}

func NewProductService(i *do.Injector) (ProductService, error) {
	return &productService{
		productRepo:       do.MustInvoke[menuRepos.ProductRepository](i),
		organizationRepo: do.MustInvoke[organizationRepos.OrganizationRepository](i),
	}, nil
}

func (s *productService) Create(ctx context.Context, accountID uuid.UUID, product *models.Product) error {
	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return err
	}

	if organization.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	return s.productRepo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, accountID uuid.UUID, productID uuid.UUID, updates map[string]interface{}) error {
	// Get product
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return err
	}

	if organization.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	// Update fields
	for key, value := range updates {
		switch key {
		case "name":
			product.Name = value.(string)
		case "description":
			product.Description = value.(string)
		case "category":
			product.Category = value.(string)
		case "price":
			product.Price = value.(float64)
		case "is_available":
			product.IsAvailable = value.(bool)
		case "is_active":
			product.IsActive = value.(bool)
		}
	}

	return s.productRepo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) error {
	// Get product
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return err
	}

	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return err
	}

	if organization.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	return s.productRepo.Delete(ctx, productID)
}

func (s *productService) GetByID(ctx context.Context, accountID uuid.UUID, productID uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, product.OrganizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return product, nil
}

func (s *productService) GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) ([]models.Product, error) {
	// Verify organization ownership
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.productRepo.FindByOrganizationID(ctx, organizationID)
}
