package services

import (
	"context"
	"kyooar/internal/auth/models"
	"kyooar/internal/auth/repositories"
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/services"
	"log"
	"net/http"
	"time"

	subscriptionServices "kyooar/internal/subscription/services"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/samber/do"
)

type RegisterData struct {
	Email     string
	Password  string
	Name      string
	FirstName string
	LastName  string
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
	MaxOrganizations       int  `json:"max_organizations"`
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
	existing, _ := s.accountRepo.FindByEmail(ctx, data.Email)
	if existing != nil {
		return nil, errors.ErrConflict
	}

	account := &models.Account{
		Email:     data.Email,
		Name:      data.Name,
		FirstName: data.FirstName,
		LastName:    data.LastName,
		IsActive:    true,
	}

	if err := account.SetPassword(data.Password); err != nil {
		return nil, err
	}

	if err := s.accountRepo.Create(ctx, account); err != nil {
		return nil, err
	}

	if data.Name != "" {
		log.Printf("New organization account created: %s. Owner team member should be created.", account.ID)
	}

	if err := s.SendEmailVerification(ctx, account.ID); err != nil {
		log.Printf("Failed to send verification email for new account %s: %v", account.ID, err)
	}

	return account, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, *models.Account, error) {
	account, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.ErrInvalidCredentials
	}

	if !account.CheckPassword(password) {
		return "", nil, errors.ErrInvalidCredentials
	}

	if !account.IsActive {
		return "", nil, errors.ErrUnauthorized
	}

	if !account.EmailVerified {
		return "", nil, errors.NewWithDetails("EMAIL_NOT_VERIFIED", "Please verify your email address before logging in", 401, nil)
	}

	if account.IsPendingDeactivation() {
		account.DeactivationRequestedAt = nil
		if err := s.accountRepo.Update(ctx, account); err != nil {
			log.Printf("Failed to cancel deactivation for account %s: %v", account.ID, err)
		}
	}

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
		Name:      account.Name,
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
		log.Printf("Could not get team member for account %s: %v", account.ID, err)
	}

	if teamMember != nil {
		claims.AccountID = teamMember.Account.ID
		claims.Role = teamMember.Role
	}

	subscription, err := s.subscriptionService.GetUserSubscription(context.Background(), claims.AccountID)
	if err != nil {
		log.Printf("Could not get subscription for account %s: %v", claims.AccountID, err)
	}

	if subscription != nil && subscription.IsActive() {
		claims.SubscriptionFeatures = &SubscriptionFeatures{
			MaxOrganizations:       subscription.Plan.MaxOrganizations,
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

	account, err := s.accountRepo.FindByID(ctx, claims.MemberID)
	if err != nil {
		return "", err
	}

	return s.generateToken(account)
}

func (s *authService) SendEmailVerification(ctx context.Context, accountID uuid.UUID) error {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return err
	}

	if account.EmailVerified {
		return errors.NewWithDetails("EMAIL_ALREADY_VERIFIED", "Email already verified", 400, nil)
	}

	token, err := models.GenerateToken()
	if err != nil {
		return err
	}

	verificationToken := &models.VerificationToken{
		AccountID: accountID,
		Token:     token,
		Type:      models.TokenTypeEmailVerification,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	if err := s.tokenRepo.Create(ctx, verificationToken); err != nil {
		return err
	}

	return s.emailService.SendVerificationEmail(ctx, account.Email, token)
}

func (s *authService) ResendVerificationEmail(ctx context.Context, email string) error {
	account, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		log.Printf("Resend verification requested for non-existent email: %s", email)
		return nil
	}

	if account.EmailVerified {
		log.Printf("Resend verification requested for already verified account: %s", account.ID)
		return nil
	}

	if err := s.tokenRepo.DeleteByAccountAndType(ctx, account.ID, models.TokenTypeEmailVerification); err != nil {
		log.Printf("Failed to delete existing verification tokens: %v", err)
	}

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

	if err := s.tokenRepo.Create(ctx, verificationToken); err != nil {
		return err
	}

	return s.emailService.SendVerificationEmail(ctx, account.Email, token)
}

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	verificationToken, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired verification token", 400, nil)
	}

	if !verificationToken.IsValid() {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired verification token", 400, nil)
	}

	if verificationToken.Type != models.TokenTypeEmailVerification {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid token type", 400, nil)
	}

	err = s.accountRepo.UpdateEmailVerification(ctx, verificationToken.AccountID, true)
	if err != nil {
		return err
	}

	return s.tokenRepo.MarkAsUsed(ctx, verificationToken.ID)
}

func (s *authService) SendPasswordReset(ctx context.Context, email string) error {
	account, err := s.accountRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil
	}

	token, err := models.GenerateToken()
	if err != nil {
		return err
	}

	resetToken := &models.VerificationToken{
		AccountID: account.ID,
		Token:     token,
		Type:      models.TokenTypePasswordReset,
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := s.tokenRepo.Create(ctx, resetToken); err != nil {
		return err
	}

	return s.emailService.SendPasswordResetEmail(ctx, email, token)
}

func (s *authService) ResetPassword(ctx context.Context, token, newPassword string) error {
	resetToken, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired reset token", 400, nil)
	}

	if !resetToken.IsValid() {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired reset token", 400, nil)
	}

	if resetToken.Type != models.TokenTypePasswordReset {
		return errors.NewWithDetails("INVALID_TOKEN", "Invalid token type", 400, nil)
	}

	account, err := s.accountRepo.FindByID(ctx, resetToken.AccountID)
	if err != nil {
		return err
	}

	if err := account.SetPassword(newPassword); err != nil {
		return err
	}

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	return s.tokenRepo.MarkAsUsed(ctx, resetToken.ID)
}

