package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	lw "github.com/joelbirchler/label-watcher/internal"
)

var (
	port, tlsCert, tlsKey string
	eventLog              []lw.LabelEvent
)

func init() {
	flag.StringVar(&port, "port", "8443", "Port to listen on")
	flag.StringVar(&tlsCert, "tls-cert", "tls-cert.pem", "TLS concatenation of the server's certificate, any intermediates, and the CA's certificate")
	flag.StringVar(&tlsKey, "tls-key", "tls-key.pem", "TLS matching private key")
	flag.Parse()
}

func main() {
	go watcher()

	http.HandleFunc("/", eventLogHandler)
	http.HandleFunc("/healthz", healthzHandler)

	log.Printf("Serving on :%s\n", port)
	err := http.ListenAndServeTLS(":"+port, tlsCert, tlsKey, nil)
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
