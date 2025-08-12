// Example: Mixed Public and Private WebSocket Channels
// This example demonstrates how to use both public and private WebSocket channels
// simultaneously to get comprehensive market data and account updates.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

// TradingSession manages both public and private WebSocket connections
type TradingSession struct {
	// WebSocket clients
	publicClient  *ws.BaseWsClient
	privateClient *ws.BaseWsClient

	// Data storage
	marketData   map[string]MarketInfo
	accountData  AccountInfo
	activeOrders map[string]OrderInfo

	// Synchronization
	mu sync.RWMutex

	// Configuration
	symbols     []string
	productType string

	// Status tracking
	publicConnected  bool
	privateConnected bool
	authenticated    bool

	logger zerolog.Logger
}

// MarketInfo stores market data for a symbol
type MarketInfo struct {
	Symbol     string
	LastPrice  string
	Volume24h  string
	Change24h  string
	BestBid    string
	BestAsk    string
	LastTrade  string
	MarkPrice  string
	LastUpdate time.Time
}

// AccountInfo stores account information
type AccountInfo struct {
	TotalBalance     string
	AvailableBalance string
	UnrealizedPnl    string
	MarginUsed       string
	MarginRatio      string
	LastUpdate       time.Time
}

// OrderInfo stores order information
type OrderInfo struct {
	OrderId    string
	Symbol     string
	Side       string
	Size       string
	Price      string
	Status     string
	OrderType  string
	LastUpdate time.Time
}

// NewTradingSession creates a new trading session
func NewTradingSession(symbols []string, logger zerolog.Logger) *TradingSession {
	return &TradingSession{
		marketData:   make(map[string]MarketInfo),
		activeOrders: make(map[string]OrderInfo),
		symbols:      symbols,
		productType:  "USDT-FUTURES",
		logger:       logger,
	}
}

// Start initializes and starts both public and private connections
func (ts *TradingSession) Start(apiKey, secretKey, passphrase string) error {
	ts.logger.Info().Msg("🚀 Starting mixed trading session...")

	// Start public connection
	if err := ts.startPublicConnection(); err != nil {
		return fmt.Errorf("failed to start public connection: %w", err)
	}

	// Start private connection
	if err := ts.startPrivateConnection(apiKey, secretKey, passphrase); err != nil {
		return fmt.Errorf("failed to start private connection: %w", err)
	}

	// Subscribe to channels
	ts.subscribeToPublicChannels()
	ts.subscribeToPrivateChannels()

	// Start monitoring
	go ts.monitorConnections()
	go ts.displaySessionInfo()

	ts.logger.Info().Msg("✅ Trading session started successfully")
	return nil
}

// startPublicConnection initializes the public WebSocket connection
func (ts *TradingSession) startPublicConnection() error {
	ts.logger.Info().Msg("🔌 Connecting to public WebSocket...")

	ts.publicClient = ws.NewBitgetBaseWsClient(
		ts.logger,
		"wss://ws.bitget.com/v2/ws/public",
		"",
	)

	ts.publicClient.SetListener(ts.publicMessageHandler, ts.publicErrorHandler)

	ts.publicClient.Connect()
	ts.publicClient.ConnectWebSocket()
	ts.publicClient.StartReadLoop()

	// Wait for connection
	time.Sleep(2 * time.Second)

	if !ts.publicClient.IsConnected() {
		return fmt.Errorf("failed to connect to public WebSocket")
	}

	ts.mu.Lock()
	ts.publicConnected = true
	ts.mu.Unlock()

	ts.logger.Info().Msg("✅ Public WebSocket connected")
	return nil
}

// startPrivateConnection initializes the private WebSocket connection
func (ts *TradingSession) startPrivateConnection(apiKey, secretKey, passphrase string) error {
	ts.logger.Info().Msg("🔐 Connecting to private WebSocket...")

	ts.privateClient = ws.NewBitgetBaseWsClient(
		ts.logger,
		"wss://ws.bitget.com/v2/ws/private",
		secretKey,
	)

	ts.privateClient.SetListener(ts.privateMessageHandler, ts.privateErrorHandler)

	ts.privateClient.Connect()
	ts.privateClient.ConnectWebSocket()
	ts.privateClient.StartReadLoop()

	// Wait for connection
	time.Sleep(2 * time.Second)

	if !ts.privateClient.IsConnected() {
		return fmt.Errorf("failed to connect to private WebSocket")
	}

	ts.mu.Lock()
	ts.privateConnected = true
	ts.mu.Unlock()

	// Authenticate
	ts.logger.Info().Msg("🔑 Authenticating...")
	ts.privateClient.Login(apiKey, passphrase, common.SHA256)

	// Wait for authentication
	maxWait := 10 * time.Second
	checkInterval := 500 * time.Millisecond
	elapsed := time.Duration(0)

	for elapsed < maxWait {
		if ts.privateClient.IsLoggedIn() {
			ts.mu.Lock()
			ts.authenticated = true
			ts.mu.Unlock()
			ts.logger.Info().Msg("✅ Successfully authenticated")
			return nil
		}
		time.Sleep(checkInterval)
		elapsed += checkInterval
		fmt.Print(".")
	}

	return fmt.Errorf("authentication timeout")
}

