package services

import (
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/shared/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"github.com/lecritique/api/internal/shared/errors"
)

type DishService interface {
	Create(accountID uuid.UUID, dish *models.Dish) error
	Update(accountID uuid.UUID, dishID uuid.UUID, updates map[string]interface{}) error
	Delete(accountID uuid.UUID, dishID uuid.UUID) error
	GetByID(accountID uuid.UUID, dishID uuid.UUID) (*models.Dish, error)
	GetByRestaurantID(accountID uuid.UUID, restaurantID uuid.UUID) ([]models.Dish, error)
}

type dishService struct {
	dishRepo       repositories.DishRepository
	restaurantRepo repositories.RestaurantRepository
}

func NewDishService(dishRepo repositories.DishRepository, restaurantRepo repositories.RestaurantRepository) DishService {
	return &dishService{
		dishRepo:       dishRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *dishService) Create(accountID uuid.UUID, dish *models.Dish) error {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(dish.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return errors.ErrForbidden
	}

	return s.dishRepo.Create(dish)
}

func (s *dishService) Update(accountID uuid.UUID, dishID uuid.UUID, updates map[string]interface{}) error {
	// Get dish
	dish, err := s.dishRepo.FindByID(dishID)
	if err != nil {
		return err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(dish.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return errors.ErrForbidden
	}

	// Update fields
	for key, value := range updates {
		switch key {
		case "name":
			dish.Name = value.(string)
		case "description":
			dish.Description = value.(string)
		case "category":
			dish.Category = value.(string)
		case "price":
			dish.Price = value.(float64)
		case "is_available":
			dish.IsAvailable = value.(bool)
		case "is_active":
			dish.IsActive = value.(bool)
		}
	}

	return s.dishRepo.Update(dish)
}

func (s *dishService) Delete(accountID uuid.UUID, dishID uuid.UUID) error {
	// Get dish
	dish, err := s.dishRepo.FindByID(dishID)
	if err != nil {
		return err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(dish.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return errors.ErrForbidden
	}

	return s.dishRepo.Delete(dishID)
}

func (s *dishService) GetByID(accountID uuid.UUID, dishID uuid.UUID) (*models.Dish, error) {
	dish, err := s.dishRepo.FindByID(dishID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(dish.RestaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, errors.ErrForbidden
	}

	return dish, nil
}

func (s *dishService) GetByRestaurantID(accountID uuid.UUID, restaurantID uuid.UUID) ([]models.Dish, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, errors.ErrForbidden
	}

	return s.dishRepo.FindByRestaurantID(restaurantID)
}
