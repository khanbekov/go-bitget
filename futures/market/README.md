# Market Data Services

This package contains services for retrieving real-time and historical market data from Bitget futures, including price data, funding rates, open interest, and recent trades.

## Services Overview

### Price & Quote Data

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `CandlestickService` | OHLCV candlestick data | `Symbol()`, `ProductType()`, `Granularity()`, `Limit()` |
| `AllTickersService` | 24hr ticker statistics for all symbols | `ProductType()` |
| `TickerService` | 24hr ticker statistics for specific symbol | `Symbol()`, `ProductType()` |
| `OrderBookService` | Order book depth data | `Symbol()`, `ProductType()`, `Limit()` |
| `SymbolPriceService` | Mark, index, and last prices | `Symbol()`, `ProductType()` |

### Trading Data

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `RecentTradesService` | Recent public trade executions | `Symbol()`, `ProductType()`, `Limit()` |
| `ContractsService` | Contract specifications and trading rules | `ProductType()`, `Symbol()` |

### Analytics Data

| Service | Description | Key Methods |
|---------|-------------|-------------|
| `CurrentFundingRateService` | Current funding rates | `Symbol()`, `ProductType()` |
| `HistoryFundingRateService` | Historical funding rates | `Symbol()`, `ProductType()`, `PageSize()` |
| `OpenInterestService` | Open interest data | `Symbol()`, `ProductType()` |

## Usage Examples

### Candlestick Data

```go
client := futures.NewClient(apiKey, secretKey, passphrase)

// Get 1-minute candlesticks for BTCUSDT
candles, err := client.NewCandlestickService().
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Granularity("1m").
    Limit("100").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, candle := range candles {
    fmt.Printf("Time: %d, Open: %f, High: %f, Low: %f, Close: %f, Volume: %f\n",
        candle.Ts, candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
}
```

### Market Ticker Data

```go
// Get ticker for specific symbol
ticker, err := client.NewTickerService().
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Last Price: %s\n", ticker.LastPr)
fmt.Printf("24h Change: %s%%\n", ticker.Change24h)
fmt.Printf("24h Volume: %s\n", ticker.BaseVolume)
```

### Order Book Data

```go
// Get order book depth
orderbook, err := client.NewOrderBookService().
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Limit("20"). // Top 20 levels
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Best Bid: %f @ %f\n", orderbook.Bids[0].Price, orderbook.Bids[0].Size)
fmt.Printf("Best Ask: %f @ %f\n", orderbook.Asks[0].Price, orderbook.Asks[0].Size)
```

### Funding Rates

```go
// Get current funding rate
fundingRates, err := client.NewCurrentFundingRateService().
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, rate := range fundingRates.FundingRates {
    fmt.Printf("Symbol: %s, Rate: %s, Next Funding: %s\n",
        rate.Symbol, rate.FundingRate, rate.FundingTime)
}
```

### Recent Trades

```go
// Get recent public trades
trades, err := client.NewRecentTradesService().
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Limit("50").
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, trade := range trades {
    fmt.Printf("Trade ID: %s, Price: %f, Size: %f, Side: %s, Time: %d\n",
        trade.TradeId, trade.Price, trade.Size, trade.Side, trade.Ts)
}
```

### Open Interest

```go
// Get open interest data
openInterest, err := client.NewOpenInterestService().
    Symbol("BTCUSDT").
    ProductType(market.ProductTypeUSDTFutures).
    Do(context.Background())

if err != nil {
    log.Fatal(err)
}

for _, oi := range openInterest.OpenInterests {
    fmt.Printf("Symbol: %s, Size: %s, Amount: %s, USDT Value: %s\n",
        oi.Symbol, oi.Size, oi.Amount, oi.OpenInterestUSDT)
}
```

## API Endpoints

This package covers the following Bitget API endpoints:

- `/api/v2/mix/market/candles` - Candlestick/OHLCV data
- `/api/v2/mix/market/tickers` - All symbol tickers
- `/api/v2/mix/market/ticker` - Single symbol ticker
- `/api/v2/mix/market/merge-depth` - Order book depth
- `/api/v2/mix/market/contracts` - Contract information
- `/api/v2/mix/market/fills` - Recent trades
- `/api/v2/mix/market/current-funding-rate` - Current funding rates
- `/api/v2/mix/market/history-funding-rate` - Historical funding rates
- `/api/v2/mix/market/open-interest` - Open interest data
- `/api/v2/mix/market/symbol-price` - Symbol prices (mark/index/last)

## Candlestick Granularities

Supported time intervals for candlestick data:

- `1m`, `3m`, `5m`, `15m`, `30m` - Minutes
- `1H`, `2H`, `4H`, `6H`, `8H`, `12H` - Hours  
- `1D`, `3D` - Days
- `1W` - Week
- `1M` - Month

## Data Types

### Market Data Structures

- **Candlestick**: OHLCV data with custom JSON unmarshaling
- **Ticker**: 24hr statistics including price, volume, change
- **OrderBook**: Bid/ask depth with price levels
- **RecentTrade**: Public trade executions
- **FundingRate**: Current and historical funding information
- **OpenInterest**: Outstanding contract positions

### Product Types

- `ProductTypeUSDTFutures` - USDT-margined futures
- `ProductTypeCoinFutures` - Coin-margined futures  
- `ProductTypeUSDCFutures` - USDC-margined futures

## Performance Notes

- **No Authentication Required**: All market data endpoints are public
- **Rate Limits**: 20 requests per 2 seconds per IP
- **WebSocket Alternative**: Consider WebSocket connections for real-time data
- **Caching**: Market data can be cached for short periods to reduce API calls

## Testing

Run market data service tests:

```bash
# Run all market service tests
go test ./futures/market/... -v

# Test specific services
go test ./futures/market/... -run TestCandlestickService -v
go test ./futures/market/... -run TestTickerService -v
```

## Error Handling

Market services include validation for required parameters:

```go
// This will return a validation error
candles, err := client.NewCandlestickService().
    // Missing required Symbol and ProductType
    Granularity("1m").
    Do(context.Background())

if err != nil {
    fmt.Printf("Validation error: %v\n", err)
}
```

## Real-time Data

For real-time market data updates, consider using the WebSocket API:

```go
// WebSocket example (see ws package)
wsClient := ws.NewWsClient()
wsClient.SubscribePublicChannel("ticker:BTCUSDT:USDT-FUTURES", callback)
```