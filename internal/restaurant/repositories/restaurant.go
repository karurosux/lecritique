package repositories

import (
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/restaurant/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	Create(restaurant *models.Restaurant) error
	FindByID(id uuid.UUID, preloads ...string) (*models.Restaurant, error)
	FindByAccountID(accountID uuid.UUID) ([]models.Restaurant, error)
	Update(restaurant *models.Restaurant) error
	Delete(id uuid.UUID) error
	CountByAccountID(accountID uuid.UUID) (int64, error)
}

type restaurantRepository struct {
	*sharedRepos.BaseRepository[models.Restaurant]
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Restaurant](db),
	}
}

func (r *restaurantRepository) FindByAccountID(accountID uuid.UUID) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.DB.Where("account_id = ?", accountID).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) CountByAccountID(accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.Model(&models.Restaurant{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}

