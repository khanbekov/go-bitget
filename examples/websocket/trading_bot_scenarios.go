// Example: Practical Trading Bot Scenarios with WebSocket Integration
// This example demonstrates real-world trading scenarios combining REST API calls
// with WebSocket subscriptions for live monitoring and automated responses.

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

// TradingScenario represents different trading bot scenarios
type TradingScenario int

const (
	ScenarioGridTrading TradingScenario = iota
	ScenarioDCABot
	ScenarioScalpingBot
	ScenarioArbitrageBot
	ScenarioRiskManager
)

// MarketData holds real-time market information
type MarketData struct {
	Symbol        string    `json:"symbol"`
	Price         float64   `json:"price"`
	Volume24h     float64   `json:"volume24h"`
	Change24h     float64   `json:"change24h"`
	MarkPrice     float64   `json:"markPrice"`
	FundingRate   float64   `json:"fundingRate"`
	LastUpdate    time.Time `json:"lastUpdate"`
	BestBid       float64   `json:"bestBid"`
	BestAsk       float64   `json:"bestAsk"`
	Spread        float64   `json:"spread"`
}

// Position represents a trading position
type Position struct {
	Symbol       string    `json:"symbol"`
	Size         float64   `json:"size"`
	AvgPrice     float64   `json:"avgPrice"`
	UnrealizedPL float64   `json:"unrealizedPL"`
	Side         string    `json:"side"`
	LastUpdate   time.Time `json:"lastUpdate"`
}

// Order represents an active order
type Order struct {
	OrderID      string    `json:"orderId"`
	Symbol       string    `json:"symbol"`
	Side         string    `json:"side"`
	Size         float64   `json:"size"`
	Price        float64   `json:"price"`
	Status       string    `json:"status"`
	Type         string    `json:"type"`
	CreatedTime  time.Time `json:"createdTime"`
}

// TradingBotManager manages multiple trading scenarios
type TradingBotManager struct {
	client         *futures.Client
	wsManager      *futures.WebSocketManager
	logger         zerolog.Logger
	ctx            context.Context
	cancel         context.CancelFunc
	
	// Data storage
	marketData     map[string]*MarketData
	positions      map[string]*Position
	activeOrders   map[string]*Order
	
	// Synchronization
	marketMutex    sync.RWMutex
	positionMutex  sync.RWMutex
	orderMutex     sync.RWMutex
	
	// Trading parameters
	symbols        []string
	maxPositions   int
	riskLimit      float64
	
	// Statistics
	tradesExecuted int64
	profitLoss     float64
	startTime      time.Time
}

// NewTradingBotManager creates a new trading bot manager
func NewTradingBotManager(apiKey, secretKey, passphrase string, logger zerolog.Logger) *TradingBotManager {
	ctx, cancel := context.WithCancel(context.Background())
	
	client := futures.NewClient(apiKey, secretKey, passphrase)
	wsManager := client.NewWebSocketManager()
	
	return &TradingBotManager{
		client:       client,
		wsManager:    wsManager,
		logger:       logger,
		ctx:          ctx,
		cancel:       cancel,
		marketData:   make(map[string]*MarketData),
		positions:    make(map[string]*Position),
		activeOrders: make(map[string]*Order),
		symbols:      []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT"},
		maxPositions: 4,
		riskLimit:    1000.0,
		startTime:    time.Now(),
	}
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found")
	}

	// Get API credentials
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("Missing API credentials. Set BITGET_API_KEY, BITGET_SECRET_KEY, and BITGET_PASSPHRASE")
	}

	// Create logger
	logger := zerolog.New(os.Stderr).
		With().
		Timestamp().
		Caller().
		Logger().
		Level(zerolog.InfoLevel)

	// Create trading bot manager
	botManager := NewTradingBotManager(apiKey, secretKey, passphrase, logger)

	// Initialize and start trading scenarios
	if err := botManager.Initialize(); err != nil {
		log.Fatal("Failed to initialize trading bot:", err)
	}

	// Run all trading scenarios
	botManager.RunAllScenarios()

	// Set up graceful shutdown
	botManager.SetupGracefulShutdown()
}

