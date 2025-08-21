// Package ws provides WebSocket channel subscription helpers for Bitget's public market data channels.
// This file implements subscription methods for real-time market data including tickers, candlesticks,
// order books, trades, mark prices, and funding information.
package ws

import (
	"fmt"
	"github.com/khanbekov/go-bitget/common"
)

// Supported candlestick timeframes
const (
	// Standard timeframes supported by Bitget WebSocket API
	Timeframe1m  = "1m"  // 1 minute
	Timeframe5m  = "5m"  // 5 minutes
	Timeframe15m = "15m" // 15 minutes
	Timeframe30m = "30m" // 30 minutes
	Timeframe1h  = "1h"  // 1 hour
	Timeframe4h  = "4h"  // 4 hours
	Timeframe6h  = "6h"  // 6 hours
	Timeframe12h = "12h" // 12 hours
	Timeframe1d  = "1d"  // 1 day
	Timeframe3d  = "3d"  // 3 days
	Timeframe1w  = "1w"  // 1 week
	Timeframe1M  = "1M"  // 1 month
)

// WebSocket channel names as defined by Bitget API
const (
	// Public channels
	ChannelTicker      = "ticker"       // Real-time ticker updates
	ChannelCandle      = "candle"       // Real-time candlestick data
	ChannelBooks       = "books"        // Full order book depth
	ChannelBooks5      = "books5"       // Top 5 order book levels
	ChannelBooks15     = "books15"      // Top 15 order book levels
	ChannelTrade       = "trade"        // Real-time trade executions
	ChannelMarkPrice   = "mark-price"   // Mark price updates
	ChannelFundingTime = "funding-time" // Funding rate and time

	// Private channels (require authentication)
	ChannelOrders    = "orders"     // Real-time order updates
	ChannelFill      = "fill"       // Real-time fill/execution updates
	ChannelPositions = "positions"  // Real-time position updates
	ChannelAccount   = "account"    // Account balance updates
	ChannelPlanOrder = "plan-order" // Trigger order updates
)

