package mocks

import (
	"github.com/stretchr/testify/mock"
)

type PDFService struct {
	mock.Mock
}

func (t *PDFService) HtmlToPdf(_ string) []byte {
	args := t.Called()
	return args.Get(0).([]byte)
}
