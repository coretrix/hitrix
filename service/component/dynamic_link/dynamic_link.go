package dynamiclink

type IGenerator interface {
	GenerateDynamicLink(string) (*GenerateResponse, error)
}

type GenerateResponse struct {
	Link        string
	PreviewLink string
}
