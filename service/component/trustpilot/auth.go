package trustpilot

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/coretrix/hitrix/pkg/entity"
	"github.com/latolukasz/beeorm"
)

type AccessToken struct {
	AccessToken           string
	AccessTokenExpiresAt  time.Time
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func (at AccessToken) accessTokenHasExpired(now time.Time) bool {
	return at.AccessTokenExpiresAt.After(now)
}

func (at AccessToken) refreshTokenHasExpired(now time.Time) bool {
	return at.RefreshTokenExpiresAt.After(now)
}

func (at AccessToken) NeedsRenew(now time.Time) bool {
	if at.accessTokenHasExpired(now) && at.refreshTokenHasExpired(now) {
		return true
	}
	return false
}

func (at AccessToken) NeedsRefresh(now time.Time) bool {
	return at.accessTokenHasExpired(now)
}

func newAccessToken(
	accessToken string,
	accessTokenTTL int,
	accessTokenIAT int,
	refreshToken string,
	refreshTokenTTL int,
	refreshTokenIAT int,
) AccessToken {
	return AccessToken{
		AccessToken: accessToken,
		AccessTokenExpiresAt: time.Unix(
			int64(accessTokenIAT+accessTokenTTL),
			int64(0),
		),
		RefreshToken: refreshToken,
		RefreshTokenExpiresAt: time.Unix(
			int64(refreshTokenIAT+refreshTokenTTL),
			int64(0),
		),
	}
}

func getNewAccessToken(
	apiKey string,
	apiSecret string,
	username string,
	password string,
) (*AccessToken, error) {
	body := url.Values{}
	body.Set("grant_type", "password")
	body.Set("username", username)
	body.Set("password", password)

	req, err := http.NewRequest(http.MethodPost, authEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	// HTTP Basic Auth
	req.SetBasicAuth(apiKey, apiSecret)

	client := http.Client{
		Timeout: time.Second * 10,

		// Add BasicAuth Header to redirect request
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.SetBasicAuth(apiKey, apiSecret)
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TrustpilotAPI accesstoken request failed with code %d", resp.StatusCode)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	respMap := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&respMap); err != nil {
		return nil, err
	}

	err = nil
	var (
		accessToken    string
		accessTokenTTL int
		accessTokenIAT int

		refreshToken    string
		refreshTokenTTL int
		refreshTokenIAT int
	)

	accessToken = respMap["access_token"]
	accessTokenTTL, err = strconv.Atoi(respMap["issued_at"])
	if err != nil {
		return nil, err
	}
	accessTokenIAT, err = strconv.Atoi(respMap["expires_in"])
	if err != nil {
		return nil, err
	}

	refreshToken = respMap["refresh_token"]
	refreshTokenTTL, err = strconv.Atoi(respMap["refresh_token_issued_at"])
	if err != nil {
		return nil, err
	}
	refreshTokenIAT, err = strconv.Atoi(respMap["refresh_token_expires_in"])
	if err != nil {
		return nil, err
	}

	AT := newAccessToken(
		accessToken,
		accessTokenTTL,
		accessTokenIAT,
		refreshToken,
		refreshTokenTTL,
		refreshTokenIAT,
	)
	return &AT, nil
}

func refreshAccessToken(
	apiKey string,
	apiSecret string,
	refreshToken string,
) (*AccessToken, error) {
	body := url.Values{}
	body.Set("grant_type", "refresh_token")
	body.Set("refresh_token", refreshToken)

	req, err := http.NewRequest(http.MethodPost, authRefreshEndpoint, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	// HTTP Basic Auth
	req.SetBasicAuth(apiKey, apiSecret)

	client := http.Client{
		Timeout: time.Second * 10,

		// Add BasicAuth Header to redirect request
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.SetBasicAuth(apiKey, apiSecret)
			return nil
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("TrustpilotAPI accesstoken-refresh request failed with code %d", resp.StatusCode)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	respMap := make(map[string]string)
	if err := json.NewDecoder(resp.Body).Decode(&respMap); err != nil {
		return nil, err
	}

	err = nil
	var (
		accessToken    string
		accessTokenTTL int
		accessTokenIAT int

		newRefreshToken    string
		newRefreshTokenTTL int
		newRefreshTokenIAT int
	)

	accessToken = respMap["access_token"]
	accessTokenTTL, err = strconv.Atoi(respMap["issued_at"])
	if err != nil {
		return nil, err
	}
	accessTokenIAT, err = strconv.Atoi(respMap["expires_in"])
	if err != nil {
		return nil, err
	}

	newRefreshToken = respMap["refresh_token"]
	newRefreshTokenTTL, err = strconv.Atoi(respMap["refresh_token_issued_at"])
	if err != nil {
		return nil, err
	}
	newRefreshTokenIAT, err = strconv.Atoi(respMap["refresh_token_expires_in"])
	if err != nil {
		return nil, err
	}

	AT := newAccessToken(
		accessToken,
		accessTokenTTL,
		accessTokenIAT,
		newRefreshToken,
		newRefreshTokenTTL,
		newRefreshTokenIAT,
	)
	return &AT, nil
}

func getSettingsAccessToken(ormService *beeorm.Engine) (*AccessToken, error) {
	var accessTokenEntity entity.SettingsEntity
	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Key", entity.HitrixSettingAll.Trustpilot)
	found := ormService.RedisSearchOne(&accessTokenEntity, query)
	if !found {
		return nil, nil
	}

	var accessToken AccessToken
	if err := json.Unmarshal([]byte(accessTokenEntity.Value), &accessToken); err != nil {
		return nil, errors.New("couldn't Unmarshal trustpilot access token setting")
	}

	return &accessToken, nil
}

func setSettingsAccessToken(ormService *beeorm.Engine, accessToken AccessToken) error {
	var accessTokenEntity entity.SettingsEntity
	query := beeorm.NewRedisSearchQuery()
	query.FilterString("Key", entity.HitrixSettingAll.Trustpilot)
	found := ormService.RedisSearchOne(&accessTokenEntity, query)
	if !found {
		accessTokenEntity = entity.SettingsEntity{
			Key: entity.HitrixSettingAll.Trustpilot,
		}
	}
	valBytes, err := json.Marshal(accessToken)
	if err != nil {
		return errors.New("couldn't Marshal trustpilot access token setting")
	}
	accessTokenEntity.Value = string(valBytes)

	ormService.Flush(&accessTokenEntity)
	return nil
}

func (tp *TrustpilotAPI) getNewAccessToken() error {
	at, err := getNewAccessToken(tp.apiKey, tp.apiSecret, tp.username, tp.password)
	if err != nil {
		return err
	}
	tp.AccessToken = at

	err = setSettingsAccessToken(tp.ormService, *at)
	if err != nil {
		return err
	}

	return nil
}

func (tp *TrustpilotAPI) refreshAccessToken() error {
	at, err := refreshAccessToken(tp.apiKey, tp.apiSecret, tp.AccessToken.RefreshToken)
	if err != nil {
		return err
	}
	if at == nil {
		return errors.New("newly refreshed Trustpilot Access-token is nil")
	}

	tp.AccessToken = at

	err = setSettingsAccessToken(tp.ormService, *at)
	if err != nil {
		return err
	}
	return nil
}

func (tp *TrustpilotAPI) authenticate() error {
	if tp.AccessToken == nil {
		if err := tp.getNewAccessToken(); err != nil {
			return err
		}
	}

	if tp.AccessToken.NeedsRenew(tp.clockService.Now()) {
		if err := tp.getNewAccessToken(); err != nil {
			return err
		}
	}

	if tp.AccessToken.NeedsRefresh(tp.clockService.Now()) {
		if err := tp.refreshAccessToken(); err != nil {
			return err
		}
	}

	return nil
}
