package delivery

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/yourusername/webhookrelay/internal/models"
)

// Worker processes delivery jobs from the queue
type Worker struct {
	id             int
	jobQueue       <-chan *models.DeliveryJob
	db             DeliveryDB
	retryConfig    *RetryConfig
	requestTimeout time.Duration
	ctx            context.Context
}

// DeliveryDB interface for database operations
type DeliveryDB interface {
	LogDeliveryAttempt(attempt *models.DeliveryAttempt) error
	UpdateEventStatus(eventID string, status string) error
}

// NewWorker creates a new worker
func NewWorker(
	id int,
	jobQueue <-chan *models.DeliveryJob,
	db DeliveryDB,
	retryConfig *RetryConfig,
	requestTimeout time.Duration,
	ctx context.Context,
) *Worker {
	return &Worker{
		id:             id,
		jobQueue:       jobQueue,
		db:             db,
		retryConfig:    retryConfig,
		requestTimeout: requestTimeout,
		ctx:            ctx,
	}
}

// Start begins processing jobs
func (w *Worker) Start(requeue func(*models.DeliveryJob)) {
	go func() {
		log.Printf("Worker %d started", w.id)
		for {
			select {
			case <-w.ctx.Done():
				log.Printf("Worker %d shutting down", w.id)
				return
			case job := <-w.jobQueue:
				if job == nil {
					continue
				}
				w.processJob(job, requeue)
			}
		}
	}()
}

// processJob handles a single delivery job
func (w *Worker) processJob(job *models.DeliveryJob, requeue func(*models.DeliveryJob)) {
	// Wait until scheduled time
	if time.Now().Before(job.ScheduledAt) {
		delay := time.Until(job.ScheduledAt)
		log.Printf("Worker %d: Waiting %v before attempt %d for event %s",
			w.id, delay, job.AttemptNumber, job.Event.ID)
		
		select {
		case <-time.After(delay):
		case <-w.ctx.Done():
			return
		}
	}

	log.Printf("Worker %d: Processing attempt %d/%d for event %s to %s",
		w.id, job.AttemptNumber, w.retryConfig.MaxAttempts,
		job.Event.ID, job.Endpoint.URL)

	// Perform the HTTP request
	startTime := time.Now()
	statusCode, err := w.deliverWebhook(job)
	latency := time.Since(startTime)

	// Create delivery attempt record
	attempt := &models.DeliveryAttempt{
		EventID:       job.Event.ID,
		EndpointID:    job.Endpoint.ID,
		AttemptNumber: job.AttemptNumber,
		LatencyMs:     int(latency.Milliseconds()),
		CreatedAt:     time.Now(),
	}

	if statusCode > 0 {
		attempt.StatusCode = &statusCode
	}

	if err != nil {
		errMsg := err.Error()
		attempt.ErrorMessage = &errMsg
	}

	// Log the attempt
	if err := w.db.LogDeliveryAttempt(attempt); err != nil {
		log.Printf("Worker %d: Failed to log attempt: %v", w.id, err)
	}

	// Determine if we should retry
	success := statusCode >= 200 && statusCode < 300
	shouldRetry := !success && (ShouldRetry(statusCode) || IsRetryableError(err))
	canRetry := job.AttemptNumber < w.retryConfig.MaxAttempts

	if success {
		// Success - update event status
		log.Printf("Worker %d: Delivery successful (status %d, %dms)",
			w.id, statusCode, latency.Milliseconds())
		
		if err := w.db.UpdateEventStatus(job.Event.ID.String(), "delivered"); err != nil {
			log.Printf("Worker %d: Failed to update event status: %v", w.id, err)
		}
	} else if shouldRetry && canRetry {
		// Retry with exponential backoff
		nextAttempt := job.AttemptNumber + 1
		backoff := w.retryConfig.CalculateBackoff(nextAttempt)
		
		log.Printf("Worker %d: Attempt %d failed (status %d), retrying in %v (attempt %d/%d)",
			w.id, job.AttemptNumber, statusCode, backoff, nextAttempt, w.retryConfig.MaxAttempts)
		
		// Update event status to retrying
		if err := w.db.UpdateEventStatus(job.Event.ID.String(), "retrying"); err != nil {
			log.Printf("Worker %d: Failed to update event status: %v", w.id, err)
		}

		// Schedule retry
		retryJob := &models.DeliveryJob{
			Event:         job.Event,
			Endpoint:      job.Endpoint,
			AttemptNumber: nextAttempt,
			ScheduledAt:   time.Now().Add(backoff),
		}
		requeue(retryJob)
	} else {
		// Permanent failure - move to dead letter queue
		status := "failed"
		if job.AttemptNumber >= w.retryConfig.MaxAttempts {
			status = "dead_letter"
			log.Printf("Worker %d: Max attempts reached for event %s, moving to dead letter",
				w.id, job.Event.ID)
		} else {
			log.Printf("Worker %d: Non-retryable error (status %d) for event %s",
				w.id, statusCode, job.Event.ID)
		}
		
		if err := w.db.UpdateEventStatus(job.Event.ID.String(), status); err != nil {
			log.Printf("Worker %d: Failed to update event status: %v", w.id, err)
		}
	}
}

// deliverWebhook performs the actual HTTP POST request
func (w *Worker) deliverWebhook(job *models.DeliveryJob) (int, error) {
	// Create request context with timeout
	ctx, cancel := context.WithTimeout(w.ctx, w.requestTimeout)
	defer cancel()

	// Sign the payload
	signature := SignPayload(job.Event.Payload, job.Endpoint.Secret)
	timestamp := time.Now().Unix()

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", job.Endpoint.URL, bytes.NewReader(job.Event.Payload))
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Webhook-Signature", signature)
	req.Header.Set("X-Webhook-ID", job.Event.ID.String())
	req.Header.Set("X-Webhook-Timestamp", fmt.Sprintf("%d", timestamp))
	req.Header.Set("X-Webhook-Event-Type", job.Event.EventType)

	// Perform request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body (for logging purposes)
	_, _ = io.ReadAll(resp.Body)

	return resp.StatusCode, nil
}
