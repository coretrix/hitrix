package account

import (
	// #nosec
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/coretrix/hitrix"

	"github.com/gin-gonic/gin"
	"github.com/summer-solutions/orm"
)

func GenerateDevTokenAndRefreshToken(ormService *orm.Engine) (string, string, error) {
	token := "test"
	refreshToken := "test"

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

func IsValidDevRefreshToken(c *gin.Context, token string) error {
	return nil

	//securityService := spring.IoC.Get(service.LoginSecurityService).(*loginsecurity.LoginSecurity)
	//if _, err := isValid(token, securityService.RefreshTokenSecret, securityService.RefreshTokenExpire); err != nil {
	//	return err
	//}
	//
	//return verifyDevUser(c, token)
}

func IsValidDevToken(c *gin.Context, token string) error {
	return nil
	//securityService := spring.IoC.Get(service.LoginSecurityService).(*loginsecurity.LoginSecurity)
	//if _, err := isValid(token, securityService.LoginTokenSecret, securityService.LoginTokenExpire); err != nil {
	//	return err
	//}
	//
	//return verifyDevUser(c, token)
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

//func isValid(token, tokenSecret string, tokenExpire int64) (uint64, error) {
//	err := loginsecurity.VerifyJWT(tokenSecret, token, tokenExpire)
//
//	if err != nil {
//		return 0, err
//	}
//
//	data := strings.Split(token, ".")
//	dbyte, err := base64.StdEncoding.DecodeString(data[1])
//	if err != nil {
//		return 0, err
//	}
//	payload := make(map[string]string)
//
//	err = json.Unmarshal(dbyte, &payload)
//	if err != nil {
//		return 0, err
//	}
//
//	userID, ok := payload["user"]
//
//	if !ok {
//		return 0, fmt.Errorf("invalid token payload")
//	}
//
//	id, err := strconv.ParseUint(userID, 10, 64)
//	if err != nil {
//		return 0, err
//	}
//
//	return id, nil
//}
