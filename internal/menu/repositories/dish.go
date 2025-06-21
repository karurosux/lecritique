package repositories

import (
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/menu/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type DishRepository interface {
	FindByID(id uuid.UUID, preloads ...string) (*models.Dish, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.Dish, error)
	Create(dish *models.Dish) error
	Update(dish *models.Dish) error
	Delete(id uuid.UUID) error
	FindAll(limit, offset int) ([]models.Dish, error)
	Count() (int64, error)
}

type dishRepository struct {
	*repositories.BaseRepository[models.Dish]
}

func NewDishRepository(db *gorm.DB) DishRepository {
	return &dishRepository{
		BaseRepository: repositories.NewBaseRepository[models.Dish](db),
	}
}

func (r *dishRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.Dish, error) {
	var dishes []models.Dish
	err := r.DB.Where("restaurant_id = ? AND is_active = ?", restaurantID, true).
		Order("display_order ASC, name ASC").
		Find(&dishes).Error
	return dishes, err
}