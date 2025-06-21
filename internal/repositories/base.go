package repositories

import (
	"errors"
	"math"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidID      = errors.New("invalid ID")
)

type BaseRepository[T any] struct {
	db *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) *BaseRepository[T] {
	return &BaseRepository[T]{db: db}
}

func (r *BaseRepository[T]) Create(entity *T) error {
	return r.db.Create(entity).Error
}

func (r *BaseRepository[T]) FindByID(id uuid.UUID, preloads ...string) (*T, error) {
	var entity T
	query := r.db
	
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	err := query.First(&entity, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}
	
	return &entity, nil
}

func (r *BaseRepository[T]) Update(entity *T) error {
	return r.db.Save(entity).Error
}

func (r *BaseRepository[T]) Delete(id uuid.UUID) error {
	var entity T
	result := r.db.Delete(&entity, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

func (r *BaseRepository[T]) SoftDelete(id uuid.UUID) error {
	var entity T
	result := r.db.Where("id = ?", id).Delete(&entity)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

type PageRequest struct {
	Page  int
	Limit int
	Sort  string
	Order string
}

type PageResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

func (r *BaseRepository[T]) FindAll(req PageRequest, preloads ...string) (*PageResponse[T], error) {
	var entities []T
	var total int64
	
	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 {
		req.Limit = 20
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	
	// Count total
	r.db.Model(new(T)).Count(&total)
	
	// Build query
	query := r.db.Limit(req.Limit).Offset((req.Page - 1) * req.Limit)
	
	for _, preload := range preloads {
		query = query.Preload(preload)
	}
	
	if req.Sort != "" && req.Order != "" {
		query = query.Order(req.Sort + " " + req.Order)
	}
	
	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	
	return &PageResponse[T]{
		Data:       entities,
		Page:       req.Page,
		Limit:      req.Limit,
		Total:      int(total),
		TotalPages: int(math.Ceil(float64(total) / float64(req.Limit))),
	}, nil
}
