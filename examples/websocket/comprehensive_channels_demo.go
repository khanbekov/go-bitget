// Example: Comprehensive WebSocket Channels Demo
// This example demonstrates ALL available WebSocket channels with detailed examples
// for each timeframe, product type, and subscription pattern.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

// MessageStats tracks statistics for different message types
type MessageStats struct {
	Ticker      int64 `json:"ticker"`
	Candle      int64 `json:"candle"`
	OrderBook   int64 `json:"orderbook"`
	Trade       int64 `json:"trade"`
	MarkPrice   int64 `json:"mark_price"`
	Funding     int64 `json:"funding"`
	Orders      int64 `json:"orders"`
	Fills       int64 `json:"fills"`
	Positions   int64 `json:"positions"`
	Account     int64 `json:"account"`
	PlanOrders  int64 `json:"plan_orders"`
}

// ComprehensiveDemo manages comprehensive WebSocket channel demonstrations
type ComprehensiveDemo struct {
	publicClient  *ws.BaseWsClient
	privateClient *ws.BaseWsClient
	logger        zerolog.Logger
	stats         MessageStats
	ctx           context.Context
	cancel        context.CancelFunc
}

// NewComprehensiveDemo creates a new comprehensive demo instance
func NewComprehensiveDemo(logger zerolog.Logger) *ComprehensiveDemo {
	ctx, cancel := context.WithCancel(context.Background())
	
	return &ComprehensiveDemo{
		logger: logger,
		ctx:    ctx,
		cancel: cancel,
	}
}

func main() {
	// Create enhanced logger
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Warn().Msg("No .env file found, using environment variables")
	}

	// Create demo instance
	demo := NewComprehensiveDemo(logger)

	// Run comprehensive demonstration
	if err := demo.RunComprehensiveDemo(); err != nil {
		log.Fatal("Comprehensive demo failed:", err)
	}

	// Set up graceful shutdown
	demo.setupGracefulShutdown()
}

// RunComprehensiveDemo runs the complete demonstration
func (d *ComprehensiveDemo) RunComprehensiveDemo() error {
	d.logger.Info().Msg("🚀 Starting Comprehensive WebSocket Channels Demo")
	d.logger.Info().Msg("📡 This demo covers ALL available channels with multiple examples")

	// Phase 1: Public Channels Demo
	if err := d.runPublicChannelsDemo(); err != nil {
		return fmt.Errorf("public channels demo failed: %w", err)
	}

	// Wait a bit before private channels
	time.Sleep(5 * time.Second)

	// Phase 2: Private Channels Demo (if credentials available)
	if err := d.runPrivateChannelsDemo(); err != nil {
		d.logger.Warn().Err(err).Msg("Private channels demo failed (this is expected without valid credentials)")
	}

	// Phase 3: Start monitoring and stats
	go d.monitorConnectionHealth()
	go d.displayLiveStats()

	return nil
}

// runPublicChannelsDemo demonstrates all public channels
func (d *ComprehensiveDemo) runPublicChannelsDemo() error {
	d.logger.Info().Msg("🌍 Phase 1: Public Channels Demonstration")

	// Create public WebSocket client
	d.publicClient = ws.NewBitgetBaseWsClient(
		d.logger,
		"wss://ws.bitget.com/v2/ws/public",
		"",
	)

	// Set handlers
	d.publicClient.SetListener(d.generalMessageHandler, d.errorHandler)

	// Connect
	d.publicClient.Connect()
	d.publicClient.ConnectWebSocket()
	d.publicClient.StartReadLoop()

	// Wait for connection
	time.Sleep(2 * time.Second)
	if !d.publicClient.IsConnected() {
		return fmt.Errorf("failed to connect to public WebSocket")
	}

	d.logger.Info().Msg("✅ Connected to public WebSocket")

	// Demonstrate all public channels
	d.demonstrateTickerChannels()
	d.demonstrateCandlestickChannels() 
	d.demonstrateOrderBookChannels()
	d.demonstrateTradeChannels()
	d.demonstrateMarkPriceChannels()
	d.demonstrateFundingChannels()

	return nil
}

