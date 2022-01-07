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

You can register OTP providers in 2 ways:
1. Provide setting in DB for key: `otp_sms_provider` with value either of `Twilio` or `Sinch`. You can pass both providers as well separated by comma - `Twilio,Sinch`.
2. Call `registry.ServiceProviderOTP(otp.SMSOTPProviderTwilio, otp.SMSOTPProviderSinch)` with 1 or more parameters for force provider.

Access the service:
```go
service.DI().OTP()
```
# Use case
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
