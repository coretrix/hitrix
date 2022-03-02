package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/stripe/stripe-go/v72"
)

type FakeStripeClient struct {
	mock.Mock
}

func (t *FakeStripeClient) CreateAccount(accountParams *stripe.AccountParams) (*stripe.Account, error) {
	args := t.Called(accountParams)
	return args.Get(0).(*stripe.Account), args.Error(1)
}

func (t *FakeStripeClient) CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error) {
	args := t.Called(accountLinkParams)
	return args.Get(0).(*stripe.AccountLink), args.Error(1)
}

func (t *FakeStripeClient) CreatePaymentIntentMultiparty(paymentIntentParams *stripe.PaymentIntentParams, linkedAccountID string) (*stripe.PaymentIntent, error) {
	args := t.Called(paymentIntentParams, linkedAccountID)
	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
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
