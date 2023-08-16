package authentication

import (
	"context"
	"errors"
	"fmt"
	"math"
	mail2 "net/mail"
	"strconv"
	"strings"
	"time"

	redisearch "github.com/coretrix/beeorm-redisearch-plugin"
	"github.com/latolukasz/beeorm/v2"
	"github.com/redis/go-redis/v9"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
	"github.com/coretrix/hitrix/service/component/generator"
	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/mail"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/coretrix/hitrix/service/component/social"
	"github.com/coretrix/hitrix/service/component/uuid"
)

const (
	separator                = ":"
	accessListSeparator      = ";"
	accessKeyPrefix          = "ACCESS"
	userAccessListPrefix     = "USER_KEYS"
	maxUserAccessKeysAllowed = 10

	SocialLoginGoogle   = "google"
	SocialLoginFacebook = "facebook"
	SocialLoginApple    = "apple"
)

type AuthenticatableEntity interface {
	beeorm.Entity
	CanAuthenticate() bool
}

type OTPProviderEntity interface {
	beeorm.Entity
	AuthenticatableEntity
	GetPhoneFieldName() string
	GetEmailFieldName() string
}

type AuthProviderEntity interface {
	beeorm.Entity
	AuthenticatableEntity
	GetUniqueFieldName() string
	GetPassword() string
}

type EmailAuthEntity interface {
	beeorm.Entity
	AuthenticatableEntity
	GetPassword() string
	GetEmailFieldName() string
}

