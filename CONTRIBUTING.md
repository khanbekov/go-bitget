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

### Running Tests
```bash
go test ./...           # All tests
go test -cover ./...    # With coverage
go test -race ./...     # Race detection
```

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