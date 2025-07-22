package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"kyooar/internal/shared/config"
	"kyooar/internal/subscription/models"
	subscriptionRepos "kyooar/internal/subscription/repositories"
	"github.com/samber/do"
)

// PaymentService handles payment operations using the configured provider
type PaymentService interface {
	// Provider management
	GetProvider() PaymentProvider
	GetProviderName() string

	// Checkout operations
	CreateCheckout(ctx context.Context, accountID uuid.UUID, planID uuid.UUID) (*CheckoutSession, error)
	CompleteCheckout(ctx context.Context, sessionID string) error

	// Customer operations
	CreateOrGetCustomer(ctx context.Context, accountID uuid.UUID, email string) (*Customer, error)
	GetCustomerPortalURL(ctx context.Context, accountID uuid.UUID) (string, error)

	// Subscription operations
	SyncSubscription(ctx context.Context, providerSubscriptionID string) error
	UpgradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error
	DowngradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error
	CancelSubscription(ctx context.Context, accountID uuid.UUID, immediately bool) error

	// Payment method operations
	ListPaymentMethods(ctx context.Context, accountID uuid.UUID) ([]*PaymentMethod, error)
	SetDefaultPaymentMethod(ctx context.Context, accountID uuid.UUID, paymentMethodID string) error

	// Webhook handling
	HandleWebhook(ctx context.Context, payload []byte, signature string) error

	// Invoice operations
	GetInvoices(ctx context.Context, accountID uuid.UUID, limit int) ([]*Invoice, error)
}

type paymentService struct {
	provider          PaymentProvider
	subscriptionRepo  subscriptionRepos.SubscriptionRepository
	planRepo          subscriptionRepos.SubscriptionPlanRepository
	config            *config.Config
	customerCache     map[uuid.UUID]string // accountID -> customerID cache
	customerCacheMux  sync.RWMutex
}

// NewPaymentService creates a new payment service with the specified provider
func NewPaymentService(i *do.Injector) (PaymentService, error) {
	config := do.MustInvoke[*config.Config](i)
	subscriptionRepo := do.MustInvoke[subscriptionRepos.SubscriptionRepository](i)
	planRepo := do.MustInvoke[subscriptionRepos.SubscriptionPlanRepository](i)
	
	// Default to stripe provider (for now, only stripe is supported)
	providerName := "stripe"
	
	var provider PaymentProvider

	switch providerName {
	case "stripe":
		provider = NewStripeProvider()
	// Add more providers here as needed
	// case "paypal":
	//     provider = NewPayPalProvider()
	default:
		return nil, fmt.Errorf("unsupported payment provider: %s", providerName)
	}

	// Initialize provider with config
	providerConfig := PaymentConfig{
		SecretKey:     config.Stripe.SecretKey,
		WebhookSecret: config.Stripe.WebhookSecret,
	}

	if err := provider.Initialize(providerConfig); err != nil {
		return nil, fmt.Errorf("failed to initialize payment provider: %w", err)
	}

	return &paymentService{
		provider:         provider,
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		config:           config,
		customerCache:    make(map[uuid.UUID]string),
	}, nil
}

func (s *paymentService) GetProvider() PaymentProvider {
	return s.provider
}

func (s *paymentService) GetProviderName() string {
	return s.provider.GetProviderName()
}

func (s *paymentService) CreateCheckout(ctx context.Context, accountID uuid.UUID, planID uuid.UUID) (*CheckoutSession, error) {
	// Get the plan
	plan, err := s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return nil, fmt.Errorf("plan not found: %w", err)
	}

	if plan.StripePriceID == "" {
		return nil, fmt.Errorf("plan does not have a payment provider price ID")
	}

	// Get or create customer
	subscription, _ := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	var customerID string
	if subscription != nil && subscription.StripeCustomerID != "" {
		customerID = subscription.StripeCustomerID
	}

	// Create checkout session
	options := CheckoutOptions{
		CustomerID:          customerID,
		PriceID:             plan.StripePriceID,
		SuccessURL:          fmt.Sprintf("%s/subscription/success?session_id={CHECKOUT_SESSION_ID}", s.config.App.FrontendURL),
		CancelURL:           fmt.Sprintf("%s/subscription/cancel", s.config.App.FrontendURL),
		TrialPeriodDays:     plan.TrialDays,
		AllowPromotionCodes: true,
		Metadata: map[string]string{
			"account_id": accountID.String(),
			"plan_id":    planID.String(),
		},
	}

	return s.provider.CreateCheckoutSession(ctx, options)
}

