package generator

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"github.com/AmirSoleimani/VoucherCodeGenerator/vcgen"
	"math/big"
	"strings"

	"github.com/coretrix/hitrix/pkg/helper"
)

type Generator interface {
	GenerateRandomRangeNumber(int64, int64) int64
	GenerateSha256Hash(string) string
	GenerateUUID() string
	GenerateRandomCode(*vcgen.Generator) string
}

type SimpleGenerator struct {
}

func (g *SimpleGenerator) GenerateRandomRangeNumber(min, max int64) int64 {
	bg := big.NewInt(max - min)

	n, err := rand.Int(rand.Reader, bg)
	if err != nil {
		panic(err)
	}
	return n.Int64() + min
}

func (g *SimpleGenerator) GenerateSha256Hash(input string) string {
	sha256Hash := sha256.New()
	sha256Hash.Write([]byte(input))
	return base64.StdEncoding.EncodeToString(sha256Hash.Sum(nil))
}

func (g *SimpleGenerator) GenerateUUID() string {
	return helper.GenerateUUID()
}

func (g *SimpleGenerator) GenerateRandomCode(generator *vcgen.Generator) string {
	vc := vcgen.New(generator)
	result, err := vc.Run()
	if err != nil {
		panic(err)
	}

	return strings.Join(*result, "")
}
