package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourusername/webhookrelay/internal/delivery"
	"github.com/yourusername/webhookrelay/internal/models"
)

type EventHandler struct {
	db     *sql.DB
	engine *delivery.Engine
}

func NewEventHandler(db *sql.DB, engine *delivery.Engine) *EventHandler {
	return &EventHandler{
		db:     db,
		engine: engine,
	}
}

// CreateEvent publishes a new event
func (h *EventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if req.EventType == "" {
		respondError(w, http.StatusBadRequest, "Event type is required")
		return
	}

	// Create event
	event := &models.Event{
		ID:        uuid.New(),
		EventType: req.EventType,
		Payload:   req.Payload,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	// Insert into database
	query := `
		INSERT INTO events (id, event_type, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := h.db.Exec(query,
		event.ID, event.EventType, event.Payload, event.Status, event.CreatedAt,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create event")
		return
	}

	// Find matching active endpoints
	endpoints, err := h.findMatchingEndpoints(event.EventType)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to find endpoints")
		return
	}

	// Queue delivery jobs for each matching endpoint
	if len(endpoints) > 0 {
		// Update event status to delivering
		h.db.Exec("UPDATE events SET status = 'delivering' WHERE id = $1", event.ID)
		
		for _, endpoint := range endpoints {
			h.engine.QueueDelivery(event, endpoint)
		}
	}

	respondJSON(w, http.StatusCreated, event)
}

// ListEvents returns all events
func (h *EventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, event_type, payload, status, created_at
		FROM events
		ORDER BY created_at DESC
		LIMIT 100
	`
	rows, err := h.db.Query(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch events")
		return
	}
	defer rows.Close()

	events := []*models.Event{}
	for rows.Next() {
		var ev models.Event
		err := rows.Scan(&ev.ID, &ev.EventType, &ev.Payload, &ev.Status, &ev.CreatedAt)
		if err != nil {
			continue
		}
		events = append(events, &ev)
	}

	respondJSON(w, http.StatusOK, events)
}

// GetEvent returns a single event with its delivery attempts
func (h *EventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get event
	query := `
		SELECT id, event_type, payload, status, created_at
		FROM events
		WHERE id = $1
	`
	var event models.Event
	err := h.db.QueryRow(query, id).Scan(
		&event.ID, &event.EventType, &event.Payload,
		&event.Status, &event.CreatedAt,
	)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Event not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch event")
		return
	}

	// Get delivery attempts
	deliveriesQuery := `
		SELECT id, event_id, endpoint_id, attempt_number, status_code, 
		       latency_ms, error_message, created_at
		FROM delivery_attempts
		WHERE event_id = $1
		ORDER BY attempt_number ASC
	`
	rows, err := h.db.Query(deliveriesQuery, id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch deliveries")
		return
	}
	defer rows.Close()

	deliveries := []*models.DeliveryAttempt{}
	for rows.Next() {
		var da models.DeliveryAttempt
		err := rows.Scan(
			&da.ID, &da.EventID, &da.EndpointID, &da.AttemptNumber,
			&da.StatusCode, &da.LatencyMs, &da.ErrorMessage, &da.CreatedAt,
		)
		if err != nil {
			continue
		}
		deliveries = append(deliveries, &da)
	}

	result := &models.EventWithDeliveries{
		Event:      &event,
		Deliveries: deliveries,
	}

	respondJSON(w, http.StatusOK, result)
}

// ReplayEvent re-queues a failed event for delivery
func (h *EventHandler) ReplayEvent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	// Get original event
	query := `
		SELECT id, event_type, payload, status, created_at
		FROM events
		WHERE id = $1
	`
	var originalEvent models.Event
	err := h.db.QueryRow(query, id).Scan(
		&originalEvent.ID, &originalEvent.EventType,
		&originalEvent.Payload, &originalEvent.Status, &originalEvent.CreatedAt,
	)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Event not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch event")
		return
	}

	// Create new event with same payload
	newEvent := &models.Event{
		ID:        uuid.New(),
		EventType: originalEvent.EventType,
		Payload:   originalEvent.Payload,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	insertQuery := `
		INSERT INTO events (id, event_type, payload, status, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = h.db.Exec(insertQuery,
		newEvent.ID, newEvent.EventType, newEvent.Payload,
		newEvent.Status, newEvent.CreatedAt,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to replay event")
		return
	}

	// Find matching endpoints and queue
	endpoints, err := h.findMatchingEndpoints(newEvent.EventType)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to find endpoints")
		return
	}

	if len(endpoints) > 0 {
		h.db.Exec("UPDATE events SET status = 'delivering' WHERE id = $1", newEvent.ID)
		
		for _, endpoint := range endpoints {
			h.engine.QueueDelivery(newEvent, endpoint)
		}
	}

	respondJSON(w, http.StatusOK, newEvent)
}

// findMatchingEndpoints finds all active endpoints subscribed to an event type
func (h *EventHandler) findMatchingEndpoints(eventType string) ([]*models.Endpoint, error) {
	query := `
		SELECT id, name, url, secret, event_types, is_active, created_at, updated_at
		FROM endpoints
		WHERE is_active = true AND $1 = ANY(event_types)
	`
	rows, err := h.db.Query(query, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	endpoints := []*models.Endpoint{}
	for rows.Next() {
		var ep models.Endpoint
		err := rows.Scan(
			&ep.ID, &ep.Name, &ep.URL, &ep.Secret,
			&ep.EventTypes, &ep.IsActive,
			&ep.CreatedAt, &ep.UpdatedAt,
		)
		if err != nil {
			continue
		}
		endpoints = append(endpoints, &ep)
	}

	return endpoints, nil
}
