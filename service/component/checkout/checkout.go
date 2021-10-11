package checkout

import (
	"github.com/checkout/checkout-sdk-go/payments"
)

type ICheckout interface {
	CheckWebhookKey(keyCode, key string) bool
	RequestPayment(source interface{}, amount uint64, currency string, reference string, customer *payments.Customer, metadata map[string]string) *payments.Response
	RequestRefunds(amount uint64, reference string, metadata map[string]string) *payments.RefundsResponse
}
