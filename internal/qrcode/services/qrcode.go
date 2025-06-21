package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lecritique/api/internal/shared/models"
	"github.com/lecritique/api/internal/shared/repositories"
	"github.com/lecritique/api/internal/shared/errors"
)

type QRCodeService interface {
	Generate(accountID uuid.UUID, restaurantID uuid.UUID, qrType models.QRCodeType, label string) (*models.QRCode, error)
	GetByCode(code string) (*models.QRCode, error)
	GetByRestaurantID(accountID uuid.UUID, restaurantID uuid.UUID) ([]models.QRCode, error)
	Delete(accountID uuid.UUID, qrCodeID uuid.UUID) error
	RecordScan(code string) error
}

type qrCodeService struct {
	qrCodeRepo     repositories.QRCodeRepository
	restaurantRepo repositories.RestaurantRepository
}

func NewQRCodeService(qrCodeRepo repositories.QRCodeRepository, restaurantRepo repositories.RestaurantRepository) QRCodeService {
	return &qrCodeService{
		qrCodeRepo:     qrCodeRepo,
		restaurantRepo: restaurantRepo,
	}
}

func (s *qrCodeService) Generate(accountID uuid.UUID, restaurantID uuid.UUID, qrType models.QRCodeType, label string) (*models.QRCode, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, errors.ErrForbidden
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
		IsActive:     true,
	}

	if err := s.qrCodeRepo.Create(qrCode); err != nil {
		return nil, err
	}

	return qrCode, nil
}

func (s *qrCodeService) GetByCode(code string) (*models.QRCode, error) {
	qrCode, err := s.qrCodeRepo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if !qrCode.IsValid() {
		return nil, errors.ErrNotFound
	}

	return qrCode, nil
}

func (s *qrCodeService) GetByRestaurantID(accountID uuid.UUID, restaurantID uuid.UUID) ([]models.QRCode, error) {
	// Verify restaurant ownership
	restaurant, err := s.restaurantRepo.FindByID(restaurantID)
	if err != nil {
		return nil, err
	}

	if restaurant.AccountID != accountID {
		return nil, errors.ErrForbidden
	}

	return s.qrCodeRepo.FindByRestaurantID(restaurantID)
}

func (s *qrCodeService) Delete(accountID uuid.UUID, qrCodeID uuid.UUID) error {
	// Get QR code
	qrCode, err := s.qrCodeRepo.FindByID(qrCodeID)
	if err != nil {
		return err
	}

	// Verify ownership
	restaurant, err := s.restaurantRepo.FindByID(qrCode.RestaurantID)
	if err != nil {
		return err
	}

	if restaurant.AccountID != accountID {
		return errors.ErrForbidden
	}

	return s.qrCodeRepo.Delete(qrCodeID)
}

func (s *qrCodeService) RecordScan(code string) error {
	qrCode, err := s.qrCodeRepo.FindByCode(code)
	if err != nil {
		return err
	}

	if !qrCode.IsValid() {
		return errors.ErrNotFound
	}

	return s.qrCodeRepo.IncrementScanCount(qrCode.ID)
}

func generateUniqueCode() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("LCQ-%s-%d", hex.EncodeToString(bytes)[:8], time.Now().Unix()), nil
}
