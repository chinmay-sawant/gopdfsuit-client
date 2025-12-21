// Package builder provides configuration builder helpers.
package builder

import (
	"fmt"

	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
)

// ConfigBuilder provides a fluent API for building Config.
type ConfigBuilder struct {
	config domain.Config
}

// NewConfigBuilder creates a new ConfigBuilder with defaults.
func NewConfigBuilder() *ConfigBuilder {
	return &ConfigBuilder{
		config: domain.Config{
			Page:          string(domain.PageSizeA4),
			PageBorder:    "1:1:1:1",
			PageAlignment: 1,
		},
	}
}

// WithPage sets the page size.
func (b *ConfigBuilder) WithPage(page domain.PageSize) *ConfigBuilder {
	b.config.Page = string(page)
	return b
}

// WithPageBorder sets the page border (top:right:bottom:left).
func (b *ConfigBuilder) WithPageBorder(top, right, bottom, left int) *ConfigBuilder {
	b.config.PageBorder = fmt.Sprintf("%d:%d:%d:%d", top, right, bottom, left)
	return b
}

// WithPageAlignment sets the page alignment.
func (b *ConfigBuilder) WithPageAlignment(alignment int) *ConfigBuilder {
	b.config.PageAlignment = alignment
	return b
}

// WithWatermark sets the watermark text.
func (b *ConfigBuilder) WithWatermark(watermark string) *ConfigBuilder {
	b.config.Watermark = watermark
	return b
}

// Build returns the final Config.
func (b *ConfigBuilder) Build() domain.Config {
	return b.config
}

// PropsBuilder provides a fluent API for building cell property strings.
type PropsBuilder struct {
	font      string
	size      int
	weight    string
	alignment string
	borders   [4]int
}

// NewPropsBuilder creates a new PropsBuilder with defaults.
func NewPropsBuilder() *PropsBuilder {
	return &PropsBuilder{
		font:      "font1",
		size:      9,
		weight:    "000",
		alignment: "left",
		borders:   [4]int{1, 1, 1, 1},
	}
}

// WithFont sets the font name.
func (b *PropsBuilder) WithFont(font string) *PropsBuilder {
	b.font = font
	return b
}

// WithSize sets the font size.
func (b *PropsBuilder) WithSize(size int) *PropsBuilder {
	b.size = size
	return b
}

// Bold sets the font weight to bold.
func (b *PropsBuilder) Bold() *PropsBuilder {
	b.weight = "100"
	return b
}

// Italic sets the font style to italic.
func (b *PropsBuilder) Italic() *PropsBuilder {
	b.weight = "010"
	return b
}

// BoldItalic sets the font to bold and italic.
func (b *PropsBuilder) BoldItalic() *PropsBuilder {
	b.weight = "110"
	return b
}

// Normal sets the font weight to normal.
func (b *PropsBuilder) Normal() *PropsBuilder {
	b.weight = "000"
	return b
}

// WithAlignment sets the text alignment.
func (b *PropsBuilder) WithAlignment(alignment domain.Alignment) *PropsBuilder {
	b.alignment = string(alignment)
	return b
}

// Left sets left alignment.
func (b *PropsBuilder) Left() *PropsBuilder {
	b.alignment = "left"
	return b
}

// Center sets center alignment.
func (b *PropsBuilder) Center() *PropsBuilder {
	b.alignment = "center"
	return b
}

// Right sets right alignment.
func (b *PropsBuilder) Right() *PropsBuilder {
	b.alignment = "right"
	return b
}

// WithBorders sets all borders.
func (b *PropsBuilder) WithBorders(top, right, bottom, left int) *PropsBuilder {
	b.borders = [4]int{top, right, bottom, left}
	return b
}

// NoBorders removes all borders.
func (b *PropsBuilder) NoBorders() *PropsBuilder {
	b.borders = [4]int{0, 0, 0, 0}
	return b
}

// AllBorders sets all borders to 1.
func (b *PropsBuilder) AllBorders() *PropsBuilder {
	b.borders = [4]int{1, 1, 1, 1}
	return b
}

// Build returns the final props string.
func (b *PropsBuilder) Build() string {
	return fmt.Sprintf("%s:%d:%s:%s:%d:%d:%d:%d",
		b.font, b.size, b.weight, b.alignment,
		b.borders[0], b.borders[1], b.borders[2], b.borders[3])
}
