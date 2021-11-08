package main

import (
	"context"
	"time"

	"github.com/coretrix/hitrix/service/component/app"

	"github.com/coretrix/hitrix/service/registry"

	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix"
)

type testScript struct {
	description string
	unique      bool
}

func (script *testScript) Run(_ context.Context, _ app.IExit) {
}

func (script *testScript) Unique() bool {
	return script.unique
}

func (script *testScript) Description() string {
	return script.description
}

func (script *testScript) Active() bool {
	return true
}

func (script *testScript) Interval() time.Duration {
	return 3 * time.Second
}

func main() {
	r := hitrix.New("test_script", "secret")
	r.RegisterDIGlobalService(&service.DefinitionGlobal{
		Name: "aa",

		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &testScript{"takie tam", false}, nil
		},
	})
	r.RegisterDIGlobalService(&service.DefinitionGlobal{
		Name: "bb",

		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &testScript{"takie tam dwa", true}, nil
		},
	})
	r.RegisterDIGlobalService(registry.ServiceProviderConfigDirectory("../config"))
	r.RegisterDIGlobalService().Build()
}