func (s *paymentService) CompleteCheckout(ctx context.Context, sessionID string) error {
	// Get checkout session details
	session, err := s.provider.GetCheckoutSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("failed to get checkout session: %w", err)
	}

	if session.Status != "complete" {
		return fmt.Errorf("checkout session not completed")
	}

	// Get account ID from metadata
	accountIDStr, ok := session.Metadata["account_id"]
	if !ok {
		return fmt.Errorf("account ID not found in session metadata")
	}

	accountID, err := uuid.Parse(accountIDStr)
	if err != nil {
		return fmt.Errorf("invalid account ID: %w", err)
	}

	// Get plan ID from metadata
	planIDStr, ok := session.Metadata["plan_id"]
	if !ok {
		return fmt.Errorf("plan ID not found in session metadata")
	}

	planID, err := uuid.Parse(planIDStr)
	if err != nil {
		return fmt.Errorf("invalid plan ID: %w", err)
	}

	// Verify the plan exists
	_, err = s.planRepo.FindByID(ctx, planID)
	if err != nil {
		return fmt.Errorf("plan not found: %w", err)
	}

	// Create or update subscription
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		// Create new subscription
		subscription = &models.Subscription{
			AccountID:            accountID,
			PlanID:               planID,
			Status:               models.SubscriptionActive,
			StripeCustomerID:     session.CustomerID,
			StripeSubscriptionID: session.SubscriptionID,
		}
		err = s.subscriptionRepo.Create(ctx, subscription)
	} else {
		// Update existing subscription
		subscription.PlanID = planID
		subscription.Status = models.SubscriptionActive
		subscription.StripeCustomerID = session.CustomerID
		subscription.StripeSubscriptionID = session.SubscriptionID
		err = s.subscriptionRepo.Update(ctx, subscription)
	}

	if err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}

	// Sync subscription details from provider
	return s.SyncSubscription(ctx, session.SubscriptionID)
}

