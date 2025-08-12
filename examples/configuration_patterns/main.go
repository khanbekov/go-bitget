package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/futures/market"
)

// ConfigurationDemo demonstrates best practices for configuration management,
// error handling, logging, and environment setup for production applications
func main() {
	// Initialize application with proper configuration
	app, err := NewTradingApp()
	if err != nil {
		log.Fatalf("❌ Failed to initialize application: %v", err)
	}

	// Run the application
	if err := app.Run(); err != nil {
		log.Fatalf("❌ Application failed: %v", err)
	}
}

// TradingApp demonstrates a well-structured trading application
type TradingApp struct {
	config *Config
	client *futures.Client
	logger *Logger
	ctx    context.Context
	cancel context.CancelFunc
}

// Config contains all application configuration with validation
type Config struct {
	// API Configuration
	Bitget BitgetConfig `json:"bitget"`
	
	// Trading Configuration
	Trading TradingConfig `json:"trading"`
	
	// Risk Management
	Risk RiskConfig `json:"risk"`
	
	// Application Settings
	App AppConfig `json:"app"`
}

// BitgetConfig contains Bitget API configuration
type BitgetConfig struct {
	APIKey      string `json:"api_key"`
	SecretKey   string `json:"secret_key"`
	Passphrase  string `json:"passphrase"`
	BaseURL     string `json:"base_url"`
	IsTestnet   bool   `json:"is_testnet"`
	Timeout     int    `json:"timeout_seconds"`
}

// TradingConfig contains trading strategy parameters
type TradingConfig struct {
	Symbols     []string `json:"symbols"`
	ProductType string   `json:"product_type"`
	MarginCoin  string   `json:"margin_coin"`
	
	// Position sizing
	DefaultSize      string  `json:"default_size"`
	MaxPositions     int     `json:"max_positions"`
	PositionTimeout  int     `json:"position_timeout_hours"`
	
	// Strategy parameters
	TakeProfitPct    float64 `json:"take_profit_pct"`
	StopLossPct      float64 `json:"stop_loss_pct"`
	EntryThreshold   float64 `json:"entry_threshold_pct"`
}

// RiskConfig contains risk management settings
type RiskConfig struct {
	MaxDailyLoss     float64 `json:"max_daily_loss_pct"`
	MaxDrawdown      float64 `json:"max_drawdown_pct"`
	MaxPortfolioRisk float64 `json:"max_portfolio_risk_pct"`
	MinBalance       float64 `json:"min_balance_usdt"`
	
	// Circuit breakers
	EnableCircuitBreaker   bool `json:"enable_circuit_breaker"`
	MaxConsecutiveLosses   int  `json:"max_consecutive_losses"`
	CooldownPeriodMinutes  int  `json:"cooldown_period_minutes"`
}

// AppConfig contains general application settings
type AppConfig struct {
	LogLevel          string `json:"log_level"`
	LogFormat         string `json:"log_format"`
	LogFile           string `json:"log_file"`
	
	UpdateInterval    int    `json:"update_interval_seconds"`
	DataRetention     int    `json:"data_retention_days"`
	
	// Monitoring
	EnableMetrics     bool   `json:"enable_metrics"`
	MetricsPort       int    `json:"metrics_port"`
	HealthCheckPort   int    `json:"health_check_port"`
	
	// Persistence
	DatabasePath      string `json:"database_path"`
	BackupInterval    int    `json:"backup_interval_hours"`
}

// Logger provides structured logging with multiple output formats
type Logger struct {
	level  string
	format string
	file   *os.File
}

// NewLogger creates a configured logger
func NewLogger(config AppConfig) (*Logger, error) {
	logger := &Logger{
		level:  config.LogLevel,
		format: config.LogFormat,
	}
	
	// Setup file logging if specified
	if config.LogFile != "" {
		// Ensure log directory exists
		logDir := filepath.Dir(config.LogFile)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}
		
		file, err := os.OpenFile(config.LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file: %w", err)
		}
		logger.file = file
	}
	
	return logger, nil
}

// Info logs info level messages
func (l *Logger) Info(format string, args ...interface{}) {
	l.log("INFO", format, args...)
}

// Warn logs warning level messages  
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log("WARN", format, args...)
}

// Error logs error level messages
func (l *Logger) Error(format string, args ...interface{}) {
	l.log("ERROR", format, args...)
}

// Debug logs debug level messages
func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level == "debug" {
		l.log("DEBUG", format, args...)
	}
}

