package handlers

import (
	"github.com/labstack/echo/v4"
	qrcodeServices "kyooar/internal/qrcode/services"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/response"
	"github.com/samber/do"
	"github.com/sirupsen/logrus"
)

type QRCodePublicHandler struct {
	qrCodeService qrcodeServices.QRCodeService
}

func NewQRCodePublicHandler(i *do.Injector) (*QRCodePublicHandler, error) {
	return &QRCodePublicHandler{
		qrCodeService: do.MustInvoke[qrcodeServices.QRCodeService](i),
	}, nil
}

// @Summary Validate QR code
// @Description Validate a QR code and return associated data
// @Tags public
// @Accept json
// @Produce json
// @Param code path string true "QR Code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/qr/{code} [get]
func (h *QRCodePublicHandler) ValidateQRCode(c echo.Context) error {
	ctx := c.Request().Context()
	code := c.Param("code")
	if code == "" {
		return response.Error(c, errors.BadRequest("QR code parameter is required"))
	}

	qrCode, err := h.qrCodeService.GetByCode(ctx, code)
	if err != nil {
		return response.Error(c, errors.NotFound("QR code"))
	}

	if err := h.qrCodeService.RecordScan(ctx, code); err != nil {
		logger.Error("Failed to record QR scan", err, logrus.Fields{
			"qr_code_id": qrCode.ID,
			"code":       code,
		})
	}

	return response.Success(c, qrCode)
}