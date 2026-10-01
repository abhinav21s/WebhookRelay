package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	Environment      string
	SupabaseURL      string
	SupabaseKey      string
	SupabaseDBURL    string
	WorkerPoolSize   int
	MaxRetryAttempts int
	BaseRetryDelay   time.Duration
	MaxRetryDelay    time.Duration
	RequestTimeout   time.Duration
	MockReceiverPort string
}

func Load() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	return &Config{
		Port:             getEnv("PORT", "8080"),
		Environment:      getEnv("ENV", "development"),
		SupabaseURL:      getEnv("SUPABASE_URL", ""),
		SupabaseKey:      getEnv("SUPABASE_KEY", ""),
		SupabaseDBURL:    getEnv("SUPABASE_DB_URL", ""),
		WorkerPoolSize:   getEnvAsInt("WORKER_POOL_SIZE", 10),
		MaxRetryAttempts: getEnvAsInt("MAX_RETRY_ATTEMPTS", 5),
		BaseRetryDelay:   getEnvAsDuration("BASE_RETRY_DELAY", time.Second),
		MaxRetryDelay:    getEnvAsDuration("MAX_RETRY_DELAY", 30*time.Second),
		RequestTimeout:   getEnvAsDuration("REQUEST_TIMEOUT", 5*time.Second),
		MockReceiverPort: getEnv("MOCK_RECEIVER_PORT", "9090"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer for %s, using default %d", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid duration for %s, using default %v", key, defaultValue)
		return defaultValue
	}
	return value
}
