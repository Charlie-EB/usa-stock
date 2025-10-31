package main

import (
	"fmt"
	"m/fetcher"
	"m/utils"
)

func main() {

	println("hello world")

	env, _ := utils.GetEnv()

	fmt.Println("from main", env["FILENAME"])

	fetcher.GetDir("/")
}
