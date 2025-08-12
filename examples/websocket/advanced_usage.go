// Example: Advanced WebSocket Usage Patterns
// This example demonstrates advanced patterns including error handling,
// reconnection strategies, data processing, and performance optimization.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

// AdvancedTicker represents structured ticker data
type AdvancedTicker struct {
	Symbol    string    `json:"instId"`
	LastPrice string    `json:"last"`
	Volume24h string    `json:"vol24h"`
	Change24h string    `json:"change24h"`
	Timestamp time.Time `json:"ts"`
}

// MarketDataProcessor handles real-time market data processing
type MarketDataProcessor struct {
	mu             sync.RWMutex
	symbols        map[string]*AdvancedTicker
	updateCount    int64
	errorCount     int64
	reconnectCount int64

	// Channels for data streaming
	tickerChan chan *AdvancedTicker
	errorChan  chan error

	// Context for cancellation
	ctx    context.Context
	cancel context.CancelFunc

	// WebSocket client
	client *ws.BaseWsClient
	logger zerolog.Logger
}

// NewMarketDataProcessor creates a new market data processor
func NewMarketDataProcessor(logger zerolog.Logger) *MarketDataProcessor {
	ctx, cancel := context.WithCancel(context.Background())

	return &MarketDataProcessor{
		symbols:    make(map[string]*AdvancedTicker),
		tickerChan: make(chan *AdvancedTicker, 1000), // Buffered channel
		errorChan:  make(chan error, 100),
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
	}
}

// Start initializes and starts the market data processor
func (m *MarketDataProcessor) Start() error {
	// Create WebSocket client with retry configuration
	m.client = ws.NewBitgetBaseWsClient(
		m.logger,
		"wss://ws.bitget.com/v2/ws/public",
		"",
	)

	// Set custom reconnection settings
	m.client.SetReconnectionTimeout(30 * time.Second)
	m.client.SetCheckConnectionInterval(5 * time.Second)

	// Set message handlers with error recovery
	m.client.SetListener(m.messageHandler, m.errorHandler)

	// Connect with retry logic
	if err := m.connectWithRetry(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	// Start background goroutines
	go m.processTickerData()
	go m.processErrors()
	go m.monitorConnection()
	go m.printPerformanceStats()

	return nil
}

// connectWithRetry attempts to connect with exponential backoff
func (m *MarketDataProcessor) connectWithRetry() error {
	maxRetries := 5
	backoff := time.Second

	for i := 0; i < maxRetries; i++ {
		m.logger.Info().Msgf("Connection attempt %d/%d", i+1, maxRetries)

		m.client.Connect()
		m.client.ConnectWebSocket()
		m.client.StartReadLoop()

		// Wait for connection
		time.Sleep(2 * time.Second)

		if m.client.IsConnected() {
			m.logger.Info().Msg("Successfully connected to WebSocket")
			return nil
		}

		m.logger.Warn().Msgf("Connection attempt %d failed, retrying in %v", i+1, backoff)
		time.Sleep(backoff)
		backoff *= 2 // Exponential backoff
	}

	return fmt.Errorf("failed to connect after %d attempts", maxRetries)
}

// SubscribeToSymbols subscribes to multiple symbols with optimized handlers
func (m *MarketDataProcessor) SubscribeToSymbols(symbols []string) {
	productType := "USDT-FUTURES"

	for _, symbol := range symbols {
		// Initialize symbol data
		m.mu.Lock()
		m.symbols[symbol] = &AdvancedTicker{Symbol: symbol}
		m.mu.Unlock()

		// Create optimized ticker handler
		tickerHandler := m.createOptimizedTickerHandler(symbol)

		// Subscribe to ticker with error handling
		func(sym string) {
			defer func() {
				if r := recover(); r != nil {
					m.logger.Error().Msgf("Panic in subscription for %s: %v", sym, r)
				}
			}()

			m.client.SubscribeTicker(sym, productType, tickerHandler)
			m.logger.Info().Msgf("Subscribed to ticker for %s", sym)
		}(symbol)

		// Add small delay to avoid overwhelming the server
		time.Sleep(100 * time.Millisecond)
	}

	m.logger.Info().Msgf("Subscribed to %d symbols", len(symbols))
}

// createOptimizedTickerHandler creates an optimized ticker handler
func (m *MarketDataProcessor) createOptimizedTickerHandler(symbol string) ws.OnReceive {
	return func(message string) {
		// Increment update counter atomically
		m.mu.Lock()
		m.updateCount++
		m.mu.Unlock()

		// Parse ticker data in background
		go func() {
			defer func() {
				if r := recover(); r != nil {
					m.errorChan <- fmt.Errorf("ticker handler panic for %s: %v", symbol, r)
				}
			}()

			ticker, err := m.parseTickerMessage(symbol, message)
			if err != nil {
				m.errorChan <- fmt.Errorf("failed to parse ticker for %s: %w", symbol, err)
				return
			}

			// Send to processing channel (non-blocking)
			select {
			case m.tickerChan <- ticker:
			default:
				m.logger.Warn().Msg("Ticker channel full, dropping message")
			}
		}()
	}
}

// parseTickerMessage parses incoming ticker messages
func (m *MarketDataProcessor) parseTickerMessage(symbol, message string) (*AdvancedTicker, error) {
	// This is a simplified parser - in reality you'd parse the actual JSON structure
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(message), &data); err != nil {
		return nil, err
	}

	ticker := &AdvancedTicker{
		Symbol:    symbol,
		LastPrice: "50000.00", // Placeholder - would extract from actual data
		Volume24h: "1000.0",   // Placeholder
		Change24h: "+2.5%",    // Placeholder
		Timestamp: time.Now(),
	}

	return ticker, nil
}