// Initialize sets up WebSocket connections and initial data
func (tbm *TradingBotManager) Initialize() error {
	tbm.logger.Info().Msg("🤖 Initializing Trading Bot Manager")

	// Connect to public WebSocket for market data
	if err := tbm.wsManager.ConnectPublic(); err != nil {
		return fmt.Errorf("failed to connect to public WebSocket: %w", err)
	}

	// Set up market data subscriptions
	tbm.setupMarketDataSubscriptions()

	// Try to connect to private WebSocket (may fail with demo accounts)
	if err := tbm.setupPrivateSubscriptions(); err != nil {
		tbm.logger.Warn().Err(err).Msg("Private WebSocket failed - continuing with market data only")
	}

	// Initialize market data
	tbm.initializeMarketData()

	// Start background monitoring
	go tbm.monitorMarketData()
	go tbm.displayTradingDashboard()

	tbm.logger.Info().Msg("✅ Trading Bot Manager initialized successfully")
	return nil
}

// setupMarketDataSubscriptions sets up subscriptions for market data
func (tbm *TradingBotManager) setupMarketDataSubscriptions() {
	for _, symbol := range tbm.symbols {
		// Subscribe to ticker for price updates
		tbm.wsManager.SubscribeToTicker(symbol, tbm.createTickerHandler(symbol))
		
		// Subscribe to order book for spread calculation
		tbm.wsManager.SubscribeToOrderBook(symbol, 5, tbm.createOrderBookHandler(symbol))
		
		// Subscribe to mark price
		tbm.wsManager.SubscribeToMarkPrice(symbol, tbm.createMarkPriceHandler(symbol))
		
		// Subscribe to funding rate
		tbm.wsManager.SubscribeToFunding(symbol, tbm.createFundingHandler(symbol))
		
		time.Sleep(200 * time.Millisecond)
	}
	
	tbm.logger.Info().
		Int("symbols", len(tbm.symbols)).
		Int("subscriptions", tbm.wsManager.GetSubscriptionCount()).
		Msg("📡 Market data subscriptions setup complete")
}

// setupPrivateSubscriptions sets up private channel subscriptions
func (tbm *TradingBotManager) setupPrivateSubscriptions() error {
	// This will fail with demo accounts - that's expected
	config := futures.TradingStreamConfig{
		EnableOrders:    true,
		EnableFills:     true,
		EnablePositions: true,
		EnableAccount:   true,
		
		OrderHandler:    tbm.createOrderUpdateHandler(),
		FillHandler:     tbm.createFillUpdateHandler(),
		PositionHandler: tbm.createPositionUpdateHandler(),
		AccountHandler:  tbm.createAccountUpdateHandler(),
	}
	
	return tbm.wsManager.CreateTradingStream(tbm.ctx, os.Getenv("BITGET_API_KEY"), os.Getenv("BITGET_PASSPHRASE"), config)
}

// Message Handlers

func (tbm *TradingBotManager) createTickerHandler(symbol string) ws.OnReceive {
	return func(message string) {
		// Parse ticker data (simplified - would need actual JSON parsing)
		tbm.updateMarketPrice(symbol, tbm.extractPriceFromMessage(message))
	}
}

func (tbm *TradingBotManager) createOrderBookHandler(symbol string) ws.OnReceive {
	return func(message string) {
		// Parse order book and update spread
		bid, ask := tbm.extractBidAskFromMessage(message)
		tbm.updateBidAsk(symbol, bid, ask)
	}
}

func (tbm *TradingBotManager) createMarkPriceHandler(symbol string) ws.OnReceive {
	return func(message string) {
		markPrice := tbm.extractMarkPriceFromMessage(message)
		tbm.updateMarkPrice(symbol, markPrice)
	}
}

