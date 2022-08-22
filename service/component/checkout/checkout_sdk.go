package checkout

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/checkout/checkout-sdk-go"
	"github.com/checkout/checkout-sdk-go/instruments"
	"github.com/checkout/checkout-sdk-go/payments"
	"github.com/checkout/checkout-sdk-go/tokens"
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

func (c *Checkout) RequestPayment(request *payments.Request) *payments.Response {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	idempotencyKey := checkout.NewIdempotencyKey()
	params := checkout.Params{
		IdempotencyKey: &idempotencyKey,
	}

	var client = payments.NewClient(*config)
	response, err := client.Request(request, &params)

	if err != nil {
		panic("checkout.com new payment request error: " + err.Error())
	}

	return response
}

func (c *Checkout) RequestRefunds(amount uint64, paymentID, reference string, metadata map[string]string) *payments.RefundsResponse {
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

	response, err := client.Refunds(paymentID, request, &params)
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

func (c *Checkout) DeleteCustomerInstrument(instrumentID string) bool {
	config, err := checkout.Create(c.secretKey, nil)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", *config.URI+"/instruments/"+instrumentID, nil)

	req.Header.Set("Authorization", c.secretKey)

	resp, _ := client.Do(req)

	if resp.StatusCode == 404 {
		return false
	}

	if resp.StatusCode != 204 {
		panic(fmt.Sprintf("error calling delete instrument api checkout : %s", resp.Body))
	}

	return true
}

func (c *Checkout) GetCustomer(idOrEmail string) (bool, *CustomerResponse) {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", *config.URI+"/customers/"+idOrEmail, nil)

	req.Header.Set("Authorization", c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(req)
	if resp.StatusCode == 404 {
		return false, nil
	} else if resp.StatusCode == 200 {
		res := &CustomerResponse{}

		err := json.NewDecoder(resp.Body).Decode(res)
		if err != nil {
			panic(err)
		}

		return true, res
	}

	data, _ := io.ReadAll(resp.Body)
	panic(fmt.Sprintf("wrong status checkout get customer code: %d, body %s", resp.StatusCode, string(data)))
}

func (c *Checkout) SaveGetClient(customerData *SaveCustomerRequest) (created bool, customer *CustomerResponse) {
	exists, customerRes := c.GetCustomer(customerData.Email)
	if exists {
		return false, customerRes
	}

	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	client := &http.Client{}
	jsonReq, _ := json.Marshal(customerData)
	req, _ := http.NewRequest("POST", *config.URI+"/customers/", bytes.NewBuffer(jsonReq))

	req.Header.Set("Authorization", c.secretKey)
	req.Header.Set("Content-Type", "application/json")

	resp, _ := client.Do(req)
	if resp.StatusCode == 201 {
		_, customerRes := c.GetCustomer(customerData.Email)

		return true, customerRes
	}

	data, _ := io.ReadAll(resp.Body)
	panic(fmt.Sprintf("wrong status checkout create customer code: %d, body %s", resp.StatusCode, string(data)))
}

func (c *Checkout) CreateToken(request *tokens.Request) (string, error) {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	token := tokens.NewClient(*config)

	resp, err := token.Request(request)
	if err != nil {
		return "", err
	}

	return resp.Created.Token, nil
}

func (c *Checkout) CreateInstrument(request *instruments.Request) (*instruments.Response, error) {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	client := instruments.NewClient(*config)

	res, err := client.Create(request)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Checkout) GetInstrument(sourceID string) (*instruments.Response, error) {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}

	client := instruments.NewClient(*config)

	res, err := client.Get(sourceID)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Checkout) GetPaymentDetail(paymentID string) (*payments.PaymentResponse, error) {
	config, err := checkout.Create(c.secretKey, c.publicKey)
	if err != nil {
		return nil, err
	}
	var client = payments.NewClient(*config)
	response, err := client.Get(paymentID)

	if err != nil {
		return nil, err
	}

	return response, nil
}
