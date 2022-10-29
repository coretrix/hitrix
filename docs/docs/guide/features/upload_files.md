# Upload files

Hitrix allow us to upload files using OSS service and assign them to FileEntity.
Whenever we want to reassign them to the right entity we need to do something like that in separate endpoint 
```go
    //...
	fileEntity := &entity.FileEntity{}
	found := ormService.LoadByID(fileID, fileEntity)

	if !found {
		return fmt.Errorf("file with FileID %v not found", *fileID)
	}

	if fileEntity.Namespace != oss.NamespaceUserAvatar.String() {
		return goErrors.New("wrong file category")
	}

	userEntity.Avatar = fileEntity.File
    //...
```

If you want to enable this feature you should call `middleware.FileRouter(ginEngine)`
This will add `/v1/file/upload/` endpoint where the customers can upload their files
