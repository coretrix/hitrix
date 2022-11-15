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

func (t *FakeStripeClient) UpdateAccount(accountID string, accountParams *stripe.AccountParams) (*stripe.Account, error) {
	args := t.Called(accountID, accountParams)

	return args.Get(0).(*stripe.Account), args.Error(1)
}

func (t *FakeStripeClient) GetAccount(accountID string, accountParams *stripe.AccountParams) (*stripe.Account, error) {
	args := t.Called(accountID, accountParams)

	return args.Get(0).(*stripe.Account), args.Error(1)
}

func (t *FakeStripeClient) CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	args := t.Called(customerParams)

	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (t *FakeStripeClient) UpdateCustomer(customerID string, customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	args := t.Called(customerID, customerParams)

	return args.Get(0).(*stripe.Customer), args.Error(1)
}

func (t *FakeStripeClient) CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	args := t.Called(checkoutSessionParams)

	return args.Get(0).(*stripe.CheckoutSession), args.Error(1)
}

func (t *FakeStripeClient) GetSubscription(subscriptionID string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	args := t.Called(subscriptionID, params)

	return args.Get(0).(*stripe.Subscription), args.Error(1)
}

func (t *FakeStripeClient) CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	args := t.Called(subscriptionParams)

	return args.Get(0).(*stripe.Subscription), args.Error(1)
}

func (t *FakeStripeClient) UpdateSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	args := t.Called(subscriptionID, subscriptionParams)

	return args.Get(0).(*stripe.Subscription), args.Error(1)
}

func (t *FakeStripeClient) CancelSubscription(subscriptionID string, cancelParams *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	args := t.Called(subscriptionID, cancelParams)

	return args.Get(0).(*stripe.Subscription), args.Error(1)
}

func (t *FakeStripeClient) CreateSetupIntent(setupIntentParams *stripe.SetupIntentParams) (*stripe.SetupIntent, error) {
	args := t.Called(setupIntentParams)

	return args.Get(0).(*stripe.SetupIntent), args.Error(1)
}

func (t *FakeStripeClient) CreateBillingPortalSession(billingPortalSParams *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error) {
	args := t.Called(billingPortalSParams)

	return args.Get(0).(*stripe.BillingPortalSession), args.Error(1)
}

func (t *FakeStripeClient) CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error) {
	args := t.Called(accountLinkParams)

	return args.Get(0).(*stripe.AccountLink), args.Error(1)
}

func (t *FakeStripeClient) GetPaymentIntent(paymentIntentID string, paymentIntentParams *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	args := t.Called(paymentIntentID, paymentIntentParams)

	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (t *FakeStripeClient) CreatePaymentIntentMultiparty(
	paymentIntentParams *stripe.PaymentIntentParams,
	linkedAccountID string,
) (*stripe.PaymentIntent, error) {
	args := t.Called(paymentIntentParams, linkedAccountID)

	return args.Get(0).(*stripe.PaymentIntent), args.Error(1)
}

func (t *FakeStripeClient) CreateRefundMultiparty(refundParams *stripe.RefundParams, linkedAccountID string) (*stripe.Refund, error) {
	args := t.Called(refundParams, linkedAccountID)

	return args.Get(0).(*stripe.Refund), args.Error(1)
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
