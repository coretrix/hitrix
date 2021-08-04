package account

import (
	"github.com/coretrix/hitrix/pkg/binding"
	"github.com/coretrix/hitrix/pkg/view/account"
	"github.com/coretrix/hitrix/service"
	"github.com/gin-gonic/gin"
	"github.com/juju/errors"
)

type LoginDevForm struct {
	Username string `binding:"required,min=6,max=60" json:"Username"  form:"username"`
	Password string `binding:"required,min=8,max=60" json:"Password"  form:"password"`
}

func (l *LoginDevForm) Login(c *gin.Context) (string, string, error) {
	err := binding.ShouldBindJSON(c, l)
	if err != nil {
		return "", "", err
	}

	ormService := service.DI().OrmEngineForContext(c.Request.Context())

	passwordService, has := service.DI().Password()
	if !has {
		return "", "", errors.New("Please load Password service")
	}

	devPanelUserEntity := service.DI().App().DevPanel.UserEntity
	ok := ormService.CachedSearchOne(devPanelUserEntity, "UserEmailIndex", l.Username)

	if !ok {
		return "", "", errors.New("invalid username or password")
	}

	if !passwordService.VerifyPassword(l.Password, devPanelUserEntity.GetPassword()) {
		return "", "", errors.New("invalid username or password")
	}

	token, refreshToken, err := account.GenerateDevTokenAndRefreshToken(ormService, devPanelUserEntity.GetID())
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}
