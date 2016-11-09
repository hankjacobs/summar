package statsd

import (
	"bytes"
	"testing"
)

func TestWrite(t *testing.T) {
	metric1 := Metric{Name: "metric1", Value: 2, Type: Set}   // metric1:2|s
	metric2 := Metric{Name: "metric2", Value: 100, Type: Set} // metric2:100|s
	metric3 := Metric{Name: "metric3", Value: 0, Type: Set}   // metric3:0|s
	metrics := []Metric{metric1, metric2, metric3}
	message := NewMessage(metrics...)

	var buf bytes.Buffer
	writer := NewIOMessageWriter(&buf)
	writer.Write(message)

	expected := "metric1:2|s\nmetric2:100|s\nmetric3:0|s\n"
	got := buf.String()

	if expected != got {
		t.Fatalf("Got %s but expected %s", got, expected)
	}
}

func TestWriteMessageEmpty(t *testing.T) {

	var buf bytes.Buffer

	writer := NewIOMessageWriter(&buf)
	writer.Write(Message{})

	expected := "" // NO newline in empty write
	got := buf.String()

	if expected != got {
		t.Fatalf("Got %s but expected %s", got, expected)
	}
}
