package file

import (
	"github.com/coretrix/hitrix/pkg/dto/file"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"
)

type CDNInterface interface {
	GetImageURLWithTemplate(image string) string
}

func GetFileTypeCounter(fileOss *entity.FileObject, bucket oss.Bucket, private bool) *file.File {
	if fileOss == nil {
		return nil
	}

	objectURL, err := service.DI().OSService().GetObjectURL(bucket, fileOss)

	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:     fileOss.ID,
		URL:    objectURL,
		IDType: file.FileIDTypeCounterID,
	}
}

func GetFileTypeCounterWithCDN(fileOss *entity.FileObject, bucket oss.Bucket, private bool, cdn CDNInterface) *file.File {
	if fileOss == nil {
		return nil
	}

	objectURL, err := service.DI().OSService().GetObjectURL(bucket, fileOss)

	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:     fileOss.ID,
		URL:    cdn.GetImageURLWithTemplate(objectURL),
		IDType: file.FileIDTypeCounterID,
	}
}

func GetFileTypeID(fileEntity *entity.FileEntity, bucket oss.Bucket) *file.File {
	objectURL, err := service.DI().OSService().GetObjectURL(bucket, fileEntity.File)

	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:        fileEntity.ID,
		URL:       objectURL,
		Namespace: oss.Namespace(fileEntity.Namespace),
		IDType:    file.FileIDTypeFileID,
	}
}
