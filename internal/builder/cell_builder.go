// Package builder provides cell builder implementation.
package builder

import (
	"github.com/chinmay/gopdfsuit-client/internal/domain"
)

// cellBuilder implements the CellBuilder interface.
type cellBuilder struct {
	cell domain.Cell
}

// NewCellBuilder creates a new CellBuilder instance.
func NewCellBuilder() domain.CellBuilder {
	return &cellBuilder{
		cell: domain.Cell{},
	}
}

// WithProps sets the cell properties.
func (b *cellBuilder) WithProps(props string) domain.CellBuilder {
	b.cell.Props = props
	return b
}

// WithText sets the cell text.
func (b *cellBuilder) WithText(text string) domain.CellBuilder {
	b.cell.Text = text
	return b
}

// WithTextField adds a text form field.
func (b *cellBuilder) WithTextField(name, value string) domain.CellBuilder {
	b.cell.FormField = &domain.FormField{
		Type:  domain.FormFieldText,
		Name:  name,
		Value: value,
	}
	return b
}

// WithCheckbox adds a checkbox form field.
func (b *cellBuilder) WithCheckbox(name, value string, checked bool) domain.CellBuilder {
	b.cell.FormField = &domain.FormField{
		Type:    domain.FormFieldCheckbox,
		Name:    name,
		Value:   value,
		Checked: checked,
	}
	return b
}

// WithRadio adds a radio button form field.
func (b *cellBuilder) WithRadio(name, value, groupName string, checked bool) domain.CellBuilder {
	b.cell.FormField = &domain.FormField{
		Type:      domain.FormFieldRadio,
		Name:      name,
		Value:     value,
		GroupName: groupName,
		Checked:   checked,
		Shape:     "round",
	}
	return b
}

// Build constructs and returns the final cell.
func (b *cellBuilder) Build() domain.Cell {
	return b.cell
}

// Cell is a helper function to quickly create a cell.
func Cell(props, text string) domain.Cell {
	return domain.Cell{
		Props: props,
		Text:  text,
	}
}

// TextFieldCell creates a cell with a text form field.
func TextFieldCell(props, text, name, value string) domain.Cell {
	return domain.Cell{
		Props: props,
		Text:  text,
		FormField: &domain.FormField{
			Type:  domain.FormFieldText,
			Name:  name,
			Value: value,
		},
	}
}

// CheckboxCell creates a cell with a checkbox form field.
func CheckboxCell(props, name, value string, checked bool) domain.Cell {
	return domain.Cell{
		Props: props,
		Text:  "",
		FormField: &domain.FormField{
			Type:    domain.FormFieldCheckbox,
			Name:    name,
			Value:   value,
			Checked: checked,
		},
	}
}

// RadioCell creates a cell with a radio button form field.
func RadioCell(props, name, value, groupName string, checked bool) domain.Cell {
	return domain.Cell{
		Props: props,
		Text:  "",
		FormField: &domain.FormField{
			Type:      domain.FormFieldRadio,
			Name:      name,
			Value:     value,
			GroupName: groupName,
			Checked:   checked,
			Shape:     "round",
		},
	}
}
