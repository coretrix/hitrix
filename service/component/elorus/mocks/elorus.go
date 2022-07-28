package mocks

import (
	"io"

	"github.com/stretchr/testify/mock"

	"github.com/coretrix/hitrix/service/component/elorus"
)

type FakeElorus struct {
	mock.Mock
}

func (e *FakeElorus) CreateContact(request *elorus.CreateContactRequest) (*elorus.Response, error) {
	args := e.Called(request)

	return args.Get(0).(*elorus.Response), args.Error(1)
}

func (e *FakeElorus) CreateInvoice(request *elorus.CreateInvoiceRequest) (*elorus.Response, error) {
	args := e.Called(request)

	return args.Get(0).(*elorus.Response), args.Error(1)
}

func (e *FakeElorus) GetInvoiceList(request *elorus.GetInvoiceListRequest) (*elorus.InvoiceListResponse, error) {
	args := e.Called(request)

	return args.Get(0).(*elorus.InvoiceListResponse), args.Error(1)
}

func (e *FakeElorus) DownloadInvoice(request *elorus.DownloadInvoiceRequest) (*io.ReadCloser, error) {
	args := e.Called(request)

	return args.Get(0).(*io.ReadCloser), args.Error(1)
}
