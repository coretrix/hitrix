package mocks

import (
  "io"
	"time"

	"github.com/stretchr/testify/mock"

  "github.com/coretrix/hitrix/service/component/elorus"
)

type FakeElorus struct {
	mock.Mock
}

func (c *FakeElorus) CreateContact(request *elorus.CreateContactRequest) (*elorus.Response, error) {
  args := t.Called(request)
	return args.Get(0).(*elorus.Response), args.Error(1)
}

func (c *FakeElorus) CreateInvoice(request *elorus.CreateInvoiceRequest) (*elorus.Response, error) {
  args := t.Called(request)
	return args.Get(0).(*elorus.Response), args.Error(1)
}

func (c *FakeElorus) GetInvoiceList(request *elorus.GetInvoiceListRequest) (*elorus.InvoiceListResponse, error) {
  args := t.Called(request)
	return args.Get(0).(*elorus.InvoiceListResponse), args.Error(1)
}

func (c *FakeElorus) DownloadInvoice(request *elorus.DownloadInvoiceRequest) (*io.ReadCloser, error) {
  args := t.Called(request)
	return args.Get(0).(*io.ReadCloser), args.Error(1)
}
