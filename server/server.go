package server

import (
	"fmt"
	"io"
	"m/sentry"
	"m/utils"
	"net"
	"os"
	"path/filepath"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func Server() error {
	fmt.Println("sever func called")

	// read docker secret which returns a string
	privateKeyStr, err := utils.ReadDockerSecret("ssh_host_rsa_key_go_usa")
	if err != nil {
		return fmt.Errorf("failed to load private key: %v", err)
	}
	// Convert string to bytes and parse
	private, err := ssh.ParsePrivateKey([]byte(privateKeyStr))
	if err != nil {
		return fmt.Errorf("failed to parse private key: %v", err)
	}

	// Load the authorized public key from repo
	authorizedKeyBytes, err := os.ReadFile("./authorised/go_usa_stock.pub")
	if err != nil {
		return fmt.Errorf("failed to load authorized key: %v", err)
	}
	// Parse it
	authorizedPubKey, _, _, _, err := ssh.ParseAuthorizedKey(authorizedKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to parse authorized key: %v", err)
	}
	authorizedFingerprint := ssh.FingerprintSHA256(authorizedPubKey)
	fmt.Printf("✅ Loaded authorized key: %s\n", authorizedFingerprint)

	// just a reminder- sftp (file operations) > ssh (encryption) > tcp (network connection)

	config := &ssh.ServerConfig{
		PublicKeyCallback: func(c ssh.ConnMetadata, pubKey ssh.PublicKey) (*ssh.Permissions, error) {
			clientFingerprint := ssh.FingerprintSHA256(pubKey)

			if clientFingerprint == authorizedFingerprint {
				fmt.Printf("✅ Authorized user '%s' with key %s\n", c.User(), clientFingerprint)
				return &ssh.Permissions{
					Extensions: map[string]string{
						"pubkey-fp": clientFingerprint,
					},
				}, nil
			}

			fmt.Printf("❌ Rejected unauthorized user '%s' with key %s\n", c.User(), clientFingerprint)
			return nil, fmt.Errorf("unauthorized key")
		},
	}
	config.AddHostKey(private)

	listener, err := net.Listen("tcp", ":22")
	if err != nil {
		sentry.Notify(err, "failed to listen on sftp port")
		return fmt.Errorf("failed to listen: %v", err)
	}
	defer listener.Close()

	fmt.Println("SFTP server listening on :22 internally (2223)")

	// Step 4: Accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			sentry.Notify(err, "failed to accept connection")
			fmt.Printf("failed to accept connection: %v\n", err)
			continue
		}

		// ✅ ADD THIS LOG
		fmt.Printf("🔗 TCP connection accepted from %s\n", conn.RemoteAddr())

		// Handle each connection in a goroutine
		go handleConnection(conn, config)
	}
}

// this func is the ssh layer. its given a raw tcp connection via netConn
func handleConnection(netConn net.Conn, config *ssh.ServerConfig) {
	fmt.Printf("handle connection fun running")
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

	downloadsPath := "/app/downloads" // Absolute path
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
