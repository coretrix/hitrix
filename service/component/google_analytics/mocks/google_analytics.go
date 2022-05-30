package mocks

import (
	"github.com/stretchr/testify/mock"

	googleanalytics "github.com/coretrix/hitrix/service/component/google_analytics"
)

type FakeGoogleAnalytics struct {
	mock.Mock
}

func (f *FakeGoogleAnalytics) GetProvider(_ googleanalytics.Provider) googleanalytics.IProvider {
	return f.Called().Get(0).(googleanalytics.IProvider)
}
