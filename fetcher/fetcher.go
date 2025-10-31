package fetcher

import (
	"fmt"
	"log"
	"m/utils"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func GetDir(dir string) error {
	client, err := connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	// List files in remote directory
	files, err := client.ReadDir(dir)
	if err != nil {
		log.Fatal("Failed to read directory:", err)
	}

	fmt.Println("Files on server:")
	for _, file := range files {
		fmt.Println(" -", file.Name())
	}
}

func DlSanmar() error {
	env, err := utils.GetEnv()
	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}
	path := env["DIR"]
	filename := env["FILENAME"]

	client, err := connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	client.Open(path + "/" + filename)

}

func connect() (*sftp.Client, error) {

	env, err := utils.GetEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}

	url := env["URL"]
	port := env["PORT"]
	username := env["USERNAME"]
	password := env["PASSWORD"]
	host := fmt.Sprintf("%s:%s", url, port)

	// SSH client config
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // ⚠️ only for testing
	}

	// Connect to the server
	conn, err := ssh.Dial("tcp", host, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)

	}

	// Create SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to start client: %w", err)
	}

	return client, nil
}
