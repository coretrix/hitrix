package googleanalytics

import (
	"os"

	"github.com/coretrix/hitrix/service/component/config"
	errorlogger "github.com/coretrix/hitrix/service/component/error_logger"
)

type NewProviderFunc func(configFolder string, configService config.IConfig, errorLogger errorlogger.ErrorLogger) (IProvider, error)

type IAPIManager interface {
	GetProvider(name Provider) IProvider
}

type APIManager struct {
	Providers        map[string]IProvider
	ProvidersByIndex map[int]IProvider
}

func NewAPIManager(
	localConfigFolder string,
	configService config.IConfig,
	errorLogger errorlogger.ErrorLogger,
	newProviderFunctions ...NewProviderFunc,
) (IAPIManager, error) {
	providers := map[string]IProvider{}
	providersByIndex := map[int]IProvider{}

	var configFolder string

	appFolder, hasConfigFolder := os.LookupEnv("APP_FOLDER")
	if !hasConfigFolder {
		configFolder = localConfigFolder
	} else {
		configFolder = appFolder + "/config"
	}

	for i, newProviderFunc := range newProviderFunctions {
		provider, err := newProviderFunc(configFolder, configService, errorLogger)
		if err != nil {
			return nil, err
		}
		providers[provider.GetName().String()] = provider
		providersByIndex[i] = provider
	}

	return &APIManager{
		Providers:        providers,
		ProvidersByIndex: providersByIndex,
	}, nil
}

func (a *APIManager) GetProvider(name Provider) IProvider {
	return a.Providers[name.String()]
}
