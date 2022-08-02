package requestlogger

import (
	"context"

	"github.com/latolukasz/beeorm"
	"github.com/xorcare/pointer"

	listDto "github.com/coretrix/hitrix/pkg/dto/list"
	"github.com/coretrix/hitrix/pkg/dto/requestlogger"
	"github.com/coretrix/hitrix/pkg/entity"
	crudView "github.com/coretrix/hitrix/pkg/view/crud"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/crud"
)

const (
	pageSizeMin = 10
	pageSizeMax = 100
)

func columns() []crud.Column {
	return []crud.Column{
		{
			Key:            "UserID",
			Type:           crud.NumberType,
			Label:          "UserID",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
		{
			Key:            "URL",
			Type:           crud.StringType,
			Label:          "URL",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
		{
			Key:            "AppName",
			Type:           crud.StringType,
			Label:          "AppName",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
		{
			Key:            "Request",
			Type:           crud.StringType,
			Label:          "Request",
			Searchable:     false,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
		{
			Key:            "Response",
			Type:           crud.StringType,
			Label:          "Response",
			Searchable:     false,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
		{
			Key:            "Status",
			Type:           crud.NumberType,
			Label:          "Status",
			Searchable:     false,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
		{
			Key:            "CreatedAt",
			Type:           crud.DateTimeType,
			Label:          "CreatedAt",
			Searchable:     false,
			Sortable:       false,
			Visible:        true,
			Filterable:     false,
			FilterValidMap: nil,
		},
	}
}

func RequestsLogger(ctx context.Context, userListRequest listDto.RequestDTOList) (*requestlogger.ResponseDTORequestLoggerListDevPanel, error) {
	request, err := crudView.ValidateListRequest(userListRequest, pageSizeMin, pageSizeMax)
	if err != nil {
		return nil, err
	}

	cols := columns()
	crudService := service.DI().Crud()

	searchParams := crudService.ExtractListParams(cols, request)
	query := crudService.GenerateListRedisSearchQuery(searchParams)

	if len(searchParams.Search) == 0 && len(searchParams.NumberFilters) == 0 && len(searchParams.StringFilters) == 0 {
		query = query.Sort("ID", true)
	}

	ormService := service.DI().OrmEngineForContext(ctx)
	var requestLoggerEntities []*entity.RequestLoggerEntity

	total := ormService.RedisSearch(&requestLoggerEntities, query, beeorm.NewPager(searchParams.Page, searchParams.PageSize))

	requestLoggerEntityList := make([]*requestlogger.ResponseDTORequestLogger, len(requestLoggerEntities))

	for i, requestLoggerEntity := range requestLoggerEntities {
		requestLoggerEntityList[i] = &requestlogger.ResponseDTORequestLogger{
			ID:        requestLoggerEntity.ID,
			UserID:    requestLoggerEntity.UserID,
			URL:       requestLoggerEntity.URL,
			AppName:   requestLoggerEntity.AppName,
			Request:   requestLoggerEntity.RequestText,
			Response:  requestLoggerEntity.ResponseText,
			Log:       pointer.String(string(requestLoggerEntity.Log)),
			Status:    requestLoggerEntity.Status,
			CreatedAt: requestLoggerEntity.CreatedAt,
		}
	}

	return &requestlogger.ResponseDTORequestLoggerListDevPanel{
		Rows:    requestLoggerEntityList,
		Total:   int(total),
		Columns: cols,
	}, nil
}
