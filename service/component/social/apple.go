package social

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Timothylock/go-signin-with-apple/apple"

	"github.com/coretrix/hitrix/pkg/helper"
	"github.com/coretrix/hitrix/service/component/config"
)

type Apple struct {
	TeamID          string
	ClientID        string
	AndroidClientID string
	KeyID           string
	PrivateKey      string
}

func NewAppleSocial(
	configService config.IConfig,
) (IUserData, error) {
	if !helper.ExistsInDir(".apple-id.json", configService.GetFolderPath()) {
		return nil, errors.New(configService.GetFolderPath() + "/.apple-id.json does not exists")
	}

	credentialsFile := configService.GetFolderPath() + "/.apple-id.json"

	var dat []byte
	var configApple = &Apple{}
	var err error

	if dat, err = os.ReadFile(credentialsFile); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dat, configApple); err != nil {
		return nil, err
	}

	return configApple, nil
}

func (a *Apple) GetUserData(ctx context.Context, token string, isAndroid bool) (*UserData, error) {
	if isAndroid && a.AndroidClientID == "" {
		return nil, fmt.Errorf("you must set AndroidClientID")
	}

	if !isAndroid && a.ClientID == "" {
		return nil, fmt.Errorf("you must set ClientID")
	}

	clientID := a.ClientID
	if isAndroid {
		clientID = a.AndroidClientID
	}

	secret, err := apple.GenerateClientSecret(a.PrivateKey, a.TeamID, clientID, a.KeyID)
	if err != nil {
		return nil, err
	}

	client := apple.New()

	req := apple.AppValidationTokenRequest{
		ClientID:     clientID,
		ClientSecret: secret,
		Code:         token,
	}

	var resp apple.ValidationResponse

	err = client.VerifyAppToken(ctx, req, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	claim, err := apple.GetClaims(resp.IDToken)
	if err != nil {
		return nil, err
	}

	idClaim, ok := (*claim)["sub"]
	if !ok {
		return nil, fmt.Errorf("apple returned claims with 'sub' missling")
	}

	emailClaim, ok := (*claim)["email"]
	if !ok {
		return nil, fmt.Errorf("apple returned claims with 'email' missling")
	}

	return &UserData{ID: idClaim.(string), Email: emailClaim.(string)}, nil
}
