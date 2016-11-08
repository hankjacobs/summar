package statsd

import "bytes"

// Message statsd message
type Message struct {
	metrics []Metric
}

// NewMessage creates a new message with the given metrics
func NewMessage(metrics ...Metric) Message {
	return Message{metrics: metrics}
}

// Bytes returns the message as a slice of bytes
func (m Message) Bytes() []byte {
	var msgBuffer bytes.Buffer
	for _, metric := range m.metrics {
		msgBuffer.WriteString(metric.String())
		msgBuffer.WriteString("\n")
	}

	return msgBuffer.Bytes()
}

func (m Message) String() string {
	return string(m.Bytes())
}
