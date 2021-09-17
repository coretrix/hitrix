package authentication

import (
	"errors"
	"fmt"
	mail2 "net/mail"

	"github.com/coretrix/hitrix/service/component/uuid"

	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"

	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/social"

	"github.com/coretrix/hitrix/service/component/clock"

	"github.com/coretrix/hitrix/service/component/generator"

	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix/service/component/sms"

	"github.com/dongri/phonenumber"

	"github.com/go-redis/redis/v8"

	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/latolukasz/beeorm"
)

const (
	separator                = ":"
	accessListSeparator      = ";"
	accessKeyPrefix          = "ACCESS"
	userAccessListPrefix     = "USER_KEYS"
	maxUserAccessKeysAllowed = 10

	SocialLoginGoogle   = "google"
	SocialLoginFacebook = "facebook"
)

type OTPProviderEntity interface {
	beeorm.Entity
	GetPhoneFieldName() string
	GetEmailFieldName() string
}

type AuthProviderEntity interface {
	beeorm.Entity
	GetUniqueFieldName() string
	GetPassword() string
}

type Authentication struct {
	accessTokenTTL       int
	refreshTokenTTL      int
	otpTTL               int
	passwordService      password.Password
	errorLoggerService   errorlogger.ErrorLogger
	jwtService           *jwt.JWT
	smsService           sms.ISender
	mailService          *mail.Sender
	socialServiceMapping map[string]social.IUserData
	generatorService     generator.IGenerator
	clockService         clock.IClock
	uuidService          uuid.IUUID
	secret               string
}

func NewAuthenticationService(
	secret string,
	accessTokenTTL int,
	refreshTokenTTL int,
	otpTTL int,
	smsService sms.ISender,
	generatorService generator.IGenerator,
	errorLoggerService errorlogger.ErrorLogger,
	clockService clock.IClock,
	passwordService password.Password,
	jwtService *jwt.JWT,
	mailService *mail.Sender,
	socialServiceMapping map[string]social.IUserData,
	uuidService uuid.IUUID,
) *Authentication {
	return &Authentication{
		secret:               secret,
		accessTokenTTL:       accessTokenTTL,
		refreshTokenTTL:      refreshTokenTTL,
		otpTTL:               otpTTL,
		passwordService:      passwordService,
		errorLoggerService:   errorLoggerService,
		jwtService:           jwtService,
		smsService:           smsService,
		clockService:         clockService,
		generatorService:     generatorService,
		mailService:          mailService,
		socialServiceMapping: socialServiceMapping,
		uuidService:          uuidService,
	}
}

type GenerateOTP struct {
	Mobile         string
	ExpirationTime string
	Token          string
}

type GenerateOTPEmail struct {
	Email          string
	ExpirationTime string
	Token          string
}

func (t *Authentication) GenerateAndSendOTP(ormService *beeorm.Engine, mobile string, countryCodeAlpha2 string) (*GenerateOTP, error) {
	// validate mobile number
	if len(countryCodeAlpha2) != 2 {
		return nil, errors.New("use alpha2 code for country")
	}

	number := phonenumber.Parse(mobile, countryCodeAlpha2)
	if number == "" {
		return nil, errors.New("phone number not valid")
	}

	code := t.generatorService.GenerateRandomRangeNumber(10000, 99999)

	err := t.smsService.SendOTPSMS(ormService, t.errorLoggerService, &sms.OTP{
		OTP: fmt.Sprint(code),
		Phone: &sms.Phone{
			Number:  number,
			ISO3166: phonenumber.GetISO3166ByNumber(number, false),
		},
		Provider: factorySMSProviders(countryCodeAlpha2),
		// TODO : replace with the desired message or get as a argument
		Template: "ICE-STORM",
	})
	if err != nil {
		return nil, err
	}

	expirationTime := t.clockService.Now().Add(time.Duration(t.otpTTL) * time.Second).Unix()
	token := t.generatorService.GenerateSha256Hash(fmt.Sprint(strconv.FormatInt(expirationTime, 10), number, fmt.Sprint(code)))

	return &GenerateOTP{
		Mobile:         number,
		ExpirationTime: strconv.FormatInt(expirationTime, 10),
		Token:          token,
	}, nil
}

func (t *Authentication) GenerateAndSendOTPEmail(ormService *beeorm.Engine, email string, template string, from string, title string) (*GenerateOTPEmail, error) {
	_, err := mail2.ParseAddress(email)

	if err != nil {
		return nil, errors.New("mail address not valid")
	}

	code := t.generatorService.GenerateRandomRangeNumber(10000, 99999)
	if t.mailService == nil {
		panic("mail service is not registered")
	}
	mailService := *t.mailService

	err = mailService.SendTemplateAsync(ormService, &mail.Message{
		From:         from,
		To:           email,
		Subject:      title,
		TemplateName: template,
		TemplateData: map[string]interface{}{"code": strconv.FormatInt(code, 10)},
	})

	if err != nil {
		return nil, err
	}

	expirationTime := t.clockService.Now().Add(time.Duration(t.otpTTL) * time.Second).Unix()
	token := t.generatorService.GenerateSha256Hash(fmt.Sprint(strconv.FormatInt(expirationTime, 10), email, fmt.Sprint(code)))

	return &GenerateOTPEmail{
		Email:          email,
		ExpirationTime: strconv.FormatInt(expirationTime, 10),
		Token:          token,
	}, nil
}

