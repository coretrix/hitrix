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

type APIManager struct {
	Providers map[uint64]IProvider
}

func NewAPIManager(configService config.IConfig, newProviderFunctions ...NewProviderFunc) (*APIManager, error) {
	providers := map[uint64]IProvider{}

	for i, newProviderFunc := range newProviderFunctions {
		provider, err := newProviderFunc(configService)
		if err != nil {
			return nil, err
		}
		providers[uint64(i)] = provider
	}

	return &APIManager{
		Providers: providers,
	}, nil
}

func (a *APIManager) GetRandomProvider() IProvider {
	return a.Providers[0]
}

func (a *APIManager) GetProvider() IProvider {
	return a.Providers[0]
}
