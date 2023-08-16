package requestlogger

import (
	"context"

	"github.com/latolukasz/beeorm/v2"
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

func columns() []*crud.Column {
	return []*crud.Column{
		{
			Key:                      "ID",
			FilterType:               crud.InputTypeNumber,
			Label:                    "ID",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "UserID",
			FilterType:               crud.InputTypeNumber,
			Label:                    "UserID",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "URL",
			FilterType:               crud.InputTypeString,
			Label:                    "URL",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "AppName",
			FilterType:               crud.InputTypeString,
			Label:                    "AppName",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "Request",
			FilterType:               crud.InputTypeString,
			Label:                    "Request",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "Response",
			FilterType:               crud.InputTypeString,
			Label:                    "Response",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "Status",
			FilterType:               crud.InputTypeNumber,
			Label:                    "Status",
			Searchable:               false,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "RequestDuration",
			FilterType:               crud.InputTypeNumber,
			Label:                    "RequestDuration",
			Searchable:               false,
			Sortable:                 true,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "CreatedAt",
			FilterType:               crud.DateTimePickerTypeDateTime,
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
	query := crudService.GenerateListMysqlQuery(searchParams)

	if len(searchParams.Sort) == 0 {
		query.Append("ORDER BY ID DESC")
	}

	ormService := service.DI().OrmEngineForContext(ctx)
	var requestLoggerEntities []*entity.RequestLoggerEntity

	total := ormService.SearchWithCount(
		query,
		beeorm.NewPager(searchParams.Page, searchParams.PageSize),
		&requestLoggerEntities,
	)

	requestLoggerEntityList := make([]*requestlogger.ResponseDTORequestLogger, len(requestLoggerEntities))

	for i, requestLoggerEntity := range requestLoggerEntities {
		requestLoggerEntityList[i] = &requestlogger.ResponseDTORequestLogger{
			ID:              requestLoggerEntity.ID,
			UserID:          requestLoggerEntity.UserID,
			URL:             requestLoggerEntity.URL,
			AppName:         requestLoggerEntity.AppName,
			Request:         string(requestLoggerEntity.Request),
			Response:        string(requestLoggerEntity.Response),
			Log:             pointer.String(string(requestLoggerEntity.Log)),
			Status:          requestLoggerEntity.Status,
			RequestDuration: requestLoggerEntity.RequestDuration,
			CreatedAt:       requestLoggerEntity.CreatedAt,
		}
	}

	return &requestlogger.ResponseDTORequestLoggerListDevPanel{
		Rows:    requestLoggerEntityList,
		Total:   total,
		Columns: cols,
	}, nil
}
