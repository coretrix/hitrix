package entity

import "github.com/latolukasz/orm"

func Init(registry *orm.Registry) {
	registry.RegisterEntity(
		&AdminUserEntity{},
	)
}
