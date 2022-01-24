# Slack
Gives you ability to send Slack messages using slack bot. You can define multiple Slack bots. 
Also, it's used to send messages if you use our `ErrorLogger` service using `errors` bot.
The config that needs to be set in hitrix.yaml is:

```yaml
slack:
    error_channel: "test" #optional, used by ErrorLogger
    dev_panel_url: "test" #optional, used by ErrorLogger
    bot_tokens:
      errors: "your token"
      another_bot: "second token"
```

> NOTE: `bot_tokens.errors` is a must-have when using `ErrorLogger` service.

Register the service into your `main.go` file:
```go 
registry.ServiceProviderSlack()
```

Access the service:
```go
service.DI().Slack()
```