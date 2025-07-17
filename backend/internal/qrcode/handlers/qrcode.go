package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/qrcode/models"
	"github.com/lecritique/api/internal/qrcode/services"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/lecritique/api/internal/shared/middleware"
	"github.com/lecritique/api/internal/shared/validator"
	"github.com/sirupsen/logrus"
)

type QRCodeHandler struct {
	qrCodeService services.QRCodeService
	validator     *validator.Validator
}

func NewQRCodeHandler(qrCodeService services.QRCodeService) *QRCodeHandler {
	return &QRCodeHandler{
		qrCodeService: qrCodeService,
		validator:     validator.New(),
	}
}

type GenerateQRCodeRequest struct {
	RestaurantID uuid.UUID          `json:"restaurant_id" validate:"required"`
	Type         models.QRCodeType  `json:"type" validate:"required,oneof=table location takeaway delivery general"`
	Label        string             `json:"label" validate:"required,min=1,max=100"`
	Location     *string            `json:"location" validate:"omitempty,max=200"`
}

type GenerateQRCodeResponse struct {
	Success bool           `json:"success"`
	Data    *models.QRCode `json:"data"`
}

// Generate creates a new QR code
// @Summary Generate QR code
// @Description Generate a new QR code for a restaurant
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param qr_code body GenerateQRCodeRequest true "QR code information"
// @Success 201 {object} GenerateQRCodeResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/restaurants/{restaurantId}/qr-codes [post]
func (h *QRCodeHandler) Generate(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req GenerateQRCodeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	qrCode, err := h.qrCodeService.Generate(ctx, resourceAccountID, req.RestaurantID, req.Type, req.Label, req.Location)
	if err != nil {
		logger.Error("Failed to generate QR code", err, logrus.Fields{
			"account_id":    resourceAccountID,
			"restaurant_id": req.RestaurantID,
			"type":          req.Type,
			"label":         req.Label,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate QR code")
	}

	return c.JSON(http.StatusCreated, GenerateQRCodeResponse{
		Success: true,
		Data:    qrCode,
	})
}

type QRCodeListResponse struct {
	Success bool             `json:"success"`
	Data    []models.QRCode  `json:"data"`
}

// GetByRestaurant gets all QR codes for a restaurant
// @Summary Get QR codes by restaurant
// @Description Get all QR codes for a specific restaurant
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param restaurantId path string true "Restaurant ID"
// @Success 200 {object} QRCodeListResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/restaurants/{restaurantId}/qr-codes [get]
func (h *QRCodeHandler) GetByRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	qrCodes, err := h.qrCodeService.GetByRestaurantID(ctx, resourceAccountID, restaurantID)
	if err != nil {
		logger.Error("Failed to get QR codes", err, logrus.Fields{
			"account_id":    resourceAccountID,
			"restaurant_id": restaurantID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get QR codes")
	}

	return c.JSON(http.StatusOK, QRCodeListResponse{
		Success: true,
		Data:    qrCodes,
	})
}

// Delete removes a QR code
// @Summary Delete QR code
// @Description Delete a QR code from the system
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "QR Code ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/qr-codes/{id} [delete]
func (h *QRCodeHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	
	qrCodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid QR code ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	if err := h.qrCodeService.Delete(ctx, resourceAccountID, qrCodeID); err != nil {
		logger.Error("Failed to delete QR code", err, logrus.Fields{
			"account_id": resourceAccountID,
			"qr_code_id": qrCodeID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete QR code")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "QR code deleted successfully",
	})
}

type UpdateQRCodeRequest struct {
	IsActive *bool   `json:"is_active"`
	Label    *string `json:"label" validate:"omitempty,min=1,max=100"`
	Location *string `json:"location" validate:"omitempty,max=200"`
}

type UpdateQRCodeResponse struct {
	Success bool           `json:"success"`
	Data    *models.QRCode `json:"data"`
}

// Update updates a QR code
// @Summary Update QR code
// @Description Update QR code details like active status, label, or location
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "QR Code ID"
// @Param qr_code body UpdateQRCodeRequest true "QR code update information"
// @Success 200 {object} UpdateQRCodeResponse
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/qr-codes/{id} [patch]
func (h *QRCodeHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	
	qrCodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid QR code ID")
	}

	// Use resource account ID for team-aware access
	resourceAccountID := middleware.GetResourceAccountID(c)

	var req UpdateQRCodeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := h.validator.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Convert handler request to service request format
	serviceReq := &services.UpdateQRCodeRequest{
		IsActive: req.IsActive,
		Label:    req.Label,
		Location: req.Location,
	}

	updatedQRCode, err := h.qrCodeService.Update(ctx, resourceAccountID, qrCodeID, serviceReq)
	if err != nil {
		logger.Error("Failed to update QR code", err, logrus.Fields{
			"account_id": resourceAccountID,
			"qr_code_id": qrCodeID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update QR code")
	}

	return c.JSON(http.StatusOK, UpdateQRCodeResponse{
		Success: true,
		Data:    updatedQRCode,
	})
}