package repositories

import (
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/models"
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
	*BaseRepository[models.Restaurant]
}

func NewRestaurantRepository(db *gorm.DB) RestaurantRepository {
	return &restaurantRepository{
		BaseRepository: NewBaseRepository[models.Restaurant](db),
	}
}

func (r *restaurantRepository) FindByAccountID(accountID uuid.UUID) ([]models.Restaurant, error) {
	var restaurants []models.Restaurant
	err := r.db.Where("account_id = ?", accountID).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) CountByAccountID(accountID uuid.UUID) (int64, error) {
	var count int64
	err := r.db.Model(&models.Restaurant{}).
		Where("account_id = ? AND deleted_at IS NULL", accountID).
		Count(&count).Error
	return count, err
}

type DishRepository interface {
	Create(dish *models.Dish) error
	FindByID(id uuid.UUID, preloads ...string) (*models.Dish, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.Dish, error)
	Update(dish *models.Dish) error
	Delete(id uuid.UUID) error
}

type dishRepository struct {
	*BaseRepository[models.Dish]
}

func NewDishRepository(db *gorm.DB) DishRepository {
	return &dishRepository{
		BaseRepository: NewBaseRepository[models.Dish](db),
	}
}

func (r *dishRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.Dish, error) {
	var dishes []models.Dish
	err := r.db.Where("restaurant_id = ?", restaurantID).
		Order("display_order ASC, created_at DESC").
		Find(&dishes).Error
	return dishes, err
}
