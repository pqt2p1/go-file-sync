package sync

import "log"

func worker(id int, jobs <-chan SyncJob, syncer *FileSyncer) {
	for job := range jobs {
		log.Printf("Worker %d processing: %s", id, job.SourcePath)

		err := syncer.SyncFile(job.SourcePath, job.DestPath)
		if err != nil {
			log.Printf("Worker %d failed: %v", id, err)
		} else {
			log.Printf("Worker %d completed: %s", id, job.SourcePath)
		}
	}
}
