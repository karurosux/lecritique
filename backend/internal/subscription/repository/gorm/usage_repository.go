package gormrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	subscriptioninterface "kyooar/internal/subscription/interface"
	subscriptionmodel "kyooar/internal/subscription/model"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type usageRepository struct {
	*sharedRepos.BaseRepository[subscriptionmodel.SubscriptionUsage]
}

func NewUsageRepository(db *gorm.DB) subscriptioninterface.UsageRepository {
	return &usageRepository{
		BaseRepository: sharedRepos.NewBaseRepository[subscriptionmodel.SubscriptionUsage](db),
	}
}

func (r *usageRepository) FindBySubscriptionAndPeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) (*subscriptionmodel.SubscriptionUsage, error) {
	var usage subscriptionmodel.SubscriptionUsage
	err := r.DB.WithContext(ctx).
		Where("subscription_id = ? AND period_start = ? AND period_end = ?", subscriptionID, periodStart, periodEnd).
		First(&usage).Error
	if err != nil {
		return nil, err
	}
	return &usage, nil
}

func (r *usageRepository) CreateEvent(ctx context.Context, event *subscriptionmodel.UsageEvent) error {
	return r.DB.WithContext(ctx).Create(event).Error
}

func (r *usageRepository) FindEventsBySubscription(ctx context.Context, subscriptionID uuid.UUID, limit int) ([]subscriptionmodel.UsageEvent, error) {
	var events []subscriptionmodel.UsageEvent
	query := r.DB.WithContext(ctx).
		Where("subscription_id = ?", subscriptionID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	
	err := query.Find(&events).Error
	return events, err
}

func (r *usageRepository) FindEventsByPeriod(ctx context.Context, subscriptionID uuid.UUID, start, end time.Time) ([]subscriptionmodel.UsageEvent, error) {
	var events []subscriptionmodel.UsageEvent
	err := r.DB.WithContext(ctx).
		Where("subscription_id = ? AND created_at >= ? AND created_at <= ?", subscriptionID, start, end).
		Order("created_at DESC").
		Find(&events).Error
	return events, err
}

func (r *usageRepository) ResetMonthlyUsage(ctx context.Context) error {
	now := time.Now()
	firstOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	
	return r.DB.WithContext(ctx).Model(&subscriptionmodel.SubscriptionUsage{}).
		Where("period_start = ?", firstOfMonth).
		Updates(map[string]interface{}{
			"feedbacks_count": 0,
			"last_updated_at": now,
		}).Error
}