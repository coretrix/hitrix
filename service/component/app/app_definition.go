package app

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm"
)

const ModeLocal = "local"
const ModeTest = "test"
const ModeDev = "dev"
const ModeDemo = "demo"
const ModeProd = "prod"

type IDevPanelUserEntity interface {
	beeorm.Entity
	GetUsername() string
	GetPassword() string
}

type DevPanel struct {
	UserEntity IDevPanelUserEntity
	Router     func(ginEngine *gin.Engine)
	PoolStream *string
	PoolSearch *string
}

type App struct {
	Mode           string
	Name           string
	ParallelTestID string
	Secret         string
	Flags          *Flags
	Scripts        []string
	DevPanel       *DevPanel
	GlobalContext  context.Context
	CancelContext  context.CancelFunc
}

func (app *App) IsInLocalMode() bool {
	return app.IsInMode(ModeLocal)
}

func (app *App) IsInTestMode() bool {
	return app.IsInMode(ModeTest)
}

func (app *App) IsInProdMode() bool {
	return app.IsInMode(ModeProd)
}

func (app *App) IsInDevMode() bool {
	return app.Mode == ModeDev
}

func (app *App) IsInDemoMode() bool {
	return app.Mode == ModeDemo
}

func (app *App) IsInMode(mode string) bool {
	return app.Mode == mode
}
