package ws

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

// createTestClient creates a WebSocket client for testing
func createTestClient() *BaseWsClient {
	logger := zerolog.Nop() // No-op logger for tests
	return NewBitgetBaseWsClient(logger, "wss://test.example.com", "")
}

func TestTimeframeConstants(t *testing.T) {
	// Test that all timeframe constants are properly defined
	assert.Equal(t, "1m", Timeframe1m)
	assert.Equal(t, "5m", Timeframe5m)
	assert.Equal(t, "15m", Timeframe15m)
	assert.Equal(t, "30m", Timeframe30m)
	assert.Equal(t, "1h", Timeframe1h)
	assert.Equal(t, "4h", Timeframe4h)
	assert.Equal(t, "6h", Timeframe6h)
	assert.Equal(t, "12h", Timeframe12h)
	assert.Equal(t, "1d", Timeframe1d)
	assert.Equal(t, "3d", Timeframe3d)
	assert.Equal(t, "1w", Timeframe1w)
	assert.Equal(t, "1M", Timeframe1M)
}

func TestChannelConstants(t *testing.T) {
	// Test that all public channel constants are properly defined
	assert.Equal(t, "ticker", ChannelTicker)
	assert.Equal(t, "candle", ChannelCandle)
	assert.Equal(t, "books", ChannelBooks)
	assert.Equal(t, "books5", ChannelBooks5)
	assert.Equal(t, "books15", ChannelBooks15)
	assert.Equal(t, "trade", ChannelTrade)
	assert.Equal(t, "mark-price", ChannelMarkPrice)
	assert.Equal(t, "funding-time", ChannelFundingTime)

	// Test that all private channel constants are properly defined
	assert.Equal(t, "orders", ChannelOrders)
	assert.Equal(t, "fill", ChannelFill)
	assert.Equal(t, "positions", ChannelPositions)
	assert.Equal(t, "account", ChannelAccount)
	assert.Equal(t, "plan-order", ChannelPlanOrder)
}

func TestSubscribeTicker(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Test subscription
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)

	// Verify subscription is stored
	assert.True(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())

	// Verify handler is stored correctly
	subscriptions := client.GetActiveSubscriptions()
	expectedArgs := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelTicker,
		Symbol:      "BTCUSDT",
	}
	assert.Contains(t, subscriptions, expectedArgs)
}

func TestSubscribeCandles(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Test subscription with different timeframes
	testCases := []struct {
		timeframe       string
		expectedChannel string
	}{
		{Timeframe1m, "candle1m"},
		{Timeframe5m, "candle5m"},
		{Timeframe1h, "candle1h"},
		{Timeframe1d, "candle1d"},
	}

	for _, tc := range testCases {
		t.Run(tc.timeframe, func(t *testing.T) {
			client.SubscribeCandles("ETHUSDT", "USDT-FUTURES", tc.timeframe, handler)

			// Verify subscription is stored with correct channel name
			assert.True(t, client.IsSubscribed(tc.expectedChannel, "ETHUSDT", "USDT-FUTURES"))
		})
	}

	assert.Equal(t, len(testCases), client.GetSubscriptionCount())
}

