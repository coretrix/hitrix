package authentication

import (
	"errors"
	"strconv"
	"time"

	"github.com/coretrix/hitrix/service/component/jwt"
	"github.com/coretrix/hitrix/service/component/password"
	"github.com/latolukasz/orm"
)

type EmailPasswordProviderEntity interface {
	orm.Entity
	GetEmailCachedIndexName() string
	GetPassword() string
}
type Authentication struct {
	entity          EmailPasswordProviderEntity
	accessTokenTTL  int
	refreshTokenTTL int
	passwordService *password.Password
	jwtService      *jwt.JWT
	ormService      *orm.Engine
	secret          string
}

func (t *Authentication) Authenticate(email string, password string) (accessToken string, refreshToken string, err error) {
	userEntity := t.entity.(orm.Entity)
	found := t.ormService.CachedSearchOne(userEntity, t.entity.GetEmailCachedIndexName(), email)
	if !found {
		return "", "", errors.New("user_not_found")
	}

	if !t.passwordService.VerifyPassword(password, t.entity.GetPassword()) {
		return "", "", errors.New("invalid_password")
	}

	accessToken, err = t.generateTokenPair(userEntity.GetID(), t.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	refreshToken, err = t.generateTokenPair(userEntity.GetID(), t.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (t *Authentication) VerifyAccessToken(accessToken string, entity orm.Entity) error {
	payload, err := t.jwtService.VerifyJWTAndGetPayload(t.secret, accessToken, time.Now().Unix())
	if err != nil {
		return err
	}

	id, err := strconv.ParseUint(payload["sub"], 10, 64)
	if err != nil {
		return err
	}

	found := t.ormService.LoadByID(id, entity)
	if !found {
		return errors.New("user_not_found")
	}

	return nil
}

func (t *Authentication) RefreshToken(refreshToken string) (newAccessToken string, newRefreshToken string, err error) {
	payload, err := t.jwtService.VerifyJWTAndGetPayload(t.secret, refreshToken, time.Now().Unix())
	if err != nil {
		return "", "", err
	}

	id, err := strconv.ParseUint(payload["sub"], 10, 64)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err = t.generateTokenPair(id, t.accessTokenTTL)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = t.generateTokenPair(id, t.refreshTokenTTL)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, err
}

func (t *Authentication) generateTokenPair(id uint64, ttl int) (string, error) {
	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	now := time.Now().Unix()

	payload := map[string]string{
		"sub": strconv.FormatUint(id, 10),
		"exp": strconv.FormatInt(now+int64(ttl), 10),
		"iat": strconv.FormatInt(now, 10),
	}

	return t.jwtService.EncodeJWT(t.secret, headers, payload)
}

func NewAuthenticationService(
	entity EmailPasswordProviderEntity,
	secret string,
	accessTokenTTL int,
	refreshTokenTTL int,
	ormService *orm.Engine,
	passwordService *password.Password,
	jwtService *jwt.JWT,
) *Authentication {
	return &Authentication{
		entity:          entity,
		secret:          secret,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
		passwordService: passwordService,
		jwtService:      jwtService,
		ormService:      ormService,
	}
}
