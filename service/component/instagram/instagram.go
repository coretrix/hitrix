package instagram

import (
	"github.com/coretrix/hitrix/service/component/config"
)

type NewProviderFunc func(configService config.IConfig) (IProvider, error)

type Account struct {
	AccountID int64
	FullName  string
	Bio       string
	Posts     int64
	Followers int64
	Following int64
	Picture   string
	IsPrivate bool
	Website   string
}

type Post struct {
	ID        string
	Title     string
	Images    []string
	CreatedAt int64
}

type IProvider interface {
	GetName() string
	GetAccount(account string) (*Account, error)
	GetFeed(accountID int64, nextPageToken string) ([]*Post, string, error)
}

type IAPIManager interface {
	GetRandomProvider() IProvider
	GetProvider(name string) IProvider
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

		providers[provider.GetName()] = provider
		providersByIndex[i] = provider
	}

	return &APIManager{
		Providers:        providers,
		ProvidersByIndex: providersByIndex,
	}, nil
}

func (a *APIManager) GetRandomProvider() IProvider {
	//rand.Seed(time.Now().UnixNano())
	//
	//return a.ProvidersByIndex[rand.Intn(len(a.ProvidersByIndex))]
	return a.ProvidersByIndex[0]
}

func (a *APIManager) GetProvider(name string) IProvider {
	return a.Providers[name]
}
