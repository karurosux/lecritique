package repositories

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrRecordNotFound = errors.New("record not found")

type BaseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{DB: db}
}

func (r *BaseRepository[T]) FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*T, error) {
	var entity T
	query := r.DB.WithContext(ctx)
	
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	err := query.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

func (r *BaseRepository[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

func (r *BaseRepository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

func (r *BaseRepository[T]) FindAll(ctx context.Context, limit, offset int) ([]T, error) {
	var entities []T
	query := r.DB.WithContext(ctx).Limit(limit).Offset(offset)
	err := query.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Count(ctx context.Context) (int64, error) {
	var entity T
	var count int64
	err := r.DB.WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}
