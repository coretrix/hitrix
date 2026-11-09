package oss

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/trixorm"
	"golang.org/x/oauth2"

	"github.com/coretrix/hitrix/pkg/entity"
)

// IResumableUploadProvider is an optional capability. Type-assert an IProvider to
// this interface before use; providers without it retain the existing API.
type IResumableUploadProvider interface {
	StartResumableUpload(ctx context.Context, ormService *trixorm.Engine,
		namespace Namespace, options ResumableUploadOptions) (*ResumableUpload, error)
	CompleteResumableUpload(ctx context.Context, namespace Namespace, upload *ResumableUpload) (entity.FileObject, error)
}

// ResumableUploadOptions describes a file before any bytes are uploaded.
// Applications must authorize the origin, namespace and size before starting.
type ResumableUploadOptions struct {
	Extension   string // A single extension including the dot, e.g. ".apk".
	ContentType string
	Size        int64
	Origin      string // Browser origin, e.g. "https://admin.example.com".
	CRC32C      string // Base64-encoded, big-endian CRC32C (Castagnoli), as used by GCS.
}

// ResumableUpload must be stored server-side and associated with the authorized
// user. Never construct the completion record from browser-supplied fields.
// UploadURL is a bearer credential: do not log it; share it only over HTTPS.
type ResumableUpload struct {
	File      entity.FileObject
	Options   ResumableUploadOptions
	UploadURL string
}

var _ IResumableUploadProvider = (*GoogleOSS)(nil)

// StartResumableUpload reserves a normal Hitrix file ID and creates a GCS JSON API
// upload session. No file bytes pass through Hitrix. A failed initiation may leave
// a gap in the counter; reserved IDs must never be reused.
func (ossStorage *GoogleOSS) StartResumableUpload(
	ctx context.Context,
	ormService *trixorm.Engine,
	namespace Namespace,
	options ResumableUploadOptions,
) (*ResumableUpload, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := validateResumableUploadOptions(options); err != nil {
		return nil, err
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return nil, err
	}

	if ormService == nil {
		return nil, errors.New("ORM is required to reserve an upload file ID")
	}

	storageCounter := getStorageCounter(ormService, bucketConfig)
	file := entity.FileObject{
		ID: storageCounter, StorageKey: ossStorage.getObjectKey(namespace, storageCounter, options.Extension),
	}

	return ossStorage.startResumableUpload(ctx, bucketConfig.Name, file, options)
}

func (ossStorage *GoogleOSS) startResumableUpload(
	ctx context.Context,
	bucket string,
	file entity.FileObject,
	options ResumableUploadOptions,
) (*ResumableUpload, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	metadata := struct {
		Name        string `json:"name"`
		ContentType string `json:"contentType"`
		CRC32C      string `json:"crc32c"`
	}{Name: file.StorageKey, ContentType: options.ContentType, CRC32C: options.CRC32C}

	body, err := json.Marshal(metadata)
	if err != nil {
		return nil, err
	}

	endpoint := "https://storage.googleapis.com/upload/storage/v1/b/" + url.PathEscape(bucket) +
		"/o?uploadType=resumable&ifGenerationMatch=0"

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Origin", options.Origin)
	request.Header.Set("X-Upload-Content-Type", options.ContentType)
	request.Header.Set("X-Upload-Content-Length", strconv.FormatInt(options.Size, 10))

	client := ossStorage.resumableUploadClient(ctx)

	response, err := client.Do(request)
	if err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		// Transport/OAuth errors can contain URLs or credential response bodies.
		return nil, errors.New("failed to initiate Google Storage resumable upload")
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google storage upload initiation returned HTTP %d", response.StatusCode)
	}

	uploadURL := response.Header.Get("Location")

	parsedURL, err := url.Parse(uploadURL)
	if err != nil || parsedURL.Scheme != "https" || parsedURL.Host != "storage.googleapis.com" ||
		parsedURL.User != nil || parsedURL.Fragment != "" || parsedURL.Query().Get("upload_id") == "" {
		return nil, errors.New("google storage returned an invalid upload session URL")
	}

	return &ResumableUpload{File: file, Options: options, UploadURL: uploadURL}, nil
}