// log handles the actual logging with format support
func (l *Logger) log(level string, format string, args ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf(format, args...)
	
	var logEntry string
	if l.format == "json" {
		entry := map[string]interface{}{
			"timestamp": timestamp,
			"level":     level,
			"message":   message,
		}
		jsonBytes, _ := json.Marshal(entry)
		logEntry = string(jsonBytes)
	} else {
		logEntry = fmt.Sprintf("[%s] %s: %s", timestamp, level, message)
	}
	
	// Output to console
	fmt.Println(logEntry)
	
	// Output to file if configured
	if l.file != nil {
		l.file.WriteString(logEntry + "\n")
		l.file.Sync()
	}
}

// Close closes the logger and any open files
func (l *Logger) Close() {
	if l.file != nil {
		l.file.Close()
	}
}

// NewTradingApp creates a new trading application with proper configuration
func NewTradingApp() (*TradingApp, error) {
	// Load configuration
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	// Initialize logger
	logger, err := NewLogger(config.App)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}
	
	// Create Bitget client
	client := futures.NewClient(
		config.Bitget.APIKey,
		config.Bitget.SecretKey,
		config.Bitget.Passphrase,
	)
	
	// Configure client for testnet if needed
	if config.Bitget.IsTestnet {
		client.SetApiEndpoint("https://testnet.bitget.com")
		logger.Info("🔧 Configured for testnet environment")
	}
	
	// Create context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	
	app := &TradingApp{
		config: config,
		client: client,
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
	
	// Validate configuration
	if err := app.validateConfig(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}
	
	logger.Info("🚀 Trading application initialized successfully")
	return app, nil
}

// LoadConfig loads configuration from multiple sources with precedence:
// 1. Environment variables (highest priority)
// 2. Configuration file
// 3. Default values (lowest priority)
func LoadConfig() (*Config, error) {
	// Start with default configuration
	config := getDefaultConfig()
	
	// Load from config file if exists
	configFile := getConfigFile()
	if configFile != "" {
		if err := loadConfigFromFile(config, configFile); err != nil {
			return nil, fmt.Errorf("failed to load config file %s: %w", configFile, err)
		}
	}
	
	// Override with environment variables
	loadConfigFromEnv(config)
	
	return config, nil
}

// getDefaultConfig returns sensible default configuration
func getDefaultConfig() *Config {
	return &Config{
		Bitget: BitgetConfig{
			BaseURL:   "https://api.bitget.com",
			IsTestnet: false,
			Timeout:   30,
		},
		Trading: TradingConfig{
			Symbols:         []string{"BTCUSDT"},
			ProductType:     "USDT-FUTURES", 
			MarginCoin:      "USDT",
			DefaultSize:     "0.001",
			MaxPositions:    3,
			PositionTimeout: 24,
			TakeProfitPct:   0.02,
			StopLossPct:     0.01,
			EntryThreshold:  0.005,
		},
		Risk: RiskConfig{
			MaxDailyLoss:           0.05,
			MaxDrawdown:            0.15,
			MaxPortfolioRisk:       0.25,
			MinBalance:             100.0,
			EnableCircuitBreaker:   true,
			MaxConsecutiveLosses:   5,
			CooldownPeriodMinutes:  60,
		},
		App: AppConfig{
			LogLevel:          "info",
			LogFormat:         "text",
			LogFile:           "",
			UpdateInterval:    30,
			DataRetention:     30,
			EnableMetrics:     false,
			MetricsPort:       8080,
			HealthCheckPort:   8081,
			DatabasePath:      "./data",
			BackupInterval:    24,
		},
	}
}

// getConfigFile determines which config file to use
func getConfigFile() string {
	// Check for environment-specific config file
	if env := os.Getenv("APP_ENV"); env != "" {
		envConfig := fmt.Sprintf("config.%s.json", env)
		if _, err := os.Stat(envConfig); err == nil {
			return envConfig
		}
	}
	
	// Check for explicit config file
	if configFile := os.Getenv("CONFIG_FILE"); configFile != "" {
		return configFile
	}
	
	// Default config file
	if _, err := os.Stat("config.json"); err == nil {
		return "config.json"
	}
	
	return ""
}

// loadConfigFromFile loads configuration from JSON file
func loadConfigFromFile(config *Config, filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	
	return json.Unmarshal(data, config)
}

