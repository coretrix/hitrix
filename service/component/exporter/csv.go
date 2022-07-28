package exporter

import (
	"bytes"
	"encoding/csv"
	"os"
)

type ICSVExporter interface {
	exportToFile(columns []string, rows [][]interface{}, filePath string) error
	exportToByte(columns []string, rows [][]interface{}) ([]byte, error)
}

type CSVExporter struct {
}

func NewCSVExportService() *CSVExporter {
	return &CSVExporter{}
}

func (e *CSVExporter) exportToFile(columns []string, rows [][]interface{}, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err.Error())
		}
	}(f)

	writer := csv.NewWriter(f)
	defer writer.Flush()

	return e.export(writer, columns, rows)
}

func (e *CSVExporter) exportToByte(columns []string, rows [][]interface{}) ([]byte, error) {
	err := verifyRows(columns, rows)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer

	writer := csv.NewWriter(&buf)
	defer writer.Flush()

	err = e.export(writer, columns, rows)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (e *CSVExporter) export(writer *csv.Writer, columns []string, rows [][]interface{}) error {
	record := make([]string, 0)
	record = append(record, columns...)

	err := writer.Write(record)
	if err != nil {
		panic(err.Error())
	}

	record = make([]string, 0)

	for _, row := range rows {
		for columnIndex := range columns {
			record = append(record, row[columnIndex].(string))
		}

		err := writer.Write(record)
		if err != nil {
			panic(err.Error())
		}

		record = make([]string, 0)
	}

	writer.Flush()

	return nil
}
