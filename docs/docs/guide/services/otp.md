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
