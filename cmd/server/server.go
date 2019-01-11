package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	lw "github.com/joelbirchler/label-watcher/internal"
)

var (
	port     string
	eventLog []lw.LabelEvent
)

func init() {
	flag.StringVar(&port, "port", "8080", "Port to listen on")
	flag.Parse()
}

func main() {
	go watcher()

	http.HandleFunc("/", eventLogHandler)
	http.HandleFunc("/healthz", healthzHandler)

	log.Printf("Serving on :%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != http.ErrServerClosed {
		log.Fatalln("Could not start origin server:", err)
	}
}

func healthzHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func eventLogHandler(w http.ResponseWriter, req *http.Request) {
	for _, event := range eventLog {
		fmt.Fprintln(w, &event)
	}
}

func watcher() {
	// Create a simple client
	api, err := lw.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Watch for node label changes and add them to our list of events
	ch := lw.WatchNodeLabels(api)
	for event := range ch {
		eventLog = append(eventLog, event)
	}
}
