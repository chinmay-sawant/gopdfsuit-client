// Package gopdfsuit provides a client library for creating and sending PDF documents.
package gopdfsuit

import (
	"context"
	"time"

	"github.com/chinmay-sawant/gopdfsuit-client/internal/builder"
	"github.com/chinmay-sawant/gopdfsuit-client/internal/client"
	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
	"github.com/chinmay-sawant/gopdfsuit-client/internal/factory"
	"github.com/chinmay-sawant/gopdfsuit-client/internal/reader"
)

// Re-export domain types
type (
	Document      = domain.Document
	Config        = domain.Config
	Title         = domain.Title
	Table         = domain.Table
	Row           = domain.Row
	Cell          = domain.Cell
	FormField     = domain.FormField
	Image         = domain.Image
	Footer        = domain.Footer
	FormFieldType = domain.FormFieldType
	PageSize      = domain.PageSize
	Alignment     = domain.Alignment
)

// Re-export builder interfaces
type (
	DocumentBuilder = domain.DocumentBuilder
	TableBuilder    = domain.TableBuilder
	CellBuilder     = domain.CellBuilder
	DocumentReader  = domain.DocumentReader
	Logger          = domain.Logger
	RetryPolicy     = domain.RetryPolicy
)

// Re-export factory types
type (
	RadioOption    = factory.RadioOption
	CheckboxOption = factory.CheckboxOption
	DocumentType   = factory.DocumentType
)

// Form field type constants
const (
	FormFieldText     = domain.FormFieldText
	FormFieldCheckbox = domain.FormFieldCheckbox
	FormFieldRadio    = domain.FormFieldRadio
)

// Page size constants
const (
	PageSizeA4     = domain.PageSizeA4
	PageSizeLetter = domain.PageSizeLetter
	PageSizeLegal  = domain.PageSizeLegal
)

// Alignment constants
const (
	AlignLeft   = domain.AlignLeft
	AlignCenter = domain.AlignCenter
	AlignRight  = domain.AlignRight
)

// Document type constants
const (
	DocumentTypeForm    = factory.DocumentTypeForm
	DocumentTypeReport  = factory.DocumentTypeReport
	DocumentTypeInvoice = factory.DocumentTypeInvoice
	DocumentTypeCustom  = factory.DocumentTypeCustom
)

// Error variables
var (
	ErrDocumentNil        = domain.ErrDocumentNil
	ErrInvalidConfig      = domain.ErrInvalidConfig
	ErrEmptyDocument      = domain.ErrEmptyDocument
	ErrFileNotFound       = domain.ErrFileNotFound
	ErrInvalidJSON        = domain.ErrInvalidJSON
	ErrHTTPRequest        = domain.ErrHTTPRequest
	ErrTimeout            = domain.ErrTimeout
	ErrMaxRetriesExceeded = domain.ErrMaxRetriesExceeded
	ErrInvalidResponse    = domain.ErrInvalidResponse
	ErrUnauthorized       = domain.ErrUnauthorized
	ErrServerError        = domain.ErrServerError
)

// Client is the main entry point for the PDF client library.
type Client struct {
	httpClient *client.Client
	pdfClient  *client.PDFClient
}

type clientConfig struct {
	timeout    time.Duration
	endpoint   string
	maxRetries int
	headers    map[string]string
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*clientConfig)

// WithTimeout sets the client timeout.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *clientConfig) { c.timeout = timeout }
}

// WithEndpoint sets the PDF generation endpoint.
func WithEndpoint(endpoint string) ClientOption {
	return func(c *clientConfig) { c.endpoint = endpoint }
}

// WithMaxRetries sets the maximum number of retries.
func WithMaxRetries(maxRetries int) ClientOption {
	return func(c *clientConfig) { c.maxRetries = maxRetries }
}

// WithHeader adds a header to all requests.
func WithHeader(key, value string) ClientOption {
	return func(c *clientConfig) {
		if c.headers == nil {
			c.headers = make(map[string]string)
		}
		c.headers[key] = value
	}
}

