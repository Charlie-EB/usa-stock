package fetcher

import (
	"fmt"
	"log"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Getdir() {

	env := utils.

	username := 
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
	defer conn.Close()

	// Create SFTP client
	client, err := sftp.NewClient(conn)
	if err != nil {
		log.Fatal("Failed to create SFTP client:", err)
	}
	defer client.Close()

	// List files in remote directory
	remoteDir := "/"
	files, err := client.ReadDir(remoteDir)
	if err != nil {
		log.Fatal("Failed to read directory:", err)
	}

	fmt.Println("Files on server:")
	for _, file := range files {
		fmt.Println(" -", file.Name())
	}

		

}