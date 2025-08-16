# Go Bitget SDK

[![GitHub release](https://img.shields.io/github/release/khanbekov/go-bitget.svg)](https://github.com/khanbekov/go-bitget/releases)
[![Go version](https://img.shields.io/badge/go-1.23.4+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

> ⚠️ **ALPHA VERSION WARNING** ⚠️  
> This is version **v0.0.1** - an alpha release. The API is **NOT STABLE** and **WILL CHANGE** without backward compatibility guarantees until v1.0.0. Use at your own risk in production environments.

A comprehensive Go SDK for the Bitget cryptocurrency exchange API, providing both REST API and WebSocket functionality for futures trading operations.

## Features

### REST API Support
- ✅ **Futures Trading**: Complete order management (create, modify, cancel, batch)
- ✅ **Advanced Orders**: Trigger/conditional orders (plan orders, stop-loss, take-profit)
- ✅ **Account Management**: Balance queries, position management, leverage control
- ✅ **Account Configuration**: Margin mode, position mode, account list, margin adjustment
- ✅ **Market Data**: Candlesticks, tickers, order books, recent trades, contracts
- ✅ **Advanced Market Data**: Funding rates, open interest, symbol prices
- ✅ **Historical Data**: Order history, fill history, position history
- ✅ **Error Handling**: Comprehensive error handling with retry logic

### WebSocket Support (Unified Implementation)
- ✅ **Production-Ready Client**: BaseWsClient with comprehensive features
- ✅ **Public Channels**: Real-time market data (tickers, candles, order books, trades, mark price, funding)
- ✅ **Private Channels**: Account updates (orders, fills, positions, balance, plan orders)
- ✅ **Advanced Features**: Rate limiting (10 msg/sec), automatic reconnection, subscription restoration
- ✅ **Connection Management**: Health monitoring, heartbeat mechanism, graceful shutdown
- ✅ **Type Safety**: Structured data types for all message formats

## Installation

```bash
# Install the latest alpha version
go get github.com/khanbekov/go-bitget@v0.0.1

# Or get the latest development version (not recommended)
go get github.com/khanbekov/go-bitget@latest
```

### Version Compatibility

| Version | Status | Stability | Backward Compatibility |
|---------|--------|-----------|----------------------|
| v0.0.1  | Alpha  | ❌ Unstable | ❌ No guarantees |
| v0.x.x  | Alpha/Beta | ⚠️ Limited | ❌ Breaking changes expected |
| v1.0.0+ | Stable | ✅ Stable | ✅ Semantic versioning |

## Quick Start

### REST API Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/khanbekov/go-bitget/futures"
)

func main() {
    // Create futures client
    client := futures.NewClient("your_api_key", "your_secret_key", "your_passphrase")
    
    // Get 24hr ticker
    ticker, err := client.NewAllTickersService().
        ProductType(futures.ProductTypeUSDTFutures).
        Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("BTC Price: %s\n", ticker[0].LastPrice)
}
```

### WebSocket Example

```go
package main

import (
    "fmt"
    "os"
    "time"

    "github.com/khanbekov/go-bitget/ws"
    "github.com/rs/zerolog"
)

func main() {
    logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
    
    // Create WebSocket client
    client := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")
    
    // Set message handlers
    client.SetListener(
        func(msg string) { fmt.Println("Message:", msg) },
        func(err string) { fmt.Println("Error:", err) },
    )
    
    // Connect and subscribe
    client.Connect()
    client.ConnectWebSocket()
    client.StartReadLoop()
    
    time.Sleep(2 * time.Second) // Wait for connection
    
    // Subscribe to Bitcoin ticker
    client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", func(message string) {
        fmt.Println("BTC Ticker:", message)
    })
    
    // Keep running
    select {}
}
```

## Documentation

### API Documentation
- **[Futures API](futures/)** - Complete futures trading API documentation with 34+ services
- **[UTA API](uta/)** - Unified Trading Account API (recommended for new development)
- **[Common Utilities](common/)** - Shared utilities, authentication, and error handling

### WebSocket Documentation
- **[WebSocket Guide](ws/README.md)** - Comprehensive unified WebSocket implementation
- **[Channel Reference](ws/README.md#api-reference)** - All available channels and subscription methods
- **[Type Definitions](ws/types.go)** - Structured data types for all WebSocket messages

### Examples
- **[Basic Examples](examples/)** - Simple usage examples
- **[WebSocket Examples](examples/websocket/)** - Real-time data streaming examples
  - [Basic Public Channels](examples/websocket/basic_public_channels.go)
  - [Private Channels with Authentication](examples/websocket/private_channels.go)
  - [Multiple Symbols Monitoring](examples/websocket/multiple_symbols.go)
  - [Advanced Usage Patterns](examples/websocket/advanced_usage.go)
  - [Mixed Public/Private Channels](examples/websocket/mixed_channels.go)

## Configuration

### Environment Variables

Create a `.env` file in your project root:

```env
BITGET_API_KEY=your_api_key_here
BITGET_SECRET_KEY=your_secret_key_here
BITGET_PASSPHRASE=your_passphrase_here
```

### API Endpoints

| Environment | REST API Base URL | WebSocket Public | WebSocket Private |
|------------|-------------------|------------------|-------------------|
| Production | `https://api.bitget.com` | `wss://ws.bitget.com/v2/ws/public` | `wss://ws.bitget.com/v2/ws/private` |

## Supported Features

### Futures Trading Operations

#### Account & Positions
- ✅ Get account information and balances
- ✅ Get all positions (open and closed)
- ✅ Get single position details
- ✅ Get position history
- ✅ Close positions

#### Order Management
- ✅ Create single orders (market, limit, stop)
- ✅ Create batch orders (up to 20 orders)
- ✅ Modify existing orders
- ✅ Cancel single orders
- ✅ Cancel all orders for a symbol
- ✅ Get order details
- ✅ Get order history
- ✅ Get pending orders

#### Advanced Orders (Plan Orders)
- ✅ Create plan orders (trigger/conditional orders)
- ✅ Modify plan orders
- ✅ Cancel plan orders
- ✅ Get pending plan orders
- ✅ Support for stop-loss, take-profit, normal plan, track plan, stop-surplus

#### Market Data
- ✅ Get candlestick/OHLCV data
- ✅ Get all tickers (24hr statistics)
- ✅ Get single ticker
- ✅ Get order book depth
- ✅ Get recent trades
- ✅ Get contract specifications

#### Advanced Market Data
- ✅ Get current funding rates
- ✅ Get historical funding rates (with pagination)
- ✅ Get open interest data
- ✅ Get symbol prices (mark, index, last price)

#### Risk Management & Account Configuration
- ✅ Set/modify leverage
- ✅ Get bill/account history
- ✅ Get fill history
- ✅ Set margin mode (isolated/cross)
- ✅ Set position mode (one-way/hedge)
- ✅ Get account list
- ✅ Adjust position margin

### WebSocket Channels

#### Public Channels (No Authentication Required)
- 📊 **Ticker**: 24hr price statistics and volume
- 🕯️ **Candles**: Real-time OHLCV data (12 timeframes)
- 📚 **Order Book**: Live bid/ask levels (full, top 5, top 15)
- 💰 **Trades**: Real-time trade executions
- 🎯 **Mark Price**: Price used for PnL calculations
- 💸 **Funding**: Funding rate and timing information

#### Private Channels (Authentication Required)
- 📋 **Orders**: Real-time order status updates
- ✅ **Fills**: Trade execution confirmations
- 📊 **Positions**: Position changes and PnL updates
- 💰 **Account**: Balance and margin updates
- ⚡ **Plan Orders**: Trigger/conditional order updates

### Supported Product Types
- `USDT-FUTURES` - USDT-margined futures contracts
- `COIN-FUTURES` - Coin-margined futures contracts
- `USDC-FUTURES` - USDC-margined futures contracts

## Architecture

### Service-Oriented Design
The SDK uses a service-oriented architecture with organized package structure:

```go
// Futures API (Legacy) - Organized into subdirectories
client := futures.NewClient(apiKey, secretKey, passphrase)
tickers := client.NewAllTickersService().ProductType("USDT-FUTURES").Do(ctx)
candles := client.NewCandlestickService().Symbol("BTCUSDT").Granularity("1m").Do(ctx)

// UTA API (Recommended) - Unified Trading Account
utaClient := uta.NewClient(apiKey, secretKey, passphrase) // Auto-detects demo mode
assets := utaClient.NewAccountAssetsService().Do(ctx)
order := utaClient.NewPlaceOrderService().Symbol("BTCUSDT").Side("buy").Do(ctx)

// WebSocket (Unified Implementation)
wsClient := ws.NewBitgetBaseWsClient(logger, endpoint, secretKey)
wsClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", tickerHandler)
```

### Package Organization

- **`futures/`**: Legacy futures API organized into 4 subdirectories (`account/`, `market/`, `position/`, `trading/`)
- **`uta/`**: Unified Trading Account API (recommended for new development)
- **`ws/`**: Unified WebSocket implementation with production-ready features
- **`common/`**: Shared utilities, authentication, error handling, and type definitions

### Fluent API Pattern
All services support method chaining for intuitive usage:

```go
result, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    ProductType(futures.ProductTypeUSDTFutures).
    Side(futures.SideBuy).
    OrderType(futures.OrderTypeLimit).
    Size("0.001").
    Price("50000").
    Do(context.Background())
```

### Error Handling
Comprehensive error handling with structured error types:

```go
result, err := service.Do(ctx)
if err != nil {
    if apiErr, ok := err.(*common.APIError); ok {
        fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Message, apiErr.Code)
    } else {
        fmt.Printf("Network Error: %v\n", err)
    }
    return
}
```

## Development

### Building the Project

```bash
# Build the application
go build -o app .

# Run directly
go run main.go

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...
```

### Code Generation and Documentation

```bash
# Generate comprehensive documentation (creates docs/ directory)
bash generate-docs.sh    # Unix/Linux/macOS
generate-docs.bat         # Windows

# This creates:
# - docs/index.html - HTML overview and navigation
# - docs/*.txt - Full package documentation  
# - docs/*-summary.txt - Package summaries
# - Live server: godoc -http=:6060
```

### Module Management

```bash
# Update dependencies
go mod tidy

# Download dependencies
go mod download
```

## Testing

The SDK includes comprehensive test coverage:

```bash
# Run all tests
go test ./...

# Run specific package tests
go test -v ./futures/
go test -v ./ws/
go test -v ./common/

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Structure
- **Unit Tests**: All services and utilities with mock clients
- **Integration Tests**: Real API testing with your own credentials
- **WebSocket Tests**: Channel subscription and management tests
- **End-to-End Tests**: Complete workflow validation

### Integration Testing (Real API)

Test against real Bitget API endpoints with your own credentials:

```bash
# Setup
cp tests/configs/integration.example.json tests/configs/integration.json
# Edit with your API keys and enable demo trading

# Run integration tests
tests/scripts/run-integration-tests.sh                    # Unix/Linux/macOS
tests/scripts/run-integration-tests.bat                   # Windows

# Direct testing
go test -tags=integration ./tests/integration/suites -v
```

**Features**:
- ✅ **Safe Testing**: Demo trading mode and read-only operations
- ✅ **Comprehensive Coverage**: All account, market, and position endpoints
- ✅ **Detailed Reports**: JSON and HTML test reports with metrics
- ✅ **Selective Testing**: Enable/disable specific endpoints
- ✅ **Error Recovery**: Built-in retry logic and error handling

For detailed integration testing guide, see [`tests/INTEGRATION_TESTING.md`](tests/INTEGRATION_TESTING.md).

## Examples and Use Cases

### Production-Ready WebSocket Features

The unified WebSocket implementation provides enterprise-grade features:

```go
// Create production-ready WebSocket client
wsClient := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")

// Built-in features:
// ✅ Rate limiting (10 messages/second)
// ✅ Automatic reconnection with configurable timeout
// ✅ Subscription restoration after reconnection  
// ✅ Connection health monitoring
// ✅ Heartbeat mechanism with ping/pong
// ✅ Thread-safe subscription management
// ✅ Graceful shutdown handling

// Configure connection parameters
wsClient.SetCheckConnectionInterval(5 * time.Second)   // Health check interval
wsClient.SetReconnectionTimeout(120 * time.Second)    // Reconnection timeout

// Set up listeners and connect
wsClient.SetListener(messageHandler, errorHandler)
wsClient.Connect()
wsClient.ConnectWebSocket()
wsClient.StartReadLoop()

// Subscribe to multiple channels
wsClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", tickerHandler)
wsClient.SubscribeOrderBook5("ETHUSDT", "USDT-FUTURES", orderbookHandler)
wsClient.SubscribeCandles("ADAUSDT", "USDT-FUTURES", ws.Timeframe1m, candleHandler)
```

### Trading Bot Example
```go
// Monitor market and place orders based on conditions
client := futures.NewClient(apiKey, secretKey, passphrase)
wsClient := ws.NewBitgetBaseWsClient(logger, publicEndpoint, "")

// Subscribe to price updates with automatic reconnection
wsClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", func(message string) {
    // Parse structured ticker data
    var tickerData ws.TickerData
    if err := json.Unmarshal([]byte(message), &tickerData); err == nil {
        price := tickerData.LastPriceFloat
        if shouldBuy(price) {
            client.NewCreateOrderService().
                Symbol("BTCUSDT").
                Side("buy").
                Size("0.001").
                Do(context.Background())
        }
    }
})
```

### Portfolio Monitoring
```go
// Monitor multiple positions in real-time
wsClient.SubscribePositions("USDT-FUTURES", func(message string) {
    // Update portfolio dashboard
    updatePortfolio(message)
})

