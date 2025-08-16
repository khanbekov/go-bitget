# Integration Testing Guide

This guide covers how to run integration tests against the real Bitget API to verify all endpoints are working correctly with your API keys.

## ⚠️ Important Safety Notes

1. **Use Demo Trading Mode**: Always set `demo_trading: true` when testing
2. **Use Testnet When Available**: Set `use_testnet: true` 
3. **Start with Read-Only Tests**: Disable write operations initially
4. **Small Test Amounts**: Use minimal amounts for any write operations
5. **Review Configuration**: Double-check all settings before running

## Quick Start

### 1. Setup Environment

Create a `.env` file in the root directory:
```bash
BITGET_API_KEY=your_api_key_here
BITGET_SECRET_KEY=your_secret_key_here
BITGET_PASSPHRASE=your_passphrase_here
BITGET_DEMO_TRADING=true
```

### 2. Configure Tests

Copy the example configuration:
```bash
cp integration_config.example.json integration_config.json
```

Edit `integration_config.json` with your preferences:
```json
{
  "demo_trading": true,
  "endpoint_tests": {
    "account_info": true,
    "account_list": true,
    "account_bills": true,
    "set_leverage": false,
    "adjust_margin": false,
    "set_margin_mode": false,
    "set_position_mode": false
  }
}
```

### 3. Run Tests

Run integration tests with the integration build tag:

```bash
# Run all enabled account endpoint tests
go test -tags=integration ./integration -v

# Run with timeout
go test -tags=integration ./integration -v -timeout=5m

# Run specific tests
go test -tags=integration ./integration -v -run TestAccountInfo

# Generate detailed reports
go test -tags=integration ./integration -v -args -generate-report=true
```

## Configuration Options

### API Credentials
- `api_key`: Your Bitget API key
- `secret_key`: Your Bitget secret key  
- `passphrase`: Your Bitget passphrase

### Environment Settings
- `use_testnet`: Use testnet environment (recommended)
- `base_url`: API base URL
- `websocket_url`: WebSocket endpoint URL
- `demo_trading`: Enable demo trading mode (STRONGLY RECOMMENDED)

### Test Control
- `enabled_suites`: Which test suites to run (`["account", "market", "position"]`)
- `endpoint_tests`: Fine-grained control over individual endpoints
- `test_timeout_seconds`: Timeout for individual tests
- `max_retries`: Number of retries for failed requests

### Test Parameters
- `test_symbol`: Symbol to use for testing (default: "BTCUSDT")
- `test_product_type`: Product type (default: "USDT-FUTURES")
- `test_coin`: Coin for testing (default: "USDT")
- `safe_order_size`: Safe order size for write operations
- `safe_price_offset`: Price offset percentage for safe orders

### Reporting
- `generate_report`: Generate JSON and HTML reports
- `report_path`: Path for JSON report
- `log_level`: Logging level (`debug`, `info`, `warn`, `error`)

## Available Test Endpoints

### Account Endpoints

#### Safe (Read-Only) Tests
✅ **account_info** - Get account information and balances
- Endpoint: `/api/v2/mix/account/account`
- Safe to run anytime
- Tests basic API connectivity and authentication

✅ **account_list** - Get all accounts list  
- Endpoint: `/api/v2/mix/account/accounts`
- Safe to run anytime
- Returns list of all trading accounts

✅ **account_bills** - Get account transaction history
- Endpoint: `/api/v2/mix/account/bill`
- Safe to run anytime
- Returns transaction/bill history

#### Write Operations (Use with Caution)
⚠️ **set_leverage** - Set leverage for trading pair
- Endpoint: `/api/v2/mix/account/set-leverage`
- **HAS SIDE EFFECTS** - Changes account settings
- Only enable in demo trading mode
- Uses safe 20x leverage for testing

⚠️ **adjust_margin** - Adjust position margin
- Endpoint: `/api/v2/mix/account/set-margin`  
- **HAS SIDE EFFECTS** - Modifies position margins
- Requires existing positions
- Only enable in demo trading mode

⚠️ **set_margin_mode** - Set margin mode (isolated/cross)
- Endpoint: `/api/v2/mix/account/set-margin-mode`
- **HAS SIDE EFFECTS** - Changes margin mode
- Only enable in demo trading mode

