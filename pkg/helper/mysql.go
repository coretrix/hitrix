package helper

import (
	"fmt"

	"github.com/latolukasz/beeorm"
)

type transaction func() error

func DBTransaction(ormService *beeorm.Engine, callback transaction) error {
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

func Limit(pager *beeorm.Pager) string {
	return fmt.Sprintf("LIMIT %d,%d", (pager.CurrentPage-1)*pager.PageSize, pager.PageSize)
}
