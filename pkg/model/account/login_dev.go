package account

import (
	"fmt"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/pkg/view/account"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

type LoginDevForm struct {
	Username string `binding:"required,min=6,max=60" json:"Username"`
	Password string `binding:"required,min=8,max=60" json:"Password"`
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

	passwordService, has := hitrix.DIC().Password()
	if !has {
		return "", "", errors.New("Please load Password service")
	}

	devPanelUserEntity := hitrix.DIC().App().DevPanel().UserEntity
	ok := ormService.CachedSearchOne(devPanelUserEntity, "UserEmailIndex", l.Username)

	if !ok {
		return "", "", fmt.Errorf("invalid username or password")
	}

	if !passwordService.VerifyPassword(l.Password, devPanelUserEntity.GetPassword()) {
		return "", "", fmt.Errorf("invalid username or password")
	}

	token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService, devPanelUserEntity.GetID())
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}
