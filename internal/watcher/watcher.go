package watcher

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"os"
	"path/filepath"
)

type FileWatcher struct {
	watcher *fsnotify.Watcher
	Events  chan FileEvent
}

type FileEvent struct {
	Path      string
	Operation string
}

func NewFileWatcher() (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		watcher: watcher,
		Events:  make(chan FileEvent, 100),
	}

	// Start goroutine để convert fsnotify events -> FileEvent
	go fw.processEvents()

	return fw, nil
}

func (fw *FileWatcher) processEvents() {
	for {
		select {
		case event := <-fw.watcher.Events:
			var operation string

			if event.Op&fsnotify.Write == fsnotify.Write {
				operation = "write"
			} else if event.Op&fsnotify.Create == fsnotify.Create {
				operation = "create"
			} else if event.Op&fsnotify.Remove == fsnotify.Remove {
				operation = "delete"
			}

			fw.Events <- FileEvent{
				Path:      event.Name,
				Operation: operation,
			}
		case err := <-fw.watcher.Errors:
			log.Println("Error:", err)
		}
	}
}

func (fw *FileWatcher) Watch(path string) error {
	// Validate file path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	log.Printf("Watching: %s", absPath)
	return fw.watcher.Add(absPath)
}

func (fw *FileWatcher) WatchRecursive(root string) error {
	// Walk through all directories
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Only watch directories
		if info.IsDir() {
			log.Printf("Adding watch: %s", path)
			if err := fw.watcher.Add(path); err != nil {
				log.Printf("Failed to watch %s: %v", path, err)
			}
		}
		return nil
	})

	return err
}

func (fw *FileWatcher) Close() error {
	err := fw.watcher.Close()
	close(fw.Events)
	return err
}
