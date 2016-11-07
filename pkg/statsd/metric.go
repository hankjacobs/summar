package statsd

import "fmt"

type Type string

const (
	Set = Type("s")
)

type Metric struct {
	Name  string
	Value int64
	Type  Type
}

func (m Metric) String() string {
	return fmt.Sprintf("%s:%d|%s", m.Name, m.Value, string(m.Type))
}
