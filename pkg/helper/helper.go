package helper

import (
	cryptoRand "crypto/rand"
	"fmt"
)

//func RandStringBytes(n int, letters string) string {
//	b := make([]byte, n)
//	for i := range b {
//		b[i] = letters[rand.Intn(len(letters))]
//	}
//	return string(b)
//}

func BoolToInt(b bool) int8 {
	if b {
		return 1
	}
	return 0
}

func ChunkSliceUint64(slice []uint64, chunkSize int) [][]uint64 {
	var divided [][]uint64

	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize

		if end > len(slice) {
			end = len(slice)
		}

		divided = append(divided, slice[i:end])
	}

	return divided
}

func IsValidEnum(val string, enums map[string]string) bool {
	if _, ok := enums[val]; ok {
		return true
	}
	return false
}

func SliceUint64Difference(slice1 []uint64, slice2 []uint64) []uint64 {
	var diff []uint64

	// Loop two times, first to find slice1 strings not in slice2,
	// second loop to find slice2 strings not in slice1
	for i := 0; i < 2; i++ {
		for _, s1 := range slice1 {
			found := false
			for _, s2 := range slice2 {
				if s1 == s2 {
					found = true
					break
				}
			}
			// String not found. We add it to return slice
			if !found {
				diff = append(diff, s1)
			}
		}
		// Swap the slices, only if it was the first loop
		if i == 0 {
			slice1, slice2 = slice2, slice1
		}
	}

	return diff
}

func ConvertUint64ToHex(serialNumber uint64) string {
	return fmt.Sprintf("%0X", serialNumber)
}

func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := cryptoRand.Read(b)
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
