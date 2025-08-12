package qrcodecontroller

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	qrcodeinterface "kyooar/internal/qrcode/interface"
	qrcodemodel "kyooar/internal/qrcode/model"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	"kyooar/internal/shared/validator"
	"github.com/sirupsen/logrus"
)

type QRCodeController struct {
	qrCodeService qrcodeinterface.QRCodeService
	validator     *validator.Validator
}

func NewQRCodeController(qrCodeService qrcodeinterface.QRCodeService) *QRCodeController {
	return &QRCodeController{
		qrCodeService: qrCodeService,
		validator:     validator.New(),
	}
}

type GenerateQRCodeRequest struct {
	OrganizationID uuid.UUID          `json:"organization_id" validate:"required"`
	Type         qrcodemodel.QRCodeType  `json:"type" validate:"required,oneof=table location takeaway delivery general"`
	Label        string             `json:"label" validate:"required,min=1,max=100"`
	Location     *string            `json:"location" validate:"omitempty,max=200"`
}

// @Summary Generate QR code
// @Description Generate a new QR code for a organization
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param qr_code body GenerateQRCodeRequest true "QR code information"
// @Success 201 {object} response.Response{data=qrcodemodel.QRCode}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/qr-codes [post]
func (h *QRCodeController) Generate(c echo.Context) error {
	ctx := c.Request().Context()
	
	var req GenerateQRCodeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	resourceAccountID := middleware.GetResourceAccountID(c)

	qrCode, err := h.qrCodeService.Generate(ctx, resourceAccountID, req.OrganizationID, req.Type, req.Label, req.Location)
	if err != nil {
		logger.Error("Failed to generate QR code", err, logrus.Fields{
			"account_id":    resourceAccountID,
			"organization_id": req.OrganizationID,
			"type":          req.Type,
			"label":         req.Label,
		})
		return response.Error(c, err)
	}

	return response.Success(c, qrCode)
}

// @Summary Get QR codes by organization
// @Description Get all QR codes for a specific organization
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param organizationId path string true "Organization ID"
// @Success 200 {object} response.Response{data=[]qrcodemodel.QRCode}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/organizations/{organizationId}/qr-codes [get]
func (h *QRCodeController) GetByOrganization(c echo.Context) error {
	ctx := c.Request().Context()
	
	organizationID, err := uuid.Parse(c.Param("organizationId"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	resourceAccountID := middleware.GetResourceAccountID(c)

	qrCodes, err := h.qrCodeService.GetByOrganizationID(ctx, resourceAccountID, organizationID)
	if err != nil {
		logger.Error("Failed to get QR codes", err, logrus.Fields{
			"account_id":    resourceAccountID,
			"organization_id": organizationID,
		})
		return response.Error(c, err)
	}

	return response.Success(c, qrCodes)
}

// @Summary Delete QR code
// @Description Delete a QR code from the system
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "QR Code ID"
// @Success 200 {object} response.Response{data=map[string]string}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/qr-codes/{id} [delete]
func (h *QRCodeController) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	
	qrCodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	resourceAccountID := middleware.GetResourceAccountID(c)

	if err := h.qrCodeService.Delete(ctx, resourceAccountID, qrCodeID); err != nil {
		logger.Error("Failed to delete QR code", err, logrus.Fields{
			"account_id": resourceAccountID,
			"qr_code_id": qrCodeID,
		})
		return response.Error(c, err)
	}

	return response.Success(c, map[string]string{
		"message": "QR code deleted successfully",
	})
}

type UpdateQRCodeRequest struct {
	IsActive *bool   `json:"is_active"`
	Label    *string `json:"label" validate:"omitempty,min=1,max=100"`
	Location *string `json:"location" validate:"omitempty,max=200"`
}

// @Summary Update QR code
// @Description Update QR code details like active status, label, or location
// @Tags qr-codes
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "QR Code ID"
// @Param qr_code body UpdateQRCodeRequest true "QR code update information"
// @Success 200 {object} response.Response{data=qrcodemodel.QRCode}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/qr-codes/{id} [patch]
func (h *QRCodeController) Update(c echo.Context) error {
	ctx := c.Request().Context()
	
	qrCodeID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	resourceAccountID := middleware.GetResourceAccountID(c)

	var req UpdateQRCodeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.ErrBadRequest)
	}

	if err := h.validator.Validate(&req); err != nil {
		return response.Error(c, errors.NewWithDetails("VALIDATION_ERROR", "Validation failed", http.StatusBadRequest, h.validator.FormatErrors(err)))
	}

	serviceReq := &qrcodeinterface.UpdateQRCodeRequest{
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
		return response.Error(c, err)
	}

	return response.Success(c, updatedQRCode)
}