package handlers

import (
	"kyooar/internal/auth/services"
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/samber/do"
)

type AuthHandler struct {
	authService       services.AuthService
	teamMemberService services.TeamMemberServiceV2
	validator         *validator.Validator
	config            *config.Config
}

func NewAuthHandler(i *do.Injector) (*AuthHandler, error) {
	return &AuthHandler{
		authService:       do.MustInvoke[services.AuthService](i),
		teamMemberService: do.MustInvoke[services.TeamMemberServiceV2](i),
		validator:         validator.New(),
		config:            do.MustInvoke[*config.Config](i),
	}, nil
}

type RegisterRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	Name            string `json:"name" validate:"required_without=InvitationToken"`
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

type TokenResponse struct {
	Token string `json:"token"`
}

// Register godoc
// @Summary Register a new account
// @Description Create a new organization owner account
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
		Name:        req.Name,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
	}

	account, err := h.authService.Register(ctx, registerData)
	if err != nil {
		return response.Error(c, err)
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
// @Success 200 {object} response.Response{data=AuthResponse}
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

	token, _, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return response.Error(c, err)
	}

	// All necessary data (account info, role, subscription features) is now in the JWT token
	// Frontend will decode the token to get this information
	return response.Success(c, TokenResponse{
		Token: token,
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
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

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

	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return response.Error(c, errors.ErrUnauthorized)
	}

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

	accountID, ok := c.Get("member_id").(uuid.UUID)
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
		"token":   newToken,
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

	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.authService.RequestDeactivation(ctx, accountID); err != nil {
		return response.Error(c, err)
	}

	deactivationDate := time.Now().Add(15 * 24 * time.Hour)

	return response.Success(c, map[string]interface{}{
		"message":           "Account deactivation requested. Your account will be deactivated on " + deactivationDate.Format("January 2, 2006") + ".",
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

	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	if err := h.authService.CancelDeactivation(ctx, accountID); err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "Account deactivation request has been cancelled.",
	})
}

type UpdateProfileRequest struct {
	Name  string `json:"name,omitempty" validate:"omitempty,min=1"`
	Phone string `json:"phone,omitempty" validate:"omitempty"`
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

	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	var req UpdateProfileRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Invalid request", http.StatusBadRequest))
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	updatedAccount, err := h.authService.UpdateProfile(ctx, accountID, updates)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, updatedAccount)
}
