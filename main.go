package main

import (
	"fmt"
	"m/fetcher"
	"m/sentry"
	"m/server"
	"m/utils"
)

func main() {

	sentry.Setup()

	fetcher.DlSanmar()
	err := server.Server()
	if err != nil {
		sentry.Notify(err, "main server error")
	}

}
