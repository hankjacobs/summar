package acc

type Accumulator struct {
	count int64
}

func (i *Accumulator) Increment() {
	i.count++
}

func (i *Accumulator) Reset() {
	i.count = 0
}

func (i *Accumulator) Count() int64 {
	return i.count
}
