package subscriptionservice

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"kyooar/internal/shared/config"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

type paymentService struct {
	provider         subscriptioninterface.PaymentProvider
	subscriptionRepo subscriptioninterface.SubscriptionRepository
	planRepo         subscriptioninterface.SubscriptionPlanRepository
	config           *config.Config
}

func NewPaymentService(
	provider subscriptioninterface.PaymentProvider,
	subscriptionRepo subscriptioninterface.SubscriptionRepository,
	planRepo subscriptioninterface.SubscriptionPlanRepository,
	config *config.Config,
) subscriptioninterface.PaymentService {
	return &paymentService{
		provider:         provider,
		subscriptionRepo: subscriptionRepo,
		planRepo:         planRepo,
		config:           config,
	}
}

func (s *paymentService) GetProvider() subscriptioninterface.PaymentProvider {
	return s.provider
}

func (s *paymentService) GetProviderName() string {
	return s.provider.GetProviderName()
}

func (s *paymentService) CreateCheckout(ctx context.Context, accountID uuid.UUID, planID uuid.UUID) (*subscriptioninterface.CheckoutSession, error) {
	return nil, fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) CompleteCheckout(ctx context.Context, sessionID string) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) CreateOrGetCustomer(ctx context.Context, accountID uuid.UUID, email string) (*subscriptioninterface.Customer, error) {
	return nil, fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) GetCustomerPortalURL(ctx context.Context, accountID uuid.UUID) (string, error) {
	return "", fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) SyncSubscription(ctx context.Context, providerSubscriptionID string) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) UpgradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) DowngradeSubscription(ctx context.Context, accountID uuid.UUID, newPlanID uuid.UUID) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) CancelSubscription(ctx context.Context, accountID uuid.UUID, immediately bool) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) ListPaymentMethods(ctx context.Context, accountID uuid.UUID) ([]*subscriptioninterface.PaymentMethod, error) {
	return nil, fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) SetDefaultPaymentMethod(ctx context.Context, accountID uuid.UUID, paymentMethodID string) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) HandleWebhook(ctx context.Context, payload []byte, signature string) error {
	return fmt.Errorf("payment service not fully implemented yet")
}

func (s *paymentService) GetInvoices(ctx context.Context, accountID uuid.UUID, limit int) ([]*subscriptioninterface.Invoice, error) {
	return nil, fmt.Errorf("payment service not fully implemented yet")
}