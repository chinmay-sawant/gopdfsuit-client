// Package builder provides fluent builder implementations for constructing PDF documents.
package builder

import (
	"sync"

	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
)

// documentBuilder implements the DocumentBuilder interface with fluent API.
type documentBuilder struct {
	doc *domain.Document
	mu  sync.Mutex
}

// NewDocumentBuilder creates a new DocumentBuilder instance.
func NewDocumentBuilder() domain.DocumentBuilder {
	return &documentBuilder{
		doc: &domain.Document{
			Tables: make([]domain.Table, 0),
			Images: make([]domain.Image, 0),
		},
	}
}

// WithConfig sets the document configuration.
func (b *documentBuilder) WithConfig(config domain.Config) domain.DocumentBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc.Config = config
	return b
}

// WithTitle sets the document title.
func (b *documentBuilder) WithTitle(props, text string) domain.DocumentBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc.Title.Props = props
	b.doc.Title.Text = text
	return b
}

// WithTitleTable sets the document title table.
func (b *documentBuilder) WithTitleTable(table domain.Table) domain.DocumentBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc.Title.Table = &table
	return b
}

// AddTable adds a table to the document.
func (b *documentBuilder) AddTable(table domain.Table) domain.DocumentBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc.Tables = append(b.doc.Tables, table)
	return b
}

// AddImage adds an image to the document.
func (b *documentBuilder) AddImage(image domain.Image) domain.DocumentBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc.Images = append(b.doc.Images, image)
	return b
}

// WithFooter sets the document footer.
func (b *documentBuilder) WithFooter(font, text string) domain.DocumentBuilder {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc.Footer = domain.Footer{
		Font: font,
		Text: text,
	}
	return b
}

// Build constructs and returns the final document.
func (b *documentBuilder) Build() *domain.Document {
	b.mu.Lock()
	defer b.mu.Unlock()
	doc := *b.doc
	return &doc
}

// Reset clears the builder state for reuse.
func (b *documentBuilder) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.doc = &domain.Document{
		Tables: make([]domain.Table, 0),
		Images: make([]domain.Image, 0),
	}
}