// SubscribeTicker subscribes to real-time ticker updates for a specific symbol.
// Provides 24hr statistics including price, volume, and change data.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming ticker messages
//
// Example:
//
//	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Ticker update:", message)
//	})
func (c *BaseWsClient) SubscribeTicker(symbol, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelTicker,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeCandles subscribes to real-time candlestick data for a specific symbol and timeframe.
// Provides OHLCV data updates as new candlesticks are formed.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - timeframe: Candlestick interval (use Timeframe constants: Timeframe1m, Timeframe5m, etc.)
//   - handler: Callback function to handle incoming candlestick messages
//
// Example:
//
//	client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe1m, func(message string) {
//	    fmt.Println("Candlestick update:", message)
//	})
func (c *BaseWsClient) SubscribeCandles(symbol, productType, timeframe string, handler OnReceive) {
	channelName := fmt.Sprintf("%s%s", ChannelCandle, timeframe)

	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     channelName,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeOrderBook subscribes to full order book depth data for a specific symbol.
// Provides complete bid and ask levels with real-time updates.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming order book messages
//
// Example:
//
//	client.SubscribeOrderBook("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Order book update:", message)
//	})
func (c *BaseWsClient) SubscribeOrderBook(symbol, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelBooks,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeOrderBook5 subscribes to top 5 order book levels for a specific symbol.
// Provides the best 5 bid and ask levels with real-time updates.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming order book messages
//
// Example:
//
//	client.SubscribeOrderBook5("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Top 5 order book update:", message)
//	})
func (c *BaseWsClient) SubscribeOrderBook5(symbol, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelBooks5,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeOrderBook15 subscribes to top 15 order book levels for a specific symbol.
// Provides the best 15 bid and ask levels with real-time updates.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming order book messages
//
// Example:
//
//	client.SubscribeOrderBook15("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Top 15 order book update:", message)
//	})
func (c *BaseWsClient) SubscribeOrderBook15(symbol, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelBooks15,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeTrades subscribes to real-time trade execution data for a specific symbol.
// Provides information about completed trades including price, size, and side.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming trade messages
//
// Example:
//
//	client.SubscribeTrades("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Trade execution:", message)
//	})
func (c *BaseWsClient) SubscribeTrades(symbol, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelTrade,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeMarkPrice subscribes to real-time mark price updates for a specific symbol.
// Mark price is used for PnL calculations and liquidations in futures trading.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming mark price messages
//
// Example:
//
//	client.SubscribeMarkPrice("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Mark price update:", message)
//	})
func (c *BaseWsClient) SubscribeMarkPrice(symbol, productType string, handler OnReceive) {
	// Note: Bitget doesn't have a dedicated "mark-price" channel
	// Mark price data is available through the ticker channel
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelTicker,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeFundingTime subscribes to funding rate and funding time updates for a specific symbol.
// Provides information about funding rates and next funding time for perpetual contracts.
//
// Parameters:
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming funding messages
//
// Example:
//
//	client.SubscribeFundingTime("BTCUSDT", "USDT-FUTURES", func(message string) {
//	    fmt.Println("Funding update:", message)
//	})
func (c *BaseWsClient) SubscribeFundingTime(symbol, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelFundingTime,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// Unsubscribe removes a subscription for a specific channel, symbol, and product type.
// Stops receiving updates for the specified subscription.
//
// Parameters:
//   - channel: Channel name (e.g., "ticker", "candle1m", "books")
//   - symbol: Trading pair symbol (e.g., "BTCUSDT")
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//
// Example:
//
//	client.Unsubscribe("ticker", "BTCUSDT", "USDT-FUTURES")
func (c *BaseWsClient) Unsubscribe(channel, symbol, productType string) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     channel,
		Symbol:      symbol,
	}

	delete(c.subscriptions, args)
	c.unsubscribe(args)
}

// UnsubscribeTicker removes ticker subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeTicker(symbol, productType string) {
	c.Unsubscribe(ChannelTicker, symbol, productType)
}

// UnsubscribeCandles removes candlestick subscription for a specific symbol, product type, and timeframe.
func (c *BaseWsClient) UnsubscribeCandles(symbol, productType, timeframe string) {
	channelName := fmt.Sprintf("%s%s", ChannelCandle, timeframe)
	c.Unsubscribe(channelName, symbol, productType)
}

// UnsubscribeOrderBook removes order book subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeOrderBook(symbol, productType string) {
	c.Unsubscribe(ChannelBooks, symbol, productType)
}

// UnsubscribeOrderBook5 removes top 5 order book subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeOrderBook5(symbol, productType string) {
	c.Unsubscribe(ChannelBooks5, symbol, productType)
}

// UnsubscribeOrderBook15 removes top 15 order book subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeOrderBook15(symbol, productType string) {
	c.Unsubscribe(ChannelBooks15, symbol, productType)
}

// UnsubscribeTrades removes trade subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeTrades(symbol, productType string) {
	c.Unsubscribe(ChannelTrade, symbol, productType)
}

// UnsubscribeMarkPrice removes mark price subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeMarkPrice(symbol, productType string) {
	c.Unsubscribe(ChannelMarkPrice, symbol, productType)
}

// UnsubscribeFundingTime removes funding time subscription for a specific symbol and product type.
func (c *BaseWsClient) UnsubscribeFundingTime(symbol, productType string) {
	c.Unsubscribe(ChannelFundingTime, symbol, productType)
}

// subscribe sends a subscription request to the WebSocket server
func (c *BaseWsClient) subscribe(args SubscriptionArgs) {
	var argsList []interface{}
	argsList = append(argsList, args)

	baseReq := WsBaseReq{
		Op:   common.WsOpSubscribe,
		Args: argsList,
	}

	c.SendByType(baseReq)
}

// unsubscribe sends an unsubscription request to the WebSocket server
func (c *BaseWsClient) unsubscribe(args SubscriptionArgs) {
	var argsList []interface{}
	argsList = append(argsList, args)

	baseReq := WsBaseReq{
		Op:   common.WsOpUnsubscribe,
		Args: argsList,
	}

	c.SendByType(baseReq)
}

// GetActiveSubscriptions returns a copy of all active subscriptions
// Returns a map of subscription arguments to their handlers
func (c *BaseWsClient) GetActiveSubscriptions() map[SubscriptionArgs]OnReceive {
	subscriptions := make(map[SubscriptionArgs]OnReceive)
	for k, v := range c.subscriptions {
		subscriptions[k] = v
	}
	return subscriptions
}

// IsSubscribed checks if a specific channel subscription is active
func (c *BaseWsClient) IsSubscribed(channel, symbol, productType string) bool {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     channel,
		Symbol:      symbol,
	}

	_, exists := c.subscriptions[args]
	return exists
}

// =============================================================================
// PRIVATE WEBSOCKET CHANNELS (Require Authentication)
// =============================================================================

