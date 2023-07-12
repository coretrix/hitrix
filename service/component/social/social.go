package social

import "context"

type IUserData interface {
	GetUserData(ctx context.Context, token string) (*UserData, error)
	SetIsAndroid(isAndroid bool)
}

type UserData struct {
	FirstName string
	LastName  string
	Avatar    string
	Email     string
}
