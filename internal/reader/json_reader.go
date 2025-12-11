// Package reader provides document reader implementations.
package reader

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/chinmay/gopdfsuit-client/internal/domain"
)

// JSONFileReader reads documents from JSON files.
type JSONFileReader struct {
	filePath string
}

// NewJSONFileReader creates a new JSONFileReader.
func NewJSONFileReader(filePath string) *JSONFileReader {
	return &JSONFileReader{
		filePath: filePath,
	}
}

// Read reads and parses a JSON file into a Document.
func (r *JSONFileReader) Read(ctx context.Context) (*domain.Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w: %s", domain.ErrFileNotFound, r.filePath)
		}
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var doc domain.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidJSON, err)
	}

	return &doc, nil
}

// SetFilePath updates the file path.
func (r *JSONFileReader) SetFilePath(filePath string) {
	r.filePath = filePath
}

// JSONBytesReader reads documents from JSON bytes.
type JSONBytesReader struct {
	data []byte
}

// NewJSONBytesReader creates a new JSONBytesReader.
func NewJSONBytesReader(data []byte) *JSONBytesReader {
	return &JSONBytesReader{
		data: data,
	}
}

// Read parses JSON bytes into a Document.
func (r *JSONBytesReader) Read(ctx context.Context) (*domain.Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	if len(r.data) == 0 {
		return nil, domain.ErrEmptyDocument
	}

	var doc domain.Document
	if err := json.Unmarshal(r.data, &doc); err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInvalidJSON, err)
	}

	return &doc, nil
}

// SetData updates the JSON data.
func (r *JSONBytesReader) SetData(data []byte) {
	r.data = data
}
