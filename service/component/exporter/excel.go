package exporter

import (
	"bytes"

	"github.com/tealeg/xlsx"
)

type IXLSXExporter interface {
	exportToFile(sheet string, columns []string, rows [][]interface{}, filePath string) error
	exportToByte(sheet string, columns []string, rows [][]interface{}) ([]byte, error)
}

type XLSXExporter struct {
}

func NewXLSXExportService() *XLSXExporter {
	return &XLSXExporter{}
}

func (e *XLSXExporter) exportToFile(sheet string, columns []string, rows [][]interface{}, filePath string) error {
	xlsxFile, err := e.export(sheet, columns, rows)
	if err != nil {
		return err
	}

	return xlsxFile.Save(filePath)
}

func (e *XLSXExporter) exportToByte(sheet string, columns []string, rows [][]interface{}) ([]byte, error) {
	xlsxFile, err := e.export(sheet, columns, rows)

	if err != nil {
		return nil, err
	}

	byteXLSX := new(bytes.Buffer)

	err = xlsxFile.Write(byteXLSX)

	if err != nil {
		return nil, err
	}

	return byteXLSX.Bytes(), nil
}

func (e *XLSXExporter) export(sheet string, columns []string, rows [][]interface{}) (*xlsx.File, error) {
	err := verifyRows(columns, rows)

	if err != nil {
		return nil, err
	}

	var xlsxFile *xlsx.File
	var xlsxSheet *xlsx.Sheet
	var xlsxRow *xlsx.Row
	var xlsxCell *xlsx.Cell

	xlsxFile = xlsx.NewFile()
	xlsxSheet, err = xlsxFile.AddSheet(sheet)

	if err != nil {
		return nil, err
	}

	xlsxRow = xlsxSheet.AddRow()

	for _, columnTitle := range columns {
		xlsxCell = xlsxRow.AddCell()
		xlsxCell.SetValue(columnTitle)
	}

	for _, row := range rows {
		xlsxRow = xlsxSheet.AddRow()
		for columnIndex := range columns {
			xlsxCell = xlsxRow.AddCell()
			xlsxCell.SetValue(row[columnIndex])
		}
	}

	return xlsxFile, nil
}