// demonstrateTickerChannels shows ticker subscriptions for different product types
func (d *ComprehensiveDemo) demonstrateTickerChannels() {
	d.logger.Info().Msg("📊 Demonstrating Ticker Channels")

	// Test symbols for different product types
	examples := []struct {
		symbol      string
		productType string
		description string
	}{
		{"BTCUSDT", "USDT-FUTURES", "Bitcoin USDT Perpetual"},
		{"ETHUSDT", "USDT-FUTURES", "Ethereum USDT Perpetual"},
		{"BTCUSD", "COIN-FUTURES", "Bitcoin USD Perpetual"},
		{"ETHUSD", "COIN-FUTURES", "Ethereum USD Perpetual"},
	}

	for _, example := range examples {
		handler := d.createTickerHandler(example.symbol, example.description)
		d.publicClient.SubscribeTicker(example.symbol, example.productType, handler)
		
		d.logger.Info().
			Str("symbol", example.symbol).
			Str("product", example.productType).
			Str("description", example.description).
			Msg("📈 Subscribed to ticker")
		
		time.Sleep(500 * time.Millisecond) // Avoid overwhelming the server
	}
}

// demonstrateCandlestickChannels shows all timeframe subscriptions
func (d *ComprehensiveDemo) demonstrateCandlestickChannels() {
	d.logger.Info().Msg("🕯️ Demonstrating Candlestick Channels - ALL Timeframes")

	symbol := "BTCUSDT"
	productType := "USDT-FUTURES"

	// All supported timeframes
	timeframes := []struct {
		value       string
		description string
	}{
		{ws.Timeframe1m, "1 Minute"},
		{ws.Timeframe5m, "5 Minutes"},
		{ws.Timeframe15m, "15 Minutes"},
		{ws.Timeframe30m, "30 Minutes"},
		{ws.Timeframe1h, "1 Hour"},
		{ws.Timeframe4h, "4 Hours"},
		{ws.Timeframe6h, "6 Hours"},
		{ws.Timeframe12h, "12 Hours"},
		{ws.Timeframe1d, "1 Day"},
		{ws.Timeframe3d, "3 Days"},
		{ws.Timeframe1w, "1 Week"},
		{ws.Timeframe1M, "1 Month"},
	}

	for i, tf := range timeframes {
		// Subscribe to first 4 timeframes for demo purposes
		if i < 4 {
			handler := d.createCandleHandler(symbol, tf.value, tf.description)
			d.publicClient.SubscribeCandles(symbol, productType, tf.value, handler)
			
			d.logger.Info().
				Str("symbol", symbol).
				Str("timeframe", tf.value).
				Str("description", tf.description).
				Msg("🕯️ Subscribed to candlestick")
		} else {
			// Just log available timeframes
			d.logger.Info().
				Str("timeframe", tf.value).
				Str("description", tf.description).
				Msg("🕯️ Available timeframe")
		}
		
		time.Sleep(200 * time.Millisecond)
	}
}

// demonstrateOrderBookChannels shows different order book depth levels
func (d *ComprehensiveDemo) demonstrateOrderBookChannels() {
	d.logger.Info().Msg("📚 Demonstrating Order Book Channels - Multiple Depths")

	symbol := "ETHUSDT"
	productType := "USDT-FUTURES"

	// Different order book depths
	depths := []struct {
		method      string
		levels      string
		description string
	}{
		{"books5", "5", "Top 5 Bid/Ask Levels"},
		{"books15", "15", "Top 15 Bid/Ask Levels"},
		{"books", "Full", "Complete Order Book Depth"},
	}

	for _, depth := range depths {
		var handler ws.OnReceive = d.createOrderBookHandler(symbol, depth.levels, depth.description)

		switch depth.method {
		case "books5":
			d.publicClient.SubscribeOrderBook5(symbol, productType, handler)
		case "books15":
			d.publicClient.SubscribeOrderBook15(symbol, productType, handler)
		case "books":
			d.publicClient.SubscribeOrderBook(symbol, productType, handler)
		}

		d.logger.Info().
			Str("symbol", symbol).
			Str("method", depth.method).
			Str("levels", depth.levels).
			Str("description", depth.description).
			Msg("📚 Subscribed to order book")

		time.Sleep(300 * time.Millisecond)
	}
}

