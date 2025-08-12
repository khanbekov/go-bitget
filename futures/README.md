# Bitget Futures Trading SDK

A comprehensive Go SDK for Bitget futures trading API, providing organized access to all futures trading operations including account management, market data, position management, and advanced trading features.

## 🏗️ Architecture Overview

The SDK is organized into logical subdirectories for better maintainability and ease of use:

```
futures/
├── account/     📊 Account Management (7 services)
├── market/      📈 Market Data & Analytics (10 services)  
├── position/    📋 Position Management (4 services)
├── trading/     💱 Order Execution & History (13 services)
├── client.go    🔧 Main client and factory methods
├── constants.go 📍 Centralized API endpoints
├── websocket.go 🔄 WebSocket integration ⭐ NEW
├── examples/    📖 WebSocket examples ⭐ NEW
└── types.go     🏷️  Shared type definitions
```

**Total: 34 services covering all Bitget futures API endpoints + WebSocket integration**

## 🚀 Quick Start

### Installation

```bash
go get github.com/khanbekov/go-bitget
```

### Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/khanbekov/go-bitget/futures"
    "github.com/khanbekov/go-bitget/futures/account"
)

func main() {
    // Initialize client
    client := futures.NewClient(
        "your-api-key",
        "your-secret-key", 
        "your-passphrase",
    )

    // Create account service and get account information
    accountService := account.NewAccountInfoService(client)
    accountInfo, err := accountService.
        Symbol("BTCUSDT").
        ProductType(account.ProductTypeUSDTFutures).
        MarginCoin("USDT").
        Do(context.Background())

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Available Balance: %s USDT\n", accountInfo.Available)
}
```

## 📦 Service Categories

### 📊 Account Management (`account/`)

Manage your trading account, leverage, and margin settings.

```go
import "github.com/khanbekov/go-bitget/futures/account"

// Get account information
accountInfo, _ := account.NewAccountInfoService(client).
    Symbol("BTCUSDT").
    ProductType(account.ProductTypeUSDTFutures).
    Do(ctx)

// Set leverage  
account.NewSetLeverageService(client).
    Symbol("BTCUSDT").
    ProductType(account.ProductTypeUSDTFutures).
    Leverage("10").
    Do(ctx)

// Adjust margin for isolated positions
account.NewAdjustMarginService(client).
    Symbol("BTCUSDT").
    Amount("100").
    AddMargin().
    Do(ctx)
```

**Services**: AccountInfo, AccountList, SetLeverage, AdjustMargin, SetMarginMode, SetPositionMode, GetAccountBill

[📖 Full Account Documentation](account/README.md)

### 📈 Market Data (`market/`)

Access real-time and historical market data.

```go
import "github.com/khanbekov/go-bitget/futures/market"

// Get candlestick data
candles, _ := market.NewCandlestickService(client).
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Granularity("1m").
    Limit("100").
    Do(ctx)

// Get current funding rate
fundingRate, _ := market.NewCurrentFundingRateService(client).
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Do(ctx)

// Get order book
orderbook, _ := market.NewOrderBookService(client).
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Limit("20").
    Do(ctx)
```

**Services**: Candlestick, AllTickers, Ticker, OrderBook, RecentTrades, CurrentFundingRate, HistoryFundingRate, OpenInterest, SymbolPrice, Contracts

[📖 Full Market Data Documentation](market/README.md)

### 📋 Position Management (`position/`)

Monitor and manage your futures positions.

```go
import "github.com/khanbekov/go-bitget/futures/position"

// Get all open positions
positions, _ := position.NewAllPositionsService(client).
    ProductType(position.ProductTypeUSDTFutures).
    Do(ctx)

// Get specific position
positionInfo, _ := position.NewSinglePositionService(client).
    Symbol("BTCUSDT").
    ProductType(position.ProductTypeUSDTFutures).
    Do(ctx)

// Close position
position.NewClosePositionService(client).
    Symbol("BTCUSDT").
    ProductType(position.ProductTypeUSDTFutures).
    HoldSide("long").
    Do(ctx)
