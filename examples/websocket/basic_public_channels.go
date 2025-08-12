// Example: Basic Public WebSocket Channels
// This example demonstrates how to connect to Bitget's public WebSocket channels
// and subscribe to various market data streams.

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
)

func main() {
	// Create a logger for WebSocket debugging
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create WebSocket client for public channels (no authentication required)
	client := ws.NewBitgetBaseWsClient(
		logger,
		"wss://ws.bitget.com/v2/ws/public", // Public endpoint
		"",                                 // No secret key needed for public channels
	)

	// Set message handlers
	client.SetListener(defaultMessageHandlerForPublic, errorMessageHandlerForPublic)

	// Connect to WebSocket
	fmt.Println("🔌 Connecting to Bitget WebSocket...")
	client.Connect()
	client.ConnectWebSocket()

	// Start reading messages in a separate goroutine
	client.StartReadLoop()

	// Wait a moment for connection to establish
	time.Sleep(2 * time.Second)

	if !client.IsConnected() {
		log.Fatal("❌ Failed to connect to WebSocket")
	}

	fmt.Println("✅ Connected to Bitget WebSocket!")

	// Subscribe to different types of market data
	subscribeToMarketData(client)

	// Set up graceful shutdown
	setupGracefulShutdownForPublic(client)
}

func subscribeToMarketData(client *ws.BaseWsClient) {
	symbol := "BTCUSDT"
	productType := "USDT-FUTURES"

	fmt.Printf("📈 Subscribing to market data for %s...\n", symbol)

	// 1. Subscribe to ticker updates (24hr statistics)
	client.SubscribeTicker(symbol, productType, func(message string) {
		fmt.Printf("📊 TICKER: %s\n", message)
	})

	// 2. Subscribe to 1-minute candlesticks
	client.SubscribeCandles(symbol, productType, ws.Timeframe1m, func(message string) {
		fmt.Printf("🕯️  CANDLE 1m: %s\n", message)
	})

	// 3. Subscribe to top 5 order book levels
	client.SubscribeOrderBook5(symbol, productType, func(message string) {
		fmt.Printf("📚 ORDER BOOK (Top 5): %s\n", message)
	})

	// 4. Subscribe to trade executions
	client.SubscribeTrades(symbol, productType, func(message string) {
		fmt.Printf("💰 TRADE: %s\n", message)
	})

	// 5. Subscribe to mark price updates
	client.SubscribeMarkPrice(symbol, productType, func(message string) {
		fmt.Printf("🎯 MARK PRICE: %s\n", message)
	})

	// 6. Subscribe to funding rate information
	client.SubscribeFundingTime(symbol, productType, func(message string) {
		fmt.Printf("💸 FUNDING: %s\n", message)
	})

	fmt.Printf("✅ Subscribed to %d channels\n", client.GetSubscriptionCount())
}

func defaultMessageHandlerForPublic(message string) {
	// This handler receives all messages that don't have specific handlers
	fmt.Printf("📝 DEFAULT: %s\n", message)
}

func errorMessageHandlerForPublic(message string) {
	// This handler receives error messages
	fmt.Printf("❌ ERROR: %s\n", message)
}

func setupGracefulShutdownForPublic(client *ws.BaseWsClient) {
	// Create a channel to receive OS signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🎧 Listening for market data... Press Ctrl+C to stop.")

	// Wait for signal
	<-sigChan

	fmt.Println("\n🛑 Shutting down...")

	// Unsubscribe from all channels
	fmt.Println("📤 Unsubscribing from channels...")
	client.UnsubscribeTicker("BTCUSDT", "USDT-FUTURES")
	client.UnsubscribeCandles("BTCUSDT", "USDT-FUTURES", ws.Timeframe1m)
	client.UnsubscribeOrderBook5("BTCUSDT", "USDT-FUTURES")
	client.UnsubscribeTrades("BTCUSDT", "USDT-FUTURES")
	client.UnsubscribeMarkPrice("BTCUSDT", "USDT-FUTURES")
	client.UnsubscribeFundingTime("BTCUSDT", "USDT-FUTURES")

	// Close WebSocket connection
	client.Close()

	fmt.Println("✅ Graceful shutdown complete")
}
