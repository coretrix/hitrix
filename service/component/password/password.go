package password

type IPassword interface {
	VerifyPassword(password string, hash string) bool
	HashPassword(password string) (string, error)
}
