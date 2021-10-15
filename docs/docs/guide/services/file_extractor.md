# File extractor service
File extractor provides you a simple function to search in a path recursively and find terms based on a regular expression.

Register the service into your `main.go` file:
```go
registry.ServiceProviderExtractor(),
```

Access the service:
```go
service.DI().FileExtractorService()
```
Extract phrase (errors in this example):
```go
errorTerms, err := extractService.Extract(fileextractor.ExtractParams{
  SearchPath: "./",
  Excludes:   []string{},
  Expression: `errors.New[(]*\("([^)]*)"\)`,
})
if err != nil {
  // handle error
}
```
