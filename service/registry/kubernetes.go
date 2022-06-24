package registry

import (
	"errors"
	"os"
	"path"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/kubernetes"
)

func ServiceProviderKubernetes() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.KubernetesService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			configFilePath := ""
			configFile, ok := configService.String("kubernetes.config_file")
			if ok {
				var configFolder string

				if path.IsAbs(configFile) {
					configFilePath = configFile
				} else {
					appFolder, hasConfigFolder := os.LookupEnv("APP_FOLDER")
					if !hasConfigFolder {
						configFolder = ctn.Get("config_directory").(string)
					} else {
						configFolder = appFolder + "/config"
					}

					configFilePath = path.Join(configFolder, configFile)
				}
			}

			environment, ok := configService.String("kubernetes.environment")
			if !ok {
				return nil, errors.New("missing kubernetes.environment")
			}

			return kubernetes.NewKubernetes(configFilePath, environment), nil
		},
	}
}
