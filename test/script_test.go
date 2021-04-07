package main

import (
	"context"
	"testing"

	"github.com/coretrix/hitrix/service"

	"github.com/coretrix/hitrix"

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
	env := createContextMyApp(t, "my-app", nil)

	testScript2 := &testScript2{}
	env.Hitrix.RunScript(testScript2)
	assert.Equal(t, 1, testScript2.RunCounter)
	assert.Equal(t, "my-app", service.DI().App().Name)
}
