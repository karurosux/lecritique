package repositories

import (
	"errors"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/subscription/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type SubscriptionRepository interface {
	Create(subscription *models.Subscription) error
	FindByID(id uuid.UUID, preloads ...string) (*models.Subscription, error)
	FindByAccountID(accountID uuid.UUID) (*models.Subscription, error)
	Update(subscription *models.Subscription) error
}

type subscriptionRepository struct {
	*sharedRepos.BaseRepository[models.Subscription]
}

func NewSubscriptionRepository(db *gorm.DB) SubscriptionRepository {
	return &subscriptionRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.Subscription](db),
	}
}

func (r *subscriptionRepository) FindByAccountID(accountID uuid.UUID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.DB.Preload("Plan").Where("account_id = ?", accountID).First(&subscription).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &subscription, nil
}

type SubscriptionPlanRepository interface {
	FindAll() ([]models.SubscriptionPlan, error)
	FindByID(id uuid.UUID, preloads ...string) (*models.SubscriptionPlan, error)
	FindByCode(code string) (*models.SubscriptionPlan, error)
}

type subscriptionPlanRepository struct {
	*sharedRepos.BaseRepository[models.SubscriptionPlan]
}

func NewSubscriptionPlanRepository(db *gorm.DB) SubscriptionPlanRepository {
	return &subscriptionPlanRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.SubscriptionPlan](db),
	}
}

func (r *subscriptionPlanRepository) FindAll() ([]models.SubscriptionPlan, error) {
	var plans []models.SubscriptionPlan
	err := r.DB.Where("is_active = ?", true).Order("price ASC").Find(&plans).Error
	return plans, err
}

func (r *subscriptionPlanRepository) FindByCode(code string) (*models.SubscriptionPlan, error) {
	var plan models.SubscriptionPlan
	err := r.DB.Where("code = ?", code).First(&plan).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &plan, nil
}
