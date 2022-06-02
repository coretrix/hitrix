package googleanalytics

import (
	"context"

	ga "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"

	"github.com/coretrix/hitrix/service/component/config"
)

// GoogleAnalytics4 https://developers.google.com/analytics/devguides/reporting/data/v1
type GoogleAnalytics4 struct {
	client       *ga.Service
	propertyID   string
	providerName Provider
}

func NewGA4(configFolder string, configService config.IConfig) (IProvider, error) {
	configFilePath := configFolder + "/" + configService.MustString("google_analytics.config_file_name")

	client, err := ga.NewService(context.Background(), option.WithCredentialsFile(configFilePath))
	if err != nil {
		return nil, err
	}

	return &GoogleAnalytics4{
		client:       client,
		propertyID:   configService.MustString("google_analytics.property_id"),
		providerName: GA4,
	}, nil
}

func (g *GoogleAnalytics4) GetName() Provider {
	return g.providerName
}

func (g *GoogleAnalytics4) RunReport(runReportRequest *ga.RunReportRequest) (*ga.RunReportResponse, error) {
	return g.client.Properties.RunReport("properties/"+g.propertyID, runReportRequest).Do()
}

func (g *GoogleAnalytics4) GetDimensionsAndMetrics() ([]*ga.DimensionMetadata, []*ga.MetricMetadata, error) {
	metadata, err := g.client.Properties.GetMetadata("properties/" + g.propertyID + "/metadata").Do()
	if err != nil {
		return nil, nil, err
	}

	return metadata.Dimensions, metadata.Metrics, nil
}
