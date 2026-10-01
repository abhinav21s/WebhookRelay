package delivery

import (
	"math"
	"time"
)

// RetryConfig holds retry configuration
type RetryConfig struct {
	MaxAttempts     int
	BaseDelay       time.Duration
	MaxDelay        time.Duration
	Multiplier      float64
}

// DefaultRetryConfig returns default retry configuration
func DefaultRetryConfig() *RetryConfig {
	return &RetryConfig{
		MaxAttempts: 5,
		BaseDelay:   1 * time.Second,
		MaxDelay:    30 * time.Second,
		Multiplier:  2.0,
	}
}

// CalculateBackoff calculates exponential backoff delay
// Formula: min(baseDelay * multiplier^(attempt-1), maxDelay)
func (rc *RetryConfig) CalculateBackoff(attemptNumber int) time.Duration {
	if attemptNumber <= 1 {
		return rc.BaseDelay
	}

	// Calculate exponential delay: baseDelay * multiplier^(attempt-1)
	exponent := float64(attemptNumber - 1)
	delay := float64(rc.BaseDelay) * math.Pow(rc.Multiplier, exponent)

	// Cap at max delay
	if time.Duration(delay) > rc.MaxDelay {
		return rc.MaxDelay
	}

	return time.Duration(delay)
}

// ShouldRetry determines if an HTTP status code should trigger a retry
func ShouldRetry(statusCode int) bool {
	// Retry on:
	// - 5xx server errors
	// - 429 rate limiting
	// - 408 request timeout
	// - 0 (network error, no response)
	
	if statusCode == 0 {
		return true // Network error
	}
	
	if statusCode >= 500 && statusCode < 600 {
		return true // Server errors
	}
	
	if statusCode == 429 || statusCode == 408 {
		return true // Rate limit or timeout
	}
	
	return false
}

// IsRetryableError checks if an error is retryable
func IsRetryableError(err error) bool {
	if err == nil {
		return false
	}
	
	// Network errors, timeouts, connection refused, etc. are retryable
	return true
}
