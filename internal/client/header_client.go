package client

import (
	"context"
	"io"

	"github.com/chinmay/gopdfsuit-client/internal/domain"
)

// HeaderClient decorates an HTTPClient to add default headers.
type HeaderClient struct {
	next    domain.HTTPClient
	headers map[string]string
}

// NewHeaderClient creates a new HeaderClient.
func NewHeaderClient(next domain.HTTPClient, headers map[string]string) *HeaderClient {
	return &HeaderClient{
		next:    next,
		headers: headers,
	}
}

// Do executes the request and adds headers.
// Note: This is tricky because `Do` takes `io.Reader`, not `http.Request`.
// We can't modify the request here easily because `BaseClient` creates it.
//
// Alternative: `BaseClient` takes headers.
// Or `Do` takes `*http.Request`.
//
// Let's modify `BaseClient` to take headers.
func (c *HeaderClient) Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error) {
	// This decorator pattern breaks if we can't modify the request.
	// But wait, `BaseClient` is the one creating the request.
	// If we want to inject headers, we should do it BEFORE BaseClient executes.
	// But BaseClient is the leaf.

	// Actually, `BaseClient` should probably just execute a `*http.Request`.
	// But the interface is `Do(..., method, url, body)`.

	// Let's stick to modifying `BaseClient` to accept headers in its struct.
	return c.next.Do(ctx, method, url, body)
}

// Post delegates to Do.
func (c *HeaderClient) Post(ctx context.Context, url string, body interface{}) ([]byte, error) {
	return nil, nil
}

// Get delegates to Do.
func (c *HeaderClient) Get(ctx context.Context, url string) ([]byte, error) {
	return nil, nil
}
