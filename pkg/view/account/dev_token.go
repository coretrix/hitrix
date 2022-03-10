package account

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/service"
)

const LoggedDevPanelUserEntity = "logged_dev_panel_user_entity"
const expireTimeToken = 3600
const expireTimeRefreshToken = 7200

func GenerateDevTokenAndRefreshToken(ormService *beeorm.Engine, userID uint64) (string, string, error) {
	appService := service.DI().App()
	token, err := generateTokenValue(appService.Secret, userID, time.Now().Unix()+expireTimeToken)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateTokenValue(appService.Secret, userID, time.Now().Unix()+expireTimeRefreshToken)
	if err != nil {
		return "", "", err
	}

	redisService := ormService.GetRedis()
	// #nosec
	redisService.Set(
		fmt.Sprintf(
			"%x",
			md5.Sum([]byte(token)),
		),
		token,
		expireTimeToken,
	)
	// #nosec
	redisService.Set(
		fmt.Sprintf(
			"%x",
			md5.Sum([]byte(refreshToken)),
		),
		refreshToken,
		expireTimeRefreshToken,
	)

	return token, refreshToken, nil
}

func generateTokenValue(secret string, id interface{}, expire int64) (string, error) {
	jwtService := service.DI().JWT()

	app := service.DI().App()

	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	payload := map[string]string{
		"iss":  app.Secret,
		"exp":  fmt.Sprintf("%v", expire),
		"user": fmt.Sprintf("%v", id),
	}

	jwtValue, err := jwtService.EncodeJWT(secret, headers, payload)

	return jwtValue, err
}

func IsValidDevRefreshToken(c *gin.Context, token string) error {
	appService := service.DI().App()

	userID, err := isValid(token, appService.Secret, expireTimeRefreshToken)
	if err != nil {
		return err
	}

	return verifyDevUser(c, userID, token)
}

func IsValidDevToken(c *gin.Context, token string) error {
	appService := service.DI().App()
	userID, err := isValid(token, appService.Secret, expireTimeToken)
	if err != nil {
		return err
	}

	return verifyDevUser(c, userID, token)
}

func verifyDevUser(c *gin.Context, userID uint64, token string) error {
	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	redisService := ormService.GetRedis()
	// #nosec
	v, has := redisService.Get(fmt.Sprintf("%x", md5.Sum([]byte(token))))

	if !has || strings.Compare(v, token) != 0 {
		return fmt.Errorf("token doesnt match")
	}

	userEntity := service.DI().App().DevPanel.UserEntity
	has = ormService.LoadByID(userID, userEntity)

	if !has {
		return errors.New("invalid user")
	}

	c.Set(LoggedDevPanelUserEntity, userEntity)

	return nil
}

func isValid(token, tokenSecret string, tokenExpire int64) (uint64, error) {
	jwtService := service.DI().JWT()

	err := jwtService.VerifyJWT(tokenSecret, token, tokenExpire)

	if err != nil {
		return 0, err
	}

	data := strings.Split(token, ".")
	dbyte, err := base64.StdEncoding.DecodeString(data[1])
	if err != nil {
		return 0, err
	}
	payload := make(map[string]string)

	err = json.Unmarshal(dbyte, &payload)
	if err != nil {
		return 0, err
	}

	userID, ok := payload["user"]

	if !ok {
		return 0, fmt.Errorf("invalid token payload")
	}

	id, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}
