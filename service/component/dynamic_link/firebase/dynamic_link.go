package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	dynamiclink "github.com/coretrix/hitrix/service/component/dynamic_link"
)

const (
	dynamicLinkGenerateURL  = "https://firebasedynamiclinks.googleapis.com/v1/shortLinks?key="
	suffixOptionShort       = "SHORT"
	suffixOptionUnguessable = "UNGUESSABLE"
)

var ValidSuffixOptionsMap = map[string]*struct{}{
	suffixOptionShort:       {},
	suffixOptionUnguessable: {},
}

type Generator struct {
	ApiKey          string           `json:"-"`
	DynamicLinkInfo *DynamicLinkInfo `json:"dynamicLinkInfo"`
	Suffix          *Suffix          `json:"suffix,omitempty"`
}

type DynamicLinkInfo struct {
	DomainUriPrefix   string             `json:"domainUriPrefix"`
	Link              string             `json:"link"`
	AndroidInfo       *AndroidInfo       `json:"androidInfo,omitempty"`
	IosInfo           *IosInfo           `json:"iosInfo,omitempty"`
	NavigationInfo    *NavigationInfo    `json:"navigationInfo,omitempty"`
	AnalyticsInfo     *AnalyticsInfo     `json:"analyticsInfo,omitempty"`
	SocialMetaTagInfo *SocialMetaTagInfo `json:"socialMetaTagInfo,omitempty"`
}

type AndroidInfo struct {
	AndroidPackageName           string `json:"androidPackageName,omitempty"`
	AndroidFallbackLink          string `json:"androidFallbackLink,omitempty"`
	AndroidMinPackageVersionCode string `json:"androidMinPackageVersionCode,omitempty"`
}

type IosInfo struct {
	IosBundleID         string `json:"iosBundleId,omitempty"`
	IosFallbackLink     string `json:"iosFallbackLink,omitempty"`
	IosCustomScheme     string `json:"iosCustomScheme,omitempty"`
	IosIpadFallbackLink string `json:"iosIpadFallbackLink,omitempty"`
	IosIpadBundleID     string `json:"iosIpadBundleId,omitempty"`
	IosAppStoreID       string `json:"iosAppStoreId,omitempty"`
}

type NavigationInfo struct {
	EnableForcedRedirect bool `json:"enableForcedRedirect"`
}

type AnalyticsInfo struct {
	GooglePlayAnalytics    *GooglePlayAnalytics    `json:"googlePlayAnalytics"`
	ItunesConnectAnalytics *ItunesConnectAnalytics `json:"itunesConnectAnalytics"`
}

type GooglePlayAnalytics struct {
	UtmSource   string `json:"utmSource"`
	UtmMedium   string `json:"utmMedium"`
	UtmCampaign string `json:"utmCampaign"`
	UtmTerm     string `json:"utmTerm"`
	UtmContent  string `json:"utmContent"`
	GclID       string `json:"gclid"`
}

type ItunesConnectAnalytics struct {
	At string `json:"at"`
	Ct string `json:"ct"`
	Mt string `json:"mt"`
	Pt string `json:"pt"`
}

type SocialMetaTagInfo struct {
	SocialTitle       string `json:"socialTitle"`
	SocialDescription string `json:"socialDescription"`
	SocialImageLink   string `json:"socialImageLink"`
}

type Suffix struct {
	Option string `json:"option"`
}

// GenerateDynamicLink hash must be URL safe
func (g *Generator) GenerateDynamicLink(hash string) (*dynamiclink.GenerateResponse, error) {
	g.DynamicLinkInfo.Link += fmt.Sprintf("?hash=%s", hash)

	marshaled, err := json.Marshal(g)
	if err != nil {
		return &dynamiclink.GenerateResponse{}, err
	}

	respRaw, err := http.Post(fmt.Sprintf("%s%s", dynamicLinkGenerateURL, g.ApiKey), "application/json", bytes.NewReader(marshaled))
	if err != nil {
		return &dynamiclink.GenerateResponse{}, err
	}

	defer func() {
		if err := respRaw.Body.Close(); err != nil {
			panic(err)
		}
	}()

	resp := &response{}
	if err := json.NewDecoder(respRaw.Body).Decode(resp); err != nil {
		return &dynamiclink.GenerateResponse{}, err
	}

	return &dynamiclink.GenerateResponse{
		Link:        resp.ShortLink,
		PreviewLink: resp.PreviewLink,
	}, nil
}

type response struct {
	ShortLink   string `json:"shortLink"`
	PreviewLink string `json:"previewLink"`
}
