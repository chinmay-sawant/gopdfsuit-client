package domain

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
