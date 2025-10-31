package utils

import (
	"os"
	"strings"
)

func EnvGetter() ( map[string]string, error ){

	workingDir, _ := os.Getwd()
	filePath:= workingDir + "/.env"
	data , err:= os.ReadFile(filePath)

	s:= string(data)
	lines := strings.Split(s, "\n") // Split into lines

	myMap := make(map[string]string) // Initializes an empty map
	for _,line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "=",2)
		if len(parts) != 2 {
			continue }
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// remove optional surrounding quotes
		if len(value) > 1 {
			first := value[0]
			last := value[len(value)-1]
			if (first == '"' && last == '"') || (first == '\'' && last == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		myMap[key] = value

	}
	return myMap, err

}
