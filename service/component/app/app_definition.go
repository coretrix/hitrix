package app

import (
	"context"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/latolukasz/beeorm/v2"
)

const ModeLocal = "local"
const ModeTest = "test"
const ModeDev = "dev"
const ModeDemo = "demo"
const ModeProd = "prod"
const ModeQA = "qa"

type IDevPanelUserEntity interface {
	beeorm.Entity
	GetUsername() string
	GetPassword() string
}

type DevPanel struct {
	UserEntity IDevPanelUserEntity
	Router     func(ginEngine *gin.Engine)
}

type RedisPools struct {
	Cache      string
	Persistent string
	Stream     string
	Search     string
}

type App struct {
	Mode           string
	Name           string
	ParallelTestID string
	Secret         string
	Flags          *Flags
	Scripts        []string
	DevPanel       *DevPanel
	RedisPools     *RedisPools
	GlobalContext  context.Context
	CancelContext  context.CancelFunc
	waitGroup      *sync.WaitGroup
}

func (app *App) IsInLocalMode() bool {
	return app.Mode == ModeLocal
}

func (app *App) IsInTestMode() bool {
	return app.Mode == ModeTest
}

func (app *App) IsInProdMode() bool {
	return app.Mode == ModeProd
}

func (app *App) IsInDevMode() bool {
	return app.Mode == ModeDev
}

func (app *App) IsInDemoMode() bool {
	return app.Mode == ModeDemo
}

func (app *App) IsInQAMode() bool {
	return app.Mode == ModeQA
}

func (app *App) IsInMode(mode string) bool {
	return app.Mode == mode
}

func (app *App) InitWaitGroup() {
	if app.waitGroup == nil {
		app.waitGroup = &sync.WaitGroup{}
	}
}

func (app *App) Add(delta int) {
	app.waitGroup.Add(delta)
}

func (app *App) Done() {
	app.waitGroup.Done()
}
func (app *App) Wait() {
	app.waitGroup.Wait()
}
