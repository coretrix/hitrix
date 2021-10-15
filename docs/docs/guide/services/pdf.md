# PDF service
PDF service provides a generating pdf function from html code using Chrome headless.

First you need these in your app config:
```yaml
chrome_headless:
  web_socket_url: ENV[CHROME_HEADLESS_WEB_SOCKET_URL]
```
Register the service into your `main.go` file:

```go
registry.ServiceProviderPDF()
```

Access the service:
```go
service.DI().PDFService()
```
Using `HtmlToPdf()` function to generate PDF from html:
```go
pdfBytes := pdfService.HtmlToPdf("<html><p>Hi!</p></html>")
```

Recommended docker file for Chrome headless:
```
https://hub.docker.com/r/chromedp/headless-shell/
```