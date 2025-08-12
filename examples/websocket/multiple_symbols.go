// Example: Multiple Symbols Monitoring
// This example shows how to monitor multiple trading pairs simultaneously
// using public WebSocket channels.

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

	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

// SymbolData holds market data for a specific symbol
type SymbolData struct {
	Symbol     string
	LastPrice  string
	Volume24h  string
	Change24h  string
	LastUpdate time.Time
	mutex      sync.RWMutex
}

// Update safely updates the symbol data
func (s *SymbolData) Update(price, volume, change string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.LastPrice = price
	s.Volume24h = volume
	s.Change24h = change
	s.LastUpdate = time.Now()
}

// GetData safely retrieves the symbol data
func (s *SymbolData) GetData() (string, string, string, time.Time) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.LastPrice, s.Volume24h, s.Change24h, s.LastUpdate
}

func main() {
	// Create logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create WebSocket client
	client := ws.NewBitgetBaseWsClient(
		logger,
		"wss://ws.bitget.com/v2/ws/public",
		"",
	)

	// Set up message handlers
	client.SetListener(defaultHandler, errorHandler)

	// Connect
	fmt.Println("🔌 Connecting to Bitget WebSocket...")
	client.Connect()
	client.ConnectWebSocket()
	client.StartReadLoop()

	time.Sleep(2 * time.Second)

	if !client.IsConnected() {
		log.Fatal("❌ Failed to connect to WebSocket")
	}

	fmt.Println("✅ Connected successfully!")

	// Define symbols to monitor
	symbols := []string{
		"BTCUSDT",
		"ETHUSDT",
		"ADAUSDT",
		"SOLUSDT",
		"DOTUSDT",
	}

	// Create data storage for each symbol
	symbolData := make(map[string]*SymbolData)
	for _, symbol := range symbols {
		symbolData[symbol] = &SymbolData{
			Symbol: symbol,
		}
	}

	// Subscribe to multiple symbols
	subscribeToMultipleSymbols(client, symbols, symbolData)

	// Start data display routine
	go displayDataPeriodically(symbolData)

	// Graceful shutdown
	setupGracefulShutdownForMultipleSymbols(client, symbols)
}

func subscribeToMultipleSymbols(client *ws.BaseWsClient, symbols []string, symbolData map[string]*SymbolData) {
	productType := "USDT-FUTURES"

	fmt.Printf("📈 Subscribing to %d symbols...\n", len(symbols))

	for _, symbol := range symbols {
		// Create symbol-specific handlers
		tickerHandler := createTickerHandler(symbol, symbolData[symbol])
		tradeHandler := createTradeHandler(symbol)
		candleHandler := createCandleHandler(symbol)

		// Subscribe to ticker for price updates
		client.SubscribeTicker(symbol, productType, tickerHandler)

		// Subscribe to trades for real-time activity
		client.SubscribeTrades(symbol, productType, tradeHandler)

		// Subscribe to 5-minute candles
		client.SubscribeCandles(symbol, productType, ws.Timeframe5m, candleHandler)

		fmt.Printf("✅ Subscribed to %s\n", symbol)
	}

	fmt.Printf("📊 Total subscriptions: %d\n", client.GetSubscriptionCount())
}

func createTickerHandler(symbol string, data *SymbolData) ws.OnReceive {
	return func(message string) {
		// Parse ticker message and update data
		// In real implementation, you would parse the JSON message
		fmt.Printf("📊 %s TICKER UPDATE\n", symbol)

		// For demo purposes, we'll simulate data update
		data.Update("50000", "1000", "+2.5%")
	}
}

func createTradeHandler(symbol string) ws.OnReceive {
	return func(message string) {
		fmt.Printf("💰 %s TRADE: %s\n", symbol, message)
	}
}

func createCandleHandler(symbol string) ws.OnReceive {
	return func(message string) {
		fmt.Printf("🕯️  %s CANDLE 5m: %s\n", symbol, message)
	}
}

func displayDataPeriodically(symbolData map[string]*SymbolData) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			displaySummary(symbolData)
		}
	}
}

func displaySummary(symbolData map[string]*SymbolData) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 MARKET DATA SUMMARY")
	fmt.Println(strings.Repeat("=", 80))
	fmt.Printf("%-10s %-12s %-15s %-10s %-20s\n", "SYMBOL", "PRICE", "VOLUME 24H", "CHANGE", "LAST UPDATE")
	fmt.Println(strings.Repeat("-", 80))

	for symbol, data := range symbolData {
		price, volume, change, lastUpdate := data.GetData()

		timeStr := "Never"
		if !lastUpdate.IsZero() {
			timeStr = lastUpdate.Format("15:04:05")
		}

		fmt.Printf("%-10s %-12s %-15s %-10s %-20s\n",
			symbol,
			price,
			volume,
			change,
			timeStr,
		)
	}
	fmt.Println(strings.Repeat("=", 80))
}

func defaultHandler(message string) {
	// Handle general messages
}

func errorHandler(message string) {
	fmt.Printf("❌ ERROR: %s\n", message)
}

func setupGracefulShutdownForMultipleSymbols(client *ws.BaseWsClient, symbols []string) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n🎧 Monitoring multiple symbols... Press Ctrl+C to stop.")

	<-sigChan

	fmt.Println("\n🛑 Shutting down...")

	// Unsubscribe from all symbols
	productType := "USDT-FUTURES"
	for _, symbol := range symbols {
		client.UnsubscribeTicker(symbol, productType)
		client.UnsubscribeTrades(symbol, productType)
		client.UnsubscribeCandles(symbol, productType, ws.Timeframe5m)
		fmt.Printf("📤 Unsubscribed from %s\n", symbol)
	}

	client.Close()
	fmt.Println("✅ Shutdown complete")
}
