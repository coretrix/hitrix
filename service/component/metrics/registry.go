package metrics

import (
	"math"
	"sync"
	"sync/atomic"
)

var values sync.Map

type metric struct {
	bits atomic.Uint64
}

func Set(name string, value float64) {
	item := get(name)
	item.bits.Store(math.Float64bits(value))
}

func Add(name string, delta float64) float64 {
	item := get(name)
	for {
		currentBits := item.bits.Load()
		updated := math.Float64frombits(currentBits) + delta
		if item.bits.CompareAndSwap(currentBits, math.Float64bits(updated)) {
			return updated
		}
	}
}

func Snapshot() map[string]float64 {
	result := make(map[string]float64)
	values.Range(func(key, value interface{}) bool {
		result[key.(string)] = math.Float64frombits(value.(*metric).bits.Load())
		return true
	})

	return result
}

func get(name string) *metric {
	value, _ := values.LoadOrStore(name, &metric{})
	return value.(*metric)
}
