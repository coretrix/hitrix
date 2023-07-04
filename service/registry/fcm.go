package registry

import (
	"os"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/app"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/fcm"
)

const (
	fcmConfigEnvName = "FIREBASE_CONFIG"
)

func ServiceProviderFCM() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.FCMService,
		Build: func(ctn di.Container) (interface{}, error) {
			val, present := os.LookupEnv(fcmConfigEnvName)

			configService := ctn.Get(service.ConfigService).(config.IConfig)

			if present {
				if _, err := os.Stat(val); err != nil {
					if os.IsNotExist(err) {
						// specified config file doesn't exists
						credentialsFile := configService.GetFolderPath() + "/.fcm.json"
						err := os.Setenv(fcmConfigEnvName, credentialsFile)
						if err != nil {
							return nil, err
						}
					}
				}
			}

			return fcm.NewFCM(ctn.Get(service.AppService).(*app.App).GlobalContext)
		},
	}
}
