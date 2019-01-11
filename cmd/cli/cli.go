package main

import (
	"fmt"
	"log"

	lw "github.com/joelbirchler/label-watcher/internal"
)

func main() {
	// Create a simple client
	api, err := lw.Connect()
	if err != nil {
		log.Fatal(err)
	}

	// Watch for node label changes and print them out
	ch := lw.WatchNodeLabels(api)
	for event := range ch {
		fmt.Printf("%s\n", &event)
	}
}
