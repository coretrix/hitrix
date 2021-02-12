package main

import (
	"context"
	"time"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix"
)

type testScript struct {
	description string
	unique      bool
}

func (script *testScript) Run(_ context.Context, _ hitrix.Exit) {
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
	r.RegisterDIService(&hitrix.ServiceDefinition{
		Name:   "aa",
		Global: true,
		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &testScript{"takie tam", false}, nil
		},
	})
	r.RegisterDIService(&hitrix.ServiceDefinition{
		Name:   "bb",
		Global: true,
		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &testScript{"takie tam dwa", true}, nil
		},
	})
	r.RegisterDIService(hitrix.ServiceProviderConfigDirectory("../config"))
	r.RegisterDIService().Build()
}
