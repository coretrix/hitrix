package password

import (
	"crypto/sha256"
	"encoding/base64"

	"github.com/coretrix/hitrix/service/component/config"
)

type SimpleManager struct {
}

func NewSimpleManager(_ config.IConfig) IPassword {
	return &SimpleManager{}
}

func (p *SimpleManager) VerifyPassword(password string, hash string) bool {
	passwordHash, err := p.HashPassword(password)

	if err != nil {
		panic(err)
	}

	return passwordHash == hash
}

func (p *SimpleManager) HashPassword(password string) (string, error) {
	sha256Hash := sha256.New()
	_, err := sha256Hash.Write([]byte(password))

	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(sha256Hash.Sum(nil)), nil
}
