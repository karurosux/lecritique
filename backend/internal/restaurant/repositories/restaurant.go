package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/restaurant/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type RestaurantRepository interface {
	Create(ctx context.Context, restaurant *models.Restaurant) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Restaurant, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Restaurant, error)
	Update(ctx context.Context, restaurant *models.Restaurant) error
	Delete(ctx context.Context, id uuid.UUID) error
	CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error)
}

type restaurantRepository struct {
	*sharedRepos.BaseRepository[models.Restaurant]
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Restaurant](db),
	}
}

func (r *restaurantRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.DB.WithContext(ctx).Where("account_id = ?", accountID).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) CountByAccountID(ctx context.Context, accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.DB.WithContext(ctx).Model(&models.Restaurant{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}