// subscribeToPublicChannels subscribes to public market data channels
func (ts *TradingSession) subscribeToPublicChannels() {
	ts.logger.Info().Msg("📈 Subscribing to public channels...")

	for _, symbol := range ts.symbols {
		// Ticker for 24hr statistics
		ts.publicClient.SubscribeTicker(symbol, ts.productType, ts.createTickerHandler(symbol))

		// Order book for best bid/ask
		ts.publicClient.SubscribeOrderBook5(symbol, ts.productType, ts.createOrderBookHandler(symbol))

		// Recent trades
		ts.publicClient.SubscribeTrades(symbol, ts.productType, ts.createTradeHandler(symbol))

		// Mark price for PnL calculations
		ts.publicClient.SubscribeMarkPrice(symbol, ts.productType, ts.createMarkPriceHandler(symbol))

		ts.logger.Info().Msgf("📊 Subscribed to public data for %s", symbol)
	}

	ts.logger.Info().Msgf("✅ Subscribed to %d public channels", ts.publicClient.GetSubscriptionCount())
}

// subscribeToPrivateChannels subscribes to private account channels
func (ts *TradingSession) subscribeToPrivateChannels() {
	if !ts.authenticated {
		ts.logger.Warn().Msg("❌ Cannot subscribe to private channels: not authenticated")
		return
	}

	ts.logger.Info().Msg("🔒 Subscribing to private channels...")

	// Order updates
	ts.privateClient.SubscribeOrders(ts.productType, ts.createOrderUpdateHandler())

	// Fill updates
	ts.privateClient.SubscribeFills(ts.productType, ts.createFillUpdateHandler())

	// Position updates
	ts.privateClient.SubscribePositions(ts.productType, ts.createPositionUpdateHandler())

	// Account balance updates
	ts.privateClient.SubscribeAccount(ts.productType, ts.createAccountUpdateHandler())

	ts.logger.Info().Msgf("✅ Subscribed to %d private channels", ts.privateClient.GetSubscriptionCount())
}

// Create handler functions for public channels
func (ts *TradingSession) createTickerHandler(symbol string) ws.OnReceive {
	return func(message string) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		info := ts.marketData[symbol]
		info.Symbol = symbol
		info.LastPrice = "50000.00" // Parse from actual message
		info.Volume24h = "1000.0"   // Parse from actual message
		info.Change24h = "+2.5%"    // Parse from actual message
		info.LastUpdate = time.Now()
		ts.marketData[symbol] = info

		ts.logger.Debug().Msgf("📊 %s ticker update", symbol)
	}
}

func (ts *TradingSession) createOrderBookHandler(symbol string) ws.OnReceive {
	return func(message string) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		info := ts.marketData[symbol]
		info.BestBid = "49950.00" // Parse from actual message
		info.BestAsk = "50050.00" // Parse from actual message
		info.LastUpdate = time.Now()
		ts.marketData[symbol] = info

		ts.logger.Debug().Msgf("📚 %s order book update", symbol)
	}
}

func (ts *TradingSession) createTradeHandler(symbol string) ws.OnReceive {
	return func(message string) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		info := ts.marketData[symbol]
		info.LastTrade = "50000.00 @ 0.1" // Parse from actual message
		info.LastUpdate = time.Now()
		ts.marketData[symbol] = info

		ts.logger.Debug().Msgf("💰 %s trade update", symbol)
	}
}

func (ts *TradingSession) createMarkPriceHandler(symbol string) ws.OnReceive {
	return func(message string) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		info := ts.marketData[symbol]
		info.MarkPrice = "49980.00" // Parse from actual message
		info.LastUpdate = time.Now()
		ts.marketData[symbol] = info

		ts.logger.Debug().Msgf("🎯 %s mark price update", symbol)
	}
}

