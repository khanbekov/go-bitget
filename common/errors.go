package common

import (
	"encoding/json"
	"fmt"
)

// APIError represents API errors with string codes (for UTA API)
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: code=%s, msg=%s", e.Code, e.Message)
}

// MissingParameterError represents a missing required parameter error
type MissingParameterError struct {
	Parameter string
}

func (e *MissingParameterError) Error() string {
	return fmt.Sprintf("missing required parameter: %s", e.Parameter)
}

// NewMissingParameterError creates a new missing parameter error
func NewMissingParameterError(parameter string) *MissingParameterError {
	return &MissingParameterError{Parameter: parameter}
}

// UnmarshalJSON unmarshals JSON data into the provided interface
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