⚠️ **set_position_mode** - Set position mode (one-way/hedge)
- Endpoint: `/api/v2/mix/account/set-position-mode`
- **HAS SIDE EFFECTS** - Changes position mode
- Only enable in demo trading mode

## Running Tests Safely

### Phase 1: Read-Only Tests
Start with safe, read-only operations:
```json
{
  "endpoint_tests": {
    "account_info": true,
    "account_list": true, 
    "account_bills": true,
    "set_leverage": false,
    "adjust_margin": false,
    "set_margin_mode": false,
    "set_position_mode": false
  }
}
```

### Phase 2: Write Operations (Demo Only)
After verifying read operations work, enable write operations in demo mode:
```json
{
  "demo_trading": true,
  "endpoint_tests": {
    "account_info": true,
    "account_list": true,
    "account_bills": true,
    "set_leverage": true,
    "set_margin_mode": true,
    "set_position_mode": true,
    "adjust_margin": false
  }
}
```

### Phase 3: Full Testing (Extreme Caution)
Only enable all tests if you understand the risks:
```json
{
  "demo_trading": true,
  "endpoint_tests": {
    "account_info": true,
    "account_list": true,
    "account_bills": true,
    "set_leverage": true,
    "adjust_margin": true,
    "set_margin_mode": true,
    "set_position_mode": true
  }
}
```

## Test Reports

### JSON Report
Detailed machine-readable test results:
```bash
cat integration_report.json | jq '.summary'
```

### HTML Report  
Human-readable test results with charts:
```bash
# Generate HTML report
go test -tags=integration ./integration -v -args -generate-html-report=true
# Open integration_report.html in browser
```

### Console Output
Real-time test execution with structured logging:
```
INFO [2024-01-15 10:30:00] Starting account integration tests base_url=https://api.bitget.com demo_trading=true
INFO [2024-01-15 10:30:01] Testing account info endpoint test=account_info
INFO [2024-01-15 10:30:02] Test passed test=account_info account_id=123456 usdt_equity=1000.50
INFO [2024-01-15 10:30:02] Account integration test summary total_tests=3 passed_tests=3 failed_tests=0
```

## Troubleshooting

### Common Issues

**Authentication Errors (40001)**
- Verify API keys are correct
- Check timestamp synchronization
- Ensure passphrase is correct

**Permission Errors (40004)**  
- Check API key permissions
- Verify futures trading is enabled
- Confirm account has sufficient balance

**Demo Mode Limitations**
- Some endpoints may not work in demo mode
- Error responses are normal for certain operations
- Focus on testing API connectivity and data structure

**Rate Limiting (30018)**
- Reduce test frequency
- Implement retry logic
- Check rate limit configuration

### Debug Mode
Enable debug logging for detailed request/response information:
```json
{
  "log_level": "debug"
}
```

### Test Individual Endpoints
Run specific endpoint tests:
```bash
go test -tags=integration ./integration -v -run TestAccountInfo
go test -tags=integration ./integration -v -run TestAccountList
```

## Best Practices

1. **Always Start with Demo Mode**: Never disable demo trading for initial testing
2. **Test Read Operations First**: Verify connectivity with safe operations
3. **Use Version Control**: Track configuration changes
4. **Monitor API Usage**: Keep track of API call limits
5. **Gradual Rollout**: Enable one endpoint at a time
6. **Review Reports**: Analyze test results before enabling more tests
7. **Keep Credentials Secure**: Never commit API keys to version control

## Environment Variables Reference

| Variable | Description | Required |
|----------|-------------|----------|
| `BITGET_API_KEY` | Your API key | Yes |
| `BITGET_SECRET_KEY` | Your secret key | Yes |
| `BITGET_PASSPHRASE` | Your passphrase | Yes |
| `BITGET_DEMO_TRADING` | Enable demo mode (`true`/`false`) | Recommended |
| `BITGET_TESTNET` | Use testnet (`true`/`false`) | Optional |
| `INTEGRATION_CONFIG_FILE` | Path to config file | Optional |

## Next Steps

After successfully running account endpoint tests:

1. **Market Data Tests**: Test market data endpoints
2. **Position Tests**: Test position management
3. **Trading Tests**: Test order placement (extreme caution)
4. **WebSocket Tests**: Test real-time data streams
5. **Comprehensive Integration**: Full end-to-end workflows

Remember: Integration testing is about verifying API connectivity and data structure, not about making actual trades or significant account changes.