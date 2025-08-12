package subscriptioncontroller

import (
	"time"
	"github.com/google/uuid"
)

// Subscription DTOs
type CreateSubscriptionRequest struct {
	PlanID string `json:"plan_id" validate:"required,uuid"`
}

// Payment DTOs
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