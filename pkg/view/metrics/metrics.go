package metrics

import (
	"context"
	"github.com/coretrix/hitrix/pkg/dto/metrics"
	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
)

const (
	pageSizeMin = 10
	pageSizeMax = 100
)

func List(ctx context.Context) (*metrics.ResponseDTORMetrics, error) {
	query := beeorm.NewWhere("ORDER BY ID DESC")

	ormService := service.DI().OrmEngineForContext(ctx)
	var metricsEntities []*entity.MetricsEntity

	ormService.Search(
		query,
		beeorm.NewPager(1, 10000),
		&metricsEntities,
	)

	rows := make([]*metrics.Row, len(metricsEntities))

	for i, metricsEntity := range metricsEntities {
		rows[i] = &metrics.Row{
			ID:        metricsEntity.ID,
			CreatedAt: helper.GetTimestamp(&metricsEntity.CreatedAt),
		}
	}

	return &metrics.ResponseDTORMetrics{
		Rows: rows,
	}, nil
}
