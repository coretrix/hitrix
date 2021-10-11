package checkout

import (
	"context"

	"github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/payments"
)

type Checkout struct {
	ctx            context.Context
	environment    string
	secretKey      string
	publicKey      *string
	webhookSecrets map[string]string
}

func NewCheckout(secretKey string, publicKey *string, environment string, webhookSecrets map[string]string) *Checkout {
	return &Checkout{
		ctx:            context.Background(),
		environment:    environment,
		secretKey:      secretKey,
		publicKey:      publicKey,
		webhookSecrets: webhookSecrets,
	}
}

func (c *Checkout) RequestPayment(source interface{}, amount uint64, currency string, reference string, customer *payments.Customer, metadata map[string]string) *payments.Response {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	idempotencyKey := checkout.NewIdempotencyKey()
	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}
	var client = payments.NewClient(*config)
	var request = &payments.Request{
		Amount:    amount,
		Source:    source,
		Currency:  currency,
		Reference: reference,
		Customer:  customer,
		Metadata:  metadata,
	}
	response, err := client.Request(request, &params)

	if err != nil {
		panic("checkout.com new payment request error: " + err.Error())
	}

	return response
}

func (c *Checkout) RequestRefunds(amount uint64, reference string, metadata map[string]string) *payments.RefundsResponse {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	idempotencyKey := checkout.NewIdempotencyKey()

	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}

	var client = payments.NewClient(*config)

	request := &payments.RefundsRequest{
		Amount:    amount,
		Reference: reference,
		Metadata:  metadata,
	}

	response, err := client.Refunds("pay_", request, &params)
	if err != nil {
		panic("checkout.com refund error: " + err.Error())
	}

	return response
}

func (c *Checkout) CheckWebhookKey(keyCode, key string) bool {
	if value, found := c.webhookSecrets[keyCode]; found {
		if value == key {
			return true
		}
	}
	return false
}
