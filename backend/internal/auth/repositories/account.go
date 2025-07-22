package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"lecritique/internal/auth/models"
	organizationModels "lecritique/internal/organization/models"
	"lecritique/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(ctx context.Context, account *models.Account) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Account, error)
	FindByEmail(ctx context.Context, email string) (*models.Account, error)
	Update(ctx context.Context, account *models.Account) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountOrganizations(ctx context.Context, accountID uuid.UUID) (int64, error)
	UpdateEmailVerification(ctx context.Context, accountID uuid.UUID, verified bool) error
	FindAccountsPendingDeactivation(ctx context.Context) ([]models.Account, error)
}

type accountRepository struct {
	*repositories.BaseRepository[models.Account]
}

func NewAccountRepository(i *do.Injector) (AccountRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &accountRepository{
		BaseRepository: repositories.NewBaseRepository[models.Account](db),
	}, nil
}

func (r *accountRepository) FindByEmail(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account
	err := r.DB.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrRecordNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) CountOrganizations(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&organizationModels.Organization{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}

func (r *accountRepository) UpdateEmailVerification(ctx context.Context, accountID uuid.UUID, verified bool) error {
	updates := map[string]interface{}{
		"email_verified": verified,
	}
	if verified {
		now := time.Now()
		updates["email_verified_at"] = &now
	}
	
	err := r.DB.WithContext(ctx).Model(&models.Account{}).
		Where("id = ?", accountID).
		Updates(updates).Error
	return err
}

func (r *accountRepository) FindAccountsPendingDeactivation(ctx context.Context) ([]models.Account, error) {
	var accounts []models.Account
	err := r.DB.WithContext(ctx).
		Where("deactivation_requested_at IS NOT NULL AND is_active = ?", true).
		Find(&accounts).Error
	if err != nil {
		return nil, err
	}
	return accounts, nil
}
