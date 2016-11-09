package core

import (
	"fmt"
	"testing"
	"time"

	"github.com/hankjacobs/summar/pkg/nginx"
	"github.com/hankjacobs/summar/pkg/statsd"
	"github.com/hpcloud/tail"
)

func newDummyCounter() *dummyCounter {
	return &dummyCounter{}
}

type dummyCounter struct {
	countInvokeCount int
	resetInvokeCount int
}

func (d *dummyCounter) Reset() {
	d.resetInvokeCount++
}

func (d *dummyCounter) CountEntry(entry nginx.LogEntry) {
	d.countInvokeCount++
}

func (d *dummyCounter) Entries20xMetric() statsd.Metric {
	return statsd.Metric{}
}

func (d *dummyCounter) Entries30xMetric() statsd.Metric {
	return statsd.Metric{}
}

func (d *dummyCounter) Entries40xMetric() statsd.Metric {
	return statsd.Metric{}
}

func (d *dummyCounter) Entries50xMetric() statsd.Metric {
	return statsd.Metric{}
}

func (d *dummyCounter) ErrorRouteMetrics() []statsd.Metric {
	return []statsd.Metric{}
}

func newDummyMessageWriter() *dummyMessageWriter {
	return &dummyMessageWriter{}
}

type dummyMessageWriter struct {
	writeInvokeCount int
}

func (d *dummyMessageWriter) Write(message statsd.Message) error {
	d.writeInvokeCount++
	return nil
}

type dummyTailer struct {
	stopCalled bool
	lines      chan *tail.Line
}

func newDummyTailer() *dummyTailer {
	return &dummyTailer{lines: make(chan *tail.Line)}
}

func (d *dummyTailer) Lines() chan *tail.Line {
	return d.lines
}

func (d *dummyTailer) Stop() error {
	d.stopCalled = true
	return nil
}

var defaultConfig = Config{Writer: newDummyMessageWriter(), Tailer: newDummyTailer(), FlushInterval: DefaultFlushInterval}

func newDummyApp() *App {
	app, _ := NewApp(defaultConfig)
	return app
}

func TestNewApp(t *testing.T) {

	writer := newDummyMessageWriter()
	tailer := newDummyTailer()
	flushInterval := DefaultFlushInterval
	config := Config{Writer: writer, Tailer: tailer, FlushInterval: flushInterval}

	app, _ := NewApp(config)

	if tailer != app.tailer || writer != app.writer {
		t.Fatal("NewApp did not set appropriate fields")
	}
}

func TestNewAppInvalidConfig(t *testing.T) {

	_, err := NewApp(Config{})

	if err != ErrInvalidConfig {
		t.Fatal("Config was valid when it shouldn't have been")
	}

	writer := newDummyMessageWriter()
	tailer := newDummyTailer()
	_, err = NewApp(Config{Writer: writer, Tailer: tailer, FlushInterval: 0})
	if err != ErrInvalidConfig {
		t.Fatal("Config was valid when it shouldn't have been")
	}
}

func TestStop(t *testing.T) {
	app := newDummyApp()
	stopped := make(chan struct{})

	go func() {
		app.Run()
		close(stopped)
	}()

	go func() {
		time.Sleep(250 * time.Millisecond)
		app.Stop()
	}()

	select {
	case <-stopped:
		break
	case <-time.After(1 * time.Second):
		t.Fatal("App did not stop")
		break
	}
}

func TestLineIsParsed(t *testing.T) {
	writer := newDummyMessageWriter()
	tailer := newDummyTailer()
	flush := time.Duration(1 * time.Second)

	var parseInvoked bool
	parse := func(line string) (nginx.LogEntry, error) {
		parseInvoked = true

		return nginx.LogEntry{}, nil
	}

	config := Config{Writer: writer, Tailer: tailer, FlushInterval: flush, ParseFunc: parse}

	app, _ := NewApp(config)

	go app.Run()

	app.tailer.Lines() <- tail.NewLine("test entry")
	app.Stop()

	if parseInvoked != true {
		t.Fatal("Parse not invoked on new line")
	}
}

func TestMalformedLineIsSkipped(t *testing.T) {
	writer := newDummyMessageWriter()
	tailer := newDummyTailer()
	counter := newDummyCounter()
	parse := func(line string) (nginx.LogEntry, error) {
		return nginx.LogEntry{}, fmt.Errorf("Invalid")
	}

	config := Config{Writer: writer, Tailer: tailer, FlushInterval: DefaultFlushInterval, Counter: counter, ParseFunc: parse}
	app, _ := NewApp(config)

	go app.Run()

	app.tailer.Lines() <- tail.NewLine("test entry")
	app.Stop()

	if counter.countInvokeCount != 0 {
		t.Fatal("Counter invoked for malformed line")
	}
}

func TestMessageWriterInvokedEveryFlushIntervalSeconds(t *testing.T) {

	writer := newDummyMessageWriter()
	tailer := newDummyTailer()
	counter := newDummyCounter()
	flush := time.Duration(250 * time.Millisecond)

	config := Config{Writer: writer, Tailer: tailer, FlushInterval: flush, Counter: counter}
	app, _ := NewApp(config)

	go app.Run()

	select {
	case <-time.After(flush*2 + time.Duration(10*time.Millisecond)):
		break
	}

	app.Stop()

	if writer.writeInvokeCount != 2 {
		t.Fatal("Write was not invoked twice")
	}
}

func TestCounterResetAfterFlush(t *testing.T) {

	writer := newDummyMessageWriter()
	tailer := newDummyTailer()
	counter := newDummyCounter()
	flush := time.Duration(250 * time.Millisecond)

	config := Config{Writer: writer, Tailer: tailer, FlushInterval: flush, Counter: counter}
	app, _ := NewApp(config)

	go app.Run()

	select {
	case <-time.After(flush + time.Duration(10*time.Millisecond)):
		break
	}

	app.Stop()

	if counter.resetInvokeCount != 1 {
		t.Fatal("Reset never called")
	}
}
