// Example: Advanced Error Handling and Recovery Patterns for WebSocket
// This example demonstrates comprehensive error handling, reconnection strategies,
// data validation, and recovery mechanisms for production WebSocket usage.

package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

// ErrorType represents different types of WebSocket errors
type ErrorType int

const (
	ErrorTypeConnection ErrorType = iota
	ErrorTypeAuthentication
	ErrorTypeSubscription
	ErrorTypeParsing
	ErrorTypeRateLimit
	ErrorTypeNetwork
	ErrorTypeAPIError
)

// ErrorEvent represents an error event with context
type ErrorEvent struct {
	Type        ErrorType `json:"type"`
	Message     string    `json:"message"`
	Symbol      string    `json:"symbol,omitempty"`
	Channel     string    `json:"channel,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Recoverable bool      `json:"recoverable"`
	RetryCount  int       `json:"retryCount"`
}

// ConnectionState represents the current connection state
type ConnectionState int

const (
	StateDisconnected ConnectionState = iota
	StateConnecting
	StateConnected
	StateAuthenticating
	StateAuthenticated
	StateReconnecting
	StateError
)

// RecoveryStrategy defines how to handle different error types
type RecoveryStrategy struct {
	MaxRetries      int
	BackoffStrategy string // "linear", "exponential", "fixed"
	BaseDelay       time.Duration
	MaxDelay        time.Duration
	AutoReconnect   bool
}

// ErrorHandlingManager manages WebSocket errors and recovery
type ErrorHandlingManager struct {
	// WebSocket clients
	publicClient  *ws.BaseWsClient
	privateClient *ws.BaseWsClient
	
	// State management
	publicState    ConnectionState
	privateState   ConnectionState
	stateMutex     sync.RWMutex
	
	// Error tracking
	errorEvents    []ErrorEvent
	errorMutex     sync.RWMutex
	errorCount     map[ErrorType]int64
	lastErrorTime  map[ErrorType]time.Time
	
	// Recovery configuration
	strategies     map[ErrorType]RecoveryStrategy
	
	// Subscription management
	subscriptions  map[string]SubscriptionInfo
	subMutex      sync.RWMutex
	
	// Context and cancellation
	ctx           context.Context
	cancel        context.CancelFunc
	
	// Logging
	logger        zerolog.Logger
	
	// Statistics
	reconnectCount     int64
	recoveryCount      int64
	dataLossEvents     int64
	uptimeStart        time.Time
}

// SubscriptionInfo holds subscription details for recovery
type SubscriptionInfo struct {
	Channel     string        `json:"channel"`
	Symbol      string        `json:"symbol"`
	ProductType string        `json:"productType"`
	Timeframe   string        `json:"timeframe,omitempty"`
	Handler     ws.OnReceive  `json:"-"`
	LastMessage time.Time     `json:"lastMessage"`
	MessageCount int64        `json:"messageCount"`
	IsPrivate   bool          `json:"isPrivate"`
}

// NewErrorHandlingManager creates a new error handling manager
func NewErrorHandlingManager(logger zerolog.Logger) *ErrorHandlingManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	manager := &ErrorHandlingManager{
		publicState:    StateDisconnected,
		privateState:   StateDisconnected,
		errorEvents:    make([]ErrorEvent, 0, 1000),
		errorCount:     make(map[ErrorType]int64),
		lastErrorTime:  make(map[ErrorType]time.Time),
		subscriptions:  make(map[string]SubscriptionInfo),
		ctx:           ctx,
		cancel:        cancel,
		logger:        logger,
		uptimeStart:   time.Now(),
	}
	
	// Initialize recovery strategies
	manager.initializeRecoveryStrategies()
	
	return manager
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Create enhanced logger
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.DebugLevel)

	// Create error handling manager
	manager := NewErrorHandlingManager(logger)

	logger.Info().Msg("🛡️ Starting Advanced Error Handling and Recovery Demo")

	// Initialize and run demo
	if err := manager.Initialize(); err != nil {
		log.Fatal("Failed to initialize:", err)
	}

	// Start monitoring and recovery systems
	manager.StartMonitoring()

	// Run various error scenarios for demonstration
	go manager.simulateErrorScenarios()

	// Set up graceful shutdown
	manager.SetupGracefulShutdown()
}

// Initialize sets up the error handling manager
func (ehm *ErrorHandlingManager) Initialize() error {
	ehm.logger.Info().Msg("🔧 Initializing Error Handling Manager")

	// Start background monitoring
	go ehm.monitorConnections()
	go ehm.monitorSubscriptions()
	go ehm.monitorErrors()
	go ehm.displayErrorDashboard()

	// Connect to public WebSocket with error handling
	if err := ehm.connectPublicWithRetry(); err != nil {
		return fmt.Errorf("failed to connect to public WebSocket: %w", err)
	}

	// Attempt private connection (may fail with demo credentials)
	go func() {
		if err := ehm.connectPrivateWithRetry(); err != nil {
			ehm.logger.Warn().Err(err).Msg("Private WebSocket connection failed")
			ehm.recordError(ErrorTypeAuthentication, "Private connection failed", "", "", false)
		}
	}()

	ehm.logger.Info().Msg("✅ Error Handling Manager initialized")
	return nil
}

// initializeRecoveryStrategies sets up recovery strategies for different error types
func (ehm *ErrorHandlingManager) initializeRecoveryStrategies() {
	ehm.strategies = map[ErrorType]RecoveryStrategy{
		ErrorTypeConnection: {
			MaxRetries:      5,
			BackoffStrategy: "exponential",
			BaseDelay:       1 * time.Second,
			MaxDelay:        30 * time.Second,
			AutoReconnect:   true,
		},
		ErrorTypeAuthentication: {
			MaxRetries:      3,
			BackoffStrategy: "linear",
			BaseDelay:       2 * time.Second,
			MaxDelay:        10 * time.Second,
			AutoReconnect:   true,
		},
		ErrorTypeSubscription: {
			MaxRetries:      10,
			BackoffStrategy: "fixed",
			BaseDelay:       500 * time.Millisecond,
			MaxDelay:        5 * time.Second,
			AutoReconnect:   true,
		},
		ErrorTypeParsing: {
			MaxRetries:      0,
			BackoffStrategy: "fixed",
			BaseDelay:       0,
			MaxDelay:        0,
			AutoReconnect:   false,
		},
		ErrorTypeRateLimit: {
			MaxRetries:      3,
			BackoffStrategy: "exponential",
			BaseDelay:       5 * time.Second,
			MaxDelay:        60 * time.Second,
			AutoReconnect:   true,
		},
		ErrorTypeNetwork: {
			MaxRetries:      10,
			BackoffStrategy: "exponential",
			BaseDelay:       1 * time.Second,
			MaxDelay:        30 * time.Second,
			AutoReconnect:   true,
		},
		ErrorTypeAPIError: {
			MaxRetries:      2,
			BackoffStrategy: "linear",
			BaseDelay:       3 * time.Second,
			MaxDelay:        15 * time.Second,
			AutoReconnect:   false,
		},
	}
}

// Connection Management with Error Handling

func (ehm *ErrorHandlingManager) connectPublicWithRetry() error {
	ehm.setPublicState(StateConnecting)
	
	strategy := ehm.strategies[ErrorTypeConnection]
	var lastErr error
	
	for attempt := 0; attempt < strategy.MaxRetries; attempt++ {
		ehm.logger.Info().
			Int("attempt", attempt+1).
			Int("maxAttempts", strategy.MaxRetries).
			Msg("🔌 Attempting public WebSocket connection")
		
		if err := ehm.connectPublic(); err != nil {
			lastErr = err
			ehm.recordError(ErrorTypeConnection, err.Error(), "", "", true)
			
			if attempt < strategy.MaxRetries-1 {
				delay := ehm.calculateBackoffDelay(strategy, attempt)
				ehm.logger.Warn().
					Err(err).
					Dur("delay", delay).
					Msg("🔄 Connection failed, retrying...")
				
				select {
				case <-time.After(delay):
					continue
				case <-ehm.ctx.Done():
					return ehm.ctx.Err()
				}
			}
		} else {
			ehm.setPublicState(StateConnected)
			atomic.AddInt64(&ehm.reconnectCount, 1)
			ehm.logger.Info().Msg("✅ Public WebSocket connected successfully")
			return nil
		}
	}
	
	ehm.setPublicState(StateError)
	return fmt.Errorf("failed to connect after %d attempts: %w", strategy.MaxRetries, lastErr)
}

func (ehm *ErrorHandlingManager) connectPublic() error {
	ehm.publicClient = ws.NewBitgetBaseWsClient(
		ehm.logger,
		"wss://ws.bitget.com/v2/ws/public",
		"",
	)
	
	// Set enhanced error handlers
	ehm.publicClient.SetListener(ehm.createMessageHandler(false), ehm.createErrorHandler(false))
	
	ehm.publicClient.Connect()
	ehm.publicClient.ConnectWebSocket()
	ehm.publicClient.StartReadLoop()
	
	// Wait for connection with timeout
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			return errors.New("connection timeout")
		case <-ticker.C:
			if ehm.publicClient.IsConnected() {
				ehm.setupPublicSubscriptions()
				return nil
			}
		}
	}
}

func (ehm *ErrorHandlingManager) connectPrivateWithRetry() error {
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")
	
	if apiKey == "" || secretKey == "" || passphrase == "" {
		return fmt.Errorf("missing API credentials")
	}
	
	ehm.setPrivateState(StateConnecting)
	
	strategy := ehm.strategies[ErrorTypeAuthentication]
	var lastErr error
	
	for attempt := 0; attempt < strategy.MaxRetries; attempt++ {
		ehm.logger.Info().
			Int("attempt", attempt+1).
			Msg("🔐 Attempting private WebSocket connection")
		
		if err := ehm.connectPrivate(apiKey, secretKey, passphrase); err != nil {
			lastErr = err
			ehm.recordError(ErrorTypeAuthentication, err.Error(), "", "", true)
			
			if attempt < strategy.MaxRetries-1 {
				delay := ehm.calculateBackoffDelay(strategy, attempt)
				time.Sleep(delay)
			}
		} else {
			ehm.setPrivateState(StateAuthenticated)
			atomic.AddInt64(&ehm.reconnectCount, 1)
			ehm.logger.Info().Msg("✅ Private WebSocket authenticated successfully")
			return nil
		}
	}
	
	ehm.setPrivateState(StateError)
	return fmt.Errorf("failed to authenticate after %d attempts: %w", strategy.MaxRetries, lastErr)
}

func (ehm *ErrorHandlingManager) connectPrivate(apiKey, secretKey, passphrase string) error {
	ehm.privateClient = ws.NewBitgetBaseWsClient(
		ehm.logger,
		"wss://ws.bitget.com/v2/ws/private",
		secretKey,
	)
	
	ehm.privateClient.SetListener(ehm.createMessageHandler(true), ehm.createErrorHandler(true))
	
	ehm.privateClient.Connect()
	ehm.privateClient.ConnectWebSocket()
	ehm.privateClient.StartReadLoop()
	
	// Wait for connection
	time.Sleep(2 * time.Second)
	if !ehm.privateClient.IsConnected() {
		return errors.New("failed to connect")
	}
	
	ehm.setPrivateState(StateAuthenticating)
	ehm.privateClient.Login(apiKey, passphrase, common.SHA256)
	
	// Wait for authentication
	timeout := time.After(10 * time.Second)
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	
	for {
		select {
		case <-timeout:
			return errors.New("authentication timeout")
		case <-ticker.C:
			if ehm.privateClient.IsLoggedIn() {
				ehm.setupPrivateSubscriptions()
				return nil
			}
		}
	}
}

// Message and Error Handlers

func (ehm *ErrorHandlingManager) createMessageHandler(isPrivate bool) ws.OnReceive {
	return func(message string) {
		// Validate message format
		if err := ehm.validateMessage(message); err != nil {
			ehm.recordError(ErrorTypeParsing, err.Error(), "", "", false)
			return
		}
		
		// Update subscription statistics
		ehm.updateSubscriptionStats(message, isPrivate)
		
		// Process message
		ehm.processMessage(message, isPrivate)
	}
}

func (ehm *ErrorHandlingManager) createErrorHandler(isPrivate bool) ws.OnReceive {
	return func(message string) {
		ehm.logger.Error().
			Bool("private", isPrivate).
			Str("error", message).
			Msg("❌ WebSocket Error")
		
		// Analyze error message and determine type
		errorType := ehm.analyzeErrorMessage(message)
		
		// Record error
		ehm.recordError(errorType, message, "", "", true)
		
		// Trigger recovery if needed
		if ehm.shouldTriggerRecovery(errorType, message) {
			go ehm.triggerRecovery(errorType, isPrivate)
		}
	}
}

// Error Analysis and Classification

func (ehm *ErrorHandlingManager) analyzeErrorMessage(message string) ErrorType {
	messageLower := strings.ToLower(message)
	
	switch {
	case strings.Contains(messageLower, "connection") || strings.Contains(messageLower, "disconnect"):
		return ErrorTypeConnection
	case strings.Contains(messageLower, "auth") || strings.Contains(messageLower, "login"):
		return ErrorTypeAuthentication
	case strings.Contains(messageLower, "subscribe") || strings.Contains(messageLower, "channel"):
		return ErrorTypeSubscription
	case strings.Contains(messageLower, "rate limit") || strings.Contains(messageLower, "too many"):
		return ErrorTypeRateLimit
	case strings.Contains(messageLower, "network") || strings.Contains(messageLower, "timeout"):
		return ErrorTypeNetwork
	case strings.Contains(messageLower, "json") || strings.Contains(messageLower, "parse"):
		return ErrorTypeParsing
	default:
		return ErrorTypeAPIError
	}
}

func (ehm *ErrorHandlingManager) validateMessage(message string) error {
	// Basic JSON validation
	if !json.Valid([]byte(message)) {
		return fmt.Errorf("invalid JSON format")
	}
	
	// Check message structure
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	
	// Validate required fields (simplified)
	if _, exists := data["event"]; !exists {
		if _, exists := data["action"]; !exists {
			return fmt.Errorf("missing event/action field")
		}
	}
	
	return nil
}

// Subscription Management

func (ehm *ErrorHandlingManager) setupPublicSubscriptions() {
	symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}
	productType := "USDT-FUTURES"
	
	for _, symbol := range symbols {
		// Ticker subscription
		ehm.subscribeWithErrorHandling(
			symbol, "ticker", productType, "",
			ehm.createSubscriptionHandler(symbol, "ticker", false),
			false,
		)
		
		// Candlestick subscription  
		ehm.subscribeWithErrorHandling(
			symbol, "candle1m", productType, ws.Timeframe1m,
			ehm.createSubscriptionHandler(symbol, "candle1m", false),
			false,
		)
		
		// Order book subscription
		ehm.subscribeWithErrorHandling(
			symbol, "books5", productType, "",
			ehm.createSubscriptionHandler(symbol, "books5", false),
			false,
		)
		
		time.Sleep(300 * time.Millisecond)
	}
}

func (ehm *ErrorHandlingManager) setupPrivateSubscriptions() {
	productType := "USDT-FUTURES"
	
	channels := []string{"orders", "fill", "positions", "account"}
	
	for _, channel := range channels {
		ehm.subscribeWithErrorHandling(
			"", channel, productType, "",
			ehm.createSubscriptionHandler("", channel, true),
			true,
		)
		
		time.Sleep(200 * time.Millisecond)
	}
}

func (ehm *ErrorHandlingManager) subscribeWithErrorHandling(symbol, channel, productType, timeframe string, handler ws.OnReceive, isPrivate bool) {
	subscriptionKey := fmt.Sprintf("%s:%s:%s", channel, symbol, productType)
	
	// Store subscription info for recovery
	ehm.subMutex.Lock()
	ehm.subscriptions[subscriptionKey] = SubscriptionInfo{
		Channel:     channel,
		Symbol:      symbol,
		ProductType: productType,
		Timeframe:   timeframe,
		Handler:     handler,
		LastMessage: time.Now(),
		IsPrivate:   isPrivate,
	}
	ehm.subMutex.Unlock()
	
	// Attempt subscription with retry
	strategy := ehm.strategies[ErrorTypeSubscription]
	
	for attempt := 0; attempt < strategy.MaxRetries; attempt++ {
		var err error
		
		if isPrivate && ehm.privateClient != nil {
			err = ehm.subscribePrivateChannel(channel, productType, handler)
		} else if !isPrivate && ehm.publicClient != nil {
			err = ehm.subscribePublicChannel(channel, symbol, productType, timeframe, handler)
		} else {
			err = fmt.Errorf("client not available")
		}
		
		if err != nil {
			ehm.recordError(ErrorTypeSubscription, err.Error(), symbol, channel, true)
			
			if attempt < strategy.MaxRetries-1 {
				delay := ehm.calculateBackoffDelay(strategy, attempt)
				time.Sleep(delay)
			}
		} else {
			ehm.logger.Info().
				Str("channel", channel).
				Str("symbol", symbol).
				Bool("private", isPrivate).
				Msg("✅ Subscription successful")
			return
		}
	}
	
	ehm.logger.Error().
		Str("channel", channel).
		Str("symbol", symbol).
		Msg("❌ Failed to subscribe after retries")
}

func (ehm *ErrorHandlingManager) subscribePublicChannel(channel, symbol, productType, timeframe string, handler ws.OnReceive) error {
	switch channel {
	case "ticker":
		ehm.publicClient.SubscribeTicker(symbol, productType, handler)
	case "candle1m":
		ehm.publicClient.SubscribeCandles(symbol, productType, timeframe, handler)
	case "books5":
		ehm.publicClient.SubscribeOrderBook5(symbol, productType, handler)
	case "books15":
		ehm.publicClient.SubscribeOrderBook15(symbol, productType, handler)
	case "books":
		ehm.publicClient.SubscribeOrderBook(symbol, productType, handler)
	case "trade":
		ehm.publicClient.SubscribeTrades(symbol, productType, handler)
	case "mark-price":
		ehm.publicClient.SubscribeMarkPrice(symbol, productType, handler)
	case "funding-time":
		ehm.publicClient.SubscribeFundingTime(symbol, productType, handler)
	default:
		return fmt.Errorf("unknown public channel: %s", channel)
	}
	
	return nil
}

func (ehm *ErrorHandlingManager) subscribePrivateChannel(channel, productType string, handler ws.OnReceive) error {
	switch channel {
	case "orders":
		ehm.privateClient.SubscribeOrders(productType, handler)
	case "fill":
		ehm.privateClient.SubscribeFills(productType, handler)
	case "positions":
		ehm.privateClient.SubscribePositions(productType, handler)
	case "account":
		ehm.privateClient.SubscribeAccount(productType, handler)
	case "plan-order":
		ehm.privateClient.SubscribePlanOrders(productType, handler)
	default:
		return fmt.Errorf("unknown private channel: %s", channel)
	}
	
	return nil
}

func (ehm *ErrorHandlingManager) createSubscriptionHandler(symbol, channel string, isPrivate bool) ws.OnReceive {
	return func(message string) {
		subscriptionKey := fmt.Sprintf("%s:%s:USDT-FUTURES", channel, symbol)
		
		ehm.subMutex.Lock()
		if sub, exists := ehm.subscriptions[subscriptionKey]; exists {
			sub.LastMessage = time.Now()
			sub.MessageCount++
			ehm.subscriptions[subscriptionKey] = sub
		}
		ehm.subMutex.Unlock()
		
		ehm.logger.Debug().
			Str("channel", channel).
			Str("symbol", symbol).
			Bool("private", isPrivate).
			Str("data", message).
			Msg("📨 Message received")
	}
}

// Recovery System

func (ehm *ErrorHandlingManager) shouldTriggerRecovery(errorType ErrorType, message string) bool {
	strategy := ehm.strategies[errorType]
	
	if !strategy.AutoReconnect {
		return false
	}
	
	// Check if we're not already in recovery
	ehm.stateMutex.RLock()
	publicRecovering := ehm.publicState == StateReconnecting
	privateRecovering := ehm.privateState == StateReconnecting
	ehm.stateMutex.RUnlock()
	
	if publicRecovering || privateRecovering {
		return false
	}
	
	// Check error frequency to avoid recovery loops
	ehm.errorMutex.RLock()
	lastErrorTime, exists := ehm.lastErrorTime[errorType]
	ehm.errorMutex.RUnlock()
	
	if exists && time.Since(lastErrorTime) < 5*time.Second {
		return false // Too frequent, avoid recovery loop
	}
	
	return true
}

func (ehm *ErrorHandlingManager) triggerRecovery(errorType ErrorType, isPrivate bool) {
	atomic.AddInt64(&ehm.recoveryCount, 1)
	
	ehm.logger.Warn().
		Str("errorType", ehm.getErrorTypeName(errorType)).
		Bool("private", isPrivate).
		Msg("🔄 Triggering recovery")
	
	switch errorType {
	case ErrorTypeConnection:
		if isPrivate {
			ehm.recoverPrivateConnection()
		} else {
			ehm.recoverPublicConnection()
		}
	case ErrorTypeAuthentication:
		ehm.recoverPrivateConnection()
	case ErrorTypeSubscription:
		ehm.recoverSubscriptions(isPrivate)
	default:
		ehm.logger.Info().Msg("🔄 No specific recovery action for this error type")
	}
}

func (ehm *ErrorHandlingManager) recoverPublicConnection() {
	ehm.setPublicState(StateReconnecting)
	
	if ehm.publicClient != nil {
		ehm.publicClient.Close()
	}
	
	time.Sleep(2 * time.Second)
	
	if err := ehm.connectPublicWithRetry(); err != nil {
		ehm.logger.Error().Err(err).Msg("❌ Public connection recovery failed")
		ehm.setPublicState(StateError)
	}
}

func (ehm *ErrorHandlingManager) recoverPrivateConnection() {
	ehm.setPrivateState(StateReconnecting)
	
	if ehm.privateClient != nil {
		ehm.privateClient.Close()
	}
	
	time.Sleep(2 * time.Second)
	
	if err := ehm.connectPrivateWithRetry(); err != nil {
		ehm.logger.Error().Err(err).Msg("❌ Private connection recovery failed")
		ehm.setPrivateState(StateError)
	}
}

func (ehm *ErrorHandlingManager) recoverSubscriptions(isPrivate bool) {
	ehm.logger.Info().Bool("private", isPrivate).Msg("🔄 Recovering subscriptions")
	
	ehm.subMutex.RLock()
	subscriptionsToRecover := make([]SubscriptionInfo, 0)
	for _, sub := range ehm.subscriptions {
		if sub.IsPrivate == isPrivate {
			subscriptionsToRecover = append(subscriptionsToRecover, sub)
		}
	}
	ehm.subMutex.RUnlock()
	
	for _, sub := range subscriptionsToRecover {
		ehm.subscribeWithErrorHandling(
			sub.Symbol,
			sub.Channel,
			sub.ProductType,
			sub.Timeframe,
			sub.Handler,
			sub.IsPrivate,
		)
		time.Sleep(200 * time.Millisecond)
	}
}

// Utility Methods

func (ehm *ErrorHandlingManager) calculateBackoffDelay(strategy RecoveryStrategy, attempt int) time.Duration {
	switch strategy.BackoffStrategy {
	case "exponential":
		delay := strategy.BaseDelay * time.Duration(1<<uint(attempt))
		if delay > strategy.MaxDelay {
			delay = strategy.MaxDelay
		}
		return delay
	case "linear":
		delay := strategy.BaseDelay * time.Duration(attempt+1)
		if delay > strategy.MaxDelay {
			delay = strategy.MaxDelay
		}
		return delay
	case "fixed":
		return strategy.BaseDelay
	default:
		return strategy.BaseDelay
	}
}

func (ehm *ErrorHandlingManager) recordError(errorType ErrorType, message, symbol, channel string, recoverable bool) {
	ehm.errorMutex.Lock()
	defer ehm.errorMutex.Unlock()
	
	event := ErrorEvent{
		Type:        errorType,
		Message:     message,
		Symbol:      symbol,
		Channel:     channel,
		Timestamp:   time.Now(),
		Recoverable: recoverable,
		RetryCount:  0,
	}
	
	// Store error event (keep only last 1000)
	ehm.errorEvents = append(ehm.errorEvents, event)
	if len(ehm.errorEvents) > 1000 {
		ehm.errorEvents = ehm.errorEvents[1:]
	}
	
	// Update counters
	ehm.errorCount[errorType]++
	ehm.lastErrorTime[errorType] = time.Now()
}

func (ehm *ErrorHandlingManager) getErrorTypeName(errorType ErrorType) string {
	names := map[ErrorType]string{
		ErrorTypeConnection:     "Connection",
		ErrorTypeAuthentication: "Authentication",
		ErrorTypeSubscription:   "Subscription",
		ErrorTypeParsing:        "Parsing",
		ErrorTypeRateLimit:      "RateLimit",
		ErrorTypeNetwork:        "Network",
		ErrorTypeAPIError:      "APIError",
	}
	
	if name, exists := names[errorType]; exists {
		return name
	}
	return "Unknown"
}

// State Management

func (ehm *ErrorHandlingManager) setPublicState(state ConnectionState) {
	ehm.stateMutex.Lock()
	defer ehm.stateMutex.Unlock()
	ehm.publicState = state
}

func (ehm *ErrorHandlingManager) setPrivateState(state ConnectionState) {
	ehm.stateMutex.Lock()
	defer ehm.stateMutex.Unlock()
	ehm.privateState = state
}

func (ehm *ErrorHandlingManager) getConnectionStates() (ConnectionState, ConnectionState) {
	ehm.stateMutex.RLock()
	defer ehm.stateMutex.RUnlock()
	return ehm.publicState, ehm.privateState
}

func (ehm *ErrorHandlingManager) getConnectionStateString(state ConnectionState) string {
	states := map[ConnectionState]string{
		StateDisconnected:   "❌ Disconnected",
		StateConnecting:     "🔄 Connecting",
		StateConnected:      "✅ Connected",
		StateAuthenticating: "🔐 Authenticating",
		StateAuthenticated:  "✅ Authenticated",
		StateReconnecting:   "🔄 Reconnecting",
		StateError:         "❌ Error",
	}
	
	if str, exists := states[state]; exists {
		return str
	}
	return "❓ Unknown"
}

// Monitoring Systems

func (ehm *ErrorHandlingManager) StartMonitoring() {
	ehm.logger.Info().Msg("🔍 Starting monitoring systems")
	
	// All monitoring goroutines are already started in Initialize()
	ehm.logger.Info().Msg("✅ All monitoring systems active")
}

func (ehm *ErrorHandlingManager) monitorConnections() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			ehm.checkConnectionHealth()
		case <-ehm.ctx.Done():
			return
		}
	}
}

func (ehm *ErrorHandlingManager) checkConnectionHealth() {
	publicState, privateState := ehm.getConnectionStates()
	
	// Check public connection
	if ehm.publicClient != nil && publicState == StateConnected {
		if !ehm.publicClient.IsConnected() {
			ehm.logger.Warn().Msg("🔍 Public connection lost detected")
			ehm.recordError(ErrorTypeConnection, "Connection lost detected", "", "", true)
			go ehm.triggerRecovery(ErrorTypeConnection, false)
		}
	}
	
	// Check private connection
	if ehm.privateClient != nil && privateState == StateAuthenticated {
		if !ehm.privateClient.IsConnected() || !ehm.privateClient.IsLoggedIn() {
			ehm.logger.Warn().Msg("🔍 Private connection lost detected")
			ehm.recordError(ErrorTypeConnection, "Private connection lost", "", "", true)
			go ehm.triggerRecovery(ErrorTypeConnection, true)
		}
	}
}

func (ehm *ErrorHandlingManager) monitorSubscriptions() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			ehm.checkSubscriptionHealth()
		case <-ehm.ctx.Done():
			return
		}
	}
}

func (ehm *ErrorHandlingManager) checkSubscriptionHealth() {
	staleThreshold := 60 * time.Second
	now := time.Now()
	
	ehm.subMutex.RLock()
	for key, sub := range ehm.subscriptions {
		if now.Sub(sub.LastMessage) > staleThreshold && sub.MessageCount > 0 {
			ehm.logger.Warn().
				Str("subscription", key).
				Time("lastMessage", sub.LastMessage).
				Msg("🔍 Stale subscription detected")
			
			ehm.recordError(ErrorTypeSubscription, "Stale subscription", sub.Symbol, sub.Channel, true)
			atomic.AddInt64(&ehm.dataLossEvents, 1)
		}
	}
	ehm.subMutex.RUnlock()
}

func (ehm *ErrorHandlingManager) monitorErrors() {
	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			ehm.analyzeErrorPatterns()
		case <-ehm.ctx.Done():
			return
		}
	}
}

func (ehm *ErrorHandlingManager) analyzeErrorPatterns() {
	ehm.errorMutex.RLock()
	defer ehm.errorMutex.RUnlock()
	
	// Check for error spikes
	recentThreshold := 5 * time.Minute
	now := time.Now()
	recentErrors := 0
	
	for i := len(ehm.errorEvents) - 1; i >= 0; i-- {
		if now.Sub(ehm.errorEvents[i].Timestamp) <= recentThreshold {
			recentErrors++
		} else {
			break
		}
	}
	
	if recentErrors > 20 {
		ehm.logger.Warn().
			Int("recentErrors", recentErrors).
			Msg("🚨 High error rate detected")
	}
}

// Processing and Statistics

func (ehm *ErrorHandlingManager) processMessage(message string, isPrivate bool) {
	// Process the message (placeholder for actual business logic)
	ehm.logger.Debug().
		Bool("private", isPrivate).
		Msg("📨 Processing message")
}

func (ehm *ErrorHandlingManager) updateSubscriptionStats(message string, isPrivate bool) {
	// Update message count for relevant subscription
	// This is a simplified implementation
}

// Dashboard Display

func (ehm *ErrorHandlingManager) displayErrorDashboard() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			ehm.printErrorDashboard()
		case <-ehm.ctx.Done():
			return
		}
	}
}

func (ehm *ErrorHandlingManager) printErrorDashboard() {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("🛡️ ERROR HANDLING & RECOVERY DASHBOARD")
	fmt.Println(strings.Repeat("=", 100))
	
	// Connection Status
	publicState, privateState := ehm.getConnectionStates()
	fmt.Printf("🌍 Public WebSocket:   %s\n", ehm.getConnectionStateString(publicState))
	fmt.Printf("🔒 Private WebSocket:  %s\n", ehm.getConnectionStateString(privateState))
	
	// Subscription Status
	ehm.subMutex.RLock()
	activeSubscriptions := len(ehm.subscriptions)
	healthySubscriptions := 0
	staleThreshold := 60 * time.Second
	now := time.Now()
	
	for _, sub := range ehm.subscriptions {
		if now.Sub(sub.LastMessage) <= staleThreshold || sub.MessageCount == 0 {
			healthySubscriptions++
		}
	}
	ehm.subMutex.RUnlock()
	
	fmt.Printf("📡 Subscriptions:      %d active, %d healthy\n", activeSubscriptions, healthySubscriptions)
	
	// Error Statistics
	ehm.errorMutex.RLock()
	totalErrors := int64(0)
	for _, count := range ehm.errorCount {
		totalErrors += count
	}
	
	fmt.Printf("❌ Total Errors:       %d\n", totalErrors)
	fmt.Printf("🔄 Recovery Attempts:  %d\n", atomic.LoadInt64(&ehm.recoveryCount))
	fmt.Printf("🔌 Reconnections:      %d\n", atomic.LoadInt64(&ehm.reconnectCount))
	fmt.Printf("📉 Data Loss Events:   %d\n", atomic.LoadInt64(&ehm.dataLossEvents))
	
	// Error breakdown
	fmt.Println("\n📊 ERROR BREAKDOWN:")
	for errorType, count := range ehm.errorCount {
		if count > 0 {
			fmt.Printf("   %s: %d\n", ehm.getErrorTypeName(errorType), count)
		}
	}
	ehm.errorMutex.RUnlock()
	
	// Uptime
	uptime := time.Since(ehm.uptimeStart)
	fmt.Printf("\n⏱️  Total Uptime:      %v\n", uptime.Truncate(time.Second))
	
	// Recent errors
	fmt.Println("\n🔍 RECENT ERRORS (Last 5):")
	recentCount := 5
	if len(ehm.errorEvents) < recentCount {
		recentCount = len(ehm.errorEvents)
	}
	
	for i := len(ehm.errorEvents) - recentCount; i < len(ehm.errorEvents); i++ {
		event := ehm.errorEvents[i]
		fmt.Printf("   %s: %s [%s] %s\n",
			event.Timestamp.Format("15:04:05"),
			ehm.getErrorTypeName(event.Type),
			event.Symbol,
			event.Message,
		)
	}
	
	fmt.Println(strings.Repeat("=", 100))
}

// Error Simulation for Demo

func (ehm *ErrorHandlingManager) simulateErrorScenarios() {
	// Wait for initial setup
	time.Sleep(10 * time.Second)
	
	ehm.logger.Info().Msg("🎭 Starting error simulation scenarios")
	
	scenarios := []func(){
		ehm.simulateConnectionError,
		ehm.simulateSubscriptionError,
		ehm.simulateParsingError,
		ehm.simulateRateLimitError,
	}
	
	for _, scenario := range scenarios {
		time.Sleep(30 * time.Second)
		scenario()
	}
}

func (ehm *ErrorHandlingManager) simulateConnectionError() {
	ehm.logger.Info().Msg("🎭 Simulating connection error")
	ehm.recordError(ErrorTypeConnection, "Simulated connection loss", "", "", true)
}

func (ehm *ErrorHandlingManager) simulateSubscriptionError() {
	ehm.logger.Info().Msg("🎭 Simulating subscription error")
	ehm.recordError(ErrorTypeSubscription, "Simulated subscription failure", "BTCUSDT", "ticker", true)
}

func (ehm *ErrorHandlingManager) simulateParsingError() {
	ehm.logger.Info().Msg("🎭 Simulating parsing error")
	ehm.recordError(ErrorTypeParsing, "Simulated JSON parsing error", "ETHUSDT", "candle", false)
}

func (ehm *ErrorHandlingManager) simulateRateLimitError() {
	ehm.logger.Info().Msg("🎭 Simulating rate limit error")
	ehm.recordError(ErrorTypeRateLimit, "Simulated rate limit exceeded", "", "", true)
}

// Graceful Shutdown

func (ehm *ErrorHandlingManager) SetupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	ehm.logger.Info().Msg("🛡️ Error Handling Demo running...")
	ehm.logger.Info().Msg("💡 Features: Auto-recovery, Error analysis, Connection monitoring")
	ehm.logger.Info().Msg("🎭 Error scenarios will be simulated for demonstration")
	ehm.logger.Info().Msg("🎧 Press Ctrl+C to stop")

	<-sigChan

	ehm.logger.Info().Msg("🛑 Shutting down error handling demo...")

	// Cancel context
	ehm.cancel()

	// Close connections
	if ehm.publicClient != nil {
		ehm.publicClient.Close()
	}
	if ehm.privateClient != nil {
		ehm.privateClient.Close()
	}

	// Display final report
	ehm.printFinalErrorReport()

	ehm.logger.Info().Msg("✅ Error handling demo shutdown complete")
}

func (ehm *ErrorHandlingManager) printFinalErrorReport() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 FINAL ERROR HANDLING REPORT")
	fmt.Println(strings.Repeat("=", 80))
	
	uptime := time.Since(ehm.uptimeStart)
	
	ehm.errorMutex.RLock()
	totalErrors := int64(0)
	for _, count := range ehm.errorCount {
		totalErrors += count
	}
	ehm.errorMutex.RUnlock()
	
	fmt.Printf("Total Runtime:         %v\n", uptime.Truncate(time.Second))
	fmt.Printf("Total Errors:          %d\n", totalErrors)
	fmt.Printf("Recovery Attempts:     %d\n", atomic.LoadInt64(&ehm.recoveryCount))
	fmt.Printf("Successful Reconnects: %d\n", atomic.LoadInt64(&ehm.reconnectCount))
	fmt.Printf("Data Loss Events:      %d\n", atomic.LoadInt64(&ehm.dataLossEvents))
	
	if totalErrors > 0 {
		fmt.Printf("Average MTBF:          %v\n", time.Duration(uptime.Nanoseconds()/totalErrors))
	}
	
	fmt.Println(strings.Repeat("=", 80))
}