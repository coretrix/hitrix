# Template
This service can be used to render templates from your html files and also from mandrill templates

Register the service into your `main.go` file:
```go 
registry.ServiceProviderTemplate()
```

Access the service:
```go
service.DI().Template()
```