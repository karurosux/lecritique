package services

import (
	"context"
	"github.com/google/uuid"
	"lecritique/internal/menu/models"
	menuRepos "lecritique/internal/menu/repositories"
	restaurantRepos "lecritique/internal/restaurant/repositories"
	sharedRepos "lecritique/internal/shared/repositories"
	"github.com/samber/do"
)

type DishService interface {
	Create(ctx context.Context, accountID uuid.UUID, dish *models.Dish) error
	Update(ctx context.Context, accountID uuid.UUID, dishID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, accountID uuid.UUID, dishID uuid.UUID) error
	GetByID(ctx context.Context, accountID uuid.UUID, dishID uuid.UUID) (*models.Dish, error)
	GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) ([]models.Dish, error)
}

type dishService struct {
	dishRepo       menuRepos.DishRepository
	restaurantRepo restaurantRepos.RestaurantRepository
}

func NewDishService(i *do.Injector) (DishService, error) {
	return &dishService{
		dishRepo:       do.MustInvoke[menuRepos.DishRepository](i),
		restaurantRepo: do.MustInvoke[restaurantRepos.RestaurantRepository](i),
	}, nil
}

func (s *dishService) Create(ctx context.Context, accountID uuid.UUID, dish *models.Dish) error {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	return s.dishRepo.Create(ctx, dish)
}

func (s *dishService) Update(ctx context.Context, accountID uuid.UUID, dishID uuid.UUID, updates map[string]interface{}) error {
	// Get dish
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
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

	return s.dishRepo.Update(ctx, dish)
}

func (s *dishService) Delete(ctx context.Context, accountID uuid.UUID, dishID uuid.UUID) error {
	// Get dish
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	return s.dishRepo.Delete(ctx, dishID)
}

func (s *dishService) GetByID(ctx context.Context, accountID uuid.UUID, dishID uuid.UUID) (*models.Dish, error) {
	dish, err := s.dishRepo.FindByID(ctx, dishID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, dish.RestaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return dish, nil
}

func (s *dishService) GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) ([]models.Dish, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.dishRepo.FindByRestaurantID(ctx, restaurantID)
}