func (tbm *TradingBotManager) createFundingHandler(symbol string) ws.OnReceive {
	return func(message string) {
		fundingRate := tbm.extractFundingRateFromMessage(message)
		tbm.updateFundingRate(symbol, fundingRate)
	}
}

func (tbm *TradingBotManager) createOrderUpdateHandler() ws.OnReceive {
	return func(message string) {
		tbm.logger.Info().Str("data", message).Msg("📋 Order Update")
		// Parse and update order status
		tbm.handleOrderUpdate(message)
	}
}

func (tbm *TradingBotManager) createFillUpdateHandler() ws.OnReceive {
	return func(message string) {
		tbm.logger.Info().Str("data", message).Msg("✅ Fill Update")
		// Parse and update fills
		tbm.handleFillUpdate(message)
	}
}

func (tbm *TradingBotManager) createPositionUpdateHandler() ws.OnReceive {
	return func(message string) {
		tbm.logger.Info().Str("data", message).Msg("📊 Position Update")
		// Parse and update positions
		tbm.handlePositionUpdate(message)
	}
}

func (tbm *TradingBotManager) createAccountUpdateHandler() ws.OnReceive {
	return func(message string) {
		tbm.logger.Info().Str("data", message).Msg("💰 Account Update")
		// Parse and update account balance
		tbm.handleAccountUpdate(message)
	}
}

// Data Management Methods

func (tbm *TradingBotManager) initializeMarketData() {
	for _, symbol := range tbm.symbols {
		tbm.marketMutex.Lock()
		tbm.marketData[symbol] = &MarketData{
			Symbol:     symbol,
			Price:      0,
			LastUpdate: time.Now(),
		}
		tbm.marketMutex.Unlock()
	}
}

func (tbm *TradingBotManager) updateMarketPrice(symbol string, price float64) {
	tbm.marketMutex.Lock()
	defer tbm.marketMutex.Unlock()
	
	if data, exists := tbm.marketData[symbol]; exists {
		data.Price = price
		data.LastUpdate = time.Now()
	}
}

func (tbm *TradingBotManager) updateBidAsk(symbol string, bid, ask float64) {
	tbm.marketMutex.Lock()
	defer tbm.marketMutex.Unlock()
	
	if data, exists := tbm.marketData[symbol]; exists {
		data.BestBid = bid
		data.BestAsk = ask
		data.Spread = ask - bid
		data.LastUpdate = time.Now()
	}
}

func (tbm *TradingBotManager) updateMarkPrice(symbol string, markPrice float64) {
	tbm.marketMutex.Lock()
	defer tbm.marketMutex.Unlock()
	
	if data, exists := tbm.marketData[symbol]; exists {
		data.MarkPrice = markPrice
		data.LastUpdate = time.Now()
	}
}

func (tbm *TradingBotManager) updateFundingRate(symbol string, fundingRate float64) {
	tbm.marketMutex.Lock()
	defer tbm.marketMutex.Unlock()
	
	if data, exists := tbm.marketData[symbol]; exists {
		data.FundingRate = fundingRate
		data.LastUpdate = time.Now()
	}
}

// Message Parsing Methods (simplified for demo)

func (tbm *TradingBotManager) extractPriceFromMessage(message string) float64 {
	// In reality, this would parse JSON and extract the actual price
	// For demo purposes, we'll simulate price movement
	return 50000 + rand.Float64()*1000 - 500
}

func (tbm *TradingBotManager) extractBidAskFromMessage(message string) (float64, float64) {
	// Simulate bid/ask spread
	price := tbm.extractPriceFromMessage(message)
	spread := 0.1 + rand.Float64()*0.9 // $0.1 to $1.0 spread
	return price - spread/2, price + spread/2
}

