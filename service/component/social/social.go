package social

type IUserData interface {
	GetUserData(token string) (*UserData, error)
}

type UserData struct {
	FirstName string
	LastName  string
	Avatar    string
	Email     string
}
