package statsd

import "fmt"

// Type statsd metric type
type Type string

const (
	// Set statsd set metric type
	Set = Type("s")
)

// Metric Statsd metric
type Metric struct {
	Name  string
	Value int64
	Type  Type
}

func (m Metric) String() string {
	return fmt.Sprintf("%s:%d|%s", m.Name, m.Value, string(m.Type))
}
