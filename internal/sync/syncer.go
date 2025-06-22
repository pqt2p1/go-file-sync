package sync

import (
	"fmt"
	"github.com/pqt2p1/go-file-sync/internal/checksum"
	"io"
	"os"
	"path/filepath"
)

type FileSyncer struct {
}

func NewFileSyncer() *FileSyncer {
	return &FileSyncer{}
}

func (fs *FileSyncer) SyncFile(src, dest string, progress *Progress) error {
	// 1. Check if destination exists
	_, err := os.Stat(dest)
	if err == nil {
		// File exists - check if same
		same, err := checksum.CompareFiles(src, dest)
		if err != nil {
			return fmt.Errorf("failed to compare: %v", err)
		}
		if same {
			return nil
		}
	}

	// 2. Need to copy (not exist or different )
	if err := fs.copyFile(src, dest); err != nil {
		return fmt.Errorf("failed to copy: %v", err)
	}

	// 3. Verify after copy
	same, err := checksum.CompareFiles(src, dest)
	if err != nil {
		return fmt.Errorf("failed to compare: %v", err)
	}
	if !same {
		return fmt.Errorf("copy verification failed: %v", dest)
	}

	fileInfo, err := os.Stat(src)
	if err == nil {
		progress.AddBytes(fileInfo.Size())
	}
	return nil
}

func (fs *FileSyncer) copyFile(src, dest string) error {
	// Create parent dir if needed
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destDir, err)
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copy content
	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	// Sync to disk
	return destFile.Sync()
}