```

**Services**: AllPositions, SinglePosition, HistoryPositions, ClosePosition

[📖 Full Position Documentation](position/README.md)

### 💱 Trading Operations (`trading/`)

Execute trades, manage orders, and access trading history.

```go
import "github.com/khanbekov/go-bitget/futures/trading"

// Place a limit order
order, _ := trading.NewCreateOrderService(client).
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    Side("buy").
    OrderType("limit").
    Size("0.01").
    Price("45000").
    Do(ctx)

// Create conditional order
planOrder, _ := trading.NewCreatePlanOrderService(client).
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    PlanType("stop_loss").
    TriggerPrice("44000").
    Side("sell").
    OrderType("market").
    Size("0.01").
    Do(ctx)

// Batch orders
batchResult, _ := trading.NewCreateBatchOrdersService(client).
    ProductType(trading.ProductTypeUSDTFutures).
    Orders(multipleOrders).
    Do(ctx)
```

**Services**: CreateOrder, ModifyOrder, CancelOrder, CancelAllOrders, OrderDetails, PendingOrders, OrderHistory, FillHistory, CreatePlanOrder, ModifyPlanOrder, CancelPlanOrder, PendingPlanOrders, CreateBatchOrders

[📖 Full Trading Documentation](trading/README.md)

## 🔧 Configuration

### Product Types

```go
const (
    ProductTypeUSDTFutures = "USDT-FUTURES" // USDT-margined perpetual futures
    ProductTypeCoinFutures = "COIN-FUTURES" // Coin-margined futures
    ProductTypeUSDCFutures = "USDC-FUTURES" // USDC-margined futures
)
```

### Environment Setup

```bash
# Set up environment variables
export BITGET_API_KEY="your-api-key"
export BITGET_SECRET_KEY="your-secret-key" 
export BITGET_PASSPHRASE="your-passphrase"

# For testnet (demo trading)
export BITGET_BASE_URL="https://api.bitget.com" # Production
# export BITGET_BASE_URL="https://testnet.bitget.com" # Testnet
```

### Client Options

```go
client := futures.NewClient(apiKey, secretKey, passphrase)

// Set custom endpoint (e.g., for testnet)
client.SetApiEndpoint("https://testnet.bitget.com")

// Enable debug logging
client.Debug = true
```

## 📊 Supported Markets

### USDT-Margined Futures
- **BTC/USDT, ETH/USDT, ADA/USDT** and 100+ other pairs
- Perpetual contracts settled in USDT
- Up to 125x leverage available

### Coin-Margined Futures  
- **BTCUSD, ETHUSD** and other major pairs
- Contracts settled in base currency
- Quarterly and perpetual contracts

### USDC-Margined Futures
- Premium trading pairs settled in USDC
- Enhanced precision and stability

## 🛡️ Error Handling

The SDK provides comprehensive error handling with detailed validation:

```go
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    // Missing required fields
    Do(context.Background())

if err != nil {
    // Handle validation errors
    fmt.Printf("Order creation failed: %v\n", err)
    
    // Check for specific error types
    if strings.Contains(err.Error(), "productType is required") {
        // Handle missing product type
    }
}
```

## ⚡ Performance & Rate Limits

### API Rate Limits
- **Futures Trading**: 2,400 requests per minute
- **Market Data**: 20 requests per 2 seconds (no auth required)
- **WebSocket**: Real-time data with no polling limits

### Best Practices

```go
// Use batch operations when possible
batchService := client.NewCreateBatchOrdersService()
batchService.Orders(multipleOrders) // Up to 20 orders per request

// Cache market data to reduce API calls
ticker := client.NewTickerService().Symbol("BTCUSDT").Do(ctx)
// Cache ticker data for ~1-5 seconds

