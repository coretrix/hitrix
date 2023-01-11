package licenceplaterecognizer

type LicencePlateRecognizer interface {
	RecognizeFromImage(base64image string) ([]string, error)
}
