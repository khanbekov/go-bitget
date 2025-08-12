# UTA (Unified Trading Account) SDK

This package provides a comprehensive Go SDK for Bitget's UTA (Unified Trading Account) API. The UTA API is Bitget's next-generation trading API that unifies spot, margin, and futures trading into a single account structure.

## Features

### Account Management
- ✅ **Account Information**: Get account settings, asset mode, holding mode
- ✅ **Account Assets**: Query account equity, balances, and asset details
- ✅ **Funding Assets**: Access funding account balances
- ✅ **Fee Rates**: Get trading fee rates for different symbols and categories
- ✅ **Account Configuration**: Set holding mode, leverage, account switching
- ✅ **Transfer Operations**: Internal transfers between account types
- 🔄 **Sub-Account Management**: Create, manage, and query sub-accounts (stubs implemented)
- 🔄 **Deposit & Withdrawal**: Address management and transaction history (stubs implemented)

### Trading Operations
- ✅ **Order Management**: Place, cancel, modify orders
- 🔄 **Batch Operations**: Batch order operations (up to 20 orders) (stubs implemented)
- 🔄 **Strategy Orders**: TPSL orders with advanced features (stubs implemented)
- 🔄 **Position Management**: Query and manage positions (stubs implemented)

### Market Data
- ✅ **Tickers**: Real-time price data across all categories
- ✅ **Candlesticks**: OHLCV data with multiple timeframes
- 🔄 **Order Books**: Depth data (stubs implemented)
- 🔄 **Historical Data**: Extended historical data access (stubs implemented)

### Product Categories
- `SPOT` - Spot trading
- `MARGIN` - Margin trading  
- `USDT-FUTURES` - USDT perpetual futures
- `COIN-FUTURES` - Coin-margined futures
- `USDC-FUTURES` - USDC futures

## Quick Start

### Installation

```go
go get github.com/khanbekov/go-bitget
```

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/khanbekov/go-bitget/uta"
)

func main() {
    // Create UTA client
    client := uta.NewClient("your_api_key", "your_secret_key", "your_passphrase")
    
    // Get USDT futures tickers
    tickers, err := client.NewGetTickersService().
        Category(uta.CategoryUSDTFutures).
        Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    for _, ticker := range tickers {
        if ticker.Symbol == "BTCUSDT" {
            fmt.Printf("BTC Price: %s\n", ticker.LastPrice)
            break
        }
    }
    
    // Get account information
    accountInfo, err := client.NewAccountInfoService().Do(context.Background())
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Account Mode: %s, Holding Mode: %s\n", 
        accountInfo.AssetMode, accountInfo.HoldingMode)
}
```

### Advanced Examples

#### Market Data

```go
// Get single symbol ticker
ticker, err := client.NewGetTickersService().
    Category(uta.CategorySpot).
    Symbol("BTCUSDT").
    Do(ctx)

// Get candlestick data
candles, err := client.NewGetCandlesticksService().
    Category(uta.CategoryUSDTFutures).
    Symbol("BTCUSDT").
    Interval(uta.Interval1m).
    Limit("100").
    Do(ctx)
```

#### Account Operations

```go
// Get account assets
assets, err := client.NewAccountAssetsService().Do(ctx)

// Get funding assets for specific coin
fundingAssets, err := client.NewAccountFundingAssetsService().
    Coin("USDT").
    Do(ctx)

// Get fee rates
feeRate, err := client.NewAccountFeeRateService().
    Symbol("BTCUSDT").
    Category(uta.CategoryUSDTFutures).
    Do(ctx)

// Set leverage
err = client.NewSetLeverageService().
    Category(uta.CategoryUSDTFutures).
    Symbol("BTCUSDT").
    Leverage("10").
    Do(ctx)
```

#### Trading Operations

```go
// Place a limit order
order, err := client.NewPlaceOrderService().
    Symbol("BTCUSDT").
    Category(uta.CategoryUSDTFutures).
    Side(uta.SideBuy).
    OrderType(uta.OrderTypeLimit).
    Size("0.001").
    Price("50000").
    TimeInForce(uta.TimeInForceGTC).
    Do(ctx)

// Cancel an order
canceledOrder, err := client.NewCancelOrderService().
    Symbol("BTCUSDT").
    Category(uta.CategoryUSDTFutures).
    OrderId(order.OrderID).
    Do(ctx)
