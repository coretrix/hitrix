package graphqlparser

type Errors []struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}

func (e Errors) Error() string {
	return e[0].Message
}
