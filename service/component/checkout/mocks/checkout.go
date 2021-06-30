package mocks

import (
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/stretchr/testify/mock"
)

type FakeCheckoutClient struct {
	mock.Mock
}

func (c *FakeCheckoutClient) CheckWebhookKey(keyCode, key string) bool {
	return c.Called(keyCode, key).Get(0).(bool)
}

func (c *FakeCheckoutClient) RequestPayment(source interface{}, amount uint64, currency string, reference string, customer *payments.Customer, metadata map[string]string) *payments.Response {
	return c.Called(source, amount, currency, reference, customer, metadata).Get(0).(*payments.Response)
}
