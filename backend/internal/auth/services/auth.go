package services

import (
	"context"
	"lecritique/internal/auth/models"
	"lecritique/internal/auth/repositories"
	"lecritique/internal/shared/config"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/services"
	"log"
	"net/http"
	"time"

	subscriptionServices "lecritique/internal/subscription/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type RegisterData struct {
	Email       string
	Password    string
	CompanyName string
	FirstName   string
	LastName    string
}

type AuthService interface {
	Register(ctx context.Context, data RegisterData) (*models.Account, error)
	Login(ctx context.Context, email, password string) (string, *models.Account, error)
	ValidateToken(tokenString string) (*Claims, error)
	RefreshToken(ctx context.Context, oldToken string) (string, error)
	SendEmailVerification(ctx context.Context, accountID uuid.UUID) error
	ResendVerificationEmail(ctx context.Context, email string) error
	VerifyEmail(ctx context.Context, token string) error
	SendPasswordReset(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, token, newPassword string) error
	RequestEmailChange(ctx context.Context, accountID uuid.UUID, newEmail string) (string, error)
	ConfirmEmailChange(ctx context.Context, token string) (string, error)
	RequestDeactivation(ctx context.Context, accountID uuid.UUID) error
	CancelDeactivation(ctx context.Context, accountID uuid.UUID) error
	ProcessPendingDeactivations(ctx context.Context) error
	UpdateProfile(ctx context.Context, accountID uuid.UUID, updates map[string]interface{}) (*models.Account, error)
	GetAccountByEmail(ctx context.Context, email string) (*models.Account, error)
}

type authService struct {
	accountRepo         repositories.AccountRepository
	tokenRepo           repositories.TokenRepository
	emailService        services.EmailService
	teamService         TeamMemberServiceV2
	subscriptionService subscriptionServices.SubscriptionService
	config              *config.Config
}

type SubscriptionFeatures struct {
	MaxRestaurants       int  `json:"max_restaurants"`
	MaxQRCodes           int  `json:"max_qr_codes"`
	MaxFeedbacksPerMonth int  `json:"max_feedbacks_per_month"`
	MaxTeamMembers       int  `json:"max_team_members"`
	HasBasicAnalytics    bool `json:"has_basic_analytics"`
	HasAdvancedAnalytics bool `json:"has_advanced_analytics"`
	HasFeedbackExplorer  bool `json:"has_feedback_explorer"`
	HasCustomBranding    bool `json:"has_custom_branding"`
	HasPrioritySupport   bool `json:"has_priority_support"`
}

type Claims struct {
	AccountID            uuid.UUID             `json:"account_id"`
	MemberID             uuid.UUID             `json:"member_id"`
	Name                 string                `json:"name"`
	Email                string                `json:"email"`
	Role                 models.MemberRole     `json:"role"`
	SubscriptionFeatures *SubscriptionFeatures `json:"subscription_features,omitempty"`
	jwt.RegisteredClaims
}

func NewAuthService(i *do.Injector) (AuthService, error) {
	return &authService{
		accountRepo:         do.MustInvoke[repositories.AccountRepository](i),
		tokenRepo:           do.MustInvoke[repositories.TokenRepository](i),
		emailService:        do.MustInvoke[services.EmailService](i),
		teamService:         do.MustInvoke[TeamMemberServiceV2](i),
		subscriptionService: do.MustInvoke[subscriptionServices.SubscriptionService](i),
		config:              do.MustInvoke[*config.Config](i),
	}, nil
}