// demonstrateTradeChannels shows trade execution subscriptions
func (d *ComprehensiveDemo) demonstrateTradeChannels() {
	d.logger.Info().Msg("💰 Demonstrating Trade Execution Channels")

	symbols := []string{"ADAUSDT", "SOLUSDT", "DOGEUSDT"}
	productType := "USDT-FUTURES"

	for _, symbol := range symbols {
		handler := d.createTradeHandler(symbol)
		d.publicClient.SubscribeTrades(symbol, productType, handler)

		d.logger.Info().
			Str("symbol", symbol).
			Msg("💰 Subscribed to trade executions")

		time.Sleep(300 * time.Millisecond)
	}
}

// demonstrateMarkPriceChannels shows mark price subscriptions
func (d *ComprehensiveDemo) demonstrateMarkPriceChannels() {
	d.logger.Info().Msg("🎯 Demonstrating Mark Price Channels")

	symbols := []string{"BTCUSDT", "ETHUSDT"}
	productType := "USDT-FUTURES"

	for _, symbol := range symbols {
		handler := d.createMarkPriceHandler(symbol)
		d.publicClient.SubscribeMarkPrice(symbol, productType, handler)

		d.logger.Info().
			Str("symbol", symbol).
			Msg("🎯 Subscribed to mark price")

		time.Sleep(300 * time.Millisecond)
	}
}

// demonstrateFundingChannels shows funding rate subscriptions
func (d *ComprehensiveDemo) demonstrateFundingChannels() {
	d.logger.Info().Msg("💸 Demonstrating Funding Rate Channels")

	symbols := []string{"BTCUSDT", "ETHUSDT"}
	productType := "USDT-FUTURES"

	for _, symbol := range symbols {
		handler := d.createFundingHandler(symbol)
		d.publicClient.SubscribeFundingTime(symbol, productType, handler)

		d.logger.Info().
			Str("symbol", symbol).
			Msg("💸 Subscribed to funding rate")

		time.Sleep(300 * time.Millisecond)
	}
}

// runPrivateChannelsDemo demonstrates private channels (requires authentication)
func (d *ComprehensiveDemo) runPrivateChannelsDemo() error {
	d.logger.Info().Msg("🔒 Phase 2: Private Channels Demonstration")

	// Get credentials
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		return fmt.Errorf("missing API credentials for private channels")
	}

	// Create private WebSocket client
	d.privateClient = ws.NewBitgetBaseWsClient(
		d.logger,
		"wss://ws.bitget.com/v2/ws/private",
		secretKey,
	)

	// Set handlers
	d.privateClient.SetListener(d.generalMessageHandler, d.errorHandler)

	// Connect
	d.privateClient.Connect()
	d.privateClient.ConnectWebSocket()
	d.privateClient.StartReadLoop()

	// Wait for connection
	time.Sleep(2 * time.Second)
	if !d.privateClient.IsConnected() {
		return fmt.Errorf("failed to connect to private WebSocket")
	}

	// Authenticate
	d.logger.Info().Msg("🔐 Authenticating...")
	d.privateClient.Login(apiKey, passphrase, common.SHA256)

	// Wait for authentication
	time.Sleep(3 * time.Second)
	if !d.privateClient.IsLoggedIn() {
		return fmt.Errorf("authentication failed")
	}

	d.logger.Info().Msg("✅ Authenticated successfully")

	// Subscribe to private channels
	d.subscribeToPrivateChannels()

	return nil
}

// subscribeToPrivateChannels subscribes to all private channels
func (d *ComprehensiveDemo) subscribeToPrivateChannels() {
	productTypes := []string{"USDT-FUTURES", "COIN-FUTURES"}

	for _, productType := range productTypes {
		d.logger.Info().
			Str("productType", productType).
			Msg("🔒 Subscribing to private channels")

		// Order updates
		d.privateClient.SubscribeOrders(productType, d.createOrdersHandler(productType))

		// Fill updates  
		d.privateClient.SubscribeFills(productType, d.createFillsHandler(productType))

		// Position updates
		d.privateClient.SubscribePositions(productType, d.createPositionsHandler(productType))

		// Account updates
		d.privateClient.SubscribeAccount(productType, d.createAccountHandler(productType))

		// Plan order updates
		d.privateClient.SubscribePlanOrders(productType, d.createPlanOrdersHandler(productType))

		time.Sleep(500 * time.Millisecond)
	}
}

