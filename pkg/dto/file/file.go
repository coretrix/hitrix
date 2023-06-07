package file

import (
	"io"
	"mime/multipart"

	"github.com/coretrix/hitrix/service/component/oss"
)

type File struct {
	ID        uint64
	URL       string
	Filename  string
	Namespace oss.Namespace
	IDType    IDType
}

type Upload struct {
	File        io.Reader `swaggerignore:"true"`
	Filename    string
	Size        int64
	ContentType string
}

type RequestDTOUploadFile struct {
	Namespace string                `form:"namespace" json:"namespace" binding:"required"`
	File      *multipart.FileHeader `form:"file" json:"file" binding:"required" swaggerignore:"true"`
}

func (r *RequestDTOUploadFile) ToUploadImage() (*RequestDTOUploadImage, func() error, error) {
	deferFn := func() error { return nil }

	f, err := r.File.Open()
	if err != nil {
		return nil, deferFn, err
	}

	deferFn = f.Close

	return &RequestDTOUploadImage{
		Image: Upload{
			File:        f,
			Filename:    r.File.Filename,
			Size:        r.File.Size,
			ContentType: r.File.Header.Get("Content-Type"),
		},
		Namespace: oss.Namespace(r.Namespace),
	}, deferFn, nil
}

type RequestDTOUploadImage struct {
	Image     Upload
	Namespace oss.Namespace
}

type IDType string

const (
	IDTypeOSSCounterID IDType = "oss_counter_id"
	IDTypeFileID       IDType = "file_id"
)

var AllFileIDType = []IDType{
	IDTypeOSSCounterID,
	IDTypeFileID,
}

func (f IDType) IsValid() bool {
	switch f {
	case IDTypeOSSCounterID, IDTypeFileID:
		return true
	}

	return false
}
