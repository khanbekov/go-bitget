// Package integration provides integration tests for real API endpoints
//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// TestConfig defines configuration for integration testing
type TestConfig struct {
	// API Credentials
	APIKey     string `json:"api_key"`
	SecretKey  string `json:"secret_key"`
	Passphrase string `json:"passphrase"`

	// Test Environment
	UseTestnet   bool   `json:"use_testnet"`
	BaseURL      string `json:"base_url"`
	WebSocketURL string `json:"websocket_url"`
	DemoTrading  bool   `json:"demo_trading"`
	TestTimeout  int    `json:"test_timeout_seconds"`
	MaxRetries   int    `json:"max_retries"`
	RetryDelayMs int    `json:"retry_delay_ms"`

	// Test Selection
	EnabledSuites  []string          `json:"enabled_suites"`
	EndpointTests  map[string]bool   `json:"endpoint_tests"`
	TestParameters map[string]string `json:"test_parameters"`

	// Test Symbols and Parameters
	TestSymbol      string `json:"test_symbol"`
	TestProductType string `json:"test_product_type"`
	TestCoin        string `json:"test_coin"`

	// Safe Testing Parameters
	SafeOrderSize   string `json:"safe_order_size"`
	SafePriceOffset string `json:"safe_price_offset"` // Percentage offset from market price
	MaxTestAmount   string `json:"max_test_amount"`

	// Reporting
	GenerateReport bool   `json:"generate_report"`
	ReportPath     string `json:"report_path"`
	LogLevel       string `json:"log_level"`
}

// TestSuite represents a group of related tests
type TestSuite struct {
	Name        string
	Description string
	Tests       []TestEndpoint
}

// TestEndpoint represents a single endpoint test configuration
type TestEndpoint struct {
	Name        string            `json:"name"`
	Endpoint    string            `json:"endpoint"`
	Method      string            `json:"method"`
	Description string            `json:"description"`
	Enabled     bool              `json:"enabled"`
	ReadOnly    bool              `json:"read_only"` // Safe to run without side effects
	Parameters  map[string]string `json:"parameters"`
	Expected    map[string]string `json:"expected"` // Expected response fields
	Timeout     int               `json:"timeout_seconds"`
}

