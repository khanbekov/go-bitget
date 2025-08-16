# Testing Framework

This directory contains the comprehensive testing framework for the Bitget Go SDK, including unit tests, integration tests, and testing utilities.

## Directory Structure

```
tests/
├── README.md                    # This file - testing overview
├── configs/                     # Test configuration files
│   ├── integration.example.json # Example integration test config
│   └── integration.json         # Your actual config (create from example)
├── integration/                 # Integration test framework
│   ├── config.go               # Configuration management
│   ├── runner.go               # Test execution and reporting
│   └── suites/                 # Test suites by category
│       ├── account_test.go     # Account endpoints tests
│       ├── market_test.go      # Market data endpoints tests
│       └── trading_test.go     # Trading endpoints tests
├── scripts/                    # Test execution scripts
│   ├── run-integration-tests.sh  # Unix/Linux/macOS runner
│   └── run-integration-tests.bat # Windows runner
├── reports/                    # Generated test reports
│   ├── integration_report.json   # JSON test results
│   └── integration_report.html   # HTML dashboard
└── unit/                       # Unit test utilities (future)
    └── mocks/                  # Shared mock implementations
```

## Quick Start

Note: run the commands below in the `tests/` directory.

### 1. Setup Integration Testing

```bash
# Copy example configuration
cp configs/integration.example.json configs/integration.json

# Edit with your API credentials
# IMPORTANT: Set demo_trading: true for safety
```

### 2. Set Environment Variables

Create `.env` file in project root:
```bash
BITGET_API_KEY=your_api_key_here
BITGET_SECRET_KEY=your_secret_key_here
BITGET_PASSPHRASE=your_passphrase_here
BITGET_DEMO_TRADING=true
```

### 3. Run Integration Tests

```bash
# Easy method - using scripts
scripts/run-integration-tests.sh                    # Unix/Linux/macOS
scripts/run-integration-tests.bat                   # Windows

# Direct Go testing
go test -tags=integration ./tests/integration/suites -v

# Run specific test suite
go test -tags=integration ./tests/integration/suites -v -run TestAccountEndpoints
```

## Test Categories

### Integration Tests (Real API)

Tests against real Bitget API endpoints with your own credentials:

**Account Endpoints**:
- ✅ **account_info** - Get account balances (safe)
- ✅ **account_list** - Get all accounts (safe)
- ✅ **account_bills** - Get transaction history (safe)
- ⚠️ **set_leverage** - Change leverage (demo mode only)
- ⚠️ **set_margin_mode** - Change margin mode (demo mode only)
- ⚠️ **set_position_mode** - Change position mode (demo mode only)
- ⚠️ **adjust_margin** - Adjust position margin (demo mode only)

**Market Data Endpoints** (coming soon):
- Get tickers, candlesticks, order books
- Contract information and funding rates
- Open interest and symbol prices

**Trading Endpoints** (coming soon):
- Order placement, modification, cancellation
- Batch operations and plan orders
- Fill history and order details

### Unit Tests

Located alongside source code in respective packages:
- `futures/*_test.go` - Futures API unit tests
- `uta/*_test.go` - UTA API unit tests  
- `ws/*_test.go` - WebSocket unit tests
- `common/*_test.go` - Common utilities unit tests

## Safety Guidelines

⚠️ **CRITICAL SAFETY MEASURES**:

1. **Always Use Demo Mode**: Set `demo_trading: true` in configuration
2. **Start with Read-Only Tests**: Begin with safe endpoints
3. **Gradual Enablement**: Enable write operations one at a time
4. **Small Test Amounts**: Use minimal values for financial operations
5. **Review Configuration**: Always check settings before running

## Configuration Options

### Integration Test Config (`tests/configs/integration.json`)

