package crud

import (
	"errors"
	"fmt"

	"github.com/coretrix/hitrix/pkg/dto/list"
	"github.com/coretrix/hitrix/service/component/crud"
)

func ValidateListRequest(request list.RequestDTOList, pageSizeMin, pageSizeMax int) (*crud.ListRequest, error) {
	if request.Page == nil {
		return nil, errors.New("page is empty")
	}

	if request.PageSize == nil {
		return nil, errors.New("page size is empty")
	}

	if *request.Page < 1 {
		return nil, errors.New("page must be greater than or equal to 1")
	}

	if *request.PageSize < pageSizeMin {
		return nil, fmt.Errorf("page size must be greater than or equal to %v", pageSizeMin)
	}
	if *request.PageSize > pageSizeMax {
		return nil, fmt.Errorf("page size must be greater than or equal to %v", pageSizeMax)
	}

	return &crud.ListRequest{
		Page:     request.Page,
		PageSize: request.PageSize,
		Filter:   request.Filter,
		Search:   request.Search,
		SearchOR: request.SearchOR,
		Sort:     request.Sort,
	}, nil
}
