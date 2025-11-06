package main

import (
	"fmt"
	"m/fetcher"
	"m/server"
	"m/utils"
	"time"
)

func main() {
	env, _ := utils.GetEnv()
	fmt.Println("from main", env["FILENAME"])

	// bg file get
		go func() {
			for {
				fetcher.DlSanmar()
				time.Sleep(4 * time.Hour)
			}
		}()

	fetcher.DlSanmar()
	server.Server()
}
