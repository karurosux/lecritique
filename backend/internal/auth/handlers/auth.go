package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/models"
	"github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/shared/config"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
	subscriptionModels "github.com/lecritique/api/internal/subscription/models"
	"gorm.io/gorm"
)

type AuthHandler struct {
	authService       services.AuthService
	teamMemberService services.TeamMemberServiceV2
	validator         *validator.Validator
	config            *config.Config
	db                *gorm.DB
}

func NewAuthHandler(authService services.AuthService, teamMemberService services.TeamMemberServiceV2, config *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		authService:       authService,
		teamMemberService: teamMemberService,
		validator:         validator.New(),
		config:            config,
		db:                db,
	}
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	CompanyName     string `json:"company_name" validate:"required_without=InvitationToken"`
	FirstName       string `json:"first_name,omitempty"`
	LastName        string `json:"last_name,omitempty"`
	InvitationToken string `json:"invitation_token,omitempty"` // Optional invitation token
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token   string      `json:"token"`
	Account interface{} `json:"account"`
}

type EnhancedAuthResponse struct {
	Token        string      `json:"token"`
	Account      interface{} `json:"account"`
	Subscription interface{} `json:"subscription"`
}

// Register godoc
// @Summary Register a new account
// @Description Create a new restaurant owner account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=interface{}}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	registerData := services.RegisterData{
		Email:       req.Email,
		Password:    req.Password,
		CompanyName: req.CompanyName,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
	}
	
	account, err := h.authService.Register(ctx, registerData)
	if err != nil {
		return response.Error(c, err)
	}

	// Check if there's an invitation token
	if req.InvitationToken != "" {
		// Accept the invitation
		if err := h.teamMemberService.AcceptInvitationDuringRegistration(ctx, req.InvitationToken, account.ID); err != nil {
			// Log the error but don't fail registration
			log.Printf("Failed to accept team invitation: %v", err)
		}
	}

	// Also check for pending invitations that were already accepted via email
	invitations, err := h.teamMemberService.CheckPendingInvitations(ctx, account.Email)
	if err == nil && len(invitations) > 0 {
		for _, inv := range invitations {
			// Only auto-accept if the email link was clicked (EmailAcceptedAt is set)
			if inv.EmailAcceptedAt != nil && inv.AcceptedAt == nil {
				if err := h.teamMemberService.AcceptInvitationDuringRegistration(ctx, inv.Token, account.ID); err != nil {
					log.Printf("Failed to auto-accept team invitation: %v", err)
				}
			}
		}
	}

	return response.Success(c, map[string]interface{}{
		"account": account,
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

// Login godoc
// @Summary Login to account
// @Description Authenticate and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=EnhancedAuthResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	token, account, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return response.Error(c, err)
	}

	// Check for pending invitations that were already accepted via email
	invitations, err := h.teamMemberService.CheckPendingInvitations(ctx, account.Email)
	if err == nil && len(invitations) > 0 {
		for _, inv := range invitations {
			// Only auto-accept if the email link was clicked (EmailAcceptedAt is set)
			if inv.EmailAcceptedAt != nil && inv.AcceptedAt == nil {
				if err := h.teamMemberService.AcceptInvitationDuringRegistration(ctx, inv.Token, account.ID); err != nil {
					log.Printf("Failed to auto-accept team invitation: %v", err)
				}
			}
		}
	}

	// Fetch subscription data for the account
	var subscription *subscriptionModels.Subscription
	err = h.db.WithContext(ctx).
		Preload("Plan").
		Where("account_id = ?", account.ID).
		First(&subscription).Error
	
	if err != nil && err != gorm.ErrRecordNotFound {
		// Log error but don't fail login
		// This ensures backward compatibility
		subscription = nil
	}
	
	// Check if this user is a team member of another account
	var teamMemberships []models.TeamMember
	h.db.WithContext(ctx).
		Preload("Account").
		Where("member_id = ? AND accepted_at IS NOT NULL", account.ID).
		Find(&teamMemberships)
	
	log.Printf("Found %d team memberships for account %s", len(teamMemberships), account.ID)
	
	// Filter out memberships where the user is a member of their own account
	var externalMemberships []models.TeamMember
	for _, tm := range teamMemberships {
		if tm.AccountID != account.ID {
			externalMemberships = append(externalMemberships, tm)
		}
	}
	
	log.Printf("Found %d external team memberships for account %s", len(externalMemberships), account.ID)
	
	// If user is a team member of another organization, return their account but with the organization's subscription
	if len(externalMemberships) > 0 {
		// For now, use the first organization they're a member of
		// In the future, we might want to let users choose which organization to access
		orgAccountID := externalMemberships[0].AccountID
		
		log.Printf("Team member %s is member of org %s", account.Email, orgAccountID)
		
		// Get the organization's subscription
		var orgSubscription subscriptionModels.Subscription
		err := h.db.WithContext(ctx).
			Preload("Plan").
			Where("account_id = ?", orgAccountID).
			First(&orgSubscription).Error
		
		if err != nil {
			log.Printf("Error fetching org subscription: %v", err)
			// Return without subscription if not found
			return response.Success(c, EnhancedAuthResponse{
				Token:   token,
				Account: account,
			})
		}
		
		log.Printf("Team member %s logged in, providing organization subscription from account %s: %+v", account.Email, orgAccountID, orgSubscription)
		
		return response.Success(c, EnhancedAuthResponse{
			Token:        token,
			Account:      account, // Keep the user's own account
			Subscription: &orgSubscription, // But use the organization's subscription
		})
	}

	return response.Success(c, EnhancedAuthResponse{
		Token:        token,
		Account:      account,
		Subscription: subscription,
	})
}

