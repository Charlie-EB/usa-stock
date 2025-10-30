package utils

import "os"

func EnvGetter() ( []byte, error ){

	workingDir, _ := os.Getwd()

	filePath:= workingDir + "/.env"
	return os.ReadFile(filePath)




}
