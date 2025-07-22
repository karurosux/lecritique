package repositories

import (
	"context"
	"github.com/google/uuid"
	"lecritique/internal/menu/models"
	"lecritique/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Product, error)
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, limit, offset int) ([]models.Product, error)
	Count(ctx context.Context) (int64, error)
}

type productRepository struct {
	*repositories.BaseRepository[models.Product]
}

func NewProductRepository(i *do.Injector) (ProductRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &productRepository{
		BaseRepository: repositories.NewBaseRepository[models.Product](db),
	}, nil
}

func (r *productRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.WithContext(ctx).Where("organization_id = ? AND is_active = ?", organizationID, true).
		Order("display_order ASC, name ASC").
		Find(&products).Error
	return products, err
}
