package acc

// Accumulator accumulates a value
type Accumulator struct {
	count int64
}

// Increment increments the accumulator
func (i *Accumulator) Increment() {
	i.count++
}

// Reset resets the accumulator back to zero
func (i *Accumulator) Reset() {
	i.count = 0
}

// Count returns the current count
func (i *Accumulator) Count() int64 {
	return i.count
}
