# Authentication Service
This service is used to making the life easy by doing the whole authentication life cycle using JWT token.

Register the service into your `main.go` file:
```go
registry.ServiceProviderAuthentication(),
```

Access the service:
```go
service.DI().Authentication()
```

##### Dependencies :
`JWTService`

`PasswordService`

`ClockService`

`GeneratorService`

`GoogleService` # optional , when you need to support google login

`FacebookService` # optional , when you need to support facebook login

`AppleService` # optional , when you need to support apple login

```go
func Authenticate(ormService *datalayer.DataLayer, uniqueValue string, password string, entity AuthProviderEntity) (accessToken string, refreshToken string, err error) {}
func VerifyAccessToken(ormService *datalayer.DataLayer, accessToken string, entity beeorm.Entity) error {}
func VerifySocialLogin(ctx context.Context, source, token string, isAndroid bool)
func RefreshToken(ormService *datalayer.DataLayer, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {}
func LogoutCurrentSession(ormService *datalayer.DataLayer, accessKey string){}
func LogoutAllSessions(ormService *datalayer.DataLayer, id uint64)
func AuthenticateOTP(ormService *datalayer.DataLayer, phone string, entity OTPProviderEntity) (accessToken string, refreshToken string, err error){}
```
1. The `Authenticate` function will take an uniqueValue such as Email or Mobile, a plain password, and generates accessToken and refreshToken.
   You will also need to pass your entity as third argument, and it will give you the specific user entity related to provided access token
   The entity should implement the `AuthProviderEntity` interface :
    ```go
       type AuthProviderEntity interface {
        beeorm.Entity
        GetUniqueFieldName() string
        GetPassword() string
       }
    ```
   The example of such entity is as follows:
    ```go
    type UserEntity struct {
        beeorm.ORM  `orm:"table=users;redisCache;redisSearch=search_pool"`
        ID       uint64 `orm:"searchable;sortable"`
        Email    string `orm:"required;unique=Email;searchable"`
        Password string `orm:"required"`
    }
   
    func (user *UserEntity) GetUniqueFieldName() string {
        return "Email"
    }
    
    func (user *UserEntity) GetPassword() string {
    return user.Password
    }
    ```
2. The `VerifyAccessToken` will get the AccessToken, process the validation and expiration, and fill the entity param with the authenticated user entity in case of successful authentication.
3. The `RefreshToken` method will generate a new token pair for given user
4. The `LogoutCurrentSession` you can logout the user current session , you need to pass it the `accessKey`  that is the jwt identifier `jti` the exists in both access and refresh token.
5. The `LogoutAllSessions` you can logout the user from all sessions , you need to pass it the `id` (user id).
6. You need to have a `authentication` key in your config file for this service to work. `secret` key under `authentication` is mandatory but other options are optional:
7. The service can also support `OTP` if you want your service to support otp you should have `support_otp` key set to true under `authentication`
8. The service also needs redis to store its sessions so you need to identify the redis storage name in config , the key is `auth_redis` under `authentication`
```yaml
authentication:
  secret: "a-deep-dark-secret" #mandatory, secret to be used for JWT
  access_token_ttl: 86400 # optional, in seconds, default to 1day
  refresh_token_ttl: 31536000 #optional, in seconds, default to 1year
  auth_redis: default #optional , default is the default redis
  otp_ttl: 120 #optional ,set it when you want to use otp, It is the ttl of otp code , default is 60 seconds
  otp_length: 5 #optional, set if you want to customize the length of otp (i.e. Email OTP)
```
