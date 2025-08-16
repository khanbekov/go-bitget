package ws

import (
	"github.com/khanbekov/go-bitget/common"
	"github.com/rs/zerolog"
	"testing"
	"time"
)

// TestReconnectFunctionality tests the manual reconnection feature
func TestReconnectFunctionality(t *testing.T) {
	logger := zerolog.Nop() // Use Nop logger to suppress output during testing

	// Create a client with a mock WebSocket endpoint (this will fail to connect)
	client := NewBitgetBaseWsClient(logger, "wss://invalid-url.example.com", "")

	// Test configuration methods
	client.SetMaxReconnectAttempts(1) // Reduce attempts for faster testing
	client.SetReconnectionTimeout(1 * time.Second)

	// Test that reconnection attempts are made
	err := client.Reconnect()
	if err == nil {
		t.Error("Expected reconnection to fail with invalid URL")
	}

	// Verify error message contains expected content
	expectedErrorMsg := "maximum reconnection attempts (1) exceeded"
	if err.Error() != expectedErrorMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedErrorMsg, err.Error())
	}
}

// TestReconnectState tests the reconnection state management
func TestReconnectState(t *testing.T) {
	logger := zerolog.Nop() // Use Nop logger to suppress output during testing
	client := NewBitgetBaseWsClient(logger, "wss://invalid-url.example.com", "")

	// Initially should not be reconnecting
	if client.reconnecting {
		t.Error("Client should not be reconnecting initially")
	}

	// Test reconnect flag is set during operation
	client.SetMaxReconnectAttempts(0) // Set to 0 to test flag setting without network calls

	// This tests the method exists and basic state management
	// We can't easily test concurrent behavior without network timeouts
}

// TestLoginCredentialStorage tests that login credentials are stored for re-authentication
func TestLoginCredentialStorage(t *testing.T) {
	logger := zerolog.Nop() // Use Nop logger to suppress output during testing
	client := NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/private", "test-secret")

	// Initially no credentials should be stored
	if client.storedLoginCreds != nil {
		t.Error("No credentials should be stored initially")
	}

	// After login, credentials should be stored
	client.Login("test-api-key", "test-passphrase", common.SHA256)

	if client.storedLoginCreds == nil {
		t.Error("Credentials should be stored after login")
	}

	if client.storedLoginCreds.apiKey != "test-api-key" {
		t.Errorf("Expected API key 'test-api-key', got '%s'", client.storedLoginCreds.apiKey)
	}

	if client.storedLoginCreds.passphrase != "test-passphrase" {
		t.Errorf("Expected passphrase 'test-passphrase', got '%s'", client.storedLoginCreds.passphrase)
	}

	if !client.needLogin {
		t.Error("needLogin should be true after calling Login")
	}
}

// TestConnectionStateManagement tests connection state during reconnection
func TestConnectionStateManagement(t *testing.T) {
	logger := zerolog.Nop() // Use Nop logger to suppress output during testing
	client := NewBitgetBaseWsClient(logger, "wss://invalid-url.example.com", "")

	// Set max attempts to 1 for quick test
	client.SetMaxReconnectAttempts(1)

	// Simulate connected state
	client.connected = true
	client.loginStatus = true

	// After failed reconnection, states should be properly reset
	client.Reconnect()

	if client.connected {
		t.Error("Client should not be connected after failed reconnection")
	}

	if client.loginStatus {
		t.Error("Login status should be false after failed reconnection")
	}
}
