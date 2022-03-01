# Elorus.com API

This service is used to access Elorus platform, for creating and managing invoices

Register the service into your `main.go` file:

```go
hitrixRegistry.ServiceProviderElorus()
```

And you should put your credentials and other configs in `config/hitrix.yml`

```yml
elorus:
  url: https://api.elorus.com/v1.1
  token: secret
  organization_id: secret
```

Access the service:
```go
elorusService := service.DI().Elorus()
```


Using the service:
```go
// Request to create contact
elorusService.CreateContact(
    elorus.CreateContactRequest{
        FirstName: "name",
        Active:  true,
        Company:"company",
        VatNumber:"BG108"
        Email: []struct {   
        	Email   string `json:"email"`
            Primary bool   `json:"primary"`
        }{{
            Email:  "email@email.com",
            Primary: true,
        }},
        Phones: []struct {
            Number  string `json:"number"`
            Primary bool   `json:"primary"`
        }{{
            Number:  "0869586598",
            Primary: true,
        }},        
    },
)
      
// Request to create invoice
elorusService.CreateInvoice(
	elorus.CreateInvoiceRequest{
        Date:              time.Now().Format("2006-01-02"),
        Client:            contactId,
        ClientDisplayName: "name",
        ClientEmail:       "email@email.com",
        ClientVatNumber:   "BG108",
        Number:            "0",
        Items: []struct {
            Title                        string   `json:"title"`
            Description                  string   `json:"description"`
            Quantity                     string   `json:"quantity"`
            UnitMeasure                  int      `json:"unit_measure"`
            UnitValue                    string   `json:"unit_value"`
            Taxes                        []string `json:"Taxes"`
            UnitTotal                    string   `json:"unit_total"`
        }{{
            Title:  "title",
            Quantity: "5",
			UnitValue: "800",
            UnitMeasure: 1,
            UnitTotal: "500",
            Taxes:       []string{"2416104549958812800"},
        }},
    }
)     

// Request to get invoice list
elorusService.GetInvoiceList(
	elorus.GetInvoiceListRequest{
        Client:            contactId,
        Page: "1",
        PageSize:"10",
    }
)

// Request to get invoice list
elorusService.DownloadInvoice(
	elorus.DownloadInvoiceRequest{
        ID:            "id",
    }
)
```
