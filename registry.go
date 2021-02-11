package hitrix

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sarulabs/di"
)

type Registry struct {
	app                 *AppDefinition
	servicesDefinitions []*ServiceDefinition
	scripts             []string
}

type Hitrix struct {
	registry *Registry
	ctx      context.Context
	cancel   context.CancelFunc
	done     chan bool
	exit     chan int
}

func New(appName string) *Registry {
	mode, hasMode := os.LookupEnv("hitrix_MODE")
	if !hasMode {
		mode = ModeLocal
	}
	app := &AppDefinition{mode: mode, name: appName}
	r := &Registry{app: app}
	app.registry = r
	return r
}

func (r *Registry) Build() *Hitrix {
	r.initializeIoCHandlers()
	r.initializeLog()
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
	return s
}

func (r *Registry) AttachDerPanel() *Hitrix {
	r.initializeIoCHandlers()
	r.initializeLog()
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
	return s
}

func (r *Registry) RegisterDIService(service ...*ServiceDefinition) *Registry {
	r.servicesDefinitions = append(r.servicesDefinitions, service...)
	return r
}

func (r *Registry) initializeIoCHandlers() {
	ioCBuilder, _ := di.NewBuilder()

	defaultDefinitions := []*ServiceDefinition{
		serviceApp(r.app),
		serviceLogGlobal(),
		serviceLogForRequest(),
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
			r.scripts = append(r.scripts, def.Name)
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

func (s *Hitrix) await() {
	termChan := make(chan os.Signal, 1)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	select {
	case code := <-s.exit:
		s.cancel()
		os.Exit(code)
	case <-s.done:
		s.cancel()
	case <-termChan:
		DIC().Log().Debug("TERMINATING")
		s.cancel()
		time.Sleep(time.Millisecond * 300)
		DIC().Log().Debug("TERMINATED")
	}
}
