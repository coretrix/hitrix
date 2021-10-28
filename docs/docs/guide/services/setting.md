#Setting service
If your application requires configurations that might change or predefined, you need to use setting service. You should save your settings in `SettingsEntity`, then use this service to fetch it.


Register the service into your `main.go` file:
```go 
registry.ServiceProviderSetting()
```

Access the service:
```go
service.DI().Setting()
```