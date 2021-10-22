package main

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/registry"
	"github.com/sarulabs/di"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestExporter(t *testing.T) {
	ioCBuilder, _ := di.NewBuilder()

	def := registry.ServiceProviderExporter()

	err := ioCBuilder.Add(di.Def{
		Name:  def.Name,
		Scope: di.App,
		Build: def.Build,
		Close: def.Close,
	})

	service.SetContainer(ioCBuilder.Build())

	exporterService := service.DI().Exporter()

	xlsxFilePath := "test_xlsx_exporter.xlsx"
	csvFilePath := "test_csv_exporter.csv"

	sheet := "Sheet Name"
	headers := []string{"Header 1", "Header 2"}
	cell1 := "cell 1"
	cell2 := "cell 2"

	rows := make([][]interface{}, 0)

	var firstRow []interface{}
	firstRow = append(firstRow, cell1, cell2)
	rows = append(rows, firstRow)

	var secondRow []interface{}
	secondRow = append(secondRow, cell1, cell2)
	rows = append(rows, secondRow)

	byteSlice, err := exporterService.XLSXExportToByte(sheet, headers, rows)
	assert.Nil(t, err)
	assert.NotNil(t, byteSlice)

	err = os.Remove(xlsxFilePath)
	assert.NotNil(t, err)

	err = exporterService.XLSXExportToFile(sheet, headers, rows, xlsxFilePath)
	assert.Nil(t, err)
	assert.FileExists(t, xlsxFilePath)

	err = os.Remove(xlsxFilePath)
	assert.Nil(t, err)

	byteSlice, err = exporterService.CSVExportToByte(headers, rows)
	assert.Nil(t, err)

	err = os.Remove(csvFilePath)
	assert.NotNil(t, err)

	err = exporterService.CSVExportToFile(headers, rows, csvFilePath)
	assert.Nil(t, err)
	assert.FileExists(t, csvFilePath)

	err = os.Remove(csvFilePath)
	assert.Nil(t, err)
}
