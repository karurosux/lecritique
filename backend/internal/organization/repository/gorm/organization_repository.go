package gormrepo

import (
	"context"
	"github.com/google/uuid"
	organizationinterface "kyooar/internal/organization/interface"
	organizationmodel "kyooar/internal/organization/model"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type organizationRepository struct {
	*sharedRepos.BaseRepository[organizationmodel.Organization]
}

func NewOrganizationRepository(db *gorm.DB) organizationinterface.OrganizationRepository {
	return &organizationRepository{
		BaseRepository: sharedRepos.NewBaseRepository[organizationmodel.Organization](db),
	}
}

func (r *organizationRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]organizationmodel.Organization, error) {
	var organizations []organizationmodel.Organization
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Find(&organizations).Error
	return organizations, err
}

func (r *organizationRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&organizationmodel.Organization{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}

