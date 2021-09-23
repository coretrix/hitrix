package generator

import (
	"github.com/AmirSoleimani/VoucherCodeGenerator/vcgen"
)

type IGenerator interface {
	GenerateRandomRangeNumber(int64, int64) int64
	GenerateSha256Hash(string) string
	GenerateRandomCode(*vcgen.Generator) string
}
