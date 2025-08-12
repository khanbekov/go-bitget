# Position Management Services

This package contains services for managing futures positions, including retrieving current positions, historical positions, and closing positions.

## Services Overview

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `AllPositionsService` | Get all open positions | `ProductType()`, `MarginCoin()` |
| `SinglePositionService` | Get specific position details | `Symbol()`, `ProductType()`, `MarginCoin()` |
| `HistoryPositionsService` | Retrieve historical/closed positions | `ProductType()`, `StartTime()`, `EndTime()` |
| `ClosePositionService` | Close positions (market/limit) | `Symbol()`, `ProductType()`, `HoldSide()` |

## Usage Examples

### View All Open Positions

```go
client := futures.NewClient(apiKey, secretKey, passphrase)

// Get all open USDT futures positions
positions, err := client.NewAllPositionsService().
    ProductType(position.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, pos := range positions {
    fmt.Printf("Symbol: %s, Side: %s, Size: %s, PnL: %s\n",
        pos.Symbol, pos.HoldSide, pos.Size, pos.UnrealizedPL)
}
```

### Get Specific Position

```go
// Get BTCUSDT position details
position, err := client.NewSinglePositionService().
    Symbol("BTCUSDT").
    ProductType(position.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Position Details:\n")
fmt.Printf("  Symbol: %s\n", position.Symbol)
fmt.Printf("  Size: %s\n", position.Size)
fmt.Printf("  Entry Price: %s\n", position.AverageOpenPrice)
fmt.Printf("  Mark Price: %s\n", position.MarkPrice)
fmt.Printf("  Unrealized PnL: %s\n", position.UnrealizedPL)
fmt.Printf("  Margin: %s\n", position.Margin)
fmt.Printf("  Leverage: %s\n", position.Leverage)
```

### Close a Position

```go
// Close long position with market order
result, err := client.NewClosePositionService().
    Symbol("BTCUSDT").
    ProductType(position.ProductTypeUSDTFutures).
    HoldSide(position.HoldSideLong).
    MarginCoin("USDT").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Position closed successfully: %+v\n", result)
```

### Historical Positions

```go
// Get historical positions from last 7 days
endTime := time.Now().Unix() * 1000
startTime := time.Now().AddDate(0, 0, -7).Unix() * 1000

history, err := client.NewHistoryPositionsService().
    ProductType(position.ProductTypeUSDTFutures).
    StartTime(fmt.Sprintf("%d", startTime)).
    EndTime(fmt.Sprintf("%d", endTime)).
    PageSize("50").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, pos := range history.List {
    fmt.Printf("Closed Position: %s, Side: %s, PnL: %s, Close Time: %s\n",
        pos.Symbol, pos.HoldSide, pos.RealizedPL, pos.UTime)
}
```

## Position Data Structure

The `Position` struct contains comprehensive position information:

```go
type Position struct {
    Symbol           string // Trading pair (e.g., "BTCUSDT")
    Size             string // Position size
    HoldSide         string // Position side ("long" or "short")
    AverageOpenPrice string // Average entry price
    MarkPrice        string // Current mark price
    UnrealizedPL     string // Unrealized profit/loss
    RealizedPL       string // Realized profit/loss
    Margin           string // Position margin
    MarginMode       string // Margin mode (isolated/cross)
    MarginRatio      string // Margin ratio
    Leverage         string // Position leverage
    AutoMargin       string // Auto margin flag
    CrossMarginLeverage string // Cross margin leverage
    Available        string // Available balance
    Locked           string // Locked balance
    Total            string // Total balance
    UTime            string // Last update time
    CTime            string // Creation time
    // ... additional fields for risk management
}
```

## API Endpoints

This package covers the following Bitget API endpoints:

- `/api/v2/mix/position/all-position` - All open positions
- `/api/v2/mix/position/single-position` - Single position details
- `/api/v2/mix/position/history-position` - Historical positions
- `/api/v2/mix/order/close-positions` - Close positions

## Position Sides

- `HoldSideLong` - Long position (bought)
- `HoldSideShort` - Short position (sold)

## Margin Modes

- `MarginModeCrossed` - Cross margin (shared across positions)
- `MarginModeIsolated` - Isolated margin (individual position margin)

## Product Types

- `ProductTypeUSDTFutures` - USDT-margined perpetual futures
- `ProductTypeCoinFutures` - Coin-margined futures contracts
- `ProductTypeUSDCFutures` - USDC-margined perpetual futures

## Risk Management

### Position Monitoring

```go
// Monitor position risk
positions, err := client.NewAllPositionsService().
    ProductType(position.ProductTypeUSDTFutures).
    Do(context.Background())

for _, pos := range positions {
    marginRatio, _ := strconv.ParseFloat(pos.MarginRatio, 64)
    
    if marginRatio > 0.8 { // 80% margin ratio warning
        fmt.Printf("⚠️  High risk position: %s, Margin Ratio: %.2f%%\n",
            pos.Symbol, marginRatio*100)
    }
}
```

### Automatic Position Closing

```go
// Close positions with high risk
for _, pos := range positions {
    unrealizedPL, _ := strconv.ParseFloat(pos.UnrealizedPL, 64)
    
    if unrealizedPL < -1000 { // Close if loss > 1000 USDT
        _, err := client.NewClosePositionService().
            Symbol(pos.Symbol).
            ProductType(position.ProductTypeUSDTFutures).
            HoldSide(pos.HoldSide).
            MarginCoin("USDT").
            Do(context.Background())
            
        if err != nil {
            log.Printf("Failed to close position %s: %v\n", pos.Symbol, err)
        }
    }
}
```

## Error Handling

Position services include parameter validation:

```go
// This will return a validation error
result, err := client.NewClosePositionService().
    Symbol("BTCUSDT").
    // Missing required ProductType and HoldSide
    Do(context.Background())

if err != nil {
    fmt.Printf("Error: %v\n", err) // Parameter validation failed
}
```

## Testing

Run position service tests:

```bash
# Run all position service tests
go test ./futures/position/... -v

# Test specific services
go test ./futures/position/... -run TestAllPositionsService -v
go test ./futures/position/... -run TestClosePositionService -v
```

## Performance Tips

1. **Batch Position Queries**: Use `AllPositionsService` instead of multiple `SinglePositionService` calls
2. **Efficient Pagination**: Use `PageSize` and pagination for large historical queries
3. **WebSocket Integration**: Consider WebSocket subscriptions for real-time position updates
4. **Caching**: Cache position data for short periods to reduce API calls

## WebSocket Alternative

For real-time position updates:

```go
// Subscribe to position updates via WebSocket
wsClient := ws.NewWsClient()
wsClient.SubscribePrivateChannel("positions", positionCallback)
```