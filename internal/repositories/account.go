package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/models"
	"gorm.io/gorm"
)

type AccountRepository interface {
	Create(account *models.Account) error
	FindByID(id uuid.UUID, preloads ...string) (*models.Account, error)
	FindByEmail(email string) (*models.Account, error)
	Update(account *models.Account) error
	Delete(id uuid.UUID) error
	CountRestaurants(accountID uuid.UUID) (int64, error)
}

type accountRepository struct {
	*BaseRepository[models.Account]
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		BaseRepository: NewBaseRepository[models.Account](db),
	}
}

func (r *accountRepository) FindByEmail(email string) (*models.Account, error) {
	var account models.Account
	err := r.db.Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) CountRestaurants(accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Restaurant{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}
