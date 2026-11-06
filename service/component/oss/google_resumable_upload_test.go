package oss

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"cloud.google.com/go/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service/component/config"
)

type uploadRoundTripper func(*http.Request) (*http.Response, error)

func (transport uploadRoundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return transport(request)
}

func uploadTestOptions() ResumableUploadOptions {
	return ResumableUploadOptions{
		Extension: ".apk", ContentType: "application/vnd.android.package-archive",
		Size: 80 * 1024 * 1024, Origin: "https://admin.example.com", CRC32C: "4waSgw==",
	}
}

func uploadTestRecord() *ResumableUpload {
	return &ResumableUpload{
		File:      entity.FileObject{ID: 42, StorageKey: "mobile_apps/42.apk"},
		Options:   uploadTestOptions(),
		UploadURL: "https://storage.googleapis.com/upload/storage/v1/b/private/o?upload_id=secret",
	}
}

func uploadTestResponse(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(body))}
}

func uploadTestProvider(t *testing.T, transport uploadRoundTripper) *GoogleOSS {
	t.Helper()

	client := &http.Client{Transport: transport}
	storageClient, err := storage.NewClient(context.Background(), option.WithHTTPClient(client))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, storageClient.Close()) })

	return &GoogleOSS{
		client: storageClient, uploadClient: client,
		namespaces: namespacesConfig{"mobile_apps": {Name: "private"}},
	}
}

func TestGoogleResumableUploadInitiation(t *testing.T) {
	provider := uploadTestProvider(t, func(request *http.Request) (*http.Response, error) {
		assert.Equal(t, http.MethodPost, request.Method)
		assert.Equal(t, "https", request.URL.Scheme)
		assert.Equal(t, "storage.googleapis.com", request.URL.Host)
		assert.Equal(t, "/upload/storage/v1/b/private/o", request.URL.Path)
		assert.Equal(t, "resumable", request.URL.Query().Get("uploadType"))
		assert.Equal(t, "0", request.URL.Query().Get("ifGenerationMatch"))
		assert.Equal(t, "https://admin.example.com", request.Header.Get("Origin"))
		assert.Equal(t, "83886080", request.Header.Get("X-Upload-Content-Length"))
		assert.Equal(t, "application/vnd.android.package-archive", request.Header.Get("X-Upload-Content-Type"))
		assert.Equal(t, "application/json", request.Header.Get("Content-Type"))

		body, err := io.ReadAll(request.Body)
		require.NoError(t, err)
		assert.Less(t, len(body), 1024, "initiation must only send metadata, not the APK")
		assert.JSONEq(t, `{"name":"mobile_apps/42.apk","contentType":"application/vnd.android.package-archive","crc32c":"4waSgw=="}`, string(body))

		response := uploadTestResponse(http.StatusOK, "")
		response.Header.Set("Location", uploadTestRecord().UploadURL)

		return response, nil
	})

	upload, err := provider.startResumableUpload(context.Background(), "private", uploadTestRecord().File, uploadTestOptions())
	require.NoError(t, err)
	assert.Equal(t, uploadTestRecord(), upload)
}