```

#### Transfer Operations

```go
// Transfer between accounts
transferResult, err := client.NewTransferService().
    FromType(uta.AccountTypeSpot).
    ToType(uta.AccountTypeUSDTFutures).
    Amount("100").
    Coin("USDT").
    Do(ctx)
```

## API Structure

### Fluent API Pattern

All services support method chaining for intuitive usage:

```go
result, err := client.NewPlaceOrderService().
    Symbol("BTCUSDT").
    Category(uta.CategoryUSDTFutures).
    Side(uta.SideBuy).
    OrderType(uta.OrderTypeLimit).
    Size("0.001").
    Price("50000").
    ClientOid("my-order-123").
    TimeInForce(uta.TimeInForceGTC).
    Do(context.Background())
```

### Error Handling

The SDK provides structured error handling:

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

### Context Support

All operations support `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := service.Do(ctx)
```

## Testing

The package includes comprehensive tests:

```bash
# Run unit tests
go test ./uta -v

# Run specific tests
go test ./uta -run TestAccountInfoService -v

# Run integration tests (requires API credentials in .env)
go test ./uta -run TestIntegration -v
```

### Integration Tests

For integration tests, create a `.env` file in the project root:

```env
BITGET_API_KEY=your_api_key
BITGET_SECRET_KEY=your_secret_key
BITGET_PASSPHRASE=your_passphrase
```

## Rate Limits

The UTA API has the following rate limits:

- **Individual Operations**: 10/sec/UID (place, cancel, modify orders)
- **Batch Operations**: 5/sec/UID (max 20 orders/batch)
- **Mass Operations**: 1-5/sec/UID (cancel all, close all)
- **Data Queries**: 20/sec/UID (account, orders, positions)
- **Public Data**: 20/sec/IP (market data)

## Constants

### Order Sides
```go
uta.SideBuy   // "buy"
uta.SideSell  // "sell"
```

### Order Types
```go
uta.OrderTypeLimit   // "limit"
uta.OrderTypeMarket  // "market"
```

### Time in Force
```go
uta.TimeInForceGTC      // "gtc" - Good 'til canceled
uta.TimeInForceIOC      // "ioc" - Immediate or cancel
uta.TimeInForceFOK      // "fok" - Fill or kill
uta.TimeInForcePostOnly // "post_only" - Post only (maker only)
```

### Product Categories
```go
uta.CategorySpot        // "SPOT"
uta.CategoryMargin      // "MARGIN"
uta.CategoryUSDTFutures // "USDT-FUTURES"
uta.CategoryCoinFutures // "COIN-FUTURES"
uta.CategoryUSDCFutures // "USDC-FUTURES"
```

### Candlestick Intervals
```go
uta.Interval1m   // "1m"
uta.Interval5m   // "5m"
uta.Interval15m  // "15m"
uta.Interval1H   // "1H"
uta.Interval1D   // "1D"
// ... and more
```

## Implementation Status

- ✅ **Fully Implemented**: Core functionality with tests
- 🔄 **Stub Implemented**: Interface defined, implementation pending
- ❌ **Not Implemented**: Not yet started

### Fully Implemented Services
- Account Information & Settings
- Account Assets & Funding Assets
- Fee Rate Queries
- Basic Account Configuration (holding mode, leverage)
- Internal Transfers
- Market Data (tickers, candlesticks)
- Basic Order Operations (place, cancel)

### Partially Implemented (Stubs)
- Advanced trading features (batch orders, strategy orders)
- Complete position management
- Deposit & withdrawal operations
- Sub-account management
- Complete market data (order books, historical data)

## UTA vs Futures API Differences

| Feature | Futures API | UTA API |
|---------|-------------|---------|
| **API Version** | `/api/v2/` | `/api/v3/` |
| **Account Type** | Futures-specific | Unified (spot/margin/futures) |
| **Product Support** | USDT/COIN/USDC futures only | All trading categories |
| **WebSocket Trading** | Market data only | Order operations via WebSocket |
| **Batch Operations** | REST API only | WebSocket batch (20 orders) |
| **Strategy Orders** | Basic plan orders | Advanced TPSL with partial position support |

## Contributing

1. Follow the existing patterns and architecture
2. Add comprehensive tests for new functionality
3. Update documentation and examples
4. Ensure backward compatibility where possible

## License

This project is licensed under the MIT License.

## Disclaimer

This SDK is for educational and development purposes. Always test thoroughly in a sandbox environment before using with real funds.