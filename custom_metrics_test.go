package hitrix

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCustomMetricsAreThreadSafeAndSnapshotValues(t *testing.T) {
	const name = "test_concurrent_metric"
	SetCustomMetric(name, 0)

	var waitGroup sync.WaitGroup
	for range 100 {
		waitGroup.Add(1)
		go func() {
			defer waitGroup.Done()
			AddCustomMetric(name, 1)
		}()
	}
	waitGroup.Wait()

	require.Equal(t, float64(100), customMetricsSnapshot()[name])

	SetCustomMetric(name, 7.5)
	require.Equal(t, 7.5, customMetricsSnapshot()[name])
}
