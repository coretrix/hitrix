package googleanalytics

import (
	"context"
	"fmt"
	"strconv"

	ga "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"

	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

// GoogleAnalytics4 https://developers.google.com/analytics/devguides/reporting/data/v1
type GoogleAnalytics4 struct {
	client       *ga.Service
	propertyID   string
	providerName Provider
	errorLogger  errorlogger.ErrorLogger
}

func NewGA4(configFolder string, configService config.IConfig, errorlogger errorlogger.ErrorLogger) (IProvider, error) {
	configFilePath := configFolder + "/" + configService.MustString("google_analytics.config_file_name")

	client, err := ga.NewService(context.Background(), option.WithCredentialsFile(configFilePath))
	if err != nil {
		return nil, err
	}

	return &GoogleAnalytics4{
		client:       client,
		propertyID:   configService.MustString("google_analytics.property_id"),
		providerName: GA4,
		errorLogger:  errorlogger,
	}, nil
}

func (g *GoogleAnalytics4) GetName() Provider {
	return g.providerName
}

func (g *GoogleAnalytics4) RunReport(ctx context.Context, runReportRequest *ga.RunReportRequest) (*ga.RunReportResponse, error) {
	return g.client.Properties.RunReport("properties/"+g.propertyID, runReportRequest).Context(ctx).Do()
}

func (g *GoogleAnalytics4) GetDimensionsAndMetrics(ctx context.Context) ([]*ga.DimensionMetadata, []*ga.MetricMetadata, error) {
	metadata, err := g.client.Properties.GetMetadata("properties/" + g.propertyID + "/metadata").Context(ctx).Do()
	if err != nil {
		return nil, nil, err
	}

	return metadata.Dimensions, metadata.Metrics, nil
}

func (g *GoogleAnalytics4) GetMetrics(ctx context.Context, dateFrom, dateTo string, metrics []string, dimensions []string) (map[uint64]map[string]interface{}, error) {
	offset := int64(0)
	headers := make([]string, 0)
	types := make([]string, 0)
	rows := make([]*ga.Row, 0)

	runReportRequest := &ga.RunReportRequest{
		DateRanges: []*ga.DateRange{
			{
				StartDate: dateFrom,
				EndDate:   dateTo,
			},
		},
		Metrics:    []*ga.Metric{},
		Dimensions: []*ga.Dimension{},
		Offset:     offset,
	}

	for _, metric := range metrics {
		runReportRequest.Metrics = append(runReportRequest.Metrics, &ga.Metric{
			Name: metric,
		})
	}

	for _, dimension := range dimensions {
		runReportRequest.Dimensions = append(runReportRequest.Dimensions, &ga.Dimension{
			Name: dimension,
		})
	}

	for {
		runReportRequest.Offset = offset

		resp, err := g.RunReport(ctx, runReportRequest)
		if err != nil {
			return nil, err
		}

		if cap(rows) == 0 {
			headers = make([]string, len(resp.MetricHeaders))
			types = make([]string, len(resp.MetricHeaders))

			for i, header := range resp.MetricHeaders {
				headers[i] = header.Name
				types[i] = header.Type
			}

			rows = make([]*ga.Row, 0, resp.RowCount)
		}

		offset += int64(len(resp.Rows))
		rows = append(rows, resp.Rows...)

		if offset >= resp.RowCount {
			break
		}
	}

	result := make(map[uint64]map[string]interface{})

	for _, row := range rows {
		if len(row.DimensionValues) < 1 {
			continue
		}

		dimension := row.DimensionValues[0].Value
		if dimension == "(not set)" {
			continue
		}

		dimensionResourceID, err := strconv.ParseUint(dimension, 10, 64)
		if err != nil {
			g.errorLogger.LogError(fmt.Sprintf("could not parse %q as dimensionResourceID (Uint64) : %v", dimension, err))

			continue
		}

		data := make(map[string]interface{})

		for i, metric := range row.MetricValues {
			var castErr error
			var value interface{}

			switch types[i] {
			case "TYPE_FLOAT":
				value, castErr = strconv.ParseFloat(metric.Value, 64)
			default:
				value, castErr = strconv.ParseInt(metric.Value, 10, 64)
			}

			if castErr != nil {
				g.errorLogger.LogError(fmt.Sprintf("could not parse %q as %s (%s) : %v", metric.Value, headers[i], types[i], err))

				continue
			}

			data[headers[i]] = value
		}

		result[dimensionResourceID] = data
	}

	return result, nil
}
