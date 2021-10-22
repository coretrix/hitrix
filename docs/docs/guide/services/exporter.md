# Exporter service
This service is able to export business data to various file formats. 
Currently, we support 2 file formats: `XLSX` and `CSV`.

Register the service into your `main.go` file:

```go
registry.ServiceProviderExporter()
```

Access the service:

```go
service.DI().Exporter()
```

Input data should be filled in as follows:

```go
sheet := "sheet 1"
headers := []string{"Header 1", "Header 2"}

rows := make([][]interface{}, 0)

var firstRow []interface{}
firstRow = append(firstRow, "cell 1", "cell 2")
rows = append(rows, firstRow)

var secondRow []interface{}
secondRow = append(secondRow, "cell 1", "cell 2")
rows = append(rows, secondRow)
```

Use `XLSXExportToByte()` function to convert raw data to Excel file and return it as a byte slice:
```go
xlsxBytes, err := exporterService.XLSXExportToByte(sheet, headers, rows)
```

Use `XLSXExportToFile()` function for converting raw data to Excel file and save it in the given path:
```go
err := exporterService.XLSXExportToFile(sheet, headers, rows, filePath)
```

Use `CSVExportToByte()` function to convert raw data to CSV  file and return it as a byte slice:
```go
csvBytes, err := exporterService.CSVExportToByte(headers, rows)
```

Using `CSVExportToFile()` function for converting raw data to CSV file and save it in the given path:
```go
err := exporterService.XLSXExportToFile(headers, rows, filePath)
```