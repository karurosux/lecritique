package gormrepo

import (
	"context"
	"errors"

	subscriptioninterface "kyooar/internal/subscription/interface"
	subscriptionmodel "kyooar/internal/subscription/model"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type subscriptionPlanRepository struct {
	*sharedRepos.BaseRepository[subscriptionmodel.SubscriptionPlan]
}

func NewSubscriptionPlanRepository(db *gorm.DB) subscriptioninterface.SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		BaseRepository: sharedRepos.NewBaseRepository[subscriptionmodel.SubscriptionPlan](db),
	}
}

func (r *subscriptionPlanRepository) FindAll(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error) {
	var plans []subscriptionmodel.SubscriptionPlan
	err := r.DB.WithContext(ctx).Where("is_active = ? AND is_visible = ?", true, true).Order("price ASC").Find(&plans).Error
	return plans, err
}

func (r *subscriptionPlanRepository) FindAllIncludingHidden(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error) {
	var plans []subscriptionmodel.SubscriptionPlan
	err := r.DB.WithContext(ctx).Where("is_active = ?", true).Order("is_visible DESC, price ASC").Find(&plans).Error
	return plans, err
}

func (r *subscriptionPlanRepository) FindByCode(ctx context.Context, code string) (*subscriptionmodel.SubscriptionPlan, error) {
	var plan subscriptionmodel.SubscriptionPlan
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &plan, nil
}