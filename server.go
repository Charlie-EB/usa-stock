package main

import (
	"fmt"
	"os"
)


func getEnv() {
	var hello = "test"
	hello = "new assigned val"
	println(hello)
}

func main() {
	getEnv()
	
	var test string = "test"
	println("hello world")
	println("hello new line")
	fmt.Fprint(os.Stdout, "hellofmt")
	fmt.Fprint(os.Stdout, test)

	testFunc()

	// sftp.NewClient()
}

