package main

import (
	"log"
	"m/fetcher"
	"m/sentry"
	"m/server"
)

func main() {

	sentry.Setup()

	err := fetcher.DlSanmar()
	if err != nil {
		log.Printf("ERROR: dl error: %v", err)
		sentry.Notify(err, "download error in main func ")
	}
	err = server.Server()
	if err != nil {
		log.Printf("ERROR: main server error: %v", err)
		sentry.Notify(err, "main server error")
	}

}
