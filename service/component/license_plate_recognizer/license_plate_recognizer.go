package licenseplaterecognizer

type LicensePlateRecognizer interface {
	RecognizeFromImage(base64image string, mainRegion string) ([]string, error)
}
