# SMS Service
This service is capable of sending sms messages using different providers .
We supports next providers `twilio` `sinch` `kavenegar` `link mobility` `mobica`

Register the service into your `main.go` file:
```go 
registry.ServiceProviderSMS(sms.NewTwilioProvider, sms.NewSinchProvider)
```

Access the service:
```go
service.DI().SMS()
```

##### Dependencies :
`ClockService`, `ORMConfigService`

Every request is saved in SmsTrackerEntity. So you will be able to track every sms

We support defining `primary` and `secondary` providers. If primary fails we try to send a message with the secondary provider.

We can set default providers and providers per message as well

Default providers are set at `ServiceProviderSMS()`
Providers per message are set at `SendMessage()`

The method `SendMessage` used to send simple message
```go
type Message struct {
	Text     string
	Number   string
	Provider *Provider
}
```
##### configs
```yaml
sms:
  twilio:
    sid: ENV[SMS_TWILIO_SID]
    token: ENV[SMS_TWILIO_TOKEN]
    from_number: ENV[SMS_TWILIO_FROM_NUMBER]
  kavenegar:
    api_key: ENV[SMS_KAVENEGAR_API_KEY]
    sender: ENV[SMS_KAVENEGAR_SENDER]
  sinch:
    app_id: ENV[SMS_SINCH_APP_ID]
    app_secret: ENV[SMS_SINCH_APP_SECRET]
    msg_url: ENV[SMS_SINCH_MSG_URL]
    from_number: ENV[SMS_SINCH_FROM_NUMBER]
  mobica:
    email: ENV[SMS_MOBICA_EMAIL]
    password: ENV[SMS_MOBICA_PASSWORD]
    route: ENV[SMS_MOBICA_ROUTE]
    from: ENV[SMS_MOBICA_FROM]
    endpoint: ENV[SMS_MOBICA_ENDPOINT]
  link_mobility:
    service: ENV[SMS_LINK_MOBILITY_SERVICE]
    key: ENV[SMS_LINK_MOBILITY_KEY]
    secret: ENV[SMS_LINK_MOBILITY_SECRET]
    endpoint: ENV[SMS_LINK_MOBILITY_ENDPOINT]
    shortcode: ENV[SMS_LINK_MOBILITY_SHORTCODE]
```
