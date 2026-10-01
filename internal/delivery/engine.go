package delivery

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/webhookrelay/internal/models"
)

// Engine manages the delivery worker pool and job queue
type Engine struct {
	jobQueue    chan *models.DeliveryJob
	workers     []*Worker
	db          DeliveryDB
	retryConfig *RetryConfig
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewEngine creates a new delivery engine
func NewEngine(
	workerCount int,
	queueSize int,
	db DeliveryDB,
	retryConfig *RetryConfig,
	requestTimeout time.Duration,
) *Engine {
	ctx, cancel := context.WithCancel(context.Background())
	
	jobQueue := make(chan *models.DeliveryJob, queueSize)
	workers := make([]*Worker, workerCount)

	for i := 0; i < workerCount; i++ {
		workers[i] = NewWorker(i+1, jobQueue, db, retryConfig, requestTimeout, ctx)
	}

	return &Engine{
		jobQueue:    jobQueue,
		workers:     workers,
		db:          db,
		retryConfig: retryConfig,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Start starts all workers in the pool
func (e *Engine) Start() {
	log.Printf("Starting delivery engine with %d workers", len(e.workers))
	
	for _, worker := range e.workers {
		worker.Start(e.Enqueue)
	}
	
	log.Println("Delivery engine started successfully")
}

// Enqueue adds a delivery job to the queue
func (e *Engine) Enqueue(job *models.DeliveryJob) {
	select {
	case e.jobQueue <- job:
		log.Printf("Job enqueued: event %s to %s (attempt %d)",
			job.Event.ID, job.Endpoint.URL, job.AttemptNumber)
	case <-e.ctx.Done():
		log.Println("Engine is shutting down, job not enqueued")
	default:
		log.Println("Warning: Job queue is full, dropping job")
	}
}

// Shutdown gracefully shuts down the engine
func (e *Engine) Shutdown() {
	log.Println("Shutting down delivery engine...")
	e.cancel()
	close(e.jobQueue)
	
	// Wait a bit for workers to finish current jobs
	time.Sleep(2 * time.Second)
	
	log.Println("Delivery engine shut down successfully")
}

// QueueDelivery creates and enqueues a delivery job
func (e *Engine) QueueDelivery(event *models.Event, endpoint *models.Endpoint) {
	job := &models.DeliveryJob{
		Event:         event,
		Endpoint:      endpoint,
		AttemptNumber: 1,
		ScheduledAt:   time.Now(),
	}
	e.Enqueue(job)
}
