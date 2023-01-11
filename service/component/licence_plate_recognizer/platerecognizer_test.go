package licenceplaterecognizer

import (
	"bufio"
	"encoding/base64"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlateRecognizerRecognizeFromImage(t *testing.T) {
	f, _ := os.Open("moskvich.jpg")
	reader := bufio.NewReader(f)
	content, _ := io.ReadAll(reader)

	encoded := base64.StdEncoding.EncodeToString(content)

	want := []string{"BP9716CX"}

	plateRecognizer := NewPlateRecognizer("7e1e7d418814c67da456418f87e66e03d7591b49")
	got, err := plateRecognizer.RecognizeFromImage(encoded)

	assert.Nil(t, err)
	assert.Equal(t, want, got)
}
