package helper

import (
	"fmt"

	"github.com/latolukasz/orm"
)

type transaction func() error

func DBTransaction(ormService *orm.Engine, callback transaction) error {
	dbService := ormService.GetMysql()

	dbService.Begin()

	err := callback()
	if err != nil {
		dbService.Rollback()
		return err
	}
	dbService.Commit()

	return nil
}

func Limit(pager *orm.Pager) string {
	return fmt.Sprintf("LIMIT %d,%d", (pager.CurrentPage-1)*pager.PageSize, pager.PageSize)
}
