package gormqrcode

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	qrcodeinterface "kyooar/internal/qrcode/interface"
	qrcodemodel "kyooar/internal/qrcode/model"
	sharedRepos "kyooar/internal/shared/repositories"
	"gorm.io/gorm"
)

type qrCodeRepository struct {
	*sharedRepos.BaseRepository[qrcodemodel.QRCode]
}

func NewQRCodeRepository(db *gorm.DB) qrcodeinterface.QRCodeRepository {
	return &qrCodeRepository{
		BaseRepository: sharedRepos.NewBaseRepository[qrcodemodel.QRCode](db),
	}
}

func (r *qrCodeRepository) FindByIDs(ctx context.Context, ids []uuid.UUID) ([]qrcodemodel.QRCode, error) {
	if len(ids) == 0 {
		return []qrcodemodel.QRCode{}, nil
	}
	
	var qrCodes []qrcodemodel.QRCode
	err := r.DB.WithContext(ctx).
		Where("id IN ?", ids).
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) FindByCode(ctx context.Context, code string) (*qrcodemodel.QRCode, error) {
	var qrCode qrcodemodel.QRCode
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

func (r *qrCodeRepository) FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]qrcodemodel.QRCode, error) {
	var qrCodes []qrcodemodel.QRCode
	err := r.DB.WithContext(ctx).
		Where("organization_id = ?", organizationID).
		Order("created_at DESC").
		Find(&qrCodes).Error
	return qrCodes, err
}

func (r *qrCodeRepository) IncrementScanCount(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.DB.WithContext(ctx).Model(&qrcodemodel.QRCode{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"scans_count":     gorm.Expr("scans_count + ?", 1),
			"last_scanned_at": now,
		}).Error
}