// NewClient creates a new PDF Client with the given base URL and options.
func NewClient(baseURL string, opts ...ClientOption) *Client {
	cfg := &clientConfig{
		timeout:    30 * time.Second,
		endpoint:   "/api/v1/generate/template-pdf",
		maxRetries: 3,
		headers:    make(map[string]string),
	}
	for _, opt := range opts {
		opt(cfg)
	}

	clientOpts := []client.Option{
		client.WithTimeout(cfg.timeout),
		client.WithMaxRetries(cfg.maxRetries),
	}
	for k, v := range cfg.headers {
		clientOpts = append(clientOpts, client.WithHeader(k, v))
	}

	httpClient := client.New(baseURL, clientOpts...)
	return &Client{
		httpClient: httpClient,
		pdfClient:  client.NewPDFClient(httpClient, cfg.endpoint),
	}
}

// Send sends a document to the PDF service.
func (c *Client) Send(ctx context.Context, doc *Document) ([]byte, error) {
	return c.pdfClient.Send(ctx, doc)
}

// SendAndSave sends a document and saves the PDF to a file.
func (c *Client) SendAndSave(ctx context.Context, doc *Document, outputPath string) error {
	return c.pdfClient.SendAndSave(ctx, doc, outputPath)
}

// ReadFromFile reads a document from a JSON file.
func (c *Client) ReadFromFile(ctx context.Context, filePath string) (*Document, error) {
	return reader.NewJSONFileReader(filePath).Read(ctx)
}

// ReadFromBytes reads a document from JSON bytes.
func (c *Client) ReadFromBytes(ctx context.Context, data []byte) (*Document, error) {
	return reader.NewJSONBytesReader(data).Read(ctx)
}

// NewDocumentBuilder creates a new DocumentBuilder.
func NewDocumentBuilder() DocumentBuilder {
	return builder.NewDocumentBuilder()
}

// NewTableBuilder creates a new TableBuilder.
func NewTableBuilder() TableBuilder {
	return builder.NewTableBuilder()
}

// NewCellBuilder creates a new CellBuilder.
func NewCellBuilder() CellBuilder {
	return builder.NewCellBuilder()
}

// NewConfigBuilder creates a new ConfigBuilder.
func NewConfigBuilder() *builder.ConfigBuilder {
	return builder.NewConfigBuilder()
}

// NewPropsBuilder creates a new PropsBuilder.
func NewPropsBuilder() *builder.PropsBuilder {
	return builder.NewPropsBuilder()
}

// DefaultConfig returns a default document configuration.
func DefaultConfig() Config {
	return builder.NewConfigBuilder().Build()
}

// NewCell creates a simple cell with props and text.
func NewCell(props, text string) Cell {
	return builder.Cell(props, text)
}

// NewTextFieldCell creates a cell with a text form field.
func NewTextFieldCell(props, text, name, value string) Cell {
	return builder.TextFieldCell(props, text, name, value)
}

// NewCheckboxCell creates a cell with a checkbox form field.
func NewCheckboxCell(props, name, value string, checked bool) Cell {
	return builder.CheckboxCell(props, name, value, checked)
}

// NewRadioCell creates a cell with a radio button form field.
func NewRadioCell(props, name, value, groupName string, checked bool) Cell {
	return builder.RadioCell(props, name, value, groupName, checked)
}

// NewJSONFileReader creates a new JSON file reader.
func NewJSONFileReader(filePath string) *reader.JSONFileReader {
	return reader.NewJSONFileReader(filePath)
}

// NewJSONBytesReader creates a new JSON bytes reader.
func NewJSONBytesReader(data []byte) *reader.JSONBytesReader {
	return reader.NewJSONBytesReader(data)
}

// NewDocumentFactory creates a new document factory.
func NewDocumentFactory() *factory.DocumentFactory {
	return factory.NewDocumentFactory()
}

// NewFormBuilder creates a new form builder with default config.
func NewFormBuilder() *factory.FormBuilder {
	return factory.NewFormBuilder(DefaultConfig())
}

// NewFormBuilderWithConfig creates a new form builder with custom config.
func NewFormBuilderWithConfig(config Config) *factory.FormBuilder {
	return factory.NewFormBuilder(config)
}
