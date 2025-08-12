package subscriptionservice

import (
	subscriptioninterface "kyooar/internal/subscription/interface"
)

func NewStripeProvider() subscriptioninterface.PaymentProvider {
	return &stripeProvider{}
}