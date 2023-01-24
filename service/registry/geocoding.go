package registry

import (
	"errors"
	"fmt"

	"github.com/sarulabs/di"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/clock"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/geocoding"
)

const (
	GeocodingProviderGoogleMaps = "google_maps"
)

func ServiceProviderGeocoding(provider string) *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.GeocodingService,
		Build: func(ctn di.Container) (interface{}, error) {
			providerConstructor, ok := providerConstructorFactory[provider]
			if !ok {
				return nil, fmt.Errorf("provider constructor not found by key: %s", provider)
			}

			configService := ctn.Get(service.ConfigService).(config.IConfig)

			useCaching, okUseCaching := configService.Bool("geocoding.use_caching")

			cacheExpiryDays, okCacheExpiryDays := configService.Int64("geocoding.cache_expiry_days")

			if okUseCaching && !okCacheExpiryDays {
				return nil, fmt.Errorf("you must specify geocoding.cache_expiry_days")
			}

			provider, err := providerConstructor(configService)
			if err != nil {
				return nil, err
			}

			return geocoding.NewGeocoding(useCaching, cacheExpiryDays, ctn.Get(service.ClockService).(clock.IClock), provider), nil
		},
	}
}

var providerConstructorFactory = map[string]func(configService config.IConfig) (geocoding.Provider, error){
	GeocodingProviderGoogleMaps: providerConstructorGoogleMaps,
}

func providerConstructorGoogleMaps(configService config.IConfig) (geocoding.Provider, error) {
	apiKey, ok := configService.String("geocoding.google_maps.api_key")
	if !ok {
		return nil, errors.New("missing geocoding.google_maps.api_key")
	}

	return geocoding.NewGoogleMapsProvider(apiKey), nil
}
