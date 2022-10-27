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
			Key:                      "InputTypeNumber",
			Type:                     InputTypeNumber,
			Label:                    "Number Type",
			Searchable:               true,
			Sortable:                 true,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "MultiSelectTypeArrayNumber",
			Type:                     MultiSelectTypeArrayNumber,
			Label:                    "ArrayNumber Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
		{
			Key:                      "RangeSliderTypeArrayNumber",
			Type:                     RangeSliderTypeArrayNumber,
			Label:                    "RangeNumber Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "StringTypeFilterable",
			Type:                     InputTypeString,
			Label:                    "String Type Filterable",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "StringTypeSearchable",
			Type:                     InputTypeString,
			Label:                    "String Type Searchable",
			Searchable:               true,
			Sortable:                 true,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "ArrayStringType",
			Type:                     ArrayStringType,
			Label:                    "Array String Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:        "SelectTypeStringString",
			Type:       SelectTypeStringString,
			Label:      "SelectTypeStringString",
			Searchable: true,
			Sortable:   false,
			Visible:    true,
			DataStringKeyStringValue: []StringKeyStringValue{{
				Key:   "active",
				Label: "Active",
			}},
		}, {
			Key:        "SelectTypeIntString",
			Type:       SelectTypeIntString,
			Label:      "SelectTypeIntString",
			Searchable: true,
			Sortable:   false,
			Visible:    true,
			DataIntKeyStringValue: []IntKeyStringValue{{
				Key:   1,
				Label: "Sofia",
			}},
		}, {
			Key:                      "DateTimePickerTypeDateTime",
			Type:                     DateTimePickerTypeDateTime,
			Label:                    "Date Time Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "DatePickerTypeDate",
			Type:                     DatePickerTypeDate,
			Label:                    "Date Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "RangeDateTimePickerTypeArrayDateTime",
			Type:                     RangeDateTimePickerTypeArrayDateTime,
			Label:                    "Range Date Time Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "RangeDatePickerTypeArrayDate",
			Type:                     RangeDatePickerTypeArrayDate,
			Label:                    "Range Date Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		}, {
			Key:                      "CheckboxTypeBoolean",
			Type:                     CheckboxTypeBoolean,
			Label:                    "Boolean Type",
			Searchable:               true,
			Sortable:                 false,
			Visible:                  true,
			DataStringKeyStringValue: nil,
		},
	}
}
func TestExtractListParams(t *testing.T) {
	crud := &Crud{}

	t.Run("Filter By SelectTypeIntString", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"SelectTypeIntString":    int64(1),
				"SelectTypeStringString": "active",
			},
		})
		expected := SearchParams{
			Page:                 1,
			PageSize:             15,
			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{"SelectTypeStringString": "active"},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"SelectTypeIntString": 1},
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

	t.Run("Filter By InputTypeNumber", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"InputTypeNumber": int64(1),
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"InputTypeNumber": 1},
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

	t.Run("Filter By InputTypeNumber - float", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"InputTypeNumber": float64(1),
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"InputTypeNumber": 1},
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
	t.Run("Filter By InputTypeNumber - float", func(t *testing.T) {
		searchParam := crud.ExtractListParams(columns(), &ListRequest{
			Page:     pointer.Int(1),
			PageSize: pointer.Int(15),
			Search: map[string]interface{}{
				"MultiSelectTypeArrayNumber": []int64{1, 100, 112},
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{},
			ArrayNumberFilters:   map[string][]int64{"MultiSelectTypeArrayNumber": {1, 100, 112}},
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
				"RangeSliderTypeArrayNumber": []int64{1, 100},
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{"RangeSliderTypeArrayNumber": {1, 100}},
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
				"RangeSliderTypeArrayNumber": []int64{1, 100},
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{},
			ArrayNumberFilters:   map[string][]int64{},
			RangeNumberFilters:   map[string][]int64{"RangeSliderTypeArrayNumber": {1, 100}},
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
				"InputTypeNumber": json.Number("1"),
			},
		})
		expected := SearchParams{
			Page:     1,
			PageSize: 15,

			StringORFilters:      map[string]string{},
			TagFilters:           map[string]string{},
			StringFilters:        map[string]string{},
			ArrayStringFilters:   map[string][]string{},
			NumberFilters:        map[string]int64{"InputTypeNumber": int64(1)},
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
