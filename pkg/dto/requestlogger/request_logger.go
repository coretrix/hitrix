package requestlogger

import (
	"time"

	"github.com/coretrix/hitrix/service/component/crud"
)

type ResponseDTORequestLoggerListDevPanel struct {
	Rows    []*ResponseDTORequestLogger
	Total   int
	Columns []crud.Column
}

type ResponseDTORequestLogger struct {
	UserID          uint64
	URL             string
	AppName         string
	Content         string
	ResponseContent string
	Text            string
	ResponseText    string
	Log             *string
	Status          int
	CreatedAt       time.Time
}
