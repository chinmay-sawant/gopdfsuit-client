package domain

// Document represents the complete PDF document structure.
type Document struct {
	Config Config  `json:"config"`
	Title  Title   `json:"title"`
	Tables []Table `json:"table"`
	Images []Image `json:"image"`
	Footer Footer  `json:"footer"`
}

// Title represents the document title.
type Title struct {
	Props string `json:"props"`
	Text  string `json:"text"`
	Table *Table `json:"table,omitempty"`
}
