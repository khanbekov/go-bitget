# Development Guide  

This comprehensive guide covers development practices, testing, contribution guidelines, and project architecture for the Bitget Go SDK.

## Development Setup

### Prerequisites

- **Go 1.23.4 or later** - Required for module support and latest features
- **Git** - For version control
- **Make** (optional) - For build automation

### Environment Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/khanbekov/go-bitget.git
   cd go-bitget
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Set up environment variables** (create `.env` file):
   ```bash
   BITGET_API_KEY=your_api_key_here
   BITGET_SECRET_KEY=your_secret_key_here
   BITGET_PASSPHRASE=your_passphrase_here
   BITGET_DEMO_TRADING=true
   ```

4. **Verify setup**:
   ```bash
   go test ./...
   go build .
   ```

## Testing

### Unit Testing

```bash
go test ./...           # Run all unit tests
go test -cover ./...    # Run with coverage
go test -race ./...     # Run with race detection

# Test specific packages
go test -v ./futures/
go test -v ./uta/
go test -v ./ws/
go test -v ./common/
```

### Integration Testing (Real API)

The repository includes comprehensive integration tests that verify functionality against the real Bitget API using your own credentials.

#### Quick Start

1. **Setup credentials**:
   ```bash
   # Create .env file
   BITGET_API_KEY=your_api_key
   BITGET_SECRET_KEY=your_secret_key  
   BITGET_PASSPHRASE=your_passphrase
   BITGET_DEMO_TRADING=true
   ```

2. **Configure tests**:
   ```bash
   cp tests/configs/integration.example.json tests/configs/integration.json
   # Edit tests/configs/integration.json with your preferences
   ```

3. **Run integration tests**:
   ```bash
   # Easy method
   tests/scripts/run-integration-tests.sh                    # Unix/Linux/macOS
   tests/scripts/run-integration-tests.bat                   # Windows
   
   # Direct Go testing
   go test -tags=integration ./tests/integration/suites -v
   ```

#### Available Test Suites

**Account Endpoints** (Recommended starting point):
- ✅ account_info - Get account balances (safe)
- ✅ account_list - Get all accounts (safe)
- ✅ account_bills - Get transaction history (safe)
- ⚠️ set_leverage - Change leverage (demo mode only)
- ⚠️ set_margin_mode - Change margin mode (demo mode only)
- ⚠️ set_position_mode - Change position mode (demo mode only)

#### Safety Features

- **Demo Trading Mode**: Write operations only run in demo mode
- **Selective Testing**: Enable/disable individual endpoints
- **Comprehensive Reporting**: JSON and HTML test reports
- **Error Recovery**: Built-in retry logic and error handling
- **Configuration Validation**: Automatic safety checks

#### Test Reports

Integration tests generate detailed reports:
```bash
# Generated files
tests/reports/integration_report.json      # Machine-readable results
tests/reports/integration_report.html      # Human-readable dashboard
```

#### Advanced Usage

```bash
# Test specific suite
tests/scripts/run-integration-tests.sh -s account

# Use custom config
tests/scripts/run-integration-tests.sh -c tests/configs/my_config.json

# Check safety configuration
tests/scripts/run-integration-tests.sh -d

# Run with custom timeout
go test -tags=integration ./tests/integration/suites -v -timeout=10m
```

For comprehensive integration testing documentation, see [`tests/INTEGRATION_TESTING.md`](../tests/INTEGRATION_TESTING.md).

## Common Development Commands

### Building & Running

```bash
# Build the application
go build -o app .

# Run directly
go run main.go

# Run with .env file (requires API credentials)
go run main.go
```

### Module Management

```bash
# Update dependencies
go mod tidy

# Download dependencies
go mod download

# Verify dependencies
go mod verify

# View dependency graph
go mod graph
```

### Documentation Generation

```bash
# Generate comprehensive documentation (creates docs/ directory)
bash generate-docs.sh    # Unix/Linux/macOS
generate-docs.bat         # Windows

# Start live documentation server
godoc -http=:6060
# Visit: http://localhost:6060/pkg/github.com/khanbekov/go-bitget/
```

## Project Architecture

### Package Structure

The SDK uses a service-oriented architecture with 4 main packages:

- **`futures/`** - Legacy futures API organized into subdirectories:
  - `account/` - Account management (balances, positions, leverage, margin)
  - `market/` - Market data (tickers, candles, order books, contracts)
  - `position/` - Position operations (get, close, history)  
  - `trading/` - Trading operations (orders, batch operations, plan orders)
- **`uta/`** - Unified Trading Account API (recommended for new development)
- **`ws/`** - Unified WebSocket implementation with production-ready features
- **`common/`** - Shared utilities, authentication, error handling, constants
- **`tests/`** - Comprehensive testing framework with integration tests

### API Design Patterns

#### Service-Oriented with Fluent API

All services support method chaining for intuitive usage:

```go
// Futures API pattern
client := futures.NewClient(apiKey, secretKey, passphrase)
result, err := client.NewCandlestickService().
    Symbol("BTCUSDT").
    ProductType(futures.ProductTypeUSDTFutures).
    Granularity("1m").
    Limit("100").
    Do(context.Background())

// UTA API pattern (recommended)
utaClient := uta.NewClient(apiKey, secretKey, passphrase)
order, err := utaClient.NewPlaceOrderService().
    Symbol("BTCUSDT").
    Side("buy").
    OrderType("limit").
    Do(context.Background())
```

#### Context Support

