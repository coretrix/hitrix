# ORM Engine
Used to access ORM in background scripts. It is one instance for the whole script

Register the service into your `main.go` file:
```go
registry.ServiceProviderOrmEngine()
```

Access the service:
```go
service.DI().ORMEngine()
```

Never use that service in API. It is not thread safe!