package generator

type IGenerator interface {
	GenerateRandomRangeNumber(int64, int64) int64
	GenerateSha256Hash(string) string
	RandomPasswordGenerator(int) string
	RandomPINCodeGenerator(int) string
}
