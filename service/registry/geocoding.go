package registry

import (
	"errors"
	"fmt"

	"github.com/latolukasz/beeorm"
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

			cacheTTLMinDays := 0
			cacheTTLMaxDays := 0

			if okUseCaching && useCaching {
				ormConfig := ctn.Get(service.ORMConfigService).(beeorm.ValidatedRegistry)
				entities := ormConfig.GetEntities()
				if _, ok := entities["entity.GeocodingCacheEntity"]; !ok {
					return nil, errors.New("you should register GeocodingCacheEntity")
				}

				if _, ok := entities["entity.GeocodingReverseCacheEntity"]; !ok {
					return nil, errors.New("you should register GeocodingReverseCacheEntity")
				}

				var has bool
				cacheTTLMinDays, has = configService.Int("geocoding.cache_ttl_min_days")
				if !has {
					return nil, fmt.Errorf("you must specify geocoding.cache_ttl_min_days")
				}

				cacheTTLMaxDays, has = configService.Int("geocoding.cache_ttl_max_days")
				if !has {
					return nil, fmt.Errorf("you must specify geocoding.cache_ttl_max_days")
				}
			}

			provider, err := providerConstructor(configService)
			if err != nil {
				return nil, err
			}

			return geocoding.NewGeocoding(
				useCaching,
				cacheTTLMinDays,
				cacheTTLMaxDays,
				ctn.Get(service.ClockService).(clock.IClock),
				provider,
			), nil
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
