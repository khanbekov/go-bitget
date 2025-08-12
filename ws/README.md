# WebSocket Client Documentation

This document provides comprehensive guidance for using the Bitget WebSocket client to stream real-time market data and account updates.

## Table of Contents

- [Overview](#overview)
- [Getting Started](#getting-started)
- [Public Channels](#public-channels)
- [Private Channels](#private-channels)
- [Examples](#examples)
- [Error Handling](#error-handling)
- [Best Practices](#best-practices)
- [API Reference](#api-reference)

## Overview

The Bitget WebSocket client provides real-time streaming of:

### Public Market Data (No Authentication Required)
- **Ticker Updates**: 24hr price statistics and trading volume
- **Candlestick Data**: OHLCV data for multiple timeframes
- **Order Book**: Live bid/ask levels with configurable depth
- **Trade Executions**: Real-time trade data
- **Mark Price**: Price used for PnL calculations
- **Funding Information**: Funding rates and timing

### Private Account Data (Authentication Required)
- **Order Updates**: Real-time order status changes
- **Fill Updates**: Trade execution confirmations
- **Position Updates**: Position changes and PnL updates
- **Account Updates**: Balance and margin changes
- **Plan Order Updates**: Trigger order status updates

## Getting Started

### Basic Setup

```go
package main

import (
    "fmt"
    "log"
    "os"
    "time"
    
    "github.com/khanbekov/go-bitget/ws"
    "github.com/rs/zerolog"
)

func main() {
    // Create logger
    logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
    
    // Create WebSocket client for public data
    client := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")
    
    // Set message handlers
    client.SetListener(messageHandler, errorHandler)
    
    // Connect and start reading
    client.Connect()
    client.ConnectWebSocket()
    client.StartReadLoop()
    
    // Subscribe to a ticker
    client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", tickerHandler)
    
    // Keep the program running
    select {}
}

func messageHandler(message string) {
    fmt.Println("Received:", message)
}

func errorHandler(message string) {
    fmt.Println("Error:", message)
}

func tickerHandler(message string) {
    fmt.Println("Ticker update:", message)
}
```

### Connection Endpoints

| Endpoint | Purpose | Authentication |
|----------|---------|----------------|
| `wss://ws.bitget.com/v2/ws/public` | Public market data | Not required |
| `wss://ws.bitget.com/v2/ws/private` | Private account data | Required |

## Public Channels

Public channels provide market data and don't require authentication.

### Available Channels

#### Ticker Channel
Real-time 24hr ticker statistics including price, volume, and change data.

```go
client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", func(message string) {
    fmt.Println("Ticker update:", message)
})
```

#### Candlestick Channels
Real-time OHLCV data for various timeframes.

**Supported Timeframes:**
- `ws.Timeframe1m`, `ws.Timeframe5m`, `ws.Timeframe15m`, `ws.Timeframe30m`
- `ws.Timeframe1h`, `ws.Timeframe4h`, `ws.Timeframe6h`, `ws.Timeframe12h`
- `ws.Timeframe1d`, `ws.Timeframe3d`, `ws.Timeframe1w`, `ws.Timeframe1M`

```go
client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", ws.Timeframe1m, func(message string) {
    fmt.Println("1m candle update:", message)
})
```

#### Order Book Channels
Real-time order book depth data.

```go
// Full order book
client.SubscribeOrderBook("BTCUSDT", "USDT-FUTURES", orderBookHandler)

// Top 5 levels
client.SubscribeOrderBook5("BTCUSDT", "USDT-FUTURES", orderBook5Handler)

// Top 15 levels
client.SubscribeOrderBook15("BTCUSDT", "USDT-FUTURES", orderBook15Handler)
```

#### Trade Channel
Real-time trade execution data.

```go
client.SubscribeTrades("BTCUSDT", "USDT-FUTURES", func(message string) {
    fmt.Println("Trade execution:", message)
})
```

#### Mark Price Channel
Mark price updates used for PnL calculations. 
**Note**: Bitget doesn't have a dedicated mark-price channel. This method subscribes to the ticker channel and extracts mark price data.

```go
client.SubscribeMarkPrice("BTCUSDT", "USDT-FUTURES", func(message string) {
    fmt.Println("Ticker update with mark price:", message)
})
```

#### Funding Channel
Funding rate and time information.

```go
client.SubscribeFundingTime("BTCUSDT", "USDT-FUTURES", func(message string) {
    fmt.Println("Funding update:", message)
})
```

### Product Types

| Product Type | Description |
|-------------|-------------|
| `"USDT-FUTURES"` | USDT-margined futures contracts |
| `"COIN-FUTURES"` | Coin-margined futures contracts |
| `"USDC-FUTURES"` | USDC-margined futures contracts |

## Private Channels

Private channels provide account-specific data and require authentication.

### Authentication Required

Before subscribing to private channels, you must authenticate:

```go
// Create client with secret key for signing
client := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/private", secretKey)

// Connect
client.Connect()
client.ConnectWebSocket()
client.StartReadLoop()

// Authenticate
client.Login(apiKey, passphrase, common.SHA256)

// Wait for login confirmation before subscribing
time.Sleep(2 * time.Second)

// Now you can subscribe to private channels
client.SubscribeOrders("USDT-FUTURES", orderHandler)
```

### Available Private Channels

#### Orders Channel
Real-time order status updates.

```go
client.SubscribeOrders("USDT-FUTURES", func(message string) {
    fmt.Println("Order update:", message)
})
```

#### Fills Channel
Real-time trade execution confirmations.

```go
client.SubscribeFills("USDT-FUTURES", func(message string) {
    fmt.Println("Fill update:", message)
})
```

#### Positions Channel
Real-time position updates.

```go
client.SubscribePositions("USDT-FUTURES", func(message string) {
    fmt.Println("Position update:", message)
})
```

#### Account Channel
Account balance and margin updates.

```go
client.SubscribeAccount("USDT-FUTURES", func(message string) {
    fmt.Println("Account update:", message)
})
```

#### Plan Orders Channel
Trigger/conditional order updates.

```go
client.SubscribePlanOrders("USDT-FUTURES", func(message string) {
    fmt.Println("Plan order update:", message)
})
```

## Examples

### Multiple Subscriptions

```go
// Subscribe to multiple channels for the same symbol
client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", tickerHandler)
client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", ws.Timeframe1m, candleHandler)
client.SubscribeOrderBook5("BTCUSDT", "USDT-FUTURES", orderBookHandler)
client.SubscribeTrades("BTCUSDT", "USDT-FUTURES", tradeHandler)
```

### Multiple Symbols

```go
symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}
for _, symbol := range symbols {
    client.SubscribeTicker(symbol, "USDT-FUTURES", tickerHandler)
}
```

### Mixed Public and Private

```go
// Public data (no auth needed)
publicClient := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")
publicClient.SetListener(messageHandler, errorHandler)
publicClient.Connect()
publicClient.ConnectWebSocket()
publicClient.StartReadLoop()

publicClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", tickerHandler)

// Private data (auth required)
privateClient := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/private", secretKey)
privateClient.SetListener(messageHandler, errorHandler)
privateClient.Connect()
privateClient.ConnectWebSocket()
privateClient.StartReadLoop()

privateClient.Login(apiKey, passphrase, common.SHA256)
time.Sleep(2 * time.Second) // Wait for authentication
privateClient.SubscribeOrders("USDT-FUTURES", orderHandler)
```

### Unsubscribing

```go
// Unsubscribe from specific channels
client.UnsubscribeTicker("BTCUSDT", "USDT-FUTURES")
client.UnsubscribeCandles("BTCUSDT", "USDT-FUTURES", ws.Timeframe1m)
client.UnsubscribeOrders("USDT-FUTURES")

// Or use generic unsubscribe
client.Unsubscribe("ticker", "BTCUSDT", "USDT-FUTURES")
```

## Error Handling

### Connection Monitoring

```go
// Check connection status
if !client.IsConnected() {
    log.Println("WebSocket not connected")
    client.ConnectWebSocket()
}

// Check authentication status for private channels
if !client.IsLoggedIn() {
    log.Println("Not authenticated")
    client.Login(apiKey, passphrase, common.SHA256)
}
```

### Automatic Reconnection

The client includes automatic reconnection with configurable timeouts:

```go
// Set custom reconnection timeout (default: 60 seconds)
client.SetReconnectionTimeout(30 * time.Second)

// Set custom connection check interval (default: 5 seconds)
client.SetCheckConnectionInterval(10 * time.Second)
```

### Error Handler

```go
func errorHandler(message string) {
    log.Printf("WebSocket error: %s", message)
    
    // Parse error and handle specific cases
    var errorData map[string]interface{}
    if err := json.Unmarshal([]byte(message), &errorData); err == nil {
        if code, ok := errorData["code"]; ok {
            switch code {
            case "30001":
                log.Println("Authentication failed")
            case "30002":
                log.Println("Invalid channel subscription")
            default:
                log.Printf("Unknown error code: %v", code)
            }
        }
    }
}
```

## Best Practices

### 1. Connection Management

```go
// Use separate clients for public and private data
publicClient := ws.NewBitgetBaseWsClient(logger, publicEndpoint, "")
privateClient := ws.NewBitgetBaseWsClient(logger, privateEndpoint, secretKey)
```

### 2. Subscription Management

```go
// Track your subscriptions
subscriptionCount := client.GetSubscriptionCount()
if subscriptionCount > 50 {
    log.Println("Warning: High number of subscriptions may impact performance")
}

// Check before subscribing
if !client.IsSubscribed("ticker", "BTCUSDT", "USDT-FUTURES") {
    client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
}
```

### 3. Authentication Flow

```go
// Always check login status before private subscriptions
if ws.RequiresAuth("orders") && !client.IsLoggedIn() {
    client.Login(apiKey, passphrase, common.SHA256)
    
    // Wait for authentication confirmation
    for !client.IsLoggedIn() {
        time.Sleep(100 * time.Millisecond)
    }
}
```

### 4. Message Processing

```go
func efficientHandler(message string) {
    // Process messages in goroutines to avoid blocking
    go func() {
        defer func() {
            if r := recover(); r != nil {
                log.Printf("Handler panic: %v", r)
            }
        }()
        
        // Your message processing logic here
        processMessage(message)
    }()
}
```

### 5. Resource Cleanup

```go
// Gracefully close connections
defer func() {
    client.Close()
}()

// Unsubscribe before closing
client.UnsubscribeTicker("BTCUSDT", "USDT-FUTURES")
client.UnsubscribeOrders("USDT-FUTURES")
```

## API Reference

### Client Creation

```go
func NewBitgetBaseWsClient(logger zerolog.Logger, url, secretKey string) *BaseWsClient
```

### Connection Methods

```go
func (c *BaseWsClient) Connect()
func (c *BaseWsClient) ConnectWebSocket()
func (c *BaseWsClient) StartReadLoop()
func (c *BaseWsClient) Close()
```

### Authentication

```go
func (c *BaseWsClient) Login(apiKey, passphrase string, signType common.SignType)
func (c *BaseWsClient) IsLoggedIn() bool
func (c *BaseWsClient) IsConnected() bool
```

### Public Channel Subscriptions

```go
func (c *BaseWsClient) SubscribeTicker(symbol, productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeCandles(symbol, productType, timeframe string, handler OnReceive)
func (c *BaseWsClient) SubscribeOrderBook(symbol, productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeOrderBook5(symbol, productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeOrderBook15(symbol, productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeTrades(symbol, productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeMarkPrice(symbol, productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeFundingTime(symbol, productType string, handler OnReceive)
```

### Private Channel Subscriptions

```go
func (c *BaseWsClient) SubscribeOrders(productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeFills(productType string, handler OnReceive)
func (c *BaseWsClient) SubscribePositions(productType string, handler OnReceive)
func (c *BaseWsClient) SubscribeAccount(productType string, handler OnReceive)
func (c *BaseWsClient) SubscribePlanOrders(productType string, handler OnReceive)
```

### Subscription Management

```go
func (c *BaseWsClient) Unsubscribe(channel, symbol, productType string)
func (c *BaseWsClient) GetActiveSubscriptions() map[SubscriptionArgs]OnReceive
func (c *BaseWsClient) IsSubscribed(channel, symbol, productType string) bool
func (c *BaseWsClient) GetSubscriptionCount() int
```

### Utility Functions

```go
func RequiresAuth(channel string) bool
```

### Configuration

```go
func (c *BaseWsClient) SetListener(msgListener OnReceive, errorListener OnReceive)
func (c *BaseWsClient) SetReconnectionTimeout(timeout time.Duration)
func (c *BaseWsClient) SetCheckConnectionInterval(interval time.Duration)
```

## Constants

### Timeframes
```go
const (
    Timeframe1m  = "1m"   // 1 minute
    Timeframe5m  = "5m"   // 5 minutes
    Timeframe15m = "15m"  // 15 minutes
    Timeframe30m = "30m"  // 30 minutes
    Timeframe1h  = "1h"   // 1 hour
    Timeframe4h  = "4h"   // 4 hours
    Timeframe6h  = "6h"   // 6 hours
    Timeframe12h = "12h"  // 12 hours
    Timeframe1d  = "1d"   // 1 day
    Timeframe3d  = "3d"   // 3 days
    Timeframe1w  = "1w"   // 1 week
    Timeframe1M  = "1M"   // 1 month
)
```

### Channels
```go
const (
    // Public channels
    ChannelTicker      = "ticker"
    ChannelCandle      = "candle"
    ChannelBooks       = "books"
    ChannelBooks5      = "books5"
    ChannelBooks15     = "books15"
    ChannelTrade       = "trade"
    ChannelMarkPrice   = "mark-price"
    ChannelFundingTime = "funding-time"
    
    // Private channels
    ChannelOrders    = "orders"
    ChannelFill      = "fill"
    ChannelPositions = "positions"
    ChannelAccount   = "account"
    ChannelPlanOrder = "plan-order"
)
```

## Support

For issues and questions:
- Check the [examples](../examples/) directory for more usage patterns
- Review the test files for additional implementation details
- Open an issue on the GitHub repository

## Related Documentation

- [Main SDK Documentation](../README.md)
- [Futures API Documentation](../futures/README.md)
- [Examples](../examples/)