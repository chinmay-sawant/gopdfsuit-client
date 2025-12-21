package client

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
)

// BaseClient implements the basic HTTP request execution.
type BaseClient struct {
	client  *http.Client
	headers map[string]string
}

// NewBaseClient creates a new BaseClient.
func NewBaseClient(client *http.Client, headers map[string]string) *BaseClient {
	return &BaseClient{
		client:  client,
		headers: headers,
	}
}

// Do executes the HTTP request using the underlying http.Client.
func (c *BaseClient) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	for k, v := range c.headers {
		req.Header.Set(k, v)
	}
	// Content-Type for JSON is usually set by the caller (Post) or we can set it here if we know.
	// But `Post` in `Client` sets it. Wait, `Client.Post` calls `Do`.
	// `Client.Post` sets `Content-Type` on the request? No, it can't access the request.
	// `Client.Post` in the OLD code set it on the request.

	// We need to handle Content-Type.
	if method == http.MethodPost {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrHTTPRequest, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return responseBody, nil
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, domain.ErrUnauthorized
	}

	return nil, domain.NewHTTPError(resp.StatusCode, fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(responseBody)), nil)
}

// Post is not implemented in BaseClient as it's a convenience method.
func (c *BaseClient) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}

// Get is not implemented in BaseClient as it's a convenience method.
func (c *BaseClient) Get(ctx context.Context, url string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
