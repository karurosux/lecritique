package subscriptioncontroller

import (
	"io"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"kyooar/internal/shared/errors"
	"kyooar/internal/shared/middleware"
	"kyooar/internal/shared/response"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

type PaymentController struct {
	paymentService subscriptioninterface.PaymentService
}

func NewPaymentController(paymentService subscriptioninterface.PaymentService) *PaymentController {
	return &PaymentController{
		paymentService: paymentService,
	}
}

// @Summary Create a checkout session
// @Description Create a payment checkout session for a subscription plan
// @Tags payment
// @Accept json
// @Produce json
// @Param request body CreateCheckoutRequest true "Checkout request"
// @Success 200 {object} response.Response{data=CheckoutResponse}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/checkout [post]
// @Security BearerAuth
func (h *PaymentController) CreateCheckoutSession(c echo.Context) error {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	var req CreateCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Invalid request", http.StatusBadRequest))
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, errors.New("VALIDATION_ERROR", err.Error(), http.StatusBadRequest))
	}

	session, err := h.paymentService.CreateCheckout(c.Request().Context(), accountID, req.PlanID)
	if err != nil {
		return response.Error(c, errors.New("INTERNAL_ERROR", err.Error(), http.StatusInternalServerError))
	}

	return response.Success(c, CheckoutResponse{
		SessionID:   session.ID,
		CheckoutURL: session.URL,
	})
}

// @Summary Complete a checkout session
// @Description Complete a checkout session after payment
// @Tags payment
// @Accept json
// @Produce json
// @Param request body CompleteCheckoutRequest true "Complete checkout request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/payment/checkout/complete [post]
func (h *PaymentController) CompleteCheckout(c echo.Context) error {
	var req CompleteCheckoutRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Invalid request", http.StatusBadRequest))
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, errors.New("VALIDATION_ERROR", err.Error(), http.StatusBadRequest))
	}

	err := h.paymentService.CompleteCheckout(c.Request().Context(), req.SessionID)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	return response.Success(c, map[string]string{"message": "Checkout completed successfully"})
}

// @Summary Create customer portal session
// @Description Create a customer portal session for self-service subscription management
// @Tags payment
// @Produce json
// @Success 200 {object} response.Response{data=PortalResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/portal [post]
// @Security BearerAuth
func (h *PaymentController) CreatePortalSession(c echo.Context) error {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	portalURL, err := h.paymentService.GetCustomerPortalURL(c.Request().Context(), accountID)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	return response.Success(c, PortalResponse{
		PortalURL: portalURL,
	})
}

// @Summary Handle payment webhook
// @Description Handle webhook events from payment provider
// @Tags payment
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/payment/webhook [post]
func (h *PaymentController) HandleWebhook(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Failed to read request body", http.StatusBadRequest))
	}

	signature := c.Request().Header.Get("Stripe-Signature")
	if signature == "" {
		return response.Error(c, errors.New("BAD_REQUEST", "Missing signature", http.StatusBadRequest))
	}

	err = h.paymentService.HandleWebhook(c.Request().Context(), body, signature)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	return response.Success(c, map[string]string{"message": "Webhook processed"})
}

// @Summary List payment methods
// @Description Get list of user's payment methods
// @Tags payment
// @Produce json
// @Success 200 {object} response.Response{data=[]PaymentMethodResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/methods [get]
// @Security BearerAuth
func (h *PaymentController) GetPaymentMethods(c echo.Context) error {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	methods, err := h.paymentService.ListPaymentMethods(c.Request().Context(), accountID)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	var resp []PaymentMethodResponse
	for _, method := range methods {
		r := PaymentMethodResponse{
			ID:        method.ID,
			Type:      method.Type,
			IsDefault: method.IsDefault,
		}
		if method.Card != nil {
			r.Card = &CardDetailsResponse{
				Brand:    method.Card.Brand,
				Last4:    method.Card.Last4,
				ExpMonth: method.Card.ExpMonth,
				ExpYear:  method.Card.ExpYear,
			}
		}
		resp = append(resp, r)
	}

	return response.Success(c, resp)
}

// @Summary Set default payment method
// @Description Set a payment method as default
// @Tags payment
// @Accept json
// @Produce json
// @Param request body SetDefaultPaymentRequest true "Set default payment request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/methods/default [post]
// @Security BearerAuth
func (h *PaymentController) SetDefaultPaymentMethod(c echo.Context) error {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	var req SetDefaultPaymentRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Invalid request", http.StatusBadRequest))
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, errors.New("VALIDATION_ERROR", err.Error(), http.StatusBadRequest))
	}

	err = h.paymentService.SetDefaultPaymentMethod(c.Request().Context(), accountID, req.PaymentMethodID)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	return response.Success(c, map[string]string{"message": "Default payment method updated"})
}

// @Summary Get invoices
// @Description Get user's invoice history
// @Tags payment
// @Produce json
// @Param limit query int false "Limit number of invoices" default(10)
// @Success 200 {object} response.Response{data=[]InvoiceResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/invoices [get]
// @Security BearerAuth
func (h *PaymentController) GetInvoices(c echo.Context) error {
	accountID, err := middleware.GetAccountID(c)
	if err != nil {
		return response.Error(c, err)
	}

	limit := 10
	if l := c.QueryParam("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		}
	}

	invoices, err := h.paymentService.GetInvoices(c.Request().Context(), accountID, limit)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	var resp []InvoiceResponse
	for _, invoice := range invoices {
		resp = append(resp, InvoiceResponse{
			ID:               invoice.ID,
			Number:           invoice.Number,
			Status:           invoice.Status,
			AmountDue:        invoice.AmountDue,
			AmountPaid:       invoice.AmountPaid,
			Currency:         invoice.Currency,
			InvoicePDF:       invoice.InvoicePDF,
			HostedInvoiceURL: invoice.HostedInvoiceURL,
			CreatedAt:        invoice.CreatedAt,
			PaidAt:           invoice.PaidAt,
		})
	}

	return response.Success(c, resp)
}