package repositories

import (
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

func (r *BaseRepository[T]) FindByID(id uuid.UUID, preloads ...string) (*T, error) {
	var entity T
	query := r.DB
	
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	err := query.First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.DB.Create(entity).Error
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.DB.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(id uuid.UUID) error {
	var entity T
	return r.DB.Delete(&entity, "id = ?", id).Error
}

func (r *BaseRepository[T]) FindAll(limit, offset int) ([]T, error) {
	var entities []T
	query := r.DB.Limit(limit).Offset(offset)
	err := query.Find(&entities).Error
	return entities, err
}

func (r *BaseRepository[T]) Count() (int64, error) {
	var entity T
	var count int64
	err := r.DB.Model(&entity).Count(&count).Error
	return count, err
}