func TestGoogleResumableUploadInvalidOptionsBeforeAllocation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*ResumableUploadOptions)
	}{
		{"zero size", func(options *ResumableUploadOptions) { options.Size = 0 }},
		{"negative size", func(options *ResumableUploadOptions) { options.Size = -1 }},
		{"invalid extension", func(options *ResumableUploadOptions) { options.Extension = "/other.apk" }},
		{"missing extension", func(options *ResumableUploadOptions) { options.Extension = "" }},
		{"empty type", func(options *ResumableUploadOptions) { options.ContentType = "" }},
		{"invalid type", func(options *ResumableUploadOptions) { options.ContentType = "not a MIME type" }},
		{"header injection", func(options *ResumableUploadOptions) { options.ContentType = "text/plain\r\nX-Test: value" }},
		{"missing origin", func(options *ResumableUploadOptions) { options.Origin = "" }},
		{"origin path", func(options *ResumableUploadOptions) { options.Origin += "/path" }},
		{"origin credentials", func(options *ResumableUploadOptions) { options.Origin = "https://user@admin.example.com" }},
		{"origin query", func(options *ResumableUploadOptions) { options.Origin += "?query=1" }},
		{"origin fragment", func(options *ResumableUploadOptions) { options.Origin += "#fragment" }},
		{"origin scheme", func(options *ResumableUploadOptions) { options.Origin = "file://admin.example.com" }},
		{"missing checksum", func(options *ResumableUploadOptions) { options.CRC32C = "" }},
		{"invalid checksum", func(options *ResumableUploadOptions) { options.CRC32C = "invalid" }},
		{"wrong checksum length", func(options *ResumableUploadOptions) { options.CRC32C = "YQ==" }},
	}
	provider := &GoogleOSS{namespaces: namespacesConfig{"mobile_apps": {Name: "private"}}}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			options := uploadTestOptions()
			testCase.change(&options)
			upload, err := provider.StartResumableUpload(context.Background(), nil, "mobile_apps", options)
			require.Error(t, err)
			assert.NotContains(t, err.Error(), "ORM", "options must be rejected before checking the allocator")
			assert.Nil(t, upload)
		})
	}

	upload, err := provider.StartResumableUpload(context.Background(), nil, "unknown", uploadTestOptions())
	require.Error(t, err)
	assert.Nil(t, upload)
}

type uploadConfig struct {
	config.IConfig
	folder string
}

func (configuration uploadConfig) GetFolderPath() string {
	return configuration.folder
}

func (configuration uploadConfig) Get(_ string, _ ...bool) (interface{}, bool) {
	return map[interface{}]interface{}{
		"private": map[interface{}]interface{}{"name": "private"},
	}, true
}

func TestGoogleResumableUploadOAuthCancellation(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	credentials, err := json.Marshal(map[string]string{
		"type": "service_account", "client_email": "test@example.com",
		"private_key": string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})),
		"token_uri":   "https://oauth2.example/token",
	})
	require.NoError(t, err)
	folder := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(folder, ".oss.json"), credentials, 0600))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	requests := 0
	originalTransport := http.DefaultTransport
	http.DefaultTransport = uploadRoundTripper(func(request *http.Request) (*http.Response, error) {
		requests++

		assert.Equal(t, "oauth2.example", request.URL.Host)
		cancel()
		assert.ErrorIs(t, request.Context().Err(), context.Canceled, "the OAuth request must share the caller's cancellation")

		return nil, context.Canceled
	})

	t.Cleanup(func() { http.DefaultTransport = originalTransport })

	storageProvider, err := NewGoogleOSS(uploadConfig{folder: folder}, nil, Namespaces{"mobile_apps": BucketPrivate})
	require.NoError(t, err)

	provider := storageProvider.(*GoogleOSS)

	t.Cleanup(func() { require.NoError(t, provider.client.Close()) })

	_, err = provider.startResumableUpload(ctx, "private", uploadTestRecord().File, uploadTestOptions())
	require.ErrorIs(t, err, context.Canceled)
	assert.Equal(t, 1, requests, "no storage request should follow a canceled token exchange")
}

func TestGoogleResumableUploadInitiationErrors(t *testing.T) {
	tests := []struct {
		name     string
		status   int
		location string
	}{
		{"unauthorized", 403, ""}, {"conflict", 412, ""}, {"server error", 503, ""},
		{"missing location", 200, ""}, {"insecure location", 200, "http://storage.googleapis.com/?upload_id=secret"},
		{"foreign location", 200, "https://example.com/?upload_id=secret"},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			provider := uploadTestProvider(t, func(_ *http.Request) (*http.Response, error) {
				response := uploadTestResponse(testCase.status, "sensitive provider response secret")
				response.Header.Set("Location", testCase.location)

				return response, nil
			})
			upload, err := provider.startResumableUpload(context.Background(), "private", uploadTestRecord().File, uploadTestOptions())
			require.Error(t, err)
			assert.Nil(t, upload)
			assert.NotContains(t, err.Error(), "secret")
		})
	}

	t.Run("transport failure", func(t *testing.T) {
		provider := uploadTestProvider(t, func(_ *http.Request) (*http.Response, error) {
			return nil, errors.New("transport error containing secret")
		})
		_, err := provider.startResumableUpload(context.Background(), "private", uploadTestRecord().File, uploadTestOptions())
		require.Error(t, err)
		assert.NotContains(t, err.Error(), "secret")
	})
}

func TestGoogleResumableUploadCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	provider := &GoogleOSS{}
	_, err := provider.StartResumableUpload(ctx, nil, "mobile_apps", uploadTestOptions())
	require.ErrorIs(t, err, context.Canceled)
	_, err = provider.CompleteResumableUpload(ctx, "mobile_apps", uploadTestRecord())
	require.ErrorIs(t, err, context.Canceled)

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	provider = uploadTestProvider(t, func(request *http.Request) (*http.Response, error) {
		cancel()

		return nil, request.Context().Err()
	})
	_, err = provider.startResumableUpload(ctx, "private", uploadTestRecord().File, uploadTestOptions())
	require.ErrorIs(t, err, context.Canceled)
}

func TestGoogleResumableUploadCompletion(t *testing.T) {
	tests := []struct {
		name   string
		status int
		body   string
		valid  bool
	}{
		{"complete", 200, `{"size":"83886080","contentType":"application/vnd.android.package-archive","crc32c":"4waSgw==","generation":"1"}`, true},
		{"wrong size", 200, `{"size":"1","contentType":"application/vnd.android.package-archive","crc32c":"4waSgw=="}`, false},
		{"wrong checksum", 200, `{"size":"83886080","contentType":"application/vnd.android.package-archive","crc32c":"AAAAAA=="}`, false},
		{"wrong type", 200, `{"size":"83886080","contentType":"text/plain","crc32c":"4waSgw=="}`, false},
		{"not completed", 404, `{"error":{"code":404,"message":"Not Found"}}`, false},
		{"forbidden", 403, `{"error":{"code":403,"message":"Forbidden"}}`, false},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			provider := uploadTestProvider(t, func(request *http.Request) (*http.Response, error) {
				assert.Equal(t, http.MethodGet, request.Method)
				assert.Contains(t, request.URL.Path, "/b/private/o/mobile_apps/42.apk")

				return uploadTestResponse(testCase.status, testCase.body), nil
			})
			for attempt := 0; attempt < 2; attempt++ {
				file, err := provider.CompleteResumableUpload(context.Background(), "mobile_apps", uploadTestRecord())
				if testCase.valid {
					require.NoError(t, err)
					assert.Equal(t, entity.FileObject{ID: 42, StorageKey: "mobile_apps/42.apk"}, file)
				} else {
					require.Error(t, err)
					assert.Empty(t, file)
				}
			}
		})
	}
}

func TestGoogleResumableUploadRejectsInvalidRecord(t *testing.T) {
	provider := uploadTestProvider(t, func(_ *http.Request) (*http.Response, error) {
		t.Fatal("invalid records must fail before querying storage")

		return nil, nil
	})

	tests := []struct {
		name   string
		change func(*ResumableUpload)
	}{
		{"other namespace", func(upload *ResumableUpload) { upload.File.StorageKey = "other/42.apk" }},
		{"other ID", func(upload *ResumableUpload) { upload.File.ID = 43 }},
		{"zero ID", func(upload *ResumableUpload) { upload.File.ID = 0 }},
		{"other extension", func(upload *ResumableUpload) { upload.File.StorageKey = "mobile_apps/42.exe" }},
		{"invalid size", func(upload *ResumableUpload) { upload.Options.Size = 0 }},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			upload := uploadTestRecord()
			testCase.change(upload)
			_, err := provider.CompleteResumableUpload(context.Background(), "mobile_apps", upload)
			require.Error(t, err)
		})
	}

	_, err := provider.CompleteResumableUpload(context.Background(), "mobile_apps", nil)
	require.Error(t, err)
	_, err = provider.CompleteResumableUpload(context.Background(), "unknown", uploadTestRecord())
	require.Error(t, err)
}

func TestGoogleResumableUploadRecordCanBeStored(t *testing.T) {
	encoded, err := json.Marshal(uploadTestRecord())
	require.NoError(t, err)
	var restored ResumableUpload
	require.NoError(t, json.Unmarshal(encoded, &restored))
	assert.Equal(t, uploadTestRecord(), &restored)
}
