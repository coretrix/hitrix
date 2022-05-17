package elorus

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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
	IsClient  bool   `json:"is_client"`
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
	Number            string `json:"number"`
	DueDays           int    `json:"due_days"`
	Items             []struct {
		Product     string   `json:"product"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Quantity    string   `json:"quantity"`
		UnitMeasure int      `json:"unit_measure"`
		UnitValue   string   `json:"unit_value"`
		Taxes       []string `json:"Taxes"`
		UnitTotal   string   `json:"unit_total"`
	} `json:"items"`
}

type GetInvoiceListRequest struct {
	Client   string `json:"client"`
	Status   string `json:"status"`
	Page     string `json:"page"`
	PageSize string `json:"page_size"`
}

type DownloadInvoiceRequest struct {
	ID string `json:"id"`
}

type InvoiceListResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		ID      string `json:"id"`
		Status  string `json:"status"`
		Date    string `json:"date"`
		DueDate string `json:"due_date"`
		Total   string `json:"total"`
	} `json:"results"`
}

type Response struct {
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

func (e *Elorus) CreateContact(request *CreateContactRequest) (*Response, error) {
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

	if resp.StatusCode > 201 {
		var failedResponse interface{}
		err = json.NewDecoder(resp.Body).Decode(&failedResponse)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("not successful request with status code : %v , response : %v", resp.StatusCode, failedResponse)
	}

	response := new(Response)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (e *Elorus) CreateInvoice(request *CreateInvoiceRequest) (*Response, error) {
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

	if resp.StatusCode > 201 {
		var failedResponse interface{}
		err = json.NewDecoder(resp.Body).Decode(&failedResponse)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("not successful request with status code : %v , response : %v", resp.StatusCode, failedResponse)
	}

	response := new(Response)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (e *Elorus) GetInvoiceList(request *GetInvoiceListRequest) (*InvoiceListResponse, error) {
	client := &http.Client{}

	requestURL, err := url.Parse(e.url + "/invoices/")
	if err != nil {
		return nil, err
	}
	query := requestURL.Query()
	query.Set("client", request.Client)
	query.Set("page", request.Page)
	query.Set("page_size", request.PageSize)
	if len(request.Status) > 0 {
		query.Set("status", request.Status)
	}
	requestURL.RawQuery = query.Encode()

	req, err := http.NewRequest("GET", requestURL.String(), nil)
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

	if resp.StatusCode > 201 {
		var failedResponse interface{}
		err = json.NewDecoder(resp.Body).Decode(&failedResponse)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("not successful request with status code : %v , response : %v", resp.StatusCode, failedResponse)
	}

	response := new(InvoiceListResponse)
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (e *Elorus) DownloadInvoice(request *DownloadInvoiceRequest) (*io.ReadCloser, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", fmt.Sprintf(e.url+"/invoices/%s/pdf", request.ID), nil)
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
	return &resp.Body, nil
}
