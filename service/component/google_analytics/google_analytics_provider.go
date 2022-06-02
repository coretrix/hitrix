package googleanalytics

import (
	ga "google.golang.org/api/analyticsdata/v1beta"
)

type Provider string

const (
	UA  Provider = "ua" // TODO: Implement later
	GA4 Provider = "ga4"
)

func (p Provider) String() string {
	return string(p)
}

type IProvider interface {
	GetName() Provider
	RunReport(runReportRequest *ga.RunReportRequest) (*ga.RunReportResponse, error)
	GetDimensionsAndMetrics() ([]*ga.DimensionMetadata, []*ga.MetricMetadata, error)
}
