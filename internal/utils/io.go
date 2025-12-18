package utils

import (
	"bytes"
)

// NewBytesReader creates a new bytes.Reader.
func NewBytesReader(b []byte) *bytes.Reader {
	return bytes.NewReader(b)
}