func TestSubscribeOrderBookVariants(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Test all order book subscription variants
	client.SubscribeOrderBook("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeOrderBook5("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeOrderBook15("BTCUSDT", "USDT-FUTURES", handler)

	// Verify all subscriptions
	assert.True(t, client.IsSubscribed(ChannelBooks, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelBooks5, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelBooks15, "BTCUSDT", "USDT-FUTURES"))
	assert.Equal(t, 3, client.GetSubscriptionCount())
}

func TestSubscribeTrades(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribeTrades("ADAUSDT", "USDT-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelTrade, "ADAUSDT", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestSubscribeFundingTime(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribeFundingTime("ETHUSDT", "USDT-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelFundingTime, "ETHUSDT", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestUnsubscribe(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Add some subscriptions
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe1m, handler)
	client.SubscribeOrderBook("BTCUSDT", "USDT-FUTURES", handler)

	assert.Equal(t, 3, client.GetSubscriptionCount())

	// Test generic unsubscribe
	client.Unsubscribe(ChannelTicker, "BTCUSDT", "USDT-FUTURES")
	assert.False(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))
	assert.Equal(t, 2, client.GetSubscriptionCount())

	// Test specific unsubscribe methods
	client.UnsubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe1m)
	assert.False(t, client.IsSubscribed("candle1m", "BTCUSDT", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())

	client.UnsubscribeOrderBook("BTCUSDT", "USDT-FUTURES")
	assert.False(t, client.IsSubscribed(ChannelBooks, "BTCUSDT", "USDT-FUTURES"))
	assert.Equal(t, 0, client.GetSubscriptionCount())
}

func TestSpecificUnsubscribeMethods(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Add subscriptions for all channel types
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe5m, handler)
	client.SubscribeOrderBook("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeOrderBook5("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeOrderBook15("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeTrades("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeMarkPrice("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeFundingTime("BTCUSDT", "USDT-FUTURES", handler)

	initialCount := client.GetSubscriptionCount()
	assert.Equal(t, 8, initialCount)

	// Test all specific unsubscribe methods
	client.UnsubscribeTicker("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, initialCount-1, client.GetSubscriptionCount())

	client.UnsubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe5m)
	assert.Equal(t, initialCount-2, client.GetSubscriptionCount())

	client.UnsubscribeOrderBook("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, initialCount-3, client.GetSubscriptionCount())

	client.UnsubscribeOrderBook5("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, initialCount-4, client.GetSubscriptionCount())

	client.UnsubscribeOrderBook15("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, initialCount-5, client.GetSubscriptionCount())

	client.UnsubscribeTrades("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, initialCount-6, client.GetSubscriptionCount())

	client.UnsubscribeMarkPrice("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, initialCount-7, client.GetSubscriptionCount())

	client.UnsubscribeFundingTime("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, 0, client.GetSubscriptionCount())
}

func TestMultipleSymbolsAndProductTypes(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Test subscriptions with different symbols and product types
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeTicker("ETHUSDT", "USDT-FUTURES", handler)
	client.SubscribeTicker("BTCUSD", "COIN-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelTicker, "ETHUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelTicker, "BTCUSD", "COIN-FUTURES"))

	// These should not be subscribed
	assert.False(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "COIN-FUTURES"))
	assert.False(t, client.IsSubscribed(ChannelTicker, "ETHUSDT", "COIN-FUTURES"))

	assert.Equal(t, 3, client.GetSubscriptionCount())
}

func TestGetActiveSubscriptions(t *testing.T) {
	client := createTestClient()

	handler1 := func(message string) {}
	handler2 := func(message string) {}

	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler1)
	client.SubscribeCandles("ETHUSDT", "USDT-FUTURES", Timeframe1h, handler2)

	subscriptions := client.GetActiveSubscriptions()

	// Verify we get a copy and both subscriptions are present
	assert.Equal(t, 2, len(subscriptions))

	expectedArgs1 := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelTicker,
		Symbol:      "BTCUSDT",
	}

	expectedArgs2 := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     "candle1h",
		Symbol:      "ETHUSDT",
	}

	assert.Contains(t, subscriptions, expectedArgs1)
	assert.Contains(t, subscriptions, expectedArgs2)
}

func TestIsSubscribed(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Initially not subscribed
	assert.False(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))

	// After subscription
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	assert.True(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))

	// Different symbol should not be subscribed
	assert.False(t, client.IsSubscribed(ChannelTicker, "ETHUSDT", "USDT-FUTURES"))

	// Different product type should not be subscribed
	assert.False(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "COIN-FUTURES"))

	// Different channel should not be subscribed
	assert.False(t, client.IsSubscribed(ChannelTrade, "BTCUSDT", "USDT-FUTURES"))
}

func TestGetSubscriptionCount(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Initially empty
	assert.Equal(t, 0, client.GetSubscriptionCount())

	// Add subscriptions
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	assert.Equal(t, 1, client.GetSubscriptionCount())

	client.SubscribeCandles("ETHUSDT", "USDT-FUTURES", Timeframe1m, handler)
	assert.Equal(t, 2, client.GetSubscriptionCount())

	client.SubscribeOrderBook("ADAUSDT", "USDT-FUTURES", handler)
	assert.Equal(t, 3, client.GetSubscriptionCount())

	// Remove subscription
	client.UnsubscribeTicker("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, 2, client.GetSubscriptionCount())
}

func TestSubscriptionArgs(t *testing.T) {
	// Test that SubscriptionArgs struct works correctly as map key
	args1 := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelTicker,
		Symbol:      "BTCUSDT",
	}

	args2 := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelTicker,
		Symbol:      "BTCUSDT",
	}

	args3 := SubscriptionArgs{
		ProductType: "COIN-FUTURES",
		Channel:     ChannelTicker,
		Symbol:      "BTCUSDT",
	}

	// Same args should be equal
	assert.Equal(t, args1, args2)

	// Different product type should not be equal
	assert.NotEqual(t, args1, args3)

	// Test as map keys
	testMap := make(map[SubscriptionArgs]string)
	testMap[args1] = "handler1"
	testMap[args3] = "handler2"

	assert.Equal(t, "handler1", testMap[args2]) // args2 should match args1
	assert.Equal(t, "handler2", testMap[args3])
	assert.Equal(t, 2, len(testMap))
}

