package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/restaurant/models"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	subscriptionRepos "github.com/lecritique/api/internal/subscription/repositories"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
)

type RestaurantService interface {
	Create(ctx context.Context, accountID uuid.UUID, restaurant *models.Restaurant) error
	Update(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, updates map[string]interface{}) error
	Delete(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) error
	GetByID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) (*models.Restaurant, error)
	GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Restaurant, error)
}

type restaurantService struct {
	restaurantRepo   restaurantRepos.RestaurantRepository
	subscriptionRepo subscriptionRepos.SubscriptionRepository
}

func NewRestaurantService(restaurantRepo restaurantRepos.RestaurantRepository, subscriptionRepo subscriptionRepos.SubscriptionRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo:   restaurantRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *restaurantService) Create(ctx context.Context, accountID uuid.UUID, restaurant *models.Restaurant) error {
	// Check subscription limits
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return sharedRepos.ErrRecordNotFound
	}

	currentCount, err := s.restaurantRepo.CountByAccountID(ctx, accountID)
	if err != nil {
		return err
	}

	if !subscription.CanAddRestaurant(int(currentCount)) {
		return sharedRepos.ErrRecordNotFound
	}

	// Set account ID and create
	restaurant.AccountID = accountID
	return s.restaurantRepo.Create(ctx, restaurant)
}

func (s *restaurantService) Update(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, updates map[string]interface{}) error {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
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
			restaurant.Name = value.(string)
		case "description":
			restaurant.Description = value.(string)
		case "phone":
			restaurant.Phone = value.(string)
		case "email":
			restaurant.Email = value.(string)
		case "website":
			restaurant.Website = value.(string)
		case "is_active":
			restaurant.IsActive = value.(bool)
		}
	}

	return s.restaurantRepo.Update(ctx, restaurant)
}

func (s *restaurantService) Delete(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) error {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	return s.restaurantRepo.Delete(ctx, restaurantID)
}

func (s *restaurantService) GetByID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return restaurant, nil
}

func (s *restaurantService) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Restaurant, error) {
	return s.restaurantRepo.FindByAccountID(ctx, accountID)
}
