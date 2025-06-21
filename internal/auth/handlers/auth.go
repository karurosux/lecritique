package handlers

import (
	"net/http"

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
// @Router /auth/register [post]
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
// @Router /auth/login [post]
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
