// Package client provides the PDF client for document operations.
package client

import (
	"context"
	"os"

	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
)

// PDFClient handles PDF document operations.
type PDFClient struct {
	httpClient *Client
	endpoint   string
}

// NewPDFClient creates a new PDFClient.
func NewPDFClient(httpClient *Client, endpoint string) *PDFClient {
	return &PDFClient{
		httpClient: httpClient,
		endpoint:   endpoint,
	}
}

// Send sends a document to the PDF service and returns the response.
func (c *PDFClient) Send(ctx context.Context, doc *domain.Document) ([]byte, error) {
	if doc == nil {
		return nil, domain.ErrDocumentNil
	}

	return c.httpClient.Post(ctx, c.endpoint, doc)
}

// SendAndSave sends a document and saves the PDF response to the specified path.
func (c *PDFClient) SendAndSave(ctx context.Context, doc *domain.Document, outputPath string) error {
	response, err := c.Send(ctx, doc)
	if err != nil {
		return err
	}

	return saveToFile(outputPath, response)
}

// saveToFile saves data to a file.
func saveToFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
