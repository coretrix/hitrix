package generator

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"math/big"
	"strings"

	"github.com/AmirSoleimani/VoucherCodeGenerator/vcgen"
)

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

func (g *SimpleGenerator) GenerateRandomCode(generator *vcgen.Generator) string {
	vc := vcgen.New(generator)
	result, err := vc.Run()
	if err != nil {
		panic(err)
	}

	return strings.Join(*result, "")
}

const (
	lowerCharSet   = "abcdefghijklmnopqrstuvwxyz"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%*-+.?"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

func (g *SimpleGenerator) RandomPasswordGenerator(passwordLength int) string {
	var password strings.Builder
	max := big.NewInt(int64(len(allCharSet)))

	for i := 0; i < passwordLength; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		password.WriteRune(rune(allCharSet[num.Int64()]))
	}

	return password.String()
}

func (g *SimpleGenerator) RandomPINCodeGenerator(codeLength int) string {
	charSet := lowerCharSet + numberSet
	var code strings.Builder
	max := big.NewInt(int64(len(charSet)))

	for i := 0; i < codeLength; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}
		code.WriteRune(rune(charSet[num.Int64()]))
	}

	return code.String()
}
