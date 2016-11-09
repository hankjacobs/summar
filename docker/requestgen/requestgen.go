package main

import (
	"math/rand"
	"net/http"
	"time"
)

func main() {

	routes := []string{"/200", "/300", "/400", "/500", "/bad1", "/bad2", "/bad3"}
	for {
		_, _ = http.Get("http://localhost" + routes[rand.Intn(len(routes))])
		<-time.After(time.Duration(10) * time.Millisecond)
	}
}
