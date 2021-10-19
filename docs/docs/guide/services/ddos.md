# DDOS Protection
This service contains DDOS protection features

Register the service into your `main.go` file:
```go
registry.ServiceProviderDDOS()
```

Access the service:
```go
service.DI().DDOS()
```

You can protect for example login endpoint from many attempts  by using method `ProtectManyAttempts`
