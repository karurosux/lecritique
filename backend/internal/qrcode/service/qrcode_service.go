package qrcodeservice

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	organizationinterface "kyooar/internal/organization/interface"
	qrcodeinterface "kyooar/internal/qrcode/interface"
	qrcodemodel "kyooar/internal/qrcode/model"
	sharedRepos "kyooar/internal/shared/repositories"
)

type qrCodeService struct {
	qrCodeRepo       qrcodeinterface.QRCodeRepository
	organizationRepo organizationinterface.OrganizationRepository
}

func NewQRCodeService(
	qrCodeRepo qrcodeinterface.QRCodeRepository,
	organizationRepo organizationinterface.OrganizationRepository,
) qrcodeinterface.QRCodeService {
	return &qrCodeService{
		qrCodeRepo:       qrCodeRepo,
		organizationRepo: organizationRepo,
	}
}

func (s *qrCodeService) Generate(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID, qrType qrcodemodel.QRCodeType, label string, location *string) (*qrcodemodel.QRCode, error) {
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

	qrCode := &qrcodemodel.QRCode{
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

func (s *qrCodeService) GetByCode(ctx context.Context, code string) (*qrcodemodel.QRCode, error) {
	qrCode, err := s.qrCodeRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if !qrCode.IsValid() {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return qrCode, nil
}

func (s *qrCodeService) GetByOrganizationID(ctx context.Context, accountID uuid.UUID, organizationID uuid.UUID) ([]qrcodemodel.QRCode, error) {
	organization, err := s.organizationRepo.FindByID(ctx, organizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.qrCodeRepo.FindByOrganizationID(ctx, organizationID)
}

func (s *qrCodeService) Update(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID, updateReq *qrcodeinterface.UpdateQRCodeRequest) (*qrcodemodel.QRCode, error) {
	qrCode, err := s.qrCodeRepo.FindByID(ctx, qrCodeID)
	if err != nil {
		return nil, err
	}

	organization, err := s.organizationRepo.FindByID(ctx, qrCode.OrganizationID)
	if err != nil {
		return nil, err
	}

	if organization.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

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

	if err := s.qrCodeRepo.Update(ctx, qrCode); err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (s *qrCodeService) Delete(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID) error {
	qrCode, err := s.qrCodeRepo.FindByID(ctx, qrCodeID)
	if err != nil {
		return err
	}

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