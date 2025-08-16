package suites

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/tests/integration"
)

// TestDevIntegration - Development test for GoLand that runs without build tags
func TestDevIntegration(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("BITGET_API_KEY", "mock_test_key")
	os.Setenv("BITGET_SECRET_KEY", "mock_test_secret")
	os.Setenv("BITGET_PASSPHRASE", "mock_test_passphrase")
	os.Setenv("BITGET_MOCK_MODE", "true")
	os.Setenv("BITGET_DEMO_TRADING", "true")

	// Load test configuration
	config, err := integration.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Create futures client
	client := futures.NewClient(config.APIKey, config.SecretKey, config.Passphrase)
	if config.BaseURL != "" {
		client.SetApiEndpoint(config.BaseURL)
	}

	// Test account info service (this will fail with mock credentials, which is expected)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	service := account.NewAccountInfoService(client).
		ProductType(futures.ProductType(config.TestProductType))

	response, err := service.Do(ctx)

	// We expect this to fail with mock credentials
	if err != nil {
		t.Logf("Expected error with mock credentials: %v", err)
		// This is expected - mock credentials should fail
		if response == nil {
			t.Log("✅ Mock credentials properly rejected by API (expected behavior)")
		}
	} else {
		t.Error("❌ Unexpected success with mock credentials - this shouldn't happen")
	}

	t.Log("🎯 Integration test structure is working correctly")
	t.Log("💡 To run with real API keys, set BITGET_MOCK_MODE=false and provide real credentials")
}
