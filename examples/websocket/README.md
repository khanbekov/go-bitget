# WebSocket Examples for Bitget Go SDK

This directory contains comprehensive examples demonstrating various WebSocket usage patterns with the Bitget Go SDK. These examples cover everything from basic subscriptions to advanced trading scenarios and error handling patterns.

## 📁 Available Examples

### 1. **basic_public_channels.go** - Basic Public Channel Subscriptions
**Purpose**: Introduction to public WebSocket channels  
**Features**:
- ✅ Ticker updates (24hr statistics)
- ✅ 1-minute candlestick data  
- ✅ Top 5 order book levels
- ✅ Real-time trade executions
- ✅ Mark price updates
- ✅ Funding rate information
- ✅ Graceful shutdown handling

**Run**: `go run basic_public_channels.go`

---

### 2. **private_channels.go** - Private Channel Authentication & Subscriptions  
**Purpose**: Demonstrate private channel authentication and account updates  
**Features**:
- 🔐 API key authentication
- 📋 Order status updates
- ✅ Trade fill notifications  
- 📊 Position updates
- 💰 Account balance changes
- ⚡ Plan order (trigger order) updates
- 🔍 Connection health monitoring

**Requirements**: Valid API credentials in `.env` file  
**Run**: `go run private_channels.go`

---

### 3. **advanced_usage.go** - Production-Grade WebSocket Management
**Purpose**: Advanced patterns for high-performance trading applications  
**Features**:
- 🏗️ Structured data processing with goroutines
- 🔄 Automatic reconnection with exponential backoff  
- 📈 Real-time performance statistics
- 🧵 Thread-safe data management
- 📊 Market data aggregation
- 🔍 Connection health monitoring
- ⚡ High-frequency message processing
- 📱 Live dashboard with market summary

**Run**: `go run advanced_usage.go`

---

### 4. **comprehensive_channels_demo.go** - Complete Channel Coverage
**Purpose**: Demonstrate ALL available WebSocket channels  
**Features**:
- 📊 **ALL Ticker Types**: USDT-FUTURES, COIN-FUTURES, USDC-FUTURES
- 🕯️ **ALL Candlestick Timeframes**: 1m, 5m, 15m, 30m, 1h, 4h, 6h, 12h, 1d, 3d, 1w, 1M
- 📚 **ALL Order Book Depths**: books5, books15, books (full depth)  
- 💰 **ALL Trade Channels**: Real-time executions across multiple symbols
- 🎯 **Mark Price & Funding**: For futures contracts
- 🔒 **ALL Private Channels**: Orders, fills, positions, account, plan orders
- 📈 **Live Statistics**: Message counts, subscription health, connection status
- 🔍 **Multi-Product Support**: USDT-FUTURES, COIN-FUTURES, USDC-FUTURES

**Run**: `go run comprehensive_channels_demo.go`

---

### 5. **trading_bot_scenarios.go** - Real-World Trading Scenarios  
**Purpose**: Practical trading bot implementations using WebSocket + REST API  
**Features**:
- 🤖 **4 Trading Scenarios**:
  - 📊 **Grid Trading Bot**: Automated grid orders based on price levels
  - 💰 **DCA Bot**: Dollar-cost averaging with scheduled buys
  - ⚡ **Scalping Bot**: High-frequency trading on spread opportunities  
  - 🛡️ **Risk Manager**: Real-time position and exposure monitoring
- 📱 **Live Trading Dashboard**: Real-time P&L, positions, market data
- 📊 **Market Data Integration**: Price feeds, order books, mark prices
- 💼 **Simulated Trading**: Safe demonstration with simulated orders
- 📈 **Performance Metrics**: Trade counts, profitability, runtime stats

**Run**: `go run trading_bot_scenarios.go`

---

### 6. **error_handling_patterns.go** - Production Error Handling
**Purpose**: Comprehensive error handling and recovery patterns  
**Features**:
- 🛡️ **Error Classification**: 7 error types with specific handling
- 🔄 **Recovery Strategies**: Exponential backoff, linear retry, fixed delay
- 🔌 **Auto-Reconnection**: Intelligent reconnection with state management
- 📊 **Error Analytics**: Pattern detection, frequency analysis
- 🔍 **Health Monitoring**: Connection, subscription, and data staleness checks
- 📈 **Error Dashboard**: Real-time error statistics and recent events
- 🎭 **Error Simulation**: Demonstration of recovery mechanisms
- 📱 **State Management**: Connection state tracking and transitions

