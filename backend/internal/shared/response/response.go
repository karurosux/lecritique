package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"lecritique/internal/shared/errors"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorData  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorData struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Meta struct {
	Timestamp time.Time   `json:"timestamp"`
	Version   string      `json:"version"`
	RequestID string      `json:"request_id,omitempty"`
	Pagination *Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
	Pages int `json:"pages"`
}

// Helper function to send JSON response with pretty formatting in debug mode
func sendJSON(c echo.Context, code int, data interface{}) error {
	if c.Echo().Debug {
		c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
		encoder := json.NewEncoder(c.Response())
		encoder.SetIndent("", "  ")
		c.Response().WriteHeader(code)
		return encoder.Encode(data)
	}
	return c.JSON(code, data)
}

func Success(c echo.Context, data interface{}) error {
	response := Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp: time.Now(),
			Version:   "1.0",
			RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
		},
	}
	return sendJSON(c, http.StatusOK, response)
}

func SuccessWithPagination(c echo.Context, data interface{}, pagination *Pagination) error {
	response := Response{
		Success: true,
		Data:    data,
		Meta: &Meta{
			Timestamp:  time.Now(),
			Version:    "1.0",
			RequestID:  c.Response().Header().Get(echo.HeaderXRequestID),
			Pagination: pagination,
		},
	}
	return sendJSON(c, http.StatusOK, response)
}

func Error(c echo.Context, err error) error {
	appErr, ok := errors.IsAppError(err)
	if !ok {
		appErr = errors.ErrInternalServer
	}

	response := Response{
		Success: false,
		Error: &ErrorData{
			Code:    appErr.Code,
			Message: appErr.Message,
			Details: appErr.Details,
		},
		Meta: &Meta{
			Timestamp: time.Now(),
			Version:   "1.0",
			RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
		},
	}
	return sendJSON(c, appErr.Status, response)
}
