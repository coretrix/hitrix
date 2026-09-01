package hitrix

import componentmetrics "github.com/coretrix/hitrix/service/component/metrics"

// SetCustomMetric stores an application-specific gauge or counter value so it
// is included in the regular Hitrix metrics snapshots.
func SetCustomMetric(name string, value float64) {
	componentmetrics.Set(name, value)
}

// AddCustomMetric atomically adds delta to an application-specific metric and
// returns the updated value.
func AddCustomMetric(name string, delta float64) float64 {
	return componentmetrics.Add(name, delta)
}

func customMetricsSnapshot() map[string]float64 {
	return componentmetrics.Snapshot()
}