func (tbm *TradingBotManager) extractMarkPriceFromMessage(message string) float64 {
	// Mark price is usually very close to spot price
	return tbm.extractPriceFromMessage(message) + rand.Float64()*2 - 1
}

func (tbm *TradingBotManager) extractFundingRateFromMessage(message string) float64 {
	// Funding rates are typically small percentages
	return (rand.Float64() - 0.5) * 0.01 // -0.5% to +0.5%
}

// Trading Scenarios Implementation

func (tbm *TradingBotManager) RunAllScenarios() {
	tbm.logger.Info().Msg("🚀 Starting All Trading Scenarios")
	
	// Start different trading scenarios as goroutines
	go tbm.runGridTradingScenario("BTCUSDT")
	go tbm.runDCAScenario("ETHUSDT")
	go tbm.runScalpingScenario("ADAUSDT")
	go tbm.runRiskManagerScenario()
	
	// Wait a bit for scenarios to initialize
	time.Sleep(2 * time.Second)
	
	tbm.logger.Info().Msg("✅ All trading scenarios are running")
}

// Scenario 1: Grid Trading Bot
func (tbm *TradingBotManager) runGridTradingScenario(symbol string) {
	tbm.logger.Info().Str("symbol", symbol).Msg("📊 Starting Grid Trading Scenario")
	
	gridLevels := 10
	gridSpacing := 50.0 // $50 between grid levels
	basePrice := 50000.0
	
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			tbm.executeGridTradingLogic(symbol, basePrice, gridLevels, gridSpacing)
		case <-tbm.ctx.Done():
			return
		}
	}
}

func (tbm *TradingBotManager) executeGridTradingLogic(symbol string, basePrice float64, levels int, spacing float64) {
	currentPrice := tbm.getCurrentPrice(symbol)
	if currentPrice == 0 {
		return
	}
	
	tbm.logger.Info().
		Str("symbol", symbol).
		Float64("currentPrice", currentPrice).
		Float64("basePrice", basePrice).
		Msg("🔲 Grid Trading: Analyzing grid positions")
	
	// Simulate grid trading logic
	if currentPrice < basePrice - spacing {
		tbm.logger.Info().Msg("🔲 Grid Trading: Price below grid - would place buy orders")
		tbm.simulateOrder(symbol, "buy", currentPrice, 0.001)
	} else if currentPrice > basePrice + spacing {
		tbm.logger.Info().Msg("🔲 Grid Trading: Price above grid - would place sell orders")
		tbm.simulateOrder(symbol, "sell", currentPrice, 0.001)
	}
}

// Scenario 2: Dollar Cost Averaging (DCA) Bot
func (tbm *TradingBotManager) runDCAScenario(symbol string) {
	tbm.logger.Info().Str("symbol", symbol).Msg("💰 Starting DCA Scenario")
	
	dcaAmount := 100.0 // $100 per interval
	dcaInterval := 60 * time.Second
	
	ticker := time.NewTicker(dcaInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			tbm.executeDCALogic(symbol, dcaAmount)
		case <-tbm.ctx.Done():
			return
		}
	}
}

func (tbm *TradingBotManager) executeDCALogic(symbol string, amount float64) {
	currentPrice := tbm.getCurrentPrice(symbol)
	if currentPrice == 0 {
		return
	}
	
	quantity := amount / currentPrice
	
	tbm.logger.Info().
		Str("symbol", symbol).
		Float64("amount", amount).
		Float64("price", currentPrice).
		Float64("quantity", quantity).
		Msg("💰 DCA: Executing dollar cost averaging buy")
	
	tbm.simulateOrder(symbol, "buy", currentPrice, quantity)
}

// Scenario 3: Scalping Bot
func (tbm *TradingBotManager) runScalpingScenario(symbol string) {
	tbm.logger.Info().Str("symbol", symbol).Msg("⚡ Starting Scalping Scenario")
	
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			tbm.executeScalpingLogic(symbol)
		case <-tbm.ctx.Done():
			return
		}
	}
}

