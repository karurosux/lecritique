package qrcodeinterface

import (
	"context"

	"github.com/google/uuid"
	qrcodemodel "kyooar/internal/qrcode/model"
)

type UpdateQRCodeRequest struct {
	IsActive *bool   `json:"is_active"`
	Label    *string `json:"label"`
	Location *string `json:"location"`
}

type QRCodeRepository interface {
	Create(ctx context.Context, qrCode *qrcodemodel.QRCode) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*qrcodemodel.QRCode, error)
	FindByIDs(ctx context.Context, ids []uuid.UUID) ([]qrcodemodel.QRCode, error)
	FindByCode(ctx context.Context, code string) (*qrcodemodel.QRCode, error)
	FindByOrganizationID(ctx context.Context, organizationID uuid.UUID) ([]qrcodemodel.QRCode, error)
	Update(ctx context.Context, qrCode *qrcodemodel.QRCode) error
	Delete(ctx context.Context, id uuid.UUID) error
	IncrementScanCount(ctx context.Context, id uuid.UUID) error
}

type QRCodeService interface {
	Generate(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, qrType qrcodemodel.QRCodeType, label string, location *string) (*qrcodemodel.QRCode, error)
	GetByCode(ctx context.Context, code string) (*qrcodemodel.QRCode, error)
	GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) ([]qrcodemodel.QRCode, error)
	Update(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID, updateReq *UpdateQRCodeRequest) (*qrcodemodel.QRCode, error)
	Delete(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID) error
	RecordScan(ctx context.Context, code string) error
}