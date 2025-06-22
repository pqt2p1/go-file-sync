package sync

type WorkerPool struct {
	numWorkers int
	jobs       chan SyncJob
	syncer     *FileSyncer
}

func NewWorkerPool(numWorkers int) *WorkerPool {
	return &WorkerPool{
		numWorkers: numWorkers,
		jobs:       make(chan SyncJob, 100), // Buffer 100 jobs
		syncer:     NewFileSyncer(),
	}
}

func (wp *WorkerPool) Start() {
	for i := 1; i <= wp.numWorkers; i++ {
		go worker(i, wp.jobs, wp.syncer)
	}
}

func (wp *WorkerPool) SubmitJob(sourcePath, destPath string) {
	job := SyncJob{
		SourcePath: sourcePath,
		DestPath:   destPath,
	}
	wp.jobs <- job
}
