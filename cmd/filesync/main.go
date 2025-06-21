package main

import (
	"github.com/pqt2p1/go-file-sync/internal/watcher"
	"log"
	"os"
	"os/signal"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "watch" {
		log.Fatal("Usage: ./filesync watch <path>")
	}

	path := os.Args[2]

	// 1. Tạo FileWatcher
	fileWatcher, err := watcher.NewFileWatcher()
	if err != nil {
		log.Fatal(err)
	}

	// 2. Watch path
	err = fileWatcher.Watch(path)
	if err != nil {
		log.Fatal(err)
	}

	// 3. Listen to events trong goroutine
	go func() {
		for event := range fileWatcher.Events {
			log.Printf("File %s: %s", event.Operation, event.Path)
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