func TestHandlerInvocation(t *testing.T) {
	client := createTestClient()

	// Test that handlers are stored correctly and can be retrieved
	handlerCalled := false
	handler := func(message string) {
		handlerCalled = true
		assert.Equal(t, "test message", message)
	}

	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)

	// Simulate getting the handler (this would normally be done by GetListener)
	args := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelTicker,
		Symbol:      "BTCUSDT",
	}

	storedHandler, exists := client.subscriptions[args]
	assert.True(t, exists)
	assert.NotNil(t, storedHandler)

	// Invoke the handler
	storedHandler("test message")
	assert.True(t, handlerCalled)
}

func TestConcurrentSubscriptions(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Test multiple subscriptions to the same symbol with different channels
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe1m, handler)
	client.SubscribeCandles("BTCUSDT", "USDT-FUTURES", Timeframe5m, handler)
	client.SubscribeOrderBook("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeTrades("BTCUSDT", "USDT-FUTURES", handler)

	assert.Equal(t, 5, client.GetSubscriptionCount())

	// Verify each subscription is independent
	assert.True(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed("candle1m", "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed("candle5m", "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelBooks, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelTrade, "BTCUSDT", "USDT-FUTURES"))

	// Unsubscribe one should not affect others
	client.UnsubscribeTicker("BTCUSDT", "USDT-FUTURES")
	assert.Equal(t, 4, client.GetSubscriptionCount())
	assert.False(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed("candle1m", "BTCUSDT", "USDT-FUTURES"))
}

// =============================================================================
// PRIVATE CHANNEL TESTS
// =============================================================================

