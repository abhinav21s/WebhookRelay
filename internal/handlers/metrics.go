package handlers

import (
	"database/sql"
	"net/http"

	"github.com/yourusername/webhookrelay/internal/models"
)

type MetricsHandler struct {
	db *sql.DB
}

func NewMetricsHandler(db *sql.DB) *MetricsHandler {
	return &MetricsHandler{db: db}
}

// GetMetrics returns aggregated metrics
func (h *MetricsHandler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := &models.Metrics{}

	// Total events in last 24 hours
	query1 := `
		SELECT COUNT(*) FROM events 
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`
	h.db.QueryRow(query1).Scan(&metrics.TotalEvents24h)

	// Success rate
	query2 := `
		SELECT 
			COUNT(CASE WHEN status = 'delivered' THEN 1 END)::FLOAT / NULLIF(COUNT(*), 0) * 100
		FROM events
		WHERE created_at > NOW() - INTERVAL '24 hours'
	`
	h.db.QueryRow(query2).Scan(&metrics.SuccessRate)

	// Average latency
	query3 := `
		SELECT COALESCE(AVG(latency_ms), 0)
		FROM delivery_attempts
		WHERE created_at > NOW() - INTERVAL '24 hours' AND status_code BETWEEN 200 AND 299
	`
	h.db.QueryRow(query3).Scan(&metrics.AverageLatencyMs)

	// Active endpoints
	query4 := `SELECT COUNT(*) FROM endpoints WHERE is_active = true`
	h.db.QueryRow(query4).Scan(&metrics.ActiveEndpoints)

	respondJSON(w, http.StatusOK, metrics)
}
