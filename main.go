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

	/*
	// bg file get
	go func() {
		for {
			time.Sleep(4 * time.Hour)
			fetcher.DlSanmar()
		}
	}()
	fetcher.DlSanmar()
	*/
	
	server.Server()

}
