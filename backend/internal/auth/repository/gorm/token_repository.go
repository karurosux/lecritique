package gormrepo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) authinterface.TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(ctx context.Context, token *models.VerificationToken) error {
	return r.db.WithContext(ctx).Create(token).Error
}

func (r *TokenRepository) FindByToken(ctx context.Context, token string) (*models.VerificationToken, error) {
	var verificationToken models.VerificationToken
	err := r.db.WithContext(ctx).Where("token = ?", token).First(&verificationToken).Error
	if err != nil {
		return nil, err
	}
	return &verificationToken, nil
}

func (r *TokenRepository) FindByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) ([]*models.VerificationToken, error) {
	var tokens []*models.VerificationToken
	err := r.db.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, tokenType).Find(&tokens).Error
	return tokens, err
}

func (r *TokenRepository) Update(ctx context.Context, token *models.VerificationToken) error {
	return r.db.WithContext(ctx).Save(token).Error
}

func (r *TokenRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&models.VerificationToken{}, id).Error
}

func (r *TokenRepository) DeleteExpired(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("expires_at < ?", time.Now()).Delete(&models.VerificationToken{}).Error
}

func (r *TokenRepository) MarkAsUsed(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.VerificationToken{}).Where("id = ?", id).Update("used_at", &now).Error
}

func (r *TokenRepository) DeleteByAccountAndType(ctx context.Context, accountID uuid.UUID, tokenType models.TokenType) error {
	return r.db.WithContext(ctx).Where("account_id = ? AND type = ?", accountID, tokenType).Delete(&models.VerificationToken{}).Error
}