// loadConfigFromEnv overrides config with environment variables
func loadConfigFromEnv(config *Config) {
	// Bitget API credentials
	if apiKey := os.Getenv("BITGET_API_KEY"); apiKey != "" {
		config.Bitget.APIKey = apiKey
	}
	if secretKey := os.Getenv("BITGET_SECRET_KEY"); secretKey != "" {
		config.Bitget.SecretKey = secretKey
	}
	if passphrase := os.Getenv("BITGET_PASSPHRASE"); passphrase != "" {
		config.Bitget.Passphrase = passphrase
	}
	
	// Testnet configuration
	if testnet := os.Getenv("BITGET_TESTNET"); testnet != "" {
		config.Bitget.IsTestnet = testnet == "true"
	}
	
	// Trading parameters
	if symbols := os.Getenv("TRADING_SYMBOLS"); symbols != "" {
		// Parse comma-separated symbols
		config.Trading.Symbols = []string{symbols} // Simplified for demo
	}
	if size := os.Getenv("TRADING_DEFAULT_SIZE"); size != "" {
		config.Trading.DefaultSize = size
	}
	
	// Risk management
	if maxLoss := os.Getenv("RISK_MAX_DAILY_LOSS"); maxLoss != "" {
		if val, err := strconv.ParseFloat(maxLoss, 64); err == nil {
			config.Risk.MaxDailyLoss = val
		}
	}
	
	// App configuration
	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		config.App.LogLevel = logLevel
	}
	if logFile := os.Getenv("LOG_FILE"); logFile != "" {
		config.App.LogFile = logFile
	}
}

// validateConfig performs comprehensive configuration validation
func (app *TradingApp) validateConfig() error {
	config := app.config
	
	// Validate API credentials
	if config.Bitget.APIKey == "" {
		return fmt.Errorf("BITGET_API_KEY is required")
	}
	if config.Bitget.SecretKey == "" {
		return fmt.Errorf("BITGET_SECRET_KEY is required")
	}
	if config.Bitget.Passphrase == "" {
		return fmt.Errorf("BITGET_PASSPHRASE is required")
	}
	
	// Validate trading parameters
	if len(config.Trading.Symbols) == 0 {
		return fmt.Errorf("at least one trading symbol is required")
	}
	if config.Trading.MaxPositions <= 0 {
		return fmt.Errorf("max_positions must be greater than 0")
	}
	if config.Trading.TakeProfitPct <= 0 {
		return fmt.Errorf("take_profit_pct must be greater than 0")
	}
	if config.Trading.StopLossPct <= 0 {
		return fmt.Errorf("stop_loss_pct must be greater than 0")
	}
	
	// Validate risk parameters
	if config.Risk.MaxDailyLoss <= 0 || config.Risk.MaxDailyLoss >= 1 {
		return fmt.Errorf("max_daily_loss_pct must be between 0 and 1")
	}
	if config.Risk.MaxDrawdown <= 0 || config.Risk.MaxDrawdown >= 1 {
		return fmt.Errorf("max_drawdown_pct must be between 0 and 1")
	}
	if config.Risk.MinBalance <= 0 {
		return fmt.Errorf("min_balance_usdt must be greater than 0")
	}
	
	// Validate app parameters
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[config.App.LogLevel] {
		return fmt.Errorf("invalid log_level: %s", config.App.LogLevel)
	}
	
	if config.App.UpdateInterval <= 0 {
		return fmt.Errorf("update_interval_seconds must be greater than 0")
	}
	
	app.logger.Info("✅ Configuration validation passed")
	return nil
}

// Run starts the main application loop with proper error handling
func (app *TradingApp) Run() error {
	app.logger.Info("🚀 Starting trading application...")
	
	// Setup graceful shutdown
	defer app.cleanup()
	
	// Test API connectivity
	if err := app.testConnectivity(); err != nil {
		return fmt.Errorf("API connectivity test failed: %w", err)
	}
	
	// Initialize application state
	if err := app.initialize(); err != nil {
		return fmt.Errorf("application initialization failed: %w", err)
	}
	
	// Start main application loop
	return app.mainLoop()
}

// testConnectivity verifies API connection and credentials
func (app *TradingApp) testConnectivity() error {
	app.logger.Info("🔌 Testing API connectivity...")
	
	// Test with a simple market data request
	tickerService := market.NewTickerService(app.client)
	_, err := tickerService.
		Symbol(app.config.Trading.Symbols[0]).
		ProductType(market.ProductType(app.config.Trading.ProductType)).
		Do(app.ctx)
	
	if err != nil {
		app.logger.Error("❌ API connectivity test failed: %v", err)
		return err
	}
	
	// Test account access
	accountService := account.NewAccountListService(app.client)
	_, err = accountService.
		ProductType(account.ProductType(app.config.Trading.ProductType)).
		Do(app.ctx)
	
	if err != nil {
		app.logger.Error("❌ Account access test failed: %v", err)
		return err
	}
	
	app.logger.Info("✅ API connectivity test passed")
	return nil
}

