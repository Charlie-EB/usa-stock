package utils

import (
	"fmt"
	"os"
	"strings"

	"golang.org/x/crypto/ssh"
)

// LoadAuthorizedKeys reads the authorized_keys file from Docker secrets
// and returns a map of authorized public keys for quick lookup
func LoadAuthorizedKeys(secretName string) (map[string]bool, error) {
	path := fmt.Sprintf("/run/secrets/%s", secretName)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read authorized_keys: %w", err)
	}

	authorizedKeys := make(map[string]bool)
	lines := strings.Split(string(data), "\n")

	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse the public key
		pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(line))
		if err != nil {
			// Log but don't fail on individual bad keys
			fmt.Printf("Warning: failed to parse key on line %d: %v\n", lineNum+1, err)
			continue
		}

		// Store the key fingerprint as the map key
		fingerprint := ssh.FingerprintSHA256(pubKey)
		authorizedKeys[fingerprint] = true
	}

	if len(authorizedKeys) == 0 {
		return nil, fmt.Errorf("no valid authorized keys found")
	}

	return authorizedKeys, nil
}

// IsKeyAuthorized checks if a given public key is in the authorized keys map
func IsKeyAuthorized(clientKey ssh.PublicKey, authorizedKeys map[string]bool) bool {
	fingerprint := ssh.FingerprintSHA256(clientKey)
	return authorizedKeys[fingerprint]
}

// LoadHostKey loads an SSH host key from Docker secrets
func LoadHostKey(secretName string) (ssh.Signer, error) {
	path := fmt.Sprintf("/run/secrets/%s", secretName)
	privateKeyBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read host key %s: %w", secretName, err)
	}

	privateKey, err := ssh.ParsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse host key %s: %w", secretName, err)
	}

	return privateKey, nil
}
