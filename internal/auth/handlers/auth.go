package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/auth/services"
	"github.com/lecritique/api/internal/shared/errors"
	"github.com/lecritique/api/internal/shared/response"
	"github.com/lecritique/api/internal/shared/validator"
)

type AuthHandler struct {
	authService services.AuthService
	validator   *validator.Validator
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	CompanyName string `json:"company_name" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	Token   string      `json:"token"`
	Account interface{} `json:"account"`
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

	account, err := h.authService.Register(ctx, req.Email, req.Password, req.CompanyName)
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

	token, account, err := h.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, AuthResponse{
		Token:   token,
		Account: account,
	})
}

type SendEmailVerificationRequest struct {
	AccountID string `json:"account_id" validate:"required,uuid"`
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
