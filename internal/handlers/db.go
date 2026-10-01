package handlers

import (
	"database/sql"

	"github.com/yourusername/webhookrelay/internal/models"
)

// DBAdapter adapts sql.DB to the DeliveryDB interface
type DBAdapter struct {
	db *sql.DB
}

func NewDBAdapter(db *sql.DB) *DBAdapter {
	return &DBAdapter{db: db}
}

func (dba *DBAdapter) LogDeliveryAttempt(attempt *models.DeliveryAttempt) error {
	query := `
		INSERT INTO delivery_attempts (id, event_id, endpoint_id, attempt_number, status_code, latency_ms, error_message, created_at)
		VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
	`
	_, err := dba.db.Exec(query,
		attempt.EventID, attempt.EndpointID, attempt.AttemptNumber,
		attempt.StatusCode, attempt.LatencyMs, attempt.ErrorMessage, attempt.CreatedAt,
	)
	return err
}

func (dba *DBAdapter) UpdateEventStatus(eventID string, status string) error {
	query := "UPDATE events SET status = $1 WHERE id = $2"
	_, err := dba.db.Exec(query, status, eventID)
	return err
}
