# Slack
Gives you ability to send slack messages using slack bot. Also it's used to send messages if you use our ErrorLogger service.
The config that needs to be set in hitrix.yaml is:

```yaml
slack:
    token: "your token"
    error_channel: "test" #optional, used by ErrorLogger
    dev_panel_url: "test" #optional, used by ErrorLogger

```

Register the service into your `main.go` file:
```go 
registry.ServiceProviderSlack()
```

Access the service:
```go
service.DI().Slack()
```