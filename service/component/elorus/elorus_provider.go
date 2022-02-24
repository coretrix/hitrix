package elorus

type IProvider interface {
	CreateContact(request *CreateContactRequest) (*Response, error)
	CreateInvoice(request *CreateInvoiceRequest) (*Response, error)
}
