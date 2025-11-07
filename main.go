package main

import (
	"fmt"
	"m/sentry"
	"m/server"
	"m/utils"
)

func main() {
	env, _ := utils.GetEnv()
	fmt.Println("from main", env["FILENAME"])

	sentry.Setup()

	// fetcher.DlSanmar()
	err := server.Server()
	if err != nil {
		sentry.Notify(err, "main server error")
	}

}
