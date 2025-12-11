// Package client provides the HTTP client implementation for the PDF service.
package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/chinmay/gopdfsuit-client/internal/domain"
)

// Config holds the client configuration.
type Config struct {
	BaseURL     string
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	Headers     map[string]string
	Logger      domain.Logger
	RetryPolicy domain.RetryPolicy
}

// DefaultConfig returns a default configuration.
func DefaultConfig() *Config {
	return &Config{
		Timeout:    30 * time.Second,
		MaxRetries: 3,
		RetryDelay: time.Second,
		Headers:    make(map[string]string),
	}
}

// Option is a functional option for configuring the client.
type Option func(*Client)

// WithTimeout sets the client timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(c *Client) {
		c.config.Timeout = timeout
	}
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(maxRetries int) Option {
	return func(c *Client) {
		c.config.MaxRetries = maxRetries
	}
}

// WithRetryDelay sets the delay between retries.
func WithRetryDelay(delay time.Duration) Option {
	return func(c *Client) {
		c.config.RetryDelay = delay
	}
}

// WithHeader adds a header to all requests.
func WithHeader(key, value string) Option {
	return func(c *Client) {
		c.config.Headers[key] = value
	}
}

// WithLogger sets the logger.
func WithLogger(logger domain.Logger) Option {
	return func(c *Client) {
		c.config.Logger = logger
	}
}

// WithRetryPolicy sets a custom retry policy.
func WithRetryPolicy(policy domain.RetryPolicy) Option {
	return func(c *Client) {
		c.config.RetryPolicy = policy
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// Client is the HTTP client for the PDF service.
type Client struct {
	config     *Config
	httpClient *http.Client
}

// New creates a new Client with the given options.
func New(baseURL string, opts ...Option) *Client {
	config := DefaultConfig()
	config.BaseURL = baseURL

	c := &Client{
		config: config,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{
			Timeout: c.config.Timeout,
		}
	}

	return c
}

// Do executes an HTTP request with retry logic.
func (c *Client) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			c.logDebug("Retry attempt %d for %s %s", attempt, method, url)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.getRetryDelay(attempt)):
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		for key, value := range c.config.Headers {
			req.Header.Set(key, value)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if c.shouldRetry(attempt, err) {
				continue
			}
			return nil, fmt.Errorf("%w: %v", domain.ErrHTTPRequest, err)
		}

		defer resp.Body.Close()

		responseBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return responseBody, nil
		}

		httpErr := domain.NewHTTPError(resp.StatusCode, fmt.Sprintf("HTTP %d", resp.StatusCode), nil)
		lastErr = httpErr

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, domain.ErrUnauthorized
		}

		if resp.StatusCode >= 500 && c.shouldRetry(attempt, httpErr) {
			continue
		}

		return nil, httpErr
	}

	return nil, fmt.Errorf("%w: %v", domain.ErrMaxRetriesExceeded, lastErr)
}

// Post sends a POST request with JSON body.
func (c *Client) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fullURL := c.config.BaseURL + url
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	return c.doWithRetry(ctx, req)
}

// Get sends a GET request.
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	fullURL := c.config.BaseURL + url
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range c.config.Headers {
		req.Header.Set(key, value)
	}

	return c.doWithRetry(ctx, req)
}

// doWithRetry executes an HTTP request with retry logic.
func (c *Client) doWithRetry(ctx context.Context, req *http.Request) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.config.MaxRetries; attempt++ {
		if attempt > 0 {
			c.logDebug("Retry attempt %d for %s %s", attempt, req.Method, req.URL.String())
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.getRetryDelay(attempt)):
			}

			if req.GetBody != nil {
				body, err := req.GetBody()
				if err == nil {
					req.Body = body
				}
			}
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			if c.shouldRetry(attempt, err) {
				continue
			}
			return nil, fmt.Errorf("%w: %v", domain.ErrHTTPRequest, err)
		}

		responseBody, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return responseBody, nil
		}

		httpErr := domain.NewHTTPError(resp.StatusCode, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(responseBody)), nil)
		lastErr = httpErr

		if resp.StatusCode == http.StatusUnauthorized {
			return nil, domain.ErrUnauthorized
		}

		if resp.StatusCode >= 500 && c.shouldRetry(attempt, httpErr) {
			continue
		}

		return nil, httpErr
	}

	return nil, fmt.Errorf("%w: %v", domain.ErrMaxRetriesExceeded, lastErr)
}

// shouldRetry determines if a retry should be attempted.
func (c *Client) shouldRetry(attempt int, err error) bool {
	if attempt >= c.config.MaxRetries {
		return false
	}

	if c.config.RetryPolicy != nil {
		return c.config.RetryPolicy.ShouldRetry(attempt, err)
	}

	return true
}

// getRetryDelay returns the delay before the next retry.
func (c *Client) getRetryDelay(attempt int) time.Duration {
	if c.config.RetryPolicy != nil {
		return time.Duration(c.config.RetryPolicy.WaitDuration(attempt)) * time.Millisecond
	}

	return c.config.RetryDelay * time.Duration(1<<uint(attempt))
}

// logDebug logs a debug message if a logger is configured.
func (c *Client) logDebug(msg string, args ...interface{}) {
	if c.config.Logger != nil {
		c.config.Logger.Debug(msg, args...)
	}
}
