package domain

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
