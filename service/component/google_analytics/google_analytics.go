package googleanalytics

import (
	"github.com/coretrix/hitrix/service/component/config"
)

type NewProviderFunc func(configService config.IConfig) (IProvider, error)

type IAPIManager interface {
	GetProvider(name Provider) IProvider
}

type APIManager struct {
	Providers        map[string]IProvider
	ProvidersByIndex map[int]IProvider
}

func NewAPIManager(configService config.IConfig, newProviderFunctions ...NewProviderFunc) (IAPIManager, error) {
	providers := map[string]IProvider{}
	providersByIndex := map[int]IProvider{}

	for i, newProviderFunc := range newProviderFunctions {
		provider, err := newProviderFunc(configService)
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
