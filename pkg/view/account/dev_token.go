package account

import (
	// #nosec
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/coretrix/hitrix"

	"github.com/gin-gonic/gin"
	"github.com/summer-solutions/orm"
)

const LoggedDevPanelUserEntity = "logged_dev_panel_user_entity"

func GenerateDevTokenAndRefreshToken(ormService *orm.Engine, userID uint64) (string, string, error) {
	appService := hitrix.DIC().App()
	token, err := generateTokenValue(appService.Secret(), userID, 3600)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateTokenValue(appService.Secret(), userID, 7200)
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
		7200,
	)
	// #nosec
	redisService.Set(
		fmt.Sprintf(
			"%x",
			md5.Sum([]byte(refreshToken)),
		),
		refreshToken,
		7200,
	)

	return token, refreshToken, nil
}

func generateTokenValue(secret string, id interface{}, expire int64) (string, error) {
	jwtToken, has := hitrix.DIC().JWT()
	if !has {
		panic("Please load JWT service")
	}

	app := hitrix.DIC().App()

	headers := map[string]string{
		"algo": "HS256",
		"type": "JWT",
	}

	payload := map[string]string{
		"iss":  app.Secret(),
		"exp":  fmt.Sprintf("%v", expire),
		"user": fmt.Sprintf("%v", id),
	}

	jwtValue, err := jwtToken.EncodeJWT(secret, headers, payload)

	return jwtValue, err
}

func IsValidDevRefreshToken(c *gin.Context, token string) error {
	appService := hitrix.DIC().App()
	if _, err := isValid(token, appService.Secret(), 7200); err != nil {
		return err
	}

	return verifyDevUser(c, token)
}

func IsValidDevToken(c *gin.Context, token string) error {
	appService := hitrix.DIC().App()
	if _, err := isValid(token, appService.Secret(), 3600); err != nil {
		return err
	}

	return verifyDevUser(c, token)
}

func verifyDevUser(c *gin.Context, token string) error {
	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		panic("orm is not registered")
	}

	redisService := ormService.GetRedis()
	// #nosec
	v, has := redisService.Get(fmt.Sprintf("%x", md5.Sum([]byte(token))))

	if !has || strings.Compare(v, token) != 0 {
		return fmt.Errorf("token doesnt match")
	}

	return nil
}

func isValid(token, tokenSecret string, tokenExpire int64) (uint64, error) {
	jwtService, has := hitrix.DIC().JWT()
	if !has {
		panic("Please load JWT service")
	}

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
