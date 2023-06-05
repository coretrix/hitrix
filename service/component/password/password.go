package password

import (
	"github.com/coretrix/hitrix/service/component/config"
)

type NewPasswordManagerFunc func(configService config.IConfig) IPassword

type IPassword interface {
	VerifyPassword(password string, hash string) bool
	HashPassword(password string) (string, error)
}