func (s *authService) Register(ctx context.Context, data RegisterData) (*models.Account, error) {
	// Check if account already exists
	existing, _ := s.accountRepo.FindByEmail(ctx, data.Email)
	if existing != nil {
		return nil, errors.ErrConflict
	}

	// Create new account
	account := &models.Account{
		Email:       data.Email,
		CompanyName: data.CompanyName,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		IsActive:    true,
	}

	// Hash password
	if err := account.SetPassword(data.Password); err != nil {
		return nil, err
	}

	// Save account
	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	// If this is not a team member registration (no invitation), create owner team member record
	// This ensures the owner appears in the team members list
	// The migration script will handle existing accounts
	// For new accounts, the frontend should handle creating the owner team member
	if data.CompanyName != "" {
		log.Printf("New organization account created: %s. Owner team member should be created.", account.ID)
	}

	// Send verification email
	if err := s.SendEmailVerification(ctx, account.ID); err != nil {
		// Log error but don't fail registration
		log.Printf("Failed to send verification email for new account %s: %v", account.ID, err)
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

	// Check if email is verified
	if !account.EmailVerified {
		return "", nil, errors.NewWithDetails("EMAIL_NOT_VERIFIED", "Please verify your email address before logging in", 401, nil)
	}

	// If account has pending deactivation, cancel it
	if account.IsPendingDeactivation() {
		account.DeactivationRequestedAt = nil
		if err := s.accountRepo.Update(ctx, account); err != nil {
			// Log error but don't fail login
			log.Printf("Failed to cancel deactivation for account %s: %v", account.ID, err)
		}
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
		MemberID:  account.ID,
		Name:      account.CompanyName,
		Email:     account.Email,
		Role:      models.RoleOwner,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.config.JWT.Expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    s.config.App.Name,
			Subject:   account.ID.String(),
		},
	}

	teamMember, err := s.teamService.GetMemberByMemberID(context.Background(), account.ID)
	if err != nil {
		return "", err
	}

	if teamMember != nil {
		claims.AccountID = teamMember.Account.ID
		claims.Role = teamMember.Role
	}

	// Get subscription features
	subscription, err := s.subscriptionService.GetUserSubscription(context.Background(), claims.AccountID)
	if err != nil {
		return "", err
	}

	if subscription != nil && subscription.IsActive() {
		claims.SubscriptionFeatures = &SubscriptionFeatures{
			MaxRestaurants:       subscription.Plan.MaxRestaurants,
			MaxQRCodes:           subscription.Plan.MaxQRCodes,
			MaxFeedbacksPerMonth: subscription.Plan.MaxFeedbacksPerMonth,
			MaxTeamMembers:       subscription.Plan.MaxTeamMembers,
			HasBasicAnalytics:    subscription.Plan.HasBasicAnalytics,
			HasAdvancedAnalytics: subscription.Plan.HasAdvancedAnalytics,
			HasFeedbackExplorer:  subscription.Plan.HasFeedbackExplorer,
			HasCustomBranding:    subscription.Plan.HasCustomBranding,
			HasPrioritySupport:   subscription.Plan.HasPrioritySupport,
		}
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
	account, err := s.accountRepo.FindByID(ctx, claims.MemberID)
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

func (s *authService) ResendVerificationEmail(ctx context.Context, email string) error {
	// Find account by email
	account, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists - just log and return nil
		log.Printf("Resend verification requested for non-existent email: %s", email)
		return nil
	}

	// Check if already verified
	if account.EmailVerified {
		log.Printf("Resend verification requested for already verified account: %s", account.ID)
		return nil
	}

	// Delete any existing verification tokens for this account
	if err := s.tokenRepo.DeleteByAccountAndType(ctx, account.ID, models.TokenTypeEmailVerification); err != nil {
		log.Printf("Failed to delete existing verification tokens: %v", err)
		// Continue anyway
	}

	// Generate new verification token
	token, err := models.GenerateToken()
	if err != nil {
		return err
	}
	verificationToken := &models.VerificationToken{
		AccountID: account.ID,
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

func (s *authService) RequestEmailChange(ctx context.Context, accountID uuid.UUID, newEmail string) (string, error) {
	// Check if new email is already in use
	existingAccount, _ := s.accountRepo.FindByEmail(ctx, newEmail)
	if existingAccount != nil {
		return "", errors.NewWithDetails("EMAIL_EXISTS", "Email already in use", 400, nil)
	}

	// Get current account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return "", err
	}

	// Check if new email is same as current
	if account.Email == newEmail {
		return "", errors.NewWithDetails("SAME_EMAIL", "New email must be different from current email", 400, nil)
	}

	// In development mode without SMTP, change email immediately
	if s.config.IsDevMode() && !s.config.IsSMTPConfigured() {
		oldEmail := account.Email

		// Update email directly
		account.Email = newEmail
		account.EmailVerified = true
		now := time.Now()
		account.EmailVerifiedAt = &now

		// Save account
		if err := s.accountRepo.Update(ctx, account); err != nil {
			return "", err
		}

		// Generate new JWT token with updated email
		newToken, err := s.generateToken(account)
		if err != nil {
			return "", err
		}

		// Log the change
		log.Printf("DEV MODE: Email changed directly from %s to %s for account %s", oldEmail, newEmail, accountID)

		return newToken, nil
	}

	// Production mode or SMTP configured - use token verification
	// Generate token
	token, err := models.GenerateToken()
	if err != nil {
		return "", err
	}

	// Create email change token (24 hour expiry)
	changeToken := &models.VerificationToken{
		AccountID: accountID,
		Token:     token,
		Type:      models.TokenTypeEmailChange,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		NewEmail:  newEmail,
	}

	// Save token
	if err := s.tokenRepo.Create(ctx, changeToken); err != nil {
		return "", err
	}

	// Send verification email to new address
	if err := s.emailService.SendEmailChangeVerification(ctx, newEmail, token); err != nil {
		return "", err
	}

	return "", nil
}

func (s *authService) ConfirmEmailChange(ctx context.Context, token string) (string, error) {
	// Find token
	changeToken, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return "", errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired token", 400, nil)
	}

	// Check if token is valid
	if !changeToken.IsValid() {
		return "", errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired token", 400, nil)
	}

	// Check token type
	if changeToken.Type != models.TokenTypeEmailChange {
		return "", errors.NewWithDetails("INVALID_TOKEN", "Invalid token type", 400, nil)
	}

	// Check if new email is still available
	existingAccount, _ := s.accountRepo.FindByEmail(ctx, changeToken.NewEmail)
	if existingAccount != nil {
		return "", errors.NewWithDetails("EMAIL_EXISTS", "Email is no longer available", 400, nil)
	}

	// Get account
	account, err := s.accountRepo.FindByID(ctx, changeToken.AccountID)
	if err != nil {
		return "", err
	}

	// Update email
	account.Email = changeToken.NewEmail
	account.EmailVerified = true
	now := time.Now()
	account.EmailVerifiedAt = &now

	// Save account
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return "", err
	}

	// Mark token as used
	if err := s.tokenRepo.MarkAsUsed(ctx, changeToken.ID); err != nil {
		return "", err
	}

	// Generate new JWT token with updated email
	newToken, err := s.generateToken(account)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *authService) RequestDeactivation(ctx context.Context, accountID uuid.UUID) error {
	// Find account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return errors.ErrNotFound
	}

	// Check if already has pending deactivation
	if account.IsPendingDeactivation() {
		return errors.NewWithDetails("DEACTIVATION_EXISTS", "Account already has a pending deactivation request", 400, nil)
	}

	// Set deactivation request timestamp
	now := time.Now()
	account.DeactivationRequestedAt = &now

	// Update account
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	// Send notification email
	deactivationDate := account.GetDeactivationDate()
	if deactivationDate != nil {
		err = s.emailService.SendDeactivationRequest(ctx, account.Email, deactivationDate.Format("January 2, 2006"))
		if err != nil {
			// Log error but don't fail the deactivation request
			log.Printf("Failed to send deactivation email for account %s: %v", account.ID, err)
		}
	}

	return nil
}

