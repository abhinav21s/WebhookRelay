package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// StringArray is a custom type for PostgreSQL text[] arrays
type StringArray []string

// Scan implements the sql.Scanner interface for StringArray
func (a *StringArray) Scan(src interface{}) error {
	if src == nil {
		*a = []string{}
		return nil
	}
	
	switch v := src.(type) {
	case []byte:
		// PostgreSQL returns arrays as: {value1,value2,value3}
		return a.parsePostgresArray(string(v))
	case string:
		// Handle string format
		return a.parsePostgresArray(v)
	case []interface{}:
		// Handle if it comes as a slice of interfaces
		*a = make([]string, len(v))
		for i, val := range v {
			if str, ok := val.(string); ok {
				(*a)[i] = str
			}
		}
		return nil
	}
	
	return nil
}

// parsePostgresArray parses PostgreSQL array format: {value1,value2}
func (a *StringArray) parsePostgresArray(s string) error {
	// Remove surrounding braces
	if len(s) < 2 {
		*a = []string{}
		return nil
	}
	
	if s[0] == '{' && s[len(s)-1] == '}' {
		s = s[1 : len(s)-1]
	}
	
	if s == "" {
		*a = []string{}
		return nil
	}
	
	// Split by comma (simple parsing - doesn't handle quoted commas)
	parts := []string{}
	current := ""
	inQuotes := false
	
	for i := 0; i < len(s); i++ {
		c := s[i]
		
		if c == '"' {
			inQuotes = !inQuotes
			continue
		}
		
		if c == ',' && !inQuotes {
			if current != "" {
				parts = append(parts, current)
				current = ""
			}
			continue
		}
		
		current += string(c)
	}
	
	if current != "" {
		parts = append(parts, current)
	}
	
	*a = parts
	return nil
}

// Value implements the driver.Valuer interface for StringArray
func (a StringArray) Value() (driver.Value, error) {
	if a == nil || len(a) == 0 {
		return "{}", nil
	}
	
	// Return as PostgreSQL array format: {value1,value2}
	result := "{"
	for i, s := range a {
		if i > 0 {
			result += ","
		}
		// Quote values that contain special characters
		if containsSpecialChars(s) {
			result += `"` + s + `"`
		} else {
			result += s
		}
	}
	result += "}"
	
	return result, nil
}

func containsSpecialChars(s string) bool {
	for _, c := range s {
		if c == ',' || c == '{' || c == '}' || c == '"' || c == '\\' || c == ' ' {
			return true
		}
	}
	return false
}

// Endpoint represents a webhook endpoint
type Endpoint struct {
	ID         uuid.UUID   `json:"id"`
	Name       string      `json:"name"`
	URL        string      `json:"url"`
	Secret     string      `json:"secret"`
	EventTypes StringArray `json:"event_types"`
	IsActive   bool        `json:"is_active"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
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
