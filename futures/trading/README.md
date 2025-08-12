# Trading Services

This package contains all services for executing trades, managing orders, and retrieving trading history on Bitget futures markets.

## Services Overview

### Core Order Management

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `CreateOrderService` | Place new orders (limit/market) | `Symbol()`, `Size()`, `Side()`, `OrderType()`, `Price()` |
| `ModifyOrderService` | Modify existing orders | `OrderId()`, `NewPrice()`, `NewSize()` |
| `CancelOrderService` | Cancel individual orders | `Symbol()`, `OrderId()` |
| `CancelAllOrdersService` | Cancel all orders | `ProductType()`, `MarginCoin()` |
| `OrderDetailsService` | Get detailed order information | `Symbol()`, `OrderId()` |

### Order History & Tracking

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `PendingOrdersService` | Get open/pending orders | `Symbol()`, `ProductType()` |
| `OrderHistoryService` | Retrieve historical orders | `ProductType()`, `StartTime()`, `EndTime()` |
| `FillHistoryService` | Get order execution history | `ProductType()`, `OrderId()`, `StartTime()` |

### Advanced Order Types

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `CreatePlanOrderService` | Place trigger/conditional orders | `TriggerPrice()`, `PlanType()`, `TriggerType()` |
| `ModifyPlanOrderService` | Modify existing plan orders | `OrderId()`, `TriggerPrice()` |
| `CancelPlanOrderService` | Cancel plan orders | `OrderId()`, `PlanType()` |
| `PendingPlanOrdersService` | Get pending plan orders | `Symbol()`, `PlanType()` |

### Batch Operations

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `CreateBatchOrdersService` | Place up to 20 orders simultaneously | `AddOrder()`, `Orders()` |

## Usage Examples

### Basic Order Placement

```go
client := futures.NewClient(apiKey, secretKey, passphrase)

// Place a limit buy order
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Side(trading.SideBuy).
    OrderType(trading.OrderTypeLimit).
    Size("0.01").
    Price("45000").
    TimeInForce(trading.TimeInForceGTC).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Order placed: ID=%s, ClientOid=%s\n", order.OrderId, order.ClientOid)
```

### Market Orders

```go
// Place a market sell order
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    Side(trading.SideSell).
    OrderType(trading.OrderTypeMarket).
    Size("0.01").
    Do(context.Background())
```

### Order Management

```go
// Modify an existing order
modifiedOrder, err := client.NewModifyOrderService().
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    OrderId("123456789").
    NewClientOrderId("my-new-order-id").
    NewPrice("46000").
    NewSize("0.02").
    Do(context.Background())

// Cancel a specific order
result, err := client.NewCancelOrderService().
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    MarginCoin("USDT").
    OrderId("123456789").
    Do(context.Background())
```

### Plan Orders (Conditional Orders)

```go
// Create a stop-loss plan order
planOrder, err := client.NewCreatePlanOrderService().
    Symbol("BTCUSDT").
    ProductType(trading.ProductTypeUSDTFutures).
    PlanType(trading.PlanTypeStopLoss).
    TriggerPrice("44000").
    TriggerType(trading.TriggerTypeMarkPrice).
    Side(trading.SideSell).
    OrderType(trading.OrderTypeMarket).
    Size("0.01").
    Do(context.Background())

fmt.Printf("Stop-loss order created: %s\n", planOrder.OrderId)
```

### Batch Order Operations

```go
// Create multiple orders in one request
batchService := client.NewCreateBatchOrdersService().
    ProductType(trading.ProductTypeUSDTFutures)

// Add multiple orders
orders := []trading.BatchOrderInfo{
    {
        Symbol: "BTCUSDT",
        MarginMode: trading.MarginModeCrossed,
        MarginCoin: "USDT",
        Size: "0.01",
        Price: "45000",
        SideType: trading.SideBuy,
        OrderType: trading.OrderTypeLimit,
        ClientOrderId: "batch-order-1",
    },
    {
        Symbol: "ETHUSDT",
        MarginMode: trading.MarginModeCrossed,
        MarginCoin: "USDT",
        Size: "0.1",
        Price: "3000",
        SideType: trading.SideBuy,
        OrderType: trading.OrderTypeLimit,
        ClientOrderId: "batch-order-2",
    },
}

result, err := batchService.Orders(orders).Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Successful orders: %d, Failed orders: %d\n",
    len(result.SuccessList), len(result.FailureList))
```

### Trading History

```go
// Get order history for the last 24 hours
endTime := time.Now().Unix() * 1000
startTime := time.Now().AddDate(0, 0, -1).Unix() * 1000

history, err := client.NewOrderHistoryService().
    ProductType(trading.ProductTypeUSDTFutures).
    StartTime(fmt.Sprintf("%d", startTime)).
    EndTime(fmt.Sprintf("%d", endTime)).
    PageSize("50").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, order := range history.List {
    fmt.Printf("Order: %s, Symbol: %s, Side: %s, Size: %s, Price: %s, Status: %s\n",
        order.OrderId, order.Symbol, order.Side, order.Size, order.Price, order.State)
}
```

### Fill History

```go
// Get execution history
fills, err := client.NewFillHistoryService().
    ProductType(trading.ProductTypeUSDTFutures).
    Symbol("BTCUSDT").
    PageSize("100").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, fill := range fills.List {
    fmt.Printf("Fill: TradeId=%s, Size=%s, Price=%s, Fee=%s, Role=%s\n",
        fill.TradeId, fill.Size, fill.Price, fill.Fee, fill.Role)
}
```

