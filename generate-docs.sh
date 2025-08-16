#!/bin/bash

# Modern Go SDK Documentation Generator
# Uses industry-standard tools and practices

set -e

echo "Generating modern documentation for Bitget Go SDK..."

# Ensure we have required tools
command -v go >/dev/null 2>&1 || { echo "ERROR: go is required but not installed." >&2; exit 1; }

# Create docs directory
mkdir -p docs

echo "Generating package documentation..."

# Generate comprehensive package documentation
cat > docs/API_REFERENCE.md << 'EOF'
# API Reference

This document provides comprehensive API reference for the Bitget Go SDK.

## Quick Navigation

- [Futures API](#futures-api) - Legacy futures trading API
- [UTA API](#uta-api) - Unified Trading Account (Recommended)
- [WebSocket API](#websocket-api) - Real-time data streaming
- [Common Utilities](#common-utilities) - Shared utilities and types

---

## Futures API

The Futures API provides access to futures trading operations organized into 4 main categories:

### Account Operations
- Account information and balances
- Position management and history  
- Leverage and margin configuration
- Bill/transaction history

### Market Data Operations
- Candlestick/OHLCV data
- Ticker information (24hr stats)
- Order book depth data
- Recent trade history
- Contract specifications

### Position Operations
- Get all positions (open/closed)
- Single position details
- Position history
- Close positions

### Trading Operations  
- Create/modify/cancel orders
- Batch order operations
- Plan orders (conditional/trigger)
- Order and fill history

---

## UTA API (Recommended)

The Unified Trading Account API is Bitget's next-generation API that supports spot, margin, and futures trading in a single account.

### Key Features
- Auto-detection of demo trading mode
- Unified account management
- Simplified API structure
- Better error handling

### Main Operations
- Account asset management
- Order placement and management
- Market data access
- Account configuration

---

## WebSocket API

Production-ready WebSocket implementation with enterprise features.

### Features
- Rate limiting (10 messages/second)
- Automatic reconnection with configurable timeout
- Subscription restoration after reconnection
- Connection health monitoring  
- Heartbeat mechanism with ping/pong
- Thread-safe subscription management

### Public Channels
- Ticker (24hr price statistics)
- Candles (real-time OHLCV data)
- Order Book (live bid/ask levels)
- Trades (real-time executions)
- Mark Price (PnL calculation price)
- Funding (rates and timing)

### Private Channels  
- Orders (real-time status updates)
- Fills (execution confirmations)
- Positions (changes and PnL)
- Account (balance updates)
- Plan Orders (trigger order updates)

---

## Common Utilities

Shared utilities used across all packages:

### Authentication
- HMAC-SHA256 request signing
- API key management
- Timestamp handling

### Error Handling
- Structured API error types
- Network error handling
- Retry logic with exponential backoff

### Data Types
- Common constants and enums
- Utility functions
- Type conversions

EOF

echo "Generating usage examples..."

# Generate comprehensive examples
cat > docs/EXAMPLES.md << 'EOF'
# Usage Examples

This document provides comprehensive usage examples for the Bitget Go SDK.

## Table of Contents

- [Quick Start](#quick-start)
- [Futures API Examples](#futures-api-examples)
- [UTA API Examples](#uta-api-examples)
- [WebSocket Examples](#websocket-examples)
- [Error Handling](#error-handling)

---

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/khanbekov/go-bitget/futures"
    "github.com/khanbekov/go-bitget/uta" 
    "github.com/khanbekov/go-bitget/ws"
    "github.com/rs/zerolog"
)

func main() {
    // Load credentials
    apiKey := os.Getenv("BITGET_API_KEY")
    secretKey := os.Getenv("BITGET_SECRET_KEY") 
    passphrase := os.Getenv("BITGET_PASSPHRASE")

    // Create clients
    futuresClient := futures.NewClient(apiKey, secretKey, passphrase)
    utaClient := uta.NewClient(apiKey, secretKey, passphrase)
    
    logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
    wsClient := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")
}
```

---

## Futures API Examples

### Market Data

```go
// Get all tickers
tickers, err := futuresClient.NewAllTickersService().
    ProductType(futures.ProductTypeUSDTFutures).
    Do(context.Background())

// Get candlestick data  
candles, err := futuresClient.NewCandlestickService().
    Symbol("BTCUSDT").
    ProductType(futures.ProductTypeUSDTFutures).
    Granularity("1m").
    Limit("100").
    Do(context.Background())

// Get order book
orderbook, err := futuresClient.NewOrderBookService().
    Symbol("BTCUSDT").
    ProductType(futures.ProductTypeUSDTFutures).
    Limit("20").
    Do(context.Background())
```

### Trading Operations

```go
// Create a limit order
order, err := futuresClient.NewCreateOrderService().
    Symbol("BTCUSDT").
    ProductType(futures.ProductTypeUSDTFutures).
    Side(futures.SideBuy).
    OrderType(futures.OrderTypeLimit).
    Size("0.001").
    Price("50000").
    Do(context.Background())

// Create batch orders
orders := []futures.BatchOrderRequest{
    {
        Symbol:      "BTCUSDT",
        ProductType: futures.ProductTypeUSDTFutures,
        Side:        futures.SideBuy,
        OrderType:   futures.OrderTypeLimit,
        Size:        "0.001",
        Price:       "49000",
    },
    // Add more orders...
}

batchResult, err := futuresClient.NewCreateBatchOrdersService().
    Orders(orders).
    Do(context.Background())
```

---

## UTA API Examples

### Account Management

```go
// Get account assets
assets, err := utaClient.NewAccountAssetsService().Do(context.Background())

// Get account information
accountInfo, err := utaClient.NewAccountInfoService().Do(context.Background())
```

### Trading

```go
// Place an order (UTA)
order, err := utaClient.NewPlaceOrderService().
    Symbol("BTCUSDT").
    Side("buy").
    OrderType("limit").
    Size("0.001").
    Price("50000").
    Do(context.Background())
```

---

## WebSocket Examples

### Public Channel Subscriptions

```go
// Set up WebSocket client
logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
wsClient := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")

// Configure and connect
wsClient.SetListener(
    func(msg string) { fmt.Println("Message:", msg) },
    func(err string) { fmt.Println("Error:", err) },
)

wsClient.Connect()
wsClient.ConnectWebSocket()
wsClient.StartReadLoop()

time.Sleep(2 * time.Second) // Wait for connection

// Subscribe to ticker
wsClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", func(message string) {
    var ticker ws.TickerData
    if err := json.Unmarshal([]byte(message), &ticker); err == nil {
        fmt.Printf("BTC Price: %f\n", ticker.LastPriceFloat)
    }
})

// Subscribe to candlesticks
wsClient.SubscribeCandles("ETHUSDT", "USDT-FUTURES", ws.Timeframe1m, func(message string) {
    fmt.Println("ETH 1m Candle:", message)
})

// Subscribe to order book
wsClient.SubscribeOrderBook5("ADAUSDT", "USDT-FUTURES", func(message string) {
    fmt.Println("ADA Order Book:", message)
})
```

### Private Channel Subscriptions

```go
// Private WebSocket requires authentication
privateWs := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/private", secretKey)

privateWs.SetListener(messageHandler, errorHandler)
privateWs.Connect()
privateWs.ConnectWebSocket()
privateWs.StartReadLoop()

// Authenticate
privateWs.Login(apiKey, passphrase, common.SHA256)
time.Sleep(2 * time.Second) // Wait for login

// Subscribe to private channels
privateWs.SubscribeOrders("USDT-FUTURES", func(message string) {
    fmt.Println("Order Update:", message)
})

privateWs.SubscribePositions("USDT-FUTURES", func(message string) {
    fmt.Println("Position Update:", message)
})

privateWs.SubscribeAccount("USDT-FUTURES", func(message string) {
    fmt.Println("Account Update:", message)  
})
```

---

## Error Handling

### API Error Handling

```go
result, err := client.NewSomeService().Do(context.Background())
if err != nil {
    if apiErr, ok := err.(*common.APIError); ok {
        // Handle API-specific errors
        fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Message, apiErr.Code)
        switch apiErr.Code {
        case "40001":
            // Handle invalid parameters
        case "40004":
            // Handle insufficient balance
        default:
            // Handle other API errors
        }
    } else {
        // Handle network/connection errors
        fmt.Printf("Network Error: %v\n", err)
    }
    return
}
```

### WebSocket Error Handling

```go
wsClient.SetListener(
    func(message string) {
        // Handle successful messages
        fmt.Println("Received:", message)
    },
    func(error string) {
        // Handle WebSocket errors
        fmt.Printf("WebSocket Error: %s\n", error)
        
        // Implement custom error recovery if needed
        // The client will automatically attempt reconnection
    },
)
```

EOF

echo "Generating development guide..."

# Generate development and contribution guide
cat > docs/DEVELOPMENT.md << 'EOF'
# Development Guide

This guide covers development practices, testing, and contribution guidelines for the Bitget Go SDK.

## Development Setup

### Prerequisites

- Go 1.23.4 or later
- Git
- Make (optional, for build automation)

### Quick Start

```bash
# Clone repository
git clone https://github.com/khanbekov/go-bitget.git
cd go-bitget

# Install dependencies
go mod download

# Run tests
go test ./...

# Build
go build .
```

### Environment Configuration

Create a `.env` file:

```env
BITGET_API_KEY=your_api_key
BITGET_SECRET_KEY=your_secret_key  
BITGET_PASSPHRASE=your_passphrase
```

## Project Structure

```
├── futures/              # Futures API (legacy, organized in subdirectories)
│   ├── account/         # Account management services
│   ├── market/          # Market data services
│   ├── position/        # Position management services
│   └── trading/         # Trading operations services
├── uta/                 # Unified Trading Account API (recommended)
├── ws/                  # Unified WebSocket implementation
├── common/              # Shared utilities and types
├── examples/            # Usage examples
└── docs/                # Generated documentation
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test -v ./futures/
go test -v ./uta/
go test -v ./ws/
go test -v ./common/

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run with race detection
go test -race ./...
```

### Test Structure

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test service interactions with mock clients
- **WebSocket Tests**: Test subscription management and message handling

### Writing Tests

Follow the established patterns:

```go
func TestServiceName_Method_Scenario(t *testing.T) {
    // Setup
    mockClient := &MockClient{}
    service := NewServiceName(mockClient)
    
    // Configure expectations
    mockClient.On("callAPI", mock.Anything, "GET", endpoint, params, []byte(nil), false).
        Return(mockResponse, &fasthttp.ResponseHeader{}, nil)
    
    // Execute
    result, err := service.Do(context.Background())
    
    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockClient.AssertExpectations(t)
}
```

## Code Style and Standards

### Go Standards

- Follow Go best practices and idioms
- Use `gofmt` for code formatting
- Use `golint` and `go vet` for code quality
- Write comprehensive GoDoc comments

### SDK Patterns

- **Service Pattern**: Each API endpoint group has its own service
- **Fluent API**: Support method chaining for intuitive usage
- **Context Support**: All operations support context.Context
- **Error Handling**: Use structured error types

### Documentation

- Add GoDoc comments for all public functions and types
- Include usage examples in comments
- Update README.md for API changes
- Add examples to the `examples/` directory

## Contributing

### Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Update documentation
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Pull Request Guidelines

- Include comprehensive description
- Add tests for new functionality
- Update documentation as needed
- Ensure CI passes
- Follow semantic versioning for breaking changes

### Code Review

- All changes require code review
- Focus on correctness, performance, and maintainability
- Check test coverage and documentation
- Verify examples and usage patterns

EOF

echo "Updating HTML index with modern content..."

# Generate modern HTML index
cat > docs/index.html << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Bitget Go SDK Documentation</title>
    <style>
        :root {
            --primary-color: #0066cc;
            --secondary-color: #f8f9fa;
            --text-color: #333;
            --border-color: #dee2e6;
            --success-color: #28a745;
            --warning-color: #ffc107;
        }
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            line-height: 1.6;
            color: var(--text-color);
            background-color: #fff;
        }
        
        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }
        
        .header {
            text-align: center;
            margin-bottom: 50px;
            padding: 40px 0;
            background: linear-gradient(135deg, var(--primary-color), #004499);
            color: white;
            border-radius: 10px;
        }
        
        .header h1 {
            font-size: 3rem;
            margin-bottom: 10px;
            font-weight: 300;
        }
        
        .header p {
            font-size: 1.2rem;
            opacity: 0.9;
        }
        
        .badge {
            display: inline-block;
            padding: 4px 12px;
            background: var(--success-color);
            color: white;
            border-radius: 20px;
            font-size: 0.8rem;
            font-weight: 600;
            margin: 0 5px;
        }
        
        .badge.warning {
            background: var(--warning-color);
            color: #333;
        }
        
        .features {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
            gap: 30px;
            margin: 50px 0;
        }
        
        .feature-card {
            background: var(--secondary-color);
            padding: 30px;
            border-radius: 10px;
            border: 1px solid var(--border-color);
            transition: transform 0.3s ease, box-shadow 0.3s ease;
        }
        
        .feature-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 10px 25px rgba(0,0,0,0.1);
        }
        
        .feature-card h3 {
            color: var(--primary-color);
            margin-bottom: 15px;
            font-size: 1.5rem;
        }
        
        .feature-card ul {
            list-style: none;
            padding: 0;
        }
        
        .feature-card li {
            padding: 5px 0;
            position: relative;
            padding-left: 20px;
        }
        
        .feature-card li:before {
            content: "✅";
            position: absolute;
            left: 0;
        }
        
        .quick-start {
            background: #f8f9fa;
            padding: 30px;
            border-radius: 10px;
            margin: 40px 0;
            border-left: 4px solid var(--primary-color);
        }
        
        .quick-start h3 {
            color: var(--primary-color);
            margin-bottom: 20px;
        }
        
        pre {
            background: #2d3748;
            color: #e2e8f0;
            padding: 20px;
            border-radius: 8px;
            overflow-x: auto;
            font-size: 0.9rem;
            line-height: 1.4;
        }
        
        .docs-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
            margin: 40px 0;
        }
        
        .doc-link {
            display: block;
            padding: 20px;
            background: white;
            border: 2px solid var(--border-color);
            border-radius: 8px;
            text-decoration: none;
            color: var(--text-color);
            transition: all 0.3s ease;
        }
        
        .doc-link:hover {
            border-color: var(--primary-color);
            background: var(--secondary-color);
            transform: translateY(-2px);
        }
        
        .doc-link h4 {
            color: var(--primary-color);
            margin-bottom: 10px;
            font-size: 1.2rem;
        }
        
        .footer {
            text-align: center;
            margin-top: 60px;
            padding: 30px 0;
            border-top: 1px solid var(--border-color);
            color: #6c757d;
        }
        
        .version-info {
            background: #fff3cd;
            border: 1px solid #ffeaa7;
            border-radius: 8px;
            padding: 20px;
            margin: 30px 0;
        }
        
        .version-info h4 {
            color: #856404;
            margin-bottom: 10px;
        }
        
        @media (max-width: 768px) {
            .header h1 {
                font-size: 2rem;
            }
            
            .container {
                padding: 10px;
            }
            
            .features {
                grid-template-columns: 1fr;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Bitget Go SDK</h1>
            <p>Production-ready Go SDK for Bitget cryptocurrency exchange</p>
            <div style="margin-top: 20px;">
                <span class="badge">Go 1.23.4+</span>
                <span class="badge">37+ Services</span>
                <span class="badge">WebSocket Real-time</span>
                <span class="badge warning">Alpha v0.0.1</span>
            </div>
        </div>

        <div class="version-info">
            <h4>⚠️ Alpha Version Notice</h4>
            <p>This is version v0.0.1 - an alpha release. The API is <strong>NOT STABLE</strong> and <strong>WILL CHANGE</strong> without backward compatibility guarantees until v1.0.0. Use with caution in production environments.</p>
        </div>

        <div class="features">
            <div class="feature-card">
                <h3>Futures API</h3>
                <p>Complete futures trading implementation with 37+ services organized into 4 directories.</p>
                <ul>
                    <li>Account & Position Management</li>
                    <li>Advanced Order Types & Batch Operations</li>
                    <li>Real-time Market Data</li>
                    <li>Risk Management & Configuration</li>
                </ul>
            </div>

            <div class="feature-card">
                <h3>UTA API (Recommended)</h3>
                <p>Unified Trading Account API - Bitget's next-generation API for all trading operations.</p>
                <ul>
                    <li>Auto-detection of Demo Trading Mode</li>
                    <li>Unified Spot, Margin & Futures</li>
                    <li>Simplified API Structure</li>
                    <li>Enhanced Error Handling</li>
                </ul>
            </div>

            <div class="feature-card">
                <h3>WebSocket (Unified)</h3>
                <p>Production-ready WebSocket implementation with enterprise features.</p>
                <ul>
                    <li>Rate Limiting (10 msg/sec)</li>
                    <li>Automatic Reconnection</li>
                    <li>Subscription Restoration</li>
                    <li>Health Monitoring & Heartbeat</li>
                    <li>Type-safe Data Structures</li>
                </ul>
            </div>
        </div>

        <div class="quick-start">
            <h3>Quick Start</h3>
            <pre><code>package main

import (
    "context"
    "fmt" 
    "log"

    "github.com/khanbekov/go-bitget/futures"
    "github.com/khanbekov/go-bitget/uta"
    "github.com/khanbekov/go-bitget/ws"
)

func main() {
    // Futures API (Legacy)
    client := futures.NewClient(apiKey, secretKey, passphrase)
    tickers, err := client.NewAllTickersService().
        ProductType(futures.ProductTypeUSDTFutures).
        Do(context.Background())
    
    // UTA API (Recommended) 
    utaClient := uta.NewClient(apiKey, secretKey, passphrase)
    assets, err := utaClient.NewAccountAssetsService().Do(context.Background())
    
    // WebSocket Real-time Data
    wsClient := ws.NewBitgetBaseWsClient(logger, publicEndpoint, "")
    wsClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", tickerHandler)
}</code></pre>
        </div>

        <h2 style="text-align: center; margin: 50px 0 30px 0; color: var(--primary-color);">Documentation</h2>

        <div class="docs-grid">
            <a href="API_REFERENCE.md" class="doc-link">
                <h4>API Reference</h4>
                <p>Comprehensive API documentation for all packages and services.</p>
            </a>

            <a href="EXAMPLES.md" class="doc-link">
                <h4>Examples</h4>
                <p>Practical usage examples and code snippets for common operations.</p>
            </a>

            <a href="DEVELOPMENT.md" class="doc-link">
                <h4>Development Guide</h4>
                <p>Development setup, testing, and contribution guidelines.</p>
            </a>

            <a href="../README.md" class="doc-link">
                <h4>Main README</h4>
                <p>Project overview, installation, and getting started guide.</p>
            </a>

            <a href="../CONTRIBUTING.md" class="doc-link">
                <h4>Contributing</h4>
                <p>Guidelines for contributing to the project.</p>
            </a>

            <a href="../LICENSE" class="doc-link">
                <h4>License</h4>
                <p>MIT License - Open source licensing terms.</p>
            </a>

            <a href="../examples/" class="doc-link">
                <h4>Code Examples</h4>
                <p>Working code examples and demo applications.</p>
            </a>
        </div>

        <div style="background: var(--secondary-color); padding: 30px; border-radius: 10px; margin: 40px 0;">
            <h3 style="color: var(--primary-color); margin-bottom: 20px;">Live Documentation</h3>
            <p>For interactive Go documentation with full package browsing:</p>
            <pre style="margin-top: 15px;"><code>godoc -http=:6060</code></pre>
            <p style="margin-top: 15px;">Then visit: <a href="http://localhost:6060/pkg/github.com/khanbekov/go-bitget/" style="color: var(--primary-color);">http://localhost:6060/pkg/github.com/khanbekov/go-bitget/</a></p>
        </div>

        <div class="footer">
            <p>Generated by modern documentation tools | Bitget Go SDK v0.0.1 (Alpha)</p>
            <p style="margin-top: 10px;">
                <a href="https://github.com/khanbekov/go-bitget" style="color: var(--primary-color); text-decoration: none;">
                    GitHub Repository
                </a>
            </p>
        </div>
    </div>
</body>
</html>
EOF

echo "Modern documentation generated successfully!"
echo ""
echo "Files created:"
echo "  docs/API_REFERENCE.md - Comprehensive API documentation"
echo "  docs/EXAMPLES.md - Usage examples and patterns"  
echo "  docs/DEVELOPMENT.md - Development and contribution guide"
echo "  docs/index.html - Modern HTML overview"
echo ""
echo "Next steps:"
echo "  1. Open docs/index.html in your browser"
echo "  2. Review the generated documentation"
echo "  3. Run: godoc -http=:6060 for live docs"
echo "  4. Consider setting up GitHub Pages for hosting"
echo ""
echo "Modern improvements:"
echo "  - Responsive design with modern CSS"
echo "  - Comprehensive API reference in Markdown"
echo "  - Detailed usage examples"
echo "  - Development and contribution guidelines"
echo "  - Mobile-friendly interface"
echo "  - Professional styling and navigation"