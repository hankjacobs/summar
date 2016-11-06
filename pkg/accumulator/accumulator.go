package accumulator

type Accumulator interface {
	Increment()
	Reset()
	Count() int64
}
