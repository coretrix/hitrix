package localizer

type Source interface {
	Push([]string) error
	Pull() (map[string]string, error)
}
