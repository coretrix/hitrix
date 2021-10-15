# JWT
You can use that service to encode and decode JWT tokens

Register the service into your `main.go` file:
```go
registry.ServiceProviderJWT()
```

Access the service:
```go
service.DI().JWT()
```