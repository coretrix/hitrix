# OTP service
If you want to authenticate your user using OTP you may need to use OTP service.
This service can send the code using SMS or even call and verifying it later.

Register the service into your `main.go` file:
```go
registry.ServiceProviderOTP()
```

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