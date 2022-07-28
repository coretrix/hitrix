package helper

func UniqueString(a []string) []string {
	var res = make([]string, 0)
	var found = make(map[string]bool)
	for i := range a {
		if !found[a[i]] {
			res = append(res, a[i])
			found[a[i]] = true
		}
	}
	return res
}

func UniqueInt64(a []int64) []int64 {
	var res = make([]int64, 0)
	var found = make(map[int64]bool)
	for i := range a {
		if !found[a[i]] {
			res = append(res, a[i])
			found[a[i]] = true
		}
	}
	return res
}

func UniqueInt32(a []int32) []int32 {
	var res = make([]int32, 0)
	var found = make(map[int32]bool)
	for i := range a {
		if !found[a[i]] {
			res = append(res, a[i])
			found[a[i]] = true
		}
	}
	return res
}

func UniqueInt(a []int) []int {
	var res = make([]int, 0)
	var found = make(map[int]bool)
	for i := range a {
		if !found[a[i]] {
			res = append(res, a[i])
			found[a[i]] = true
		}
	}
	return res
}

func UniqueUInt64(a []uint64) []uint64 {
	var res = make([]uint64, 0)
	var found = make(map[uint64]bool)
	for i := range a {
		if !found[a[i]] {
			res = append(res, a[i])
			found[a[i]] = true
		}
	}
	return res
}

func UniqueUInt32(a []uint32) []uint32 {
	var res = make([]uint32, 0)
	var found = make(map[uint32]bool)
	for i := range a {
		if !found[a[i]] {
			res = append(res, a[i])
			found[a[i]] = true
		}
	}
	return res
}

func StringInArray(s string, arr ...string) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

func Int64InArray(s int64, arr ...int64) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

func Int32InArray(s int32, arr ...int32) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

func IntInArray(s int, arr ...int) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

func UIn64tInArray(s uint64, arr ...uint64) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

func UIn32tInArray(s uint32, arr ...uint32) bool {
	for i := range arr {
		if arr[i] == s {
			return true
		}
	}
	return false
}

func HasIntersectionInt64(a []int64, b []int64) bool {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	found := false
mainLoop:
	for _, l := range low {
		for _, h := range high {
			if l == h {
				found = true

				break mainLoop
			}
		}
	}
	return found
}

func HasIntersectionInt32(a []int32, b []int32) bool {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	found := false
mainLoop:
	for _, l := range low {
		for _, h := range high {
			if l == h {
				found = true

				break mainLoop
			}
		}
	}
	return found
}

func HasIntersectionInt(a []int, b []int) bool {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	found := false
mainLoop:
	for _, l := range low {
		for _, h := range high {
			if l == h {
				found = true

				break mainLoop
			}
		}
	}
	return found
}

func HasIntersectionUInt64(a []uint64, b []uint64) bool {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	found := false
mainLoop:
	for _, l := range low {
		for _, h := range high {
			if l == h {
				found = true

				break mainLoop
			}
		}
	}
	return found
}

func HasIntersectionUInt32(a []uint32, b []uint32) bool {
	low, high := a, b
	if len(a) > len(b) {
		low = b
		high = a
	}

	found := false
mainLoop:
	for _, l := range low {
		for _, h := range high {
			if l == h {
				found = true

				break mainLoop
			}
		}
	}
	return found
}

func SubtractUInt64Slice(a []uint64, b []uint64) []uint64 { // a-b
	var res = make([]uint64, 0)
	var bMap = make(map[uint64]bool)

	for i := range b {
		bMap[b[i]] = true
	}

	for i := range a {
		exist := bMap[a[i]]
		if !exist {
			res = append(res, a[i])
		}
	}
	return res
}

func SubtractInt64Slice(a []int64, b []int64) []int64 { // a-b
	var res = make([]int64, 0)
	var bMap = make(map[int64]bool)

	for i := range b {
		bMap[b[i]] = true
	}

	for i := range a {
		exist := bMap[a[i]]
		if !exist {
			res = append(res, a[i])
		}
	}
	return res
}

func SubtractInt32Slice(a []int32, b []int32) []int32 { // a-b
	var res = make([]int32, 0)
	var bMap = make(map[int32]bool)

	for i := range b {
		bMap[b[i]] = true
	}

	for i := range a {
		exist := bMap[a[i]]
		if !exist {
			res = append(res, a[i])
		}
	}
	return res
}

func SubtractUInt32Slice(a []uint32, b []uint32) []uint32 { // a-b
	var res = make([]uint32, 0)
	var bMap = make(map[uint32]bool)

	for i := range b {
		bMap[b[i]] = true
	}

	for i := range a {
		exist := bMap[a[i]]
		if !exist {
			res = append(res, a[i])
		}
	}
	return res
}

func SubtractIntSlice(a []int, b []int) []int { // a-b
	var res = make([]int, 0)
	var bMap = make(map[int]bool)

	for i := range b {
		bMap[b[i]] = true
	}

	for i := range a {
		exist := bMap[a[i]]
		if !exist {
			res = append(res, a[i])
		}
	}
	return res
}

func SubtractUIntSlice(a []uint, b []uint) []uint { // a-b
	var res = make([]uint, 0)
	var bMap = make(map[uint]bool)

	for i := range b {
		bMap[b[i]] = true
	}

	for i := range a {
		exist := bMap[a[i]]
		if !exist {
			res = append(res, a[i])
		}
	}
	return res
}
