package hitrix

import (
	"context"
	"flag"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/registry"
)

type Registry struct {
	app                        *app.App
	servicesDefinitionsGlobal  []*service.DefinitionGlobal
	servicesDefinitionsRequest []*service.DefinitionRequest
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

func (r *Registry) SetParallelTestID(parallelTestID string) *Registry {
	r.app.ParallelTestID = parallelTestID

	return r
}
func (r *Registry) Build() (*Hitrix, func()) {
	globalContext, cancel := context.WithCancel(context.Background())
	r.app.GlobalContext = globalContext
	r.app.CancelContext = cancel
	r.app.InitWaitGroup()

	r.initializeIoCHandlers()

	flags := service.DI().App().Flags
	if flags.Bool("list-scripts") {
		listScrips()
	}
	scriptToRun := flags.String("run-script")
	h := &Hitrix{done: make(chan bool), exit: make(chan int)}
	if scriptToRun != "" {
		h.runDynamicScrips(globalContext, scriptToRun)
	}
	h.startupOnBuild()

	return h, func() {
		if r := recover(); r != nil {
			service.DI().ErrorLogger().LogError(r)
		}
	}
}

func (r *Registry) RegisterDevPanel(devPanelUserEntity app.IDevPanelUserEntity, router func(ginEngine *gin.Engine)) *Registry {
	if devPanelUserEntity == nil {
		panic("devPanelUserEntity cannot be nil")
	}
	if router == nil {
		panic("router cannot be nil")
	}

	r.app.DevPanel = &app.DevPanel{UserEntity: devPanelUserEntity, Router: router}

	return r
}

func (r *Registry) RegisterRedisPools(pools *app.RedisPools) *Registry {
	r.app.RedisPools = pools

	return r
}

func (r *Registry) RegisterDIGlobalService(service ...*service.DefinitionGlobal) *Registry {
	r.servicesDefinitionsGlobal = append(r.servicesDefinitionsGlobal, service...)

	return r
}

func (r *Registry) RegisterDIRequestService(service ...*service.DefinitionRequest) *Registry {
	r.servicesDefinitionsRequest = append(r.servicesDefinitionsRequest, service...)

	return r
}

func (r *Registry) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*service.DefinitionGlobal{
		registry.ServiceProviderApp(r.app),
		registry.ServiceProviderConfig(),
		registry.ServiceProviderGoroutine(),
	}

	flagsRegistry := &app.FlagsRegistry{Flags: make(map[string]interface{})}
	for _, def := range append(defaultDefinitions, r.servicesDefinitionsGlobal...) {
		if def == nil {
			continue
		}

		if def.Script {
			r.app.Scripts = append(r.app.Scripts, def.Name)
		}

		err := ioCBuilder.Add(di.Def{
			Name:  def.Name,
			Scope: di.App,
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
		flagsRegistry.Bool("force-alters", false, "Execute all alters")
		flagsRegistry.Bool("pre-deploy", false, "Check for alters and exit")
		flagsRegistry.Bool("list-scripts", false, "list all available scripts")
		flagsRegistry.String("run-script", "", "run script")
	}

	err := ioCBuilder.Add()

	if err != nil {
		panic(err)
	}
	service.SetContainer(ioCBuilder.Build())
	service.SetRequestServices(r.servicesDefinitionsRequest)
	if !flag.Parsed() {
		flag.Parse()
	}
	r.app.Flags = &app.Flags{Registry: flagsRegistry}
}
