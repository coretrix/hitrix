package hitrix

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"github.com/summer-solutions/orm"

	"github.com/sarulabs/di"
)

type DevPanelUserEntity interface {
	orm.Entity
	GetUsername() string
	GetPassword() string
}

type Registry struct {
	app                 *AppDefinition
	servicesDefinitions []*ServiceDefinition
}

func New(appName string, secret string) *Registry {
	mode, hasMode := os.LookupEnv("SPRING_MODE")
	if !hasMode {
		mode = ModeLocal
	}
	app := &AppDefinition{mode: mode, name: appName, secret: secret}
	r := &Registry{app: app}
	return r
}

func (r *Registry) Build() (*Hitrix, func()) {
	r.initializeIoCHandlers()

	flags := DIC().App().Flags()
	if flags.Bool("list-scripts") {
		listScrips()
	}
	scriptToRun := flags.String("run-script")
	ctx, cancel := context.WithCancel(context.Background())
	s := &Hitrix{registry: r, ctx: ctx, cancel: cancel, done: make(chan bool), exit: make(chan int)}
	if scriptToRun != "" {
		s.runDynamicScrips(ctx, scriptToRun)
	}

	return s, func() {
		if r := recover(); r != nil {
			errorLogger, has := DIC().ErrorLogger()
			if has {
				errorLogger.LogRecover(r)
			} else {
				log.Print(r.(string))
			}
		}
	}
}

func (r *Registry) RegisterDevPanel(devPanelUserEntity DevPanelUserEntity, router func(ginEngine *gin.Engine), poolStream *string) *Registry {
	if devPanelUserEntity == nil {
		panic("devPanelUserEntity cannot be nil")
	}
	if router == nil {
		panic("router cannot be nil")
	}

	r.app.devPanel = &DevPanel{UserEntity: devPanelUserEntity, Router: router, PoolStream: poolStream}
	return r
}

func (r *Registry) RegisterDIService(service ...*ServiceDefinition) *Registry {
	r.servicesDefinitions = append(r.servicesDefinitions, service...)
	return r
}

func (r *Registry) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*ServiceDefinition{
		serviceApp(r.app),
		serviceConfig(),
	}

	flagsRegistry := &FlagsRegistry{flags: make(map[string]interface{})}
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
			r.app.scripts = append(r.app.scripts, def.Name)
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
	container = ioCBuilder.Build()
	dicInstance = &dic{}
	if !flag.Parsed() {
		flag.Parse()
	}
	r.app.flags = &Flags{flagsRegistry}
}
