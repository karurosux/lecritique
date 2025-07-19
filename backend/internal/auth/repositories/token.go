package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"
	"lecritique/internal/auth/models"
	"lecritique/internal/shared/errors"
	"github.com/samber/do"
	"gorm.io/gorm"
)

type TokenRepository interface {
	Create(ctx context.Context, token *models.VerificationToken) error
	FindByToken(ctx context.Context, token string) (*models.VerificationToken, error)
	FindByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) (*models.VerificationToken, error)
	MarkAsUsed(ctx context.Context, tokenID uuid.UUID) error
	DeleteExpired(ctx context.Context) error
	DeleteByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) error
}

type tokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(i *do.Injector) (TokenRepository, error) {
	db := do.MustInvoke[*gorm.DB](i)
	return &tokenRepository{db: db}, nil
}

func (r *tokenRepository) Create(ctx context.Context, token *models.VerificationToken) error {
	if err := r.db.WithContext(ctx).Create(token).Error; err != nil {
		return errors.ErrInternalServer
	}
	return nil
}

func (r *tokenRepository) FindByToken(ctx context.Context, token string) (*models.VerificationToken, error) {
	var verificationToken models.VerificationToken
	err := r.db.WithContext(ctx).Preload("Account").Where("token = ?", token).First(&verificationToken).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return &verificationToken, nil
}

func (r *tokenRepository) FindByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) (*models.VerificationToken, error) {
	var token models.VerificationToken
	err := r.db.WithContext(ctx).Where("account_id = ? AND type = ? AND used_at IS NULL AND expires_at > ?", 
		accountID, tokenType, time.Now()).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.ErrNotFound
		}
		return nil, errors.ErrInternalServer
	}
	return &token, nil
}

func (r *tokenRepository) MarkAsUsed(ctx context.Context, tokenID uuid.UUID) error {
	err := r.db.WithContext(ctx).Model(&models.VerificationToken{}).
		Where("id = ?", tokenID).
		Update("used_at", time.Now()).Error
	if err != nil {
		return errors.ErrInternalServer
	}
	return nil
}

func (r *tokenRepository) DeleteExpired(ctx context.Context) error {
	err := r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.VerificationToken{}).Error
	if err != nil {
		return errors.ErrInternalServer
	}
	return nil
}

func (r *tokenRepository) DeleteByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) error {
	err := r.db.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, tokenType).Delete(&models.VerificationToken{}).Error
	if err != nil {
		return errors.ErrInternalServer
	}
	return nil
}