// Package factory provides factory pattern implementations for creating documents.
package factory

import (
	"github.com/chinmay/gopdfsuit-client/internal/builder"
	"github.com/chinmay/gopdfsuit-client/internal/domain"
)

// DocumentType represents the type of document to create.
type DocumentType string

const (
	DocumentTypeForm    DocumentType = "form"
	DocumentTypeReport  DocumentType = "report"
	DocumentTypeInvoice DocumentType = "invoice"
	DocumentTypeCustom  DocumentType = "custom"
)

// RadioOption represents a radio button option.
type RadioOption struct {
	Label   string
	Value   string
	Checked bool
}

// CheckboxOption represents a checkbox option.
type CheckboxOption struct {
	Name    string
	Label   string
	Value   string
	Checked bool
}

// DocumentFactory creates documents based on templates or configurations.
type DocumentFactory struct {
	defaultConfig domain.Config
}

// NewDocumentFactory creates a new DocumentFactory with default configuration.
func NewDocumentFactory() *DocumentFactory {
	return &DocumentFactory{
		defaultConfig: domain.Config{
			Page:          string(domain.PageSizeA4),
			PageBorder:    "1:1:1:1",
			PageAlignment: 1,
		},
	}
}

// WithDefaultConfig sets the default configuration for all created documents.
func (f *DocumentFactory) WithDefaultConfig(config domain.Config) *DocumentFactory {
	f.defaultConfig = config
	return f
}

// CreateDocument creates a new document based on the specified type.
func (f *DocumentFactory) CreateDocument(docType DocumentType) domain.DocumentBuilder {
	docBuilder := builder.NewDocumentBuilder()
	docBuilder.WithConfig(f.defaultConfig)

	switch docType {
	case DocumentTypeForm:
		return docBuilder.WithFooter("font1:7:000:center", "")
	case DocumentTypeReport:
		return docBuilder.WithFooter("font1:8:000:center", "")
	case DocumentTypeInvoice:
		return docBuilder.WithFooter("font1:7:000:right", "")
	default:
		return docBuilder
	}
}

// CreateFormBuilder creates a new form builder for patient registration forms.
func (f *DocumentFactory) CreateFormBuilder() *FormBuilder {
	return NewFormBuilder(f.defaultConfig)
}

// FormBuilder provides a specialized builder for form documents.
type FormBuilder struct {
	docBuilder domain.DocumentBuilder
	config     domain.Config
}

// NewFormBuilder creates a new FormBuilder.
func NewFormBuilder(config domain.Config) *FormBuilder {
	return &FormBuilder{
		docBuilder: builder.NewDocumentBuilder(),
		config:     config,
	}
}

// WithTitle sets the form title.
func (fb *FormBuilder) WithTitle(text string) *FormBuilder {
	props := builder.NewPropsBuilder().WithSize(16).Bold().Left().WithBorders(0, 0, 0, 1).Build()
	fb.docBuilder.WithTitle(props, text)
	return fb
}

// AddSection adds a section header to the form.
func (fb *FormBuilder) AddSection(title string) *FormBuilder {
	props := builder.NewPropsBuilder().WithSize(10).Bold().Left().AllBorders().Build()
	table := builder.NewTableBuilder().WithColumns(1, []float64{1}).AddRow(builder.Cell(props, title)).Build()
	fb.docBuilder.AddTable(table)
	return fb
}

// AddTextField adds a text field row to the form.
func (fb *FormBuilder) AddTextField(label, name, value string) *FormBuilder {
	labelProps := builder.NewPropsBuilder().WithSize(9).Bold().Left().AllBorders().Build()
	valueProps := builder.NewPropsBuilder().WithSize(9).Normal().Left().AllBorders().Build()
	table := builder.NewTableBuilder().
		WithColumns(2, []float64{1, 3}).
		AddRow(builder.Cell(labelProps, label), builder.TextFieldCell(valueProps, value, name, value)).
		Build()
	fb.docBuilder.AddTable(table)
	return fb
}

// AddTwoColumnTextField adds two text fields in a row.
func (fb *FormBuilder) AddTwoColumnTextField(label1, name1, value1, label2, name2, value2 string) *FormBuilder {
	labelProps := builder.NewPropsBuilder().WithSize(9).Bold().Left().AllBorders().Build()
	valueProps := builder.NewPropsBuilder().WithSize(9).Normal().Left().AllBorders().Build()
	table := builder.NewTableBuilder().
		WithColumns(4, []float64{1, 2, 1, 2}).
		AddRow(
			builder.Cell(labelProps, label1),
			builder.TextFieldCell(valueProps, value1, name1, value1),
			builder.Cell(labelProps, label2),
			builder.TextFieldCell(valueProps, value2, name2, value2),
		).
		Build()
	fb.docBuilder.AddTable(table)
	return fb
}

// WithFooter sets the form footer.
func (fb *FormBuilder) WithFooter(text string) *FormBuilder {
	fb.docBuilder.WithFooter("font1:7:000:center", text)
	return fb
}

// Build constructs the final document.
func (fb *FormBuilder) Build() *domain.Document {
	fb.docBuilder.WithConfig(fb.config)
	return fb.docBuilder.Build()
}