func (tbm *TradingBotManager) executeScalpingLogic(symbol string) {
	spread := tbm.getCurrentSpread(symbol)
	if spread == 0 {
		return
	}
	
	// Scalp if spread is wide enough
	minSpread := 0.5
	if spread > minSpread {
		tbm.logger.Info().
			Str("symbol", symbol).
			Float64("spread", spread).
			Msg("⚡ Scalping: Wide spread detected - would place both buy and sell orders")
		
		currentPrice := tbm.getCurrentPrice(symbol)
		tbm.simulateOrder(symbol, "buy", currentPrice-spread/4, 0.01)
		tbm.simulateOrder(symbol, "sell", currentPrice+spread/4, 0.01)
	}
}

// Scenario 4: Risk Manager
func (tbm *TradingBotManager) runRiskManagerScenario() {
	tbm.logger.Info().Msg("🛡️ Starting Risk Manager Scenario")
	
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			tbm.executeRiskManagement()
		case <-tbm.ctx.Done():
			return
		}
	}
}

func (tbm *TradingBotManager) executeRiskManagement() {
	totalExposure := tbm.calculateTotalExposure()
	
	if totalExposure > tbm.riskLimit {
		tbm.logger.Warn().
			Float64("exposure", totalExposure).
			Float64("limit", tbm.riskLimit).
			Msg("🛡️ Risk Manager: Exposure limit exceeded - would reduce positions")
		
		tbm.simulateRiskReduction()
	} else {
		tbm.logger.Debug().
			Float64("exposure", totalExposure).
			Float64("limit", tbm.riskLimit).
			Msg("🛡️ Risk Manager: Risk levels normal")
	}
}

// Utility Methods

func (tbm *TradingBotManager) getCurrentPrice(symbol string) float64 {
	tbm.marketMutex.RLock()
	defer tbm.marketMutex.RUnlock()
	
	if data, exists := tbm.marketData[symbol]; exists {
		return data.Price
	}
	return 0
}

func (tbm *TradingBotManager) getCurrentSpread(symbol string) float64 {
	tbm.marketMutex.RLock()
	defer tbm.marketMutex.RUnlock()
	
	if data, exists := tbm.marketData[symbol]; exists {
		return data.Spread
	}
	return 0
}

func (tbm *TradingBotManager) calculateTotalExposure() float64 {
	tbm.positionMutex.RLock()
	defer tbm.positionMutex.RUnlock()
	
	total := 0.0
	for _, position := range tbm.positions {
		total += position.Size * position.AvgPrice
	}
	return total
}

func (tbm *TradingBotManager) simulateOrder(symbol, side string, price, quantity float64) {
	// In a real implementation, this would use the futures client to place orders
	orderID := fmt.Sprintf("%s_%s_%d", symbol, side, time.Now().Unix())
	
	order := &Order{
		OrderID:     orderID,
		Symbol:      symbol,
		Side:        side,
		Size:        quantity,
		Price:       price,
		Status:      "filled", // Simulate immediate fill for demo
		Type:        "limit",
		CreatedTime: time.Now(),
	}
	
	tbm.orderMutex.Lock()
	tbm.activeOrders[orderID] = order
	tbm.tradesExecuted++
	tbm.orderMutex.Unlock()
	
	tbm.logger.Info().
		Str("orderId", orderID).
		Str("symbol", symbol).
		Str("side", side).
		Float64("price", price).
		Float64("quantity", quantity).
		Msg("📝 SIMULATED ORDER PLACED")
	
	// Simulate position update
	go func() {
		time.Sleep(1 * time.Second)
		tbm.updateSimulatedPosition(symbol, side, quantity, price)
	}()
}