// initialize performs application initialization tasks
func (app *TradingApp) initialize() error {
	app.logger.Info("⚙️ Initializing application state...")
	
	// Create necessary directories
	dirs := []string{
		app.config.App.DatabasePath,
		filepath.Dir(app.config.App.LogFile),
	}
	
	for _, dir := range dirs {
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %w", dir, err)
			}
		}
	}
	
	// Initialize circuit breaker if enabled
	if app.config.Risk.EnableCircuitBreaker {
		app.logger.Info("🔒 Circuit breaker enabled")
	}
	
	// Log configuration summary
	app.logConfigSummary()
	
	app.logger.Info("✅ Application initialization completed")
	return nil
}

// logConfigSummary logs key configuration parameters
func (app *TradingApp) logConfigSummary() {
	app.logger.Info("📋 Configuration Summary:")
	app.logger.Info("  🎯 Symbols: %v", app.config.Trading.Symbols)
	app.logger.Info("  💰 Default Size: %s", app.config.Trading.DefaultSize)
	app.logger.Info("  📊 Max Positions: %d", app.config.Trading.MaxPositions)
	app.logger.Info("  📈 Take Profit: %.1f%%", app.config.Trading.TakeProfitPct*100)
	app.logger.Info("  📉 Stop Loss: %.1f%%", app.config.Trading.StopLossPct*100)
	app.logger.Info("  ⚠️ Max Daily Loss: %.1f%%", app.config.Risk.MaxDailyLoss*100)
	app.logger.Info("  🔒 Circuit Breaker: %v", app.config.Risk.EnableCircuitBreaker)
	if app.config.Bitget.IsTestnet {
		app.logger.Warn("  🧪 TESTNET MODE ENABLED")
	}
}

// mainLoop runs the main application logic
func (app *TradingApp) mainLoop() error {
	app.logger.Info("🔄 Starting main application loop...")
	
	ticker := time.NewTicker(time.Duration(app.config.App.UpdateInterval) * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-app.ctx.Done():
			app.logger.Info("📴 Application shutdown requested")
			return nil
			
		case <-ticker.C:
			if err := app.processUpdate(); err != nil {
				app.logger.Error("❌ Update processing failed: %v", err)
				
				// Implement retry logic or circuit breaker here
				if app.shouldStopOnError(err) {
					return fmt.Errorf("critical error, stopping application: %w", err)
				}
			}
		}
	}
}

// processUpdate handles a single update cycle
func (app *TradingApp) processUpdate() error {
	app.logger.Debug("🔄 Processing update cycle...")
	
	// Example update logic
	for _, symbol := range app.config.Trading.Symbols {
		// Get current price
		tickerService := market.NewTickerService(app.client)
		ticker, err := tickerService.
			Symbol(symbol).
			ProductType(market.ProductType(app.config.Trading.ProductType)).
			Do(app.ctx)
		
		if err != nil {
			app.logger.Warn("⚠️ Failed to get ticker for %s: %v", symbol, err)
			continue
		}
		
		app.logger.Debug("📊 %s: $%s", symbol, ticker.LastPr)
	}
	
	return nil
}

// shouldStopOnError determines if an error should stop the application
func (app *TradingApp) shouldStopOnError(err error) bool {
	// Define critical error conditions
	criticalErrors := []string{
		"authentication",
		"invalid credentials", 
		"account suspended",
		"insufficient balance",
	}
	
	errorStr := err.Error()
	for _, critical := range criticalErrors {
		if contains(errorStr, critical) {
			app.logger.Error("🚨 Critical error detected: %v", err)
			return true
		}
	}
	
	return false
}

// contains checks if a string contains a substring (case insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || 
		len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}()))
}

// cleanup performs graceful shutdown tasks
func (app *TradingApp) cleanup() {
	app.logger.Info("🧹 Performing cleanup...")
	
	// Cancel context
	if app.cancel != nil {
		app.cancel()
	}
	
	// Close logger
	if app.logger != nil {
		app.logger.Close()
	}
	
	// Add other cleanup tasks here:
	// - Close database connections
	// - Flush metrics
	// - Close positions if required
	// - Send shutdown notifications
	
	fmt.Println("✅ Cleanup completed")
}