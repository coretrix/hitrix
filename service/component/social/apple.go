package social

import (
	"context"
	"fmt"

	"github.com/Timothylock/go-signin-with-apple/apple"
)

type Apple struct {
	teamID          string
	clientID        string
	androidClientID string
	keyID           string
	privateKey      string
}

func NewAppleSocial(
	teamID string,
	clientID string,
	androidClientID string,
	keyID string,
	privateKey string,
) IUserData {
	return &Apple{
		teamID:          teamID,
		clientID:        clientID,
		androidClientID: androidClientID,
		keyID:           keyID,
		privateKey:      privateKey,
	}
}

func (a *Apple) GetUserData(ctx context.Context, token string, isAndroid bool) (*UserData, error) {
	if isAndroid && a.androidClientID == "" {
		return nil, fmt.Errorf("you must set androidClientID")
	}

	if !isAndroid && a.clientID == "" {
		return nil, fmt.Errorf("you must set clientID")
	}

	clientID := a.clientID
	if isAndroid {
		clientID = a.androidClientID
	}

	secret, err := apple.GenerateClientSecret(a.privateKey, a.teamID, clientID, a.keyID)
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
