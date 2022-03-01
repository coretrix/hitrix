# Object Storage Service

This service is used for storing files into any amazon s3 or google cloud compatible services

Register the service into your `main.go`. You need to provide function that init the provider and list of public/private buckets

```go
    registry.ServiceProviderOSS(oss.NewAmazonOSS, []oss.Namespace{"product_images"}, []oss.Namespace{"invoices"}),
```

and you should register the entity `OSSBucketCounterEntity` into the ORM

Also, you should put your credentials and other configs in `config/hitrix.yml`

S3 example:
```yml
oss:
  amazon_s3:
    endpoint: "https://somestorage.com" # set to "" if you're using https://s3.amazonaws.com
    access_key_id: ENV[S3_ACCESS_KEY_ID]
    secret_access_key: ENV[S3_SECRET_ACCESS_KEY_ID]
    disable_ssl: false
    region: us-east-1
  buckets:
    public: # config for public bucket
     name: bucket-name # bucket name
     cdn_url: "https://somesite.com/{{.StorageKey}}/" # Available variables is: .StorageKey (Namespace is part of StorageKey)
    private: # config for private bucket
     name: bucket-name-private # bucket name
     local: "http://127.0.0.1/{{.StorageKey}}" # Will output "http://127.0.0.1/product/1.jpeg"
```
Google example:

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
  google:
    anyvar: 1
  buckets:
    public: # config for public bucket
      name: bucket-name # bucket name
      cdn_url: "https://somesite.com/{{.StorageKey}}/" # Available variables are: .Namespace, .CounterID, and, .StorageKey
    private: # config for private bucket
      name: bucket-name-private # bucket name
      local: "http://127.0.0.1/{{.Namespace}}/{{.StorageKey}}/{{.CounterID}}" # Will output "http://127.0.0.1/product/1.jpeg/1"
```

Access the service:
```go
service.DI().OSService()
```