// TestResult represents the result of a single test
type TestResult struct {
	TestName  string        `json:"test_name"`
	Endpoint  string        `json:"endpoint"`
	Success   bool          `json:"success"`
	Duration  time.Duration `json:"duration"`
	Error     string        `json:"error,omitempty"`
	Response  string        `json:"response,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// TestReport represents the overall test execution report
type TestReport struct {
	Config          TestConfig        `json:"config"`
	StartTime       time.Time         `json:"start_time"`
	EndTime         time.Time         `json:"end_time"`
	Duration        time.Duration     `json:"duration"`
	TotalTests      int               `json:"total_tests"`
	PassedTests     int               `json:"passed_tests"`
	FailedTests     int               `json:"failed_tests"`
	SkippedTests    int               `json:"skipped_tests"`
	TestResults     []TestResult      `json:"test_results"`
	Summary         string            `json:"summary"`
	EnvironmentInfo map[string]string `json:"environment_info"`
}

// DefaultConfig returns a default configuration for integration testing
func DefaultConfig() *TestConfig {
	return &TestConfig{
		UseTestnet:   true,
		BaseURL:      "https://api.bitget.com",
		WebSocketURL: "wss://ws.bitget.com/v2/ws/public",
		DemoTrading:  true,
		TestTimeout:  30,
		MaxRetries:   3,
		RetryDelayMs: 1000,

		EnabledSuites: []string{"account", "market", "position"},
		EndpointTests: map[string]bool{
			"account_info":      true,
			"account_list":      true,
			"account_bills":     true,
			"set_leverage":      false, // Disabled by default - has side effects
			"adjust_margin":     false, // Disabled by default - has side effects
			"set_margin_mode":   false, // Disabled by default - has side effects
			"set_position_mode": false, // Disabled by default - has side effects
		},

		TestSymbol:      "BTCUSDT",
		TestProductType: "USDT-FUTURES",
		TestCoin:        "USDT",
		SafeOrderSize:   "0.001",
		SafePriceOffset: "0.3", // 30% below market price
		MaxTestAmount:   "10",

		GenerateReport: true,
		ReportPath:     "tests/reports/integration_report.json",
		LogLevel:       "info",

		TestParameters: map[string]string{
			"limit":       "10",
			"granularity": "1m",
			"coin":        "USDT",
		},
	}
}

// LoadConfig loads configuration from environment and config file with fallback mechanisms
func LoadConfig() (*TestConfig, error) {
	config := DefaultConfig()

	// Try to load .env files in order of preference
	loadEnvFiles()

	// Load from environment variables
	if apiKey := os.Getenv("BITGET_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}
	if secretKey := os.Getenv("BITGET_SECRET_KEY"); secretKey != "" {
		config.SecretKey = secretKey
	}
	if passphrase := os.Getenv("BITGET_PASSPHRASE"); passphrase != "" {
		config.Passphrase = passphrase
	}

	// Override with environment-specific settings
	if os.Getenv("BITGET_TESTNET") == "true" {
		config.UseTestnet = true
		config.DemoTrading = true
	}

	if os.Getenv("BITGET_DEMO_TRADING") == "false" {
		config.DemoTrading = false
	}

	// Load from config file if exists
	if configFile := os.Getenv("INTEGRATION_CONFIG_FILE"); configFile != "" {
		if err := config.LoadFromFile(configFile); err != nil {
			return nil, fmt.Errorf("failed to load config file: %w", err)
		}
	} else if _, err := os.Stat("tests/configs/integration.json"); err == nil {
		if err := config.LoadFromFile("tests/configs/integration.json"); err != nil {
			return nil, fmt.Errorf("failed to load default config file: %w", err)
		}
	}

	// Check for mock testing mode
	if os.Getenv("BITGET_MOCK_MODE") == "true" {
		config = setupMockConfig(config)
		return config, nil
	}

	// Validate required fields with helpful error messages
	if config.APIKey == "" || config.SecretKey == "" || config.Passphrase == "" {
		return nil, fmt.Errorf(`API credentials are required for integration tests.

Options to provide credentials:
1. Set environment variables: BITGET_API_KEY, BITGET_SECRET_KEY, BITGET_PASSPHRASE
2. Create .env file (copy from .env.example)
3. Use mock mode: export BITGET_MOCK_MODE=true (for compilation testing only)

Current values:
- BITGET_API_KEY: %q
- BITGET_SECRET_KEY: %q 
- BITGET_PASSPHRASE: %q`, config.APIKey, config.SecretKey, config.Passphrase)
	}

	return config, nil
}

// loadEnvFiles tries to load environment files in order of preference
func loadEnvFiles() {
	envFiles := []string{
		".env.local",       // Local overrides (git-ignored)
		".env",             // Main env file
		".env.development", // Development defaults
	}

	for _, envFile := range envFiles {
		if _, err := os.Stat(envFile); err == nil {
			_ = godotenv.Load(envFile)
			break // Use first found file
		}
	}
}

// setupMockConfig configures mock credentials for compilation testing
func setupMockConfig(config *TestConfig) *TestConfig {
	config.APIKey = "mock_api_key_for_testing"
	config.SecretKey = "mock_secret_key_for_testing"
	config.Passphrase = "mock_passphrase_for_testing"
	config.DemoTrading = true
	config.UseTestnet = true
	return config
}

// LoadFromFile loads configuration from a JSON file
func (c *TestConfig) LoadFromFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, c)
}

// SaveToFile saves configuration to a JSON file
func (c *TestConfig) SaveToFile(filename string) error {
	// Create directory if it doesn't exist
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// IsEndpointEnabled checks if a specific endpoint test is enabled
func (c *TestConfig) IsEndpointEnabled(endpoint string) bool {
	if enabled, exists := c.EndpointTests[endpoint]; exists {
		return enabled
	}
	return false
}

// GetLogger creates a configured logger for integration tests
func (c *TestConfig) GetLogger() zerolog.Logger {
	level, err := zerolog.ParseLevel(c.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	}

	return zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("component", "integration-test").
		Logger()
}

// AccountTestEndpoints defines all account-related endpoint tests
var AccountTestEndpoints = []TestEndpoint{
	{
		Name:        "account_info",
		Endpoint:    "/api/v2/mix/account/account",
		Method:      "GET",
		Description: "Get account information and balances",
		Enabled:     true,
		ReadOnly:    true,
		Timeout:     10,
		Expected: map[string]string{
			"accountId":  "string",
			"usdtEquity": "string",
			"btcEquity":  "string",
		},
	},
	{
		Name:        "account_list",
		Endpoint:    "/api/v2/mix/account/accounts",
		Method:      "GET",
		Description: "Get all accounts list",
		Enabled:     true,
		ReadOnly:    true,
		Timeout:     10,
		Expected: map[string]string{
			"accountId":  "string",
			"marginCoin": "string",
		},
	},
	{
		Name:        "account_bills",
		Endpoint:    "/api/v2/mix/account/bill",
		Method:      "GET",
		Description: "Get account transaction history",
		Enabled:     true,
		ReadOnly:    true,
		Parameters: map[string]string{
			"limit": "10",
		},
		Timeout: 15,
		Expected: map[string]string{
			"billId": "string",
			"amount": "string",
			"fee":    "string",
		},
	},
	{
		Name:        "set_leverage",
		Endpoint:    "/api/v2/mix/account/set-leverage",
		Method:      "POST",
		Description: "Set leverage for symbol (has side effects)",
		Enabled:     false, // Disabled by default
		ReadOnly:    false,
		Parameters: map[string]string{
			"leverage": "20",
			"holdSide": "long",
		},
		Timeout: 10,
		Expected: map[string]string{
			"symbol":     "string",
			"marginCoin": "string",
			"longLever":  "string",
		},
	},
	{
		Name:        "set_margin_mode",
		Endpoint:    "/api/v2/mix/account/set-margin-mode",
		Method:      "POST",
		Description: "Set margin mode (isolated/cross) - has side effects",
		Enabled:     false, // Disabled by default
		ReadOnly:    false,
		Parameters: map[string]string{
			"marginMode": "isolated",
		},
		Timeout: 10,
		Expected: map[string]string{
			"symbol":     "string",
			"marginMode": "string",
		},
	},
	{
		Name:        "set_position_mode",
		Endpoint:    "/api/v2/mix/account/set-position-mode",
		Method:      "POST",
		Description: "Set position mode (one-way/hedge) - has side effects",
		Enabled:     false, // Disabled by default
		ReadOnly:    false,
		Parameters: map[string]string{
			"posMode": "one_way",
		},
		Timeout: 10,
		Expected: map[string]string{
			"productType":  "string",
			"positionMode": "string",
		},
	},
	{
		Name:        "adjust_margin",
		Endpoint:    "/api/v2/mix/account/set-margin",
		Method:      "POST",
		Description: "Adjust position margin - has side effects",
		Enabled:     false, // Disabled by default
		ReadOnly:    false,
		Parameters: map[string]string{
			"amount":   "1",
			"holdSide": "long",
			"type":     "add",
		},
		Timeout: 10,
		Expected: map[string]string{
			"symbol":     "string",
			"marginCoin": "string",
		},
	},
}
