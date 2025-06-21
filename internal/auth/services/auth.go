package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/auth/repositories"
	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/errors"
)

type AuthService interface {
	Register(email, password, companyName string) (*models.Account, error)
	Login(email, password string) (string, *models.Account, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(oldToken string) (string, error)
}

type authService struct {
	accountRepo repositories.AccountRepository
	config      *config.Config
}

type Claims struct {
	AccountID uuid.UUID `json:"account_id"`
	Email     string    `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(accountRepo repositories.AccountRepository, config *config.Config) AuthService {
	return &authService{
		accountRepo: accountRepo,
		config:      config,
	}
}

func (s *authService) Register(email, password, companyName string) (*models.Account, error) {
	// Check if account already exists
	existing, _ := s.accountRepo.FindByEmail(email)
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
	if err := s.accountRepo.Create(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *authService) Login(email, password string) (string, *models.Account, error) {
	// Find account
	account, err := s.accountRepo.FindByEmail(email)
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

func (s *authService) RefreshToken(oldToken string) (string, error) {
	claims, err := s.ValidateToken(oldToken)
	if err != nil {
		return "", err
	}

	// Get account
	account, err := s.accountRepo.FindByID(claims.AccountID)
	if err != nil {
		return "", err
	}

	// Generate new token
	return s.generateToken(account)
}
