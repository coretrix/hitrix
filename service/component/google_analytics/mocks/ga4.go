package mocks

import (
	"context"

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

func (f *FakeGoogleAnalytics4) RunReport(ctx context.Context, runReportRequest *ga.RunReportRequest) (*ga.RunReportResponse, error) {
	args := f.Called(ctx, runReportRequest)
	return args.Get(0).(*ga.RunReportResponse), args.Error(1)
}

func (f *FakeGoogleAnalytics4) GetDimensionsAndMetrics(ctx context.Context) ([]*ga.DimensionMetadata, []*ga.MetricMetadata, error) {
	args := f.Called(ctx)
	return args.Get(0).([]*ga.DimensionMetadata), args.Get(1).([]*ga.MetricMetadata), args.Error(2)
}
