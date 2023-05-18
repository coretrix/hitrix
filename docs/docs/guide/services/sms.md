# SMS Service
This service is capable of sending sms messages using different providers providers .
We supports next providers `twilio` `sinch` `kavenegar`

Register the service into your `main.go` file:
```go 
registry.ServiceProviderSMS(sms.NewTwilioProvider, sms.NewSinchProvider)
```

Access the service:
```go
service.DI().SMS()
```

##### Dependencies :
`ClockService`

and also when registering the service you need to pass it the `LogEntity` that is responsible to log every action made by sms service :
```go
type LogEntity interface {
    beeorm.Entity
    SetStatus(string)
    SetTo(string)
    SetText(string)
    SetFromPrimaryProvider(string)
    SetFromSecondaryProvider(string)
    SetPrimaryProviderError(string)
    SetSecondaryProviderError(string)
    SetType(string)
    SetSentAt(time time.Time)
}
```
for example :
```go
const (
	SMSTrackerTypeSMS     = "sms"
	SMSTrackerTypeCallout = "callout"
)

type smsTrackerTypeAll struct {
	SMSTrackerTypeSMS     string
	SMSTrackerTypeCallout string
}

var SMSTrackerTypeAll = smsTrackerTypeAll{
	SMSTrackerTypeSMS:     SMSTrackerTypeSMS,
	SMSTrackerTypeCallout: SMSTrackerTypeCallout,
}

type SmsTrackerEntity struct {
	beeorm.ORM               `orm:"table=sms_tracker"`
	ID                    uint64
	Status                string
	To                    string `orm:"varchar=15"`
	Text                  string
	FromPrimaryProvider    string
	FromSecondaryProvider  string
	PrimaryProviderError   string
	SecondaryProviderError string
	Type                  string    `orm:"enum=entity.SMSTrackerTypeAll;required"`
	SentAt                time.Time `orm:"time"`
}
```
For every separate message we can provide 2 providers - `primary` `secondary` and if primary fails we try to send with the secondary provider.

3. The `SendMessage` used to send simple message
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
```
