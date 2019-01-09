package main

import (
	"fmt"
	"log"

	lw "github.com/joelbirchler/label-watcher/internal"
)

func main() {
	api, err := lw.Connect()
	if err != nil {
		log.Fatal(err)
	}

	ch := lw.WatchNodeLabels(api)
	for event := range ch {
		fmt.Printf("-- %v\n", event)
	}
}
