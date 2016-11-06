package tailer

import "github.com/hpcloud/tail"

type Tailer interface {
	Lines() chan *tail.Line
	Stop() error
}
