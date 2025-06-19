package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ./filesync <source> <dest>")
		os.Exit(1)
	}

	src, dst := os.Args[1], os.Args[2]

	if err := copyFile(src, dst); err != nil {
		fmt.Println("Error copying file:", err)
		os.Exit(1)
	}

	fmt.Println("Successfully copied file", src, "to", dst)
}

func copyFile(src, dst string) error {
	// Open source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	// Create destination file
	destFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	return nil
}
