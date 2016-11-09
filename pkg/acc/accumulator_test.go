package acc

import "testing"

func TestIncrement(t *testing.T) {
	acc := Accumulator{}
	old := acc.Count()
	exp := old + 1

	acc.Increment()
	if acc.Count() != exp {
		t.Fatalf("Count was %v but should have been %v", acc.Count(), exp)
	}
}

func TestReset(t *testing.T) {
	acc := Accumulator{}
	acc.Increment()
	acc.Increment()
	acc.Increment()

	acc.Reset()
	if acc.Count() != 0 {
		t.Fatalf("Count was %v but should have been %v", acc.Count(), 0)
	}
}
