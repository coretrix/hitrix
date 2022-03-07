package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"
	"github.com/stripe/stripe-go/v72/accountlink"
	portalsession "github.com/stripe/stripe-go/v72/billingportal/session"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/setupintent"
	"github.com/stripe/stripe-go/v72/sub"
	"github.com/stripe/stripe-go/v72/webhook"
)

type Stripe struct {
	webhookSecrets map[string]string
}

func NewStripe(token string, webhookSecrets map[string]string) *Stripe {
	stripe.Key = token

	return &Stripe{
		webhookSecrets: webhookSecrets,
	}
}

func (s *Stripe) CreateAccount(accountParams *stripe.AccountParams) (*stripe.Account, error) {
	return account.New(accountParams)
}

func (s *Stripe) CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	return customer.New(customerParams)
}

func (s *Stripe) UpdateCustomer(customerID string, customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	return customer.Update(customerID, customerParams)
}

func (s *Stripe) CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	return session.New(checkoutSessionParams)
}

func (s *Stripe) CreateBillingPortalSession(billingPortalSessionParams *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error) {
	return portalsession.New(billingPortalSessionParams)
}

func (s *Stripe) CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return sub.New(subscriptionParams)
}

func (s *Stripe) UpdateSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	return sub.Update(subscriptionID, subscriptionParams)
}

func (s *Stripe) CancelSubscription(subscriptionID string, subscriptionCancelParams *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	return sub.Cancel(subscriptionID, subscriptionCancelParams)
}

func (s *Stripe) CreateSetupIntent(setupIntentParams *stripe.SetupIntentParams) (*stripe.SetupIntent, error) {
	return setupintent.New(setupIntentParams)
}

func (s *Stripe) CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error) {
	return accountlink.New(accountLinkParams)
}

func (s *Stripe) CreatePaymentIntentMultiparty(paymentIntentParams *stripe.PaymentIntentParams, linkedAccountID string) (*stripe.PaymentIntent, error) {
	paymentIntentParams.SetStripeAccount(linkedAccountID)
	return paymentintent.New(paymentIntentParams)
}

func (s *Stripe) NewCheckoutSession(paymentMethods []string, mode, successURL, CancelURL string, lineItems []*stripe.CheckoutSessionLineItemParams, discounts []*stripe.CheckoutSessionDiscountParams) *stripe.CheckoutSession {
	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice(paymentMethods),
		LineItems:          lineItems,
		Mode:               stripe.String(mode),
		SuccessURL:         stripe.String(successURL),
		CancelURL:          stripe.String(CancelURL),
		Discounts:          discounts,
	}
	checkoutSession, err := session.New(params)

	if err != nil {
		panic("failed creating new session for stripe checkout" + err.Error())
	}

	return checkoutSession
}

func (s *Stripe) ConstructWebhookEvent(reqBody []byte, signature string, webhookKey string) (stripe.Event, error) {
	secret, ok := s.webhookSecrets[webhookKey]
	if !ok {
		panic("stripe webhook secret [" + webhookKey + "] not found")
	}
	event, err := webhook.ConstructEvent(reqBody, signature, secret)
	return event, err
}

type IStripe interface {
	CreateAccount(accountParams *stripe.AccountParams) (*stripe.Account, error)
	CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error)
	UpdateCustomer(customerID string, customerParams *stripe.CustomerParams) (*stripe.Customer, error)
	CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error)
	CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	UpdateSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
	CancelSubscription(subscriptionID string, subscriptionCancelParams *stripe.SubscriptionCancelParams) (*stripe.Subscription, error)
	CreateSetupIntent(setupIntentParams *stripe.SetupIntentParams) (*stripe.SetupIntent, error)
	CreateBillingPortalSession(billingPortalSessionParams *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error)
	CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error)
	CreatePaymentIntentMultiparty(paymentIntentParams *stripe.PaymentIntentParams, linkedAccountID string) (*stripe.PaymentIntent, error)
	ConstructWebhookEvent(reqBody []byte, signature string, webhookKey string) (stripe.Event, error)
	NewCheckoutSession(paymentMethods []string, mode, successURL, CancelURL string, lineItems []*stripe.CheckoutSessionLineItemParams, discounts []*stripe.CheckoutSessionDiscountParams) *stripe.CheckoutSession
}
