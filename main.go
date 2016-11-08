package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/hankjacobs/summar/pkg/core"
	"github.com/hankjacobs/summar/pkg/statsd"
	"github.com/hankjacobs/summar/pkg/tailer"
)

func main() {

	var in = flag.String("in", "/var/log/nginx/access.log", "Input file")
	var out = flag.String("out", "/var/log/stats.log", "Output file location")
	var verbose = flag.Bool("v", false, "Verbose logging")
	flag.Parse()

	logger := core.DiscardingLogger
	if *verbose != false {
		logger = core.DefaultLogger
	}

	tailer, tailErr := tailer.NewTailer(*in, logger)
	if tailErr != nil {
		log.Fatalf("Could not open %v", *in)
	}

	file, fileERr := os.OpenFile(*out, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if fileERr != nil {
		log.Fatalf("Could not open %v", *out)
	}
	writer := statsd.NewIOMessageWriter(file)

	config := core.Config{Tailer: tailer, Writer: writer, FlushInterval: core.DefaultFlushInterval, Logger: logger}
	app, err := core.NewApp(config)

	if err != nil {
		log.Fatalf("INTERNAL ERROR -- Invalid app config")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan struct{})

	go func() {
		app.Run()
		close(done)
	}()

	select {
	case <-c:
		logger.Printf("Received sigterm")
		break
	case <-done:
		logger.Printf("Premature death of app")
		break
	}

	logger.Printf("Stopping")
	app.Stop()
	_ = tailer.Stop()
	_ = file.Close()
	logger.Printf("Stopped")
}
