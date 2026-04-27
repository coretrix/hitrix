package locker

import (
	"github.com/coretrix/trixorm"
	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"

	"github.com/coretrix/hitrix/service"
)

type GetLockerFunc func(ctn di.Container) tusd.Locker

func GetRedisLocker(ctn di.Container) tusd.Locker {
	ormService := ctn.Get(service.ORMEngineGlobalService).(*trixorm.Engine)

	return &RedisLocker{ormService: ormService}
}
