package tailer

import (
	"os"

	"github.com/hpcloud/tail"
)

type logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Fatalln(v ...interface{})
	Panic(v ...interface{})
	Panicf(format string, v ...interface{})
	Panicln(v ...interface{})
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// Tailer is used to tail a file
type Tailer interface {
	// ideally, we would use an interface instead of *tail.Line
	// but casting between chan types isn't possible
	Lines() chan *tail.Line
	Stop() error
}

var _ Tailer = &tailer{}

type tailer struct {
	impl *tail.Tail
}

// NewTailer creates a new tailer for the specified file
func NewTailer(filename string, logger logger) (Tailer, error) {
	config := tail.Config{Follow: true, ReOpen: true, Logger: logger}

	// Only use location if the file already exists.
	// If we use location, and the file does not exist,
	// then it will seek to the end of the file as soon
	// as it does exist, effectively ignoring the first line
	// written
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		config.Location = nil
	} else {
		config.Location = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
	}

	impl, err := tail.TailFile(filename, config)

	if err != nil {
		return nil, err
	}

	return &tailer{impl}, nil
}

func (t *tailer) Lines() chan *tail.Line {
	return t.impl.Lines
}

func (t *tailer) Stop() error {
	return t.impl.Stop()
}
