package authentication

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/jwt"
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
)

type AuthProviderEntity interface {
	beeorm.Entity
	GetUniqueFieldName() string
	GetPassword() string
}

type Authentication struct {
	accessTokenTTL       int
	refreshTokenTTL      int
	passwordService      password.IPassword
	appService           *app.App
	jwtService           *jwt.JWT
	socialServiceMapping map[string]social.IUserData
	clockService         clock.IClock
	uuidService          uuid.IUUID
	secret               string
}

func NewAuthenticationService(
	secret string,
	accessTokenTTL int,
	refreshTokenTTL int,
	appService *app.App,
	clockService clock.IClock,
	passwordService password.IPassword,
	jwtService *jwt.JWT,
	socialServiceMapping map[string]social.IUserData,
	uuidService uuid.IUUID,
) *Authentication {
	return &Authentication{
		secret:               secret,
		accessTokenTTL:       accessTokenTTL,
		refreshTokenTTL:      refreshTokenTTL,
		passwordService:      passwordService,
		jwtService:           jwtService,
		appService:           appService,
		clockService:         clockService,
		socialServiceMapping: socialServiceMapping,
		uuidService:          uuidService,
	}
}

type GenerateOTPEmail struct {
	Email          string
	ExpirationTime string
	Token          string
}

func (t *Authentication) VerifySocialLogin(source, token string) (*social.UserData, error) {
	socialProvider, ok := t.socialServiceMapping[source]
	if !ok {
		return nil, errors.New("not supported social provider: " + source)
	}
	return socialProvider.GetUserData(token)
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

func (t *Authentication) LogoutCurrentSession(ormService *beeorm.Engine, accessKey string) {
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

func (t *Authentication) LogoutAllSessions(ormService *beeorm.Engine, id uint64) {
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

func (t *Authentication) generateAndStoreAccessKey(ormService *beeorm.Engine, id uint64, ttl int) string {
	key := generateAccessKey(id, t.uuidService.Generate())
	ormService.GetRedis(t.appService.RedisPools.Persistent).Set(key, "", ttl)
	return key
}

func (t *Authentication) addUserAccessKeyList(ormService *beeorm.Engine, id uint64, accessKey, oldAccessKey string, ttl int) {
	key := generateUserTokenListKey(id)
	cacheService := ormService.GetRedis(t.appService.RedisPools.Persistent)
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
