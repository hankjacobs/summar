package core

import (
	"bytes"
	"io"
	"time"

	"github.com/hankjacobs/summar/pkg/nginx"
	"github.com/hankjacobs/summar/pkg/tailer"
	"github.com/hpcloud/tail"
)

// App app
type App struct {
	tailer  tailer.Tailer
	writer  io.Writer
	counter *MetricCounter
	stop    chan struct{}
}

// NewApp creates a new app
func NewApp(tailer tailer.Tailer, writer io.WriteCloser) *App {
	return &App{tailer, writer, NewMetricCounter(), make(chan struct{})}
}

// Run runs an app
func (a *App) Run() {
	for {
		select {
		case line := <-a.tailer.Lines():
			a.handleLine(line)
		case <-time.After(5 * time.Second):
			a.writeCounts()
		case <-a.stop:
			return
		}
	}
}

// Stop stops an app
func (a *App) Stop() {
	close(a.stop)
}

func (a *App) handleLine(line *tail.Line) {
	entry, err := nginx.ParseLogEntry(line.Text)
	if err != nil {
		return
	}

	a.counter.CountEntry(entry)
}

func (a *App) writeCounts() {

	metric20x := a.counter.Entries20xMetric()
	metric30x := a.counter.Entries30xMetric()
	metric40x := a.counter.Entries40xMetric()
	metric50x := a.counter.Entries50xMetric()
	errorRouteMetrics := a.counter.ErrorRouteMetrics()
	a.counter.Reset()

	var msgBuffer bytes.Buffer
	msgBuffer.WriteString(metric20x.String())
	msgBuffer.WriteString("\n")
	msgBuffer.WriteString(metric30x.String())
	msgBuffer.WriteString("\n")
	msgBuffer.WriteString(metric40x.String())
	msgBuffer.WriteString("\n")
	msgBuffer.WriteString(metric50x.String())
	msgBuffer.WriteString("\n")

	for _, metric := range errorRouteMetrics {
		msgBuffer.WriteString(metric.String())
		msgBuffer.WriteString("\n")
	}

	writer := a.writer
	go func() {
		writer.Write(msgBuffer.Bytes())
	}()
}