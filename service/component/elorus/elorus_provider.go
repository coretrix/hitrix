package elorus

type IProvider interface {
	CreateContact(request *CreateContactRequest) (*ElorusResponse, error)
	CreateInvoice(request *CreateInvoiceRequest) (*ElorusResponse, error)
}
