package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/restaurant/repositories"
	"github.com/lecritique/api/internal/subscription/models"
	subscriptionRepos "github.com/lecritique/api/internal/subscription/repositories"
)

type SubscriptionService interface {
	GetUserSubscription(ctx context.Context, accountID uuid.UUID) (*models.Subscription, error)
	GetAvailablePlans(ctx context.Context) ([]models.SubscriptionPlan, error)
	CanUserCreateRestaurant(ctx context.Context, accountID uuid.UUID) (*PermissionResponse, error)
	CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*models.Subscription, error)
	CancelSubscription(ctx context.Context, accountID uuid.UUID) error
}

type subscriptionService struct {
	subscriptionRepo     subscriptionRepos.SubscriptionRepository
	planRepo            subscriptionRepos.SubscriptionPlanRepository
	restaurantRepo      repositories.RestaurantRepository
}

func NewSubscriptionService(
	subscriptionRepo subscriptionRepos.SubscriptionRepository,
	planRepo subscriptionRepos.SubscriptionPlanRepository,
	restaurantRepo repositories.RestaurantRepository,
) SubscriptionService {
	return &subscriptionService{
		subscriptionRepo: subscriptionRepo,
		planRepo:        planRepo,
		restaurantRepo:  restaurantRepo,
	}
}

type PermissionResponse struct {
	CanCreate         bool   `json:"can_create"`
	Reason           string `json:"reason"`
	CurrentCount     int    `json:"current_count"`
	MaxAllowed       int    `json:"max_allowed"`
	SubscriptionStatus string `json:"subscription_status"`
}

type CreateSubscriptionRequest struct {
	AccountID uuid.UUID `json:"account_id"`
	PlanID    uuid.UUID `json:"plan_id"`
}

func (s *subscriptionService) GetUserSubscription(ctx context.Context, accountID uuid.UUID) (*models.Subscription, error) {
	return s.subscriptionRepo.FindByAccountID(ctx, accountID)
}

func (s *subscriptionService) GetAvailablePlans(ctx context.Context) ([]models.SubscriptionPlan, error) {
	return s.planRepo.FindAll(ctx)
}

func (s *subscriptionService) CanUserCreateRestaurant(ctx context.Context, accountID uuid.UUID) (*PermissionResponse, error) {
	// Get user's active subscription
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return &PermissionResponse{
			CanCreate:          false,
			Reason:            "No active subscription found",
			SubscriptionStatus: "none",
		}, nil
	}

	// Check if subscription is active
	if !subscription.IsActive() {
		return &PermissionResponse{
			CanCreate:          false,
			Reason:            "Subscription is not active or has expired",
			SubscriptionStatus: string(subscription.Status),
		}, nil
	}

	// Get current restaurant count
	restaurants, err := s.restaurantRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return nil, fmt.Errorf("failed to get restaurant count: %w", err)
	}

	currentCount := len(restaurants)
	maxRestaurants := subscription.Plan.Features.MaxRestaurants

	// Check if user can create more restaurants
	canCreate := subscription.CanAddRestaurant(currentCount)
	
	var reason string
	if maxRestaurants == -1 {
		reason = "Unlimited restaurants allowed"
	} else {
		reason = fmt.Sprintf("Current: %d/%d restaurants", currentCount, maxRestaurants)
	}

	return &PermissionResponse{
		CanCreate:          canCreate,
		Reason:            reason,
		CurrentCount:      currentCount,
		MaxAllowed:        maxRestaurants,
		SubscriptionStatus: string(subscription.Status),
	}, nil
}

func (s *subscriptionService) CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*models.Subscription, error) {
	// Get the plan
	plan, err := s.planRepo.FindByID(ctx, req.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	// Create subscription
	subscription := &models.Subscription{
		AccountID: req.AccountID,
		PlanID:    req.PlanID,
		Status:    models.SubscriptionActive,
		Plan:      *plan,
	}

	err = s.subscriptionRepo.Create(ctx, subscription)
	if err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	return subscription, nil
}

func (s *subscriptionService) CancelSubscription(ctx context.Context, accountID uuid.UUID) error {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = models.SubscriptionCanceled
	return s.subscriptionRepo.Update(ctx, subscription)
}