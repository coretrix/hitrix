package file

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coretrix/hitrix/pkg/dto/file"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/pkg/errors"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/oss"
)

func CreateFile(ctx context.Context, newFile *file.RequestDTOUploadImage) (*file.File, error) {
	ormService := service.DI().OrmEngineForContext(ctx)

	now := service.DI().Clock().Now()

	ext := strings.Replace(filepath.Ext(newFile.Image.Filename), ".", "", 1)

	tempFile, err := os.CreateTemp("", fmt.Sprintf("*.%s", ext))
	clean := func() {
		err = os.Remove(tempFile.Name())
		if err != nil {
			service.DI().ErrorLogger().LogError(fmt.Sprintf("failed deleting temp file %s\nError: %s", tempFile.Name(), err.Error()))
		}
	}

	defer clean()

	buf := make([]byte, 1024)

	for {
		n, err := newFile.Image.File.Read(buf)

		if err != nil && err != io.EOF {
			panic(err)
		}

		if n == 0 {
			break
		}

		if _, err := tempFile.Write(buf[:n]); err != nil {
			panic(err)
		}
	}

	_ = tempFile.Close()

	namespace := oss.Namespace(newFile.Namespace.String())

	if namespace == "" {
		return nil, errors.HandleCustomErrors(map[string]string{"Namespace": "namespace invalid"})
	}

	obj, err := service.DI().OSService().UploadImageFromFile(ormService, namespace, tempFile.Name())

	if err != nil {
		return nil, err
	}

	fileEntity := &entity.FileEntity{
		File:      &obj,
		Namespace: newFile.Namespace.String(),
		Status:    entity.FileStatusNew.String(),
		CreatedAt: service.DI().Clock().Now(),
	}

	ormService.Flush(fileEntity)

	bucketConfig, err := service.DI().OSService().GetNamespaceBucketConfig(namespace)

	if err != nil {
		panic(err)
	}

	objectURL := ""

	switch bucketConfig.Type {
	case oss.BucketPublic:
		objectURL, err = service.DI().OSService().GetObjectURL(namespace, fileEntity.File)
	case oss.BucketPrivate:
		objectURL, err = service.DI().OSService().GetObjectSignedURL(namespace, fileEntity.File, now.Add(30*time.Minute)) // TODO make this time dynamic
	}

	if err != nil {
		panic(err)
	}

	return &file.File{
		ID:        fileEntity.ID,
		URL:       objectURL,
		Namespace: oss.Namespace(fileEntity.Namespace),
		IDType:    file.FileIDTypeFileID,
	}, nil
}
