# GoPDFSuit Client

A Go client library for creating and sending PDF documents to a PDF generation service (GoPdfSuit). This library provides a fluent API for building PDF documents programmatically or loading them from JSON in Go.

## Requirements

- Go 1.23.0 or higher
- A running PDF generation service (default endpoint: `http://localhost:8080`)

## Installation

```bash
go get github.com/chinmay/gopdfsuit-client
```

## Features

- Fluent builder pattern for constructing PDF documents
- JSON file/bytes reader for loading document definitions
- HTTP client with retry support and configurable timeouts
- Form field support (text fields, checkboxes, radio buttons)
- Customizable page configuration (size, borders, alignment, watermark)
- Table-based layout with flexible column widths

## Design Patterns Used

- **Builder Pattern** - Document, Table, Cell, Config builders
- **Factory Pattern** - DocumentFactory, FormBuilder
- **Functional Options** - Client configuration
- **Strategy Pattern** - RetryPolicy interface

## Project Structure

```
gopdfsuit-client/
├── gopdfsuit.go           # Main package entry point with public API
├── go.mod                 # Go module definition
├── sample.json            # Sample JSON document definition
├── makefile               # Build and run commands
├── internal/
│   ├── builder/           # Builder implementations
│   │   ├── document_builder.go
│   │   ├── table_builder.go
│   │   ├── cell_builder.go
│   │   └── config_builder.go
│   ├── client/            # HTTP client implementations
│   │   ├── base_client.go
│   │   ├── http_client.go
│   │   ├── pdf_client.go
│   │   ├── header_client.go
│   │   └── retry_client.go
│   ├── domain/            # Domain types and interfaces
│   │   ├── document.go
│   │   ├── config.go
│   │   ├── table.go
│   │   ├── form.go
│   │   ├── common.go
│   │   ├── interfaces.go
│   │   └── errors.go
│   ├── factory/           # Factory implementations
│   │   └── document_factory.go
│   ├── reader/            # Document readers
│   │   └── json_reader.go
│   └── utils/             # Utility functions
│       ├── io.go
│       └── retry.go
└── samplecode/            # Example implementations
    ├── builder/
    │   └── main.go        # Builder pattern example
    └── reader/
        └── main.go        # JSON reader example
```

## Quick Start

### Import the Package

```go
import pdf "github.com/chinmay/gopdfsuit-client"
```

### Create a Client

```go
client := pdf.NewClient(
    "http://localhost:8080",
    pdf.WithTimeout(60*time.Second),
    pdf.WithMaxRetries(3),
    pdf.WithHeader("Authorization", "Bearer your-token"), //as of now no support added for tokens 
)
```

### Option 1: Read from JSON File

```go
ctx := context.Background()

// Read document from JSON file
doc, err := client.ReadFromFile(ctx, "sample.json")
if err != nil {
    log.Fatal(err)
}

// Send to PDF service and save
err = client.SendAndSave(ctx, doc, "output.pdf")
```

### Option 2: Build with Builder Pattern

```go
ctx := context.Background()

// Create configuration
config := pdf.NewConfigBuilder().
    WithPage(pdf.PageSizeA4).
    WithPageBorder(1, 1, 1, 1).
    WithPageAlignment(1).
    Build()

// Build a table
table := pdf.NewTableBuilder().
    WithColumns(2, []float64{1, 3}).
    AddRow(
        pdf.NewCell("font1:9:100:left:1:1:1:1", "Label:"),
        pdf.NewTextFieldCell("font1:9:000:left:1:1:1:1", "Value", "field_name", "Value"),
    ).
    Build()

// Build the document
doc := pdf.NewDocumentBuilder().
    WithConfig(config).
    WithTitle("font1:16:100:center:0:0:0:1", "My Document").
    AddTable(table).
    WithFooter("font1:7:000:center", "Footer text").
    Build()

// Send and save
err := client.SendAndSave(ctx, doc, "output.pdf")
```

### Option 3: Quick Form with FormBuilder

```go
doc := pdf.NewFormBuilder().
    WithTitle("Contact Form").
    AddSection("Personal Information").
    AddTextField("Name:", "name", "John Doe").
    AddTwoColumnTextField("Email:", "email", "john@example.com", "Phone:", "phone", "555-1234").
    WithFooter("Confidential").
    Build()
```

## API Reference

### Client Options

| Option | Description |
|--------|-------------|
| `WithTimeout(duration)` | Sets the HTTP client timeout (default: 30s) |
| `WithMaxRetries(n)` | Sets maximum retry attempts (default: 3) |
| `WithEndpoint(path)` | Sets the PDF generation endpoint |
| `WithHeader(key, value)` | Adds a custom header to all requests |

### Page Sizes

- `pdf.PageSizeA4`
- `pdf.PageSizeLetter`
- `pdf.PageSizeLegal`

### Form Field Types

- `pdf.FormFieldText` - Text input field
- `pdf.FormFieldCheckbox` - Checkbox field
- `pdf.FormFieldRadio` - Radio button field

### Cell Helpers

```go
// Simple text cell
pdf.NewCell(props, text)

// Text field cell
pdf.NewTextFieldCell(props, text, name, value)

// Checkbox cell
pdf.NewCheckboxCell(props, name, value, checked)

// Radio button cell
pdf.NewRadioCell(props, name, value, groupName, checked)
```

### Cell Props Format

Cell properties follow this format: `font:size:weight:alignment:borderTop:borderRight:borderBottom:borderLeft`

Example: `"font1:9:100:left:1:1:1:1"`
- `font1` - Font family
- `9` - Font size
- `100` - Font weight (000=normal, 100=bold, 010=italic)
- `left` - Alignment (left, center, right)
- `1:1:1:1` - Borders (top:right:bottom:left, 1=visible, 0=hidden)

## Running Examples

Use the makefile to run sample code:

```bash
# Run the builder pattern example
make run-builder

# Run the JSON reader example
make run-reader
```

## Error Handling

The library provides typed errors for common scenarios:

```go
var (
    pdf.ErrDocumentNil        // Document is nil
    pdf.ErrInvalidConfig      // Invalid configuration
    pdf.ErrEmptyDocument      // Document has no content
    pdf.ErrFileNotFound       // JSON file not found
    pdf.ErrInvalidJSON        // Invalid JSON format
    pdf.ErrHTTPRequest        // HTTP request failed
    pdf.ErrTimeout            // Request timed out
    pdf.ErrMaxRetriesExceeded // Max retries exceeded
    pdf.ErrInvalidResponse    // Invalid server response
    pdf.ErrUnauthorized       // Authentication failed
    pdf.ErrServerError        // Server error
)
```

## JSON Document Format

Documents can be defined in JSON format. See [sample.json](sample.json) for a complete example.

Basic structure:

```json
{
  "config": {
    "pageBorder": "1:1:1:1",
    "page": "A4",
    "pageAlignment": 1,
    "watermark": ""
  },
  "title": {
    "props": "font1:16:100:left:0:0:0:1",
    "text": "Document Title"
  },
  "table": [
    {
      "maxcolumns": 2,
      "columnwidths": [1, 3],
      "rows": [
        {
          "row": [
            {"props": "font1:9:100:left:1:1:1:1", "text": "Label:"},
            {"props": "font1:9:000:left:1:1:1:1", "text": "Value"}
          ]
        }
      ]
    }
  ],
  "footer": {
    "font": "font1:7:000:center",
    "text": "Footer text"
  }
}
```

## License
[MIT]