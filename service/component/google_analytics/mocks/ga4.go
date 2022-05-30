package mocks

import (
	"github.com/stretchr/testify/mock"
	ga "google.golang.org/api/analyticsdata/v1beta"

	googleanalytics "github.com/coretrix/hitrix/service/component/google_analytics"
)

type FakeGoogleAnalytics4 struct {
	mock.Mock
}

func (f *FakeGoogleAnalytics4) GetName() googleanalytics.Provider {
	return f.Called().Get(0).(googleanalytics.Provider)
}

func (f *FakeGoogleAnalytics4) RunReport(runReportRequest *ga.RunReportRequest) (*ga.RunReportResponse, error) {
	args := f.Called(runReportRequest)
	return args.Get(0).(*ga.RunReportResponse), args.Error(1)
}
