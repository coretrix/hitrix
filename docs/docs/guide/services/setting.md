# Setting service
If your application requires configurations that might change or predefined, you need to use setting service. You should save your settings in `SettingsEntity`, then use this service to fetch it.


Register the service into your `main.go` file:
```go 
registry.ServiceProviderSetting()
```

Access the service:
```go
service.DI().Setting()
```

# Use case

Imagine you need to restrict access to login page after certain number of failed login attempts. You can simply store this value in `SettingsEntity` and fetch it using this service:
```go
package save 
import (
 "service"
)

func SaveConfig(){
    ormService := service.DI().ORMEngine()
    ormService.Flush(&entity.SettingsEntity{
        Key: "user.login.threshold",
        Value: "3",
    })
}
```

Then later in your login package, you can retrieve this value and use it:

```go
package login
import (
    "errors"
    "service"
)

func Login(currentCount uint64) error {
    ormService := service.DI().ORMEngine()
    allowed, found := service.DI().Setting().GetUint64(ormService, "user.login.threshold")
    if found && currentCount> allowed{
        return errors.New("too many login attempt")
    }  
    return nil  
}
```