func (t *Authentication) VerifyOTP(code string, input *GenerateOTP) error {
	token := t.generatorService.GenerateSha256Hash(fmt.Sprint(input.ExpirationTime, input.Mobile, code))
	if token != input.Token {
		return errors.New("wrong code provided")
	}

	timeInt, err := strconv.ParseInt(input.ExpirationTime, 10, 64)
	if err != nil {
		panic("wrong time format")
	}
	expirationTime := time.Unix(timeInt, 0)

	if expirationTime.Before(t.clockService.Now()) {
		return errors.New("code expired")
	}

	return nil
}

func (t *Authentication) VerifyOTPEmail(code string, input *GenerateOTPEmail) error {
	token := t.generatorService.GenerateSha256Hash(fmt.Sprint(input.ExpirationTime, input.Email, code))
	if token != input.Token {
		return errors.New("wrong code provided")
	}

	timeInt, err := strconv.ParseInt(input.ExpirationTime, 10, 64)
	if err != nil {
		panic("wrong time format")
	}
	expirationTime := time.Unix(timeInt, 0)

	if expirationTime.Before(t.clockService.Now()) {
		return errors.New("code expired")
	}

	return nil
}

func (t *Authentication) VerifySocialLogin(source, token string) (*social.UserData, error) {
	socialProvider, ok := t.socialServiceMapping[source]
	if !ok {
		return nil, errors.New("not supported social provider: " + source)
	}
	return socialProvider.GetUserData(token)
}

