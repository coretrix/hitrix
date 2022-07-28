package mocks

import (
	"github.com/stretchr/testify/mock"
)

type FakeTemplateService struct {
	mock.Mock
}

func (t *FakeTemplateService) RenderTemplate(html string, data interface{}) (string, error) {
	args := t.Called(html, data)

	err := args.Get(1)
	if err == nil {
		return args.Get(0).(string), nil
	}

	return args.Get(0).(string), err.(error)
}

func (t *FakeTemplateService) RenderMandrillTemplate(template string, data interface{}) (string, error) {
	args := t.Called(template, data)

	err := args.Get(1)
	if err == nil {
		return args.Get(0).(string), nil
	}

	return args.Get(0).(string), err.(error)
}
