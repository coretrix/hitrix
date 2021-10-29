# Checkout.com API

This service is used to create a payment, check webhook key, request a refund and manage user cards for checkout.com payment API

Register the service into your `main.go` file:

```go
hitrixRegistry.ServiceProviderCheckout()
```

And you should put your credentials and other configs in `config/hitrix.yml`

```yml
checkout:
  secret_key: secret
  public_key: public
  currency: USD
  webhook_keys:
    main: somekey
```

Access the service:
```go
checkoutService := service.DI().Checkout()
```


Using the service:
```go
// Request a payment
checkoutService.RequestPayment(
    payments.IDSource{
        Type: "id",
        ID:  "sometoken",
    },
    100,
    "USD",
    "Order-1000",
    &payments.Customer{Email: "email@email.com"},
    map[string]string{"OrderId": "Order-1000"}
)
      
// Request a refund
checkoutService.RequestRefunds(1000, "PaymentId", "Order-1000", map[string]string{"OrderId": "Order-1000", "RefundsID": "Order-1000"})
      
// Validating incoming webhook
checkoutService.CheckWebhookKey("main", "value of Authorization header")

// Get user cards
checkoutService.GetCustomerInstruments("cus_someid")

// Delete user card
checkoutService.DeleteCustomerInstrument("src_someid")
```
