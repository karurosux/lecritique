package qrcodecontroller

import (
	"github.com/labstack/echo/v4"
	qrcodeinterface "kyooar/internal/qrcode/interface"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/logger"
	"kyooar/internal/shared/response"
	"github.com/sirupsen/logrus"
)

type PublicController struct {
	qrCodeService qrcodeinterface.QRCodeService
}

func NewPublicController(qrCodeService qrcodeinterface.QRCodeService) *PublicController {
	return &PublicController{
		qrCodeService: qrCodeService,
	}
}

// @Summary Validate QR code
// @Description Validate a QR code and return associated data
// @Tags public
// @Accept json
// @Produce json
// @Param code path string true "QR Code"
// @Success 200 {object} response.Response{data=qrcodemodel.QRCode}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/public/qr/{code} [get]
func (h *PublicController) ValidateQRCode(c echo.Context) error {
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