func (t *Authentication) AuthenticateOTP(ormService *beeorm.Engine, phone string, entity OTPProviderEntity) (accessToken string, refreshToken string, err error) {
	q := &beeorm.RedisSearchQuery{}
	q.FilterString(entity.GetPhoneFieldName(), phone)
	found := ormService.RedisSearchOne(entity, q)
	if !found {
		return "", "", errors.New("invalid credentials")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) AuthenticateOTPEmail(ormService *beeorm.Engine, email string, entity OTPProviderEntity) (accessToken string, refreshToken string, err error) {
	q := &beeorm.RedisSearchQuery{}
	q.FilterString(entity.GetEmailFieldName(), email)
	found := ormService.RedisSearchOne(entity, q)
	if !found {
		return "", "", errors.New("invalid credentials")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) Authenticate(ormService *beeorm.Engine, uniqueValue string, password string, entity AuthProviderEntity) (accessToken string, refreshToken string, err error) {
	q := &beeorm.RedisSearchQuery{}
	q.FilterString(entity.GetUniqueFieldName(), uniqueValue)
	found := ormService.RedisSearchOne(entity, q)
	if !found {
		return "", "", errors.New("invalid user/pass")
	}

	if !t.passwordService.VerifyPassword(password, entity.GetPassword()) {
		return "", "", errors.New("invalid user/pass")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) AuthenticateByID(ormService *beeorm.Engine, id uint64, entity AuthProviderEntity) (accessToken string, refreshToken string, err error) {
	exists := ormService.LoadByID(id, entity)
	if !exists {
		return "", "", errors.New("id_does_not_exists")
	}
	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) generateUserTokens(ormService *beeorm.Engine, ID uint64) (accessToken string, refreshToken string, err error) {
	accessKey := t.generateAndStoreAccessKey(ormService, ID, t.refreshTokenTTL)

	accessToken, err = t.GenerateTokenPair(ID, accessKey, t.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = t.GenerateTokenPair(ID, accessKey, t.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	t.addUserAccessKeyList(ormService, ID, accessKey, "", t.refreshTokenTTL)
	return accessToken, refreshToken, nil
}

func (t *Authentication) VerifyAccessToken(ormService *beeorm.Engine, accessToken string, entity beeorm.Entity) (map[string]string, error) {
	payload, err := t.jwtService.VerifyJWTAndGetPayload(t.secret, accessToken, t.clockService.Now().Unix())
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(payload["sub"], 10, 64)
	if err != nil {
		return nil, err
	}

	accessKey := payload["jti"]

	_, has := ormService.GetRedis().Get(accessKey)
	if !has {
		return nil, errors.New("access key not found")
	}

	found := ormService.LoadByID(id, entity)
	if !found {
		return nil, errors.New("user_not_found")
	}

	return payload, nil
}

func (t *Authentication) RefreshToken(ormService *beeorm.Engine, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	payload, err := t.jwtService.VerifyJWTAndGetPayload(t.secret, refreshToken, t.clockService.Now().Unix())
	if err != nil {
		return "", "", err
	}

	id, err := strconv.ParseUint(payload["sub"], 10, 64)
	if err != nil {
		return "", "", err
	}

	//check the access key
	oldAccessKey := payload["jti"]
	_, has := ormService.GetRedis().Get(oldAccessKey)
	if !has {
		return "", "", errors.New("refresh token not valid")
	}

	ormService.GetRedis().Del(oldAccessKey)

	newAccessKey := t.generateAndStoreAccessKey(ormService, id, t.accessTokenTTL)

	newAccessToken, err = t.GenerateTokenPair(id, newAccessKey, t.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = t.GenerateTokenPair(id, newAccessKey, t.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	t.addUserAccessKeyList(ormService, id, newAccessKey, oldAccessKey, t.refreshTokenTTL)

	return newAccessToken, newRefreshToken, err
}

func (t *Authentication) LogoutCurrentSession(ormService *beeorm.Engine, accessKey string) {
	cacheService := ormService.GetRedis()

	cacheService.Del(accessKey)

	tokenListKey := generateUserTokenListKey(getUserIDFromAccessKey(accessKey))
	tokenList, has := cacheService.Get(tokenListKey)
	if has {
		var newTokenList = make([]string, 0)
		tokenArr := strings.Split(tokenList, accessListSeparator)
		if len(tokenArr) != 0 {
			for i := range tokenArr {
				if tokenArr[i] != accessKey {
					newTokenList = append(newTokenList, tokenArr[i])
				}
			}
			if len(newTokenList) != 0 {
				cacheService.Set(tokenListKey, strings.Join(newTokenList, accessListSeparator), redis.KeepTTL)
			}
		}
	}
}

func (t *Authentication) LogoutAllSessions(ormService *beeorm.Engine, id uint64) {
	tokenListKey := generateUserTokenListKey(id)
	cacheService := ormService.GetRedis()

	tokenList, has := cacheService.Get(tokenListKey)
	if has && tokenList != "" {
		tokenArr := strings.Split(tokenList, accessListSeparator)
		if len(tokenArr) != 0 {
			cacheService.Del(tokenArr...)
		}
	}
	cacheService.Del(tokenListKey)
}

func (t *Authentication) GenerateTokenPair(id uint64, accessKey string, ttl int) (string, error) {
	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	now := t.clockService.Now().Unix()

	payload := map[string]string{
		"jti": accessKey,
		"sub": strconv.FormatUint(id, 10),
		"exp": strconv.FormatInt(now+int64(ttl), 10),
		"iat": strconv.FormatInt(now, 10),
	}

	return t.jwtService.EncodeJWT(t.secret, headers, payload)
}

func (t *Authentication) generateAndStoreAccessKey(ormService *beeorm.Engine, id uint64, ttl int) string {
	key := generateAccessKey(id, t.uuidService.Generate())
	ormService.GetRedis().Set(key, "", ttl)
	return key
}

func (t *Authentication) addUserAccessKeyList(ormService *beeorm.Engine, id uint64, accessKey, oldAccessKey string, ttl int) {
	key := generateUserTokenListKey(id)
	cacheService := ormService.GetRedis()
	res, has := cacheService.Get(key)
	if !has {
		cacheService.Set(key, accessKey, ttl)
		return
	}

	currentTokenArr := strings.Split(res, accessListSeparator)
	if len(currentTokenArr) >= maxUserAccessKeysAllowed {
		cacheService.Del(currentTokenArr[0])
		currentTokenArr = currentTokenArr[1:]
	}

	if oldAccessKey == "" {
		currentTokenArr = append(currentTokenArr, accessKey)
		cacheService.Set(key, strings.Join(currentTokenArr, accessListSeparator), ttl)
		return
	}

	var finalTokenArr = make([]string, 0)
	finalTokenArr = append(finalTokenArr, accessKey)

	if oldAccessKey != "" {
		for i := range currentTokenArr {
			if currentTokenArr[i] != oldAccessKey {
				finalTokenArr = append(finalTokenArr, currentTokenArr[i])
			}
		}
	}
	if len(finalTokenArr) == 0 {
		cacheService.Del(key)
	} else {
		cacheService.Set(key, strings.Join(finalTokenArr, accessListSeparator), ttl)
	}
}

func generateAccessKey(id uint64, uuid string) string {
	return fmt.Sprintf("%s%s%d%s%s", accessKeyPrefix, separator, id, separator, uuid)
}

func getUserIDFromAccessKey(accessKey string) uint64 {
	accessArr := strings.Split(accessKey, separator)
	if len(accessArr) == 3 {
		userIDInt, _ := strconv.ParseInt(accessArr[1], 10, 0)
		return uint64(userIDInt)
	}
	return 0
}

func generateUserTokenListKey(id uint64) string {
	return fmt.Sprintf("%s%s%d", userAccessListPrefix, separator, id)
}

func factorySMSProviders(country string) *sms.Provider {
	providers := &sms.Provider{
		Primary:   sms.Twilio,
		Secondary: sms.Sinch,
	}

	if country == "IR" {
		providers = &sms.Provider{
			Primary:   sms.Kavenegar,
			Secondary: sms.Twilio,
		}
	}
	return providers
}
