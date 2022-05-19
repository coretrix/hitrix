# OTP service
If you want to authenticate your user using OTP you may need to use OTP service.
This service can send the code using SMS or even call and verifying it later.

Register the service into your `main.go` file:
```go
registry.ServiceProviderOTP()
```

Supported OTP providers:
1. Twilio
2. Sinch
3. Mada

Now it is possible to provide phone number prefixes for each OTP provider.

For example if I provide the following setting `Twilio;Mada:+35987,+35988`, this means if phone numer
starts with either of `+35987` or `+35988`, then we will use `Mada`, for all others numbers we will use `Twilio`.

You can register OTP providers in 2 ways:
1. Provide setting in DB for key: `otp_sms_provider` with value either of `Twilio` or `Sinch` or `Mada`. You can pass all providers as well separated by semicolon - `Twilio;Mada:+35987,+35988;Sinch`.
2. Call `registry.ServiceProviderOTP(otp.SMSOTPProviderTwilio, otp.SMSOTPProviderSinch)` with 1 or more parameters for force provider.

## Retry feature:
You can set up the service to retry failed OTP send attempts.
In order to do this, you need to add in the config:
```
sms:
  retry: true
  max_retries: 20
``` 

For retry feature you also need to start in your app this consumer:

```go
    // add this if you want to use send OTP retry feature
    s.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
	    go b.RunScript(&scripts.RetryOTPConsumer{})
    })
```
Retry feature uses exponential backoff to retry OTP requests, starting from 0.5 seconds.
If `max_retries` is reached, the consumer will drop the OTP request and mark it unsendable in DB.

Access the service:
```go
service.DI().OTP()
```
## Use case
You can send OTP to user phone using SMS  or call like this:
```go
package auth
import (
    "context"
    "service"
    "github.com/coretrix/hitrix/service/component/otp"
)

func SendOTP(){
    ormService := service.DI().OrmEngineForContext(context.Background())
    OTPService := service.DI().OTP()

    // add this if you want to use send OTP retry feature
    s.RunBackgroundProcess(func(b *hitrix.BackgroundProcessor) {
    	go b.RunScript(&scripts.RetryOTPConsumer{})
    })

	//SMS
    code, err := OTPService.SendSMS(ormService, &otp.Phone{
        Number: "+123456789",
    })

    //call
    code, err := OTPService.Call(ormService, &otp.Phone{
        Number: "+123456789",
    })
}
```

Then you can verify OTP like this:
```go
package auth
import (
    "context"
    "service"
    "github.com/coretrix/hitrix/service/component/otp"
)

func Verify(){
    ormService := service.DI().OrmEngineForContext(context.Background())
    OTPService := service.DI().OTP()
    code:="1234" //the code user entered

    otpRequestValid, otpCodeValid, err := OTPService.Verify(
    ormService,
    &otp.Phone{Number: "+123456789"},
    code,
    )
}
```
