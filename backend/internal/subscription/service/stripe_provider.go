package subscriptionservice

import (
	"fmt"
	subscriptioninterface "kyooar/internal/subscription/interface"
)

type stripeProvider struct {
	config subscriptioninterface.PaymentConfig
}

func (s *stripeProvider) Initialize(config subscriptioninterface.PaymentConfig) error {
	s.config = config
	return nil
}

func (s *stripeProvider) GetProviderName() string {
	return "stripe"
}

func (s *stripeProvider) CreateCheckoutSession(params subscriptioninterface.CheckoutParams) (*subscriptioninterface.CheckoutSession, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) RetrieveCheckoutSession(sessionID string) (*subscriptioninterface.CheckoutSession, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) CreateCustomer(params subscriptioninterface.CustomerParams) (*subscriptioninterface.Customer, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) GetCustomer(customerID string) (*subscriptioninterface.Customer, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) UpdateCustomer(customerID string, params subscriptioninterface.CustomerParams) error {
	return fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) CreateSubscription(params subscriptioninterface.SubscriptionParams) (*subscriptioninterface.Subscription, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) GetSubscription(subscriptionID string) (*subscriptioninterface.Subscription, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) UpdateSubscription(subscriptionID string, params subscriptioninterface.SubscriptionParams) (*subscriptioninterface.Subscription, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) CancelSubscription(subscriptionID string, immediately bool) error {
	return fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) ListPaymentMethods(customerID string) ([]*subscriptioninterface.PaymentMethod, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) SetDefaultPaymentMethod(customerID string, paymentMethodID string) error {
	return fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) GetPortalURL(customerID string) (string, error) {
	return "", fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) ValidateWebhook(payload []byte, signature string) (map[string]interface{}, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}

func (s *stripeProvider) ListInvoices(customerID string, limit int) ([]*subscriptioninterface.Invoice, error) {
	return nil, fmt.Errorf("stripe provider not fully implemented yet")
}