**Run**: `go run error_handling_patterns.go`

---

### 7. **mixed_channels.go** - Combined Public & Private Channels
**Purpose**: Demonstrate simultaneous public and private channel usage  
**Features**: 
- 🌍 Public market data subscriptions
- 🔒 Private account update subscriptions  
- 🔄 Dual connection management
- 📊 Combined data processing

**Run**: `go run mixed_channels.go`

---

### 8. **multiple_symbols.go** - Multi-Symbol Market Data
**Purpose**: Efficient handling of multiple trading pairs  
**Features**:
- 📈 Batch symbol subscriptions
- 🔄 Parallel data processing
- 📊 Cross-symbol analysis
- ⚡ Optimized message handling

**Run**: `go run multiple_symbols.go`

## 🚀 Quick Start

### Prerequisites

1. **Environment Setup**:
   ```bash
   # Copy environment template
   cp .env.example .env
   
   # Edit with your API credentials
   BITGET_API_KEY=your_api_key
   BITGET_SECRET_KEY=your_secret_key  
   BITGET_PASSPHRASE=your_passphrase
   ```

2. **Dependencies**:
   ```bash
   go mod tidy
   ```

### Running Examples

#### Option 1: Run Individual Examples
```bash
# Basic public channels
go run basic_public_channels.go

# Private channels (requires valid API keys)
go run private_channels.go

# Advanced usage patterns  
go run advanced_usage.go

# Comprehensive channel demo
go run comprehensive_channels_demo.go

# Trading bot scenarios
go run trading_bot_scenarios.go

# Error handling patterns
go run error_handling_patterns.go
```

#### Option 2: Use Test Scripts
```bash
# On Unix/Linux/macOS
bash test_examples.sh

# On Windows
test_examples.bat
```

## 📋 Channel Reference

### Public Channels (No Authentication Required)

| Channel | Description | Symbols | Example |
|---------|-------------|---------|---------|
| `ticker` | 24hr price statistics | All symbols | Price, volume, change |
| `candle{timeframe}` | OHLCV candlestick data | All symbols | `candle1m`, `candle1h` |  
| `books` | Full order book depth | All symbols | Complete bid/ask levels |
| `books5` | Top 5 order book levels | All symbols | Best 5 bids/asks |
| `books15` | Top 15 order book levels | All symbols | Best 15 bids/asks |
| `trade` | Real-time trade executions | All symbols | Price, size, side, timestamp |
| `mark-price` | Mark price updates | Futures only | Contract mark price |
| `funding-time` | Funding rate & time | Perpetual futures | Rate, next funding time |

### Private Channels (Authentication Required)

| Channel | Description | Product Types | Example |
|---------|-------------|---------------|---------|
| `orders` | Order status updates | All | New, filled, canceled orders |
| `fill` | Trade execution updates | All | Fill price, quantity, fees |
| `positions` | Position updates | Futures/Margin | Size, PnL, margin changes |  
| `account` | Balance updates | All | Available balance changes |
| `plan-order` | Trigger order updates | All | Stop-loss, take-profit orders |

### Supported Timeframes

| Timeframe | Constant | Description |
|-----------|----------|-------------|
| `1m` | `ws.Timeframe1m` | 1 minute |
| `5m` | `ws.Timeframe5m` | 5 minutes |
| `15m` | `ws.Timeframe15m` | 15 minutes |
| `30m` | `ws.Timeframe30m` | 30 minutes |
| `1h` | `ws.Timeframe1h` | 1 hour |
| `4h` | `ws.Timeframe4h` | 4 hours |
| `6h` | `ws.Timeframe6h` | 6 hours |
| `12h` | `ws.Timeframe12h` | 12 hours |
| `1d` | `ws.Timeframe1d` | 1 day |
| `3d` | `ws.Timeframe3d` | 3 days |
| `1w` | `ws.Timeframe1w` | 1 week |
| `1M` | `ws.Timeframe1M` | 1 month |

### Product Types

| Product Type | Description | Examples |
|-------------|-------------|----------|
| `USDT-FUTURES` | USDT-margined perpetual | `BTCUSDT`, `ETHUSDT` |
| `COIN-FUTURES` | Coin-margined perpetual | `BTCUSD`, `ETHUSD` |  
| `USDC-FUTURES` | USDC-margined perpetual | `BTCUSDC`, `ETHUSDC` |
| `SPOT` | Spot trading | `BTCUSDT`, `ETHUSDT` |

