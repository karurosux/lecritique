package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/qrcode/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type QRCodeRepository interface {
	Create(ctx context.Context, qrCode *models.QRCode) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.QRCode, error)
	FindByCode(ctx context.Context, code string) (*models.QRCode, error)
	FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.QRCode, error)
	Update(ctx context.Context, qrCode *models.QRCode) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementScanCount(ctx context.Context, id uuid.UUID) error
}

type qrCodeRepository struct {
	*sharedRepos.BaseRepository[models.QRCode]
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
	return &qrCodeRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.QRCode](db),
	}
}

func (r *qrCodeRepository) FindByCode(ctx context.Context, code string) (*models.QRCode, error) {
	var qrCode models.QRCode
	err := r.DB.WithContext(ctx).Preload("Restaurant").Preload("Location").
		Where("code = ?", code).First(&qrCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &qrCode, nil
}

func (r *qrCodeRepository) FindByRestaurantID(ctx context.Context, restaurantID uuid.UUID) ([]models.QRCode, error) {
	var qrCodes []models.QRCode
	err := r.DB.WithContext(ctx).Preload("Location").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) IncrementScanCount(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&models.QRCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"scans_count":     gorm.Expr("scans_count + ?", 1),
			"last_scanned_at": now,
		}).Error
}
