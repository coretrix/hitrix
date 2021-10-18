package checkout

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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

type Instrument struct {
	ID            string `json:"id"`
	Type          string `json:"type"`
	Fingerprint   string `json:"fingerprint"`
	ExpiryMonth   int    `json:"expiry_month"`
	ExpiryYear    int    `json:"expiry_year"`
	Name          string `json:"name"`
	Scheme        string `json:"scheme"`
	Last4         string `json:"last4"`
	Bin           string `json:"bin"`
	CardType      string `json:"card_type"`
	CardCategory  string `json:"card_category"`
	Issuer        string `json:"issuer"`
	IssuerCountry string `json:"issuer_country"`
	ProductID     string `json:"product_id"`
	ProductType   string `json:"product_type"`
	AccountHolder struct {
		BillingAddress struct {
			AddressLine1 string `json:"address_line1"`
			AddressLine2 string `json:"address_line2"`
			City         string `json:"city"`
			State        string `json:"state"`
			Zip          string `json:"zip"`
			Country      string `json:"country"`
		} `json:"billing_address"`
		Phone struct {
			CountryCode string `json:"country_code"`
			Number      string `json:"number"`
		} `json:"phone"`
	} `json:"account_holder"`
}

type CustomerResponse struct {
	ID          string       `json:"id"`
	Instruments []Instrument `json:"instruments"`
}

func (c *Checkout) GetCustomerInstruments(customerId string) *CustomerResponse {
	config, err := checkout.Create(c.secretKey, nil)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("GET", *config.URI+"/customers/"+customerId, nil)
	req.Header.Set("Authorization", c.secretKey)
	resp, _ := client.Do(req)

	if resp.StatusCode != 200 {
		panic(fmt.Sprintf("error calling customer api checkout : %s", resp.Body))
	}

	decoder := json.NewDecoder(resp.Body)
	var jres CustomerResponse
	err = decoder.Decode(&jres)
	if err != nil {
		panic(fmt.Sprintf("failed parsing json from customer api checkout : %s", err.Error()))
	}
	return &jres
}

func (c *Checkout) DeleteCustomerInstrument(instrumentId string) bool {
	config, err := checkout.Create(c.secretKey, nil)
	if err != nil {
		panic("failed creating checkout client: " + err.Error())
	}
	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", *config.URI+"/instruments/"+instrumentId, nil)
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
