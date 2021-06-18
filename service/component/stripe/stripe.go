package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/checkout/session"
	"github.com/stripe/stripe-go/v72/webhook"
)

type Stripe struct {
	ctx            context.Context
	environment    string
	webhookSecrets map[string]string
	stripeConfig   map[string]interface{}
}

func NewStripe(token string, webhookSecrets map[string]string, environment string, config map[string]interface{}) *Stripe {
	stripe.Key = token

	return &Stripe{
		ctx:            context.Background(),
		environment:    environment,
		stripeConfig:   config,
		webhookSecrets: webhookSecrets,
	}
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
	ConstructWebhookEvent(reqBody []byte, signature string, webhookKey string) (stripe.Event, error)
	NewCheckoutSession(paymentMethods []string, mode, successURL, CancelURL string, lineItems []*stripe.CheckoutSessionLineItemParams, discounts []*stripe.CheckoutSessionDiscountParams) *stripe.CheckoutSession
}
