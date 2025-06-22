package repositories

import (
	"context"
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/menu/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type DishRepository interface {
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Dish, error)
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Dish, error)
	Create(ctx context.Context, dish *models.Dish) error
	Update(ctx context.Context, dish *models.Dish) error
	Delete(ctx context.Context, id uuid.UUID) error
	FindAll(ctx context.Context, limit, offset int) ([]models.Dish, error)
	Count(ctx context.Context) (int64, error)
}

type dishRepository struct {
	*repositories.BaseRepository[models.Dish]
}

func NewDishRepository(db *gorm.DB) DishRepository {
	return &dishRepository{
		BaseRepository: repositories.NewBaseRepository[models.Dish](db),
	}
}

func (r *dishRepository) FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.Dish, error) {
	var dishes []models.Dish
	err := r.DB.WithContext(ctx).Where("restaurant_id = ? AND is_active = ?", restaurantID, true).
		Order("display_order ASC, name ASC").
		Find(&dishes).Error
	return dishes, err
}