## API Endpoints

This package covers the following Bitget API endpoints:

### Basic Trading
- `/api/v2/mix/order/place-order` - Place orders
- `/api/v2/mix/order/modify-order` - Modify orders
- `/api/v2/mix/order/cancel-order` - Cancel orders
- `/api/v2/mix/order/cancel-all-orders` - Cancel all orders
- `/api/v2/mix/order/detail` - Order details

### Order Management
- `/api/v2/mix/order/current` - Pending orders
- `/api/v2/mix/order/history` - Order history
- `/api/v2/mix/order/fills` - Fill history

### Plan Orders
- `/api/v2/mix/order/place-plan-order` - Create plan orders
- `/api/v2/mix/order/modify-plan-order` - Modify plan orders
- `/api/v2/mix/order/cancel-plan-order` - Cancel plan orders
- `/api/v2/mix/order/plan-current` - Pending plan orders

### Batch Operations
- `/api/v2/mix/order/batch-orders` - Batch order placement

## Order Types

### Basic Order Types
- `OrderTypeLimit` - Limit orders with specified price
- `OrderTypeMarket` - Market orders (immediate execution)

### Plan Order Types
- `PlanTypeNormalPlan` - Normal trigger orders
- `PlanTypeTrackPlan` - Trailing stop orders
- `PlanTypeStopLoss` - Stop loss orders
- `PlanTypeTakeProfit` - Take profit orders
- `PlanTypeStopSurplus` - Stop surplus orders

### Trigger Types
- `TriggerTypeFillPrice` - Trigger based on fill price
- `TriggerTypeMarkPrice` - Trigger based on mark price

## Order Sides

- `SideBuy` - Buy orders (long position)
- `SideSell` - Sell orders (short position)

## Time in Force Options

- `TimeInForceGTC` - Good Till Cancel (default)
- `TimeInForceIOC` - Immediate or Cancel
- `TimeInForceFOK` - Fill or Kill

## Advanced Features

### Self-Trade Prevention

```go
// Prevent self-trading
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    // ... other parameters ...
    SelfTradePrevention(trading.STPCancelTaker).
    Do(context.Background())
```

Available STP modes:
- `STPNone` - No self-trade prevention
- `STPCancelTaker` - Cancel taker order
- `STPCancelMaker` - Cancel maker order  
- `STPCancelBoth` - Cancel both orders

### Reduce-Only Orders

```go
// Reduce-only order (only reduces position size)
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    // ... other parameters ...
    ReduceOnly(true).
    Do(context.Background())
```

### Stop-Loss and Take-Profit

```go
// Order with built-in stop-loss and take-profit
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    // ... basic parameters ...
    PresetStopLossPrice("44000").       // Stop-loss at $44,000
    PresetStopSurplusPrice("50000").    // Take-profit at $50,000
    Do(context.Background())
```

## Error Handling

All trading services include comprehensive validation:

```go
// This will return validation errors
order, err := client.NewCreateOrderService().
    Symbol("BTCUSDT").
    // Missing required parameters
    Do(context.Background())

if err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```

## Risk Management

### Position-Based Trading

```go
// Check position before placing order
positions, _ := client.NewAllPositionsService().
    ProductType(trading.ProductTypeUSDTFutures).
    Do(context.Background())

for _, pos := range positions {
    if pos.Symbol == "BTCUSDT" && pos.HoldSide == "long" {
        // Close position if unrealized loss > 5%
        unrealizedPL, _ := strconv.ParseFloat(pos.UnrealizedPL, 64)
        if unrealizedPL < -500 {
            client.NewCreateOrderService().
                Symbol("BTCUSDT").
                Side(trading.SideSell).
                OrderType(trading.OrderTypeMarket).
                Size(pos.Size).
                ReduceOnly(true).
                Do(context.Background())
        }
    }
}
```

### Order Size Validation

```go
// Get contract specifications for size validation
contracts, err := client.NewContractsService().
    ProductType(trading.ProductTypeUSDTFutures).
    Symbol("BTCUSDT").
    Do(context.Background())

for _, contract := range contracts {
    minSize, _ := strconv.ParseFloat(contract.MinTradeNum, 64)
    maxSize, _ := strconv.ParseFloat(contract.MaxTradeNum, 64)
    
    // Validate order size
    orderSize := 0.001
    if orderSize < minSize || orderSize > maxSize {
        fmt.Printf("Invalid order size for %s: %f (min: %f, max: %f)\n",
            contract.Symbol, orderSize, minSize, maxSize)
    }
}
```

## Testing

Run trading service tests:

```bash
# Run all trading service tests
go test ./futures/trading/... -v

# Test specific services
go test ./futures/trading/... -run TestCreateOrderService -v
go test ./futures/trading/... -run TestBatchOrdersService -v
```

## Performance Tips

1. **Batch Operations**: Use batch orders for multiple simultaneous placements
2. **WebSocket Integration**: Subscribe to order updates for real-time status
3. **Efficient Polling**: Use appropriate intervals for order status checks
4. **Rate Limiting**: Respect API rate limits (2400 requests/minute for futures)

## WebSocket Integration

For real-time order updates:

```go
wsClient := ws.NewWsClient()
wsClient.SubscribePrivateChannel("orders", orderUpdateCallback)
```