// Message Handler Creation Methods

func (d *ComprehensiveDemo) createTickerHandler(symbol, description string) ws.OnReceive {
	return func(message string) {
		d.stats.Ticker++
		d.logger.Debug().
			Str("symbol", symbol).
			Str("description", description).
			Str("data", message).
			Msg("📊 TICKER")
	}
}

func (d *ComprehensiveDemo) createCandleHandler(symbol, timeframe, description string) ws.OnReceive {
	return func(message string) {
		d.stats.Candle++
		d.logger.Debug().
			Str("symbol", symbol).
			Str("timeframe", timeframe).
			Str("description", description).
			Str("data", message).
			Msg("🕯️ CANDLE")
	}
}

func (d *ComprehensiveDemo) createOrderBookHandler(symbol, levels, description string) ws.OnReceive {
	return func(message string) {
		d.stats.OrderBook++
		d.logger.Debug().
			Str("symbol", symbol).
			Str("levels", levels).
			Str("description", description).
			Str("data", message).
			Msg("📚 ORDER_BOOK")
	}
}

func (d *ComprehensiveDemo) createTradeHandler(symbol string) ws.OnReceive {
	return func(message string) {
		d.stats.Trade++
		d.logger.Debug().
			Str("symbol", symbol).
			Str("data", message).
			Msg("💰 TRADE")
	}
}

func (d *ComprehensiveDemo) createMarkPriceHandler(symbol string) ws.OnReceive {
	return func(message string) {
		d.stats.MarkPrice++
		d.logger.Debug().
			Str("symbol", symbol).
			Str("data", message).
			Msg("🎯 MARK_PRICE")
	}
}

func (d *ComprehensiveDemo) createFundingHandler(symbol string) ws.OnReceive {
	return func(message string) {
		d.stats.Funding++
		d.logger.Debug().
			Str("symbol", symbol).
			Str("data", message).
			Msg("💸 FUNDING")
	}
}

func (d *ComprehensiveDemo) createOrdersHandler(productType string) ws.OnReceive {
	return func(message string) {
		d.stats.Orders++
		d.logger.Info().
			Str("productType", productType).
			Str("data", message).
			Msg("📋 ORDER_UPDATE")
	}
}

func (d *ComprehensiveDemo) createFillsHandler(productType string) ws.OnReceive {
	return func(message string) {
		d.stats.Fills++
		d.logger.Info().
			Str("productType", productType).
			Str("data", message).
			Msg("✅ FILL_UPDATE")
	}
}

func (d *ComprehensiveDemo) createPositionsHandler(productType string) ws.OnReceive {
	return func(message string) {
		d.stats.Positions++
		d.logger.Info().
			Str("productType", productType).
			Str("data", message).
			Msg("📊 POSITION_UPDATE")
	}
}

func (d *ComprehensiveDemo) createAccountHandler(productType string) ws.OnReceive {
	return func(message string) {
		d.stats.Account++
		d.logger.Info().
			Str("productType", productType).
			Str("data", message).
			Msg("💰 ACCOUNT_UPDATE")
	}
}

func (d *ComprehensiveDemo) createPlanOrdersHandler(productType string) ws.OnReceive {
	return func(message string) {
		d.stats.PlanOrders++
		d.logger.Info().
			Str("productType", productType).
			Str("data", message).
			Msg("⚡ PLAN_ORDER_UPDATE")
	}
}

// Utility Methods

func (d *ComprehensiveDemo) generalMessageHandler(message string) {
	d.logger.Debug().Str("data", message).Msg("ℹ️ GENERAL")
}

func (d *ComprehensiveDemo) errorHandler(message string) {
	d.logger.Error().Str("error", message).Msg("❌ ERROR")
}

