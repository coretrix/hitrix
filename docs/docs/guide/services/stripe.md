# Stripe

Stripe payment integration

Register the service into your `main.go` file:
```go
registry.ServiceProviderStripe(),
```

Access the service:
```go
service.DI().Stripe()
```

Config sample:

```yml
stripe:
  key: "api_key"
  webhook_secrets: # map of your webhook secrets
    checkout: "key"
```
