package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtractUInt64Slice(t *testing.T) {
	a := []uint64{1, 2, 3}
	b := []uint64{1}
	res := SubtractUInt64Slice(a, b)
	assert.Equal(t, []uint64{2, 3}, res)
}

func TestSubtractInt64Slice(t *testing.T) {
	a := []int64{1, 2, 3}
	b := []int64{1}
	res := SubtractInt64Slice(a, b)
	assert.Equal(t, []int64{2, 3}, res)
}

func TestSubtractUInt32Slice(t *testing.T) {
	var a []uint32
	b := []uint32{1}
	res := SubtractUInt32Slice(a, b)
	assert.Equal(t, []uint32{}, res)
}

func TestSubtractInt32Slice(t *testing.T) {
	a := []int32{1, 2, 3}
	var b []int32
	res := SubtractInt32Slice(a, b)
	assert.Equal(t, a, res)
}

func TestSubtractUIntSlice(t *testing.T) {
	a := []uint{1, 2, 3}
	b := []uint{10, 20}
	res := SubtractUIntSlice(a, b)
	assert.Equal(t, a, res)
}

func TestSubtractIntSlice(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{10, 20}
	res := SubtractIntSlice(a, b)
	assert.Equal(t, a, res)
}
