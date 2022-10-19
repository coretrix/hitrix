package crud

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xorcare/pointer"
)

func columns() []Column {
	return []Column{
		{
			Key:            "NumberType",
			Type:           NumberType,
			Label:          "Number Type",
			Searchable:     true,
			Sortable:       true,
			Visible:        true,
			FilterValidMap: nil,
		},
		{
			Key:            "ArrayNumberType",
			Type:           ArrayNumberType,
			Label:          "ArrayNumber Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		},
		{
			Key:            "RangeNumberType",
			Type:           RangeNumberType,
			Label:          "RangeNumber Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "StringTypeFilterable",
			Type:           StringType,
			Label:          "String Type Filterable",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "StringTypeSearchable",
			Type:           StringType,
			Label:          "String Type Searchable",
			Searchable:     true,
			Sortable:       true,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "ArrayStringType",
			Type:           ArrayStringType,
			Label:          "Array String Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "EnumType",
			Type:           EnumType,
			Label:          "EnumType",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "DateTimeType",
			Type:           DateTimeType,
			Label:          "Date Time Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "DateType",
			Type:           DateType,
			Label:          "Date Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "RangeDateTimeType",
			Type:           RangeDateTimeType,
			Label:          "Range Date Time Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "RangeDateType",
			Type:           RangeDateType,
			Label:          "Range Date Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		}, {
			Key:            "BooleanType",
			Type:           BooleanType,
			Label:          "Boolean Type",
			Searchable:     true,
			Sortable:       false,
			Visible:        true,
			FilterValidMap: nil,
		},
	}
}
func TestExtractListParams(t *testing.T) {
	crud := &Crud{}

	t.Run("Filter By NumberType", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"NumberType": int64(1),
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			EnumFilters:          map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"NumberType": 1},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{},
			DateTimeFilters:      map[string]time.Time{},
			DateFilters:          map[string]time.Time{},
			RangeDateTimeFilters: map[string][]time.Time{},
			RangeDateFilters:     map[string][]time.Time{},
			BooleanFilters:       map[string]bool{},
			Sort:                 map[string]bool{},
		}
		assert.Equal(t, expected, searchParam)
	})

	t.Run("Filter By NumberType - float", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"NumberType": float64(1),
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			EnumFilters:          map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"NumberType": 1},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{},
			DateTimeFilters:      map[string]time.Time{},
			DateFilters:          map[string]time.Time{},
			RangeDateTimeFilters: map[string][]time.Time{},
			RangeDateFilters:     map[string][]time.Time{},
			BooleanFilters:       map[string]bool{},
			Sort:                 map[string]bool{},
		}
		assert.Equal(t, expected, searchParam)
	})
	t.Run("Filter By NumberType - float", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"ArrayNumberType": []int64{1, 100, 112},
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			EnumFilters:          map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{},
			ArrayNumberFilters:   map[string][]int64{"ArrayNumberType": {1, 100, 112}},
			RangeNumberFilters:   map[string][]int64{},
			DateTimeFilters:      map[string]time.Time{},
			DateFilters:          map[string]time.Time{},
			RangeDateTimeFilters: map[string][]time.Time{},
			RangeDateFilters:     map[string][]time.Time{},
			BooleanFilters:       map[string]bool{},
			Sort:                 map[string]bool{},
		}
		assert.Equal(t, expected, searchParam)
	})
	t.Run("Filter By Range number", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"RangeNumberType": []int64{1, 100},
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			EnumFilters:          map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{"RangeNumberType": {1, 100}},
			DateTimeFilters:      map[string]time.Time{},
			DateFilters:          map[string]time.Time{},
			RangeDateTimeFilters: map[string][]time.Time{},
			RangeDateFilters:     map[string][]time.Time{},
			BooleanFilters:       map[string]bool{},
			Sort:                 map[string]bool{},
		}
		assert.Equal(t, expected, searchParam)
	})
	t.Run("Filter By string", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"RangeNumberType": []int64{1, 100},
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			EnumFilters:          map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{"RangeNumberType": {1, 100}},
			DateTimeFilters:      map[string]time.Time{},
			DateFilters:          map[string]time.Time{},
			RangeDateTimeFilters: map[string][]time.Time{},
			RangeDateFilters:     map[string][]time.Time{},
			BooleanFilters:       map[string]bool{},
			Sort:                 map[string]bool{},
		}
		assert.Equal(t, expected, searchParam)
	})
	t.Run("Filter By json.Number", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"NumberType": json.Number("1"),
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			EnumFilters:          map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"NumberType": int64(1)},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{},
			DateTimeFilters:      map[string]time.Time{},
			DateFilters:          map[string]time.Time{},
			RangeDateTimeFilters: map[string][]time.Time{},
			RangeDateFilters:     map[string][]time.Time{},
			BooleanFilters:       map[string]bool{},
			Sort:                 map[string]bool{},
		}
		assert.Equal(t, expected, searchParam)
	})
}
