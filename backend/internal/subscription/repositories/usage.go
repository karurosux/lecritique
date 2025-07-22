package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/subscription/models"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type UsageRepository interface {
	Create(ctx context.Context, usage *models.SubscriptionUsage) error
	Update(ctx context.Context, usage *models.SubscriptionUsage) error
	FindByID(ctx context.Context, id uuid.UUID) (*models.SubscriptionUsage, error)
	FindBySubscriptionAndPeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) (*models.SubscriptionUsage, error)
	FindBySubscription(ctx context.Context, subscriptionID uuid.UUID) ([]*models.SubscriptionUsage, error)
	CreateEvent(ctx context.Context, event *models.UsageEvent) error
	FindEventsBySubscription(ctx context.Context, subscriptionID uuid.UUID, limit int) ([]*models.UsageEvent, error)
}

type usageRepository struct {
	db *gorm.DB
}

func NewUsageRepository(i *do.Injector) (UsageRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &usageRepository{db: db}, nil
}

func (r *usageRepository) Create(ctx context.Context, usage *models.SubscriptionUsage) error {
	return r.db.WithContext(ctx).Create(usage).Error
}

func (r *usageRepository) Update(ctx context.Context, usage *models.SubscriptionUsage) error {
	return r.db.WithContext(ctx).Save(usage).Error
}

func (r *usageRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.SubscriptionUsage, error) {
	var usage models.SubscriptionUsage
	err := r.db.WithContext(ctx).Preload("Subscription").First(&usage, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (r *usageRepository) FindBySubscriptionAndPeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) (*models.SubscriptionUsage, error) {
	var usage models.SubscriptionUsage
	err := r.db.WithContext(ctx).
		Where("subscription_id = ? AND period_start = ? AND period_end = ?", subscriptionID, periodStart, periodEnd).
		First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (r *usageRepository) FindBySubscription(ctx context.Context, subscriptionID uuid.UUID) ([]*models.SubscriptionUsage, error) {
	var usages []*models.SubscriptionUsage
	err := r.db.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		Order("period_start DESC").
		Find(&usages).Error
	return usages, err
}

func (r *usageRepository) CreateEvent(ctx context.Context, event *models.UsageEvent) error {
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *usageRepository) FindEventsBySubscription(ctx context.Context, subscriptionID uuid.UUID, limit int) ([]*models.UsageEvent, error) {
	var events []*models.UsageEvent
	query := r.db.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&events).Error
	return events, err
}