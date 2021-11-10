package mocks

import (
	"github.com/checkout/checkout-sdk-go/instruments"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
	"github.com/coretrix/hitrix/service/component/checkout"
	"github.com/stretchr/testify/mock"
)

type FakeCheckoutClient struct {
	mock.Mock
}

func (c *FakeCheckoutClient) CheckWebhookKey(keyCode, key string) bool {
	return c.Called(keyCode, key).Get(0).(bool)
}

func (c *FakeCheckoutClient) RequestPayment(request *payments.Request) *payments.Response {
	return c.Called(request).Get(0).(*payments.Response)
}

func (c *FakeCheckoutClient) RequestRefunds(amount uint64, paymentID, reference string, metadata map[string]string) *payments.RefundsResponse {
	return c.Called(amount, paymentID, reference, metadata).Get(0).(*payments.RefundsResponse)
}

func (c *FakeCheckoutClient) DeleteCustomerInstrument(instrumentID string) bool {
	return c.Called(instrumentID).Get(0).(bool)
}

func (c *FakeCheckoutClient) GetCustomer(idOrEmail string) (bool, *checkout.CustomerResponse) {
	return c.Called(idOrEmail).Get(0).(bool), c.Called(idOrEmail).Get(1).(*checkout.CustomerResponse)
}

func (c *FakeCheckoutClient) SaveGetClient(customerData *checkout.SaveCustomerRequest) (created bool, customer *checkout.CustomerResponse) {
	return c.Called(customerData).Get(0).(bool), c.Called(customerData).Get(1).(*checkout.CustomerResponse)
}

func (c *FakeCheckoutClient) CreateToken(request *tokens.Request) (string, error) {
	return c.Called(request).Get(0).(string), c.Called(request).Error(1)
}

func (c *FakeCheckoutClient) CreateInstrument(request *instruments.Request) (*instruments.Response, error) {
	return c.Called(request).Get(0).(*instruments.Response), c.Called(request).Error(1)
}

func (c *FakeCheckoutClient) GetInstrument(sourceID string) (*instruments.Response, error) {
	return c.Called(sourceID).Get(0).(*instruments.Response), c.Called(sourceID).Error(1)
}

func (c *FakeCheckoutClient) GetPaymentDetail(paymentID string) (*payments.PaymentResponse, error) {
	return c.Called(paymentID).Get(0).(*payments.PaymentResponse), c.Called(paymentID).Error(1)
}
