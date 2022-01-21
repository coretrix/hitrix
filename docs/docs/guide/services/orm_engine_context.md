# ORM Engine Context
Used to access ORM in foreground scripts like API. It is one instance per every request

Register the service into your `main.go` file as context service:
```go
registry.ServiceProviderOrmEngineForContext()
```

Access the service:
```go
service.DI().ORMEngineForContext()
```