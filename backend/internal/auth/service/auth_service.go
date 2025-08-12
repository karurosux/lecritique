package authservice

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	authinterface "kyooar/internal/auth/interface"
	"kyooar/internal/auth/models"
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/services"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

type AuthService struct {
	accountRepo         authinterface.AccountRepository
	tokenRepo           authinterface.TokenRepository
	emailService        services.EmailService
	teamService         authinterface.TeamMemberService
	subscriptionService subscriptioninterface.SubscriptionService
	config              *config.Config
}

func NewAuthService(
	accountRepo authinterface.AccountRepository,
	tokenRepo authinterface.TokenRepository,
	emailService services.EmailService,
	teamService authinterface.TeamMemberService,
	subscriptionService subscriptioninterface.SubscriptionService,
	config *config.Config,
) authinterface.AuthService {
	return &AuthService{
		accountRepo:         accountRepo,
		tokenRepo:           tokenRepo,
		emailService:        emailService,
		teamService:         teamService,
		subscriptionService: subscriptionService,
		config:              config,
	}
}

func (s *AuthService) Register(ctx context.Context, data authinterface.RegisterData) (*models.Account, error) {
	existing, _ := s.accountRepo.FindByEmail(ctx, data.Email)
	if existing != nil {
		return nil, errors.ErrConflict
	}

	account := &models.Account{
		Email:     data.Email,
		Name:      data.Name,
		FirstName: data.FirstName,
		LastName:  data.LastName,
		IsActive:  true,
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

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *models.Account, error) {
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

func (s *AuthService) generateToken(account *models.Account) (string, error) {
	claims := &authinterface.Claims{
		AccountID: account.ID,
		MemberID:  account.ID,
		Name:      account.Name,
		Email:     account.Email,
		Role:      models.RoleOwner,
	}

	teamMember, err := s.teamService.GetMemberByMemberID(context.Background(), account.ID)
	if err != nil {
		log.Printf("Could not get team member for account %s: %v", account.ID, err)
	}

	if teamMember != nil {
		claims.AccountID = teamMember.Account.ID
		claims.Role = teamMember.Role
	}

	if s.subscriptionService != nil {
		subscription, err := s.subscriptionService.GetUserSubscription(context.Background(), claims.AccountID)
		if err != nil {
			log.Printf("Could not get subscription for account %s: %v", claims.AccountID, err)
		}

		if subscription != nil && subscription.IsActive() {
			claims.SubscriptionFeatures = &authinterface.SubscriptionFeatures{
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
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"account_id":             claims.AccountID.String(),
		"member_id":              claims.MemberID.String(),
		"name":                   claims.Name,
		"email":                  claims.Email,
		"role":                   string(claims.Role),
		"subscription_features":  claims.SubscriptionFeatures,
		"exp":                    time.Now().Add(s.config.JWT.Expiration).Unix(),
		"iat":                    time.Now().Unix(),
		"nbf":                    time.Now().Unix(),
		"iss":                    s.config.App.Name,
		"sub":                    account.ID.String(),
	})

	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthService) ValidateToken(tokenString string) (*authinterface.Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return nil, errors.ErrTokenInvalid
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		accountID, _ := uuid.Parse(claims["account_id"].(string))
		memberID, _ := uuid.Parse(claims["member_id"].(string))
		
		return &authinterface.Claims{
			AccountID: accountID,
			MemberID:  memberID,
			Name:      claims["name"].(string),
			Email:     claims["email"].(string),
			Role:      models.MemberRole(claims["role"].(string)),
		}, nil
	}

	return nil, errors.ErrTokenInvalid
}

func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (string, error) {
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

func (s *AuthService) SendEmailVerification(ctx context.Context, accountID uuid.UUID) error {
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

func (s *AuthService) ResendVerificationEmail(ctx context.Context, email string) error {
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

func (s *AuthService) VerifyEmail(ctx context.Context, token string) error {
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

func (s *AuthService) SendPasswordReset(ctx context.Context, email string) error {
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

func (s *AuthService) ResetPassword(ctx context.Context, token, newPassword string) error {
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

func (s *AuthService) RequestEmailChange(ctx context.Context, accountID uuid.UUID, newEmail string) (string, error) {
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

func (s *AuthService) ConfirmEmailChange(ctx context.Context, token string) (string, error) {
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

func (s *AuthService) RequestDeactivation(ctx context.Context, accountID uuid.UUID) error {
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

func (s *AuthService) CancelDeactivation(ctx context.Context, accountID uuid.UUID) error {
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

func (s *AuthService) ProcessPendingDeactivations(ctx context.Context) error {
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

func (s *AuthService) UpdateProfile(ctx context.Context, accountID uuid.UUID, updates map[string]interface{}) (*models.Account, error) {
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

func (s *AuthService) GetAccountByEmail(ctx context.Context, email string) (*models.Account, error) {
	return s.accountRepo.FindByEmail(ctx, email)
}