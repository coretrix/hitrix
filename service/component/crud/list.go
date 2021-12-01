package crud

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/latolukasz/beeorm"
)

type SearchParams struct {
	Page           int
	PageSize       int
	Search         map[string]string
	SearchOR       map[string]string
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
	FilterValidMap []FilterValue
}

type FilterValue struct {
	Key   string
	Label string
}

type ListRequest struct {
	Page     *int
	PageSize *int
	Filter   map[string]interface{}
	Search   map[string]interface{}
	SearchOR map[string]interface{}
	Sort     map[string]interface{}
}

type Crud struct {
}

func (c *Crud) ExtractListParams(cols []Column, request *ListRequest) SearchParams {
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
	var stringEnumFilters = make(map[string][]FilterValue)
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
	var selectedORSearches = make(map[string]string)

mainLoop:
	for field, value := range request.Filter {
		jsonIntValue, ok := value.(json.Number)
		if ok {
			jsonInt, err := jsonIntValue.Int64()
			if err == nil {
				if helper.StringInArray(field, numberFilters...) {
					selectedNumberFilters[field] = jsonInt
					continue mainLoop
				}
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
					for _, filterValue := range stringEnumFilters[enumFiledName] {
						if filterValue.Key == stringValue {
							selectedStringFilters[field] = stringValue
							continue mainLoop
						}
					}
				}
			}
		}
	}

	for field, value := range request.Search {
		stringValue, ok := value.(string)

		if ok && len(stringValue) >= 2 {
			if helper.StringInArray(field, searchable...) {
				selectedSearches[field] = stringValue
			}
		}
	}

	for field, value := range request.SearchOR {
		stringValue, ok := value.(string)

		if ok && len(stringValue) >= 2 {
			if helper.StringInArray(field, searchable...) {
				selectedORSearches[field] = stringValue
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
		SearchOR:       selectedORSearches,
		StringFilters:  selectedStringFilters,
		NumberFilters:  selectedNumberFilters,
		BooleanFilters: selectedBooleanFilters,
		Sort:           selectedSort,
	}
}

// TODO : add full text queries when supported by hitrix
func (c *Crud) GenerateListRedisSearchQuery(params SearchParams) *beeorm.RedisSearchQuery {
	query := &beeorm.RedisSearchQuery{}
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
			"@%s:%v* ",
			field, strings.TrimSpace(beeorm.EscapeRedisSearchString(value)),
		))
	}

	orStatements := make([]string, 0)
	for field, value := range params.SearchOR {
		orStatements = append(orStatements, fmt.Sprintf(
			"(@%s:%v*)",
			field, strings.TrimSpace(beeorm.EscapeRedisSearchString(value)),
		))
	}
	if len(orStatements) > 0 {
		query.AppendQueryRaw("(" + strings.Join(orStatements, "|") + ")")
	}

	if len(params.Sort) == 1 {
		for field, mode := range params.Sort {
			query.Sort(field, !mode)
		}
	}

	return query
}

func (c *Crud) GenerateListMysqlQuery(params SearchParams) *beeorm.Where {
	where := beeorm.NewWhere("1")
	for field, value := range params.NumberFilters {
		where.Append("AND "+field+" = ?", value)
	}

	for field, value := range params.StringFilters {
		where.Append("AND "+field+" = ?", value)
	}

	for field, value := range params.BooleanFilters {
		where.Append("AND "+field+" = ?", value)
	}

	// TODO : use full text search
	for field, value := range params.Search {
		where.Append("AND "+field+" LIKE ?", value+"%")
	}

	orStatements := make([]string, 0)
	orStatementsVariables := make([]interface{}, 0)

	for field, value := range params.SearchOR {
		orStatements = append(orStatements, fmt.Sprintf(
			"%s LIKE ?",
			field,
		))
		orStatementsVariables = append(orStatementsVariables, value+"%")
	}
	if len(orStatements) > 0 {
		where.Append(
			"AND ("+strings.Join(orStatements, " OR ")+")",
			orStatementsVariables...,
		)
	}

	if len(params.Sort) == 1 {
		sortQuery := "ORDER BY "
		for field, mode := range params.Sort {
			sort := "ASC"
			if !mode {
				sort = "DESC"
			}

			sortQuery += field + " " + sort
		}

		where.Append(sortQuery)
	}

	return where
}
