package services

import (
	"context"

	"github.com/google/uuid"
	"lecritique/internal/restaurant/models"
	restaurantRepos "lecritique/internal/restaurant/repositories"
	"lecritique/internal/shared/errors"
	sharedRepos "lecritique/internal/shared/repositories"
	subscriptionRepos "lecritique/internal/subscription/repositories"
	"github.com/samber/do"
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

func NewRestaurantService(i *do.Injector) (RestaurantService, error) {
	return &restaurantService{
		restaurantRepo:   do.MustInvoke[restaurantRepos.RestaurantRepository](i),
		subscriptionRepo: do.MustInvoke[subscriptionRepos.SubscriptionRepository](i),
	}, nil
}

func (s *restaurantService) Create(ctx context.Context, accountID uuid.UUID, restaurant *models.Restaurant) error {
	// Check subscription limits
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		// For now, allow creation without subscription (development mode)
		// TODO: Re-enable subscription check in production
		// return errors.New("SUBSCRIPTION_REQUIRED", "No active subscription found for account", 402)
	}

	currentCount, err := s.restaurantRepo.CountByAccountID(ctx, accountID)
	if err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to verify restaurant count", 500)
	}

	// Only check limits if subscription exists
	if subscription != nil && !subscription.CanAddRestaurant(int(currentCount)) {
		return errors.NewWithDetails("SUBSCRIPTION_LIMIT",
			"Restaurant limit exceeded for current subscription plan",
			402,
			map[string]interface{}{
				"current_count": currentCount,
				"max_allowed":   subscription.Plan.MaxRestaurants,
			})
	}

	// Set account ID and create
	restaurant.AccountID = accountID
	if err := s.restaurantRepo.Create(ctx, restaurant); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to create restaurant", 500)
	}

	return nil
}

func (s *restaurantService) Update(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, updates map[string]interface{}) error {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		if err == sharedRepos.ErrRecordNotFound {
			return errors.NotFound("Restaurant")
		}
		return errors.ErrDatabaseOperation
	}

	if restaurant.AccountID != accountID {
		return errors.Forbidden("update this restaurant")
	}

	// Update fields
	for key, value := range updates {
		switch key {
		case "name":
			if v, ok := value.(string); ok {
				restaurant.Name = v
			}
		case "description":
			if v, ok := value.(string); ok {
				restaurant.Description = v
			}
		case "phone":
			if v, ok := value.(string); ok {
				restaurant.Phone = v
			}
		case "email":
			if v, ok := value.(string); ok {
				restaurant.Email = v
			}
		case "website":
			if v, ok := value.(string); ok {
				restaurant.Website = v
			}
		case "is_active":
			if v, ok := value.(bool); ok {
				restaurant.IsActive = v
			}
		}
	}

	if err := s.restaurantRepo.Update(ctx, restaurant); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to update restaurant", 500)
	}

	return nil
}

func (s *restaurantService) Delete(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) error {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		if err == sharedRepos.ErrRecordNotFound {
			return errors.NotFound("Restaurant")
		}
		return errors.ErrDatabaseOperation
	}

	if restaurant.AccountID != accountID {
		return errors.Forbidden("delete this restaurant")
	}

	if err := s.restaurantRepo.Delete(ctx, restaurantID); err != nil {
		return errors.Wrap(err, "DATABASE_ERROR", "Unable to delete restaurant", 500)
	}

	return nil
}

func (s *restaurantService) GetByID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		if err == sharedRepos.ErrRecordNotFound {
			return nil, errors.NotFound("Restaurant")
		}
		return nil, errors.New("DATABASE_ERROR", "Failed to fetch restaurant", 500)
	}

	if restaurant.AccountID != accountID {
		return nil, errors.Forbidden("access this restaurant")
	}

	return restaurant, nil
}

func (s *restaurantService) GetByAccountID(ctx context.Context, accountID uuid.UUID) ([]models.Restaurant, error) {
	restaurants, err := s.restaurantRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, errors.Wrap(err, "DATABASE_ERROR", "Unable to retrieve restaurants", 500)
	}
	return restaurants, nil
}