// Create handler functions for private channels
func (ts *TradingSession) createOrderUpdateHandler() ws.OnReceive {
	return func(message string) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		// Parse order update and store
		orderInfo := OrderInfo{
			OrderId:    "123456789", // Parse from actual message
			Symbol:     "BTCUSDT",   // Parse from actual message
			Side:       "buy",       // Parse from actual message
			Size:       "0.001",     // Parse from actual message
			Price:      "50000.00",  // Parse from actual message
			Status:     "filled",    // Parse from actual message
			OrderType:  "limit",     // Parse from actual message
			LastUpdate: time.Now(),
		}

		ts.activeOrders[orderInfo.OrderId] = orderInfo
		ts.logger.Info().Msgf("📋 Order update: %s %s %s @ %s",
			orderInfo.Side, orderInfo.Size, orderInfo.Symbol, orderInfo.Price)
	}
}

func (ts *TradingSession) createFillUpdateHandler() ws.OnReceive {
	return func(message string) {
		ts.logger.Info().Msgf("✅ Fill update: %s", message)
	}
}

func (ts *TradingSession) createPositionUpdateHandler() ws.OnReceive {
	return func(message string) {
		ts.logger.Info().Msgf("📊 Position update: %s", message)
	}
}

func (ts *TradingSession) createAccountUpdateHandler() ws.OnReceive {
	return func(message string) {
		ts.mu.Lock()
		defer ts.mu.Unlock()

		ts.accountData = AccountInfo{
			TotalBalance:     "10000.00", // Parse from actual message
			AvailableBalance: "8500.00",  // Parse from actual message
			UnrealizedPnl:    "+150.00",  // Parse from actual message
			MarginUsed:       "1500.00",  // Parse from actual message
			MarginRatio:      "15.0%",    // Parse from actual message
			LastUpdate:       time.Now(),
		}

		ts.logger.Info().Msg("💰 Account balance updated")
	}
}

// Message and error handlers
func (ts *TradingSession) publicMessageHandler(message string) {
	ts.logger.Debug().Msgf("📝 Public: %s", message)
}

func (ts *TradingSession) publicErrorHandler(message string) {
	ts.logger.Error().Msgf("❌ Public error: %s", message)
}

func (ts *TradingSession) privateMessageHandler(message string) {
	ts.logger.Debug().Msgf("🔒 Private: %s", message)
}

func (ts *TradingSession) privateErrorHandler(message string) {
	ts.logger.Error().Msgf("❌ Private error: %s", message)
}

// monitorConnections monitors the health of both connections
func (ts *TradingSession) monitorConnections() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ts.mu.RLock()
		publicOk := ts.publicConnected && ts.publicClient.IsConnected()
		privateOk := ts.privateConnected && ts.privateClient.IsConnected() && ts.authenticated
		ts.mu.RUnlock()

		status := "🟢"
		if !publicOk || !privateOk {
			status = "🔴"
		}

		ts.logger.Info().Msgf("%s Connection Status - Public: %v, Private: %v, Auth: %v",
			status, publicOk, privateOk, ts.authenticated)
	}
}

// displaySessionInfo displays comprehensive session information
func (ts *TradingSession) displaySessionInfo() {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		ts.displayDashboard()
	}
}