// Use WebSocket for real-time data
wsClient := ws.NewWsClient()
wsClient.SubscribePublicChannel("ticker:BTCUSDT:USDT-FUTURES", callback)
```

## 🧪 Testing

### Run All Tests

```bash
# Test all packages
go test ./futures/... -v

# Test specific categories  
go test ./futures/account/... -v
go test ./futures/market/... -v
go test ./futures/position/... -v
go test ./futures/trading/... -v
```

### Mock Testing

All services support mock testing:

```go
// Example test with mock client
mockClient := &MockClient{}
service := &account.AccountInfoService{C: mockClient}

mockClient.On("CallAPI", mock.Anything, "GET", "/api/v2/mix/account/account", 
    mock.Anything, mock.Anything, true).Return(mockResponse, nil, nil)

result, err := service.Do(context.Background())
```

## 🌐 WebSocket Integration ⭐ **NEW**

The futures package now includes integrated WebSocket functionality for seamless real-time data streaming:

### Easy WebSocket Integration

```go
import "github.com/khanbekov/go-bitget/futures"

func main() {
    // Initialize client and WebSocket manager
    client := futures.NewClient("api_key", "secret_key", "passphrase")
    wsManager := client.NewWebSocketManager()
    
    // Connect to public WebSocket
    wsManager.ConnectPublic()
    
    // Subscribe to real-time ticker updates
    wsManager.SubscribeToTicker("BTCUSDT", func(message string) {
        fmt.Printf("Price Update: %s\n", message)
    })
    
    // Subscribe to 1-minute candlesticks
    wsManager.SubscribeToCandlesticks("BTCUSDT", "1m", func(message string) {
        fmt.Printf("New Candle: %s\n", message)
    })
    
    // Subscribe to order book updates
    wsManager.SubscribeToOrderBook("BTCUSDT", 5, func(message string) {
        fmt.Printf("Order Book: %s\n", message)
    })
    
    // Keep running
    select {}
}
```

### High-Level Market Data Stream

```go
// Define symbols and create configuration
symbols := []string{"BTCUSDT", "ETHUSDT", "SOLUSDT"}
config := futures.MarketDataConfig{
    EnableTicker:    true,
    EnableCandles:   true,
    EnableOrderBook: true,
    EnableTrades:    true,
    
    CandleTimeframe: "5m",
    OrderBookLevels: 15,
    
    TickerHandler: func(message string) {
        fmt.Printf("📊 Ticker: %s\n", message)
    },
    CandleHandler: func(message string) {
        fmt.Printf("🕯️ Candle: %s\n", message)
    },
    // ... more handlers
}

// Create complete market data stream
wsManager.CreateMarketDataStream(ctx, symbols, config)
```

### Private Trading Stream

```go
// Connect to private channels (requires valid API credentials)
config := futures.TradingStreamConfig{
    EnableOrders:    true,
    EnableFills:     true,
    EnablePositions: true,
    EnableAccount:   true,
    
    OrderHandler: func(message string) {
        fmt.Printf("📋 Order Update: %s\n", message)
    },
    FillHandler: func(message string) {
        fmt.Printf("✅ Fill Update: %s\n", message)
    },
    // ... more handlers
}

wsManager.CreateTradingStream(ctx, "api_key", "passphrase", config)
```

### Combined REST + WebSocket Trading Bot

```go
// Monitor price via WebSocket, execute trades via REST API
wsManager.SubscribeToTicker("BTCUSDT", func(message string) {
    // Parse price from WebSocket message
    currentPrice := parsePrice(message)
    
    // Use REST API for actual trading
    if shouldBuyCondition(currentPrice) {
        orderService := trading.NewCreateOrderService(client)
        order, _ := orderService.
            Symbol("BTCUSDT").
            Side("buy").
            OrderType("limit").
            Size("0.001").
            Price(fmt.Sprintf("%.2f", currentPrice*0.99)).
            Do(ctx)
        
        fmt.Printf("Order placed: %s\n", order.OrderID)
    }
})

