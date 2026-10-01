package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

// MockReceiver is a controllable webhook receiver for testing
type MockReceiver struct {
	mode string
	mu   sync.RWMutex
}

type MockControlRequest struct {
	Mode string `json:"mode"` // success, error500, timeout, rate_limit, bad_request
}

func NewMockReceiver() *MockReceiver {
	return &MockReceiver{
		mode: "success",
	}
}

// HandleWebhook receives webhook POSTs and responds based on current mode
func (m *MockReceiver) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	currentMode := m.mode
	m.mu.RUnlock()

	// Read and log the webhook
	body, _ := io.ReadAll(r.Body)
	signature := r.Header.Get("X-Webhook-Signature")
	eventID := r.Header.Get("X-Webhook-ID")
	timestamp := r.Header.Get("X-Webhook-Timestamp")

	log.Printf("[Mock Receiver] Received webhook - Event: %s, Signature: %s, Timestamp: %s, Body: %s",
		eventID, signature, timestamp, string(body))

	// Respond based on mode
	switch currentMode {
	case "error500":
		log.Printf("[Mock Receiver] Responding with 500 Internal Server Error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))

	case "timeout":
		log.Printf("[Mock Receiver] Simulating timeout (sleeping 10 seconds)")
		time.Sleep(10 * time.Second)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))

	case "rate_limit":
		log.Printf("[Mock Receiver] Responding with 429 Rate Limited")
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte(`{"error": "Rate limit exceeded"}`))

	case "bad_request":
		log.Printf("[Mock Receiver] Responding with 400 Bad Request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Bad request"}`))

	default: // success
		log.Printf("[Mock Receiver] Responding with 200 OK")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "received", "event_id": "` + eventID + `"}`))
	}
}

// SetMode changes the mock receiver's behavior mode
func (m *MockReceiver) SetMode(w http.ResponseWriter, r *http.Request) {
	var req MockControlRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	m.mu.Lock()
	m.mode = req.Mode
	m.mu.Unlock()

	log.Printf("[Mock Receiver] Mode changed to: %s", req.Mode)
	respondJSON(w, http.StatusOK, map[string]string{"mode": req.Mode})
}

// GetMode returns the current mode
func (m *MockReceiver) GetMode(w http.ResponseWriter, r *http.Request) {
	m.mu.RLock()
	currentMode := m.mode
	m.mu.RUnlock()

	respondJSON(w, http.StatusOK, map[string]string{"mode": currentMode})
}
