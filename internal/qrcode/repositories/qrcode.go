package repositories

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/qrcode/models"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
	"gorm.io/gorm"
)

type QRCodeRepository interface {
	Create(qrCode *models.QRCode) error
	FindByID(id uuid.UUID, preloads ...string) (*models.QRCode, error)
	FindByCode(code string) (*models.QRCode, error)
	FindByRestaurantID(restaurantID uuid.UUID) ([]models.QRCode, error)
	Update(qrCode *models.QRCode) error
	Delete(id uuid.UUID) error
	IncrementScanCount(id uuid.UUID) error
}

type qrCodeRepository struct {
	*sharedRepos.BaseRepository[models.QRCode]
}

func NewQRCodeRepository(db *gorm.DB) QRCodeRepository {
	return &qrCodeRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.QRCode](db),
	}
}

func (r *qrCodeRepository) FindByCode(code string) (*models.QRCode, error) {
	var qrCode models.QRCode
	err := r.DB.Preload("Restaurant").Preload("Location").
		Where("code = ?", code).First(&qrCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &qrCode, nil
}

func (r *qrCodeRepository) FindByRestaurantID(restaurantID uuid.UUID) ([]models.QRCode, error) {
	var qrCodes []models.QRCode
	err := r.DB.Preload("Location").
		Where("restaurant_id = ?", restaurantID).
		Order("created_at DESC").
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) IncrementScanCount(id uuid.UUID) error {
	now := time.Now()
	return r.DB.Model(&models.QRCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"scans_count":     gorm.Expr("scans_count + ?", 1),
			"last_scanned_at": now,
		}).Error
}
