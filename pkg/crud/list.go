package crud

import (
	"fmt"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/latolukasz/orm"
)

type SearchParams struct {
	Page           int
	PageSize       int
	Search         map[string]string
	StringFilters  map[string]string
	NumberFilters  map[string]int64
	BooleanFilters map[string]bool
	Sort           map[string]bool
}

type Column struct {
	Key            string
	Label          string
	Type           string
	Searchable     bool
	Sortable       bool
	Filterable     bool
	Visible        bool
	FilterValidMap []string
}

type ListRequest struct {
	Page     *int
	PageSize *int
	Filter   map[string]interface{}
	Search   map[string]interface{}
	Sort     map[string]interface{}
}

func ExtractListParams(cols []Column, request *ListRequest) SearchParams {
	finalPage := 1
	finalPageSize := 10
	if request.Page != nil && *request.Page > 0 {
		finalPage = *request.Page
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		finalPageSize = *request.PageSize
	}

	var searchable = make([]string, 0)
	var stringFilters = make([]string, 0)
	var booleanFilters = make([]string, 0)
	var stringEnumFilters = make(map[string][]string)
	var numberFilters = make([]string, 0)
	var sortables = make([]string, 0)
	for i := range cols {
		if cols[i].Sortable {
			sortables = append(sortables, cols[i].Key)
		}
		if cols[i].Searchable {
			searchable = append(searchable, cols[i].Key)
			continue
		}
		if cols[i].Filterable && cols[i].Type == StringType {
			stringFilters = append(stringFilters, cols[i].Key)
			continue
		}
		if cols[i].Filterable && cols[i].Type == BooleanType {
			booleanFilters = append(booleanFilters, cols[i].Key)
			continue
		}
		if cols[i].Filterable && cols[i].Type == NumberType {
			numberFilters = append(numberFilters, cols[i].Key)
			continue
		}
		if cols[i].Filterable && cols[i].Type == EnumType {
			stringEnumFilters[cols[i].Key] = cols[i].FilterValidMap
			continue
		}
	}

	var selectedStringFilters = make(map[string]string)
	var selectedNumberFilters = make(map[string]int64)
	var selectedBooleanFilters = make(map[string]bool)
	var selectedSort = make(map[string]bool)
	var selectedSearches = make(map[string]string)

mainLoop:
	for field, value := range request.Filter {
		intValue, ok := value.(int64)
		if ok {
			if helper.StringInArray(field, numberFilters...) {
				selectedNumberFilters[field] = intValue
				continue mainLoop
			}
		}

		booleanValue, ok := value.(bool)
		if ok {
			if helper.StringInArray(field, booleanFilters...) {
				selectedBooleanFilters[field] = booleanValue
				continue mainLoop
			}
		}

		stringValue, ok := value.(string)
		if ok {
			if helper.StringInArray(field, stringFilters...) {
				selectedStringFilters[field] = stringValue
				continue mainLoop
			}

			for enumFiledName := range stringEnumFilters {
				if field == enumFiledName {
					if helper.StringInArray(stringValue, stringEnumFilters[enumFiledName]...) {
						selectedStringFilters[field] = stringValue
						continue mainLoop
					}
				}
			}
		}
	}

	for field, value := range request.Search {
		stringValue, ok := value.(string)

		if ok && len(stringValue) >= 2 {
			if helper.StringInArray(field, searchable...) {
				selectedSearches[field] = stringValue + "*"
			}
		}
	}

	if len(request.Sort) == 1 {
		for field, mode := range request.Sort {
			stringVal := mode.(string)
			if helper.StringInArray(field, sortables...) && helper.StringInArray(stringVal, "asc", "desc") {
				selectedSort[field] = stringVal == "asc"
			}
		}
	}

	return SearchParams{
		Page:           finalPage,
		PageSize:       finalPageSize,
		Search:         selectedSearches,
		StringFilters:  selectedStringFilters,
		NumberFilters:  selectedNumberFilters,
		BooleanFilters: selectedBooleanFilters,
		Sort:           selectedSort,
	}
}

// TODO : add full text queries when supported by hitrix
func GenerateListRedisSearchQuery(params SearchParams) *orm.RedisSearchQuery {
	query := &orm.RedisSearchQuery{}
	for field, value := range params.NumberFilters {
		query.FilterInt(field, value)
	}

	for field, value := range params.StringFilters {
		query.FilterTag(field, value)
	}

	for field, value := range params.BooleanFilters {
		query.FilterBool(field, value)
	}

	// TODO : use full text search
	for field, value := range params.Search {
		query.QueryRaw(fmt.Sprintf(
			"@%s:%v",
			field, value,
		))
	}

	if len(params.Sort) == 1 {
		for field, mode := range params.Sort {
			query.Sort(field, !mode)
		}
	}

	return query
}
