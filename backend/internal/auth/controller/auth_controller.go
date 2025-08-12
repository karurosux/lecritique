package authcontroller

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	authconstants "kyooar/internal/auth/constants"
	authinterface "kyooar/internal/auth/interface"
	authmodel "kyooar/internal/auth/model"
	"kyooar/internal/shared/config"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
)

type AuthController struct {
	authService       authinterface.AuthService
	teamMemberService authinterface.TeamMemberService
	validator         *validator.Validator
	config            *config.Config
}

func NewAuthController(
	authService authinterface.AuthService,
	teamMemberService authinterface.TeamMemberService,
	validator *validator.Validator,
	config *config.Config,
) *AuthController {
	return &AuthController{
		authService:       authService,
		teamMemberService: teamMemberService,
		validator:         validator,
		config:            config,
	}
}

func (c *AuthController) handleError(ctx echo.Context, err error) error {
	switch err.Error() {
	case authconstants.ErrAccountAlreadyExists:
		return response.Error(ctx, errors.ErrConflict)
	case authconstants.ErrAccountNotFound:
		return response.Error(ctx, errors.ErrNotFound)
	case authconstants.ErrInvalidCredentials:
		return response.Error(ctx, errors.ErrInvalidCredentials)
	case authconstants.ErrInvalidToken:
		return response.Error(ctx, errors.ErrTokenInvalid)
	case authconstants.ErrEmailNotVerified:
		return response.Error(ctx, errors.NewWithDetails("EMAIL_NOT_VERIFIED", "Please verify your email address before logging in", 401, nil))
	default:
		return response.Error(ctx, err)
	}
}

// @Summary Register a new account
// @Description Create a new organization owner account
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.RegisterRequest true "Registration details"
// @Success 201 {object} response.Response{data=interface{}}
// @Failure 400 {object} response.Response
// @Failure 409 {object} response.Response
// @Router /api/v1/auth/register [post]
func (c *AuthController) Register(ctx echo.Context) error {
	var req authmodel.RegisterRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	registerData := authinterface.RegisterData{
		Email:     req.Email,
		Password:  req.Password,
		Name:      req.Name,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	account, err := c.authService.Register(ctx.Request().Context(), registerData)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]interface{}{
		"account": account,
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

// @Summary Login to account
// @Description Authenticate and get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.LoginRequest true "Login credentials"
// @Success 200 {object} response.Response{data=authmodel.TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/login [post]
func (c *AuthController) Login(ctx echo.Context) error {
	var req authmodel.LoginRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	token, _, err := c.authService.Login(ctx.Request().Context(), req.Email, req.Password)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, authmodel.TokenResponse{
		Token: token,
	})
}

// @Summary Refresh JWT token
// @Description Refresh an existing JWT token to get a new one
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=authmodel.TokenResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/refresh [post]
func (c *AuthController) RefreshToken(ctx echo.Context) error {
	tokenString := ctx.Request().Header.Get("Authorization")
	if tokenString == "" {
		return response.Error(ctx, errors.ErrUnauthorized)
	}

	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	newToken, err := c.authService.RefreshToken(ctx.Request().Context(), tokenString)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, authmodel.TokenResponse{
		Token: newToken,
	})
}

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
func (c *AuthController) SendEmailVerification(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	if err := c.authService.SendEmailVerification(ctx.Request().Context(), accountID); err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Verification email sent successfully",
	})
}

// @Summary Resend email verification
// @Description Resend verification email to the specified email address (public endpoint)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.ResendVerificationRequest true "Email address"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/resend-verification [post]
func (c *AuthController) ResendVerificationEmail(ctx echo.Context) error {
	var req authmodel.ResendVerificationRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	if err := c.authService.ResendVerificationEmail(ctx.Request().Context(), req.Email); err != nil {
	}

	return response.Success(ctx, map[string]string{
		"message": "If an account with that email exists, a verification email has been sent",
	})
}

// @Summary Verify email address
// @Description Verify email address using verification token
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Verification token"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/verify-email [get]
func (c *AuthController) VerifyEmail(ctx echo.Context) error {
	token := ctx.QueryParam("token")
	if token == "" {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.authService.VerifyEmail(ctx.Request().Context(), token); err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Email verified successfully",
	})
}

// @Summary Send password reset email
// @Description Send password reset email to the specified email address
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.PasswordResetRequest true "Email address"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/forgot-password [post]
func (c *AuthController) SendPasswordReset(ctx echo.Context) error {
	var req authmodel.PasswordResetRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	if err := c.authService.SendPasswordReset(ctx.Request().Context(), req.Email); err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "If an account with that email exists, a password reset email has been sent",
	})
}

// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/reset-password [post]
func (c *AuthController) ResetPassword(ctx echo.Context) error {
	var req authmodel.ResetPasswordRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	if err := c.authService.ResetPassword(ctx.Request().Context(), req.Token, req.NewPassword); err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Password reset successfully",
	})
}

