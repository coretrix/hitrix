package locker

import (
	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/service"
)

type GetLockerFunc func(ctn di.Container) tusd.Locker

func GetRedisLocker(ctn di.Container) tusd.Locker {
	ormService := ctn.Get(service.ORMEngineGlobalService).(*datalayer.DataLayer)

	return &RedisLocker{ormService: ormService}
}
