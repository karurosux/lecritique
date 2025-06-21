package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/auth/models"
	restaurantModels "github.com/lecritique/api/internal/restaurant/models"
	"github.com/lecritique/api/internal/shared/repositories"
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
	*repositories.BaseRepository[models.Account]
}

func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{
		BaseRepository: repositories.NewBaseRepository[models.Account](db),
	}
}

func (r *accountRepository) FindByEmail(email string) (*models.Account, error) {
	var account models.Account
	err := r.DB.Where("email = ?", email).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, repositories.ErrRecordNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *accountRepository) CountRestaurants(accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&restaurantModels.Restaurant{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}
