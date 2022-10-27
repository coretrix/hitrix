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
			Key:                      "ID",
			Type:                     crud.InputTypeNumber,
			Label:                    "ID",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "UserID",
			Type:                     crud.InputTypeNumber,
			Label:                    "UserID",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "URL",
			Type:                     crud.InputTypeString,
			Label:                    "URL",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "AppName",
			Type:                     crud.InputTypeString,
			Label:                    "AppName",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "Request",
			Type:                     crud.InputTypeString,
			Label:                    "Request",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "Response",
			Type:                     crud.InputTypeString,
			Label:                    "Response",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "Status",
			Type:                     crud.InputTypeNumber,
			Label:                    "Status",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "CreatedAt",
			Type:                     crud.DateTimePickerTypeDateTime,
			Label:                    "CreatedAt",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
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

	if len(searchParams.NumberFilters) == 0 && len(searchParams.TagFilters) == 0 {
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
