package domain

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