func (s *authService) CancelDeactivation(ctx context.Context, accountID uuid.UUID) error {
	// Find account
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return errors.ErrNotFound
	}

	// Check if has pending deactivation
	if !account.IsPendingDeactivation() {
		return errors.NewWithDetails("NO_DEACTIVATION", "No pending deactivation request found", 400, nil)
	}

	// Clear deactivation request
	account.DeactivationRequestedAt = nil

	// Update account
	if err := s.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	// Send confirmation email
	err = s.emailService.SendDeactivationCancelled(ctx, account.Email)
	if err != nil {
		// Log error but don't fail the cancellation
		log.Printf("Failed to send deactivation cancellation email for account %s: %v", account.ID, err)
	}

	return nil
}

func (s *authService) ProcessPendingDeactivations(ctx context.Context) error {
	// Find all accounts that should be deactivated
	accounts, err := s.accountRepo.FindAccountsPendingDeactivation(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.ShouldBeDeactivated() {
			// Deactivate the account
			account.IsActive = false
			account.DeactivationRequestedAt = nil

			if err := s.accountRepo.Update(ctx, &account); err != nil {
				log.Printf("Failed to deactivate account %s: %v", account.ID, err)
				continue
			}

			// Send final notification
			err = s.emailService.SendAccountDeactivated(ctx, account.Email)
			if err != nil {
				log.Printf("Failed to send deactivation completion email for account %s: %v", account.ID, err)
			}

			log.Printf("Successfully deactivated account %s", account.ID)
		}
	}

	return nil
}

func (s *authService) UpdateProfile(ctx context.Context, accountID uuid.UUID, updates map[string]interface{}) (*models.Account, error) {
	// Validate that we have at least one update
	if len(updates) == 0 {
		return nil, errors.New("BAD_REQUEST", "No updates provided", http.StatusBadRequest)
	}

	// Get the account first
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	// Update the account fields manually
	if companyName, ok := updates["company_name"].(string); ok {
		account.CompanyName = companyName
	}
	if phone, ok := updates["phone"].(string); ok {
		account.Phone = phone
	}

	// Save the updated account
	err = s.accountRepo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *authService) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	return s.accountRepo.FindByEmail(ctx, email)
}
