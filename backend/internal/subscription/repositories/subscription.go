package repositories

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"kyooar/internal/subscription/models"
	sharedRepos "kyooar/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *models.Subscription) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.Subscription, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) (*models.Subscription, error)
	FindByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*models.Subscription, error)
	Update(ctx context.Context, subscription *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type subscriptionRepository struct {
	*sharedRepos.BaseRepository[models.Subscription]
}

func NewSubscriptionRepository(i *do.Injector) (SubscriptionRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &subscriptionRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Subscription](db),
	}, nil
}

func (r *subscriptionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.DB.WithContext(ctx).Preload("Plan").Where("account_id = ?", accountID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) FindByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.DB.WithContext(ctx).Preload("Plan").Where("stripe_subscription_id = ?", stripeSubscriptionID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

type SubscriptionPlanRepository interface {
	FindAll(ctx context.Context) ([]models.SubscriptionPlan, error)
	FindAllIncludingHidden(ctx context.Context) ([]models.SubscriptionPlan, error)
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.SubscriptionPlan, error)
	FindByCode(ctx context.Context, code string) (*models.SubscriptionPlan, error)
}

type subscriptionPlanRepository struct {
	*sharedRepos.BaseRepository[models.SubscriptionPlan]
}

func NewSubscriptionPlanRepository(i *do.Injector) (SubscriptionPlanRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &subscriptionPlanRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.SubscriptionPlan](db),
	}, nil
}

func (r *subscriptionPlanRepository) FindAll(ctx context.Context) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	// Only return visible and active plans for public listing
	err := r.DB.WithContext(ctx).Where("is_active = ? AND is_visible = ?", true, true).Order("price ASC").Find(&plans).Error
	return plans, err
}

func (r *subscriptionPlanRepository) FindAllIncludingHidden(ctx context.Context) ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	// Return all active plans, including hidden ones (for admin use)
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).Order("is_visible DESC, price ASC").Find(&plans).Error
	return plans, err
}

func (r *subscriptionPlanRepository) FindByCode(ctx context.Context, code string) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &plan, nil
}
