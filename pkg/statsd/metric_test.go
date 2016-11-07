package statsd

import "testing"

func TestString(t *testing.T) {
	m := Metric{"50x", 15, Set}
	str := m.String()

	valid := "50x:15|s"
	if str != valid {
		t.Fatalf("Got %s but expected %s", str, valid)
	}

}
