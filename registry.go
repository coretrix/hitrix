package hitrix

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/coretrix/hitrix/service/registry"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"

	"github.com/gin-gonic/gin"

	"github.com/sarulabs/di"
)

type Registry struct {
	app                 *app.App
	servicesDefinitions []*service.Definition
}

func New(appName string, secret string) *Registry {
	mode, hasMode := os.LookupEnv("APP_MODE")
	if !hasMode {
		mode = app.ModeLocal
	}

	r := &Registry{
		app: &app.App{Mode: mode, Name: appName, Secret: secret},
	}
	return r
}

func (r *Registry) Build() (*Hitrix, func()) {
	r.initializeIoCHandlers()

	flags := service.DI().App().Flags
	if flags.Bool("list-scripts") {
		listScrips()
	}
	scriptToRun := flags.String("run-script")
	ctx, cancel := context.WithCancel(context.Background())
	h := &Hitrix{ctx: ctx, cancel: cancel, done: make(chan bool), exit: make(chan int)}
	if scriptToRun != "" {
		h.runDynamicScrips(ctx, scriptToRun)
	}

	return h, func() {
		if r := recover(); r != nil {
			errorLogger, has := service.DI().ErrorLogger()
			if has {
				errorLogger.LogError(r)
			} else {
				log.Print(r.(string))
			}
		}
	}
}

func (r *Registry) RegisterDevPanel(devPanelUserEntity app.DevPanelUserEntity, router func(ginEngine *gin.Engine), poolStream *string) *Registry {
	if devPanelUserEntity == nil {
		panic("devPanelUserEntity cannot be nil")
	}
	if router == nil {
		panic("router cannot be nil")
	}

	r.app.DevPanel = &app.DevPanel{UserEntity: devPanelUserEntity, Router: router, PoolStream: poolStream}
	return r
}

func (r *Registry) RegisterDIService(service ...*service.Definition) *Registry {
	r.servicesDefinitions = append(r.servicesDefinitions, service...)
	return r
}

func (r *Registry) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*service.Definition{
		registry.ServiceApp(r.app),
		registry.ServiceConfig(),
	}

	flagsRegistry := &app.FlagsRegistry{Flags: make(map[string]interface{})}
	for _, def := range append(defaultDefinitions, r.servicesDefinitions...) {
		if def == nil {
			continue
		}

		var scope string
		if def.Global {
			scope = di.App
		} else {
			scope = di.Request
		}
		if def.Script {
			r.app.Scripts = append(r.app.Scripts, def.Name)
		}

		err := ioCBuilder.Add(di.Def{
			Name:  def.Name,
			Scope: scope,
			Build: def.Build,
			Close: def.Close,
		})
		if err != nil {
			panic(err)
		}
		if def.Flags != nil && !flag.Parsed() {
			def.Flags(flagsRegistry)
		}
	}
	if !flag.Parsed() {
		flagsRegistry.Bool("list-scripts", false, "list all available scripts")
		flagsRegistry.String("run-script", "", "run script")
	}

	err := ioCBuilder.Add()

	if err != nil {
		panic(err)
	}
	service.SetContainer(ioCBuilder.Build())
	if !flag.Parsed() {
		flag.Parse()
	}
	r.app.Flags = &app.Flags{Registry: flagsRegistry}
}
