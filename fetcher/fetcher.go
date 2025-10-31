package fetcher

import (
	"fmt"
	"log"
	"m/utils"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func GetDir(dir string) {

	client, err := connect()
	if err != nil {
		log.Fatal(err)
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

func connect() (*sftp.Client, error) {

	env, err := utils.GetEnv()
	if err != nil {
		log.Fatal("Failed to load .env:", err)
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
		log.Fatal("Failed to dial:", err)
	}

	// Create SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal("Failed to create SFTP client:", err)
	}

	return client, nil
}
