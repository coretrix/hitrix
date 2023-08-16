package registry

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/latolukasz/beeorm/v2"
	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
)

type ORMRegistryInitFunc func(registry *beeorm.Registry)

var sequence int

func ServiceProviderOrmRegistry(init ORMRegistryInitFunc) *service.DefinitionGlobal {
	var err error
	var ormConfig beeorm.ValidatedRegistry
	var appService *app.App
	var configService config.IConfig

	return &service.DefinitionGlobal{
		Name: service.ORMConfigService,
		Build: func(ctn di.Container) (interface{}, error) {
			appService = ctn.Get(service.AppService).(*app.App)
			configService = ctn.Get(service.ConfigService).(config.IConfig)

			registry := beeorm.NewRegistry()

			configuration, ok := configService.Get("orm")
			if !ok {
				return nil, errors.New("no orm config")
			}

			yamlConfig := map[string]interface{}{}
			for k, v := range configuration.(map[interface{}]interface{}) {
				yamlConfig[fmt.Sprint(k)] = v
			}

			if appService.IsInTestMode() {
				overwriteORMConfig(appService, configService, yamlConfig)
			}

			registry.InitByYaml(yamlConfig)

			// TODO: check why is removed
			//if appService.IsInTestMode() {
			//registry.ForceEntityLogInAllEntities("")
			//}

			init(registry)

			ormConfig, err = registry.Validate()

			return ormConfig, err
		},
		Close: func(obj interface{}) error {
			return nil
		},
	}
}

func overwriteORMConfig(appService *app.App, configService config.IConfig, yamlConfig map[string]interface{}) {
	mysqlConnection := strings.Split(configService.MustString("orm.default.mysql.uri"), "/")

	db, err := sql.Open("mysql", mysqlConnection[0]+"/?multiStatements=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	newDBName := "t_" + appService.ParallelTestID
	color.Blue("DB name: %s", newDBName)

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + newDBName + "`")

	if err != nil {
		panic(err)
	}

	yamlConfig["default"].(map[interface{}]interface{})["mysql"] = map[interface{}]interface{}{
		"uri": mysqlConnection[0] +
			"/" +
			newDBName +
			"?multiStatements=true",
	}

	connectionString, has := configService.String("orm.log_db_pool.mysql.uri")
	if has {
		mysqlLogConnection := strings.Split(connectionString, "/")

		dbLog, err := sql.Open("mysql", mysqlLogConnection[0]+"/?multiStatements=true")
		if err != nil {
			panic(err)
		}

		defer dbLog.Close()

		newDBLogName := newDBName + "_log"

		_, err = db.Exec("CREATE DATABASE IF NOT EXISTS `" + newDBLogName + "`")

		if err != nil {
			panic(err)
		}

		yamlConfig["log_db_pool"].(map[interface{}]interface{})["mysql"] = map[interface{}]interface{}{
			"uri": mysqlLogConnection[0] +
				"/" +
				newDBLogName +
				"?multiStatements=true",
		}
	}

	for _, value := range yamlConfig {
		if _, ok := value.(map[interface{}]interface{})["sentinel"]; ok {
			for masterConf := range value.(map[interface{}]interface{})["sentinel"].(map[interface{}]interface{}) {
				settings := strings.Split(fmt.Sprint(masterConf), ":")

				_, has = os.LookupEnv("REDIS_TEST")
				if !has {
					panic("Please set `REDIS_TEST` ENV variable")
				}

				sequence++
				//host:dbIndex:namespace
				value.(map[interface{}]interface{})["redis"] = os.Getenv("REDIS_TEST") + ":" + settings[1] + ":" + newDBName + fmt.Sprint(sequence)
				delete(value.(map[interface{}]interface{}), "sentinel")
			}
		}
	}
}
