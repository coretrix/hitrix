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

func GetFileTypeCounter(fileObject *entity.FileObject, namespace oss.Namespace) *file.File {
	objectURL, err := service.DI().OSService().GetObjectURL(namespace, fileObject)
	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:       fileObject.ID,
		Filename: fileObject.StorageKey,
		URL:      objectURL,
		IDType:   file.IDTypeOSSCounterID,
	}
}

func GetFileTypeCounterWithCDN(fileObject *entity.FileObject, namespace oss.Namespace, cdn CDNInterface) *file.File {
	objectURL, err := service.DI().OSService().GetObjectURL(namespace, fileObject)
	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:     fileObject.ID,
		Filename: fileObject.StorageKey,
		URL:    cdn.GetImageURLWithTemplate(objectURL),
		IDType: file.IDTypeOSSCounterID,
	}
}

func GetFileTypeID(fileEntity *entity.FileEntity, namespace oss.Namespace) *file.File {
	objectURL, err := service.DI().OSService().GetObjectURL(namespace, fileEntity.File)
	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:        fileEntity.ID,
		URL:       objectURL,
		Namespace: oss.Namespace(fileEntity.Namespace),
		IDType:    file.IDTypeFileID,
	}
}
