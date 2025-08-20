package common

import (
	"testing"
)

func TestSafeStringCast(t *testing.T) {
	// Test valid string
	result := SafeStringCast("hello")
	if result != "hello" {
		t.Errorf("Expected 'hello', got '%s'", result)
	}

	// Test nil value
	result = SafeStringCast(nil)
	if result != "" {
		t.Errorf("Expected empty string for nil, got '%s'", result)
	}

	// Test non-string value
	result = SafeStringCast(123)
	if result != "" {
		t.Errorf("Expected empty string for int, got '%s'", result)
	}

	// Test interface{} with string
	var value interface{} = "test"
	result = SafeStringCast(value)
	if result != "test" {
		t.Errorf("Expected 'test', got '%s'", result)
	}
}

func TestSafeFloat64Cast(t *testing.T) {
	// Test valid float64
	result := SafeFloat64Cast(123.45)
	if result != 123.45 {
		t.Errorf("Expected 123.45, got %f", result)
	}

	// Test nil value
	result = SafeFloat64Cast(nil)
	if result != 0.0 {
		t.Errorf("Expected 0.0 for nil, got %f", result)
	}

	// Test non-float64 value
	result = SafeFloat64Cast("not a float")
	if result != 0.0 {
		t.Errorf("Expected 0.0 for string, got %f", result)
	}

	// Test interface{} with float64
	var value interface{} = 99.99
	result = SafeFloat64Cast(value)
	if result != 99.99 {
		t.Errorf("Expected 99.99, got %f", result)
	}
}
