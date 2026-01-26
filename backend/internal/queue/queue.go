package queue

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Job represents a scan job
type Job struct {
	ID       string
	Domain   string
	DomainID string
	Status   string // "pending", "running", "completed", "failed"
	Result   interface{}
	Error    error
	CreatedAt time.Time
	StartedAt *time.Time
	FinishedAt *time.Time
}

// Queue is a simple in-memory job queue
type Queue struct {
	jobs    map[string]*Job
	mu      sync.RWMutex
	workers chan *Job
}

var (
	globalQueue *Queue
	once        sync.Once
)

// GetQueue returns the global queue instance
func GetQueue() *Queue {
	once.Do(func() {
		globalQueue = &Queue{
			jobs:    make(map[string]*Job),
			workers: make(chan *Job, 100), // Buffer for 100 jobs
		}
	})
	return globalQueue
}

// Enqueue adds a job to the queue
func (q *Queue) Enqueue(domain, domainID string) string {
	jobID := uuid.New().String()
	job := &Job{
		ID:        jobID,
		Domain:    domain,
		DomainID:  domainID,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	q.mu.Lock()
	q.jobs[jobID] = job
	q.mu.Unlock()

	// Send to worker channel
	select {
	case q.workers <- job:
	default:
		// Channel full, but job is still queued
	}

	return jobID
}

// GetJob retrieves a job by ID
func (q *Queue) GetJob(jobID string) (*Job, bool) {
	q.mu.RLock()
	defer q.mu.RUnlock()
	job, exists := q.jobs[jobID]
	return job, exists
}

// UpdateJobStatus updates a job's status
func (q *Queue) UpdateJobStatus(jobID, status string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if job, exists := q.jobs[jobID]; exists {
		job.Status = status
		if status == "running" && job.StartedAt == nil {
			now := time.Now()
			job.StartedAt = &now
		} else if (status == "completed" || status == "failed") && job.FinishedAt == nil {
			now := time.Now()
			job.FinishedAt = &now
		}
	}
}

// SetJobResult sets the result of a completed job
func (q *Queue) SetJobResult(jobID string, result interface{}, err error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if job, exists := q.jobs[jobID]; exists {
		job.Result = result
		job.Error = err
		if err != nil {
			job.Status = "failed"
		} else {
			job.Status = "completed"
		}
		now := time.Now()
		job.FinishedAt = &now
	}
}

// GetNextJob retrieves the next pending job (for workers)
func (q *Queue) GetNextJob(ctx context.Context) (*Job, bool) {
	select {
	case job := <-q.workers:
		return job, true
	case <-ctx.Done():
		return nil, false
	}
}

// ListJobs returns all jobs (for debugging/admin)
func (q *Queue) ListJobs() []*Job {
	q.mu.RLock()
	defer q.mu.RUnlock()
	jobs := make([]*Job, 0, len(q.jobs))
	for _, job := range q.jobs {
		jobs = append(jobs, job)
	}
	return jobs
}
