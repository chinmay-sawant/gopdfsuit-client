// Package domain contains custom errors for the PDF client.
package domain

import "errors"

var (
	// ErrDocumentNil is returned when a nil document is provided.
	ErrDocumentNil = errors.New("document cannot be nil")

	// ErrInvalidConfig is returned when the document configuration is invalid.
	ErrInvalidConfig = errors.New("invalid document configuration")

	// ErrEmptyDocument is returned when attempting to send an empty document.
	ErrEmptyDocument = errors.New("document is empty")

	// ErrFileNotFound is returned when the specified file is not found.
	ErrFileNotFound = errors.New("file not found")

	// ErrInvalidJSON is returned when JSON parsing fails.
	ErrInvalidJSON = errors.New("invalid JSON format")

	// ErrHTTPRequest is returned when an HTTP request fails.
	ErrHTTPRequest = errors.New("HTTP request failed")

	// ErrTimeout is returned when a request times out.
	ErrTimeout = errors.New("request timed out")

	// ErrMaxRetriesExceeded is returned when maximum retry attempts are exceeded.
	ErrMaxRetriesExceeded = errors.New("maximum retry attempts exceeded")

	// ErrInvalidResponse is returned when the server response is invalid.
	ErrInvalidResponse = errors.New("invalid server response")

	// ErrUnauthorized is returned when authentication fails.
	ErrUnauthorized = errors.New("unauthorized: invalid or missing credentials")

	// ErrServerError is returned when the server returns an error.
	ErrServerError = errors.New("server error")
)

// HTTPError represents an HTTP error with status code.
type HTTPError struct {
	StatusCode int
	Message    string
	Err        error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *HTTPError) Unwrap() error {
	return e.Err
}

// NewHTTPError creates a new HTTPError.
func NewHTTPError(statusCode int, message string, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Message:    message,
		Err:        err,
	}
}
