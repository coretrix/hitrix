package uploader

import (
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/tus/tusd/pkg/gcsstore"
	tusd "github.com/tus/tusd/pkg/handler"
	"github.com/tus/tusd/pkg/s3store"
)

type Uploader interface {
	PostFile(w http.ResponseWriter, r *http.Request)
	PatchFile(w http.ResponseWriter, r *http.Request)
	DelFile(w http.ResponseWriter, r *http.Request)
	GetFile(w http.ResponseWriter, r *http.Request)
	HeadFile(w http.ResponseWriter, r *http.Request)
	Middleware(h http.Handler) http.Handler
	SupportedExtensions() string
	Metrics() tusd.Metrics
	GetCreatedUploadsChan() chan tusd.HookEvent
	GetCompletedUploadsChan() chan tusd.HookEvent
	GetTerminatedUploadsChan() chan tusd.HookEvent
	GetUploadProgressChan() chan tusd.HookEvent
	GetBucket() string
}

type TUSDUploader struct {
	handler *tusd.UnroutedHandler
	bucket  string
}

func GetStore(OSSClient interface{}, bucket string) Store {
	if _, ok := OSSClient.(storage.Client); ok {
		return &GoogleStore{store: gcsstore.New(bucket, OSSClient.(gcsstore.GCSAPI))}
	}

	if _, ok := OSSClient.(*s3.S3); ok {
		return &AmazonStore{store: s3store.New(bucket, OSSClient.(s3store.S3API))}
	}

	panic("OSSClient store not found")
}

func NewTUSDUploader(c tusd.Config, bucket string) Uploader {
	uploader, err := tusd.NewUnroutedHandler(c)
	if err != nil {
		panic(err)
	}

	return &TUSDUploader{handler: uploader, bucket: bucket}
}

func (u *TUSDUploader) GetBucket() string {
	return u.bucket
}

func (u *TUSDUploader) PostFile(w http.ResponseWriter, r *http.Request) {
	u.handler.PostFile(w, r)
}

func (u *TUSDUploader) PatchFile(w http.ResponseWriter, r *http.Request) {
	u.handler.PatchFile(w, r)
}
func (u *TUSDUploader) DelFile(w http.ResponseWriter, r *http.Request) {
	u.handler.DelFile(w, r)
}

func (u *TUSDUploader) GetFile(w http.ResponseWriter, r *http.Request) {
	u.handler.GetFile(w, r)
}

func (u *TUSDUploader) HeadFile(w http.ResponseWriter, r *http.Request) {
	u.handler.HeadFile(w, r)
}

func (u *TUSDUploader) Middleware(h http.Handler) http.Handler {
	return u.handler.Middleware(h)
}

func (u *TUSDUploader) SupportedExtensions() string {
	return u.handler.SupportedExtensions()
}

func (u *TUSDUploader) Metrics() tusd.Metrics {
	return u.handler.Metrics
}

func (u *TUSDUploader) GetCreatedUploadsChan() chan tusd.HookEvent {
	return u.handler.CreatedUploads
}
func (u *TUSDUploader) GetCompletedUploadsChan() chan tusd.HookEvent {
	return u.handler.CompleteUploads
}
func (u *TUSDUploader) GetTerminatedUploadsChan() chan tusd.HookEvent {
	return u.handler.TerminatedUploads
}
func (u *TUSDUploader) GetUploadProgressChan() chan tusd.HookEvent {
	return u.handler.UploadProgress
}
