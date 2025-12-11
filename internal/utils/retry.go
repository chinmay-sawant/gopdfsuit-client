package utils

import (
	"time"
)

// CalculateBackoff returns the delay before the next retry using exponential backoff.
func CalculateBackoff(attempt int, baseDelay time.Duration) time.Duration {
	return baseDelay * time.Duration(1<<uint(attempt))
}
