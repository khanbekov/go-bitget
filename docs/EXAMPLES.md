# Usage Examples

This document provides comprehensive usage examples for the Bitget Go SDK.

## Quick Start

### Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/khanbekov/go-bitget/futures"
    "github.com/khanbekov/go-bitget/uta" 
    "github.com/khanbekov/go-bitget/ws"
    "github.com/rs/zerolog"
)

func main() {
    // Load credentials
    apiKey := os.Getenv("BITGET_API_KEY")
    secretKey := os.Getenv("BITGET_SECRET_KEY") 
    passphrase := os.Getenv("BITGET_PASSPHRASE")

    // Create clients
    futuresClient := futures.NewClient(apiKey, secretKey, passphrase)
    utaClient := uta.NewClient(apiKey, secretKey, passphrase)
ECHO is off.
    logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
    wsClient := ws.NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")
}
```

For more comprehensive examples, see the examples/ directory.
