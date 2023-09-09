package metrics

import (
	"context"
	"encoding/json"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/dto/metrics"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

func Get(ctx context.Context) map[string]map[string][]metrics.Row {
	query := beeorm.NewWhere("1 ORDER BY ID DESC")

	ormService := service.DI().OrmEngineForContext(ctx)
	var metricsEntities []*entity.MetricsEntity

	ormService.Search(
		query,
		beeorm.NewPager(1, 10000),
		&metricsEntities,
	)

	result := map[string]map[string][]metrics.Row{}
	//{
	//    "Memory" :  [{"Name": "admin-api", "Data": [{date, val}]}]
	// }
	for _, metricsEntity := range metricsEntities {
		memStats := &map[string]interface{}{}

		err := json.Unmarshal([]byte(metricsEntity.Metrics), memStats)
		if err != nil {
			panic(err)
		}

		for k, v := range *memStats {
			if _, ok := result[k]; !ok {
				result[k] = map[string][]metrics.Row{}
			}

			if _, ok := result[k][metricsEntity.AppName]; !ok {
				result[k][metricsEntity.AppName] = make([]metrics.Row, 0)
			} else {
				result[k][metricsEntity.AppName] = append(result[k][metricsEntity.AppName], metrics.Row{
					Value:     v,
					CreatedAt: metricsEntity.CreatedAt.UnixMilli(),
				})
			}
		}
	}

	return result
}