wsClient.SubscribeAccount("USDT-FUTURES", func(message string) {
    // Update account balance display
    updateAccountBalance(message)
})
```

### Market Data Analysis
```go
// Collect and analyze market data
symbols := []string{"BTCUSDT", "ETHUSDT", "ADAUSDT"}
for _, symbol := range symbols {
    wsClient.SubscribeCandles(symbol, "USDT-FUTURES", ws.Timeframe1m, func(msg string) {
        // Store candle data for analysis
        analyzeMarketData(symbol, msg)
    })
}
```

## Dependencies

Core dependencies:
- `github.com/valyala/fasthttp` - High-performance HTTP client
- `github.com/json-iterator/go` - Fast JSON processing
- `github.com/rs/zerolog` - Structured logging
- `github.com/gorilla/websocket` - WebSocket implementation

Testing dependencies:
- `github.com/stretchr/testify` - Testing framework with assertions and mocking

Utility dependencies:
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/robfig/cron/v3` - Cron job scheduling
- `github.com/google/uuid` - UUID generation

## Support and Contributing

### Getting Help
- 📖 Check the package documentation and [examples](examples/)
- 📋 Read the [Development Guide](docs/DEVELOPMENT.md) for comprehensive development guidance
- 🐛 Report issues on [GitHub Issues](https://github.com/khanbekov/go-bitget/issues)
- 💬 Join discussions in [GitHub Discussions](https://github.com/khanbekov/go-bitget/discussions)

### Contributing
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines
- Follow Go best practices and idioms
- Add tests for new functionality
- Update documentation for API changes
- Use the existing code style and patterns

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Disclaimer

This SDK is for educational and development purposes. Always test thoroughly in a sandbox environment before using with real funds. The authors are not responsible for any financial losses incurred through the use of this software.

## Changelog

### Latest Release - v0.0.1 (Alpha Release)
- ✅ Complete futures REST API implementation (37+ services across 4 organized directories)
- ✅ Unified Trading Account (UTA) API with demo trading auto-detection
- ✅ Advanced trading features (plan orders, account configuration, batch operations)
- ✅ Advanced market data (funding rates, open interest, symbol prices, historical data)
- ✅ **Unified WebSocket implementation**: Production-ready BaseWsClient with 11+ subscription methods
- ✅ **WebSocket features**: Rate limiting, automatic reconnection, subscription restoration, health monitoring
- ✅ Comprehensive error handling with retry logic and structured error types
- ✅ Extensive test coverage with mock clients and integration tests
- ✅ Type-safe WebSocket data structures and comprehensive documentation

**⚠️ Alpha Limitations:**
- API may change without notice
- Limited production testing
- Breaking changes expected before v1.0.0
- Use with caution in live trading

For detailed changes, see [CHANGELOG.md](CHANGELOG.md).