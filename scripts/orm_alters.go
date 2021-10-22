package scripts

import (
	"context"

	"github.com/coretrix/hitrix/service"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix"
	"github.com/fatih/color"
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
	return service.DI().OrmConfig() != nil
}

func (script *ORMAltersScript) Unique() bool {
	return true
}

func (script *ORMAltersScript) Description() string {
	return "show all MySQL schema changes"
}

func (script *ORMAltersScript) Run(_ context.Context, exit hitrix.Exit) {
	ormEngine := service.DI().OrmEngine()
	alters := ormEngine.GetAlters()
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
