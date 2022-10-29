package entity

import (
	"time"

	"github.com/latolukasz/beeorm"
)

type FileStatus string

func (f FileStatus) String() string {
	return string(f)
}

const (
	FileStatusNew       FileStatus = "new"
	FileStatusProcessed FileStatus = "processed"
)

type fileStatus struct {
	FileStatusNew       string
	FileStatusProcessed string
}

var FileStatusAll = fileStatus{
	FileStatusNew:       FileStatusNew.String(),
	FileStatusProcessed: FileStatusProcessed.String(),
}

type FileObject struct {
	ID         uint64
	StorageKey string
	Data       interface{}
}

type FileEntity struct {
	beeorm.ORM `orm:"table=files;redisSearch=search_pool"`
	ID         uint64 `orm:"searchable;sortable"`
	File       *FileObject
	Status     string    `orm:"required;enum=entity.FileStatusAll;searchable"`
	Namespace  string    `orm:"required;unique=Namespace;searchable"`
	CreatedAt  time.Time `orm:"time=true;searchable"`
}
