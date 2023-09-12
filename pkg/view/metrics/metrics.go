package metrics

import (
	"context"
	"encoding/json"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/dto/metrics"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func Get(ctx context.Context) map[string]metrics.Series {
	configService := service.DI().Config()

	xAxisTitle, has := configService.StringMap("metrics.xAxis.title")

	if !has {
		panic("Metrics xAxisTitle are required")
	}

	ormService := service.DI().OrmEngineForContext(ctx)

	query := beeorm.NewWhere("1 ORDER BY ID DESC")
	pager := beeorm.NewPager(1, 10000)

	var allMetricsEntities []*entity.MetricsEntity

	for {
		var metricsEntities []*entity.MetricsEntity

		ormService.Search(query, pager, &metricsEntities)

		allMetricsEntities = append(allMetricsEntities, metricsEntities...)

		if len(metricsEntities) < pager.PageSize {
			break
		}

		pager.IncrementPage()
	}

	result := map[string]metrics.Series{}
	//{
	//    "Memory" :  {"admin-api: [{date, val}]
	// }
	for _, metricsEntity := range allMetricsEntities {
		memStats := &map[string]interface{}{}

		err := json.Unmarshal([]byte(metricsEntity.Metrics), memStats)
		if err != nil {
			panic(err)
		}

		for k, v := range *memStats {
			if _, ok := result[k]; !ok {
				result[k] = metrics.Series{
					XAxisTitle: xAxisTitle[k],
					Data:       map[string][]metrics.Row{},
				}
			}

			if _, ok := result[k].Data[metricsEntity.AppName]; !ok {
				result[k].Data[metricsEntity.AppName] = make([]metrics.Row, 0)
			} else {
				result[k].Data[metricsEntity.AppName] = append(result[k].Data[metricsEntity.AppName], metrics.Row{
					Value:     v,
					CreatedAt: metricsEntity.CreatedAt.UnixMilli(),
				})
			}
		}
	}

	return result
}
