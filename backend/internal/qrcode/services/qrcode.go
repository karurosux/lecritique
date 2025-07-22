package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"lecritique/internal/qrcode/models"
	qrcodeRepos "lecritique/internal/qrcode/repositories"
	organizationRepos "lecritique/internal/organization/repositories"
	sharedRepos "lecritique/internal/shared/repositories"
	"github.com/samber/do"
)

type UpdateQRCodeRequest struct {
	IsActive *bool   `json:"is_active"`
	Label    *string `json:"label"`
	Location *string `json:"location"`
}

type QRCodeService interface {
	Generate(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, qrType models.QRCodeType, label string, location *string) (*models.QRCode, error)
	GetByCode(ctx context.Context, code string) (*models.QRCode, error)
	GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) ([]models.QRCode, error)
	Update(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID, updateReq *UpdateQRCodeRequest) (*models.QRCode, error)
	Delete(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID) error
	RecordScan(ctx context.Context, code string) error
}

type qrCodeService struct {
	qrCodeRepo     qrcodeRepos.QRCodeRepository
	organizationRepo organizationRepos.OrganizationRepository
}

func NewQRCodeService(i *do.Injector) (QRCodeService, error) {
	return &qrCodeService{
		qrCodeRepo:     do.MustInvoke[qrcodeRepos.QRCodeRepository](i),
		organizationRepo: do.MustInvoke[organizationRepos.OrganizationRepository](i),
	}, nil
}

func (s *qrCodeService) Generate(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, qrType models.QRCodeType, label string, location *string) (*models.QRCode, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	code, err := generateUniqueCode()
	if err != nil {
		return nil, err
	}

	qrCode := &models.QRCode{
		OrganizationID: organizationID,
		Code:         code,
		Type:         qrType,
		Label:        label,
		Location:     location,
		IsActive:     true,
	}

	if err := s.qrCodeRepo.Create(ctx, qrCode); err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (s *qrCodeService) GetByCode(ctx context.Context, code string) (*models.QRCode, error) {
	qrCode, err := s.qrCodeRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if !qrCode.IsValid() {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return qrCode, nil
}

func (s *qrCodeService) GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) ([]models.QRCode, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.qrCodeRepo.FindByOrganizationID(ctx, organizationID)
}

func (s *qrCodeService) Update(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID, updateReq *UpdateQRCodeRequest) (*models.QRCode, error) {
	// Get QR code
	qrCode, err := s.qrCodeRepo.FindByID(ctx, qrCodeID)
	if err != nil {
		return nil, err
	}

	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, qrCode.OrganizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	// Update fields if provided
	if updateReq.IsActive != nil {
		qrCode.IsActive = *updateReq.IsActive
	}
	if updateReq.Label != nil {
		qrCode.Label = *updateReq.Label
	}
	if updateReq.Location != nil {
		qrCode.Location = updateReq.Location
	}

	qrCode.UpdatedAt = time.Now()

	// Save to repository
	if err := s.qrCodeRepo.Update(ctx, qrCode); err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (s *qrCodeService) Delete(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID) error {
	// Get QR code
	qrCode, err := s.qrCodeRepo.FindByID(ctx, qrCodeID)
	if err != nil {
		return err
	}

	// Verify ownership
	organization, err := s.organizationRepo.FindByID(ctx, qrCode.OrganizationID)
	if err != nil {
		return err
	}

	if organization.AccountID != accountID {
		return sharedRepos.ErrRecordNotFound
	}

	return s.qrCodeRepo.Delete(ctx, qrCodeID)
}

func (s *qrCodeService) RecordScan(ctx context.Context, code string) error {
	qrCode, err := s.qrCodeRepo.FindByCode(ctx, code)
	if err != nil {
		return err
	}

	if !qrCode.IsValid() {
		return sharedRepos.ErrRecordNotFound
	}

	return s.qrCodeRepo.IncrementScanCount(ctx, qrCode.ID)
}

func generateUniqueCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("LCQ-%s-%d", hex.EncodeToString(bytes)[:8], time.Now().Unix()), nil
}
