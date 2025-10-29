package main

import (
	"fmt"
	"os"
)

func main() {

	var test string = "test"
	print("hello world")
	println("hello new line")
	fmt.Fprint(os.Stdout, "hellofmt")
	fmt.Fprint(os.Stdout, test)


	// sftp.NewClient()

}