type Authentication struct {
	accessTokenTTL       int
	refreshTokenTTL      int
	otpTTL               int
	otpLength            int
	passwordService      password.IPassword
	errorLoggerService   errorlogger.ErrorLogger
	appService           *app.App
	jwtService           *jwt.JWT
	mailService          *mail.ISender
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
	otpLength int,
	appService *app.App,
	generatorService generator.IGenerator,
	errorLoggerService errorlogger.ErrorLogger,
	clockService clock.IClock,
	passwordService password.IPassword,
	jwtService *jwt.JWT,
	mailService *mail.ISender,
	socialServiceMapping map[string]social.IUserData,
	uuidService uuid.IUUID,
) *Authentication {
	return &Authentication{
		secret:               secret,
		accessTokenTTL:       accessTokenTTL,
		refreshTokenTTL:      refreshTokenTTL,
		otpTTL:               otpTTL,
		otpLength:            otpLength,
		passwordService:      passwordService,
		errorLoggerService:   errorLoggerService,
		jwtService:           jwtService,
		appService:           appService,
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

func (t *Authentication) GenerateAndSendOTPEmail(ormService *datalayer.ORM, email, template, from, title string) (*GenerateOTPEmail, error) {
	_, err := mail2.ParseAddress(email)

	if err != nil {
		return nil, errors.New("mail address not valid")
	}

	var code int64
	if t.otpLength == 0 {
		code = t.generatorService.GenerateRandomRangeNumber(10000, 99999)
	} else {
		// example, if t.otpLength = 5 (the default)
		// min = 10 ^ (5-1) 	==> min = 10 ^ 4 		==> min = 10000
		// max = (10 ^ 5) - 1   ==> max = 100000 - 1	==> max = 99999
		min := int64(math.Pow(10, float64(t.otpLength-1)))
		max := int64(math.Pow(10, float64(t.otpLength))) - 1
		code = t.generatorService.GenerateRandomRangeNumber(min, max)
	}

	if t.mailService == nil {
		panic("mail service is not registered")
	}

	mailService := *t.mailService

	err = mailService.SendTemplate(ormService, &mail.Message{
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

func (t *Authentication) VerifySocialLogin(ctx context.Context, source, token string, isAndroid bool) (*social.UserData, error) {
	socialProvider, ok := t.socialServiceMapping[source]
	if !ok {
		return nil, errors.New("not supported social provider: " + source)
	}

	return socialProvider.GetUserData(ctx, token, isAndroid)
}

func (t *Authentication) AuthenticateOTP(
	ormService *datalayer.ORM,
	phone string,
	entity OTPProviderEntity,
) (accessToken string, refreshToken string, err error) {
	q := &redisearch.RedisSearchQuery{}
	q.FilterString(entity.GetPhoneFieldName(), phone)

	found := ormService.RedisSearchOne(entity, q)
	if !found {
		return "", "", errors.New("invalid credentials")
	}

	if !entity.CanAuthenticate() {
		return "", "", errors.New("cannot authenticate this entity")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) AuthenticateOTPEmail(
	ormService *datalayer.ORM,
	email string,
	entity OTPProviderEntity,
) (accessToken string, refreshToken string, err error) {
	q := &redisearch.RedisSearchQuery{}
	q.FilterString(entity.GetEmailFieldName(), email)

	found := ormService.RedisSearchOne(entity, q)
	if !found {
		return "", "", errors.New("invalid credentials")
	}

	if !entity.CanAuthenticate() {
		return "", "", errors.New("cannot authenticate this entity")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) Authenticate(
	ormService *datalayer.ORM,
	uniqueValue string,
	password string,
	entity AuthProviderEntity,
) (accessToken string, refreshToken string, err error) {
	q := &redisearch.RedisSearchQuery{}
	q.FilterString(entity.GetUniqueFieldName(), uniqueValue)

	found := ormService.RedisSearchOne(entity, q)
	if !found {
		return "", "", errors.New("invalid user/pass")
	}

	if !t.passwordService.VerifyPassword(password, entity.GetPassword()) {
		return "", "", errors.New("invalid user/pass")
	}

	if !entity.CanAuthenticate() {
		return "", "", errors.New("cannot authenticate this entity")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) AuthenticateEmail(
	ormService *datalayer.ORM,
	email string,
	password string,
	entity EmailAuthEntity,
) (accessToken string, refreshToken string, err error) {
	found := ormService.CachedSearchOne(entity, "CachedQueryEmail", email)
	if !found {
		return "", "", errors.New("invalid credentials")
	}

	if !t.passwordService.VerifyPassword(password, entity.GetPassword()) {
		return "", "", errors.New("invalid user/pass")
	}

	if !entity.CanAuthenticate() {
		return "", "", errors.New("cannot authenticate this entity")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) AuthenticateByID(
	ormService *datalayer.ORM,
	id uint64,
	entity AuthProviderEntity,
) (accessToken string, refreshToken string, err error) {
	exists := ormService.LoadByID(id, entity)

	if !exists {
		return "", "", errors.New("id_does_not_exists")
	}

	if !entity.CanAuthenticate() {
		return "", "", errors.New("cannot authenticate this entity")
	}

	return t.generateUserTokens(ormService, entity.GetID())
}

func (t *Authentication) generateUserTokens(ormService *datalayer.ORM, ID uint64) (accessToken string, refreshToken string, err error) {
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

func (t *Authentication) VerifyAccessToken(ormService *datalayer.ORM, accessToken string, entity beeorm.Entity) (map[string]string, error) {
	payload, err := t.jwtService.VerifyJWTAndGetPayload(t.secret, accessToken, t.clockService.Now().Unix())
	if err != nil {
		return nil, err
	}

	id, err := strconv.ParseUint(payload["sub"], 10, 64)
	if err != nil {
		return nil, err
	}

	accessKey := payload["jti"]

	_, has := ormService.GetRedis(t.appService.RedisPools.Persistent).Get(accessKey)
	if !has {
		return nil, errors.New("access key not found")
	}

	found := ormService.LoadByID(id, entity)
	if !found {
		return nil, errors.New("user_not_found")
	}

	return payload, nil
}

// nolint // info
func (t *Authentication) RefreshToken(ormService *datalayer.ORM, refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
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

	_, has := ormService.GetRedis(t.appService.RedisPools.Persistent).Get(oldAccessKey)
	if !has {
		return "", "", errors.New("refresh token not valid")
	}

	ormService.GetRedis(t.appService.RedisPools.Persistent).Del(oldAccessKey)

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

func (t *Authentication) LogoutCurrentSession(ormService *datalayer.ORM, accessKey string) {
	cacheService := ormService.GetRedis(t.appService.RedisPools.Persistent)

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

func (t *Authentication) LogoutAllSessions(ormService *datalayer.ORM, id uint64) {
	tokenListKey := generateUserTokenListKey(id)
	cacheService := ormService.GetRedis(t.appService.RedisPools.Persistent)

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

func (t *Authentication) generateAndStoreAccessKey(ormService *datalayer.ORM, id uint64, ttl int) string {
	key := generateAccessKey(id, t.uuidService.Generate())
	ormService.GetRedis(t.appService.RedisPools.Persistent).Set(key, "", time.Second*time.Duration(ttl))

	return key
}

func (t *Authentication) addUserAccessKeyList(ormService *datalayer.ORM, id uint64, accessKey, oldAccessKey string, ttl int) {
	key := generateUserTokenListKey(id)
	cacheService := ormService.GetRedis(t.appService.RedisPools.Persistent)

	res, has := cacheService.Get(key)
	if !has {
		cacheService.Set(key, accessKey, time.Second*time.Duration(ttl))

		return
	}

	currentTokenArr := strings.Split(res, accessListSeparator)
	if len(currentTokenArr) >= maxUserAccessKeysAllowed {
		cacheService.Del(currentTokenArr[0])
		currentTokenArr = currentTokenArr[1:]
	}

	if oldAccessKey == "" {
		currentTokenArr = append(currentTokenArr, accessKey)
		cacheService.Set(key, strings.Join(currentTokenArr, accessListSeparator), time.Second*time.Duration(ttl))

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
		cacheService.Set(key, strings.Join(finalTokenArr, accessListSeparator), time.Second*time.Duration(ttl))
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
