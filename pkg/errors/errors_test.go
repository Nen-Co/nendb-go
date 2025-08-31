package errors

import (
	"testing"
)

func TestNenDBError(t *testing.T) {
	// Test basic error creation
	err := New("Test error message", map[string]interface{}{"key": "value"})
	if err.Message != "Test error message" {
		t.Errorf("Expected message 'Test error message', got '%s'", err.Message)
	}
	if err.Details["key"] != "value" {
		t.Errorf("Expected details key 'value', got '%v'", err.Details["key"])
	}
	if err.Time.IsZero() {
		t.Error("Expected time to be set, got zero time")
	}

	// Test error string representation
	errStr := err.Error()
	expectedStr := "Test error message - Details: map[key:value]"
	if errStr != expectedStr {
		t.Errorf("Expected error string '%s', got '%s'", expectedStr, errStr)
	}

	// Test error without details
	errNoDetails := New("Simple error", nil)
	if errNoDetails.Details == nil {
		t.Error("Expected details to be initialized as empty map, got nil")
	}
	if len(errNoDetails.Details) != 0 {
		t.Errorf("Expected empty details map, got map with %d items", len(errNoDetails.Details))
	}

	errStrNoDetails := errNoDetails.Error()
	if errStrNoDetails != "Simple error" {
		t.Errorf("Expected error string 'Simple error', got '%s'", errStrNoDetails)
	}
}

func TestNenDBConnectionError(t *testing.T) {
	err := NewConnectionError("Connection failed", map[string]interface{}{"url": "http://localhost:8080"})
	
	// Test that it's a NenDBError
	if err.Message != "Connection failed" {
		t.Errorf("Expected message 'Connection failed', got '%s'", err.Message)
	}
	if err.Details["url"] != "http://localhost:8080" {
		t.Errorf("Expected details url 'http://localhost:8080', got '%v'", err.Details["url"])
	}

	// Test error string
	errStr := err.Error()
	expectedStr := "Connection failed - Details: map[url:http://localhost:8080]"
	if errStr != expectedStr {
		t.Errorf("Expected error string '%s', got '%s'", expectedStr, errStr)
	}
}

func TestNenDBTimeoutError(t *testing.T) {
	err := NewTimeoutError("Request timed out", map[string]interface{}{"timeout": "30s"})
	
	if err.Message != "Request timed out" {
		t.Errorf("Expected message 'Request timed out', got '%s'", err.Message)
	}
	if err.Details["timeout"] != "30s" {
		t.Errorf("Expected details timeout '30s', got '%v'", err.Details["timeout"])
	}
}

func TestNenDBValidationError(t *testing.T) {
	err := NewValidationError("Invalid input", map[string]interface{}{"field": "age", "value": -5})
	
	if err.Message != "Invalid input" {
		t.Errorf("Expected message 'Invalid input', got '%s'", err.Message)
	}
	if err.Details["field"] != "age" {
		t.Errorf("Expected details field 'age', got '%v'", err.Details["field"])
	}
	if err.Details["value"] != -5 {
		t.Errorf("Expected details value -5, got '%v'", err.Details["value"])
	}
}

func TestNenDBAlgorithmError(t *testing.T) {
	err := NewAlgorithmError("Algorithm failed", map[string]interface{}{"algorithm": "BFS", "reason": "No path found"})
	
	if err.Message != "Algorithm failed" {
		t.Errorf("Expected message 'Algorithm failed', got '%s'", err.Message)
	}
	if err.Details["algorithm"] != "BFS" {
		t.Errorf("Expected details algorithm 'BFS', got '%v'", err.Details["algorithm"])
	}
	if err.Details["reason"] != "No path found" {
		t.Errorf("Expected details reason 'No path found', got '%v'", err.Details["reason"])
	}
}

func TestNenDBResponseError(t *testing.T) {
	err := NewResponseError("Server error", map[string]interface{}{"status": 500, "code": "INTERNAL_ERROR"})
	
	if err.Message != "Server error" {
		t.Errorf("Expected message 'Server error', got '%s'", err.Message)
	}
	if err.Details["status"] != 500 {
		t.Errorf("Expected details status 500, got '%v'", err.Details["status"])
	}
	if err.Details["code"] != "INTERNAL_ERROR" {
		t.Errorf("Expected details code 'INTERNAL_ERROR', got '%v'", err.Details["code"])
	}
}

func TestErrorInheritance(t *testing.T) {
	// Test that all error types properly inherit from NenDBError
	connectionErr := NewConnectionError("test", nil)
	timeoutErr := NewTimeoutError("test", nil)
	validationErr := NewValidationError("test", nil)
	algorithmErr := NewAlgorithmError("test", nil)
	responseErr := NewResponseError("test", nil)

	// All should have the same base structure
	if connectionErr.Message != "test" {
		t.Error("ConnectionError should inherit Message from NenDBError")
	}
	if timeoutErr.Message != "test" {
		t.Error("TimeoutError should inherit Message from NenDBError")
	}
	if validationErr.Message != "test" {
		t.Error("ValidationError should inherit Message from NenDBError")
	}
	if algorithmErr.Message != "test" {
		t.Error("AlgorithmError should inherit Message from NenDBError")
	}
	if responseErr.Message != "test" {
		t.Error("ResponseError should inherit Message from NenDBError")
	}
}
