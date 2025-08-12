# API Reference

This document provides comprehensive API reference for the Bitget Go SDK.

## Quick Navigation

- [Futures API](#futures-api) - Legacy futures trading API
- [UTA API](#uta-api) - Unified Trading Account (Recommended)
- [WebSocket API](#websocket-api) - Real-time data streaming
- [Common Utilities](#common-utilities) - Shared utilities and types

---

## Futures API

The Futures API provides access to futures trading operations organized into 4 main categories:

### Account Operations
- Account information and balances
- Position management and history  
- Leverage and margin configuration
- Bill/transaction history

### Market Data Operations
- Candlestick/OHLCV data
- Ticker information (24hr stats)
- Order book depth data
- Recent trade history
- Contract specifications

### Position Operations
- Get all positions (open/closed)
- Single position details
- Position history
- Close positions

### Trading Operations  
- Create/modify/cancel orders
- Batch order operations
- Plan orders (conditional/trigger)
- Order and fill history

---

## UTA API (Recommended)

The Unified Trading Account API is Bitget's next-generation API that supports spot, margin, and futures trading in a single account.

### Key Features
- Auto-detection of demo trading mode
- Unified account management
- Simplified API structure
- Better error handling

### Main Operations
- Account asset management
- Order placement and management
- Market data access
- Account configuration

---

## WebSocket API

Production-ready WebSocket implementation with enterprise features.

### Features
- ✅ Rate limiting (10 messages/second)
- ✅ Automatic reconnection with configurable timeout
- ✅ Subscription restoration after reconnection
- ✅ Connection health monitoring  
- ✅ Heartbeat mechanism with ping/pong
- ✅ Thread-safe subscription management

### Public Channels
- Ticker (24hr price statistics)
- Candles (real-time OHLCV data)
- Order Book (live bid/ask levels)
- Trades (real-time executions)
- Mark Price (PnL calculation price)
- Funding (rates and timing)

### Private Channels  
- Orders (real-time status updates)
- Fills (execution confirmations)
- Positions (changes and PnL)
- Account (balance updates)
- Plan Orders (trigger order updates)

---

## Common Utilities

Shared utilities used across all packages:

### Authentication
- HMAC-SHA256 request signing
- API key management
- Timestamp handling

### Error Handling
- Structured API error types
- Network error handling
- Retry logic with exponential backoff

### Data Types
- Common constants and enums
- Utility functions
- Type conversions
