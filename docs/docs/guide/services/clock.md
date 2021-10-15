# Clock service
This service is used for `time` operations. It is better to use it everywhere instead of `time.Now()` because it can be mocked and you can set whatever time you want in your tests

Register the service into your `main.go` file:
```go 
registry.ServiceClock(),
```

Access the service:
```go
service.DI().Clock()
```

The methods that this service provide are:
```Now() and NowPointer()```
