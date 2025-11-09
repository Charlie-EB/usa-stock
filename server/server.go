package server

import (
	"fmt"
	"io"
	"m/sentry"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Server() error {
	fmt.Println("sever func called")

	// Load the server's private host key
	privateBytes, err := os.ReadFile("./keys/ssh_host_rsa_key_go_usa")
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// just a reminder- sftp (file operations) > ssh (encryption) > tcp (network connection)

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			// TODO: add a Check if this public key is authorized
			// For now, let's accept any key (we'll fix this next)
			return nil, nil
		},
	}
	config.AddHostKey(private)

	// Step 3: Listen on SFTP port (usually 22, but let's use 2022 to avoid conflicts)
	listener, err := net.Listen("tcp", ":2022")
	if err != nil {
		sentry.Notify(err, "failed to listen on sftp port")
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer listener.Close()

	fmt.Println("SFTP server listening on :2022")

	// Step 4: Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		// Handle each connection in a goroutine
		go handleConnection(conn, config)
	}
}

// this func is the ssh layer. its given a raw tcp connection via netConn
func handleConnection(netConn net.Conn, config *ssh.ServerConfig) {
	defer netConn.Close()

	// Perform SSH handshakem and returns sshConn = encrypted SSH connection tunnel, chans = channel that will receive ssh channels (what flows through the tunnel. like a stream of data)
	sshConn, chans, reqs, err := ssh.NewServerConn(netConn, config)
	if err != nil {
		sentry.Notify(err, "SSH handshake failed")
		fmt.Printf("SSH handshake failed: %v\n", err)
		return
	}
	defer sshConn.Close()

	fmt.Printf("User %s connected\n", sshConn.User())

	// Handle out-of-band requests (throw them away)
	go ssh.DiscardRequests(reqs)

	// Handle channels. when client wants to start sftp (doing stuff with files), client will request a channel. session is sftp
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Printf("could not accept channel: %v\n", err)
			continue
		}

		go handleChannel(channel, requests)
	}
}

// this func is the sftp layer where files are served?
func handleChannel(channel ssh.Channel, requests <-chan *ssh.Request) {
	defer channel.Close()

	downloadsPath := "./downloads"
	absPath, err := filepath.Abs(downloadsPath)
	if err != nil {
		sentry.Notify(err, "failed to get absolute path for downloads dir")
		fmt.Printf("failed to get absolute path: %v\n", err)
		return
	}
	// debug
	fmt.Println("📂 Serving SFTP from:", absPath)

	for req := range requests {
		// wait for net ss to say connection.download()
		if req.Type == "subsystem" && string(req.Payload[4:]) == "sftp" {
			// reply yes to net ss
			req.Reply(true, nil)

			// Use custom restricted filesystem
			fs := &restrictedFS{root: absPath}
			handlers := sftp.Handlers{
				FileGet:  fs,
				FilePut:  fs,
				FileCmd:  fs,
				FileList: fs,
			}

			// create an SFTP server on this channel, serving files
			server := sftp.NewRequestServer(
				channel,
				handlers,
			)

			if err := server.Serve(); err != nil && err != io.EOF {
				sentry.Notify(err, "sftp server error")
				fmt.Printf("sftp server error: %v\n", err)
			}
			return
		}
		req.Reply(false, nil)
	}
}