// processTickerData processes incoming ticker data
func (m *MarketDataProcessor) processTickerData() {
	for {
		select {
		case ticker := <-m.tickerChan:
			m.updateSymbolData(ticker)
		case <-m.ctx.Done():
			return
		}
	}
}

// updateSymbolData updates symbol data thread-safely
func (m *MarketDataProcessor) updateSymbolData(ticker *AdvancedTicker) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if existing, exists := m.symbols[ticker.Symbol]; exists {
		existing.LastPrice = ticker.LastPrice
		existing.Volume24h = ticker.Volume24h
		existing.Change24h = ticker.Change24h
		existing.Timestamp = ticker.Timestamp
	}
}

// processErrors handles and logs errors
func (m *MarketDataProcessor) processErrors() {
	for {
		select {
		case err := <-m.errorChan:
			m.mu.Lock()
			m.errorCount++
			m.mu.Unlock()

			m.logger.Error().Err(err).Msg("Processing error")

			// Implement error-specific handling
			if m.shouldReconnect(err) {
				go m.handleReconnection()
			}

		case <-m.ctx.Done():
			return
		}
	}
}

// shouldReconnect determines if an error requires reconnection
func (m *MarketDataProcessor) shouldReconnect(err error) bool {
	// Implement logic to determine if reconnection is needed
	// This is a simplified example
	return false
}

// handleReconnection handles reconnection logic
func (m *MarketDataProcessor) handleReconnection() {
	m.mu.Lock()
	m.reconnectCount++
	m.mu.Unlock()

	m.logger.Warn().Msg("Attempting reconnection...")

	// Close existing connection
	m.client.Close()

	// Wait before reconnecting
	time.Sleep(5 * time.Second)

	// Reconnect
	if err := m.connectWithRetry(); err != nil {
		m.logger.Error().Err(err).Msg("Reconnection failed")
		return
	}

	// Resubscribe to all symbols
	symbols := make([]string, 0, len(m.symbols))
	m.mu.RLock()
	for symbol := range m.symbols {
		symbols = append(symbols, symbol)
	}
	m.mu.RUnlock()

	m.SubscribeToSymbols(symbols)
	m.logger.Info().Msg("Reconnection successful")
}

// monitorConnection monitors connection health
func (m *MarketDataProcessor) monitorConnection() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !m.client.IsConnected() {
				m.logger.Warn().Msg("Connection lost, attempting reconnection")
				go m.handleReconnection()
			}
		case <-m.ctx.Done():
			return
		}
	}
}

// printPerformanceStats prints performance statistics periodically
func (m *MarketDataProcessor) printPerformanceStats() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			m.mu.RLock()
			stats := fmt.Sprintf(
				"Updates: %d, Errors: %d, Reconnects: %d, Subscriptions: %d",
				m.updateCount,
				m.errorCount,
				m.reconnectCount,
				m.client.GetSubscriptionCount(),
			)
			m.mu.RUnlock()

			m.logger.Info().Msg("Performance Stats: " + stats)

		case <-m.ctx.Done():
			return
		}
	}
}

