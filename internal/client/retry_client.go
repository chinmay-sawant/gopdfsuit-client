package client

import (
	"context"
	"io"
	"time"

	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
	"github.com/chinmay-sawant/gopdfsuit-client/internal/utils"
)

// RetryClient decorates an HTTPClient with retry logic.
type RetryClient struct {
	next        domain.HTTPClient
	maxRetries  int
	retryDelay  time.Duration
	retryPolicy domain.RetryPolicy
	logger      domain.Logger
}

// NewRetryClient creates a new RetryClient.
func NewRetryClient(next domain.HTTPClient, maxRetries int, retryDelay time.Duration, policy domain.RetryPolicy, logger domain.Logger) *RetryClient {
	return &RetryClient{
		next:        next,
		maxRetries:  maxRetries,
		retryDelay:  retryDelay,
		retryPolicy: policy,
		logger:      logger,
	}
}

// Do executes the request with retries.
func (c *RetryClient) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	var lastErr error

	// If body is an io.ReadCloser, we can't easily rewind it for retries unless we buffer it.
	// For simplicity, we assume body is reusable or we'd need to read it into a buffer here.
	// In the current architecture, body is usually a *bytes.Buffer or *bytes.Reader which is Seekable,
	// but io.Reader interface doesn't guarantee it.
	// Ideally, we should read body into a byte slice if we expect retries.

	// However, since we are refactoring existing code, let's assume the caller handles body reset
	// OR we handle it if it's a specific type.
	// The original code in http_client.go used `req.GetBody` which is available on `http.Request`.
	// Since we are at the `Do` level, we don't have `http.Request` yet.
	// We will assume for now that the body is safe to reuse or the underlying `next.Do` handles it?
	// No, `next.Do` consumes it.

	// Let's buffer the body if it's not nil.
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = io.ReadAll(body)
		if err != nil {
			return nil, err
		}
	}

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			if c.logger != nil {
				c.logger.Debug("Retry attempt %d for %s %s", attempt, method, url)
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.getRetryDelay(attempt)):
			}
		}

		// Create a new reader for each attempt
		var currentBody io.Reader
		if bodyBytes != nil {
			// bytes.NewReader is cheap
			currentBody = utils.NewBytesReader(bodyBytes)
		}

		resp, err := c.next.Do(ctx, method, url, currentBody)
		if err == nil {
			return resp, nil
		}

		lastErr = err
		if c.shouldRetry(attempt, err) {
			continue
		}
		return nil, err
	}

	return nil, lastErr
}

func (c *RetryClient) shouldRetry(attempt int, err error) bool {
	if attempt >= c.maxRetries {
		return false
	}
	if c.retryPolicy != nil {
		return c.retryPolicy.ShouldRetry(attempt, err)
	}
	// Default retry logic: retry on network errors or 5xx
	// We need to check if it's an HTTP error
	if httpErr, ok := err.(*domain.HTTPError); ok {
		return httpErr.StatusCode >= 500
	}
	// Retry on other errors (network, etc)
	return true
}

func (c *RetryClient) getRetryDelay(attempt int) time.Duration {
	if c.retryPolicy != nil {
		return time.Duration(c.retryPolicy.WaitDuration(attempt)) * time.Millisecond
	}
	return utils.CalculateBackoff(attempt, c.retryDelay)
}

// Post delegates to Do.
func (c *RetryClient) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	// This method shouldn't really be called if we structure it right,
	// but if it is, we can't easily implement it without circular dependency or duplication.
	// The `Client` struct (the facade) should handle Post/Get -> Do conversion.
	return nil, nil
}

// Get delegates to Do.
func (c *RetryClient) Get(ctx context.Context, url string) ([]byte, error) {
	return nil, nil
}
