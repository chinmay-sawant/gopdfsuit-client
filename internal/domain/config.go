package domain

// Config holds the page configuration settings.
type Config struct {
	PageBorder    string `json:"pageBorder"`
	Page          string `json:"page"`
	PageAlignment int    `json:"pageAlignment"`
	Watermark     string `json:"watermark"`
}
