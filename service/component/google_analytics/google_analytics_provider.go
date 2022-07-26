package googleanalytics

import (
	"golang.org/x/net/context"
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
	RunReport(ctx context.Context, runReportRequest *ga.RunReportRequest) (*ga.RunReportResponse, error)
	GetDimensionsAndMetrics(ctx context.Context) ([]*ga.DimensionMetadata, []*ga.MetricMetadata, error)
	GetMetrics(ctx context.Context, dateFrom, dateTo string, metrics []string, dimensions []string) (map[uint64]map[string]interface{}, error)
}
