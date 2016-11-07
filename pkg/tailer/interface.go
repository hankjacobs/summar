package tailer

import "github.com/hpcloud/tail"

// Tailer is used to tail a file
type Tailer interface {
	Lines() chan *tail.Line
	Stop() error
}