All operations accept `context.Context` for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := client.NewAccountInfoService().Do(ctx)
```

#### Structured Error Handling

Use structured error types from `common/types`:

```go
if apiErr, ok := err.(*common.APIError); ok {
    // Handle API-specific errors with codes
    fmt.Printf("API Error: %s (Code: %s)\n", apiErr.Message, apiErr.Code)
    switch apiErr.Code {
    case "40001":
        // Handle invalid parameters
    case "40004":
        // Handle insufficient balance
    }
} else {
    // Handle network/connection errors
    fmt.Printf("Network Error: %v\n", err)
}
```

### WebSocket Implementation

The unified WebSocket client (`ws/`) provides enterprise-grade features:

- **Rate Limiting**: 10 messages/second automatic throttling
- **Automatic Reconnection**: Configurable timeout and retry logic
- **Subscription Restoration**: Maintains subscriptions after reconnection
- **Health Monitoring**: Connection health checks and heartbeat mechanism
- **Thread-Safe**: Concurrent subscription management
- **Type Safety**: Structured data types for all message formats

```go
// WebSocket client initialization
logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
wsClient := ws.NewBitgetBaseWsClient(logger, endpoint, secretKey)

// Subscribe with automatic reconnection
wsClient.SubscribeTicker("BTCUSDT", "USDT-FUTURES", func(msg string) {
    // Handle ticker updates
})
```

## Key Implementation Details

### Client Initialization

```go
// Futures client (legacy)
client := futures.NewClient(apiKey, secretKey, passphrase)

// UTA client (recommended) - auto-detects demo mode
utaClient := uta.NewClient(apiKey, secretKey, passphrase)
utaClient.SetDemoTrading(false) // explicit control

// WebSocket client
wsClient := ws.NewBitgetBaseWsClient(logger, endpoint, secretKey)
```

### Environment Configuration

The main.go and tests expect these environment variables:

| Variable | Description | Required |
|----------|-------------|----------|
| `BITGET_API_KEY` | Your Bitget API key | Yes |
| `BITGET_SECRET_KEY` | Your Bitget secret key | Yes |
| `BITGET_PASSPHRASE` | Your Bitget passphrase | Yes |
| `BITGET_DEMO_TRADING` | Enable demo mode (`true`/`false`) | Recommended |
| `BITGET_TESTNET` | Use testnet (`true`/`false`) | Optional |

### Product Types & Constants

Key constants defined across packages:

**Futures Package**:
- `futures.ProductTypeUSDTFutures`, `futures.ProductTypeCoinFutures`
- `futures.MarginModeIsolated`, `futures.MarginModeCrossed`
- `futures.PositionModeOneWay`, `futures.PositionModeHedge`

**UTA Package**:
- `uta.CategorySpot`, `uta.CategoryUSDTFutures`, `uta.CategoryMargin`
- `uta.SideBuy`, `uta.SideSell`
- `uta.OrderTypeLimit`, `uta.OrderTypeMarket`

**WebSocket Package**:
- Channels and timeframes in `ws/channels.go`
- Message types in `ws/types.go`

## Development Guidelines

### Version Status

This is **v0.0.1 (Alpha)** - API is not stable and will change without backward compatibility until v1.0.0.

### Code Style Standards

- **Follow Go Conventions**: Use `gofmt`, `golint`, and `go vet`
- **Method Chaining**: Support fluent API patterns for all services
- **GoDoc Comments**: All public APIs require comprehensive documentation
- **Structured Logging**: Use zerolog for consistent logging
- **Error Handling**: Implement proper error handling with custom error types
- **Context Support**: All operations must support context cancellation
- **Testing**: Write comprehensive unit tests and integration tests

### Dependencies

Core dependencies managed in `go.mod`:

- **HTTP Client**: `github.com/valyala/fasthttp` - High-performance HTTP client
- **JSON Processing**: `github.com/json-iterator/go` - Fast JSON serialization
- **Logging**: `github.com/rs/zerolog` - Structured logging framework
- **WebSocket**: `github.com/gorilla/websocket` - WebSocket implementation
- **Testing**: `github.com/stretchr/testify` - Testing framework with assertions
- **Environment**: `github.com/joho/godotenv` - Environment variable loading

### Common Development Workflows

1. **Use UTA API for New Development**: Recommended over legacy futures API
2. **WebSocket Auto-Reconnection**: Don't manually reconnect, let the client handle it
3. **Error Type Checking**: Always distinguish between API errors and network errors
4. **Demo Mode Testing**: Use demo trading mode detection for safer testing
5. **Documentation Generation**: Run docs generation after significant changes
6. **Integration Testing**: Test against real API with proper safety measures

### Testing Architecture

- **Unit Tests**: Mock clients in each package (`*_test.go` files)
- **Integration Tests**: Real API testing framework in `tests/integration/`
- **Mock Patterns**: Follow testify patterns with expectations
- **WebSocket Testing**: Cover subscription management and channel handling
- **Safety First**: All write operations use demo trading mode by default

### File Organization

**Key Files to Check for Examples**:
- `main.go` - Complete UTA workflow with error handling
- `examples/` - Various usage patterns and demo applications  
- `*_test.go` files - Service usage patterns and mock client setup
- `tests/` directory - Comprehensive testing framework
- `tests/README.md` - Testing framework overview
- `tests/INTEGRATION_TESTING.md` - Detailed integration testing guide

## Contribution Guidelines

See [`CONTRIBUTING.md`](../CONTRIBUTING.md) for detailed contribution guidelines, including:

- Code style standards and review process
- Testing requirements and integration testing
- Pull request guidelines and branching strategy
- Issue reporting and feature request process
