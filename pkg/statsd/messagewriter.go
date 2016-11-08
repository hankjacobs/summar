package statsd

import "io"

// MessageWriter writes a list of metrics
type MessageWriter interface {
	Write(Message) error
}

type ioWriterImpl struct {
	writer io.Writer
}

// NewIOMessageWriter creates a new MessageWriter that writes using
// the specified io.Writer
func NewIOMessageWriter(writer io.Writer) MessageWriter {
	return &ioWriterImpl{writer: writer}
}

func (i *ioWriterImpl) Write(message Message) error {
	_, err := i.writer.Write(message.Bytes())

	return err
}
