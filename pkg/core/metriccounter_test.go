package core

import (
	"sort"
	"testing"

	"github.com/hankjacobs/summar/pkg/nginx"
	"github.com/hankjacobs/summar/pkg/statsd"
)

func TestReset(t *testing.T) {
	counter := NewMetricCounter()
	entry200 := nginx.LogEntry{Route: "/test", StatusCode: 200}
	entry300 := nginx.LogEntry{Route: "/test", StatusCode: 200}
	entry400 := nginx.LogEntry{Route: "/test", StatusCode: 200}
	entry500 := nginx.LogEntry{Route: "/test", StatusCode: 200}

	counter.CountEntry(entry200)
	counter.CountEntry(entry300)
	counter.CountEntry(entry400)
	counter.CountEntry(entry500)

	counter.Reset()

	metric20x := counter.Entries20xMetric()
	metric30x := counter.Entries30xMetric()
	metric40x := counter.Entries40xMetric()
	metric50x := counter.Entries50xMetric()
	errorRouteMetrics := counter.ErrorRouteMetrics()

	if metric20x.Value != 0 {
		t.Errorf("Expect 0 but got %v", metric20x.Value)
	}

	if metric30x.Value != 0 {
		t.Errorf("Expect 0 but got %v", metric30x.Value)
	}

	if metric40x.Value != 0 {
		t.Errorf("Expect 0 but got %v", metric40x.Value)
	}

	if metric50x.Value != 0 {
		t.Errorf("Expect 0 but got %v", metric50x.Value)
	}

	if len(errorRouteMetrics) != 0 {
		t.Errorf("Expect 0 error routes but got %v", len(errorRouteMetrics))
	}

}

func TestCount20xMetric(t *testing.T) {
	counter := NewMetricCounter()
	entry := nginx.LogEntry{Route: "/test", StatusCode: 200}

	counter.CountEntry(entry)
	counter.CountEntry(entry)

	expected := statsd.Metric{Name: metric20xName, Value: 2, Type: statsd.Set}
	got := counter.Entries20xMetric()

	if expected != got {
		t.Fatalf("Got %v but expected %v", expected, got)
	}
}

func TestCount30xMetric(t *testing.T) {
	counter := NewMetricCounter()
	entry := nginx.LogEntry{Route: "/test", StatusCode: 300}

	counter.CountEntry(entry)
	counter.CountEntry(entry)

	expected := statsd.Metric{Name: metric30xName, Value: 2, Type: statsd.Set}
	got := counter.Entries30xMetric()

	if expected != got {
		t.Fatalf("Got %v but expected %v", expected, got)
	}
}

func TestCount40xMetric(t *testing.T) {
	counter := NewMetricCounter()
	entry := nginx.LogEntry{Route: "/test", StatusCode: 400}

	counter.CountEntry(entry)
	counter.CountEntry(entry)

	expected := statsd.Metric{Name: metric40xName, Value: 2, Type: statsd.Set}
	got := counter.Entries40xMetric()

	if expected != got {
		t.Fatalf("Got %v but expected %v", expected, got)
	}
}

func TestCount50xMetric(t *testing.T) {
	counter := NewMetricCounter()
	entry := nginx.LogEntry{Route: "/test", StatusCode: 500}

	counter.CountEntry(entry)
	counter.CountEntry(entry)

	expected := statsd.Metric{Name: metric50xName, Value: 2, Type: statsd.Set}
	got := counter.Entries50xMetric()

	if expected != got {
		t.Fatalf("Got %v but expected %v", expected, got)
	}
}

func TestErrorRoutesMetric(t *testing.T) {
	counter := NewMetricCounter()
	entryErrorTestRoute := nginx.LogEntry{Route: "/test", StatusCode: 500}
	entryErrorOtherRoute := nginx.LogEntry{Route: "/other", StatusCode: 501}
	entryGoodRoute := nginx.LogEntry{Route: "/good", StatusCode: 200}

	counter.CountEntry(entryErrorTestRoute)
	counter.CountEntry(entryErrorTestRoute)
	counter.CountEntry(entryErrorOtherRoute)
	counter.CountEntry(entryGoodRoute)

	errorTestRouteMetric := statsd.Metric{Name: entryErrorTestRoute.Route, Value: 2, Type: statsd.Set}
	errorOtherRouteMetric := statsd.Metric{Name: entryErrorOtherRoute.Route, Value: 1, Type: statsd.Set}

	expectedMetrics := []statsd.Metric{errorTestRouteMetric, errorOtherRouteMetric}
	sort.Sort(metrics(expectedMetrics))

	gotMetrics := counter.ErrorRouteMetrics()
	sort.Sort(metrics(gotMetrics))

	for i, expected := range expectedMetrics {
		if expected != gotMetrics[i] {
			t.Errorf("Got %v but expected %v", expected, gotMetrics[i])
		}
	}

}

func TestErrorRoutesMetricEquals50xMetrics(t *testing.T) {
	counter := NewMetricCounter()
	entryErrorTestRoute := nginx.LogEntry{Route: "/test", StatusCode: 500}
	entryErrorOtherRoute := nginx.LogEntry{Route: "/other", StatusCode: 501}
	entryGoodRoute := nginx.LogEntry{Route: "/good", StatusCode: 200}

	counter.CountEntry(entryErrorTestRoute)
	counter.CountEntry(entryErrorTestRoute)
	counter.CountEntry(entryErrorOtherRoute)
	counter.CountEntry(entryErrorOtherRoute)
	counter.CountEntry(entryGoodRoute)

	errorRouteMetrics := counter.ErrorRouteMetrics()
	count := int64(0)
	for _, metric := range errorRouteMetrics {
		count += metric.Value
	}

	entries50xMetric := counter.Entries50xMetric()

	if count != entries50xMetric.Value {
		t.Fatalf("Got %v but expected %v", count, entries50xMetric.Value)
	}
}

type metrics []statsd.Metric

// sort.Interface
func (m metrics) Len() int           { return len(m) }
func (m metrics) Less(i, j int) bool { return m[i].Name < m[j].Name }
func (m metrics) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
