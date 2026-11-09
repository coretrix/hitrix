# Object Storage Service

This service is used for storing files into any amazon s3 or google cloud compatible services

Register the service into your `main.go`. You need to provide function that init the provider and list of  buckets

```go
registry.ServiceProviderOSS(oss.NewAmazonOSS, oss.Namespaces{"products": oss.BucketPublic})
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

## Direct browser uploads (Google Cloud Storage)

GoogleOSS implements the optional `IResumableUploadProvider` interface. Existing
`IProvider` methods and other storage providers are unchanged. Check support first:

```go
uploads, supported := service.DI().OSService().(oss.IResumableUploadProvider)
// Handle !supported before calling uploads.
upload, err := uploads.StartResumableUpload(ctx, ormService, namespace, oss.ResumableUploadOptions{
    Extension: ".apk", ContentType: "application/vnd.android.package-archive",
    Size: fileSize, Origin: allowedBrowserOrigin, CRC32C: fileCRC32C,
})
// Handle err, then save upload server-side with its owner, namespace and expiry.
// Return only your pending-upload ID and upload.UploadURL to the browser.
// After upload, load the saved record using that ID and check owner and expiry:
fileObject, err := uploads.CompleteResumableUpload(ctx, namespace, upload)
// Handle err, then publish fileObject exactly once in your application.
```

- Authorize the user, namespace, browser origin and size before initiation. CRC32C
  is the Castagnoli checksum encoded as four big-endian bytes in standard base64;
  it verifies transfer integrity, not APK validity or publisher identity.
- Initiation reserves the normal Hitrix file ID/key and prevents overwrites.
  File bytes go directly to GCS. Failed attempts can leave gaps in file IDs.
- Upload with `PUT` to the returned URL using a client without your API's auth
  headers/cookies. For chunked uploads, use `Content-Range`; non-final chunks must
  be multiples of 256 KiB. An HTTP 308 is incomplete, not failure. After a network
  interruption, query the persisted offset with an empty PUT and
  `Content-Range: bytes */TOTAL`, then resume from Google's `Range` response.
- Completion verifies the stored object's size, content type and checksum before
  returning a FileObject. Pass the original server-stored record, never browser
  fields. Completion is retryable; the application must prevent duplicate publishing.
- Session URLs are bearer credentials: share over HTTPS and exclude them from logs.
  GCS sessions expire after one week. Applications own pending-record expiry,
  cancellation and cleanup of unpublished objects; keep published objects immutable.
- Reuses `.oss.json`; requires object-create/get permissions and delete permission
  for cleanup. The bucket stays private. This JSON API flow includes Origin and
  does not require bucket CORS changes; XML API uploads have different CORS rules.

See [Google's resumable upload protocol](https://cloud.google.com/storage/docs/performing-resumable-uploads)
and [CORS behavior](https://cloud.google.com/storage/docs/cross-origin).
