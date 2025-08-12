// Example: Private WebSocket Channels with Authentication
// This example demonstrates how to connect to Bitget's private WebSocket channels
// to receive real-time account updates including orders, fills, positions, and balances.

package main

import (
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

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Warning: No .env file found, using environment variables")
	}

	// Get API credentials from environment
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("❌ Missing API credentials. Please set BITGET_API_KEY, BITGET_SECRET_KEY, and BITGET_PASSPHRASE environment variables")
	}

	// Create logger
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Create WebSocket client for private channels
	client := ws.NewBitgetBaseWsClient(
		logger,
		"wss://ws.bitget.com/v2/ws/private", // Private endpoint
		secretKey,                           // Secret key required for private channels
	)

	// Set message handlers
	client.SetListener(defaultMessageHandler, errorMessageHandler)

	// Connect to WebSocket
	fmt.Println("🔌 Connecting to Bitget Private WebSocket...")
	client.Connect()
	client.ConnectWebSocket()

	// Start reading messages
	client.StartReadLoop()

	// Wait for connection
	time.Sleep(2 * time.Second)

	if !client.IsConnected() {
		log.Fatal("❌ Failed to connect to WebSocket")
	}

	fmt.Println("✅ Connected to WebSocket!")

	// Authenticate (REQUIRED for private channels)
	fmt.Println("🔐 Authenticating...")
	client.Login(apiKey, passphrase, common.SHA256)

	// Wait for authentication to complete
	fmt.Println("⏳ Waiting for authentication...")
	waitForAuthentication(client)

	if !client.IsLoggedIn() {
		log.Fatal("❌ Authentication failed")
	}

	fmt.Println("✅ Successfully authenticated!")

	// Subscribe to private channels
	subscribeToPrivateChannels(client)

	// Set up graceful shutdown
	setupGracefulShutdownForPrivate(client)
}

func waitForAuthentication(client *ws.BaseWsClient) {
	maxWait := 10 * time.Second
	checkInterval := 500 * time.Millisecond
	elapsed := time.Duration(0)

	for elapsed < maxWait {
		if client.IsLoggedIn() {
			return
		}
		time.Sleep(checkInterval)
		elapsed += checkInterval
		fmt.Print(".")
	}

	fmt.Println()
	log.Fatal("❌ Authentication timeout")
}

func subscribeToPrivateChannels(client *ws.BaseWsClient) {
	productType := "USDT-FUTURES"

	fmt.Printf("📡 Subscribing to private channels for %s...\n", productType)

	// 1. Subscribe to order updates
	client.SubscribeOrders(productType, func(message string) {
		fmt.Printf("📋 ORDER UPDATE: %s\n", message)
	})

	// 2. Subscribe to fill/execution updates
	client.SubscribeFills(productType, func(message string) {
		fmt.Printf("✅ FILL UPDATE: %s\n", message)
	})

	// 3. Subscribe to position updates
	client.SubscribePositions(productType, func(message string) {
		fmt.Printf("📊 POSITION UPDATE: %s\n", message)
	})

	// 4. Subscribe to account balance updates
	client.SubscribeAccount(productType, func(message string) {
		fmt.Printf("💰 ACCOUNT UPDATE: %s\n", message)
	})

	// 5. Subscribe to plan order (trigger order) updates
	client.SubscribePlanOrders(productType, func(message string) {
		fmt.Printf("⚡ PLAN ORDER UPDATE: %s\n", message)
	})

	fmt.Printf("✅ Subscribed to %d private channels\n", client.GetSubscriptionCount())

	// Display subscription status
	displaySubscriptionStatus(client, productType)
}

func displaySubscriptionStatus(client *ws.BaseWsClient, productType string) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📡 PRIVATE CHANNEL SUBSCRIPTION STATUS")
	fmt.Println(strings.Repeat("=", 60))

	channels := map[string]string{
		"orders":     "📋 Order Updates",
		"fill":       "✅ Fill Updates",
		"positions":  "📊 Position Updates",
		"account":    "💰 Account Updates",
		"plan-order": "⚡ Plan Order Updates",
	}

	for channel, description := range channels {
		status := "❌ Not Subscribed"
		if client.IsSubscribed(channel, "", productType) {
			status = "✅ Subscribed"
		}
		fmt.Printf("%-25s %s\n", description, status)
	}

	fmt.Printf("\nTotal Active Subscriptions: %d\n", client.GetSubscriptionCount())
	fmt.Println(strings.Repeat("=", 60))
}

func defaultMessageHandler(message string) {
	// Handle general messages (login confirmations, etc.)
	fmt.Printf("ℹ️  SYSTEM: %s\n", message)
}

func errorMessageHandler(message string) {
	// Handle error messages
	fmt.Printf("❌ ERROR: %s\n", message)
}

func setupGracefulShutdownForPrivate(client *ws.BaseWsClient) {
	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("\n🎧 Listening for account updates... Press Ctrl+C to stop.")
	fmt.Println("💡 Try placing orders or making trades to see real-time updates!")

	// Wait for shutdown signal
	<-sigChan

	fmt.Println("\n🛑 Shutting down...")

	// Unsubscribe from all private channels
	productType := "USDT-FUTURES"
	fmt.Println("📤 Unsubscribing from private channels...")

	client.UnsubscribeOrders(productType)
	client.UnsubscribeFills(productType)
	client.UnsubscribePositions(productType)
	client.UnsubscribeAccount(productType)
	client.UnsubscribePlanOrders(productType)

	fmt.Printf("✅ Unsubscribed from all channels\n")

	// Close connection
	client.Close()
	fmt.Println("✅ Graceful shutdown complete")
}
