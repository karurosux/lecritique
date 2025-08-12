package gormrepo

import (
	"context"
	"errors"

	"github.com/google/uuid"
	subscriptioninterface "kyooar/internal/subscription/interface"
	subscriptionmodel "kyooar/internal/subscription/model"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type subscriptionRepository struct {
	*sharedRepos.BaseRepository[subscriptionmodel.Subscription]
}

func NewSubscriptionRepository(db *gorm.DB) subscriptioninterface.SubscriptionRepository {
	return &subscriptionRepository{
		BaseRepository: sharedRepos.NewBaseRepository[subscriptionmodel.Subscription](db),
	}
}

func (r *subscriptionRepository) FindByAccountID(ctx context.Context, accountID uuid.UUID) (*subscriptionmodel.Subscription, error) {
	var subscription subscriptionmodel.Subscription
	err := r.DB.WithContext(ctx).Preload("Plan").Where("account_id = ?", accountID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

func (r *subscriptionRepository) FindByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*subscriptionmodel.Subscription, error) {
	var subscription subscriptionmodel.Subscription
	err := r.DB.WithContext(ctx).Preload("Plan").Where("stripe_subscription_id = ?", stripeSubscriptionID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &subscription, nil
}