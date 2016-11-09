package core

import (
	"github.com/hankjacobs/summar/pkg/acc"
	"github.com/hankjacobs/summar/pkg/nginx"
	"github.com/hankjacobs/summar/pkg/statsd"
)

const (
	metric20xName = "20x"
	metric30xName = "30x"
	metric40xName = "40x"
	metric50xName = "50x"
)

// MetricCounter counts the number of log entries with 2xx, 3xx, 4xx,
// 5xx status codes
type MetricCounter interface {
	Reset()
	CountEntry(entry nginx.LogEntry)
	Entries20xMetric() statsd.Metric
	Entries30xMetric() statsd.Metric
	Entries40xMetric() statsd.Metric
	Entries50xMetric() statsd.Metric
	ErrorRouteMetrics() []statsd.Metric
}

type metricCounterImpl struct {
	entries20x acc.Accumulator // 20x entries
	entries30x acc.Accumulator // 30x entries
	entries40x acc.Accumulator // 40x entries
	entries50x acc.Accumulator // 50x entries

	errorRoutes map[string]acc.Accumulator
}

// NewMetricCounter creates a new metric counter
func NewMetricCounter() MetricCounter {
	return &metricCounterImpl{errorRoutes: make(map[string]acc.Accumulator)}
}

// Reset resets a metric counter to zero state
func (c *metricCounterImpl) Reset() {
	c.entries20x.Reset()
	c.entries30x.Reset()
	c.entries40x.Reset()
	c.entries50x.Reset()
	c.errorRoutes = make(map[string]acc.Accumulator)
}

// CountEntry uses entry to increment appropriate metrics
func (c *metricCounterImpl) CountEntry(entry nginx.LogEntry) {
	switch {
	case entry.Has20xStatusCode():
		c.entries20x.Increment()
	case entry.Has30xStatusCode():
		c.entries30x.Increment()
	case entry.Has40xStatusCode():
		c.entries40x.Increment()
	case entry.Has50xStatusCode():
		c.entries50x.Increment()
		c.incrementErrorRouteEntry(entry)
	}
}

// Entries20xMetric metric for 20x entries
func (c *metricCounterImpl) Entries20xMetric() statsd.Metric {
	return statsd.Metric{Name: metric20xName, Value: c.entries20x.Count(), Type: statsd.Set}
}

// Entries30xMetric metric for 30x entries
func (c *metricCounterImpl) Entries30xMetric() statsd.Metric {
	return statsd.Metric{Name: metric30xName, Value: c.entries30x.Count(), Type: statsd.Set}
}

// Entries40xMetric metric for 40x entries
func (c *metricCounterImpl) Entries40xMetric() statsd.Metric {
	return statsd.Metric{Name: metric40xName, Value: c.entries40x.Count(), Type: statsd.Set}
}

// Entries50xMetric metric for 50x entries
func (c *metricCounterImpl) Entries50xMetric() statsd.Metric {
	return statsd.Metric{Name: metric50xName, Value: c.entries50x.Count(), Type: statsd.Set}
}

// ErrorRouteMetrics metrics for routes that had a 50x status code
func (c *metricCounterImpl) ErrorRouteMetrics() []statsd.Metric {
	metrics := []statsd.Metric{}

	for key, value := range c.errorRoutes {
		m := statsd.Metric{Name: key, Value: value.Count(), Type: statsd.Set}
		metrics = append(metrics, m)
	}

	return metrics
}

func (c *metricCounterImpl) incrementErrorRouteEntry(entry nginx.LogEntry) {
	acc := c.errorRoutes[entry.Route]
	acc.Increment()
	c.errorRoutes[entry.Route] = acc
}
