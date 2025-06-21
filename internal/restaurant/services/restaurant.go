package services

import (
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/shared/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"github.com/lecritique/api/internal/shared/errors"
)

type RestaurantService interface {
	Create(accountID uuid.UUID, restaurant *models.Restaurant) error
	Update(accountID uuid.UUID, restaurantID uuid.UUID, updates map[string]interface{}) error
	Delete(accountID uuid.UUID, restaurantID uuid.UUID) error
	GetByID(accountID uuid.UUID, restaurantID uuid.UUID) (*models.Restaurant, error)
	GetByAccountID(accountID uuid.UUID) ([]models.Restaurant, error)
}

type restaurantService struct {
	restaurantRepo   repositories.RestaurantRepository
	subscriptionRepo repositories.SubscriptionRepository
}

func NewRestaurantService(restaurantRepo repositories.RestaurantRepository, subscriptionRepo repositories.SubscriptionRepository) RestaurantService {
	return &restaurantService{
		restaurantRepo:   restaurantRepo,
		subscriptionRepo: subscriptionRepo,
	}
}

func (s *restaurantService) Create(accountID uuid.UUID, restaurant *models.Restaurant) error {
	// Check subscription limits
	subscription, err := s.subscriptionRepo.FindByAccountID(accountID)
	if err != nil {
		return errors.ErrSubscriptionLimit
	}

	currentCount, err := s.restaurantRepo.CountByAccountID(accountID)
	if err != nil {
		return err
	}

	if !subscription.CanAddRestaurant(int(currentCount)) {
		return errors.ErrSubscriptionLimit
	}

	// Set account ID and create
	restaurant.AccountID = accountID
	return s.restaurantRepo.Create(restaurant)
}

func (s *restaurantService) Update(accountID uuid.UUID, restaurantID uuid.UUID, updates map[string]interface{}) error {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
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

	return s.restaurantRepo.Update(restaurant)
}

func (s *restaurantService) Delete(accountID uuid.UUID, restaurantID uuid.UUID) error {
	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return errors.ErrForbidden
	}

	return s.restaurantRepo.Delete(restaurantID)
}

func (s *restaurantService) GetByID(accountID uuid.UUID, restaurantID uuid.UUID) (*models.Restaurant, error) {
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, errors.ErrForbidden
	}

	return restaurant, nil
}

func (s *restaurantService) GetByAccountID(accountID uuid.UUID) ([]models.Restaurant, error) {
	return s.restaurantRepo.FindByAccountID(accountID)
}
