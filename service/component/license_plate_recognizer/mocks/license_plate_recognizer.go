package mocks

import "github.com/stretchr/testify/mock"

type FakeLicensePlateRecognizer struct {
	mock.Mock
}

func (f *FakeLicensePlateRecognizer) RecognizeFromImage(base64image string) ([]string, error) {
	args := f.Called(base64image)

	return args.Get(0).([]string), args.Error(1)
}
