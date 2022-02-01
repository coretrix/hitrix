# Gql
This service can be used to return GQL errors and translate them if localize service is registered

Register the service into your `main.go` file:
```go 
registry.ServiceProviderGql()
```

Access the service:
```go
service.DI().Gql()
```

The methods that can be used are `GraphqlErr` and `GraphqlErrPath`