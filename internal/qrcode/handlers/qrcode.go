package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/lecritique/api/internal/qrcode/models"
	"github.com/lecritique/api/internal/qrcode/services"
	"github.com/lecritique/api/internal/shared/logger"
	"github.com/sirupsen/logrus"
)

type QRCodeHandler struct {
	qrCodeService services.QRCodeService
}

func NewQRCodeHandler(qrCodeService services.QRCodeService) *QRCodeHandler {
	return &QRCodeHandler{
		qrCodeService: qrCodeService,
	}
}

type GenerateQRCodeRequest struct {
	RestaurantID uuid.UUID          `json:"restaurant_id" validate:"required"`
	Type         models.QRCodeType  `json:"type" validate:"required,oneof=table location takeaway delivery general"`
	Label        string             `json:"label" validate:"required,min=1,max=100"`
}

type GenerateQRCodeResponse struct {
	Success bool           `json:"success"`
	Data    *models.QRCode `json:"data"`
}

func (h *QRCodeHandler) Generate(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req GenerateQRCodeRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request body")
	}

	if err := c.Validate(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

	qrCode, err := h.qrCodeService.Generate(ctx, accountID, req.RestaurantID, req.Type, req.Label)
	if err != nil {
		logger.Error("Failed to generate QR code", err, logrus.Fields{
			"account_id":    accountID,
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

func (h *QRCodeHandler) GetByRestaurant(c echo.Context) error {
	ctx := c.Request().Context()
	
	restaurantID, err := uuid.Parse(c.Param("restaurantId"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid restaurant ID")
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

	qrCodes, err := h.qrCodeService.GetByRestaurantID(ctx, accountID, restaurantID)
	if err != nil {
		logger.Error("Failed to get QR codes", err, logrus.Fields{
			"account_id":    accountID,
			"restaurant_id": restaurantID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to get QR codes")
	}

	return c.JSON(http.StatusOK, QRCodeListResponse{
		Success: true,
		Data:    qrCodes,
	})
}

func (h *QRCodeHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	
	qrCodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid QR code ID")
	}

	accountID, ok := c.Get("account_id").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authentication")
	}

	if err := h.qrCodeService.Delete(ctx, accountID, qrCodeID); err != nil {
		logger.Error("Failed to delete QR code", err, logrus.Fields{
			"account_id": accountID,
			"qr_code_id": qrCodeID,
		})
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to delete QR code")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "QR code deleted successfully",
	})
}