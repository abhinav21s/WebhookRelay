package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/yourusername/webhookrelay/internal/config"
	"github.com/yourusername/webhookrelay/internal/handlers"
)

func main() {
	cfg := config.Load()
	log.Printf("Starting Mock Receiver on :%s", cfg.MockReceiverPort)

	mockReceiver := handlers.NewMockReceiver()

	r := chi.NewRouter()

	// Control endpoint
	r.Post("/control", mockReceiver.SetMode)
	r.Get("/mode", mockReceiver.GetMode)

	// Webhook endpoints
	r.Post("/webhook", mockReceiver.HandleWebhook)
	r.Post("/webhook-a", mockReceiver.HandleWebhook)
	r.Post("/webhook-b", mockReceiver.HandleWebhook)
	r.Post("/webhook-c", mockReceiver.HandleWebhook)

	log.Printf("Mock receiver started on http://localhost:%s", cfg.MockReceiverPort)
	log.Println("Default mode: success")
	log.Fatal(http.ListenAndServe(":"+cfg.MockReceiverPort, r))
}
