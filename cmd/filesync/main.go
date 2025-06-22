package main

import (
	"github.com/pqt2p1/go-file-sync/internal/sync"
	"github.com/pqt2p1/go-file-sync/internal/watcher"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	stdsync "sync"
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
	err = fileWatcher.WatchRecursive(sourceDir)
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

	// 5. Debounce map
	lastSync := make(map[string]time.Time)
	var mu stdsync.Mutex // Protect concurrent map access

	// 6. Listen to events trong goroutine
	go func() {
		for event := range fileWatcher.Events {
			log.Printf("DEBUG: Got event %s on %s", event.Operation, event.Path)
			if event.Operation == "create" {
				info, err := os.Stat(event.Path)
				if err == nil && info.IsDir() {
					log.Printf("New directory detected: %s", event.Path)
					fileWatcher.Watch(event.Path)
					continue
				}
			}

			if event.Operation == "write" || event.Operation == "create" {
				mu.Lock()
				if lastTime, exists := lastSync[event.Path]; exists {
					if time.Since(lastTime) < 500*time.Millisecond {
						mu.Unlock()
						continue
					}
				}
				lastSync[event.Path] = time.Now()
				mu.Unlock()

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

			if event.Operation == "delete" {
				// Remove from destination
				relPath, err := filepath.Rel(sourceDir, event.Path)
				if err != nil {
					log.Printf("Failed to get relative path: %v", err)
					continue
				}

				destPath := filepath.Join(destDir, relPath)

				// Check if it's file or directory
				if err := os.Remove(destPath); err != nil {
					if err := os.RemoveAll(destPath); err != nil {
						log.Printf("Failed to remove %s: %v", destPath, err)
					} else {
						log.Printf("Removed directory: %s", destPath)
					}
				} else {
					log.Printf("Removed file: %s", destPath)
				}
			}
		}
	}()

	// 7. Graceful shutdown
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
