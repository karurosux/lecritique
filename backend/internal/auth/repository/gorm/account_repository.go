package gormrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) authinterface.AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Create(account).Error
}

func (r *AccountRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) FindByEmail(ctx context.Context, email string) (*models.Account, error) {
	var account models.Account
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&account).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) Update(ctx context.Context, account *models.Account) error {
	return r.db.WithContext(ctx).Save(account).Error
}

func (r *AccountRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.Account{}, id).Error
}

func (r *AccountRepository) Exists(ctx context.Context, id uuid.UUID) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Account{}).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *AccountRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Account{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *AccountRepository) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&models.Account{}).Where("id = ?", id).Update("last_login_at", time.Now()).Error
}

func (r *AccountRepository) UpdateEmailVerification(ctx context.Context, accountID uuid.UUID, verified bool) error {
	updates := map[string]interface{}{
		"email_verified": verified,
	}
	if verified {
		updates["email_verified_at"] = time.Now()
	}
	return r.db.WithContext(ctx).Model(&models.Account{}).Where("id = ?", accountID).Updates(updates).Error
}

func (r *AccountRepository) FindAccountsPendingDeactivation(ctx context.Context) ([]models.Account, error) {
	var accounts []models.Account
	err := r.db.WithContext(ctx).Where("deactivation_requested_at IS NOT NULL AND is_active = ?", true).Find(&accounts).Error
	return accounts, err
}