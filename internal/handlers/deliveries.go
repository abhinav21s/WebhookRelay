package handlers

import (
	"database/sql"
	"net/http"

	"github.com/yourusername/webhookrelay/internal/models"
)

type DeliveryHandler struct {
	db *sql.DB
}

func NewDeliveryHandler(db *sql.DB) *DeliveryHandler {
	return &DeliveryHandler{db: db}
}

// ListDeliveries returns delivery attempts, optionally filtered by event_id
func (h *DeliveryHandler) ListDeliveries(w http.ResponseWriter, r *http.Request) {
	eventID := r.URL.Query().Get("event_id")

	var query string
	var args []interface{}

	if eventID != "" {
		query = `
			SELECT id, event_id, endpoint_id, attempt_number, status_code,
			       latency_ms, error_message, created_at
			FROM delivery_attempts
			WHERE event_id = $1
			ORDER BY created_at DESC
		`
		args = append(args, eventID)
	} else {
		query = `
			SELECT id, event_id, endpoint_id, attempt_number, status_code,
			       latency_ms, error_message, created_at
			FROM delivery_attempts
			ORDER BY created_at DESC
			LIMIT 100
		`
	}

	rows, err := h.db.Query(query, args...)
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

	respondJSON(w, http.StatusOK, deliveries)
}