// Monitor order status via WebSocket
wsManager.SubscribeToOrders(func(message string) {
    fmt.Printf("Order status: %s\n", message)
    // Handle order fills, cancellations, etc.
})
```

### WebSocket Features

- **Market Data Streaming**: Real-time tickers, candlesticks, order books, trades
- **Private Data Streaming**: Order updates, fill notifications, position changes  
- **Easy Integration**: Seamless connection between REST API and WebSocket streams
- **High-Level Configurations**: Pre-built setups for common trading patterns
- **Connection Management**: Auto-reconnection, health monitoring, error handling
- **Authentication**: Automatic login for private channels

For complete examples, see `futures/examples/websocket_example.go`

## 🛠️ Development

### Project Structure

```
futures/
├── account/          # Account management services
│   ├── *.go         # Service implementations
│   ├── *_test.go    # Comprehensive test suites  
│   ├── types.go     # Local type definitions
│   └── README.md    # Package documentation
├── market/           # Market data services
├── position/         # Position management services  
├── trading/          # Trading operations
├── client.go        # Main client implementation
├── constants.go     # API endpoints and constants
└── README.md        # This file
```

### Contributing

1. **Follow Existing Patterns**: Each package follows the same service pattern
2. **Comprehensive Testing**: Include test coverage for all new features
3. **Documentation**: Update README files for any API changes
4. **Type Safety**: Use strong typing for all API parameters

## 📋 API Coverage

### ✅ Fully Implemented (53/53 endpoints)

- **Account Management**: 7/7 endpoints ✅
- **Market Data**: 10/10 endpoints ✅  
- **Position Management**: 4/4 endpoints ✅
- **Order Operations**: 13/13 endpoints ✅
- **Plan Orders**: 4/4 endpoints ✅
- **Batch Operations**: 1/1 endpoints ✅
- **History & Analytics**: 14/14 endpoints ✅

## 🔐 Security

- **Never commit API keys** to version control
- **Use environment variables** for credentials
- **Test with small amounts** before production use
- **Enable IP restrictions** in Bitget API settings
- **Use testnet** for development and testing

## 📞 Support

- **GitHub Issues**: [Report bugs and feature requests](https://github.com/khanbekov/go-bitget/issues)
- **Documentation**: See individual package README files
- **Bitget API Docs**: [Official API documentation](https://bitgetlimited.github.io/apidoc/en/mix/)

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🔄 Migration Guide

### For Existing Code

If you're upgrading from an older version, here's how to update your code:

#### Before Refactoring
```go
// Old monolithic approach - all services in one package
import "github.com/khanbekov/go-bitget/futures"

client := futures.NewClient(apiKey, secretKey, passphrase)
candles, err := client.NewCandlestickService().Do(ctx)
account, err := client.NewAccountInfoService().Do(ctx)
order, err := client.NewCreateOrderService().Do(ctx)
```

#### After Refactoring (Current)
```go
// New organized approach - services grouped by functionality
import (
    "github.com/khanbekov/go-bitget/futures"
    "github.com/khanbekov/go-bitget/futures/market"
    "github.com/khanbekov/go-bitget/futures/account"
    "github.com/khanbekov/go-bitget/futures/trading"
)

client := futures.NewClient(apiKey, secretKey, passphrase)

// Market data
candles, err := market.NewCandlestickService(client).Do(ctx)

// Account operations  
account, err := account.NewAccountInfoService(client).Do(ctx)

// Trading operations
order, err := trading.NewCreateOrderService(client).Do(ctx)
```

#### Migration Required

⚠️ **Factory methods moved**: Services are now created through subdirectory packages.

✅ **Same API**: All method signatures and functionality remain identical.

✅ **Better organization**: Services are now logically grouped and documented.

### Benefits of New Organization

1. **🎯 Focused Imports**: Only import what you need
2. **📚 Better Documentation**: Each package has specific docs  
3. **🔍 Easier Discovery**: Services grouped by functionality
4. **⚡ Improved Performance**: Smaller compile times with selective imports

---

**Happy Trading! 🚀**