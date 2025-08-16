@echo off
setlocal enabledelayedexpansion

REM Modern Go SDK Documentation Generator  
REM Uses industry-standard tools and practices

echo Generating modern documentation for Bitget Go SDK...

REM Ensure we have required tools
where go >nul 2>&1
if errorlevel 1 (
    echo ERROR: go is required but not installed.
    exit /b 1
)

REM Create docs directory
if not exist docs mkdir docs

echo Generating package documentation...

REM Generate comprehensive API reference
(
echo # API Reference
echo.
echo This document provides comprehensive API reference for the Bitget Go SDK.
echo.
echo ## Quick Navigation
echo.
echo - [Futures API]^(#futures-api^) - Legacy futures trading API
echo - [UTA API]^(#uta-api^) - Unified Trading Account ^(Recommended^)
echo - [WebSocket API]^(#websocket-api^) - Real-time data streaming
echo - [Common Utilities]^(#common-utilities^) - Shared utilities and types
echo.
echo ---
echo.
echo ## Futures API
echo.
echo The Futures API provides access to futures trading operations organized into 4 main categories:
echo.
echo ### Account Operations
echo - Account information and balances
echo - Position management and history  
echo - Leverage and margin configuration
echo - Bill/transaction history
echo.
echo ### Market Data Operations
echo - Candlestick/OHLCV data
echo - Ticker information ^(24hr stats^)
echo - Order book depth data
echo - Recent trade history
echo - Contract specifications
echo.
echo ### Position Operations
echo - Get all positions ^(open/closed^)
echo - Single position details
echo - Position history
echo - Close positions
echo.
echo ### Trading Operations  
echo - Create/modify/cancel orders
echo - Batch order operations
echo - Plan orders ^(conditional/trigger^)
echo - Order and fill history
echo.
echo ---
echo.
echo ## UTA API ^(Recommended^)
echo.
echo The Unified Trading Account API is Bitget's next-generation API that supports spot, margin, and futures trading in a single account.
echo.
echo ### Key Features
echo - Auto-detection of demo trading mode
echo - Unified account management
echo - Simplified API structure
echo - Better error handling
echo.
echo ### Main Operations
echo - Account asset management
echo - Order placement and management
echo - Market data access
echo - Account configuration
echo.
echo ---
echo.
echo ## WebSocket API
echo.
echo Production-ready WebSocket implementation with enterprise features.
echo.
echo ### Features
echo - ✅ Rate limiting ^(10 messages/second^)
echo - ✅ Automatic reconnection with configurable timeout
echo - ✅ Subscription restoration after reconnection
echo - ✅ Connection health monitoring  
echo - ✅ Heartbeat mechanism with ping/pong
echo - ✅ Thread-safe subscription management
echo.
echo ### Public Channels
echo - Ticker ^(24hr price statistics^)
echo - Candles ^(real-time OHLCV data^)
echo - Order Book ^(live bid/ask levels^)
echo - Trades ^(real-time executions^)
echo - Mark Price ^(PnL calculation price^)
echo - Funding ^(rates and timing^)
echo.
echo ### Private Channels  
echo - Orders ^(real-time status updates^)
echo - Fills ^(execution confirmations^)
echo - Positions ^(changes and PnL^)
echo - Account ^(balance updates^)
echo - Plan Orders ^(trigger order updates^)
echo.
echo ---
echo.
echo ## Common Utilities
echo.
echo Shared utilities used across all packages:
echo.
echo ### Authentication
echo - HMAC-SHA256 request signing
echo - API key management
echo - Timestamp handling
echo.
echo ### Error Handling
echo - Structured API error types
echo - Network error handling
echo - Retry logic with exponential backoff
echo.
echo ### Data Types
echo - Common constants and enums
echo - Utility functions
echo - Type conversions
) > docs/API_REFERENCE.md

echo Generating usage examples...

REM Generate comprehensive examples (abbreviated for batch file)
(
echo # Usage Examples
echo.
echo This document provides comprehensive usage examples for the Bitget Go SDK.
echo.
echo ## Quick Start
echo.
echo ### Basic Setup
echo.
echo ```go
echo package main
echo.
echo import ^(
echo     "context"
echo     "fmt"
echo     "log"
echo     "os"
echo.
echo     "github.com/khanbekov/go-bitget/futures"
echo     "github.com/khanbekov/go-bitget/uta" 
echo     "github.com/khanbekov/go-bitget/ws"
echo     "github.com/rs/zerolog"
echo ^)
echo.
echo func main^(^) {
echo     // Load credentials
echo     apiKey := os.Getenv^("BITGET_API_KEY"^)
echo     secretKey := os.Getenv^("BITGET_SECRET_KEY"^) 
echo     passphrase := os.Getenv^("BITGET_PASSPHRASE"^)
echo.
echo     // Create clients
echo     futuresClient := futures.NewClient^(apiKey, secretKey, passphrase^)
echo     utaClient := uta.NewClient^(apiKey, secretKey, passphrase^)
echo     
echo     logger := zerolog.New^(os.Stderr^).With^(^).Timestamp^(^).Logger^(^)
echo     wsClient := ws.NewBitgetBaseWsClient^(logger, "wss://ws.bitget.com/v2/ws/public", ""^)
echo }
echo ```
echo.
echo For more comprehensive examples, see the examples/ directory.
) > docs/EXAMPLES.md

echo Generating development guide...

REM Generate development guide (abbreviated for batch file)
(
echo # Development Guide  
echo.
echo This guide covers development practices, testing, and contribution guidelines.
echo.
echo ## Development Setup
echo.
echo ### Prerequisites
echo.
echo - Go 1.23.4 or later
echo - Git
echo - Make ^(optional^)
echo.
echo ### Quick Start
echo.
echo ```bash
echo git clone https://github.com/khanbekov/go-bitget.git
echo cd go-bitget
echo go mod download
echo go test ./...
echo go build .
echo ```
echo.
echo ## Testing
echo.
echo ```bash
echo go test ./...
echo go test -cover ./...
echo go test -race ./...
echo ```
echo.
) > docs/DEVELOPMENT.md

echo Updating HTML index with modern content...

REM Generate modern HTML index (abbreviated for batch file)
(
echo ^<!DOCTYPE html^>
echo ^<html lang="en"^>
echo ^<head^>
echo     ^<meta charset="UTF-8"^>
echo     ^<meta name="viewport" content="width=device-width, initial-scale=1.0"^>
echo     ^<title^>Bitget Go SDK Documentation^</title^>
echo     ^<style^>
echo         :root {
echo             --primary-color: #0066cc;
echo             --secondary-color: #f8f9fa;
echo             --text-color: #333;
echo             --border-color: #dee2e6;
echo             --success-color: #28a745;
echo         }
echo         body {
echo             font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
echo             line-height: 1.6;
echo             color: var^(--text-color^);
echo             margin: 0;
echo             padding: 20px;
echo             max-width: 1200px;
echo             margin: 0 auto;
echo         }
echo         .header {
echo             text-align: center;
echo             margin-bottom: 50px;
echo             padding: 40px 20px;
echo             background: linear-gradient^(135deg, var^(--primary-color^), #004499^);
echo             color: white;
echo             border-radius: 10px;
echo         }
echo         .header h1 {
echo             font-size: 3rem;
echo             margin-bottom: 10px;
echo             font-weight: 300;
echo         }
echo         .badge {
echo             display: inline-block;
echo             padding: 4px 12px;
echo             background: var^(--success-color^);
echo             color: white;
echo             border-radius: 20px;
echo             font-size: 0.8rem;
echo             margin: 0 5px;
echo         }
echo         .features {
echo             display: grid;
echo             grid-template-columns: repeat^(auto-fit, minmax^(300px, 1fr^)^);
echo             gap: 30px;
echo             margin: 50px 0;
echo         }
echo         .feature-card {
echo             background: var^(--secondary-color^);
echo             padding: 30px;
echo             border-radius: 10px;
echo             border: 1px solid var^(--border-color^);
echo         }
echo         .feature-card h3 {
echo             color: var^(--primary-color^);
echo             margin-bottom: 15px;
echo         }
echo         .docs-grid {
echo             display: grid;
echo             grid-template-columns: repeat^(auto-fit, minmax^(250px, 1fr^)^);
echo             gap: 20px;
echo             margin: 40px 0;
echo         }
echo         .doc-link {
echo             display: block;
echo             padding: 20px;
echo             background: white;
echo             border: 2px solid var^(--border-color^);
echo             border-radius: 8px;
echo             text-decoration: none;
echo             color: var^(--text-color^);
echo         }
echo         .doc-link:hover {
echo             border-color: var^(--primary-color^);
echo         }
echo         .doc-link h4 {
echo             color: var^(--primary-color^);
echo             margin-bottom: 10px;
echo         }
echo         pre {
echo             background: #2d3748;
echo             color: #e2e8f0;
echo             padding: 20px;
echo             border-radius: 8px;
echo             overflow-x: auto;
echo         }
echo     ^</style^>
echo ^</head^>
echo ^<body^>
echo     ^<div class="header"^>
echo         ^<h1^>Bitget Go SDK^</h1^>
echo         ^<p^>Production-ready Go SDK for Bitget cryptocurrency exchange^</p^>
echo         ^<div style="margin-top: 20px;"^>
echo             ^<span class="badge"^>Go 1.23.4+^</span^>
echo             ^<span class="badge"^>37+ Services^</span^>
echo             ^<span class="badge"^>WebSocket Real-time^</span^>
echo         ^</div^>
echo     ^</div^>
echo.
echo     ^<div class="features"^>
echo         ^<div class="feature-card"^>
echo             ^<h3^>Futures API^</h3^>
echo             ^<p^>Complete futures trading with 37+ services organized into 4 directories.^</p^>
echo         ^</div^>
echo         ^<div class="feature-card"^>
echo             ^<h3^>UTA API ^(Recommended^)^</h3^>
echo             ^<p^>Unified Trading Account API with auto-detection of demo mode.^</p^>
echo         ^</div^>
echo         ^<div class="feature-card"^>
echo             ^<h3^>WebSocket ^(Unified^)^</h3^>
echo             ^<p^>Production-ready WebSocket with rate limiting and auto-reconnection.^</p^>
echo         ^</div^>
echo     ^</div^>
echo.
echo     ^<h2 style="text-align: center; margin: 50px 0 30px 0; color: var^(--primary-color^);"^>Documentation^</h2^>
echo.
echo     ^<div class="docs-grid"^>
echo         ^<a href="API_REFERENCE.md" class="doc-link"^>
echo             ^<h4^>API Reference^</h4^>
echo             ^<p^>Comprehensive API documentation for all packages.^</p^>
echo         ^</a^>
echo         ^<a href="EXAMPLES.md" class="doc-link"^>
echo             ^<h4^>Examples^</h4^>
echo             ^<p^>Practical usage examples and code snippets.^</p^>
echo         ^</a^>
echo         ^<a href="DEVELOPMENT.md" class="doc-link"^>
echo             ^<h4^>Development Guide^</h4^>
echo             ^<p^>Development setup and contribution guidelines.^</p^>
echo         ^</a^>
echo         ^<a href="../README.md" class="doc-link"^>
echo             ^<h4^>Main README^</h4^>
echo             ^<p^>Project overview and getting started guide.^</p^>
echo         ^</a^>
echo.
echo         ^<a href="../CONTRIBUTING.md" class="doc-link"^>
echo             ^<h4^>Contributing^</h4^>
echo             ^<p^>Guidelines for contributing to the project.^</p^>
echo         ^</a^>
echo.
echo         ^<a href="../LICENSE" class="doc-link"^>
echo             ^<h4^>License^</h4^>
echo             ^<p^>MIT License - Open source licensing terms.^</p^>
echo         ^</a^>
echo.
echo         ^<a href="../examples/" class="doc-link"^>
echo             ^<h4^>Code Examples^</h4^>
echo             ^<p^>Working code examples and demo applications.^</p^>
echo         ^</a^>
echo     ^</div^>
echo.
echo     ^<div style="text-align: center; margin-top: 50px;"^>
echo         ^<h3^>Live Documentation^</h3^>
echo         ^<pre^>^<code^>godoc -http=:6060^</code^>^</pre^>
echo         ^<p^>Visit: ^<a href="http://localhost:6060/pkg/github.com/khanbekov/go-bitget/"^>http://localhost:6060/pkg/github.com/khanbekov/go-bitget/^</a^>^</p^>
echo     ^</div^>
echo ^</body^>
echo ^</html^>
) > docs/index.html

echo Modern documentation generated successfully!
echo.
echo Files created:
echo   docs/API_REFERENCE.md - Comprehensive API documentation
echo   docs/EXAMPLES.md - Usage examples and patterns  
echo   docs/DEVELOPMENT.md - Development and contribution guide
echo   docs/index.html - Modern HTML overview
echo.
echo Next steps:
echo   1. Open docs/index.html in your browser
echo   2. Review the generated documentation
echo   3. Run: godoc -http=:6060 for live docs
echo.
echo Modern improvements:
echo   - Responsive design with modern CSS
echo   - Comprehensive API reference in Markdown
echo   - Detailed usage examples
echo   - Development and contribution guidelines
echo   - Mobile-friendly interface
echo   - Professional styling

pause