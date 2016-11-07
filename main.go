package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/hankjacobs/summar/pkg/core"
	"github.com/hankjacobs/summar/pkg/tailer"
)

func main() {

	var in = flag.String("in", "/var/log/nginx/access.log", "Input file")
	var out = flag.String("out", "/var/log/stats.log", "Output file location")
	flag.Parse()

	tailer, tailErr := tailer.NewTailer(*in)
	if tailErr != nil {
		log.Fatalf("Could not open %v", *in)
	}

	file, fileERr := os.OpenFile(*out, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if fileERr != nil {
		log.Fatalf("Could not open %v", *out)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	done := make(chan struct{})

	app := core.NewApp(tailer, file)
	go func() {
		app.Run()
		close(done)
	}()

	select {
	case <-c:
		log.Printf("Received sigterm")
		break
	case <-done:
		log.Printf("Premature death of app")
		break
	}

	log.Printf("Stopping")
	_ = tailer.Stop()
	_ = file.Close()
	log.Printf("Stopped")
}
