package main

import (
	"context"
	"testing"

	"github.com/latolukasz/beeorm"
	"github.com/stretchr/testify/assert"

	"github.com/coretrix/hitrix"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type testScript2 struct {
	RunCounter int
}

func (script *testScript2) Run(_ context.Context, _ app.IExit, _ *beeorm.Engine) {
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
	env := createContextMyApp(t, "my-app", nil, nil, nil)

	testScript2 := &testScript2{}
	b := &hitrix.BackgroundProcessor{Server: env.Hitrix}
	b.RunScript(testScript2)

	assert.Equal(t, 0, testScript2.RunCounter)
	assert.Equal(t, "my-app", service.DI().App().Name)
}
