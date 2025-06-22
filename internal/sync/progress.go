package sync

import (
	"fmt"
	"sync/atomic"
	"time"
)

type Progress struct {
	totalFiles       int64
	completedFiles   int64
	failedFiles      int64
	bytesTransferred int64
	startTime        time.Time
}

func NewProgress() *Progress {
	return &Progress{
		startTime: time.Now(),
	}
}

// Getter
func (p *Progress) GetCompleted() int64 {
	return atomic.LoadInt64(&p.completedFiles)
}

func (p *Progress) GetFailed() int64 {
	return atomic.LoadInt64(&p.failedFiles)
}

func (p *Progress) GetBytesTransferred() int64 {
	return atomic.LoadInt64(&p.bytesTransferred)
}

func (p *Progress) GetTotalFiles() int64 {
	return atomic.LoadInt64(&p.totalFiles)
}

func (p *Progress) PrintProgress() {
	completed := atomic.LoadInt64(&p.completedFiles)
	failed := atomic.LoadInt64(&p.failedFiles)
	bytes := atomic.LoadInt64(&p.bytesTransferred)

	// Calculate MB
	mb := float64(bytes) / 1024 / 1024

	// Simple version first - no bar yet
	fmt.Printf("\rFiles: %d completed, %d failed | %.2f MB",
		completed, failed, mb)
}

// Thread-safe methods
func (p *Progress) IncrementCompleted() {
	atomic.AddInt64(&p.completedFiles, 1)
}

func (p *Progress) IncrementFailed() {
	atomic.AddInt64(&p.failedFiles, 1)
}

func (p *Progress) AddBytes(bytes int64) {
	atomic.AddInt64(&p.bytesTransferred, bytes)
}
