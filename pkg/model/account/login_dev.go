package account

import (
	"github.com/coretrix/hitrix"
	hitrixErrors "github.com/coretrix/hitrix/pkg/error"
	"github.com/coretrix/hitrix/pkg/view/account"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

type LoginDevForm struct {
	Username string `binding:"required,min=3,max=60" json:"Username"`
	Password string `binding:"required,min=3,max=60" json:"Password"`
}

func (l *LoginDevForm) Login(c *gin.Context) (string, string, error) {
	err := c.ShouldBindJSON(l)
	if err != nil {
		return "", "", err
	}

	ormService, has := hitrix.DIC().OrmEngineForContext(c)

	if !has {
		return "", "", errors.New("orm is not registered")
	}
	if l.Username == "DevUser" && l.Password == "$umm3r" {
		token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService)
		if err != nil {
			return "", "", err
		}

		return token, refreshToken, err
	}

	return "", "", hitrixErrors.HandleCustomErrors(map[string]string{"Password": "Wrong password"}, c)
}