func TestSubscribeOrders(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribeOrders("USDT-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelOrders, "", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestSubscribeFills(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribeFills("COIN-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelFill, "", "COIN-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestSubscribePositions(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribePositions("USDT-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelPositions, "", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestSubscribeAccount(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribeAccount("USDT-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelAccount, "", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestSubscribePlanOrders(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	client.SubscribePlanOrders("USDT-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelPlanOrder, "", "USDT-FUTURES"))
	assert.Equal(t, 1, client.GetSubscriptionCount())
}

func TestPrivateChannelUnsubscribeMethods(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Add subscriptions for all private channel types
	client.SubscribeOrders("USDT-FUTURES", handler)
	client.SubscribeFills("USDT-FUTURES", handler)
	client.SubscribePositions("USDT-FUTURES", handler)
	client.SubscribeAccount("USDT-FUTURES", handler)
	client.SubscribePlanOrders("USDT-FUTURES", handler)

	initialCount := client.GetSubscriptionCount()
	assert.Equal(t, 5, initialCount)

	// Test all specific unsubscribe methods for private channels
	client.UnsubscribeOrders("USDT-FUTURES")
	assert.Equal(t, initialCount-1, client.GetSubscriptionCount())
	assert.False(t, client.IsSubscribed(ChannelOrders, "", "USDT-FUTURES"))

	client.UnsubscribeFills("USDT-FUTURES")
	assert.Equal(t, initialCount-2, client.GetSubscriptionCount())
	assert.False(t, client.IsSubscribed(ChannelFill, "", "USDT-FUTURES"))

	client.UnsubscribePositions("USDT-FUTURES")
	assert.Equal(t, initialCount-3, client.GetSubscriptionCount())
	assert.False(t, client.IsSubscribed(ChannelPositions, "", "USDT-FUTURES"))

	client.UnsubscribeAccount("USDT-FUTURES")
	assert.Equal(t, initialCount-4, client.GetSubscriptionCount())
	assert.False(t, client.IsSubscribed(ChannelAccount, "", "USDT-FUTURES"))

	client.UnsubscribePlanOrders("USDT-FUTURES")
	assert.Equal(t, 0, client.GetSubscriptionCount())
	assert.False(t, client.IsSubscribed(ChannelPlanOrder, "", "USDT-FUTURES"))
}

func TestMultipleProductTypesPrivateChannels(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Test subscriptions with different product types for private channels
	client.SubscribeOrders("USDT-FUTURES", handler)
	client.SubscribeOrders("COIN-FUTURES", handler)
	client.SubscribeFills("USDT-FUTURES", handler)
	client.SubscribeFills("COIN-FUTURES", handler)

	assert.True(t, client.IsSubscribed(ChannelOrders, "", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelOrders, "", "COIN-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelFill, "", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelFill, "", "COIN-FUTURES"))

	// These should not be subscribed
	assert.False(t, client.IsSubscribed(ChannelPositions, "", "USDT-FUTURES"))
	assert.False(t, client.IsSubscribed(ChannelAccount, "", "COIN-FUTURES"))

	assert.Equal(t, 4, client.GetSubscriptionCount())
}

func TestPrivateChannelSubscriptionArgs(t *testing.T) {
	// Test that private channel subscription args work correctly as map keys
	args1 := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelOrders,
		Symbol:      "", // Private channels use empty symbol
	}

	args2 := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelOrders,
		Symbol:      "", // Private channels use empty symbol
	}

	args3 := SubscriptionArgs{
		ProductType: "COIN-FUTURES",
		Channel:     ChannelOrders,
		Symbol:      "", // Private channels use empty symbol
	}

	// Same args should be equal
	assert.Equal(t, args1, args2)

	// Different product type should not be equal
	assert.NotEqual(t, args1, args3)

	// Test as map keys
	testMap := make(map[SubscriptionArgs]string)
	testMap[args1] = "handler1"
	testMap[args3] = "handler2"

	assert.Equal(t, "handler1", testMap[args2]) // args2 should match args1
	assert.Equal(t, "handler2", testMap[args3])
	assert.Equal(t, 2, len(testMap))
}

func TestMixedPublicAndPrivateSubscriptions(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Mix of public and private subscriptions
	client.SubscribeTicker("BTCUSDT", "USDT-FUTURES", handler)
	client.SubscribeOrders("USDT-FUTURES", handler)
	client.SubscribeCandles("ETHUSDT", "USDT-FUTURES", Timeframe1m, handler)
	client.SubscribeFills("USDT-FUTURES", handler)
	client.SubscribeOrderBook("ADAUSDT", "USDT-FUTURES", handler)
	client.SubscribePositions("USDT-FUTURES", handler)

	assert.Equal(t, 6, client.GetSubscriptionCount())

	// Verify each subscription type
	assert.True(t, client.IsSubscribed(ChannelTicker, "BTCUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelOrders, "", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed("candle1m", "ETHUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelFill, "", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelBooks, "ADAUSDT", "USDT-FUTURES"))
	assert.True(t, client.IsSubscribed(ChannelPositions, "", "USDT-FUTURES"))
}

func TestIsLoggedIn(t *testing.T) {
	client := createTestClient()

	// Initially not logged in
	assert.False(t, client.IsLoggedIn())

	// Manually set login status for testing (normally done by Login method)
	client.loginStatus = true
	assert.True(t, client.IsLoggedIn())

	client.loginStatus = false
	assert.False(t, client.IsLoggedIn())
}

func TestIsConnected(t *testing.T) {
	client := createTestClient()

	// Initially not connected
	assert.False(t, client.IsConnected())

	// Manually set connection status for testing
	client.connected = true
	assert.True(t, client.IsConnected())

	client.connected = false
	assert.False(t, client.IsConnected())
}

func TestRequiresAuth(t *testing.T) {
	// Test public channels (should not require auth)
	assert.False(t, RequiresAuth(ChannelTicker))
	assert.False(t, RequiresAuth(ChannelCandle))
	assert.False(t, RequiresAuth(ChannelBooks))
	assert.False(t, RequiresAuth(ChannelBooks5))
	assert.False(t, RequiresAuth(ChannelBooks15))
	assert.False(t, RequiresAuth(ChannelTrade))
	assert.False(t, RequiresAuth(ChannelMarkPrice))
	assert.False(t, RequiresAuth(ChannelFundingTime))
	assert.False(t, RequiresAuth("candle1m"))        // Constructed channel name
	assert.False(t, RequiresAuth("unknown-channel")) // Unknown channel

	// Test private channels (should require auth)
	assert.True(t, RequiresAuth(ChannelOrders))
	assert.True(t, RequiresAuth(ChannelFill))
	assert.True(t, RequiresAuth(ChannelPositions))
	assert.True(t, RequiresAuth(ChannelAccount))
	assert.True(t, RequiresAuth(ChannelPlanOrder))
}

func TestPrivateChannelHandlerInvocation(t *testing.T) {
	client := createTestClient()

	// Test that private channel handlers are stored correctly and can be retrieved
	handlerCalled := false
	handler := func(message string) {
		handlerCalled = true
		assert.Equal(t, "test private message", message)
	}

	client.SubscribeOrders("USDT-FUTURES", handler)

	// Simulate getting the handler for private channel
	args := SubscriptionArgs{
		ProductType: "USDT-FUTURES",
		Channel:     ChannelOrders,
		Symbol:      "", // Private channels use empty symbol
	}

	storedHandler, exists := client.subscriptions[args]
	assert.True(t, exists)
	assert.NotNil(t, storedHandler)

	// Invoke the handler
	storedHandler("test private message")
	assert.True(t, handlerCalled)
}

func TestAllPrivateChannelTypes(t *testing.T) {
	client := createTestClient()

	handler := func(message string) {}

	// Subscribe to all private channel types
	client.SubscribeOrders("USDT-FUTURES", handler)
	client.SubscribeFills("USDT-FUTURES", handler)
	client.SubscribePositions("USDT-FUTURES", handler)
	client.SubscribeAccount("USDT-FUTURES", handler)
	client.SubscribePlanOrders("USDT-FUTURES", handler)

	// Verify all are subscribed
	assert.Equal(t, 5, client.GetSubscriptionCount())

	subscriptions := client.GetActiveSubscriptions()
	assert.Equal(t, 5, len(subscriptions))

	// Verify specific subscriptions exist
	expectedChannels := []string{
		ChannelOrders,
		ChannelFill,
		ChannelPositions,
		ChannelAccount,
		ChannelPlanOrder,
	}

	for _, channel := range expectedChannels {
		assert.True(t, client.IsSubscribed(channel, "", "USDT-FUTURES"))

		args := SubscriptionArgs{
			ProductType: "USDT-FUTURES",
			Channel:     channel,
			Symbol:      "",
		}
		assert.Contains(t, subscriptions, args)
	}
}
