package social

import "context"

type IUserData interface {
	GetUserData(ctx context.Context, token string, isAndroid bool) (*UserData, error)
}

type UserData struct {
	FirstName string
	LastName  string
	Avatar    string
	Email     string
}
