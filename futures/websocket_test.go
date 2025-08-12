package futures

import (
	"testing"
	"time"

	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockBaseWsClient mocks the WebSocket client for testing
type MockBaseWsClient struct {
	mock.Mock
	connected      bool
	loggedIn       bool
	subscriptions  map[string]bool
	subscriptCount int
}

func (m *MockBaseWsClient) SubscribeTicker(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-ticker"] = true
}

func (m *MockBaseWsClient) SubscribeCandles(symbol, productType, timeframe string, handler ws.OnReceive) {
	m.Called(symbol, productType, timeframe, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-candles"] = true
}

func (m *MockBaseWsClient) SubscribeOrderBook5(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-orderbook5"] = true
}

func (m *MockBaseWsClient) SubscribeOrderBook15(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-orderbook15"] = true
}

func (m *MockBaseWsClient) SubscribeOrderBook(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-orderbook"] = true
}

func (m *MockBaseWsClient) SubscribeTrades(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-trades"] = true
}

func (m *MockBaseWsClient) SubscribeMarkPrice(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-markprice"] = true
}

func (m *MockBaseWsClient) SubscribeFundingTime(symbol, productType string, handler ws.OnReceive) {
	m.Called(symbol, productType, handler)
	m.subscriptCount++
	m.subscriptions[symbol+"-funding"] = true
}

func (m *MockBaseWsClient) SubscribeOrders(productType string, handler ws.OnReceive) {
	m.Called(productType, handler)
	m.subscriptCount++
	m.subscriptions["orders"] = true
}

func (m *MockBaseWsClient) SubscribeFills(productType string, handler ws.OnReceive) {
	m.Called(productType, handler)
	m.subscriptCount++
	m.subscriptions["fills"] = true
}

func (m *MockBaseWsClient) SubscribePositions(productType string, handler ws.OnReceive) {
	m.Called(productType, handler)
	m.subscriptCount++
	m.subscriptions["positions"] = true
}

func (m *MockBaseWsClient) SubscribeAccount(productType string, handler ws.OnReceive) {
	m.Called(productType, handler)
	m.subscriptCount++
	m.subscriptions["account"] = true
}

func (m *MockBaseWsClient) IsConnected() bool {
	return m.connected
}

func (m *MockBaseWsClient) IsLoggedIn() bool {
	return m.loggedIn
}

func (m *MockBaseWsClient) GetSubscriptionCount() int {
	return m.subscriptCount
}

func (m *MockBaseWsClient) Connect() {
	m.Called()
}

func (m *MockBaseWsClient) ConnectWebSocket() {
	m.Called()
	m.connected = true
}

func (m *MockBaseWsClient) StartReadLoop() {
	m.Called()
}

func (m *MockBaseWsClient) Login(apiKey, passphrase string, signType common.SignType) {
	m.Called(apiKey, passphrase, signType)
	m.loggedIn = true
}

func (m *MockBaseWsClient) SetListener(msgHandler, errorHandler ws.OnReceive) {
	m.Called(msgHandler, errorHandler)
}

func (m *MockBaseWsClient) SetReconnectionTimeout(timeout time.Duration) {
	m.Called(timeout)
}

func (m *MockBaseWsClient) SetCheckConnectionInterval(interval time.Duration) {
	m.Called(interval)
}

func (m *MockBaseWsClient) Close() {
	m.Called()
	m.connected = false
	m.loggedIn = false
}

// Create a test WebSocket manager with mock client
func createTestWebSocketManager() *WebSocketManager {
	client := &Client{
		apiKey:    "test-key",
		secretKey: "test-secret",
	}

	mockClient := &MockBaseWsClient{
		subscriptions: make(map[string]bool),
	}
	
	return &WebSocketManager{
		client:  client,
		logger:  zerolog.Nop(), // Silent logger for tests
		wsClient: mockClient,
	}
}

func TestWebSocketManager_NewWebSocketManager(t *testing.T) {
	client := NewClient("test-key", "test-secret", "test-passphrase")
	wsManager := client.NewWebSocketManager()

	assert.NotNil(t, wsManager)
	assert.Equal(t, client, wsManager.client)
	assert.False(t, wsManager.isConnected)
	assert.False(t, wsManager.isPrivate)
	assert.False(t, wsManager.isLoggedIn)
	assert.True(t, wsManager.autoReconnect)
}

func TestWebSocketManager_ConnectPublic(t *testing.T) {
	// Skip connection tests since they involve real WebSocket connections
	// These would be covered by integration tests
	t.Skip("Connection tests require real WebSocket infrastructure")
}

func TestWebSocketManager_ConnectPrivate(t *testing.T) {
	// Skip connection tests since they involve real WebSocket connections
	// These would be covered by integration tests
	t.Skip("Connection tests require real WebSocket infrastructure")
}

func TestWebSocketManager_SubscribeToTicker(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {
		// Test handler
	}

	mockClient.On("SubscribeTicker", "BTCUSDT", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToTicker("BTCUSDT", handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToTicker_NotConnected(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = false

	handler := func(message string) {}

	err := wsManager.SubscribeToTicker("BTCUSDT", handler)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "WebSocket not connected")
}

func TestWebSocketManager_SubscribeToCandlesticks(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeCandles", "BTCUSDT", string(ProductTypeUSDTFutures), "1m", mock.Anything).Return()

	err := wsManager.SubscribeToCandlesticks("BTCUSDT", "1m", handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToOrderBook_Levels(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	tests := []struct {
		name   string
		levels int
		method string
	}{
		{"OrderBook5", 5, "SubscribeOrderBook5"},
		{"OrderBook15", 15, "SubscribeOrderBook15"},
		{"OrderBookFull", 0, "SubscribeOrderBook"},
		{"OrderBookDefault", 20, "SubscribeOrderBook"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient.On(tt.method, "BTCUSDT", string(ProductTypeUSDTFutures), mock.Anything).Return()

			err := wsManager.SubscribeToOrderBook("BTCUSDT", tt.levels, handler)

			assert.NoError(t, err)
		})
	}

	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToTrades(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeTrades", "BTCUSDT", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToTrades("BTCUSDT", handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToMarkPrice(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeMarkPrice", "BTCUSDT", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToMarkPrice("BTCUSDT", handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToFunding(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeFundingTime", "BTCUSDT", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToFunding("BTCUSDT", handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// Private channel tests

func TestWebSocketManager_SubscribeToOrders(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	wsManager.isPrivate = true
	wsManager.isLoggedIn = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeOrders", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToOrders(handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToOrders_NotAuthenticated(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	wsManager.isPrivate = false
	wsManager.isLoggedIn = false

	handler := func(message string) {}

	err := wsManager.SubscribeToOrders(handler)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "private WebSocket not authenticated")
}

func TestWebSocketManager_SubscribeToFills(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	wsManager.isPrivate = true
	wsManager.isLoggedIn = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeFills", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToFills(handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToPositions(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	wsManager.isPrivate = true
	wsManager.isLoggedIn = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribePositions", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToPositions(handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SubscribeToAccount(t *testing.T) {
	wsManager := createTestWebSocketManager()
	wsManager.isConnected = true
	wsManager.isPrivate = true
	wsManager.isLoggedIn = true
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	handler := func(message string) {}

	mockClient.On("SubscribeAccount", string(ProductTypeUSDTFutures), mock.Anything).Return()

	err := wsManager.SubscribeToAccount(handler)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

// Utility method tests

func TestWebSocketManager_IsConnected(t *testing.T) {
	wsManager := createTestWebSocketManager()
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	// Test not connected
	mockClient.connected = false
	wsManager.isConnected = false
	assert.False(t, wsManager.IsConnected())

	// Test connected
	mockClient.connected = true
	wsManager.isConnected = true
	assert.True(t, wsManager.IsConnected())
}

func TestWebSocketManager_IsAuthenticated(t *testing.T) {
	wsManager := createTestWebSocketManager()
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	// Test not authenticated
	wsManager.isPrivate = false
	wsManager.isLoggedIn = false
	mockClient.loggedIn = false
	assert.False(t, wsManager.IsAuthenticated())

	// Test authenticated
	wsManager.isPrivate = true
	wsManager.isLoggedIn = true
	mockClient.loggedIn = true
	assert.True(t, wsManager.IsAuthenticated())
}

func TestWebSocketManager_GetSubscriptionCount(t *testing.T) {
	wsManager := createTestWebSocketManager()
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	// Test no subscriptions
	mockClient.subscriptCount = 0
	assert.Equal(t, 0, wsManager.GetSubscriptionCount())

	// Test with subscriptions
	mockClient.subscriptCount = 5
	assert.Equal(t, 5, wsManager.GetSubscriptionCount())

	// Test with nil client
	wsManager.wsClient = nil
	assert.Equal(t, 0, wsManager.GetSubscriptionCount())
}

func TestWebSocketManager_Close(t *testing.T) {
	wsManager := createTestWebSocketManager()
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	mockClient.On("Close").Return()

	err := wsManager.Close()

	assert.NoError(t, err)
	assert.False(t, wsManager.isConnected)
	assert.False(t, wsManager.isLoggedIn)
	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SetReconnectionTimeout(t *testing.T) {
	wsManager := createTestWebSocketManager()
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	timeout := 30 * time.Second
	mockClient.On("SetReconnectionTimeout", timeout).Return()

	wsManager.SetReconnectionTimeout(timeout)

	mockClient.AssertExpectations(t)
}

func TestWebSocketManager_SetConnectionCheckInterval(t *testing.T) {
	wsManager := createTestWebSocketManager()
	mockClient := wsManager.wsClient.(*MockBaseWsClient)

	interval := 5 * time.Second
	mockClient.On("SetCheckConnectionInterval", interval).Return()

	wsManager.SetConnectionCheckInterval(interval)

	mockClient.AssertExpectations(t)
}

// High-level configuration tests

func TestWebSocketManager_CreateMarketDataStream(t *testing.T) {
	// Skip high-level configuration tests that involve connections
	t.Skip("Configuration tests require connection infrastructure")
}

func TestWebSocketManager_CreateTradingStream(t *testing.T) {
	// Skip high-level configuration tests that involve connections
	t.Skip("Configuration tests require connection infrastructure")
}

// Default configuration tests

func TestDefaultMarketDataConfig(t *testing.T) {
	config := DefaultMarketDataConfig()

	assert.True(t, config.EnableTicker)
	assert.True(t, config.EnableCandles)
	assert.True(t, config.EnableOrderBook)
	assert.False(t, config.EnableTrades)
	assert.False(t, config.EnableMarkPrice)
	assert.False(t, config.EnableFunding)
	
	assert.Equal(t, ws.Timeframe1m, config.CandleTimeframe)
	assert.Equal(t, 5, config.OrderBookLevels)
	
	assert.NotNil(t, config.TickerHandler)
	assert.NotNil(t, config.CandleHandler)
	assert.NotNil(t, config.OrderBookHandler)
}

func TestDefaultTradingStreamConfig(t *testing.T) {
	config := DefaultTradingStreamConfig()

	assert.True(t, config.EnableOrders)
	assert.True(t, config.EnableFills)
	assert.True(t, config.EnablePositions)
	assert.False(t, config.EnableAccount)
	
	assert.NotNil(t, config.OrderHandler)
	assert.NotNil(t, config.FillHandler)
	assert.NotNil(t, config.PositionHandler)
}