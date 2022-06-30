package crud

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/latolukasz/beeorm"

	"github.com/coretrix/hitrix/pkg/helper"
)

type SearchParams struct {
	Page                 int
	PageSize             int
	Search               map[string]string
	SearchOR             map[string]string
	StringFilters        map[string]string
	FormatStringFilters  map[string]string
	ArrayStringFilters   map[string][]string
	NumberFilters        map[string]int64
	ArrayNumberFilters   map[string][]int64
	RangeNumberFilters   map[string][]int64
	DateTimeFilters      map[string]time.Time
	DateFilters          map[string]time.Time
	RangeDateTimeFilters map[string][]time.Time
	RangeDateFilters     map[string][]time.Time
	BooleanFilters       map[string]bool
	Sort                 map[string]bool
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
	var formatStringSearch = make([]string, 0)
	var arrayStringFilters = make([]string, 0)
	var booleanFilters = make([]string, 0)
	var stringEnumFilters = make(map[string][]FilterValue)
	var numberFilters = make([]string, 0)
	var rangeNumberFilters = make([]string, 0)
	var arrayNumberFilters = make([]string, 0)
	var dateTimeFilters = make([]string, 0)
	var dateFilters = make([]string, 0)
	var rangeDateTimeFilters = make([]string, 0)
	var rangeDateFilters = make([]string, 0)
	var sortables = make([]string, 0)
	for _, column := range cols {
		if column.Sortable {
			sortables = append(sortables, column.Key)
		}
		if column.Searchable && column.Type == FormatStringType {
			formatStringSearch = append(formatStringSearch, column.Key)
			continue
		}
		if column.Searchable {
			searchable = append(searchable, column.Key)
			continue
		}
		if column.Filterable {
			switch column.Type {
			case StringType:
				stringFilters = append(stringFilters, column.Key)
			case ArrayStringType:
				arrayStringFilters = append(arrayStringFilters, column.Key)
			case BooleanType:
				booleanFilters = append(booleanFilters, column.Key)
			case RangeNumberType:
				rangeNumberFilters = append(rangeNumberFilters, column.Key)
			case ArrayNumberType:
				arrayNumberFilters = append(arrayNumberFilters, column.Key)
			case NumberType:
				numberFilters = append(numberFilters, column.Key)
			case EnumType:
				stringEnumFilters[column.Key] = column.FilterValidMap
			case DateTimeType:
				dateTimeFilters = append(dateTimeFilters, column.Key)
			case DateType:
				dateFilters = append(dateFilters, column.Key)
			case RangeDateTimeType:
				rangeDateTimeFilters = append(rangeDateTimeFilters, column.Key)
			case RangeDateType:
				rangeDateFilters = append(rangeDateFilters, column.Key)
			}
		}
	}

	var selectedStringFilters = make(map[string]string)
	var selectedFormatStringFilters = make(map[string]string)
	var selectedArrayStringFilters = make(map[string][]string)
	var selectedNumberFilters = make(map[string]int64)
	var selectedRangeNumberFilters = make(map[string][]int64, 2)
	var selectedArrayNumberFilters = make(map[string][]int64)
	var selectedDateTimeFilters = make(map[string]time.Time)
	var selectedDateFilters = make(map[string]time.Time)
	var selectedRangeDateTimeFilters = make(map[string][]time.Time)
	var selectedRangeDateFilters = make(map[string][]time.Time)
	var selectedBooleanFilters = make(map[string]bool)
	var selectedSort = make(map[string]bool)
	var selectedSearches = make(map[string]string)
	var selectedORSearches = make(map[string]string)

mainLoop:
	for field, value := range request.Filter {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Int64:
			if helper.StringInArray(field, numberFilters...) {
				selectedNumberFilters[field] = value.(int64)
				continue mainLoop
			}
		case reflect.Float64:
			if helper.StringInArray(field, numberFilters...) {
				selectedNumberFilters[field] = int64(value.(float64))
				continue mainLoop
			}
		case reflect.Bool:
			if helper.StringInArray(field, booleanFilters...) {
				selectedBooleanFilters[field] = value.(bool)
				continue mainLoop
			}
		case reflect.Slice:
			s := reflect.ValueOf(value)

			if helper.StringInArray(field, rangeNumberFilters...) {
				if s.Len() != 2 {
					continue mainLoop
				}

				minRange, _ := strconv.ParseInt(fmt.Sprintf("%v", s.Index(0)), 10, 64)
				maxRange, _ := strconv.ParseInt(fmt.Sprintf("%v", s.Index(1)), 10, 64)

				selectedRangeNumberFilters[field] = []int64{minRange, maxRange}
			} else if helper.StringInArray(field, rangeDateTimeFilters...) {
				if s.Len() != 2 {
					continue mainLoop
				}

				minRange, _ := time.Parse(helper.TimeLayoutRFC3339Milli, fmt.Sprintf("%v", s.Index(0)))
				maxRange, _ := time.Parse(helper.TimeLayoutRFC3339Milli, fmt.Sprintf("%v", s.Index(1)))

				selectedRangeDateTimeFilters[field] = []time.Time{minRange, maxRange}
			} else if helper.StringInArray(field, rangeDateFilters...) {
				if s.Len() != 2 {
					continue mainLoop
				}

				minRange, _ := time.Parse(helper.TimeLayoutYMD, fmt.Sprintf("%v", s.Index(0)))
				maxRange, _ := time.Parse(helper.TimeLayoutYMD, fmt.Sprintf("%v", s.Index(1)))

				selectedRangeDateFilters[field] = []time.Time{minRange, maxRange}
			} else if helper.StringInArray(field, arrayNumberFilters...) {
				if s.Len() == 0 {
					continue mainLoop
				}

				for i := 0; i < s.Len(); i++ {
					int64Value, _ := strconv.ParseInt(fmt.Sprintf("%v", s.Index(i)), 10, 64)
					selectedArrayNumberFilters[field] = append(selectedArrayNumberFilters[field], int64Value)
				}
			} else if helper.StringInArray(field, arrayStringFilters...) {
				if s.Len() == 0 {
					continue mainLoop
				}
				for i := 0; i < s.Len(); i++ {
					selectedArrayStringFilters[field] = append(selectedArrayStringFilters[field], fmt.Sprintf("%v", s.Index(i)))
				}
			}

			continue mainLoop

		case reflect.String:
			stringValue := value.(string)
			if helper.StringInArray(field, stringFilters...) {
				selectedStringFilters[field] = stringValue
				continue mainLoop
			}

			if helper.StringInArray(field, dateTimeFilters...) {
				dateTimeValue, _ := time.Parse(helper.TimeLayoutRFC3339Milli, stringValue)
				selectedDateTimeFilters[field] = dateTimeValue
				continue mainLoop
			}

			if helper.StringInArray(field, dateFilters...) {
				dateValue, _ := time.Parse(helper.TimeLayoutYMD, stringValue)
				selectedDateFilters[field] = dateValue
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

		if ok {
			if helper.StringInArray(field, formatStringSearch...) {
				selectedFormatStringFilters[field] = stringValue
				continue
			}

			if len(stringValue) >= 2 {
				if helper.StringInArray(field, searchable...) {
					selectedSearches[field] = stringValue
				}
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
			stringVal, ok := mode.(string)
			if ok && helper.StringInArray(field, sortables...) && helper.StringInArray(stringVal, "asc", "desc") {
				selectedSort[field] = stringVal == "asc"
			}
		}
	}

	return SearchParams{
		Page:                 finalPage,
		PageSize:             finalPageSize,
		Search:               selectedSearches,
		SearchOR:             selectedORSearches,
		StringFilters:        selectedStringFilters,
		FormatStringFilters:  selectedFormatStringFilters,
		ArrayStringFilters:   selectedArrayStringFilters,
		NumberFilters:        selectedNumberFilters,
		ArrayNumberFilters:   selectedArrayNumberFilters,
		RangeNumberFilters:   selectedRangeNumberFilters,
		DateTimeFilters:      selectedDateTimeFilters,
		DateFilters:          selectedDateFilters,
		RangeDateTimeFilters: selectedRangeDateTimeFilters,
		RangeDateFilters:     selectedRangeDateFilters,
		BooleanFilters:       selectedBooleanFilters,
		Sort:                 selectedSort,
	}
}

//GenerateListRedisSearchQuery TODO : add full text queries when supported by hitrix
func (c *Crud) GenerateListRedisSearchQuery(params SearchParams) *beeorm.RedisSearchQuery {
	query := &beeorm.RedisSearchQuery{}
	for field, value := range params.NumberFilters {
		query.FilterInt(field, value)
	}

	for field, value := range params.ArrayNumberFilters {
		query.FilterInt(field, value...)
	}

	for field, value := range params.RangeNumberFilters {
		query.FilterIntMinMax(field, value[0], value[1])
	}

	for field, value := range params.DateTimeFilters {
		query.FilterDateTime(field, value)
	}

	for field, value := range params.DateFilters {
		query.FilterDate(field, value)
	}

	for field, value := range params.RangeDateTimeFilters {
		query.FilterDateTimeMinMax(field, value[0], value[1])
	}

	for field, value := range params.RangeDateFilters {
		query.FilterDateMinMax(field, value[0], value[1])
	}

	for field, value := range params.StringFilters {
		query.FilterTag(field, value)
	}

	for field, value := range params.FormatStringFilters {
		query.QueryFieldPrefixMatch(field, value)
	}

	for field, value := range params.ArrayStringFilters {
		query.FilterTag(field, value...)
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