func (tbm *TradingBotManager) updateSimulatedPosition(symbol, side string, quantity, price float64) {
	tbm.positionMutex.Lock()
	defer tbm.positionMutex.Unlock()
	
	if position, exists := tbm.positions[symbol]; exists {
		// Update existing position
		if side == "buy" {
			position.Size += quantity
		} else {
			position.Size -= quantity
		}
		position.AvgPrice = price // Simplified - would normally calculate weighted average
		position.LastUpdate = time.Now()
	} else {
		// Create new position
		size := quantity
		if side == "sell" {
			size = -quantity
		}
		
		tbm.positions[symbol] = &Position{
			Symbol:     symbol,
			Size:       size,
			AvgPrice:   price,
			Side:       side,
			LastUpdate: time.Now(),
		}
	}
}

func (tbm *TradingBotManager) simulateRiskReduction() {
	tbm.logger.Info().Msg("🛡️ Risk Manager: Simulating position reduction")
	// In reality, would close positions or hedge
}

// Event Handlers

func (tbm *TradingBotManager) handleOrderUpdate(message string) {
	// Parse order update and update internal state
	tbm.logger.Debug().Str("message", message).Msg("Handling order update")
}

func (tbm *TradingBotManager) handleFillUpdate(message string) {
	// Parse fill update and update positions
	tbm.logger.Debug().Str("message", message).Msg("Handling fill update")
}

func (tbm *TradingBotManager) handlePositionUpdate(message string) {
	// Parse position update from WebSocket
	tbm.logger.Debug().Str("message", message).Msg("Handling position update")
}

func (tbm *TradingBotManager) handleAccountUpdate(message string) {
	// Parse account balance update
	tbm.logger.Debug().Str("message", message).Msg("Handling account update")
}

// Monitoring Methods

func (tbm *TradingBotManager) monitorMarketData() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			tbm.checkMarketDataFreshness()
		case <-tbm.ctx.Done():
			return
		}
	}
}

func (tbm *TradingBotManager) checkMarketDataFreshness() {
	tbm.marketMutex.RLock()
	defer tbm.marketMutex.RUnlock()
	
	staleThreshold := 30 * time.Second
	now := time.Now()
	
	for symbol, data := range tbm.marketData {
		if now.Sub(data.LastUpdate) > staleThreshold {
			tbm.logger.Warn().
				Str("symbol", symbol).
				Time("lastUpdate", data.LastUpdate).
				Msg("⚠️ Market data is stale")
		}
	}
}

// Dashboard Display

func (tbm *TradingBotManager) displayTradingDashboard() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	
	for {
		select {
		case <-ticker.C:
			tbm.printTradingDashboard()
		case <-tbm.ctx.Done():
			return
		}
	}
}

