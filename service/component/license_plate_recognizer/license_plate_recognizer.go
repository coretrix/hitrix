package licenseplaterecognizer

type LicensePlateRecognizer interface {
	RecognizeFromImage(base64image string) ([]string, error)
}
