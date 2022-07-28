package exporter

import (
	"errors"
	"strconv"
	"strings"
)

type IExporter interface {
	XLSXExportToFile(sheet string, columns []string, rows [][]interface{}, filePath string) error
	XLSXExportToByte(sheet string, columns []string, rows [][]interface{}) ([]byte, error)
	CSVExportToFile(columns []string, rows [][]interface{}, filePath string) error
	CSVExportToByte(columns []string, rows [][]interface{}) ([]byte, error)
}

type Exporter struct {
	xlsxExporter IXLSXExporter
	csvExporter  ICSVExporter
}

func NewExportService(xlsxExporter IXLSXExporter, csvExporter ICSVExporter) *Exporter {
	return &Exporter{
		xlsxExporter: xlsxExporter,
		csvExporter:  csvExporter,
	}
}

func (e *Exporter) XLSXExportToFile(sheet string, columns []string, rows [][]interface{}, filePath string) error {
	return e.xlsxExporter.exportToFile(sheet, columns, rows, filePath)
}

func (e *Exporter) XLSXExportToByte(sheet string, columns []string, rows [][]interface{}) ([]byte, error) {
	return e.xlsxExporter.exportToByte(sheet, columns, rows)
}

func (e *Exporter) CSVExportToFile(columns []string, rows [][]interface{}, filePath string) error {
	return e.csvExporter.exportToFile(columns, rows, filePath)
}

func (e *Exporter) CSVExportToByte(columns []string, rows [][]interface{}) ([]byte, error) {
	return e.csvExporter.exportToByte(columns, rows)
}

func verifyRows(columns []string, rows [][]interface{}) error {
	dataErrors := make([]string, 0, 1)

	for rowID, row := range rows {
		if len(row) != len(columns) {
			dataErrors = append(dataErrors, "Different column count for row["+strconv.Itoa(rowID)+"]")

			continue
		}
	}

	if len(dataErrors) > 0 {
		return errors.New(strings.Join(dataErrors, "\n"))
	}

	return nil
}