## 🔧 Configuration Examples

### Basic WebSocket Setup
```go
client := ws.NewBitgetBaseWsClient(
    logger,
    "wss://ws.bitget.com/v2/ws/public",
    "", // No secret for public
)

client.SetListener(messageHandler, errorHandler)
client.Connect()
client.ConnectWebSocket() 
client.StartReadLoop()
```

### High-Level Futures Integration  
```go
futuresClient := futures.NewClient(apiKey, secretKey, passphrase)
wsManager := futuresClient.NewWebSocketManager()

// Connect and subscribe  
wsManager.ConnectPublic()
wsManager.SubscribeToTicker("BTCUSDT", tickerHandler)
wsManager.SubscribeToOrderBook("BTCUSDT", 5, orderBookHandler)
```

### Advanced Error Handling
```go
// Custom reconnection settings
client.SetReconnectionTimeout(30 * time.Second)
client.SetCheckConnectionInterval(5 * time.Second)

// Enhanced error handler
errorHandler := func(message string) {
    logger.Error().Str("error", message).Msg("WebSocket Error")
    
    // Analyze error type and trigger recovery
    if shouldReconnect(message) {
        go reconnectWithBackoff()
    }
}
```

## 📊 Performance Considerations

### Message Processing
- **Buffer Channels**: Use buffered channels for high-frequency data
- **Goroutine Pools**: Process messages in separate goroutines
- **Data Validation**: Validate JSON before processing
- **Rate Limiting**: Respect API rate limits

### Connection Management  
- **Connection Pooling**: Reuse connections when possible
- **Heartbeat Monitoring**: Implement ping/pong for health checks
- **Graceful Shutdown**: Properly close connections and unsubscribe
- **State Management**: Track connection and authentication states

### Error Handling
- **Retry Logic**: Implement exponential backoff for retries
- **Circuit Breakers**: Prevent cascade failures
- **Monitoring**: Track error rates and patterns
- **Alerting**: Set up notifications for critical errors

## 🔍 Troubleshooting

### Common Issues

**Connection Failures**:
```bash  
# Check network connectivity
curl -I https://api.bitget.com

# Verify WebSocket endpoint
wscat -c wss://ws.bitget.com/v2/ws/public
```

**Authentication Errors**:
- Verify API key permissions
- Check passphrase accuracy  
- Ensure correct signature generation
- Validate timestamp synchronization

**Subscription Issues**:
- Check symbol format (case sensitive)
- Verify product type spelling
- Ensure connection before subscribing
- Monitor rate limits

**Data Processing Errors**:
- Validate JSON structure
- Handle malformed messages gracefully
- Implement data type validation
- Add timeout handling

### Debug Mode
Enable debug logging for detailed information:
```go
logger := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
```

### Performance Monitoring
Monitor key metrics:
- Message rates per channel
- Connection uptime
- Error frequencies  
- Memory usage
- Goroutine counts

## 📚 Additional Resources

- [Bitget WebSocket API Documentation](https://bitgetlimited.github.io/apidoc/en/mix/#websocket)
- [Go SDK Documentation](../../README.md)
- [Futures Package Documentation](../../futures/README.md)
- [UTA Package Documentation](../../uta/README.md)
- [Common WebSocket Patterns](../../ws/README.md)

## 🤝 Contributing

When adding new examples:

1. **Follow Naming Convention**: `{purpose}_{type}.go`
2. **Add Documentation**: Include comprehensive comments
3. **Error Handling**: Implement proper error handling
4. **Testing**: Verify examples work with demo accounts
5. **Update README**: Add entry to this documentation

## ⚠️ Important Notes

### Security
- **Never commit API keys** to version control
- **Use environment variables** for credentials  
- **Enable demo mode** for testing when possible
- **Validate all input data** before processing

### Rate Limits
- **Public channels**: No authentication, shared rate limits
- **Private channels**: Requires authentication, account-specific limits  
- **Subscription limits**: Maximum subscriptions per connection
- **Message rates**: Vary by channel and symbol

### Demo vs Production
- **Demo accounts**: Limited functionality for some private channels
- **Testnet**: Use testnet endpoints for development
- **Production**: Only use with verified, funded accounts

---

🚀 **Happy trading with Bitget WebSocket examples!**