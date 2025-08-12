package qrcodeconstants

type QRCodeType string

const (
	QRCodeTypeTable    QRCodeType = "table"
	QRCodeTypeLocation QRCodeType = "location"
	QRCodeTypeTakeaway QRCodeType = "takeaway"
	QRCodeTypeDelivery QRCodeType = "delivery"
	QRCodeTypeGeneral  QRCodeType = "general"
)