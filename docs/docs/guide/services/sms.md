# SMS Service
This service is capable of sending simple message and otp message and also calling by different sms providers .
for now we support 3 sms providers : `twilio` `sinch` `kavenegar`

Register the service into your `main.go` file:
```go 
registry.ServiceProviderSMS()
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
    SetFromPrimaryGateway(string)
    SetFromSecondaryGateway(string)
    SetPrimaryGatewayError(string)
    SetSecondaryGatewayError(string)
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
	FromPrimaryGateway    string
	FromSecondaryGateway  string
	PrimaryGatewayError   string
	SecondaryGatewayError string
	Type                  string    `orm:"enum=entity.SMSTrackerTypeAll;required"`
	SentAt                time.Time `orm:"time"`
}
```
we have 2 providers active at the same time `primary` `secondary` and when send via primary fails we try to send with the secondary provider.
```go
func SendOTPSMS(*OTP) error{}
func SendOTPCallout(*OTP) error{}
func SendMessage(*Message) error{}
```
1. The `SendOTPSMS` send otp sms by providing the otp data
```go
type OTP struct {
	OTP      string
	Number   string
	CC       string
	Provider *Provider
	Template string
}
```
2. The `SendOTPCallout` used to call and tell the otp code
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
    authy_url: ENV[SMS_TWILIO_AUTHY_URL]
    authy_api_key: ENV[SMS_TWILIO_AUTHY_API_KEY]
    verify_url: ENV[SMS_TWILIO_VERIFY_URL]
    verify_sid: ENV[SMS_TWILIO_VERIFY_SID]
  kavenegar:
    api_key: ENV[SMS_KAVENEGAR_API_KEY]
    sender: ENV[SMS_KAVENEGAR_SENDER]
  sinch:
    app_id: ENV[SMS_SINCH_APP_ID]
    app_secret: ENV[SMS_SINCH_APP_SECRET]
    msg_url: ENV[SMS_SINCH_MSG_URL]
    from_number: ENV[SMS_SINCH_FROM_NUMBER]
    call_url: ENV[SMS_SINCH_CALL_URL]
    caller_number: ENV[SMS_SINCH_CALLER_NUMBER]
```
