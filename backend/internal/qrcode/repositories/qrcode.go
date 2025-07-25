package repositories

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"kyooar/internal/qrcode/models"
	sharedRepos "kyooar/internal/shared/repositories"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type QRCodeRepository interface {
	Create(ctx context.Context, qrCode *models.QRCode) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*models.QRCode, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.QRCode, error)
	FindByCode(ctx context.Context, code string) (*models.QRCode, error)
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.QRCode, error)
	Update(ctx context.Context, qrCode *models.QRCode) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementScanCount(ctx context.Context, id uuid.UUID) error
}

type qrCodeRepository struct {
	*sharedRepos.BaseRepository[models.QRCode]
}

func NewQRCodeRepository(i *do.Injector) (QRCodeRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &qrCodeRepository{
		BaseRepository: sharedRepos.NewBaseRepository[models.QRCode](db),
	}, nil
}

func (r *qrCodeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]models.QRCode, error) {
	if len(ids) == 0 {
		return []models.QRCode{}, nil
	}
	
	var qrCodes []models.QRCode
	err := r.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) FindByCode(ctx context.Context, code string) (*models.QRCode, error) {
	var qrCode models.QRCode
	err := r.DB.WithContext(ctx).Preload("Organization").
		Where("code = ?", code).First(&qrCode).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, sharedRepos.ErrRecordNotFound
		}
		return nil, err
	}
	return &qrCode, nil
}

func (r *qrCodeRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]models.QRCode, error) {
	var qrCodes []models.QRCode
	err := r.DB.WithContext(ctx).
		Where("organization_id = ?", organizationID).
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
