package statsd

type Type string

const (
	Set = Type("s")
)

type Metric struct {
	Name  string
	Value string
	Type  Type
}