// monitorConnectionHealth monitors both connections
func (d *ComprehensiveDemo) monitorConnectionHealth() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			publicStatus := "❌ Disconnected"
			if d.publicClient != nil && d.publicClient.IsConnected() {
				publicStatus = "✅ Connected"
			}

			privateStatus := "❌ Disconnected"
			if d.privateClient != nil && d.privateClient.IsConnected() {
				privateStatus = "✅ Connected"
				if d.privateClient.IsLoggedIn() {
					privateStatus += " & Authenticated"
				}
			}

			d.logger.Info().
				Str("public", publicStatus).
				Str("private", privateStatus).
				Msg("🔍 Connection Health Check")

		case <-d.ctx.Done():
			return
		}
	}
}

// displayLiveStats shows live statistics every 10 seconds
func (d *ComprehensiveDemo) displayLiveStats() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			d.displayCurrentStats()
		case <-d.ctx.Done():
			return
		}
	}
}

// displayCurrentStats displays current message statistics
func (d *ComprehensiveDemo) displayCurrentStats() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 COMPREHENSIVE WEBSOCKET DEMO - LIVE STATISTICS")
	fmt.Println(strings.Repeat("=", 80))
	
	fmt.Printf("📈 Ticker Messages:     %d\n", d.stats.Ticker)
	fmt.Printf("🕯️  Candle Messages:     %d\n", d.stats.Candle)
	fmt.Printf("📚 OrderBook Messages:  %d\n", d.stats.OrderBook)
	fmt.Printf("💰 Trade Messages:      %d\n", d.stats.Trade)
	fmt.Printf("🎯 MarkPrice Messages:  %d\n", d.stats.MarkPrice)
	fmt.Printf("💸 Funding Messages:    %d\n", d.stats.Funding)
	fmt.Printf("📋 Order Updates:       %d\n", d.stats.Orders)
	fmt.Printf("✅ Fill Updates:        %d\n", d.stats.Fills)
	fmt.Printf("📊 Position Updates:    %d\n", d.stats.Positions)
	fmt.Printf("💰 Account Updates:     %d\n", d.stats.Account)
	fmt.Printf("⚡ Plan Order Updates:  %d\n", d.stats.PlanOrders)
	
	totalPublic := d.publicClient.GetSubscriptionCount()
	totalPrivate := 0
	if d.privateClient != nil {
		totalPrivate = d.privateClient.GetSubscriptionCount()
	}
	
	fmt.Printf("\n🔌 Active Subscriptions: %d (Public: %d, Private: %d)\n", 
		totalPublic + totalPrivate, totalPublic, totalPrivate)
	
	// Calculate total messages
	total := d.stats.Ticker + d.stats.Candle + d.stats.OrderBook + d.stats.Trade +
		d.stats.MarkPrice + d.stats.Funding + d.stats.Orders + d.stats.Fills +
		d.stats.Positions + d.stats.Account + d.stats.PlanOrders
	
	fmt.Printf("📨 Total Messages:      %d\n", total)
	fmt.Println(strings.Repeat("=", 80))
}

// setupGracefulShutdown handles graceful shutdown
func (d *ComprehensiveDemo) setupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	d.logger.Info().Msg("🎧 Comprehensive demo running... Press Ctrl+C to stop")
	d.logger.Info().Msg("💡 All WebSocket channels are being demonstrated with live data")

	<-sigChan

	d.logger.Info().Msg("🛑 Shutting down comprehensive demo...")

	// Cancel context
	d.cancel()

	// Close connections
	if d.publicClient != nil {
		d.logger.Info().Msg("📤 Closing public WebSocket connection")
		d.publicClient.Close()
	}

	if d.privateClient != nil {
		d.logger.Info().Msg("📤 Closing private WebSocket connection")
		d.privateClient.Close()
	}

	// Display final statistics
	d.logger.Info().Msg("📊 Final Statistics:")
	statsJSON, _ := json.MarshalIndent(d.stats, "", "  ")
	fmt.Println(string(statsJSON))

	d.logger.Info().Msg("✅ Comprehensive demo shutdown complete")
}