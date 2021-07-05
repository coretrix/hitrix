package locker

import (
	"github.com/coretrix/hitrix/service"

	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"
)

type GetLockerFunc func(ctn di.Container) tusd.Locker

func GetRedisLocker(ctn di.Container) tusd.Locker {
	ormService := ctn.Get(service.ORMEngineGlobalService).(*orm.Engine)

	return &RedisLocker{ormService: ormService}
}
