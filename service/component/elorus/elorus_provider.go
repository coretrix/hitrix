package elorus

import "io"

type IProvider interface {
	CreateContact(request *CreateContactRequest) (*Response, error)
	CreateInvoice(request *CreateInvoiceRequest) (*Response, error)
	GetInvoiceList(request *GetInvoiceListRequest) (*InvoiceListResponse, error)
	DownloadInvoice(request *DownloadInvoiceRequest) (*io.ReadCloser, error)
}
