package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/yourusername/webhookrelay/internal/models"
)

type WebhookHandler struct {
	db *sql.DB
}

func NewWebhookHandler(db *sql.DB) *WebhookHandler {
	return &WebhookHandler{db: db}
}

// CreateEndpoint creates a new webhook endpoint
func (h *WebhookHandler) CreateEndpoint(w http.ResponseWriter, r *http.Request) {
	var req models.CreateEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate
	if req.Name == "" || req.URL == "" || len(req.EventTypes) == 0 {
		respondError(w, http.StatusBadRequest, "Name, URL, and event types are required")
		return
	}

	// Generate secret if not provided
	if req.Secret == "" {
		req.Secret = generateSecret(32)
	}

	// Insert into database
	endpoint := &models.Endpoint{
		ID:         uuid.New(),
		Name:       req.Name,
		URL:        req.URL,
		Secret:     req.Secret,
		EventTypes: req.EventTypes,
		IsActive:   true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	query := `
		INSERT INTO endpoints (id, name, url, secret, event_types, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := h.db.Exec(query,
		endpoint.ID, endpoint.Name, endpoint.URL, endpoint.Secret,
		pq.Array(endpoint.EventTypes), endpoint.IsActive,
		endpoint.CreatedAt, endpoint.UpdatedAt,
	)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to create endpoint")
		return
	}

	respondJSON(w, http.StatusCreated, endpoint)
}

// ListEndpoints returns all endpoints
func (h *WebhookHandler) ListEndpoints(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, name, url, secret, event_types, is_active, created_at, updated_at
		FROM endpoints
		ORDER BY created_at DESC
	`
	rows, err := h.db.Query(query)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch endpoints")
		return
	}
	defer rows.Close()

	endpoints := []*models.Endpoint{}
	for rows.Next() {
		var ep models.Endpoint
		err := rows.Scan(
			&ep.ID, &ep.Name, &ep.URL, &ep.Secret,
			pq.Array(&ep.EventTypes), &ep.IsActive,
			&ep.CreatedAt, &ep.UpdatedAt,
		)
		if err != nil {
			continue
		}
		endpoints = append(endpoints, &ep)
	}

	respondJSON(w, http.StatusOK, endpoints)
}

// GetEndpoint returns a single endpoint by ID
func (h *WebhookHandler) GetEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	
	query := `
		SELECT id, name, url, secret, event_types, is_active, created_at, updated_at
		FROM endpoints
		WHERE id = $1
	`
	var ep models.Endpoint
	err := h.db.QueryRow(query, id).Scan(
		&ep.ID, &ep.Name, &ep.URL, &ep.Secret,
		pq.Array(&ep.EventTypes), &ep.IsActive,
		&ep.CreatedAt, &ep.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		respondError(w, http.StatusNotFound, "Endpoint not found")
		return
	}
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to fetch endpoint")
		return
	}

	respondJSON(w, http.StatusOK, ep)
}

// UpdateEndpoint updates an existing endpoint
func (h *WebhookHandler) UpdateEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req models.UpdateEndpointRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Build dynamic update query
	query := "UPDATE endpoints SET updated_at = $1"
	args := []interface{}{time.Now()}
	argCount := 2

	if req.Name != nil {
		query += ", name = $" + string(rune(argCount+'0'))
		args = append(args, *req.Name)
		argCount++
	}
	if req.URL != nil {
		query += ", url = $" + string(rune(argCount+'0'))
		args = append(args, *req.URL)
		argCount++
	}
	if req.EventTypes != nil {
		query += ", event_types = $" + string(rune(argCount+'0'))
		args = append(args, pq.Array(req.EventTypes))
		argCount++
	}
	if req.IsActive != nil {
		query += ", is_active = $" + string(rune(argCount+'0'))
		args = append(args, *req.IsActive)
		argCount++
	}

	query += " WHERE id = $" + string(rune(argCount+'0'))
	args = append(args, id)

	result, err := h.db.Exec(query, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to update endpoint")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Endpoint not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Endpoint updated successfully"})
}

// DeleteEndpoint soft deletes an endpoint
func (h *WebhookHandler) DeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	query := "UPDATE endpoints SET is_active = false, updated_at = $1 WHERE id = $2"
	result, err := h.db.Exec(query, time.Now(), id)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to delete endpoint")
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		respondError(w, http.StatusNotFound, "Endpoint not found")
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{"message": "Endpoint deleted successfully"})
}

// Helper function to generate random secret
func generateSecret(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Helper functions
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
