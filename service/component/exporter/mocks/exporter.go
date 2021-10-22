package mocks

import "github.com/stretchr/testify/mock"

type FakeExporter struct {
	mock.Mock
}

func (e *FakeExporter) XLSXExportToFile(sheet string, columns []string, rows [][]interface{}, filePath string) error {
	return e.Called(sheet, columns, rows, filePath).Get(0).(error)
}

func (e *FakeExporter) XLSXExportToByte(sheet string, columns []string, rows [][]interface{}) ([]byte, error) {
	args := e.Called(sheet, columns, rows)
	return args.Get(0).([]byte), args.Error(1)
}

func (e *FakeExporter) CSVExportToFile(columns []string, rows [][]interface{}, filePath string) error {
	return e.Called(columns, rows, filePath).Get(0).(error)
}

func (e *FakeExporter) CSVExportToByte(columns []string, rows [][]interface{}) ([]byte, error) {
	args := e.Called(columns, rows)
	return args.Get(0).([]byte), args.Error(1)
}
