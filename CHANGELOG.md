# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [v0.0.1] - 2025-01-31

### ⚠️ ALPHA RELEASE WARNING
This is an alpha release. The API is **NOT STABLE** and **WILL CHANGE** without backward compatibility guarantees until v1.0.0.

### 🚀 Initial Alpha Release

This is the first public release of the Go Bitget SDK with comprehensive REST API and WebSocket support for futures trading.

### Added

#### Core Trading Features
- **Order Management**: Complete order lifecycle management
  - `CreateOrderService` - Create single orders (market, limit, stop)
  - `CreateBatchOrdersService` - Create up to 20 orders in a single request
  - `ModifyOrderService` - Modify existing orders
  - `CancelOrderService` - Cancel single orders
  - `CancelAllOrdersService` - Cancel all orders for a symbol
  - `GetOrderDetailsService` - Get detailed order information
  - `PendingOrdersService` - Get all pending orders
  - `OrderHistoryService` - Get order history with pagination
  - `FillHistoryService` - Get fill/trade history

#### Advanced Order Management (Plan Orders)
- **Plan Order Services**: Trigger/conditional orders support
  - `CreatePlanOrderService` - Create trigger orders (stop-loss, take-profit, normal, track, stop-surplus)
  - `ModifyPlanOrderService` - Modify existing plan orders
  - `CancelPlanOrderService` - Cancel plan orders
  - `PendingPlanOrdersService` - Get pending plan orders
  - Support for `PlanType` (normal_plan, track_plan, stop_loss, take_profit, stop_surplus)
  - Support for `TriggerType` (fill_price, mark_price)

#### Position Management
- **Position Services**: Complete position lifecycle management
  - `AllPositionsService` - Get all open positions
  - `SinglePositionService` - Get single position details
  - `HistoryPositionsService` - Get position history
  - `ClosePositionService` - Close positions

#### Account Management
- **Account Services**: Account information and configuration
  - `AccountService` - Get account balance and information
  - `GetAccountBillService` - Get account bill/transaction history
  - `SetLeverageService` - Set trading leverage

#### Account Configuration
- **Configuration Services**: Advanced account settings
  - `SetMarginModeService` - Set margin mode (isolated/cross)
  - `SetPositionModeService` - Set position mode (one-way/hedge)
  - `AccountListService` - Get list of all futures accounts
  - `AdjustMarginService` - Add or reduce margin for positions

#### Market Data
- **Basic Market Data**: Real-time and historical market information
  - `CandlestickService` - Get OHLCV candlestick data
  - `AllTickersService` - Get all tickers (24hr statistics)
  - `TickerService` - Get single ticker
  - `OrderBookService` - Get order book depth
  - `RecentTradesService` - Get recent trades
  - `ContractsService` - Get contract specifications

#### Advanced Market Data
- **Advanced Market Services**: Professional trading data
  - `CurrentFundingRateService` - Get current funding rates
  - `HistoryFundingRateService` - Get historical funding rates with pagination
  - `OpenInterestService` - Get open interest data
  - `SymbolPriceService` - Get mark, index, and last prices

#### WebSocket Real-time Data
- **Public Channels** (8 channels, no authentication required):
  - `ticker` - Real-time ticker updates
  - `candle{timeframe}` - Real-time candlesticks (12 timeframes)
  - `books` - Full order book depth
  - `books5` - Top 5 order book levels
  - `books15` - Top 15 order book levels
  - `trade` - Real-time trade executions
  - `mark-price` - Mark price updates
  - `funding-time` - Funding rate and time

- **Private Channels** (5 channels, authentication required):
  - `orders` - Real-time order updates
  - `fill` - Real-time fill/execution updates
  - `positions` - Real-time position updates
  - `account` - Account balance updates
  - `plan-order` - Trigger order updates

#### Type System
- **Comprehensive Type Definitions**:
  - `ProductType` - Product type constants (USDT-FUTURES, COIN-FUTURES, USDC-FUTURES)
  - `MarginModeType` - Margin mode constants (ISOLATED, CROSSED)
  - `PositionModeType` - Position mode constants (one_way, hedge)
  - `SideType` - Order side constants (BUY, SELL)
  - `OrderType` - Order type constants (MARKET, LIMIT)
  - `TimeInForceType` - Time in force constants (GTC, IOC, FOK, post_only)
  - `PlanType` - Plan order type constants
  - `TriggerType` - Trigger type constants

#### Client Architecture
- **Service-Oriented Design**: Each API endpoint group has its own service
- **Fluent API Pattern**: Method chaining for intuitive usage
- **Context Support**: All operations support `context.Context`
- **Client Constructors**: All 37 services have client constructor methods

#### Error Handling
- **Comprehensive Error Handling**:
  - Structured error types for API errors
  - Network error handling with retry logic
  - Context cancellation support
  - Detailed error messages and codes

#### Testing
- **Extensive Test Coverage**:
  - 34 test files with comprehensive unit tests
  - 89.6% code coverage
  - Mock client implementation for testing
  - Integration-style tests
  - All services include fluent API tests, success tests, error tests, and client integration tests

#### Examples and Documentation
- **WebSocket Examples**:
  - `basic_public_channels.go` - Basic public channel usage
  - `multiple_symbols.go` - Multiple symbol monitoring
  - `private_channels.go` - Private channels with authentication
  - `advanced_usage.go` - Advanced usage patterns
  - `mixed_channels.go` - Mixed public/private channels

- **Documentation**:
  - Complete API documentation
  - WebSocket channel reference
  - Architecture guide
  - Configuration instructions
  - Development guidelines

### Technical Details

#### Dependencies
- `github.com/valyala/fasthttp` - High-performance HTTP client
- `github.com/json-iterator/go` - Fast JSON processing
- `github.com/rs/zerolog` - Structured logging
- `github.com/gorilla/websocket` - WebSocket implementation
- `github.com/stretchr/testify` - Testing framework
- `github.com/joho/godotenv` - Environment variable loading
- `github.com/robfig/cron/v3` - Cron job scheduling
- `github.com/google/uuid` - UUID generation

#### Supported Environments
- Go 1.23.4+
- All major operating systems (Windows, macOS, Linux)
- Production API endpoint: `https://api.bitget.com`
- WebSocket endpoints: `wss://ws.bitget.com/v2/ws/public` and `wss://ws.bitget.com/v2/ws/private`

#### Security Features
- HMAC-SHA256 request signing
- API key authentication
- WebSocket authentication for private channels
- Secure credential handling
- No hardcoded secrets

## Statistics

- **Total REST API Services**: 34 services (futures trading complete)
- **WebSocket Channels**: 13 channels (8 public + 5 private)
- **Test Files**: 34+ test files
- **Documentation Files**: Complete documentation suite
- **Go Version**: 1.23.4+

## Alpha Release Notes

- **Stability**: ❌ Not suitable for production use
- **API Changes**: ⚠️ Breaking changes expected
- **Testing**: Limited production testing
- **Support**: Community-driven, use at your own risk
- **Backward Compatibility**: ❌ No guarantees until v1.0.0

## Roadmap to v1.0.0

- [ ] Production testing and stability improvements
- [ ] API stabilization and documentation refinement
- [ ] Performance optimizations
- [ ] Additional error handling and edge cases
- [ ] Community feedback integration
- [ ] Breaking change consolidation

## Contributors

- Initial implementation and complete API coverage
- Comprehensive testing and documentation
- WebSocket real-time data support
- Production-ready architecture

---

**Note**: This SDK is for educational and development purposes. Always test thoroughly in a sandbox environment before using with real funds.