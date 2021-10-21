package mocks

import (
	"github.com/stripe/stripe-go/v72"

	"github.com/stretchr/testify/mock"
)

type FakeStripeClient struct {
	mock.Mock
}

func (t *FakeStripeClient) ConstructWebhookEvent(reqBody []byte, signature string, webhookKey string) (stripe.Event, error) {
	return t.Called(reqBody, signature, webhookKey).Get(0).(stripe.Event), t.Called(reqBody, signature, webhookKey).Error(1)
}

func (t *FakeStripeClient) NewCheckoutSession(
	paymentMethods []string,
	mode, successURL, CancelURL string,
	lineItems []*stripe.CheckoutSessionLineItemParams,
	_ []*stripe.CheckoutSessionDiscountParams) *stripe.CheckoutSession {
	return t.Called(paymentMethods, mode, successURL, CancelURL, lineItems).Get(0).(*stripe.CheckoutSession)
}
