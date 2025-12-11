// Package domain contains the core business models for the PDF document structure.
package domain

// Document represents the complete PDF document structure.
type Document struct {
	Config Config  `json:"config"`
	Title  Title   `json:"title"`
	Tables []Table `json:"table"`
	Images []Image `json:"image"`
	Footer Footer  `json:"footer"`
}

// Config holds the page configuration settings.
type Config struct {
	PageBorder    string `json:"pageBorder"`
	Page          string `json:"page"`
	PageAlignment int    `json:"pageAlignment"`
	Watermark     string `json:"watermark"`
}

// Title represents the document title.
type Title struct {
	Props string `json:"props"`
	Text  string `json:"text"`
	Table *Table `json:"table,omitempty"`
}

// Table represents a table structure in the document.
type Table struct {
	MaxColumns   int       `json:"maxcolumns"`
	ColumnWidths []float64 `json:"columnwidths"`
	Rows         []Row     `json:"rows"`
}

// Row represents a row in a table.
type Row struct {
	Height int    `json:"height,omitempty"`
	Cells  []Cell `json:"row"`
}

// Cell represents a cell in a table row.
type Cell struct {
	Props     string     `json:"props"`
	Text      string     `json:"text"`
	FormField *FormField `json:"form_field,omitempty"`
}

// FormField represents an interactive form field.
type FormField struct {
	Type      FormFieldType `json:"type"`
	Name      string        `json:"name"`
	Value     string        `json:"value"`
	Checked   bool          `json:"checked,omitempty"`
	GroupName string        `json:"group_name,omitempty"`
	Shape     string        `json:"shape,omitempty"`
}

// FormFieldType represents the type of form field.
type FormFieldType string

const (
	FormFieldText     FormFieldType = "text"
	FormFieldCheckbox FormFieldType = "checkbox"
	FormFieldRadio    FormFieldType = "radio"
)

// Image represents an image in the document.
type Image struct {
	Path   string  `json:"path,omitempty"`
	X      float64 `json:"x,omitempty"`
	Y      float64 `json:"y,omitempty"`
	Width  float64 `json:"width,omitempty"`
	Height float64 `json:"height,omitempty"`
}

// Footer represents the document footer.
type Footer struct {
	Font string `json:"font"`
	Text string `json:"text"`
}

// PageSize represents common page sizes.
type PageSize string

const (
	PageSizeA4     PageSize = "A4"
	PageSizeLetter PageSize = "Letter"
	PageSizeLegal  PageSize = "Legal"
)

// Alignment represents text alignment options.
type Alignment string

const (
	AlignLeft   Alignment = "left"
	AlignCenter Alignment = "center"
	AlignRight  Alignment = "right"
)
