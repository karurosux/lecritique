package handlers

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"lecritique/internal/shared/errors"
	"lecritique/internal/shared/middleware"
	"lecritique/internal/shared/response"
	"lecritique/internal/subscription/services"
	"github.com/samber/do"
)

type PaymentHandler struct {
	paymentService services.PaymentService
}

func NewPaymentHandler(i *do.Injector) (*PaymentHandler, error) {
	return &PaymentHandler{
		paymentService: do.MustInvoke[services.PaymentService](i),
	}, nil
}

// CreateCheckoutSession godoc
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
func (h *PaymentHandler) CreateCheckoutSession(c echo.Context) error {
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

// CompleteCheckout godoc
// @Summary Complete a checkout session
// @Description Complete a checkout session after payment
// @Tags payment
// @Accept json
// @Produce json
// @Param request body CompleteCheckoutRequest true "Complete checkout request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/payment/checkout/complete [post]
func (h *PaymentHandler) CompleteCheckout(c echo.Context) error {
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

// CreatePortalSession godoc
// @Summary Create customer portal session
// @Description Create a customer portal session for self-service subscription management
// @Tags payment
// @Produce json
// @Success 200 {object} response.Response{data=PortalResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/portal [post]
// @Security BearerAuth
func (h *PaymentHandler) CreatePortalSession(c echo.Context) error {
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

// HandleWebhook godoc
// @Summary Handle payment webhook
// @Description Handle webhook events from payment provider
// @Tags payment
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/v1/payment/webhook [post]
func (h *PaymentHandler) HandleWebhook(c echo.Context) error {
	// Read the request body
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", "Failed to read request body", http.StatusBadRequest))
	}

	// Get the signature from headers (Stripe uses Stripe-Signature)
	signature := c.Request().Header.Get("Stripe-Signature")
	if signature == "" {
		return response.Error(c, errors.New("BAD_REQUEST", "Missing signature", http.StatusBadRequest))
	}

	// Handle the webhook
	err = h.paymentService.HandleWebhook(c.Request().Context(), body, signature)
	if err != nil {
		return response.Error(c, errors.New("BAD_REQUEST", err.Error(), http.StatusBadRequest))
	}

	return response.Success(c, map[string]string{"message": "Webhook processed"})
}

// GetPaymentMethods godoc
// @Summary List payment methods
// @Description Get list of user's payment methods
// @Tags payment
// @Produce json
// @Success 200 {object} response.Response{data=[]PaymentMethodResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/methods [get]
// @Security BearerAuth
func (h *PaymentHandler) GetPaymentMethods(c echo.Context) error {
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

// SetDefaultPaymentMethod godoc
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
func (h *PaymentHandler) SetDefaultPaymentMethod(c echo.Context) error {
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

// GetInvoices godoc
// @Summary Get invoices
// @Description Get user's invoice history
// @Tags payment
// @Produce json
// @Param limit query int false "Limit number of invoices" default(10)
// @Success 200 {object} response.Response{data=[]InvoiceResponse}
// @Failure 401 {object} response.Response
// @Router /api/v1/payment/invoices [get]
// @Security BearerAuth
func (h *PaymentHandler) GetInvoices(c echo.Context) error {
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

// Request/Response DTOs

type CreateCheckoutRequest struct {
	PlanID uuid.UUID `json:"plan_id" validate:"required"`
}

type CheckoutResponse struct {
	SessionID   string `json:"session_id"`
	CheckoutURL string `json:"checkout_url"`
}

type CompleteCheckoutRequest struct {
	SessionID string `json:"session_id" validate:"required"`
}

type PortalResponse struct {
	PortalURL string `json:"portal_url"`
}

type PaymentMethodResponse struct {
	ID        string                `json:"id"`
	Type      string                `json:"type"`
	IsDefault bool                  `json:"is_default"`
	Card      *CardDetailsResponse  `json:"card,omitempty"`
}

type CardDetailsResponse struct {
	Brand    string `json:"brand"`
	Last4    string `json:"last4"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
}

type SetDefaultPaymentRequest struct {
	PaymentMethodID string `json:"payment_method_id" validate:"required"`
}

type InvoiceResponse struct {
	ID               string     `json:"id"`
	Number           string     `json:"number"`
	Status           string     `json:"status"`
	AmountDue        int64      `json:"amount_due"`
	AmountPaid       int64      `json:"amount_paid"`
	Currency         string     `json:"currency"`
	InvoicePDF       string     `json:"invoice_pdf"`
	HostedInvoiceURL string     `json:"hosted_invoice_url"`
	CreatedAt        time.Time  `json:"created_at"`
	PaidAt           *time.Time `json:"paid_at"`
}