// @Summary Request email change
// @Description Request to change the account email address
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body authmodel.ChangeEmailRequest true "New email address"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/email-change [post]
func (c *AuthController) ChangeEmail(ctx echo.Context) error {
	accountID, ok := ctx.Get("member_id").(uuid.UUID)
	if !ok {
		return response.Error(ctx, errors.ErrUnauthorized)
	}

	var req authmodel.ChangeEmailRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	newToken, err := c.authService.RequestEmailChange(ctx.Request().Context(), accountID, req.NewEmail)
	if err != nil {
		return c.handleError(ctx, err)
	}

	message := "Email change request sent. Please check your new email for verification."
	responseData := map[string]string{
		"message": message,
	}

	if c.config.IsDevMode() && !c.config.IsSMTPConfigured() {
		message = "Email changed successfully (dev mode - no SMTP configured)."
		responseData["message"] = message
		if newToken != "" {
			responseData["token"] = newToken
		}
	}

	return response.Success(ctx, responseData)
}

// @Summary Confirm email change
// @Description Confirm email change using the token sent to the new email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.ConfirmEmailChangeRequest true "Email change token"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Router /api/v1/auth/email-change/confirm [post]
func (c *AuthController) ConfirmEmailChange(ctx echo.Context) error {
	var req authmodel.ConfirmEmailChangeRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.ErrBadRequest)
	}

	if err := c.validator.Validate(req); err != nil {
		return response.Error(ctx, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, c.validator.FormatErrors(err)))
	}

	newToken, err := c.authService.ConfirmEmailChange(ctx.Request().Context(), req.Token)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Email changed successfully.",
		"token":   newToken,
	})
}

// @Summary Request account deactivation
// @Description Request to deactivate the account with a 15-day grace period
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=authmodel.DeactivationResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/deactivate [post]
func (c *AuthController) RequestDeactivation(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	if err := c.authService.RequestDeactivation(ctx.Request().Context(), accountID); err != nil {
		return c.handleError(ctx, err)
	}

	deactivationDate := time.Now().Add(15 * 24 * time.Hour)

	return response.Success(ctx, authmodel.DeactivationResponse{
		Message:          "Account deactivation requested. Your account will be deactivated on " + deactivationDate.Format("January 2, 2006") + ".",
		DeactivationDate: &deactivationDate,
	})
}

// @Summary Cancel account deactivation
// @Description Cancel a pending account deactivation request
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/auth/deactivate/cancel [post]
func (c *AuthController) CancelDeactivation(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	if err := c.authService.CancelDeactivation(ctx.Request().Context(), accountID); err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, map[string]string{
		"message": "Account deactivation request has been cancelled.",
	})
}

// @Summary Update user profile
// @Description Update user profile information including company name and personal details
// @Tags auth
// @Accept json
// @Produce json
// @Param request body authmodel.UpdateProfileRequest true "Profile update data"
// @Success 200 {object} response.Response{data=interface{}}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/auth/profile [put]
// @Security BearerAuth
func (c *AuthController) UpdateProfile(ctx echo.Context) error {
	accountID, err := middleware.GetAccountID(ctx)
	if err != nil {
		return response.Error(ctx, err)
	}

	var req authmodel.UpdateProfileRequest
	if err := ctx.Bind(&req); err != nil {
		return response.Error(ctx, errors.New("BAD_REQUEST", "Invalid request", http.StatusBadRequest))
	}

	if err := c.validator.Validate(&req); err != nil {
		return response.Error(ctx, err)
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}

	updatedAccount, err := c.authService.UpdateProfile(ctx.Request().Context(), accountID, updates)
	if err != nil {
		return c.handleError(ctx, err)
	}

	return response.Success(ctx, updatedAccount)
}

func (c *AuthController) RegisterRoutes(v1 *echo.Group, authMiddleware echo.MiddlewareFunc) {
	auth := v1.Group("/auth")
	
	// Public endpoints
	auth.POST("/register", c.Register)
	auth.POST("/login", c.Login)
	auth.POST("/refresh", c.RefreshToken)
	auth.GET("/verify-email", c.VerifyEmail)
	auth.POST("/resend-verification", c.ResendVerificationEmail)
	auth.POST("/forgot-password", c.SendPasswordReset)
	auth.POST("/reset-password", c.ResetPassword)

	// Protected endpoints
	authProtected := v1.Group("/auth")
	authProtected.Use(authMiddleware)
	authProtected.PUT("/profile", c.UpdateProfile)
	authProtected.POST("/deactivate", c.RequestDeactivation)
	authProtected.POST("/deactivate/cancel", c.CancelDeactivation)
	authProtected.POST("/email-change", c.ChangeEmail)
	authProtected.POST("/email-change/confirm", c.ConfirmEmailChange)
	authProtected.POST("/send-verification", c.SendEmailVerification)
}