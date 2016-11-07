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

type MetricCounter struct {
	entries20x acc.Accumulator // 20x entries
	entries30x acc.Accumulator // 30x entries
	entries40x acc.Accumulator // 40x entries
	entries50x acc.Accumulator // 50x entries

	errorRoutes map[string]acc.Accumulator
}

func NewMetricCounter() *MetricCounter {
	return &MetricCounter{errorRoutes: make(map[string]acc.Accumulator)}
}

func (c *MetricCounter) Reset() {
	c.entries20x.Reset()
	c.entries30x.Reset()
	c.entries40x.Reset()
	c.entries50x.Reset()
	c.errorRoutes = make(map[string]acc.Accumulator)
}
func (c *MetricCounter) CountEntry(entry nginx.LogEntry) {
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

func (c *MetricCounter) Entries20xMetric() statsd.Metric {
	return statsd.Metric{metric20xName, c.entries20x.Count(), statsd.Set}
}

func (c *MetricCounter) Entries30xMetric() statsd.Metric {
	return statsd.Metric{metric30xName, c.entries30x.Count(), statsd.Set}
}

func (c *MetricCounter) Entries40xMetric() statsd.Metric {
	return statsd.Metric{metric40xName, c.entries40x.Count(), statsd.Set}
}

func (c *MetricCounter) Entries50xMetric() statsd.Metric {
	return statsd.Metric{metric50xName, c.entries50x.Count(), statsd.Set}
}

func (c *MetricCounter) ErrorRouteMetrics() []statsd.Metric {
	metrics := []statsd.Metric{}

	for key, value := range c.errorRoutes {
		m := statsd.Metric{key, value.Count(), statsd.Set}
		metrics = append(metrics, m)
	}

	return metrics
}

func (c *MetricCounter) incrementErrorRouteEntry(entry nginx.LogEntry) {
	acc := c.errorRoutes[entry.Route]
	acc.Increment()
	c.errorRoutes[entry.Route] = acc
}
