package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Endpoint represents a webhook endpoint
type Endpoint struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	URL        string         `json:"url"`
	Secret     string         `json:"secret"`
	EventTypes pq.StringArray `json:"event_types"`
	IsActive   bool           `json:"is_active"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
}

// Event represents a webhook event
type Event struct {
	ID        uuid.UUID       `json:"id"`
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
	Status    string          `json:"status"` // pending, delivering, delivered, failed, dead_letter
	CreatedAt time.Time       `json:"created_at"`
}

// DeliveryAttempt represents a single delivery attempt
type DeliveryAttempt struct {
	ID            uuid.UUID  `json:"id"`
	EventID       uuid.UUID  `json:"event_id"`
	EndpointID    uuid.UUID  `json:"endpoint_id"`
	AttemptNumber int        `json:"attempt_number"`
	StatusCode    *int       `json:"status_code"`
	LatencyMs     int        `json:"latency_ms"`
	ErrorMessage  *string    `json:"error_message"`
	CreatedAt     time.Time  `json:"created_at"`
}

// DeliveryJob represents a delivery job in the queue
type DeliveryJob struct {
	Event          *Event
	Endpoint       *Endpoint
	AttemptNumber  int
	ScheduledAt    time.Time
}

// EventWithDeliveries includes event and its delivery attempts
type EventWithDeliveries struct {
	Event     *Event             `json:"event"`
	Deliveries []*DeliveryAttempt `json:"deliveries"`
}

// Metrics represents aggregated metrics
type Metrics struct {
	TotalEvents24h    int     `json:"total_events_24h"`
	SuccessRate       float64 `json:"success_rate"`
	AverageLatencyMs  float64 `json:"average_latency_ms"`
	ActiveEndpoints   int     `json:"active_endpoints"`
}

// CreateEndpointRequest represents the request to create an endpoint
type CreateEndpointRequest struct {
	Name       string   `json:"name"`
	URL        string   `json:"url"`
	EventTypes []string `json:"event_types"`
	Secret     string   `json:"secret,omitempty"`
}

// UpdateEndpointRequest represents the request to update an endpoint
type UpdateEndpointRequest struct {
	Name       *string  `json:"name,omitempty"`
	URL        *string  `json:"url,omitempty"`
	EventTypes []string `json:"event_types,omitempty"`
	IsActive   *bool    `json:"is_active,omitempty"`
}

// CreateEventRequest represents the request to create an event
type CreateEventRequest struct {
	EventType string          `json:"event_type"`
	Payload   json.RawMessage `json:"payload"`
}

// JSONB is a custom type for PostgreSQL JSONB
type JSONB map[string]interface{}

// Value implements the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan implements the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}
