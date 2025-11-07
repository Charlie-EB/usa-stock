package main

import (
	"fmt"
	"m/sentry"
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
			fetcher.DlSanmar()
			time.Sleep(4 * time.Hour)
		}
	}()

	fetcher.DlSanmar()
	server.Server()

*/	
}
