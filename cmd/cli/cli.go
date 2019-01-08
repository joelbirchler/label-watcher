package main

import (
	"log"

	lw "github.com/joelbirchler/label-watcher/internal"
)

func main() {
	api, err := lw.Connect()
	if err != nil {
		log.Fatal(err)
	}

	lw.WatchNodeLabels(api)
}
