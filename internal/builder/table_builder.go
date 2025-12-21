// Package builder provides table builder implementation.
package builder

import (
	"github.com/chinmay-sawant/gopdfsuit-client/internal/domain"
)

// tableBuilder implements the TableBuilder interface.
type tableBuilder struct {
	table domain.Table
}

// NewTableBuilder creates a new TableBuilder instance.
func NewTableBuilder() domain.TableBuilder {
	return &tableBuilder{
		table: domain.Table{
			Rows: make([]domain.Row, 0),
		},
	}
}

// WithColumns sets the column configuration.
func (b *tableBuilder) WithColumns(maxColumns int, widths []float64) domain.TableBuilder {
	b.table.MaxColumns = maxColumns
	b.table.ColumnWidths = widths
	return b
}

// AddRow adds a row to the table.
func (b *tableBuilder) AddRow(cells ...domain.Cell) domain.TableBuilder {
	row := domain.Row{
		Cells: cells,
	}
	b.table.Rows = append(b.table.Rows, row)
	return b
}

// AddRowWithHeight adds a row with custom height.
func (b *tableBuilder) AddRowWithHeight(height int, cells ...domain.Cell) domain.TableBuilder {
	row := domain.Row{
		Height: height,
		Cells:  cells,
	}
	b.table.Rows = append(b.table.Rows, row)
	return b
}

// Build constructs and returns the final table.
func (b *tableBuilder) Build() domain.Table {
	return b.table
}

// Reset clears the builder state for reuse.
func (b *tableBuilder) Reset() {
	b.table = domain.Table{
		Rows: make([]domain.Row, 0),
	}
}