type SendEmailVerificationRequest struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
}

type ResendVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type PasswordResetRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

// SendEmailVerification godoc
// @Summary Send email verification
// @Description Send verification email to the authenticated account
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/send-verification [post]
func (h *AuthHandler) SendEmailVerification(c echo.Context) error {
	ctx := c.Request().Context()
	accountID := c.Get("account_id").(uuid.UUID)

	if err := h.authService.SendEmailVerification(ctx, accountID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Verification email sent successfully",
	})
}

// ResendVerificationEmail godoc
// @Summary Resend email verification
// @Description Resend verification email to the specified email address (public endpoint)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResendVerificationRequest true "Email address"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/resend-verification [post]
func (h *AuthHandler) ResendVerificationEmail(c echo.Context) error {
	ctx := c.Request().Context()

	var req ResendVerificationRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	// Note: We don't reveal whether the email exists for security reasons
	if err := h.authService.ResendVerificationEmail(ctx, req.Email); err != nil {
		// Log the error but return success to prevent email enumeration
		// The service will handle logging internally
	}

	return response.Success(c, map[string]string{
		"message": "If an account with that email exists, a verification email has been sent",
	})
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Verify email address using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/verify-email [get]
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	ctx := c.Request().Context()
	token := c.QueryParam("token")

	if token == "" {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.authService.VerifyEmail(ctx, token); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Email verified successfully",
	})
}

// SendPasswordReset godoc
// @Summary Send password reset email
// @Description Send password reset email to the specified email address
// @Tags auth
// @Accept json
// @Produce json
// @Param request body PasswordResetRequest true "Email address"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) SendPasswordReset(c echo.Context) error {
	ctx := c.Request().Context()

	var req PasswordResetRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	if err := h.authService.SendPasswordReset(ctx, req.Email); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "If an account with that email exists, a password reset email has been sent",
	})
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	ctx := c.Request().Context()

	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	if err := h.authService.ResetPassword(ctx, req.Token, req.NewPassword); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Password reset successfully",
	})
}

// RefreshToken refreshes an existing JWT token
// @Summary Refresh JWT token
// @Description Refresh an existing JWT token to get a new one
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get token from header
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return response.Error(c, errors.ErrUnauthorized)
	}

	// Remove "Bearer " prefix
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	newToken, err := h.authService.RefreshToken(ctx, tokenString)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"token": newToken,
	})
}

type ChangeEmailRequest struct {
	NewEmail string `json:"new_email" validate:"required,email"`
}

