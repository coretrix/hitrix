package registry

import (
	"errors"
	"fmt"

	"github.com/coretrix/hitrix/service"
	"github.com/coretrix/hitrix/service/component/config"
	"github.com/coretrix/hitrix/service/component/dynamic_link/firebase"

	"github.com/sarulabs/di"
)

func ServiceProviderDynamicLink() *service.DefinitionGlobal {
	return &service.DefinitionGlobal{
		Name: service.DynamicLinkService,
		Build: func(ctn di.Container) (interface{}, error) {
			configService := ctn.Get(service.ConfigService).(config.IConfig)

			apiKey, ok := configService.String("firebase.api_key")
			if !ok {
				return nil, errors.New("missing firebase.api_key key")
			}

			domainUriPrefix, ok := configService.String("firebase.dynamic_link_info.domain_uri_prefix")
			if !ok {
				return nil, errors.New("missing firebase.dynamic_link_info.domain_uri_prefix key")
			}

			link, ok := configService.String("firebase.dynamic_link_info.link")
			if !ok {
				return nil, errors.New("missing firebase.dynamic_link_info.link key")
			}

			generator := &firebase.Generator{
				ApiKey: apiKey,
				DynamicLinkInfo: &firebase.DynamicLinkInfo{
					DomainUriPrefix: domainUriPrefix,
					Link:            link,
				},
			}

			androidInfoMap, ok := configService.StringMap("firebase.dynamic_link_info.android_info")
			if ok {
				generator.DynamicLinkInfo.AndroidInfo = &firebase.AndroidInfo{
					AndroidPackageName:           androidInfoMap["package_name"],
					AndroidFallbackLink:          androidInfoMap["fallback_link"],
					AndroidMinPackageVersionCode: androidInfoMap["min_package_version_code"],
				}
			}

			iosInfoMap, ok := configService.StringMap("firebase.dynamic_link_info.ios_info")
			if ok {
				generator.DynamicLinkInfo.IosInfo = &firebase.IosInfo{
					IosBundleID:         iosInfoMap["bundle_id"],
					IosFallbackLink:     iosInfoMap["fallback_link"],
					IosCustomScheme:     iosInfoMap["custom_scheme"],
					IosIpadFallbackLink: iosInfoMap["ipad_fallback_link"],
					IosIpadBundleID:     iosInfoMap["ipad_bundle_id"],
					IosAppStoreID:       iosInfoMap["app_store_id"],
				}
			}

			enableForcedRedirect, ok := configService.Bool("firebase.dynamic_link_info.navigation_info.enable_forced_redirect")
			if ok {
				generator.DynamicLinkInfo.NavigationInfo = &firebase.NavigationInfo{
					EnableForcedRedirect: enableForcedRedirect,
				}
			}

			googlePlayAnalyticsMap, ok := configService.StringMap("firebase.dynamic_link_info.analytics_info.google_play_analytics")
			if ok {
				if generator.DynamicLinkInfo.AnalyticsInfo == nil {
					generator.DynamicLinkInfo.AnalyticsInfo = &firebase.AnalyticsInfo{}
				}

				generator.DynamicLinkInfo.AnalyticsInfo.GooglePlayAnalytics = &firebase.GooglePlayAnalytics{
					UtmSource:   googlePlayAnalyticsMap["utm_source"],
					UtmMedium:   googlePlayAnalyticsMap["utm_medium"],
					UtmCampaign: googlePlayAnalyticsMap["utm_campaign"],
					UtmTerm:     googlePlayAnalyticsMap["utm_term"],
					UtmContent:  googlePlayAnalyticsMap["utm_content"],
					GclID:       googlePlayAnalyticsMap["gcl_id"],
				}
			}

			itunesConnectAnalyticsMap, ok := configService.StringMap("firebase.dynamic_link_info.analytics_info.itunes_connect_analytics")
			if ok {
				if generator.DynamicLinkInfo.AnalyticsInfo == nil {
					generator.DynamicLinkInfo.AnalyticsInfo = &firebase.AnalyticsInfo{}
				}

				generator.DynamicLinkInfo.AnalyticsInfo.ItunesConnectAnalytics = &firebase.ItunesConnectAnalytics{
					At: itunesConnectAnalyticsMap["at"],
					Ct: itunesConnectAnalyticsMap["ct"],
					Mt: itunesConnectAnalyticsMap["mt"],
					Pt: itunesConnectAnalyticsMap["pt"],
				}
			}

			socialMetTagInfoMap, ok := configService.StringMap("firebase.dynamic_link_info.social_meta_tag_info")
			if ok {
				generator.DynamicLinkInfo.SocialMetaTagInfo = &firebase.SocialMetaTagInfo{
					SocialTitle:       socialMetTagInfoMap["social_title"],
					SocialDescription: socialMetTagInfoMap["social_description"],
					SocialImageLink:   socialMetTagInfoMap["social_image_link"],
				}
			}

			suffixOption, ok := configService.String("firebase.suffix.option")
			if ok {
				if _, ok := firebase.ValidSuffixOptionsMap[suffixOption]; !ok {
					availableSuffixOptions := make([]string, 0)
					for k := range firebase.ValidSuffixOptionsMap {
						availableSuffixOptions = append(availableSuffixOptions, k)
					}

					panic(fmt.Sprintf("invalid firebase.suffix.option value, please provide one of following: %v", availableSuffixOptions))
				}

				generator.Suffix = &firebase.Suffix{Option: suffixOption}
			}

			return generator, nil
		},
	}
}
