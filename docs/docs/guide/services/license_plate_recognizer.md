# License plate recognizer
This service allow us to upload base64, and it will return all car license plates on this image.

You should put your api key from `platerecognizer.com` in config:

```yml
platerecognizer:
  api_key: some_key
```

Access the service:
```go
service.DI().LicensePlateRecognizer()
```