type ConfirmEmailChangeRequest struct {
	Token string `json:"token" validate:"required"`
}

// ChangeEmail godoc
// @Summary Request email change
// @Description Request to change the account email address
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body ChangeEmailRequest true "New email address"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/change-email [post]
func (h *AuthHandler) ChangeEmail(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get account ID from context (set by auth middleware)
	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return response.Error(c, errors.ErrUnauthorized)
	}

	var req ChangeEmailRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	newToken, err := h.authService.RequestEmailChange(ctx, accountID, req.NewEmail)
	if err != nil {
		return response.Error(c, err)
	}

	// Check if we're in dev mode without SMTP
	message := "Email change request sent. Please check your new email for verification."
	responseData := map[string]string{
		"message": message,
	}
	
	if h.config.IsDevMode() && !h.config.IsSMTPConfigured() {
		message = "Email changed successfully (dev mode - no SMTP configured)."
		responseData["message"] = message
		// Include new token in dev mode when email is changed immediately
		if newToken != "" {
			responseData["token"] = newToken
		}
	}

	return response.Success(c, responseData)
}

// ConfirmEmailChange godoc
// @Summary Confirm email change
// @Description Confirm email change using the token sent to the new email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body ConfirmEmailChangeRequest true "Email change token"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/confirm-email-change [post]
func (h *AuthHandler) ConfirmEmailChange(c echo.Context) error {
	ctx := c.Request().Context()

	var req ConfirmEmailChangeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	newToken, err := h.authService.ConfirmEmailChange(ctx, req.Token)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Email changed successfully.",
		"token": newToken,
	})
}

// RequestDeactivation godoc
// @Summary Request account deactivation
// @Description Request to deactivate the account with a 15-day grace period
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/deactivate [post]
func (h *AuthHandler) RequestDeactivation(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get account ID from context (set by auth middleware)
	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return response.Error(c, errors.ErrUnauthorized)
	}

	if err := h.authService.RequestDeactivation(ctx, accountID); err != nil {
		return response.Error(c, err)
	}

	// Calculate deactivation date
	deactivationDate := time.Now().Add(15 * 24 * time.Hour)

	return response.Success(c, map[string]interface{}{
		"message": "Account deactivation requested. Your account will be deactivated on " + deactivationDate.Format("January 2, 2006") + ".",
		"deactivation_date": deactivationDate,
	})
}

// CancelDeactivation godoc
// @Summary Cancel account deactivation
// @Description Cancel a pending account deactivation request
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/cancel-deactivation [post]
func (h *AuthHandler) CancelDeactivation(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get account ID from context (set by auth middleware)
	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return response.Error(c, errors.ErrUnauthorized)
	}

	if err := h.authService.CancelDeactivation(ctx, accountID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Account deactivation request has been cancelled.",
	})
}

type UpdateProfileRequest struct {
	CompanyName string `json:"company_name,omitempty" validate:"omitempty,min=1"`
	Phone       string `json:"phone,omitempty" validate:"omitempty"`
	FirstName   string `json:"first_name,omitempty" validate:"omitempty,min=1"`
	LastName    string `json:"last_name,omitempty" validate:"omitempty,min=1"`
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Update user profile information including company name and personal details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body UpdateProfileRequest true "Profile update data"
// @Success 200 {object} response.Response{data=interface{}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/profile [put]
// @Security BearerAuth
func (h *AuthHandler) UpdateProfile(c echo.Context) error {
	ctx := c.Request().Context()
	
	// Get account ID from context (set by auth middleware)
	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return response.Error(c, errors.ErrUnauthorized)
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Invalid request", http.StatusBadRequest))
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, err)
	}

	// Update account fields
	updates := make(map[string]interface{})
	if req.CompanyName != "" {
		updates["company_name"] = req.CompanyName
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	// Note: FirstName and LastName would be updated on the User model
	// This requires fetching the user associated with the account
	// For now, we'll just update Account fields

	updatedAccount, err := h.authService.UpdateProfile(ctx, accountID, updates)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, updatedAccount)
}
