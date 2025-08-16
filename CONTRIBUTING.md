# Contributing to Bitget Go SDK

Thank you for your interest in contributing to the Bitget Go SDK! This document provides guidelines for contributors.

## Table of Contents

- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Contributing Process](#contributing-process)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Pull Request Process](#pull-request-process)

## Getting Started

### Prerequisites
- Go 1.23.4 or later
- Git

### Setup
1. Fork and clone the repository:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-bitget.git
   cd go-bitget
   git remote add upstream https://github.com/khanbekov/go-bitget.git
   ```

2. Install dependencies and run tests:
   ```bash
   go mod download
   go test ./...
   ```

## Contributing Process

1. **Create an Issue** (for major changes): Describe the problem or feature
2. **Create a Branch**:
   ```bash
   git checkout -b feature/descriptive-name
   ```
3. **Make Changes**: Follow coding standards, add tests, update docs
4. **Test**: Ensure `go test ./...` passes
5. **Submit PR**: Include clear description of changes

## Coding Standards

- Follow Go conventions and best practices
- Use meaningful names and clear comments for public APIs
- Support method chaining (fluent API pattern)
- All operations should accept `context.Context`
- Use structured error types from `common/types`
- Write tests for new functionality

### Package Structure
- `futures/` - Futures API 
- `uta/` - Unified Trading Account API 
- `ws/` - WebSocket implementation  
- `common/` - Common API and shared utilities

## Testing Guidelines

### Unit Tests

Write tests for all new functionality following existing patterns:

```go
func TestServiceName_Method_Scenario(t *testing.T) {
    mockClient := &MockClient{}
    service := NewServiceName(mockClient)
    
    mockClient.On("callAPI", mock.Anything, "GET", endpoint, params, []byte(nil), false).
        Return(mockResponse, &fasthttp.ResponseHeader{}, nil)
    
    result, err := service.Do(context.Background())
    
    assert.NoError(t, err)
    assert.NotNil(t, result)
    mockClient.AssertExpectations(t)
}
```

### Running Unit Tests
```bash
go test ./...           # All unit tests
go test -cover ./...    # With coverage
go test -race ./...     # Race detection
```

### Integration Testing (Real API)

For testing against real Bitget API endpoints with your own credentials:

#### Prerequisites
- Valid Bitget API credentials (API key, secret key, passphrase)
- Demo trading mode enabled (strongly recommended)
- Network connectivity to Bitget API

#### Setup
1. **Configure Environment**:
   ```bash
   # Create .env file with your credentials
   BITGET_API_KEY=your_api_key_here
   BITGET_SECRET_KEY=your_secret_key_here
   BITGET_PASSPHRASE=your_passphrase_here
   BITGET_DEMO_TRADING=true
   ```

2. **Setup Integration Config**:
   ```bash
   # Copy example configuration
   cp tests/configs/integration.example.json tests/configs/integration.json
   
   # Edit tests/configs/integration.json with your preferences
   # IMPORTANT: Ensure demo_trading is set to true for safety
   ```

3. **Review Configuration**:
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

#### Running Integration Tests

**Easy Method (Recommended)**:
```bash
# Unix/Linux/macOS
tests/scripts/run-integration-tests.sh

# Windows
tests/scripts/run-integration-tests.bat

# Check demo trading configuration
tests/scripts/run-integration-tests.sh -d
```

**Direct Go Testing**:
```bash
# Run all integration tests
go test -tags=integration ./tests/integration/suites -v

# Run specific test suite
go test -tags=integration ./tests/integration/suites -v -run TestAccountEndpoints

# Run with timeout
go test -tags=integration ./tests/integration/suites -v -timeout=5m
```

**Advanced Options**:
```bash
# Use custom config file
tests/scripts/run-integration-tests.sh -c tests/configs/my_config.json

# Run specific suite with timeout
tests/scripts/run-integration-tests.sh -s account -t 10m

# Quiet mode
tests/scripts/run-integration-tests.sh -q
```

#### Safety Guidelines for Integration Testing

⚠️ **IMPORTANT SAFETY MEASURES**:

1. **Always Use Demo Mode**: Set `demo_trading: true` in configuration
2. **Start with Read-Only**: Begin with safe endpoints (account_info, account_list)
3. **Gradual Enablement**: Enable write operations one at a time after verification
4. **Small Test Amounts**: Use minimal values for any financial operations
5. **Review Before Running**: Always check configuration before executing tests

**Safe Testing Phases**:

**Phase 1 - Read-Only Operations**:
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

**Phase 2 - Write Operations (Demo Mode Only)**:
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

#### Integration Test Reports

Tests generate comprehensive reports:
- **Console Output**: Real-time test execution with structured logging
- **JSON Report**: Machine-readable results (`tests/reports/integration_report.json`)
- **HTML Report**: Human-readable results with charts (`tests/reports/integration_report.html`)
- **Test Summary**: Pass/fail counts and performance metrics

#### When to Run Integration Tests

**Required for Contributors**:
- When adding new API endpoints
- When modifying existing service implementations
- When changing authentication or request handling
- Before submitting pull requests that affect API calls

**Optional but Recommended**:
- When updating dependencies that might affect HTTP clients
- When working on error handling improvements
- When testing against different account configurations

#### Troubleshooting Integration Tests

**Common Issues**:
- **Authentication Errors**: Verify API credentials and permissions
- **Demo Mode Limitations**: Some endpoints may not work in demo mode (this is normal)
- **Rate Limiting**: Reduce test frequency or implement longer delays
- **Network Issues**: Check connectivity to Bitget API endpoints

**Getting Help**:
- Check `INTEGRATION_TESTING.md` for detailed troubleshooting
- Review test logs for specific error messages
- Ensure API credentials have proper permissions
- Verify account is in correct mode (UTA vs Classic)

For detailed integration testing documentation, see [`tests/INTEGRATION_TESTING.md`](tests/INTEGRATION_TESTING.md).

## Pull Request Process

1. **Sync with upstream**:
   ```bash
   git fetch upstream
   git rebase upstream/master
   ```

2. **Run quality checks**:
   ```bash
   go fmt ./...
   go vet ./...
   go test ./...
   ```

3. **Create PR** with:
   - Clear description of changes
   - Type of change (bug fix, feature, docs)
   - Testing completed checklist

That's it! Keep contributions focused and follow existing patterns.