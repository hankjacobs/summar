package statsd

import "testing"

func TestNewMessage(t *testing.T) {
	message := NewMessage([]Metric{Metric{}}...)

	if len(message.metrics) != 1 {
		t.Fatal("NewMessage failed to create message with metrics")
	}
}

func TestMessageString(t *testing.T) {
	metric1 := Metric{Name: "metric1", Value: 2, Type: Set}   // metric1:2|s
	metric2 := Metric{Name: "metric2", Value: 100, Type: Set} // metric2:100|s
	metric3 := Metric{Name: "metric3", Value: 0, Type: Set}   // metric3:0|s
	metrics := []Metric{metric1, metric2, metric3}
	message := NewMessage(metrics...)

	expected := "metric1:2|s\nmetric2:100|s\nmetric3:0|s\n"
	got := message.String()

	if got != expected {
		t.Fatalf("Got %v but expected %v", got, expected)
	}
}
