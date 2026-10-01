package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/yourusername/webhookrelay/internal/config"
	"github.com/yourusername/webhookrelay/internal/database"
	"github.com/yourusername/webhookrelay/internal/delivery"
	"github.com/yourusername/webhookrelay/internal/handlers"
	custommw "github.com/yourusername/webhookrelay/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Printf("Starting WebhookRelay API on :%s", cfg.Port)

	// Connect to Supabase
	db, err := database.NewSupabaseDB(cfg.SupabaseDBURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize delivery engine
	dbAdapter := handlers.NewDBAdapter(db.DB)
	retryConfig := &delivery.RetryConfig{
		MaxAttempts: cfg.MaxRetryAttempts,
		BaseDelay:   cfg.BaseRetryDelay,
		MaxDelay:    cfg.MaxRetryDelay,
		Multiplier:  2.0,
	}
	
	engine := delivery.NewEngine(
		cfg.WorkerPoolSize,
		1000, // queue size
		dbAdapter,
		retryConfig,
		cfg.RequestTimeout,
	)
	engine.Start()
	defer engine.Shutdown()

	log.Printf("Worker pool started with %d workers", cfg.WorkerPoolSize)

	// Initialize handlers
	webhookHandler := handlers.NewWebhookHandler(db.DB)
	eventHandler := handlers.NewEventHandler(db.DB, engine)
	deliveryHandler := handlers.NewDeliveryHandler(db.DB)
	metricsHandler := handlers.NewMetricsHandler(db.DB)
	mockReceiver := handlers.NewMockReceiver()

	// Setup router
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(custommw.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Route("/api", func(r chi.Router) {
		// Webhooks/Endpoints
		r.Route("/webhooks", func(r chi.Router) {
			r.Post("/", webhookHandler.CreateEndpoint)
			r.Get("/", webhookHandler.ListEndpoints)
			r.Get("/{id}", webhookHandler.GetEndpoint)
			r.Patch("/{id}", webhookHandler.UpdateEndpoint)
			r.Delete("/{id}", webhookHandler.DeleteEndpoint)
		})

		// Events
		r.Route("/events", func(r chi.Router) {
			r.Post("/", eventHandler.CreateEvent)
			r.Get("/", eventHandler.ListEvents)
			r.Get("/{id}", eventHandler.GetEvent)
			r.Post("/{id}/replay", eventHandler.ReplayEvent)
		})

		// Deliveries
		r.Get("/deliveries", deliveryHandler.ListDeliveries)

		// Metrics
		r.Get("/metrics", metricsHandler.GetMetrics)

		// Mock receiver control (for testing)
		r.Post("/mock/control", mockReceiver.SetMode)
		r.Get("/mock/mode", mockReceiver.GetMode)
	})

	// Mock webhook receiver endpoint
	r.Post("/webhook", mockReceiver.HandleWebhook)
	r.Post("/webhook-a", mockReceiver.HandleWebhook)
	r.Post("/webhook-b", mockReceiver.HandleWebhook)
	r.Post("/webhook-c", mockReceiver.HandleWebhook)

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// Start server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		<-sigint

		log.Println("Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server started on http://localhost:%s", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}

	log.Println("Server stopped")
}
