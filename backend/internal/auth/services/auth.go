package services

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/auth/repositories"
	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/services"
)

type AuthService interface {
	Register(ctx context.Context, email, password, companyName string) (*models.Account, error)
	Login(ctx context.Context, email, password string) (string, *models.Account, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(ctx context.Context, oldToken string) (string, error)
	SendEmailVerification(ctx context.Context, accountID uuid.UUID) error
	VerifyEmail(ctx context.Context, token string) error
	SendPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
}

type authService struct {
	accountRepo repositories.AccountRepository
	tokenRepo   repositories.TokenRepository
	emailService services.EmailService
	config      *config.Config
}

type Claims struct {
	AccountID uuid.UUID `json:"account_id"`
	Email     string    `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(accountRepo repositories.AccountRepository, tokenRepo repositories.TokenRepository, emailService services.EmailService, config *config.Config) AuthService {
	return &authService{
		accountRepo:  accountRepo,
		tokenRepo:    tokenRepo,
		emailService: emailService,
		config:       config,
	}
}

func (s *authService) Register(ctx context.Context, email, password, companyName string) (*models.Account, error) {
	// Check if account already exists
	existing, _ := s.accountRepo.FindByEmail(ctx, email)
	if existing != nil {
		return nil, errors.ErrConflict
	}

	// Create new account
	account := &models.Account{
		Email:       email,
		CompanyName: companyName,
		IsActive:    true,
	}

	// Hash password
	if err := account.SetPassword(password); err != nil {
		return nil, err
	}

	// Save account
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *models.Account, error) {
	// Find account
	account, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.ErrInvalidCredentials
	}

	// Check password
	if !account.CheckPassword(password) {
		return "", nil, errors.ErrInvalidCredentials
	}

	// Check if account is active
	if !account.IsActive {
		return "", nil, errors.ErrUnauthorized
	}

	// Generate token
	token, err := s.generateToken(account)
	if err != nil {
		return "", nil, err
	}

	return token, account, nil
}

func (s *authService) generateToken(account *models.Account) (string, error) {
	claims := &Claims{
		AccountID: account.ID,
		Email:     account.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.App.Name,
			Subject:   account.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *authService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, errors.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrTokenInvalid
}

func (s *authService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
	claims, err := s.ValidateToken(oldToken)
	if err != nil {
		return "", err
	}

	// Get account
	account, err := s.accountRepo.FindByID(ctx, claims.AccountID)
	if err != nil {
		return "", err
	}

	// Generate new token
	return s.generateToken(account)
}

func (s *authService) SendEmailVerification(ctx context.Context, accountID uuid.UUID) error {
	// Get account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	// Check if already verified
	if account.EmailVerified {
		return errors.NewWithDetails("EMAIL_ALREADY_VERIFIED", "Email already verified", 400, nil)
	}

	// Generate token
	token, err := models.GenerateToken()
	if err != nil {
		return err
	}

	// Create verification token (24 hours expiry)
	verificationToken := &models.VerificationToken{
		AccountID: accountID,
		Token:     token,
		Type:      models.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	// Save token
	if err := s.tokenRepo.Create(ctx, verificationToken); err != nil {
		return err
	}

	// Send email
	return s.emailService.SendVerificationEmail(ctx, account.Email, token)
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	// Find token
	verificationToken, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired verification token", 400, nil)
	}

	// Check if token is valid
	if !verificationToken.IsValid() {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired verification token", 400, nil)
	}

	// Check token type
	if verificationToken.Type != models.TokenTypeEmailVerification {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid token type", 400, nil)
	}

	// Mark account as verified
	err = s.accountRepo.UpdateEmailVerification(ctx, verificationToken.AccountID, true)
	if err != nil {
		return err
	}

	// Mark token as used
	return s.tokenRepo.MarkAsUsed(ctx, verificationToken.ID)
}

func (s *authService) SendPasswordReset(ctx context.Context, email string) error {
	// Find account
	account, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not for security
		return nil
	}

	// Generate token
	token, err := models.GenerateToken()
	if err != nil {
		return err
	}

	// Create reset token (1 hour expiry)
	resetToken := &models.VerificationToken{
		AccountID: account.ID,
		Token:     token,
		Type:      models.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	// Save token
	if err := s.tokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	// Send email
	return s.emailService.SendPasswordResetEmail(ctx, email, token)
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	// Find token
	resetToken, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired reset token", 400, nil)
	}

	// Check if token is valid
	if !resetToken.IsValid() {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired reset token", 400, nil)
	}

	// Check token type
	if resetToken.Type != models.TokenTypePasswordReset {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid token type", 400, nil)
	}

	// Get account
	account, err := s.accountRepo.FindByID(ctx, resetToken.AccountID)
	if err != nil {
		return err
	}

	// Update password
	if err := account.SetPassword(newPassword); err != nil {
		return err
	}

	// Save account
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	// Mark token as used
	return s.tokenRepo.MarkAsUsed(ctx, resetToken.ID)
}
