package main

import (
	"fmt"
	"m/utils"
) 

func main() {
	
	println("hello world")

	env ,_ := utils.GetEnv()

	fmt.Println("from main", env["FILENAME"])


}