```json
{
  "demo_trading": true,           // REQUIRED for safety
  "endpoint_tests": {
    "account_info": true,         // Safe - read only
    "account_list": true,         // Safe - read only
    "account_bills": true,        // Safe - read only
    "set_leverage": false,        // Disabled - has side effects
    "adjust_margin": false,       // Disabled - has side effects
    "set_margin_mode": false,     // Disabled - has side effects
    "set_position_mode": false    // Disabled - has side effects
  },
  "test_symbol": "BTCUSDT",
  "test_product_type": "USDT-FUTURES",
  "generate_report": true,
  "report_path": "tests/reports/integration_report.json"
}
```

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `BITGET_API_KEY` | Your API key | Yes |
| `BITGET_SECRET_KEY` | Your secret key | Yes |
| `BITGET_PASSPHRASE` | Your passphrase | Yes |
| `BITGET_DEMO_TRADING` | Enable demo mode (`true`/`false`) | Recommended |
| `INTEGRATION_CONFIG_FILE` | Custom config file path | Optional |

## Test Reports

Integration tests generate comprehensive reports:

### JSON Report (`tests/reports/integration_report.json`)
Machine-readable test results with detailed metrics:
- Test execution times and success rates
- Error details and stack traces
- Environment and configuration information
- Individual test results and responses

### HTML Report (`tests/reports/integration_report.html`)
Human-readable dashboard with:
- Visual test result summary
- Success/failure charts
- Environment information
- Detailed test breakdown with error messages

### Console Output
Real-time structured logging during test execution:
```
INFO [2024-01-15 10:30:00] Starting account integration tests base_url=https://api.bitget.com demo_trading=true
INFO [2024-01-15 10:30:01] Testing account info endpoint test=account_info
INFO [2024-01-15 10:30:02] Test passed test=account_info account_id=123456 usdt_equity=1000.50
```

## Advanced Usage

### Custom Configuration

```bash
# Use custom config file
tests/scripts/run-integration-tests.sh -c tests/configs/my_config.json

# Run specific test suite with timeout
tests/scripts/run-integration-tests.sh -s account -t 10m

# Check demo trading configuration
tests/scripts/run-integration-tests.sh -d

# Quiet mode with no reports
tests/scripts/run-integration-tests.sh -q -r
```

### Direct Go Commands

```bash
# Run all integration tests
go test -tags=integration ./tests/integration/suites -v

# Run with custom timeout
go test -tags=integration ./tests/integration/suites -v -timeout=10m

# Run specific test
go test -tags=integration ./tests/integration/suites -v -run TestAccountInfo

# Run with race detection
go test -tags=integration ./tests/integration/suites -v -race
```

## Troubleshooting

### Common Issues

**Authentication Errors**:
- Verify API credentials in `.env` or config file
- Check API key permissions and account status
- Ensure correct environment (testnet vs production)

**Demo Mode Limitations**:
- Some endpoints may return errors in demo mode (this is normal)
- Focus on testing connectivity and data structure
- Write operations may have limited functionality

**Configuration Issues**:
- Check JSON syntax in config files
- Verify file paths and permissions
- Ensure required fields are set

**Network Issues**:
- Check internet connectivity
- Verify firewall settings
- Test API accessibility via curl/browser

### Debug Mode

Enable detailed logging by setting log level to `debug` in configuration:
```json
{
  "log_level": "debug"
}
```

### Getting Help

1. Check [`INTEGRATION_TESTING.md`](INTEGRATION_TESTING.md) for detailed documentation
2. Review test logs for specific error messages
3. Verify API credentials and permissions
4. Ensure account is in correct mode (UTA vs Classic)
5. Check Bitget API status and maintenance schedules

## Best Practices

1. **Start Small**: Begin with read-only endpoints
2. **Use Version Control**: Track configuration changes
3. **Monitor API Usage**: Keep track of API call limits
4. **Regular Testing**: Run tests after code changes
5. **Review Reports**: Analyze results before enabling more tests
6. **Keep Credentials Secure**: Never commit API keys to version control
7. **Document Changes**: Update test configurations as APIs evolve

## Contributing

When adding new tests:

1. **Follow Existing Patterns**: Use the same structure as existing test suites
2. **Add Safety Checks**: Ensure write operations are demo-mode only
3. **Update Configuration**: Add new endpoints to config templates
4. **Document Behavior**: Include expected responses and error conditions
5. **Test Thoroughly**: Verify both success and failure scenarios

For detailed development guidelines, see [`CONTRIBUTING.md`](../CONTRIBUTING.md).