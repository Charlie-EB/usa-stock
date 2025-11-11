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

// GetDockerSecret loads configuration from Docker secrets with .env fallback
func GetDockerSecret() (map[string]string, error) {
	secrets := make(map[string]string)

	// List of all secrets we need
	secretNames := []string{
		"go_remote_url",
		"go_remote_port",
		"go_remote_username",
		"go_remote_password",
		"sentry_dsn",
	}

	for _, secretName := range secretNames {
		value, err := ReadDockerSecret(secretName)
		if err != nil {
			// Fallback to environment variable for local dev
			envName := strings.ToUpper(secretName)
			value = os.Getenv(envName)
			if value == "" {
				return nil, fmt.Errorf("config %s not found in secrets or env", secretName)
			}
			fmt.Printf("Using env var for %s (secret not found)\n", secretName)
		}
		secrets[secretName] = value
	}

	return secrets, nil
}
