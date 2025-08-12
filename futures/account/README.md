# Account Management Services

This package contains services for managing Bitget futures account operations, including account information, leverage settings, margin management, and position modes.

## Services Overview

### Core Account Services

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `AccountInfoService` | Retrieve account information and balances | `Symbol()`, `ProductType()`, `MarginCoin()` |
| `AccountListService` | Get list of all futures accounts | `ProductType()` |
| `GetAccountBillService` | Fetch account transaction history | `ProductType()`, `StartTime()`, `EndTime()` |

### Leverage & Margin Management

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `SetLeverageService` | Set leverage for trading pairs | `Symbol()`, `ProductType()`, `MarginCoin()`, `Leverage()` |
| `AdjustMarginService` | Add or reduce margin for isolated positions | `Symbol()`, `ProductType()`, `Amount()`, `Type()` |
| `SetMarginModeService` | Switch between cross and isolated margin | `Symbol()`, `ProductType()`, `MarginMode()` |

### Position Configuration

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `SetPositionModeService` | Set position mode (one-way/hedge) | `ProductType()`, `PositionMode()` |

## Usage Examples

### Basic Account Information

```go
client := futures.NewClient(apiKey, secretKey, passphrase)

// Get account info for BTCUSDT
account, err := client.NewAccountInfoService().
    Symbol("BTCUSDT").
    ProductType(account.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Available Balance: %s USDT\n", account.Available)
fmt.Printf("Equity: %s USDT\n", account.Equity)
```

### Setting Leverage

```go
// Set 10x leverage for BTCUSDT
result, err := client.NewSetLeverageService().
    Symbol("BTCUSDT").
    ProductType(account.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Leverage("10").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Leverage set to: %s\n", result.LongLeverage)
```

### Margin Management

```go
// Add 100 USDT margin to isolated position
result, err := client.NewAdjustMarginService().
    Symbol("BTCUSDT").
    ProductType(account.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Amount("100").
    AddMargin(). // Helper method for type="add"
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Margin adjusted successfully\n")
```

### Position Mode Configuration

```go
// Set to hedge mode (allows both long and short positions)
err := client.NewSetPositionModeService().
    ProductType(account.ProductTypeUSDTFutures).
    PositionMode(account.PositionModeHedge).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Position mode set to hedge\n")
```

## API Endpoints

This package covers the following Bitget API endpoints:

- `/api/v2/mix/account/account` - Account information
- `/api/v2/mix/account/accounts` - Account list
- `/api/v2/mix/account/bill` - Account bills/transactions
- `/api/v2/mix/account/set-leverage` - Set leverage
- `/api/v2/mix/position/change-margin` - Adjust margin
- `/api/v2/mix/account/set-margin-mode` - Set margin mode
- `/api/v2/mix/account/set-position-mode` - Set position mode

## Types and Constants

### Product Types
- `ProductTypeUSDTFutures` - USDT-margined futures
- `ProductTypeCoinFutures` - Coin-margined futures
- `ProductTypeUSDCFutures` - USDC-margined futures

### Margin Modes
- `MarginModeCrossed` - Cross margin mode
- `MarginModeIsolated` - Isolated margin mode

### Position Modes
- `PositionModeOneWay` - One-way position mode
- `PositionModeHedge` - Hedge position mode (long + short)

### Hold Sides
- `HoldSideLong` - Long position
- `HoldSideShort` - Short position

## Error Handling

All services include comprehensive error handling with parameter validation:

```go
// Example of handling validation errors
result, err := client.NewSetLeverageService().
    Symbol("BTCUSDT").
    // Missing required ProductType
    Do(context.Background())

if err != nil {
    fmt.Printf("Error: %v\n", err) // "productType is required"
}
```

## Testing

Each service includes comprehensive tests with mock clients:

```bash
# Run all account service tests
go test ./futures/account/... -v

# Run specific service tests
go test ./futures/account/... -run TestAccountInfoService -v
```

## Notes

- All services require valid API credentials with futures trading permissions
- Rate limits apply according to Bitget API documentation
- Use testnet for development and testing
- All monetary values are returned as strings to preserve precision