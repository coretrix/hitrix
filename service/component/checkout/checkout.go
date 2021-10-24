package checkout

import (
	"github.com/checkout/checkout-sdk-go/instruments"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
)

type ICheckout interface {
	CheckWebhookKey(keyCode, key string) bool
	RequestPayment(source interface{}, amount uint64, currency string, reference string, customer *payments.Customer, metadata map[string]string) *payments.Response
	RequestRefunds(amount uint64, reference string, metadata map[string]string) *payments.RefundsResponse
	DeleteCustomerInstrument(instrumentID string) bool
	GetCustomer(idOrEmail string) (bool, *CustomerResponse)
	SaveGetClient(customerData *SaveCustomerRequest) (created bool, customer *CustomerResponse)
	CreateToken(request *tokens.Request) (string, error)
	CreateInstrument(request *instruments.Request) (*instruments.Response, error)
}
