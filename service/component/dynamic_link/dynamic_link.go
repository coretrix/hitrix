package dynamiclink

type Generator interface {
	GenerateDynamicLink(string) (*GenerateResponse, error)
}

type GenerateResponse struct {
	Link        string
	PreviewLink string
}
