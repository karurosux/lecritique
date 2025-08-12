package services

import (
	"context"
	"time"
)

type PaymentProvider interface {
	Initialize(config PaymentConfig) error
	GetProviderName() string

	CreateCheckoutSession(ctx context.Context, options CheckoutOptions) (*CheckoutSession, error)
	GetCheckoutSession(ctx context.Context, sessionID string) (*CheckoutSession, error)

	CreateCustomer(ctx context.Context, customer CustomerInfo) (*Customer, error)
	UpdateCustomer(ctx context.Context, customerID string, updates CustomerInfo) (*Customer, error)
	GetCustomer(ctx context.Context, customerID string) (*Customer, error)

	CreateSubscription(ctx context.Context, options SubscriptionOptions) (*PaymentSubscription, error)
	UpdateSubscription(ctx context.Context, subscriptionID string, options UpdateSubscriptionOptions) (*PaymentSubscription, error)
	CancelSubscription(ctx context.Context, subscriptionID string, immediately bool) error
	GetSubscription(ctx context.Context, subscriptionID string) (*PaymentSubscription, error)

	AttachPaymentMethod(ctx context.Context, customerID string, paymentMethodID string) error
	DetachPaymentMethod(ctx context.Context, paymentMethodID string) error
	ListPaymentMethods(ctx context.Context, customerID string) ([]*PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, customerID string, paymentMethodID string) error

	CreatePortalSession(ctx context.Context, customerID string, returnURL string) (*PortalSession, error)

	ConstructWebhookEvent(payload []byte, signature string) (*WebhookEvent, error)
	HandleWebhookEvent(ctx context.Context, event *WebhookEvent) error

	GetInvoice(ctx context.Context, invoiceID string) (*Invoice, error)
	ListInvoices(ctx context.Context, customerID string, limit int) ([]*Invoice, error)
}

type PaymentConfig struct {
	SecretKey      string
	PublishableKey string
	WebhookSecret  string
	Extra          map[string]interface{}
}

type CheckoutOptions struct {
	CustomerID          string
	PriceID             string
	SuccessURL          string
	CancelURL           string
	TrialPeriodDays     int
	AllowPromotionCodes bool
	Metadata            map[string]string
}

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

type CustomerInfo struct {
	Email       string
	Name        string
	Phone       string
	Description string
	Metadata    map[string]string
}

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

type SubscriptionOptions struct {
	CustomerID        string
	PriceID           string
	TrialPeriodDays   int
	DefaultPaymentID  string
	Metadata          map[string]string
}

type UpdateSubscriptionOptions struct {
	PriceID             string
	ProrationBehavior   string
	CancelAtPeriodEnd   bool
	DefaultPaymentID    string
	Metadata            map[string]string
}

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

type SubscriptionItem struct {
	ID       string
	PriceID  string
	Quantity int64
}

type PaymentMethod struct {
	ID         string
	Type       string
	IsDefault  bool
	Card       *CardDetails
	CreatedAt  time.Time
}

type CardDetails struct {
	Brand    string
	Last4    string
	ExpMonth int
	ExpYear  int
	Country  string
}

type PortalSession struct {
	ID        string
	URL       string
	ReturnURL string
	CreatedAt time.Time
}

type WebhookEvent struct {
	ID        string
	Type      string
	Data      interface{}
	CreatedAt time.Time
}

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

type InvoiceLine struct {
	Description string
	Quantity    int64
	UnitAmount  int64
	Amount      int64
}