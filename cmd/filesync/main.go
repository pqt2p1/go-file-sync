package main

import (
	"github.com/pqt2p1/go-file-sync/internal/sync"
	"github.com/pqt2p1/go-file-sync/internal/watcher"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 4 || os.Args[1] != "watch" {
		log.Fatal("Usage: ./filesync watch <source> <dest> ")
	}

	sourceDir, err := filepath.Abs(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	destDir, err := filepath.Abs(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	// 1. Tạo FileWatcher
	fileWatcher, err := watcher.NewFileWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Watch path
	err = fileWatcher.Watch(sourceDir)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Create worker pool
	workerPool := sync.NewWorkerPool(5)
	workerPool.Start()

	// 4. Start progress reporter
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			progress := workerPool.GetProgress()
			progress.PrintProgress()
		}
	}()

	// 4. Listen to events trong goroutine
	go func() {
		for event := range fileWatcher.Events {
			if event.Operation == "write" {
				log.Printf("Got event: %s on %s", event.Operation, event.Path) // Debug log!

				// Get relative path from source
				relPath, err := filepath.Rel(sourceDir, event.Path)
				if err != nil {
					log.Printf("Failed to get relative path: %v", err)
					continue
				}

				// Build destination path
				destPath := filepath.Join(destDir, relPath)

				// Submit job to pool
				workerPool.SubmitJob(event.Path, destPath)
			}
		}
	}()

	// 4. Wait for Ctrl+C
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	log.Println("Watching... Press Ctrl+C to stop")
	<-sigChan

	// 5. Cleanup
	log.Println("Stopping...")
	err = fileWatcher.Close()
	if err != nil {
		log.Fatal(err)
	}
}
