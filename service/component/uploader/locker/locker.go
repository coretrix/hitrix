package locker

import (
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"

	"github.com/latolukasz/orm"
	"github.com/sarulabs/di"
	tusd "github.com/tus/tusd/pkg/handler"
)

type GetLockerFunc func(ctn di.Container) tusd.Locker

func GetRedisLocker(ctn di.Container) tusd.Locker {
	appService := ctn.Get(service.AppService).(*app.App)
	ormService := ctn.Get(service.ORMEngineGlobalService).(*orm.Engine)

	return &RedisLocker{ctx: appService.Ctx, ormService: ormService}
}
