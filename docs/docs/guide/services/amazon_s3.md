# Amazon S3

This service is used for storing files into amazon s3

Register the service into your `main.go` file:

```go
registry.ServiceProviderAmazonS3(map[string]uint64{"products": 1}) // 1 is the bucket ID for database counter
```

and you should register the entity `S3BucketCounterEntity` into the ORM
Also, you should put your credentials and other configs in `config/hitrix.yml`

```yml
amazon_s3:
  endpoint: "https://somestorage.com" # set to "" if you're using https://s3.amazonaws.com
  access_key_id: ENV[S3_ACCESS_KEY_ID]
  secret_access_key: ENV[S3_SECRET_ACCESS_KEY_ID]
  disable_ssl: false
  region: us-east-1
  url_prefix: prefix
  domain: domain.com
  buckets: # Register your buckets here for each app mode
    products: # bucket name
      prod: bucket-name
      local: bucket-name-local
  public_urls: # Register your public urls for the GetObjectCachedURL method
    product: # bucket name
      prod: "https://somesite.com/{{.StorageKey}}/" # Available variables are: .Environment, .BucketName, .CounterID, and, .StorageKey
      local: "http://127.0.0.1/{{.Environment}}/{{.BucketName}}/{{.StorageKey}}/{{.CounterID}}" # Will output "http://127.0.0.1/local/product/1.jpeg/1"
```

Access the service:
```go
service.DI().AmazonS3()
```