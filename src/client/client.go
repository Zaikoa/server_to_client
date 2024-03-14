package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "network_ipv4:8080")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	// Send the file name to the server
	fileName := os.Args[1]
	_, err = conn.Write([]byte(fileName))
	if err != nil {
		fmt.Println("Error sending file name:", err)
		return
	}

	// Open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	gzipWriter := gzip.NewWriter(conn)
	defer gzipWriter.Close()

	// Copy the file content from the file to the gzip writer
	_, err = io.Copy(gzipWriter, file)
	if err != nil {
		fmt.Println("Error copying file content:", err)
		return
	}

	fmt.Println("File sent successfully.")
}
