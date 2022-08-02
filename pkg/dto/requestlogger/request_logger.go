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
	ID        uint64
	UserID    uint64
	URL       string
	AppName   string
	Request   string
	Response  string
	Log       *string
	Status    int
	CreatedAt time.Time
}
