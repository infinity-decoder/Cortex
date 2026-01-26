package queue

import (
	"context"
	"log"

	"cortex-backend/internal/scanner"
)

// Worker processes scan jobs from the queue
type Worker struct {
	queue       *Queue
	orchestrator *scanner.Orchestrator
}

// NewWorker creates a new worker
func NewWorker(queue *Queue, orchestrator *scanner.Orchestrator) *Worker {
	return &Worker{
		queue:       queue,
		orchestrator: orchestrator,
	}
}

// Start begins processing jobs from the queue
func (w *Worker) Start(ctx context.Context) {
	log.Println("[Worker] Starting scan worker...")
	
	for {
		select {
		case <-ctx.Done():
			log.Println("[Worker] Stopping worker...")
			return
		default:
			job, ok := w.queue.GetNextJob(ctx)
			if !ok {
				continue
			}

			// Process the job
			w.processJob(ctx, job)
		}
	}
}

func (w *Worker) processJob(ctx context.Context, job *Job) {
	log.Printf("[Worker] Processing job %s for domain %s", job.ID, job.Domain)
	
	// Update status to running
	w.queue.UpdateJobStatus(job.ID, "running")

	// Run the scan
	result, err := w.orchestrator.RunScan(ctx, job.Domain, job.DomainID)
	
	// Set result
	w.queue.SetJobResult(job.ID, result, err)
	
	if err != nil {
		log.Printf("[Worker] Job %s failed: %v", job.ID, err)
	} else {
		log.Printf("[Worker] Job %s completed successfully", job.ID)
	}
}
