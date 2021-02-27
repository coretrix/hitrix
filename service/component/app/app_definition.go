package app

import (
	"github.com/gin-gonic/gin"
	"github.com/summer-solutions/orm"
)

const ModeLocal = "local"
const ModeTest = "test"
const ModeDemo = "demo"
const ModeProd = "prod"

type DevPanelUserEntity interface {
	orm.Entity
	GetUsername() string
	GetPassword() string
}

type DevPanel struct {
	UserEntity DevPanelUserEntity
	Router     func(ginEngine *gin.Engine)
	PoolStream *string
}

type App struct {
	Mode     string
	Name     string
	Secret   string
	Flags    *Flags
	Scripts  []string
	DevPanel *DevPanel
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

func (app *App) IsInDemoMode() bool {
	return app.Mode == ModeDemo
}

func (app *App) IsInMode(mode string) bool {
	return app.Mode == mode
}