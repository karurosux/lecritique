package repositories

import (
	"context"
	"github.com/google/uuid"
	"lecritique/internal/organization/models"
	sharedRepos "lecritique/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type OrganizationRepository interface {
	Create(ctx context.Context, organization *models.Organization) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Organization, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Organization, error)
	Update(ctx context.Context, organization *models.Organization) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
}

type organizationRepository struct {
	*sharedRepos.BaseRepository[models.Organization]
}

func NewOrganizationRepository(i *do.Injector) (OrganizationRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &organizationRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Organization](db),
	}, nil
}

func (r *organizationRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Organization, error) {
	var organizations []models.Organization
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Find(&organizations).Error
	return organizations, err
}

func (r *organizationRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Organization{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}

