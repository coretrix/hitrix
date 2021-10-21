# OSS Google
This service is used for storage files into google storage

Register the service into your `main.go` file:
```go
registry.OSService(map[string]uint64{"my-bucket-name": 1})
```

Access the service:
```go
service.DI().OSService()
```

and you should register the entity `OSSBucketCounterEntity` into the ORM
You should pass parameter as a map that contains all buckets you need as a key and as a value you should pass id. This id should be unique

In your config folder you should put the .oss.json config file that you have from google
Your config file should looks like that:
```json
{
  "type": "...",
  "project_id": "...",
  "private_key_id": "...",
  "private_key": "...",
  "client_email": "...",
  "client_id": "...",
  "auth_uri": "...",
  "token_uri": "...",
  "auth_provider_x509_cert_url": "...",
  "client_x509_cert_url": "..."
}
```

The last thing you need to set in domain that gonna be used for the static files.
You can setup the domain in hitrix.yaml config file like this:
```yaml
oss: 
  domain: myapp.com
```

and the url to access your static files will looks like
`https://static-%s.myapp.com/%s/%s`
where first %s is app mode

second %s is bucket name concatenated with app mode

and last %s is the id of the file