func (s *authService) RequestEmailChange(ctx context.Context, accountID uuid.UUID, newEmail string) (string, error) {
	existingAccount, _ := s.accountRepo.FindByEmail(ctx, newEmail)
	if existingAccount != nil {
		return "", errors.NewWithDetails("EMAIL_EXISTS", "Email already in use", 400, nil)
	}

	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return "", err
	}

	if account.Email == newEmail {
		return "", errors.NewWithDetails("SAME_EMAIL", "New email must be different from current email", 400, nil)
	}

	if s.config.IsDevMode() && !s.config.IsSMTPConfigured() {
		oldEmail := account.Email

		account.Email = newEmail
		account.EmailVerified = true
		now := time.Now()
		account.EmailVerifiedAt = &now

			if err := s.accountRepo.Update(ctx, account); err != nil {
			return "", err
		}

		newToken, err := s.generateToken(account)
		if err != nil {
			return "", err
		}

		log.Printf("DEV MODE: Email changed directly from %s to %s for account %s", oldEmail, newEmail, accountID)

		return newToken, nil
	}

	token, err := models.GenerateToken()
	if err != nil {
		return "", err
	}

	changeToken := &models.VerificationToken{
		AccountID: accountID,
		Token:     token,
		Type:      models.TokenTypeEmailChange,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		NewEmail:  newEmail,
	}

	if err := s.tokenRepo.Create(ctx, changeToken); err != nil {
		return "", err
	}

	if err := s.emailService.SendEmailChangeVerification(ctx, newEmail, token); err != nil {
		return "", err
	}

	return "", nil
}

func (s *authService) ConfirmEmailChange(ctx context.Context, token string) (string, error) {
	changeToken, err := s.tokenRepo.FindByToken(ctx, token)
	if err != nil {
		return "", errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired token", 400, nil)
	}

	if !changeToken.IsValid() {
		return "", errors.NewWithDetails("INVALID_TOKEN", "Invalid or expired token", 400, nil)
	}

	if changeToken.Type != models.TokenTypeEmailChange {
		return "", errors.NewWithDetails("INVALID_TOKEN", "Invalid token type", 400, nil)
	}

	existingAccount, _ := s.accountRepo.FindByEmail(ctx, changeToken.NewEmail)
	if existingAccount != nil {
		return "", errors.NewWithDetails("EMAIL_EXISTS", "Email is no longer available", 400, nil)
	}

	account, err := s.accountRepo.FindByID(ctx, changeToken.AccountID)
	if err != nil {
		return "", err
	}

	account.Email = changeToken.NewEmail
	account.EmailVerified = true
	now := time.Now()
	account.EmailVerifiedAt = &now

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return "", err
	}

	if err := s.tokenRepo.MarkAsUsed(ctx, changeToken.ID); err != nil {
		return "", err
	}

	newToken, err := s.generateToken(account)
	if err != nil {
		return "", err
	}

	return newToken, nil
}

func (s *authService) RequestDeactivation(ctx context.Context, accountID uuid.UUID) error {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return errors.ErrNotFound
	}

	if account.IsPendingDeactivation() {
		return errors.NewWithDetails("DEACTIVATION_EXISTS", "Account already has a pending deactivation request", 400, nil)
	}

	now := time.Now()
	account.DeactivationRequestedAt = &now

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	deactivationDate := account.GetDeactivationDate()
	if deactivationDate != nil {
		err = s.emailService.SendDeactivationRequest(ctx, account.Email, deactivationDate.Format("January 2, 2006"))
		if err != nil {
			log.Printf("Failed to send deactivation email for account %s: %v", account.ID, err)
		}
	}

	return nil
}

func (s *authService) CancelDeactivation(ctx context.Context, accountID uuid.UUID) error {
	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return errors.ErrNotFound
	}

	if !account.IsPendingDeactivation() {
		return errors.NewWithDetails("NO_DEACTIVATION", "No pending deactivation request found", 400, nil)
	}

	account.DeactivationRequestedAt = nil

	if err := s.accountRepo.Update(ctx, account); err != nil {
		return err
	}

	err = s.emailService.SendDeactivationCancelled(ctx, account.Email)
	if err != nil {
		log.Printf("Failed to send deactivation cancellation email for account %s: %v", account.ID, err)
	}

	return nil
}

func (s *authService) ProcessPendingDeactivations(ctx context.Context) error {
	accounts, err := s.accountRepo.FindAccountsPendingDeactivation(ctx)
	if err != nil {
		return err
	}

	for _, account := range accounts {
		if account.ShouldBeDeactivated() {
			account.IsActive = false
			account.DeactivationRequestedAt = nil

			if err := s.accountRepo.Update(ctx, &account); err != nil {
				log.Printf("Failed to deactivate account %s: %v", account.ID, err)
				continue
			}

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
	if len(updates) == 0 {
		return nil, errors.New("BAD_REQUEST", "No updates provided", http.StatusBadRequest)
	}

	account, err := s.accountRepo.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	if name, ok := updates["name"].(string); ok {
		account.Name = name
	}
	if phone, ok := updates["phone"].(string); ok {
		account.Phone = phone
	}

	err = s.accountRepo.Update(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *authService) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	return s.accountRepo.FindByEmail(ctx, email)
}