func (tbm *TradingBotManager) printTradingDashboard() {
	fmt.Println("\n" + strings.Repeat("=", 100))
	fmt.Println("🤖 TRADING BOT SCENARIOS DASHBOARD")
	fmt.Println(strings.Repeat("=", 100))
	
	// Market Data Section
	fmt.Println("📊 MARKET DATA")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-12s %-12s %-12s %-12s %-12s %-15s\n", 
		"SYMBOL", "PRICE", "BID", "ASK", "SPREAD", "LAST UPDATE")
	
	tbm.marketMutex.RLock()
	for _, symbol := range tbm.symbols {
		if data, exists := tbm.marketData[symbol]; exists {
			fmt.Printf("%-12s $%-11.2f $%-11.2f $%-11.2f $%-11.2f %-15s\n",
				symbol,
				data.Price,
				data.BestBid,
				data.BestAsk,
				data.Spread,
				data.LastUpdate.Format("15:04:05"),
			)
		}
	}
	tbm.marketMutex.RUnlock()
	
	// Positions Section
	fmt.Println("\n📊 SIMULATED POSITIONS")
	fmt.Println(strings.Repeat("-", 100))
	fmt.Printf("%-12s %-12s %-12s %-15s %-15s\n", 
		"SYMBOL", "SIZE", "AVG PRICE", "UNREALIZED P&L", "LAST UPDATE")
	
	tbm.positionMutex.RLock()
	totalPL := 0.0
	for symbol, position := range tbm.positions {
		currentPrice := tbm.getCurrentPrice(symbol)
		unrealizedPL := (currentPrice - position.AvgPrice) * position.Size
		totalPL += unrealizedPL
		
		fmt.Printf("%-12s %-12.4f $%-11.2f $%-14.2f %-15s\n",
			symbol,
			position.Size,
			position.AvgPrice,
			unrealizedPL,
			position.LastUpdate.Format("15:04:05"),
		)
	}
	tbm.positionMutex.RUnlock()
	
	// Statistics Section
	fmt.Println("\n📈 TRADING STATISTICS")
	fmt.Println(strings.Repeat("-", 100))
	runtime := time.Since(tbm.startTime)
	
	fmt.Printf("Runtime:           %v\n", runtime.Truncate(time.Second))
	fmt.Printf("Trades Executed:   %d\n", tbm.tradesExecuted)
	fmt.Printf("Active Orders:     %d\n", len(tbm.activeOrders))
	fmt.Printf("Total P&L:         $%.2f\n", totalPL)
	fmt.Printf("WebSocket Status:  %s\n", tbm.getWebSocketStatus())
	
	fmt.Println(strings.Repeat("=", 100))
}

func (tbm *TradingBotManager) getWebSocketStatus() string {
	if tbm.wsManager.IsConnected() {
		status := "✅ Public Connected"
		if tbm.wsManager.IsAuthenticated() {
			status += " & Private Authenticated"
		}
		return status
	}
	return "❌ Disconnected"
}

// Graceful Shutdown

func (tbm *TradingBotManager) SetupGracefulShutdown() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	tbm.logger.Info().Msg("🤖 Trading Bot Scenarios running...")
	tbm.logger.Info().Msg("💡 Scenarios: Grid Trading, DCA, Scalping, Risk Management")
	tbm.logger.Info().Msg("📊 Real-time market data with simulated trading")
	tbm.logger.Info().Msg("🎧 Press Ctrl+C to stop all scenarios")

	<-sigChan

	tbm.logger.Info().Msg("🛑 Shutting down all trading scenarios...")

	// Cancel context to stop all goroutines
	tbm.cancel()

	// Close WebSocket connections
	if err := tbm.wsManager.Close(); err != nil {
		tbm.logger.Error().Err(err).Msg("Error closing WebSocket")
	}

	// Display final statistics
	tbm.printFinalReport()

	tbm.logger.Info().Msg("✅ Trading bot shutdown complete")
}

func (tbm *TradingBotManager) printFinalReport() {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 FINAL TRADING REPORT")
	fmt.Println(strings.Repeat("=", 80))
	
	runtime := time.Since(tbm.startTime)
	fmt.Printf("Total Runtime:     %v\n", runtime.Truncate(time.Second))
	fmt.Printf("Total Trades:      %d\n", tbm.tradesExecuted)
	fmt.Printf("Final Positions:   %d\n", len(tbm.positions))
	
	// Calculate final P&L
	totalPL := 0.0
	tbm.positionMutex.RLock()
	for symbol, position := range tbm.positions {
		currentPrice := tbm.getCurrentPrice(symbol)
		unrealizedPL := (currentPrice - position.AvgPrice) * position.Size
		totalPL += unrealizedPL
	}
	tbm.positionMutex.RUnlock()
	
	fmt.Printf("Final P&L:         $%.2f\n", totalPL)
	
	if tbm.tradesExecuted > 0 {
		avgPLPerTrade := totalPL / float64(tbm.tradesExecuted)
		fmt.Printf("Avg P&L per Trade: $%.2f\n", avgPLPerTrade)
	}
	
	fmt.Println(strings.Repeat("=", 80))
}