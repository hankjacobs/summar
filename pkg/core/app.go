package core

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/hankjacobs/summar/pkg/nginx"
	"github.com/hankjacobs/summar/pkg/statsd"
	"github.com/hankjacobs/summar/pkg/tailer"
	"github.com/hpcloud/tail"
)

// Flush Interval
var (
	DefaultFlushInterval = time.Duration(5) * time.Second
)

// Loggers
var (
	// DefaultLogger
	DefaultLogger = log.New(os.Stderr, "", log.LstdFlags)

	// DiscardingLogger
	DiscardingLogger = log.New(ioutil.Discard, "", 0)
)

var (
	// ErrInvalidConfig error returned for invalid app config
	ErrInvalidConfig = fmt.Errorf("invalid config")
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

// NginxParseFunc parsing func used to parse nginx log lines
type NginxParseFunc func(string) (nginx.LogEntry, error)

// Config app config
type Config struct {
	Tailer  tailer.Tailer
	Writer  statsd.MessageWriter
	Counter MetricCounter
	Logger  logger

	FlushInterval time.Duration

	ParseFunc NginxParseFunc
}

// App app
type App struct {
	tailer        tailer.Tailer
	writer        statsd.MessageWriter
	counter       MetricCounter
	flushInterval time.Duration
	logger        logger
	stop          chan struct{} // signals app to stop running
	done          chan struct{} // signals that app has stopped
	parseFunc     func(string) (nginx.LogEntry, error)
}

// NewApp creates a new app with the given config
func NewApp(config Config) (*App, error) {

	if config.Tailer == nil || config.Writer == nil {
		return nil, ErrInvalidConfig
	}

	if config.FlushInterval == 0 {
		return nil, ErrInvalidConfig
	}

	if config.Counter == nil {
		config.Counter = NewMetricCounter()
	}

	if config.Logger == nil {
		config.Logger = DefaultLogger
	}

	if config.ParseFunc == nil {
		config.ParseFunc = nginx.ParseLogEntry
	}

	return &App{tailer: config.Tailer, writer: config.Writer, counter: config.Counter, flushInterval: config.FlushInterval, logger: config.Logger, stop: make(chan struct{}), done: make(chan struct{}), parseFunc: config.ParseFunc}, nil
}

// Run runs an app
func (a *App) Run() {
	ticker := time.NewTicker(a.flushInterval)

	for {
		select {
		case line := <-a.tailer.Lines():
			a.logger.Printf("Line channel received")
			if line != nil {
				a.handleLine(line)
			} else {
				a.logger.Printf("No line present")
			}
		case <-ticker.C:
			a.logger.Printf("Flush timer fired")
			a.writeCounts()
		case <-a.stop:
			close(a.done)
			ticker.Stop()
			return
		}
	}
}

// Stop stops an app
func (a *App) Stop() {
	close(a.stop)
	a.wait()
}

// wait waits until the app has stopped
func (a *App) wait() {
	select {
	case <-a.done:
		break
	}
}

func (a *App) handleLine(line *tail.Line) {
	a.logger.Printf("Parsing new log line")
	entry, err := a.parseFunc(line.Text)
	if err != nil {
		a.logger.Printf("Could not parse new log line %v", err)
		return
	}
	a.logger.Printf("Counting log entry")
	a.counter.CountEntry(entry)
}

func (a *App) writeCounts() {

	metrics := []statsd.Metric{}
	metrics = append(metrics, a.counter.Entries20xMetric())
	metrics = append(metrics, a.counter.Entries30xMetric())
	metrics = append(metrics, a.counter.Entries40xMetric())
	metrics = append(metrics, a.counter.Entries50xMetric())
	metrics = append(metrics, a.counter.ErrorRouteMetrics()...)
	a.counter.Reset()

	message := statsd.NewMessage(metrics...)
	writer := a.writer
	go func() {
		a.logger.Printf("Writing message...")

		if err := writer.Write(message); err != nil {
			a.logger.Printf("Failed to write message %v", err)
		} else {
			a.logger.Printf("Wrote message")
		}
	}()
}
