package fetcher

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"m/utils"
	"os"
	"path/filepath"
	"sync"
	"time"

	"m/sentry"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// downloadMutex prevents multiple simultaneous downloads
var downloadMutex sync.Mutex
var downloadInProgress sync.Map // tracks if download is in progress for a file

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
	return nil
}

// EnsureFresh checks if file is fresh and triggers background download if stale
// Returns immediately - does not block
func EnsureFresh(maxAge time.Duration) {
	env, err := utils.GetEnv()
	if err != nil {
		fmt.Printf("Warning: failed to get env for freshness check: %v\n", err)
		return
	}
	filename := env["REMOTE_FILENAME"]
	finalFilePath := filepath.Join("downloads", filename)

	// Check if file exists and is fresh
	info, err := os.Stat(finalFilePath)
	if err == nil {
		age := time.Since(info.ModTime())
		if age < maxAge {
			fmt.Printf("File %s is fresh (age: %v), no download needed\n", filename, age)
			return
		}
		fmt.Printf("File %s is stale (age: %v), triggering background refresh\n", filename, age)
	} else {
		fmt.Printf("File %s does not exist, triggering download\n", filename)
	}

	// Check if download is already in progress
	if _, downloading := downloadInProgress.LoadOrStore(filename, true); downloading {
		fmt.Printf("Download for %s already in progress, skipping\n", filename)
		return
	}

	// Trigger background download (non-blocking)
	go func() {
		defer downloadInProgress.Delete(filename)
		if err := DlSanmar(); err != nil {
			errorMessage := fmt.Sprintf("background download of %s failed from sanmar", filename)
			sentry.Notify(err, errorMessage)
			fmt.Printf("Background download failed: %v\n", err)
		}
	}()
}

func DlSanmar() error {
	// Acquire lock to prevent concurrent downloads
	downloadMutex.Lock()
	defer downloadMutex.Unlock()

	start := time.Now()
	fmt.Printf("download started at %s\n", start.Format(time.RFC1123))
	env, err := utils.GetEnv()
	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}
	path := env["REMOTE_DIR"]
	filename := env["REMOTE_FILENAME"]

	// Ensure downloads directory exists
	if err := os.MkdirAll("downloads", 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %w", err)
	}

	client, err := connect()
	if err != nil {
		sentry.Notify(err, "failed to connect to sanmar ftp")
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	remoteFilePath := filepath.Join(path, filename)
	tempFilePath := filepath.Join("downloads", filename+".tmp")
	finalFilePath := filepath.Join("downloads", filename)

	// Step 1: Download the entire file efficiently to a temp file
	fmt.Println("Downloading file from SFTP...")
	downloadStart := time.Now()

	// Open remote file
	remoteFile, err := client.Open(remoteFilePath)
	if err != nil {
		sentry.Notify(err, "failed to open remote file on sanmar ftp")
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	// Create temp file for download
	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		sentry.Notify(err, "failed to create temp file")
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	written, err := io.Copy(tempFile, remoteFile)
	if err != nil {
		tempFile.Close()
		os.Remove(tempFilePath) // cleanup on error
		sentry.Notify(err, "failed to download file from sanmar ftp")
		return fmt.Errorf("failed to download file: %w", err)
	}
	tempFile.Close()

	downloadEnd := time.Now()
	fmt.Printf("Downloaded %d bytes in %v\n", written, downloadEnd.Sub(downloadStart))

	// Step 2: Process the downloaded file locally
	fmt.Println("Processing and filtering CSV...")
	processStart := time.Now()

	// Open the temp file for reading
	tempFileReader, err := os.Open(tempFilePath)
	if err != nil {
		os.Remove(tempFilePath) // cleanup on error
		sentry.Notify(err, "failed to open temp file for processing")
		return fmt.Errorf("failed to open temp file: %w", err)
	}
	defer tempFileReader.Close()

	// Create processed output file with different temp name
	// This ensures atomic replacement - we never overwrite the final file until it's complete
	processedTempPath := filepath.Join("downloads", filename+".processed.tmp")
	processedFile, err := os.Create(processedTempPath)
	if err != nil {
		os.Remove(tempFilePath) // cleanup on error
		sentry.Notify(err, "failed to create processed temp file")
		return fmt.Errorf("failed to create processed temp file: %w", err)
	}
	defer processedFile.Close()

	// Use buffered writer for output
	bufferedWriter := bufio.NewWriterSize(processedFile, 16*1024)
	defer bufferedWriter.Flush()

	// Define which columns you want to keep
	columnsToKeep := []string{"Variant SKU", "Variant Inventory Qty"}

	// Parse and filter CSV from local temp file
	if err := filterCSV(tempFileReader, bufferedWriter, columnsToKeep); err != nil {
		os.Remove(tempFilePath)      // cleanup on error
		os.Remove(processedTempPath) // cleanup on error
		sentry.Notify(err, "failed to filter csv")
		return fmt.Errorf("failed to filter CSV: %w", err)
	}

	// Ensure all data is written to disk before rename
	bufferedWriter.Flush()
	processedFile.Close()

	processEnd := time.Now()
	fmt.Printf("Processed CSV in %v\n", processEnd.Sub(processStart))

	// Step 3: Atomic rename - this is instant and atomic on Unix/Linux
	// Active file handles continue to work with the old file content
	// New opens will get the new file
	if err := os.Rename(processedTempPath, finalFilePath); err != nil {
		os.Remove(tempFilePath)      // cleanup on error
		os.Remove(processedTempPath) // cleanup on error
		sentry.Notify(err, "failed to rename processed file")
		return fmt.Errorf("failed to rename processed file: %w", err)
	}

	// Clean up raw temp file
	if err := os.Remove(tempFilePath); err != nil {
		fmt.Printf("Warning: failed to remove temp file: %v\n", err)
	}

	end := time.Now()
	fmt.Printf("Downloaded and filtered %s successfully to %s at %s (total time: %v)\n",
		filename, finalFilePath, end.Format(time.RFC1123), end.Sub(start))
	return nil
}

func connect() (*sftp.Client, error) {

	env, err := utils.GetEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to get env: %w", err)
	}

	url := env["REMOTE_URL"]
	port := env["REMOTE_PORT"]
	username := env["REMOTE_USERNAME"]
	password := env["REMOTE_PASSWORD"]
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

func filterCSV(reader io.Reader, writer io.Writer, columnsToKeep []string) error {
	csvReader := csv.NewReader(reader)
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// Read header
	header, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Find indices of columns to keep
	columnIndices := make([]int, 0)
	for _, colName := range columnsToKeep {
		for i, h := range header {
			if h == colName {
				columnIndices = append(columnIndices, i)
				break
			}
		}
	}

	// Write filtered header
	filteredHeader := make([]string, len(columnIndices))
	for i, idx := range columnIndices {
		filteredHeader[i] = header[idx]
	}
	if err := csvWriter.Write(filteredHeader); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Process rows
	rowCount := 0
	for {
		row, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read row: %w", err)
		}

		// Extract only the columns we want
		filteredRow := make([]string, len(columnIndices))
		for i, idx := range columnIndices {
			if idx < len(row) {
				filteredRow[i] = row[idx]
			}
		}
		rowCount++
		if rowCount%5000 == 0 {
			fmt.Printf("Processed %d rows...\n", rowCount)
		}

		if err := csvWriter.Write(filteredRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	fmt.Printf("Total rows processed: %d\n", rowCount)
	return nil
}
