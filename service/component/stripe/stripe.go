package stripe

import (
	"github.com/coretrix/hitrix/service/component/app"
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

const Env = "env"

type Stripe struct {
	webhookSecrets map[string]string
	appService     *app.App
}

func NewStripe(token string, webhookSecrets map[string]string, appService *app.App) *Stripe {
	stripe.Key = token

	return &Stripe{
		webhookSecrets: webhookSecrets,
		appService:     appService,
	}
}

func (s *Stripe) CreateAccount(accountParams *stripe.AccountParams) (*stripe.Account, error) {
	if accountParams.Params.Metadata == nil {
		accountParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		accountParams.Params.Metadata[Env] = s.appService.Mode
	}
	return account.New(accountParams)
}

func (s *Stripe) CreateCustomer(customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	if customerParams.Params.Metadata == nil {
		customerParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		customerParams.Params.Metadata[Env] = s.appService.Mode
	}
	return customer.New(customerParams)
}

func (s *Stripe) UpdateCustomer(customerID string, customerParams *stripe.CustomerParams) (*stripe.Customer, error) {
	if customerParams.Params.Metadata == nil {
		customerParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		customerParams.Params.Metadata[Env] = s.appService.Mode
	}
	return customer.Update(customerID, customerParams)
}

func (s *Stripe) CreateCheckoutSession(checkoutSessionParams *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	if checkoutSessionParams.Params.Metadata == nil {
		checkoutSessionParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		checkoutSessionParams.Params.Metadata[Env] = s.appService.Mode
	}
	return session.New(checkoutSessionParams)
}

func (s *Stripe) CreateBillingPortalSession(billingPortalSessionParams *stripe.BillingPortalSessionParams) (*stripe.BillingPortalSession, error) {
	return portalsession.New(billingPortalSessionParams)
}

func (s *Stripe) GetSubscription(subscriptionID string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	if params.Params.Metadata == nil {
		params.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		params.Params.Metadata[Env] = s.appService.Mode
	}
	return sub.Get(subscriptionID, params)
}

func (s *Stripe) CreateSubscription(subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	if subscriptionParams.Params.Metadata == nil {
		subscriptionParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		subscriptionParams.Params.Metadata[Env] = s.appService.Mode
	}
	return sub.New(subscriptionParams)
}

func (s *Stripe) UpdateSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	if subscriptionParams.Params.Metadata == nil {
		subscriptionParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		subscriptionParams.Params.Metadata[Env] = s.appService.Mode
	}
	return sub.Update(subscriptionID, subscriptionParams)
}

func (s *Stripe) CancelSubscription(subscriptionID string, subscriptionCancelParams *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	if subscriptionCancelParams.Params.Metadata == nil {
		subscriptionCancelParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		subscriptionCancelParams.Params.Metadata[Env] = s.appService.Mode
	}
	return sub.Cancel(subscriptionID, subscriptionCancelParams)
}

func (s *Stripe) CreateSetupIntent(setupIntentParams *stripe.SetupIntentParams) (*stripe.SetupIntent, error) {
	if setupIntentParams.Params.Metadata == nil {
		setupIntentParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		setupIntentParams.Params.Metadata[Env] = s.appService.Mode
	}
	return setupintent.New(setupIntentParams)
}

func (s *Stripe) CreateAccountLink(accountLinkParams *stripe.AccountLinkParams) (*stripe.AccountLink, error) {
	return accountlink.New(accountLinkParams)
}

func (s *Stripe) CreatePaymentIntentMultiparty(paymentIntentParams *stripe.PaymentIntentParams, linkedAccountID string) (*stripe.PaymentIntent, error) {
	if paymentIntentParams.Params.Metadata == nil {
		paymentIntentParams.Params.Metadata = map[string]string{Env: s.appService.Mode}
	} else {
		paymentIntentParams.Params.Metadata[Env] = s.appService.Mode
	}
	paymentIntentParams.SetStripeAccount(linkedAccountID)
	return paymentintent.New(paymentIntentParams)
}

func (s *Stripe) NewCheckoutSession(paymentMethods []string, mode, successURL, CancelURL string, lineItems []*stripe.CheckoutSessionLineItemParams, discounts []*stripe.CheckoutSessionDiscountParams) *stripe.CheckoutSession {
	params := &stripe.CheckoutSessionParams{
		Params:             stripe.Params{Metadata: map[string]string{Env: s.appService.Mode}},
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
	GetSubscription(subscriptionID string, subscriptionParams *stripe.SubscriptionParams) (*stripe.Subscription, error)
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
