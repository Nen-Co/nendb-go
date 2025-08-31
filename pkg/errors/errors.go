package errors

import (
	"fmt"
	"time"
)

// NenDBError represents the base error type for NenDB operations
type NenDBError struct {
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
	Time    time.Time             `json:"time"`
}

func (e *NenDBError) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("%s - Details: %v", e.Message, e.Details)
	}
	return e.Message
}

// New creates a new NenDBError
func New(message string, details map[string]interface{}) *NenDBError {
	if details == nil {
		details = make(map[string]interface{})
	}
	return &NenDBError{
		Message: message,
		Details: details,
		Time:    time.Now(),
	}
}

// NenDBConnectionError is raised when connection to NenDB server fails
type NenDBConnectionError struct {
	*NenDBError
}

func NewConnectionError(message string, details map[string]interface{}) *NenDBConnectionError {
	return &NenDBConnectionError{
		NenDBError: New(message, details),
	}
}

// NenDBTimeoutError is raised when a request times out
type NenDBTimeoutError struct {
	*NenDBError
}

func NewTimeoutError(message string, details map[string]interface{}) *NenDBTimeoutError {
	return &NenDBTimeoutError{
		NenDBError: New(message, details),
	}
}

// NenDBValidationError is raised when input validation fails
type NenDBValidationError struct {
	*NenDBError
}

func NewValidationError(message string, details map[string]interface{}) *NenDBValidationError {
	return &NenDBValidationError{
		NenDBError: New(message, details),
	}
}

// NenDBAlgorithmError is raised when graph algorithm execution fails
type NenDBAlgorithmError struct {
	*NenDBError
}

func NewAlgorithmError(message string, details map[string]interface{}) *NenDBAlgorithmError {
	return &NenDBAlgorithmError{
		NenDBError: New(message, details),
	}
}

// NenDBResponseError is raised when the server returns an error response
type NenDBResponseError struct {
	*NenDBError
}

func NewResponseError(message string, details map[string]interface{}) *NenDBResponseError {
	return &NenDBResponseError{
		NenDBError: New(message, details),
	}
}
