package main

import (
	"log"
	"m/sentry"
	"m/server"
)

func main() {

	sentry.Setup()

	// fetcher.DlSanmar()

	err := server.Server()
	if err != nil {
		log.Printf("ERROR: main server error: %v", err)
		sentry.Notify(err, "main server error")
	}

}
