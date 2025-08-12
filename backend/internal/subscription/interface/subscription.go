package subscriptioninterface

import (
	"context"
	"time"

	"github.com/google/uuid"
	subscriptionmodel "kyooar/internal/subscription/model"
)

// Repository interfaces
type SubscriptionRepository interface {
	Create(ctx context.Context, subscription *subscriptionmodel.Subscription) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*subscriptionmodel.Subscription, error)
	FindByAccountID(ctx context.Context, accountID uuid.UUID) (*subscriptionmodel.Subscription, error)
	FindByStripeSubscriptionID(ctx context.Context, stripeSubscriptionID string) (*subscriptionmodel.Subscription, error)
	Update(ctx context.Context, subscription *subscriptionmodel.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type SubscriptionPlanRepository interface {
	FindAll(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error)
	FindAllIncludingHidden(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error)
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*subscriptionmodel.SubscriptionPlan, error)
	FindByCode(ctx context.Context, code string) (*subscriptionmodel.SubscriptionPlan, error)
}

type UsageRepository interface {
	Create(ctx context.Context, usage *subscriptionmodel.SubscriptionUsage) error
	FindByID(ctx context.Context, id uuid.UUID, preloads ...string) (*subscriptionmodel.SubscriptionUsage, error)
	FindBySubscriptionAndPeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) (*subscriptionmodel.SubscriptionUsage, error)
	Update(ctx context.Context, usage *subscriptionmodel.SubscriptionUsage) error
	CreateEvent(ctx context.Context, event *subscriptionmodel.UsageEvent) error
	FindEventsBySubscription(ctx context.Context, subscriptionID uuid.UUID, limit int) ([]subscriptionmodel.UsageEvent, error)
	FindEventsByPeriod(ctx context.Context, subscriptionID uuid.UUID, start, end time.Time) ([]subscriptionmodel.UsageEvent, error)
	ResetMonthlyUsage(ctx context.Context) error
}

// Service interfaces
type SubscriptionService interface {
	GetUserSubscription(ctx context.Context, accountID uuid.UUID) (*subscriptionmodel.Subscription, error)
	GetAvailablePlans(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error)
	GetAllPlans(ctx context.Context) ([]subscriptionmodel.SubscriptionPlan, error)
	CanUserCreateOrganization(ctx context.Context, accountID uuid.UUID) (*PermissionResponse, error)
	CreateSubscription(ctx context.Context, req *CreateSubscriptionRequest) (*subscriptionmodel.Subscription, error)
	AssignCustomPlan(ctx context.Context, accountID uuid.UUID, planCode string) error
	CancelSubscription(ctx context.Context, accountID uuid.UUID) error
}

type UsageService interface {
	TrackUsage(ctx context.Context, subscriptionID uuid.UUID, resourceType string, delta int) error
	RecordUsageEvent(ctx context.Context, event *subscriptionmodel.UsageEvent) error
	
	CanAddResource(ctx context.Context, subscriptionID uuid.UUID, resourceType string) (bool, string, error)
	GetCurrentUsage(ctx context.Context, subscriptionID uuid.UUID) (*subscriptionmodel.SubscriptionUsage, error)
	GetUsageForPeriod(ctx context.Context, subscriptionID uuid.UUID, start, end time.Time) (*subscriptionmodel.SubscriptionUsage, error)
	
	InitializeUsagePeriod(ctx context.Context, subscriptionID uuid.UUID, periodStart, periodEnd time.Time) error
	ResetMonthlyUsage(ctx context.Context) error
}

type PaymentService interface {
	GetProvider() PaymentProvider
	GetProviderName() string

	CreateCheckout(ctx context.Context, accountID uuid.UUID, planID uuid.UUID) (*CheckoutSession, error)
	CompleteCheckout(ctx context.Context, sessionID string) error

	CreateOrGetCustomer(ctx context.Context, accountID uuid.UUID, email string) (*Customer, error)
	GetCustomerPortalURL(ctx context.Context, accountID uuid.UUID) (string, error)

	SyncSubscription(ctx context.Context, providerSubscriptionID string) error
	UpgradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error
	DowngradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error
	CancelSubscription(ctx context.Context, accountID uuid.UUID, immediately bool) error

	ListPaymentMethods(ctx context.Context, accountID uuid.UUID) ([]*PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, accountID uuid.UUID, paymentMethodID string) error

	HandleWebhook(ctx context.Context, payload []byte, signature string) error

	GetInvoices(ctx context.Context, accountID uuid.UUID, limit int) ([]*Invoice, error)
}

// DTOs
type PermissionResponse struct {
	CanCreate         bool   `json:"can_create"`
	Reason           string `json:"reason"`
	CurrentCount     int    `json:"current_count"`
	MaxAllowed       int    `json:"max_allowed"`
	SubscriptionStatus string `json:"subscription_status"`
}

type CreateSubscriptionRequest struct {
	AccountID uuid.UUID `json:"account_id"`
	PlanID    uuid.UUID `json:"plan_id"`
}

// Payment provider interfaces and types
type PaymentProvider interface {
	Initialize(config PaymentConfig) error
	GetProviderName() string

	CreateCheckoutSession(params CheckoutParams) (*CheckoutSession, error)
	RetrieveCheckoutSession(sessionID string) (*CheckoutSession, error)

	CreateCustomer(params CustomerParams) (*Customer, error)
	GetCustomer(customerID string) (*Customer, error)
	UpdateCustomer(customerID string, params CustomerParams) error

	CreateSubscription(params SubscriptionParams) (*Subscription, error)
	GetSubscription(subscriptionID string) (*Subscription, error)
	UpdateSubscription(subscriptionID string, params SubscriptionParams) (*Subscription, error)
	CancelSubscription(subscriptionID string, immediately bool) error

	ListPaymentMethods(customerID string) ([]*PaymentMethod, error)
	SetDefaultPaymentMethod(customerID string, paymentMethodID string) error

	GetPortalURL(customerID string) (string, error)
	ValidateWebhook(payload []byte, signature string) (map[string]interface{}, error)

	ListInvoices(customerID string, limit int) ([]*Invoice, error)
}

type PaymentConfig struct {
	SecretKey     string
	WebhookSecret string
}

type CheckoutParams struct {
	CustomerID   string
	CustomerEmail string
	PriceID      string
	SuccessURL   string
	CancelURL    string
	TrialDays    int
	Metadata     map[string]string
}

type CheckoutSession struct {
	ID          string
	URL         string
	Status      string
	CustomerID  string
	SubscriptionID string
	Metadata    map[string]string
}

type CustomerParams struct {
	Email    string
	Name     string
	Metadata map[string]string
}

type Customer struct {
	ID       string
	Email    string
	Name     string
	Metadata map[string]string
}

type SubscriptionParams struct {
	CustomerID string
	PriceID    string
	TrialDays  int
	Metadata   map[string]string
}

type Subscription struct {
	ID                string
	CustomerID        string
	Status            string
	CurrentPeriodEnd  time.Time
	CancelAtPeriodEnd bool
}

type PaymentMethod struct {
	ID        string
	Type      string
	IsDefault bool
	Card      *CardDetails
}

type CardDetails struct {
	Brand    string
	Last4    string
	ExpMonth int
	ExpYear  int
}

type Invoice struct {
	ID               string
	Number           string
	Status           string
	AmountDue        int64
	AmountPaid       int64
	Currency         string
	InvoicePDF       string
	HostedInvoiceURL string
	CreatedAt        time.Time
	PaidAt           *time.Time
}