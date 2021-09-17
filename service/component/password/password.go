package password

type Password interface {
	VerifyPassword(password string, hash string) bool
	HashPassword(password string) (string, error)
}