func (s *paymentService) CreateOrGetCustomer(ctx context.Context, accountID uuid.UUID, email string) (*Customer, error) {
	// Check cache first
	s.customerCacheMux.RLock()
	customerID, cached := s.customerCache[accountID]
	s.customerCacheMux.RUnlock()

	if cached {
		return s.provider.GetCustomer(ctx, customerID)
	}

	// Check if customer exists in subscription
	subscription, _ := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if subscription != nil && subscription.StripeCustomerID != "" {
		// Cache it
		s.customerCacheMux.Lock()
		s.customerCache[accountID] = subscription.StripeCustomerID
		s.customerCacheMux.Unlock()

		return s.provider.GetCustomer(ctx, subscription.StripeCustomerID)
	}

	// Create new customer
	customer, err := s.provider.CreateCustomer(ctx, CustomerInfo{
		Email: email,
		Metadata: map[string]string{
			"account_id": accountID.String(),
		},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Cache it
	s.customerCacheMux.Lock()
	s.customerCache[accountID] = customer.ID
	s.customerCacheMux.Unlock()

	return customer, nil
}

func (s *paymentService) GetCustomerPortalURL(ctx context.Context, accountID uuid.UUID) (string, error) {
	// Get subscription to find customer ID
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil || subscription.StripeCustomerID == "" {
		return "", fmt.Errorf("no customer found for account")
	}

	// Create portal session
	portalSession, err := s.provider.CreatePortalSession(
		ctx,
		subscription.StripeCustomerID,
		fmt.Sprintf("%s/settings/subscription", s.config.App.FrontendURL),
	)

	if err != nil {
		return "", fmt.Errorf("failed to create portal session: %w", err)
	}

	return portalSession.URL, nil
}

func (s *paymentService) SyncSubscription(ctx context.Context, providerSubscriptionID string) error {
	// Get subscription details from provider
	providerSub, err := s.provider.GetSubscription(ctx, providerSubscriptionID)
	if err != nil {
		return fmt.Errorf("failed to get subscription from provider: %w", err)
	}

	// Find our subscription
	subscription, err := s.subscriptionRepo.FindByStripeSubscriptionID(ctx, providerSubscriptionID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Update subscription with provider data
	subscription.CurrentPeriodStart = providerSub.CurrentPeriodStart
	subscription.CurrentPeriodEnd = providerSub.CurrentPeriodEnd
	subscription.CancelAt = providerSub.CancelAt
	subscription.CancelledAt = providerSub.CanceledAt

	// Map provider status to our status
	switch providerSub.Status {
	case "active", "trialing":
		subscription.Status = models.SubscriptionActive
	case "canceled":
		subscription.Status = models.SubscriptionCanceled
	case "past_due", "unpaid":
		subscription.Status = models.SubscriptionPending
	default:
		subscription.Status = models.SubscriptionExpired
	}

	return s.subscriptionRepo.Update(ctx, subscription)
}

func (s *paymentService) UpgradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error {
	// Implementation would handle plan upgrades with proper proration
	return fmt.Errorf("not implemented")
}

func (s *paymentService) DowngradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error {
	// Implementation would handle plan downgrades
	return fmt.Errorf("not implemented")
}

func (s *paymentService) CancelSubscription(ctx context.Context, accountID uuid.UUID, immediately bool) error {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	if subscription.StripeSubscriptionID == "" {
		// Just mark as canceled in our database
		subscription.Status = models.SubscriptionCanceled
		return s.subscriptionRepo.Update(ctx, subscription)
	}

	// Cancel in provider
	err = s.provider.CancelSubscription(ctx, subscription.StripeSubscriptionID, immediately)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	// Sync the updated status
	return s.SyncSubscription(ctx, subscription.StripeSubscriptionID)
}

func (s *paymentService) ListPaymentMethods(ctx context.Context, accountID uuid.UUID) ([]*PaymentMethod, error) {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil || subscription.StripeCustomerID == "" {
		return nil, fmt.Errorf("no customer found for account")
	}

	return s.provider.ListPaymentMethods(ctx, subscription.StripeCustomerID)
}

func (s *paymentService) SetDefaultPaymentMethod(ctx context.Context, accountID uuid.UUID, paymentMethodID string) error {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil || subscription.StripeCustomerID == "" {
		return fmt.Errorf("no customer found for account")
	}

	return s.provider.SetDefaultPaymentMethod(ctx, subscription.StripeCustomerID, paymentMethodID)
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	event, err := s.provider.ConstructWebhookEvent(payload, signature)
	if err != nil {
		return fmt.Errorf("failed to construct webhook event: %w", err)
	}

	// Handle different event types
	switch event.Type {
	case WebhookCheckoutCompleted:
		// Handle checkout completion
		// This is already handled by CompleteCheckout
		
	case WebhookSubscriptionUpdated:
		// Sync subscription changes
		// Extract subscription ID from event data and sync
		
	case WebhookInvoicePaymentFailed:
		// Handle payment failures
		// Mark subscription as pending, send notification
		
	// Add more event handlers as needed
	}

	return s.provider.HandleWebhookEvent(ctx, event)
}

func (s *paymentService) GetInvoices(ctx context.Context, accountID uuid.UUID, limit int) ([]*Invoice, error) {
	subscription, err := s.subscriptionRepo.FindByAccountID(ctx, accountID)
	if err != nil || subscription.StripeCustomerID == "" {
		return nil, fmt.Errorf("no customer found for account")
	}

	return s.provider.ListInvoices(ctx, subscription.StripeCustomerID, limit)
}