func (ossStorage *GoogleOSS) resumableUploadClient(ctx context.Context) *http.Client {
	client := *ossStorage.uploadClient
	if ossStorage.uploadAuth == nil {
		return &client
	}

	baseTransport := client.Transport
	if baseTransport == nil {
		baseTransport = http.DefaultTransport
	}

	// JWT token exchange uses PostForm without attaching its context. Propagate
	// the initiation context explicitly so token acquisition cancels promptly.
	tokenClient := client
	tokenClient.Transport = &uploadContextTransport{ctx: ctx, base: baseTransport}
	tokenContext := context.WithValue(ctx, oauth2.HTTPClient, &tokenClient)
	client.Transport = &oauth2.Transport{
		Source: ossStorage.uploadAuth(tokenContext), Base: baseTransport,
	}

	return &client
}

type uploadContextTransport struct {
	ctx  context.Context
	base http.RoundTripper
}

func (transport *uploadContextTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport.base.RoundTrip(request.Clone(transport.ctx))
}

// CompleteResumableUpload verifies a completed object without downloading it.
// Pass the original server-stored record, never a record supplied by the browser.
// This operation is read-only and retryable; the application publishes once.
func (ossStorage *GoogleOSS) CompleteResumableUpload(
	ctx context.Context,
	namespace Namespace,
	upload *ResumableUpload,
) (entity.FileObject, error) {
	if err := ctx.Err(); err != nil {
		return entity.FileObject{}, err
	}

	if upload == nil {
		return entity.FileObject{}, errors.New("upload record is required")
	}

	if err := validateResumableUploadOptions(upload.Options); err != nil {
		return entity.FileObject{}, err
	}

	bucketConfig, err := ossStorage.namespaces.getBucketConfig(namespace)
	if err != nil {
		return entity.FileObject{}, err
	}

	expectedKey := ossStorage.getObjectKey(namespace, upload.File.ID, upload.Options.Extension)
	if upload.File.ID == 0 || upload.File.StorageKey != expectedKey {
		return entity.FileObject{}, errors.New("upload file does not match its namespace, ID and extension")
	}

	attributes, err := ossStorage.client.Bucket(bucketConfig.Name).Object(expectedKey).Attrs(ctx)
	if err != nil {
		return entity.FileObject{}, fmt.Errorf("verify uploaded object: %w", err)
	}

	if attributes.Size != upload.Options.Size {
		return entity.FileObject{}, errors.New("uploaded object size does not match")
	}

	if attributes.ContentType != upload.Options.ContentType {
		return entity.FileObject{}, errors.New("uploaded object content type does not match")
	}

	checksum := make([]byte, 4)
	binary.BigEndian.PutUint32(checksum, attributes.CRC32C)

	if base64.StdEncoding.EncodeToString(checksum) != upload.Options.CRC32C {
		return entity.FileObject{}, errors.New("uploaded object CRC32C does not match")
	}

	return entity.FileObject{ID: upload.File.ID, StorageKey: expectedKey}, nil
}

var uploadExtensionPattern = regexp.MustCompile(`^\.[a-zA-Z0-9]+$`)

func validateResumableUploadOptions(options ResumableUploadOptions) error {
	if options.Size <= 0 || options.Size > 5*1024*1024*1024*1024 {
		return errors.New("upload size must be positive and at most 5 TiB")
	}

	if !uploadExtensionPattern.MatchString(options.Extension) {
		return errors.New("upload extension must be a dot followed by letters or digits")
	}

	contentType, _, err := mime.ParseMediaType(options.ContentType)
	if err != nil || !strings.Contains(contentType, "/") || strings.ContainsAny(options.ContentType, "\r\n") {
		return errors.New("invalid upload content type")
	}

	origin, err := url.Parse(options.Origin)
	if err != nil || origin.Host == "" || (origin.Scheme != "https" && origin.Scheme != "http") ||
		origin.User != nil || origin.Path != "" || origin.RawQuery != "" || origin.ForceQuery || origin.Fragment != "" {
		return errors.New("upload origin must be an HTTP(S) origin without credentials, path, query or fragment")
	}

	checksum, err := base64.StdEncoding.Strict().DecodeString(options.CRC32C)
	if err != nil || len(checksum) != 4 || base64.StdEncoding.EncodeToString(checksum) != options.CRC32C {
		return errors.New("upload CRC32C must be a base64-encoded four-byte checksum")
	}

	return nil
}
