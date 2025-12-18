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
	doer       domain.HTTPClient
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

	// Build the decorator chain
	var doer domain.HTTPClient = NewBaseClient(c.httpClient, c.config.Headers)

	// Add retry decorator
	if c.config.MaxRetries > 0 {
		doer = NewRetryClient(doer, c.config.MaxRetries, c.config.RetryDelay, c.config.RetryPolicy, c.config.Logger)
	}

	c.doer = doer
	return c
}

// Do executes an HTTP request with retry logic.
func (c *Client) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	return c.doer.Do(ctx, method, url, body)
}

// Post sends a POST request with JSON body.
func (c *Client) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	fullURL := c.config.BaseURL + url
	return c.Do(ctx, http.MethodPost, fullURL, bytes.NewReader(jsonBody))
}

// Get sends a GET request.
func (c *Client) Get(ctx context.Context, url string) ([]byte, error) {
	fullURL := c.config.BaseURL + url
	return c.Do(ctx, http.MethodGet, fullURL, nil)
}

// doWithRetry is removed as it is replaced by RetryClient.
