package scripts

import (
	"context"
	"github.com/latolukasz/beeorm"

	"github.com/fatih/color"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

func ORMAlters() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: "orm-alters",

		Script: true,
		Build: func(ctn di.Container) (interface{}, error) {
			return &ORMAltersScript{}, nil
		},
	}
}

type ORMAltersScript struct {
}

func (script *ORMAltersScript) Active() bool {
	_ = service.DI().OrmConfig()

	return true
}

func (script *ORMAltersScript) Unique() bool {
	return true
}

func (script *ORMAltersScript) Description() string {
	return "show all MySQL schema changes"
}

func (script *ORMAltersScript) Run(_ context.Context, exit app.IExit, ormService *beeorm.Engine) {
	alters := ormService.GetAlters()

	for _, alter := range alters {
		if alter.Safe {
			color.Green("%s\n\n", alter.SQL)
		} else {
			color.Red("%s\n\n", alter.SQL)
		}
	}

	if len(alters) > 0 {
		exit.Error()
	}
}
