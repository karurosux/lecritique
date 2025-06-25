package services

import (
	"context"
	"time"
)

// PaymentProvider defines the interface for payment processing
type PaymentProvider interface {
	// Initialization
	Initialize(config PaymentConfig) error
	GetProviderName() string

	// Checkout
	CreateCheckoutSession(ctx context.Context, options CheckoutOptions) (*CheckoutSession, error)
	GetCheckoutSession(ctx context.Context, sessionID string) (*CheckoutSession, error)

	// Customer management
	CreateCustomer(ctx context.Context, customer CustomerInfo) (*Customer, error)
	UpdateCustomer(ctx context.Context, customerID string, updates CustomerInfo) (*Customer, error)
	GetCustomer(ctx context.Context, customerID string) (*Customer, error)

	// Subscription management
	CreateSubscription(ctx context.Context, options SubscriptionOptions) (*PaymentSubscription, error)
	UpdateSubscription(ctx context.Context, subscriptionID string, options UpdateSubscriptionOptions) (*PaymentSubscription, error)
	CancelSubscription(ctx context.Context, subscriptionID string, immediately bool) error
	GetSubscription(ctx context.Context, subscriptionID string) (*PaymentSubscription, error)

	// Payment methods
	AttachPaymentMethod(ctx context.Context, customerID string, paymentMethodID string) error
	DetachPaymentMethod(ctx context.Context, paymentMethodID string) error
	ListPaymentMethods(ctx context.Context, customerID string) ([]*PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, customerID string, paymentMethodID string) error

	// Billing portal
	CreatePortalSession(ctx context.Context, customerID string, returnURL string) (*PortalSession, error)

	// Webhooks
	ConstructWebhookEvent(payload []byte, signature string) (*WebhookEvent, error)
	HandleWebhookEvent(ctx context.Context, event *WebhookEvent) error

	// Invoices
	GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error)
	ListInvoices(ctx context.Context, customerID string, limit int) ([]*Invoice, error)
}

// PaymentConfig holds provider-specific configuration
type PaymentConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	Extra          map[string]interface{} // Provider-specific config
}

// CheckoutOptions for creating a checkout session
type CheckoutOptions struct {
	CustomerID          string
	PriceID             string
	SuccessURL          string
	CancelURL           string
	TrialPeriodDays     int
	AllowPromotionCodes bool
	Metadata            map[string]string
}

// CheckoutSession represents a payment checkout session
type CheckoutSession struct {
	ID                string
	URL               string
	Status            string
	CustomerID        string
	SubscriptionID    string
	PaymentIntentID   string
	AmountTotal       int64
	Currency          string
	ExpiresAt         time.Time
	Metadata          map[string]string
}

// CustomerInfo for creating/updating customers
type CustomerInfo struct {
	Email       string
	Name        string
	Phone       string
	Description string
	Metadata    map[string]string
}

// Customer represents a payment provider customer
type Customer struct {
	ID               string
	Email            string
	Name             string
	Phone            string
	Description      string
	DefaultPaymentID string
	Metadata         map[string]string
	CreatedAt        time.Time
}

// SubscriptionOptions for creating subscriptions
type SubscriptionOptions struct {
	CustomerID        string
	PriceID           string
	TrialPeriodDays   int
	DefaultPaymentID  string
	Metadata          map[string]string
}

// UpdateSubscriptionOptions for updating subscriptions
type UpdateSubscriptionOptions struct {
	PriceID             string
	ProrationBehavior   string // "create_prorations", "none", "always_invoice"
	CancelAtPeriodEnd   bool
	DefaultPaymentID    string
	Metadata            map[string]string
}

// PaymentSubscription represents a subscription in the payment provider
type PaymentSubscription struct {
	ID                   string
	CustomerID           string
	Status               string
	CurrentPeriodStart   time.Time
	CurrentPeriodEnd     time.Time
	CancelAt             *time.Time
	CanceledAt           *time.Time
	EndedAt              *time.Time
	TrialStart           *time.Time
	TrialEnd             *time.Time
	Items                []SubscriptionItem
	DefaultPaymentID     string
	LatestInvoiceID      string
	Metadata             map[string]string
}

// SubscriptionItem represents a line item in a subscription
type SubscriptionItem struct {
	ID       string
	PriceID  string
	Quantity int64
}

// PaymentMethod represents a customer's payment method
type PaymentMethod struct {
	ID         string
	Type       string // "card", "bank_account", etc
	IsDefault  bool
	Card       *CardDetails
	CreatedAt  time.Time
}

// CardDetails for card payment methods
type CardDetails struct {
	Brand    string // visa, mastercard, amex, etc
	Last4    string
	ExpMonth int
	ExpYear  int
	Country  string
}

// PortalSession for customer self-service
type PortalSession struct {
	ID        string
	URL       string
	ReturnURL string
	CreatedAt time.Time
}

// WebhookEvent represents an incoming webhook event
type WebhookEvent struct {
	ID        string
	Type      string
	Data      interface{}
	CreatedAt time.Time
}

// Common webhook event types
const (
	WebhookCheckoutCompleted        = "checkout.session.completed"
	WebhookSubscriptionCreated      = "customer.subscription.created"
	WebhookSubscriptionUpdated      = "customer.subscription.updated"
	WebhookSubscriptionDeleted      = "customer.subscription.deleted"
	WebhookInvoicePaymentSucceeded  = "invoice.payment_succeeded"
	WebhookInvoicePaymentFailed     = "invoice.payment_failed"
	WebhookPaymentMethodAttached    = "payment_method.attached"
	WebhookPaymentMethodDetached    = "payment_method.detached"
)

// Invoice represents a billing invoice
type Invoice struct {
	ID                string
	Number            string
	CustomerID        string
	SubscriptionID    string
	Status            string
	AmountDue         int64
	AmountPaid        int64
	Currency          string
	InvoicePDF        string
	HostedInvoiceURL  string
	CreatedAt         time.Time
	PaidAt            *time.Time
	DueDate           *time.Time
	Lines             []InvoiceLine
}

// InvoiceLine represents a line item on an invoice
type InvoiceLine struct {
	Description string
	Quantity    int64
	UnitAmount  int64
	Amount      int64
}