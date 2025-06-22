package errors

import (
	"fmt"
	"net/http"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Details any    `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("code: %s, message: %s", e.Code, e.Message)
}

func New(code string, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func NewWithDetails(code string, message string, status int, details any) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
	}
}

// Common errors
var (
	ErrBadRequest          = New("BAD_REQUEST", "Bad request", http.StatusBadRequest)
	ErrUnauthorized        = New("UNAUTHORIZED", "Unauthorized", http.StatusUnauthorized)
	ErrForbidden           = New("FORBIDDEN", "Forbidden", http.StatusForbidden)
	ErrNotFound            = New("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrConflict            = New("CONFLICT", "Resource already exists", http.StatusConflict)
	ErrInternalServer      = New("INTERNAL_SERVER_ERROR", "Internal server error", http.StatusInternalServerError)
	ErrValidation          = New("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest)
	ErrInvalidCredentials  = New("INVALID_CREDENTIALS", "Invalid credentials", http.StatusUnauthorized)
	ErrTokenExpired        = New("TOKEN_EXPIRED", "Token has expired", http.StatusUnauthorized)
	ErrTokenInvalid        = New("TOKEN_INVALID", "Invalid token", http.StatusUnauthorized)
	ErrSubscriptionLimit   = New("SUBSCRIPTION_LIMIT", "Subscription limit reached", http.StatusPaymentRequired)
)

func IsAppError(err error) (*AppError, bool) {
	appErr, ok := err.(*AppError)
	return appErr, ok
}
