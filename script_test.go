package hitrix

import (
	"context"
	"testing"

	"github.com/sarulabs/di"
	"github.com/tj/assert"
)

type testScript struct {
	RunCounter int
}

func (script *testScript) Run(_ context.Context, _ Exit) {
	script.RunCounter++
}

func (script *testScript) Unique() bool {
	return false
}

func (script *testScript) Code() string {
	return "test script"
}

func (script *testScript) Description() string {
	return "test description"
}

func TestRunScript(t *testing.T) {
	r := New("test_script").RegisterDIService()
	testService := &ServiceDefinition{
		Name:   "test_service",
		Global: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return "hello", nil
		},
	}
	s := r.RegisterDIService(testService).Build()

	testScript := &testScript{}
	s.RunScript(testScript)
	assert.Equal(t, 1, testScript.RunCounter)
	assert.Equal(t, "test_script", DIC().App().Name())
	assert.Equal(t, "hello", GetServiceRequired("test_service"))
}