// GetMarketSummary returns a summary of all tracked symbols
func (m *MarketDataProcessor) GetMarketSummary() map[string]*AdvancedTicker {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to avoid race conditions
	summary := make(map[string]*AdvancedTicker)
	for symbol, ticker := range m.symbols {
		// Create a copy of the ticker
		summary[symbol] = &AdvancedTicker{
			Symbol:    ticker.Symbol,
			LastPrice: ticker.LastPrice,
			Volume24h: ticker.Volume24h,
			Change24h: ticker.Change24h,
			Timestamp: ticker.Timestamp,
		}
	}

	return summary
}

// messageHandler handles general WebSocket messages
func (m *MarketDataProcessor) messageHandler(message string) {
	// Handle general messages (like ping/pong, status updates)
	m.logger.Debug().Msgf("General message: %s", message)
}

// errorHandler handles WebSocket errors
func (m *MarketDataProcessor) errorHandler(message string) {
	err := fmt.Errorf("websocket error: %s", message)

	select {
	case m.errorChan <- err:
	default:
		m.logger.Error().Msg("Error channel full, dropping error")
	}
}

// Stop gracefully stops the market data processor
func (m *MarketDataProcessor) Stop() {
	m.logger.Info().Msg("Stopping market data processor...")

	// Cancel context to stop goroutines
	m.cancel()

	// Unsubscribe from all symbols
	m.unsubscribeAll()

	// Close WebSocket connection
	if m.client != nil {
		m.client.Close()
	}

	// Close channels
	close(m.tickerChan)
	close(m.errorChan)

	m.logger.Info().Msg("Market data processor stopped")
}

// unsubscribeAll unsubscribes from all symbols
func (m *MarketDataProcessor) unsubscribeAll() {
	productType := "USDT-FUTURES"

	m.mu.RLock()
	symbols := make([]string, 0, len(m.symbols))
	for symbol := range m.symbols {
		symbols = append(symbols, symbol)
	}
	m.mu.RUnlock()

	for _, symbol := range symbols {
		m.client.UnsubscribeTicker(symbol, productType)
	}

	m.logger.Info().Msgf("Unsubscribed from %d symbols", len(symbols))
}

func main() {
	// Create logger with detailed configuration
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

	// Create and start market data processor
	processor := NewMarketDataProcessor(logger)

	if err := processor.Start(); err != nil {
		log.Fatal("Failed to start processor:", err)
	}

	// Define symbols to monitor
	symbols := []string{
		"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT",
		"DOTUSDT", "LINKUSDT", "AVAXUSDT", "MATICUSDT",
	}

	// Subscribe to symbols
	processor.SubscribeToSymbols(symbols)

	// Start market summary display
	go displayMarketSummary(processor, logger)

	// Set up graceful shutdown
	setupAdvancedGracefulShutdown(processor, logger)
}

// displayMarketSummary displays market summary periodically
func displayMarketSummary(processor *MarketDataProcessor, logger zerolog.Logger) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		summary := processor.GetMarketSummary()

		fmt.Println("\n" + strings.Repeat("=", 100))
		fmt.Println("📊 ADVANCED MARKET DATA SUMMARY")
		fmt.Println(strings.Repeat("=", 100))
		fmt.Printf("%-12s %-15s %-15s %-15s %-20s\n",
			"SYMBOL", "PRICE", "VOLUME 24H", "CHANGE 24H", "LAST UPDATE")
		fmt.Println(strings.Repeat("-", 100))

		for _, ticker := range summary {
			timeStr := "Never"
			if !ticker.Timestamp.IsZero() {
				timeStr = ticker.Timestamp.Format("15:04:05")
			}

			fmt.Printf("%-12s %-15s %-15s %-15s %-20s\n",
				ticker.Symbol,
				ticker.LastPrice,
				ticker.Volume24h,
				ticker.Change24h,
				timeStr,
			)
		}
		fmt.Println(strings.Repeat("=", 100))
	}
}

// setupAdvancedGracefulShutdown sets up advanced graceful shutdown
func setupAdvancedGracefulShutdown(processor *MarketDataProcessor, logger zerolog.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Info().Msg("🚀 Advanced market data processor running...")
	logger.Info().Msg("💡 Features: Error handling, auto-reconnection, performance monitoring")
	logger.Info().Msg("📊 Market data being processed in real-time")
	logger.Info().Msg("🎧 Press Ctrl+C to stop")

	// Wait for shutdown signal
	<-sigChan

	logger.Info().Msg("🛑 Shutdown signal received...")

	// Create a timeout context for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Stop processor gracefully
	done := make(chan struct{})
	go func() {
		processor.Stop()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		logger.Info().Msg("✅ Graceful shutdown completed")
	case <-shutdownCtx.Done():
		logger.Warn().Msg("⚠️ Shutdown timeout reached, forcing exit")
	}
}
