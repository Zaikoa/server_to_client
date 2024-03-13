package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		fmt.Println("Error creating listener:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Listening on :0.0.0.0:8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		fmt.Println("Handling connection from", conn.RemoteAddr())
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Receive the file name
	path, err := receiveString(conn)
	if err != nil {
		fmt.Println("Error receiving file name:", err)
		return
	}
	split := strings.Split(path, "/")
	fileName := split[len(split)-1]
	fmt.Println("Received file name:", fileName)
	zipNmae := strings.Split(fileName, ".")

	zipFilePath := filepath.Join("../../uploads", zipNmae[0]+".zip")
	zipFilePath = strings.ReplaceAll(zipFilePath, "\\", "/") // windows moment

	zipFile, err := os.Create(zipFilePath)
	if err != nil {
		fmt.Println("Error creating zip file:", err)
		return
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	zipFileWriter, err := zipWriter.Create(fileName)
	if err != nil {
		fmt.Println("Error creating zip file entry:", err)
		return
	}

	_, err = io.Copy(zipFileWriter, file)
	if err != nil {
		fmt.Println("Error copying file content to zip:", err)
		return
	}

	fmt.Println("File compressed and stored in zip at:", zipFilePath)
}

func receiveString(conn net.Conn) (string, error) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}
