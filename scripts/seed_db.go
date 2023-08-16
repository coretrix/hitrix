package scripts

import (
	"context"
	"log"
	"os"

	"github.com/latolukasz/beeorm/v2"

	"github.com/coretrix/hitrix/datalayer"
	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
)

type DBSeedScript struct {
	SeedsPerProject map[string][]Seed
}

func (script *DBSeedScript) Run(_ context.Context, _ app.IExit) {
	ormService := service.DI().OrmEngine()
	appService := service.DI().App()
	Seeder(script.SeedsPerProject, ormService, appService)
}

func (script *DBSeedScript) Infinity() bool {
	return true
}

func (script *DBSeedScript) Unique() bool {
	return true
}

func (script *DBSeedScript) Description() string {
	return "Seed Database"
}

type Seed interface {
	Execute(*datalayer.ORM)
	Environments() []string
	Name() string
}

func Seeder(seedsPerProject map[string][]Seed, ormService *datalayer.ORM, appService *app.App) {
	for project, seeds := range seedsPerProject {
		if project != os.Getenv("PROJECT_NAME") {
			continue
		}

		for _, seed := range seeds {
			supportCurrentEnv := false

			for _, env := range seed.Environments() {
				if env == appService.Mode {
					supportCurrentEnv = true

					break
				}
			}

			if !supportCurrentEnv {
				continue
			}

			seederEntity := &entity.SeederEntity{}

			whereStmt := beeorm.NewWhere("`Name` = ?", seed.Name())

			found := ormService.SearchOne(whereStmt, seederEntity)
			if found {
				continue
			}

			seed.Execute(ormService)

			seederEntity.Name = seed.Name()
			seederEntity.CreatedAt = service.DI().Clock().Now()
			ormService.Flush(seederEntity)

			log.Println("Seeder " + seed.Name() + " has been executed")
		}
	}
}
