package utils

import (
	"fmt"
	"os"
	"strings"
)

// ReadDockerSecret reads a Docker secret from /run/secrets/<name>
// and returns its contents as a string.
func ReadDockerSecret(name string) (string, error) {
	path := fmt.Sprintf("/run/secrets/%s", name)

	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read secret %s: %w", name, err)
	}

	// Trim newline or trailing spaces, which are common in secret files
	return strings.TrimSpace(string(data)), nil
}