// SubscribeOrders subscribes to real-time order updates for a specific product type.
// Provides updates when orders are created, modified, filled, or cancelled.
// Requires authentication via Login() before subscription.
//
// Parameters:
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming order update messages
//
// Example:
//
//	// First authenticate
//	client.Login(apiKey, passphrase, common.SHA256)
//
//	// Then subscribe to order updates
//	client.SubscribeOrders("USDT-FUTURES", func(message string) {
//	    fmt.Println("Order update:", message)
//	})
func (c *BaseWsClient) SubscribeOrders(productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelOrders,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeFills subscribes to real-time fill/execution updates for a specific product type.
// Provides updates when orders are partially or fully executed.
// Requires authentication via Login() before subscription.
//
// Parameters:
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming fill update messages
//
// Example:
//
//	client.SubscribeFills("USDT-FUTURES", func(message string) {
//	    fmt.Println("Fill update:", message)
//	})
func (c *BaseWsClient) SubscribeFills(symbol string, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelFill,
		Symbol:      symbol,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribePositions subscribes to real-time position updates for a specific product type.
// Provides updates when positions are opened, modified, or closed.
// Requires authentication via Login() before subscription.
//
// Parameters:
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming position update messages
//
// Example:
//
//	client.SubscribePositions("USDT-FUTURES", func(message string) {
//	    fmt.Println("Position update:", message)
//	})
func (c *BaseWsClient) SubscribePositions(productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelPositions,
		Symbol:      "default",
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribeAccount subscribes to real-time account balance updates for a specific product type.
// Provides updates when account balances change due to trading, funding, or transfers.
// Requires authentication via Login() before subscription.
//
// Parameters:
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming account update messages
//
// Example:
//
//	client.SubscribeAccount("USDT-FUTURES", func(message string) {
//	    fmt.Println("Account update:", message)
//	})
func (c *BaseWsClient) SubscribeAccount(coin string, productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelAccount,
		Coin:        coin,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// SubscribePlanOrders subscribes to real-time trigger order updates for a specific product type.
// Provides updates when trigger/conditional orders are created, triggered, or cancelled.
// Requires authentication via Login() before subscription.
//
// Parameters:
//   - productType: Product type ("USDT-FUTURES", "COIN-FUTURES", etc.)
//   - handler: Callback function to handle incoming plan order update messages
//
// Example:
//
//	client.SubscribePlanOrders("USDT-FUTURES", func(message string) {
//	    fmt.Println("Plan order update:", message)
//	})
func (c *BaseWsClient) SubscribePlanOrders(productType string, handler OnReceive) {
	args := SubscriptionArgs{
		ProductType: productType,
		Channel:     ChannelPlanOrder,
	}

	c.subscriptions[args] = handler
	c.subscribe(args)
}

// =============================================================================
// PRIVATE CHANNEL UNSUBSCRIBE METHODS
// =============================================================================

// UnsubscribeOrders removes order updates subscription for a specific product type.
func (c *BaseWsClient) UnsubscribeOrders(productType string) {
	c.Unsubscribe(ChannelOrders, "", productType)
}

// UnsubscribeFills removes fill updates subscription for a specific product type.
func (c *BaseWsClient) UnsubscribeFills(productType string) {
	c.Unsubscribe(ChannelFill, "", productType)
}

// UnsubscribePositions removes position updates subscription for a specific product type.
func (c *BaseWsClient) UnsubscribePositions(productType string) {
	c.Unsubscribe(ChannelPositions, "", productType)
}

// UnsubscribeAccount removes account updates subscription for a specific product type.
func (c *BaseWsClient) UnsubscribeAccount(productType string) {
	c.Unsubscribe(ChannelAccount, "", productType)
}

// UnsubscribePlanOrders removes plan order updates subscription for a specific product type.
func (c *BaseWsClient) UnsubscribePlanOrders(productType string) {
	c.Unsubscribe(ChannelPlanOrder, "", productType)
}

// =============================================================================
// AUTHENTICATION STATUS HELPERS
// =============================================================================

// RequiresAuth checks if a given channel requires authentication.
// Returns true for private channels, false for public channels.
func RequiresAuth(channel string) bool {
	privateChannels := map[string]bool{
		ChannelOrders:    true,
		ChannelFill:      true,
		ChannelPositions: true,
		ChannelAccount:   true,
		ChannelPlanOrder: true,
	}

	return privateChannels[channel]
}
