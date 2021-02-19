package main

import (
	"context"
	"testing"

	"github.com/coretrix/hitrix/service"

	"github.com/coretrix/hitrix"

	"github.com/sarulabs/di"
	"github.com/tj/assert"
)

type testScript2 struct {
	RunCounter int
}

func (script *testScript2) Run(_ context.Context, _ hitrix.Exit) {
	script.RunCounter++
}

func (script *testScript2) Unique() bool {
	return false
}

func (script *testScript2) Code() string {
	return "test script"
}

func (script *testScript2) Description() string {
	return "test description"
}

func TestRunScript(t *testing.T) {
	r := hitrix.New("test_script", "secret").RegisterDIService()
	testService := &service.Definition{
		Name:   "test_service",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return "hello", nil
		},
	}
	s, deferFunc := r.RegisterDIService(testService).Build()
	defer deferFunc()

	testScript2 := &testScript2{}
	s.RunScript(testScript2)
	assert.Equal(t, 1, testScript2.RunCounter)
	assert.Equal(t, "test_script", service.DI().App().Name)
	assert.Equal(t, "hello", service.GetServiceRequired("test_service"))
}
