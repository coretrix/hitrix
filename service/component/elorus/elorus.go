package elorus

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

type Elorus struct {
	ctx            context.Context
	url            string
	token          string
	organizationID string
	environment    string
}

type CreateContactRequest struct {
	FirstName string `json:"first_name"`
	Active    bool   `json:"active"`
	Company   string `json:"company"`
	VatNumber string `json:"vat_number"`
	Email     []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	} `json:"email"`
	Phones []struct {
		Number  string `json:"number"`
		Primary bool   `json:"primary"`
	} `json:"phones"`
}

type CreateInvoiceRequest struct {
	Date              string `json:"date"`
	Client            string `json:"client"`
	ClientDisplayName string `json:"client_display_name"`
	ClientVatNumber   string `json:"client_vat_number"`
	ClientEmail       string `json:"client_email"`
	Items             []struct {
		Product                      string   `json:"product"`
		Title                        string   `json:"title"`
		Description                  string   `json:"description"`
		Quantity                     string   `json:"quantity"`
		UnitMeasure                  int      `json:"unit_measure"`
		UnitValue                    string   `json:"unit_value"`
		Taxes                        []string `json:"Taxes"`
		UnitTotal                    string   `json:"unit_total"`
		MydataClassificationCategory string   `json:"mydata_classification_category"`
		MydataClassificationType     string   `json:"mydata_classification_type"`
	} `json:"items"`
}

type ElorusResponse struct {
	ID string `json:"id"`
}

func NewElorus(url string, token string, organizationID string, environment string) *Elorus {
	return &Elorus{
		ctx:            context.Background(),
		url:            url,
		token:          token,
		organizationID: organizationID,
		environment:    environment,
	}
}

func (e *Elorus) CreateContact(request *CreateContactRequest) (*ElorusResponse, error) {
	client := &http.Client{}

	jsonReq, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", e.url+"/contacts/", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+e.token)
	req.Header.Set("X-Elorus-Organization", e.organizationID)
	if e.environment != "prod" {
		req.Header.Set("X-Elorus-Demo", "true")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	response := new(ElorusResponse)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (e *Elorus) CreateInvoice(request *CreateInvoiceRequest) (*ElorusResponse, error) {
	client := &http.Client{}

	jsonReq, _ := json.Marshal(request)
	req, err := http.NewRequest("POST", e.url+"/invoices/", bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+e.token)
	req.Header.Set("X-Elorus-Organization", e.organizationID)
	if e.environment != "prod" {
		req.Header.Set("X-Elorus-Demo", "true")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	response := new(ElorusResponse)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
