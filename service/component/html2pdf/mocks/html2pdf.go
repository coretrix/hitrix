package mocks

import (
	"github.com/stretchr/testify/mock"
)

type HTML2PDF struct {
	mock.Mock
}

func (t *HTML2PDF) HTMLToPdf(_ string) []byte {
	args := t.Called()

	return args.Get(0).([]byte)
}
