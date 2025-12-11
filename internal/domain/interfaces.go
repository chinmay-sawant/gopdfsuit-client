// Package domain contains the core business interfaces for the PDF client.
package domain

import (
	"context"
	"io"
)

// DocumentReader defines the interface for reading document data from various sources.
type DocumentReader interface {
	// Read reads document data from the underlying source.
	Read(ctx context.Context) (*Document, error)
}

// DocumentSender defines the interface for sending documents to a remote service.
type DocumentSender interface {
	// Send sends the document and returns the response body.
	Send(ctx context.Context, doc *Document) ([]byte, error)
}

// DocumentBuilder defines the interface for building documents using fluent API.
type DocumentBuilder interface {
	// WithConfig sets the document configuration.
	WithConfig(config Config) DocumentBuilder
	// WithTitle sets the document title.
	WithTitle(props, text string) DocumentBuilder
	// WithTitleTable sets the document title table.
	WithTitleTable(table Table) DocumentBuilder
	// AddTable adds a table to the document.
	AddTable(table Table) DocumentBuilder
	// AddImage adds an image to the document.
	AddImage(image Image) DocumentBuilder
	// WithFooter sets the document footer.
	WithFooter(font, text string) DocumentBuilder
	// Build constructs and returns the final document.
	Build() *Document
	// Reset clears the builder state for reuse.
	Reset()
}

// TableBuilder defines the interface for building tables.
type TableBuilder interface {
	// WithColumns sets the column configuration.
	WithColumns(maxColumns int, widths []float64) TableBuilder
	// AddRow adds a row to the table.
	AddRow(cells ...Cell) TableBuilder
	// AddRowWithHeight adds a row with custom height.
	AddRowWithHeight(height int, cells ...Cell) TableBuilder
	// Build constructs and returns the final table.
	Build() Table
	// Reset clears the builder state for reuse.
	Reset()
}

// CellBuilder defines the interface for building cells.
type CellBuilder interface {
	// WithProps sets the cell properties.
	WithProps(props string) CellBuilder
	// WithText sets the cell text.
	WithText(text string) CellBuilder
	// WithTextField adds a text form field.
	WithTextField(name, value string) CellBuilder
	// WithCheckbox adds a checkbox form field.
	WithCheckbox(name, value string, checked bool) CellBuilder
	// WithRadio adds a radio button form field.
	WithRadio(name, value, groupName string, checked bool) CellBuilder
	// Build constructs and returns the final cell.
	Build() Cell
}

// HTTPClient defines the interface for HTTP operations.
type HTTPClient interface {
	// Do executes an HTTP request.
	Do(ctx context.Context, method, url string, body io.Reader) ([]byte, error)
	// Post sends a POST request with JSON body.
	Post(ctx context.Context, url string, body interface{}) ([]byte, error)
	// Get sends a GET request.
	Get(ctx context.Context, url string) ([]byte, error)
}

// ClientOption defines a functional option for configuring the client.
type ClientOption func(interface{})

// RetryPolicy defines the retry behavior for failed requests.
type RetryPolicy interface {
	// ShouldRetry determines if a retry should be attempted.
	ShouldRetry(attempt int, err error) bool
	// WaitDuration returns the duration to wait before the next retry.
	WaitDuration(attempt int) int64
}

// Logger defines the interface for logging.
type Logger interface {
	// Debug logs a debug message.
	Debug(msg string, args ...interface{})
	// Info logs an info message.
	Info(msg string, args ...interface{})
	// Warn logs a warning message.
	Warn(msg string, args ...interface{})
	// Error logs an error message.
	Error(msg string, args ...interface{})
}
