package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/qrcode/models"
	qrcodeRepos "github.com/lecritique/api/internal/qrcode/repositories"
	restaurantRepos "github.com/lecritique/api/internal/restaurant/repositories"
	sharedRepos "github.com/lecritique/api/internal/shared/repositories"
)

type QRCodeService interface {
	Generate(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, qrType models.QRCodeType, label string, location *string) (*models.QRCode, error)
	GetByCode(ctx context.Context, code string) (*models.QRCode, error)
	GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) ([]models.QRCode, error)
	Delete(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID) error
	RecordScan(ctx context.Context, code string) error
}

type qrCodeService struct {
	qrCodeRepo     qrcodeRepos.QRCodeRepository
	restaurantRepo restaurantRepos.RestaurantRepository
}

func NewQRCodeService(qrCodeRepo qrcodeRepos.QRCodeRepository, restaurantRepo restaurantRepos.RestaurantRepository) QRCodeService {
	return &qrCodeService{
		qrCodeRepo:     qrCodeRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *qrCodeService) Generate(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID, qrType models.QRCodeType, label string, location *string) (*models.QRCode, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	// Generate unique code
	code, err := generateUniqueCode()
	if err != nil {
		return nil, err
	}

	// Create QR code
	qrCode := &models.QRCode{
		RestaurantID: restaurantID,
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

func (s *qrCodeService) GetByRestaurantID(ctx context.Context, accountID uuid.UUID, restaurantID uuid.UUID) ([]models.QRCode, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, sharedRepos.ErrRecordNotFound
	}

	return s.qrCodeRepo.FindByRestaurantID(ctx, restaurantID)
}

func (s *qrCodeService) Delete(ctx context.Context, accountID uuid.UUID, qrCodeID uuid.UUID) error {
	// Get QR code
	qrCode, err := s.qrCodeRepo.FindByID(ctx, qrCodeID)
	if err != nil {
		return err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(ctx, qrCode.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
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
