package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/services"
	"github.com/lecritique/api/pkg/errors"
	"github.com/lecritique/api/pkg/response"
	"github.com/lecritique/api/pkg/validator"
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

func (h *AuthHandler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	account, err := h.authService.Register(req.Email, req.Password, req.CompanyName)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]interface{}{
		"account": account,
		"message": "Registration successful. Please check your email to verify your account.",
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	token, account, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, AuthResponse{
		Token:   token,
		Account: account,
	})
}

func (h *AuthHandler) RefreshToken(c echo.Context) error {
	// Get token from header
	tokenString := c.Request().Header.Get("Authorization")
	if tokenString == "" {
		return response.Error(c, errors.ErrUnauthorized)
	}

	// Remove "Bearer " prefix
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	newToken, err := h.authService.RefreshToken(tokenString)
	if err != nil {
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"token": newToken,
	})
}