// displayDashboard shows a comprehensive dashboard
func (ts *TradingSession) displayDashboard() {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	fmt.Println("\n" + strings.Repeat("=", 120))
	fmt.Println("🚀 MIXED TRADING SESSION DASHBOARD")
	fmt.Println(strings.Repeat("=", 120))

	// Connection Status
	fmt.Printf("📡 CONNECTION STATUS\n")
	fmt.Printf("   Public WebSocket:  %s\n", formatStatus(ts.publicConnected))
	fmt.Printf("   Private WebSocket: %s\n", formatStatus(ts.privateConnected))
	fmt.Printf("   Authentication:    %s\n", formatStatus(ts.authenticated))
	fmt.Printf("   Public Subs:       %d channels\n", ts.getPublicSubCount())
	fmt.Printf("   Private Subs:      %d channels\n", ts.getPrivateSubCount())

	fmt.Println()

	// Account Information
	fmt.Printf("💰 ACCOUNT INFORMATION\n")
	fmt.Printf("   Total Balance:     %s USDT\n", ts.accountData.TotalBalance)
	fmt.Printf("   Available:         %s USDT\n", ts.accountData.AvailableBalance)
	fmt.Printf("   Unrealized PnL:    %s USDT\n", ts.accountData.UnrealizedPnl)
	fmt.Printf("   Margin Used:       %s USDT\n", ts.accountData.MarginUsed)
	fmt.Printf("   Margin Ratio:      %s\n", ts.accountData.MarginRatio)

	fmt.Println()

	// Market Data
	fmt.Printf("📈 MARKET DATA\n")
	fmt.Printf("%-12s %-12s %-12s %-12s %-12s %-12s %-15s\n",
		"SYMBOL", "PRICE", "CHANGE", "BID", "ASK", "MARK", "LAST UPDATE")
	fmt.Println(strings.Repeat("-", 120))

	for _, symbol := range ts.symbols {
		if info, exists := ts.marketData[symbol]; exists {
			timeStr := "Never"
			if !info.LastUpdate.IsZero() {
				timeStr = info.LastUpdate.Format("15:04:05")
			}

			fmt.Printf("%-12s %-12s %-12s %-12s %-12s %-12s %-15s\n",
				symbol,
				info.LastPrice,
				info.Change24h,
				info.BestBid,
				info.BestAsk,
				info.MarkPrice,
				timeStr,
			)
		}
	}

	fmt.Println()

	// Active Orders
	if len(ts.activeOrders) > 0 {
		fmt.Printf("📋 ACTIVE ORDERS (%d)\n", len(ts.activeOrders))
		fmt.Printf("%-15s %-12s %-8s %-12s %-12s %-12s\n",
			"ORDER ID", "SYMBOL", "SIDE", "SIZE", "PRICE", "STATUS")
		fmt.Println(strings.Repeat("-", 120))

		for _, order := range ts.activeOrders {
			fmt.Printf("%-15s %-12s %-8s %-12s %-12s %-12s\n",
				order.OrderId,
				order.Symbol,
				order.Side,
				order.Size,
				order.Price,
				order.Status,
			)
		}
	} else {
		fmt.Printf("📋 ACTIVE ORDERS: None\n")
	}

	fmt.Println(strings.Repeat("=", 120))
}

// Helper functions
func formatStatus(status bool) string {
	if status {
		return "✅ Connected"
	}
	return "❌ Disconnected"
}

func (ts *TradingSession) getPublicSubCount() int {
	if ts.publicClient != nil {
		return ts.publicClient.GetSubscriptionCount()
	}
	return 0
}

func (ts *TradingSession) getPrivateSubCount() int {
	if ts.privateClient != nil {
		return ts.privateClient.GetSubscriptionCount()
	}
	return 0
}

// Stop gracefully stops the trading session
func (ts *TradingSession) Stop() {
	ts.logger.Info().Msg("🛑 Stopping trading session...")

	// Unsubscribe from all channels
	if ts.publicClient != nil {
		for _, symbol := range ts.symbols {
			ts.publicClient.UnsubscribeTicker(symbol, ts.productType)
			ts.publicClient.UnsubscribeOrderBook5(symbol, ts.productType)
			ts.publicClient.UnsubscribeTrades(symbol, ts.productType)
			ts.publicClient.UnsubscribeMarkPrice(symbol, ts.productType)
		}
		ts.publicClient.Close()
	}

	if ts.privateClient != nil && ts.authenticated {
		ts.privateClient.UnsubscribeOrders(ts.productType)
		ts.privateClient.UnsubscribeFills(ts.productType)
		ts.privateClient.UnsubscribePositions(ts.productType)
		ts.privateClient.UnsubscribeAccount(ts.productType)
		ts.privateClient.Close()
	}

	ts.logger.Info().Msg("✅ Trading session stopped")
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: No .env file found, using environment variables")
	}

	// Get API credentials
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("❌ Missing API credentials. Please set BITGET_API_KEY, BITGET_SECRET_KEY, and BITGET_PASSPHRASE")
	}

	// Create logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Define symbols to track
	symbols := []string{
		"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT",
	}

	// Create trading session
	session := NewTradingSession(symbols, logger)

	// Start session
	if err := session.Start(apiKey, secretKey, passphrase); err != nil {
		log.Fatal("Failed to start trading session:", err)
	}

	// Set up graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n🎧 Mixed trading session active...")
	fmt.Println("📊 Monitoring both market data and account updates")
	fmt.Println("💡 This demonstrates combining public and private WebSocket channels")
	fmt.Println("🎯 Press Ctrl+C to stop")

	// Wait for shutdown signal
	<-sigChan

	fmt.Println("\n🛑 Shutdown signal received...")
	session.Stop()
	fmt.Println("✅ Graceful shutdown complete")
}
