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
	"time"

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
	return nil
}

func DlSanmar() error {
	start := time.Now()
	fmt.Printf("download started at %s\n", start.Format(time.RFC1123))
	env, err := utils.GetEnv()
	if err != nil {
		return fmt.Errorf("failed to get env: %w", err)
	}
	path := env["DIR"]
	filename := env["FILENAME"]

	// Ensure downloads directory exists
	if err := os.MkdirAll("downloads", 0755); err != nil {
		return fmt.Errorf("failed to create downloads directory: %w", err)
	}

	client, err := connect()
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer client.Close()

	remoteFilePath := filepath.Join(path, filename)
	localFilePath := filepath.Join("downloads", filename)

	// Open the remote file for reading
	remoteFile, err := client.Open(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %w", err)
	}
	defer remoteFile.Close()

	// Create the local file for writing
	localFile, err := os.Create(localFilePath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %w", err)
	}
	// batch writes using an in memory buffer, sized at 16KB. default was 4KB
	// accumulate writes in memory, then "flushes" automatically when full and does 1 x write. instead of 1 x write for each row
	bufferedWriter := bufio.NewWriterSize(localFile, 16*1024)
	
	// anonymous function that runs at the end to cleanup- capture the last unwritten rows and close the file 
	defer func() {
		bufferedWriter.Flush()
		localFile.Close()
	}()

	// Define which columns you want to keep (by name or index)
	columnsToKeep := []string{"Variant SKU", "Variant Inventory Qty"} // adjust to your needs

	// Parse and filter CSV
	if err := filterCSV(remoteFile, bufferedWriter, columnsToKeep); err != nil {
		return fmt.Errorf("failed to filter CSV: %w", err)
	}

	end:= time.Now()
	fmt.Printf("Downloaded and filtered %s successfully to %s at %s\n", filename, localFilePath, end.Format(time.RFC1123))
	return nil

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
		if rowCount%10000 == 0 {
			fmt.Printf("Processed %d rows...\n", rowCount)
		}

		if err := csvWriter.Write(filteredRow); err != nil {
			return fmt.Errorf("failed to write row: %w", err)
		}
	}
	fmt.Printf("Total rows processed: %d\n", rowCount)
	return nil
}
