// Package futures provides WebSocket integration for real-time futures market data and trading updates.
// This file adds convenient methods to the futures client for WebSocket operations.
package futures

import (
	"context"
	"fmt"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
	"os"
	"time"
)

// WebSocketClientInterface defines the interface for WebSocket operations
type WebSocketClientInterface interface {
	// Connection methods
	Connect()
	ConnectWebSocket()
	StartReadLoop()
	Login(apiKey, passphrase string, signType common.SignType)
	SetListener(msgListener ws.OnReceive, errorListener ws.OnReceive)
	IsConnected() bool
	IsLoggedIn() bool
	Close()

	// Configuration methods
	SetReconnectionTimeout(timeout time.Duration)
	SetCheckConnectionInterval(interval time.Duration)

	// Public channel subscriptions
	SubscribeTicker(symbol, productType string, handler ws.OnReceive)
	SubscribeCandles(symbol, productType, timeframe string, handler ws.OnReceive)
	SubscribeOrderBook(symbol, productType string, handler ws.OnReceive)
	SubscribeOrderBook5(symbol, productType string, handler ws.OnReceive)
	SubscribeOrderBook15(symbol, productType string, handler ws.OnReceive)
	SubscribeTrades(symbol, productType string, handler ws.OnReceive)
	SubscribeMarkPrice(symbol, productType string, handler ws.OnReceive)
	SubscribeFundingTime(symbol, productType string, handler ws.OnReceive)

	// Private channel subscriptions
	SubscribeOrders(productType string, handler ws.OnReceive)
	SubscribeFills(string, productType string, handler ws.OnReceive)
	SubscribePositions(productType string, handler ws.OnReceive)
	SubscribeAccount(string, productType string, handler ws.OnReceive)

	// Utility methods
	GetSubscriptionCount() int
}

// WebSocketClientAdapter adapts ws.BaseWsClient to implement WebSocketClientInterface
type WebSocketClientAdapter struct {
	BaseWsClient *ws.BaseWsClient
}

// Connection methods
func (w *WebSocketClientAdapter) Connect() {
	w.BaseWsClient.Connect()
}

func (w *WebSocketClientAdapter) ConnectWebSocket() {
	w.BaseWsClient.ConnectWebSocket()
}

func (w *WebSocketClientAdapter) StartReadLoop() {
	w.BaseWsClient.StartReadLoop()
}

func (w *WebSocketClientAdapter) Login(apiKey, passphrase string, signType common.SignType) {
	w.BaseWsClient.Login(apiKey, passphrase, signType)
}

func (w *WebSocketClientAdapter) SetListener(msgListener ws.OnReceive, errorListener ws.OnReceive) {
	w.BaseWsClient.SetListener(msgListener, errorListener)
}

func (w *WebSocketClientAdapter) IsConnected() bool {
	return w.BaseWsClient.IsConnected()
}

func (w *WebSocketClientAdapter) IsLoggedIn() bool {
	return w.BaseWsClient.IsLoggedIn()
}

func (w *WebSocketClientAdapter) Close() {
	w.BaseWsClient.Close()
}

// Configuration methods
func (w *WebSocketClientAdapter) SetReconnectionTimeout(timeout time.Duration) {
	w.BaseWsClient.SetReconnectionTimeout(timeout)
}

func (w *WebSocketClientAdapter) SetCheckConnectionInterval(interval time.Duration) {
	w.BaseWsClient.SetCheckConnectionInterval(interval)
}

// Public channel subscriptions
func (w *WebSocketClientAdapter) SubscribeTicker(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeTicker(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeCandles(symbol, productType, timeframe string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeCandles(symbol, productType, timeframe, handler)
}

func (w *WebSocketClientAdapter) SubscribeOrderBook(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeOrderBook(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeOrderBook5(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeOrderBook5(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeOrderBook15(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeOrderBook15(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeTrades(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeTrades(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeMarkPrice(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeMarkPrice(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeFundingTime(symbol, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeFundingTime(symbol, productType, handler)
}

// Private channel subscriptions
func (w *WebSocketClientAdapter) SubscribeOrders(productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeOrders(productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeFills(symbol string, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeFills(symbol, productType, handler)
}

func (w *WebSocketClientAdapter) SubscribePositions(productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribePositions(productType, handler)
}

func (w *WebSocketClientAdapter) SubscribeAccount(coin string, productType string, handler ws.OnReceive) {
	w.BaseWsClient.SubscribeAccount(coin, productType, handler)
}

// Utility methods
func (w *WebSocketClientAdapter) GetSubscriptionCount() int {
	return w.BaseWsClient.GetSubscriptionCount()
}

// WebSocketManager handles WebSocket connections for futures trading
type WebSocketManager struct {
	client        *Client
	wsClient      WebSocketClientInterface
	logger        zerolog.Logger
	isPrivate     bool
	isConnected   bool
	isLoggedIn    bool
	autoReconnect bool
}

// NewWebSocketManager creates a new WebSocket manager for the futures client
func (c *Client) NewWebSocketManager() *WebSocketManager {
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	return &WebSocketManager{
		client:        c,
		logger:        logger,
		autoReconnect: true,
	}
}

// ConnectPublic connects to public WebSocket channels for market data
func (wm *WebSocketManager) ConnectPublic() error {
	wm.logger.Info().Msg("Connecting to Bitget public WebSocket...")

	baseClient := ws.NewBitgetBaseWsClient(
		wm.logger,
		"wss://ws.bitget.com/v2/ws/public",
		"",
	)
	wm.wsClient = &WebSocketClientAdapter{BaseWsClient: baseClient}

	wm.wsClient.SetListener(wm.defaultMessageHandler, wm.errorHandler)
	wm.wsClient.Connect()
	wm.wsClient.ConnectWebSocket()
	wm.wsClient.StartReadLoop()

	// Wait for connection
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("connection timeout")
		case <-ticker.C:
			if wm.wsClient.IsConnected() {
				wm.isConnected = true
				wm.isPrivate = false
				wm.logger.Info().Msg("Successfully connected to public WebSocket")
				return nil
			}
		}
	}
}

// ConnectPrivate connects to private WebSocket channels for account updates
func (wm *WebSocketManager) ConnectPrivate(apiKey, passphrase string) error {
	wm.logger.Info().Msg("Connecting to Bitget private WebSocket...")

	baseClient := ws.NewBitgetBaseWsClient(
		wm.logger,
		"wss://ws.bitget.com/v2/ws/private",
		wm.client.secretKey,
	)
	wm.wsClient = &WebSocketClientAdapter{BaseWsClient: baseClient}

	wm.wsClient.SetListener(wm.defaultMessageHandler, wm.errorHandler)
	wm.wsClient.Connect()
	wm.wsClient.ConnectWebSocket()
	wm.wsClient.StartReadLoop()

	// Wait for connection
	timeout := time.After(5 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("connection timeout")
		case <-ticker.C:
			if wm.wsClient.IsConnected() {
				wm.isConnected = true
				wm.isPrivate = true

				// Authenticate for private channels
				wm.wsClient.Login(apiKey, passphrase, common.SHA256)

				// Wait for login
				loginTimeout := time.After(3 * time.Second)
				loginTicker := time.NewTicker(100 * time.Millisecond)
				defer loginTicker.Stop()

				for {
					select {
					case <-loginTimeout:
						return fmt.Errorf("login timeout")
					case <-loginTicker.C:
						if wm.wsClient.IsLoggedIn() {
							wm.isLoggedIn = true
							wm.logger.Info().Msg("Successfully connected and authenticated to private WebSocket")
							return nil
						}
					}
				}
			}
		}
	}
}

// Market Data Subscription Methods

// SubscribeToTicker subscribes to real-time ticker updates
func (wm *WebSocketManager) SubscribeToTicker(symbol string, handler ws.OnReceive) error {
	if !wm.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	wm.wsClient.SubscribeTicker(symbol, string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Str("symbol", symbol).Msg("Subscribed to ticker updates")
	return nil
}

// SubscribeToCandlesticks subscribes to real-time candlestick data
func (wm *WebSocketManager) SubscribeToCandlesticks(symbol, timeframe string, handler ws.OnReceive) error {
	if !wm.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	wm.wsClient.SubscribeCandles(symbol, string(ProductTypeUSDTFutures), timeframe, handler)
	wm.logger.Info().Str("symbol", symbol).Str("timeframe", timeframe).Msg("Subscribed to candlestick updates")
	return nil
}

// SubscribeToOrderBook subscribes to order book updates
func (wm *WebSocketManager) SubscribeToOrderBook(symbol string, levels int, handler ws.OnReceive) error {
	if !wm.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	switch levels {
	case 5:
		wm.wsClient.SubscribeOrderBook5(symbol, string(ProductTypeUSDTFutures), handler)
	case 15:
		wm.wsClient.SubscribeOrderBook15(symbol, string(ProductTypeUSDTFutures), handler)
	default:
		wm.wsClient.SubscribeOrderBook(symbol, string(ProductTypeUSDTFutures), handler)
	}

	wm.logger.Info().Str("symbol", symbol).Int("levels", levels).Msg("Subscribed to order book updates")
	return nil
}

// SubscribeToTrades subscribes to real-time trade executions
func (wm *WebSocketManager) SubscribeToTrades(symbol string, handler ws.OnReceive) error {
	if !wm.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	wm.wsClient.SubscribeTrades(symbol, string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Str("symbol", symbol).Msg("Subscribed to trade updates")
	return nil
}

// SubscribeToMarkPrice subscribes to mark price updates
func (wm *WebSocketManager) SubscribeToMarkPrice(symbol string, handler ws.OnReceive) error {
	if !wm.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	wm.wsClient.SubscribeMarkPrice(symbol, string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Str("symbol", symbol).Msg("Subscribed to mark price updates")
	return nil
}

// SubscribeToFunding subscribes to funding rate updates
func (wm *WebSocketManager) SubscribeToFunding(symbol string, handler ws.OnReceive) error {
	if !wm.isConnected {
		return fmt.Errorf("WebSocket not connected")
	}

	wm.wsClient.SubscribeFundingTime(symbol, string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Str("symbol", symbol).Msg("Subscribed to funding updates")
	return nil
}

// Private Channel Subscription Methods (require authentication)

// SubscribeToOrders subscribes to real-time order updates
func (wm *WebSocketManager) SubscribeToOrders(handler ws.OnReceive) error {
	if !wm.isPrivate || !wm.isLoggedIn {
		return fmt.Errorf("private WebSocket not authenticated")
	}

	wm.wsClient.SubscribeOrders(string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Msg("Subscribed to order updates")
	return nil
}

// SubscribeToFills subscribes to real-time fill updates
func (wm *WebSocketManager) SubscribeToFills(handler ws.OnReceive) error {
	if !wm.isPrivate || !wm.isLoggedIn {
		return fmt.Errorf("private WebSocket not authenticated")
	}

	wm.wsClient.SubscribeFills("default", string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Msg("Subscribed to fill updates")
	return nil
}

// SubscribeToPositions subscribes to real-time position updates
func (wm *WebSocketManager) SubscribeToPositions(handler ws.OnReceive) error {
	if !wm.isPrivate || !wm.isLoggedIn {
		return fmt.Errorf("private WebSocket not authenticated")
	}

	wm.wsClient.SubscribePositions(string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Msg("Subscribed to position updates")
	return nil
}

// SubscribeToAccount subscribes to account balance updates
func (wm *WebSocketManager) SubscribeToAccount(handler ws.OnReceive) error {
	if !wm.isPrivate || !wm.isLoggedIn {
		return fmt.Errorf("private WebSocket not authenticated")
	}

	wm.wsClient.SubscribeAccount("default", string(ProductTypeUSDTFutures), handler)
	wm.logger.Info().Msg("Subscribed to account updates")
	return nil
}

// Utility Methods

// IsConnected returns true if WebSocket is connected
func (wm *WebSocketManager) IsConnected() bool {
	return wm.isConnected && wm.wsClient != nil && wm.wsClient.IsConnected()
}

// IsAuthenticated returns true if private WebSocket is authenticated
func (wm *WebSocketManager) IsAuthenticated() bool {
	return wm.isPrivate && wm.isLoggedIn && wm.wsClient != nil && wm.wsClient.IsLoggedIn()
}

// GetSubscriptionCount returns the number of active subscriptions
func (wm *WebSocketManager) GetSubscriptionCount() int {
	if wm.wsClient == nil {
		return 0
	}
	return wm.wsClient.GetSubscriptionCount()
}

// Close closes the WebSocket connection and cleans up resources
func (wm *WebSocketManager) Close() error {
	if wm.wsClient != nil {
		wm.wsClient.Close()
		wm.isConnected = false
		wm.isLoggedIn = false
		wm.logger.Info().Msg("WebSocket connection closed")
	}
	return nil
}

// SetLogger allows custom logger configuration
func (wm *WebSocketManager) SetLogger(logger zerolog.Logger) {
	wm.logger = logger
}

// EnableAutoReconnect enables/disables automatic reconnection
func (wm *WebSocketManager) EnableAutoReconnect(enable bool) {
	wm.autoReconnect = enable
}

// SetReconnectionTimeout configures the reconnection timeout
func (wm *WebSocketManager) SetReconnectionTimeout(timeout time.Duration) {
	if wm.wsClient != nil {
		wm.wsClient.SetReconnectionTimeout(timeout)
	}
}

// SetConnectionCheckInterval configures how often to check connection health
func (wm *WebSocketManager) SetConnectionCheckInterval(interval time.Duration) {
	if wm.wsClient != nil {
		wm.wsClient.SetCheckConnectionInterval(interval)
	}
}

// Message handlers

func (wm *WebSocketManager) defaultMessageHandler(message string) {
	wm.logger.Debug().Str("message", message).Msg("WebSocket message received")
}

func (wm *WebSocketManager) errorHandler(message string) {
	wm.logger.Error().Str("error", message).Msg("WebSocket error received")
}

// High-level convenience methods

// CreateMarketDataStream creates a complete market data streaming setup
func (wm *WebSocketManager) CreateMarketDataStream(ctx context.Context, symbols []string, config MarketDataConfig) error {
	if err := wm.ConnectPublic(); err != nil {
		return fmt.Errorf("failed to connect to public WebSocket: %w", err)
	}

	for _, symbol := range symbols {
		if config.EnableTicker {
			wm.SubscribeToTicker(symbol, config.TickerHandler)
		}

		if config.EnableCandles {
			wm.SubscribeToCandlesticks(symbol, config.CandleTimeframe, config.CandleHandler)
		}

		if config.EnableOrderBook {
			wm.SubscribeToOrderBook(symbol, config.OrderBookLevels, config.OrderBookHandler)
		}

		if config.EnableTrades {
			wm.SubscribeToTrades(symbol, config.TradeHandler)
		}

		if config.EnableMarkPrice {
			wm.SubscribeToMarkPrice(symbol, config.MarkPriceHandler)
		}

		if config.EnableFunding {
			wm.SubscribeToFunding(symbol, config.FundingHandler)
		}
	}

	wm.logger.Info().Int("symbols", len(symbols)).Int("subscriptions", wm.GetSubscriptionCount()).Msg("Market data stream created")
	return nil
}

// CreateTradingStream creates a complete trading stream setup
func (wm *WebSocketManager) CreateTradingStream(ctx context.Context, apiKey, passphrase string, config TradingStreamConfig) error {
	if err := wm.ConnectPrivate(apiKey, passphrase); err != nil {
		return fmt.Errorf("failed to connect to private WebSocket: %w", err)
	}

	if config.EnableOrders {
		wm.SubscribeToOrders(config.OrderHandler)
	}

	if config.EnableFills {
		wm.SubscribeToFills(config.FillHandler)
	}

	if config.EnablePositions {
		wm.SubscribeToPositions(config.PositionHandler)
	}

	if config.EnableAccount {
		wm.SubscribeToAccount(config.AccountHandler)
	}

	wm.logger.Info().Int("subscriptions", wm.GetSubscriptionCount()).Msg("Trading stream created")
	return nil
}

// Configuration structs for easy setup

// MarketDataConfig configures market data subscriptions
type MarketDataConfig struct {
	EnableTicker    bool
	EnableCandles   bool
	EnableOrderBook bool
	EnableTrades    bool
	EnableMarkPrice bool
	EnableFunding   bool

	CandleTimeframe string
	OrderBookLevels int

	TickerHandler    ws.OnReceive
	CandleHandler    ws.OnReceive
	OrderBookHandler ws.OnReceive
	TradeHandler     ws.OnReceive
	MarkPriceHandler ws.OnReceive
	FundingHandler   ws.OnReceive
}

// TradingStreamConfig configures private trading stream subscriptions
type TradingStreamConfig struct {
	EnableOrders    bool
	EnableFills     bool
	EnablePositions bool
	EnableAccount   bool

	OrderHandler    ws.OnReceive
	FillHandler     ws.OnReceive
	PositionHandler ws.OnReceive
	AccountHandler  ws.OnReceive
}

// Default configurations

// DefaultMarketDataConfig returns a default configuration for market data streaming
func DefaultMarketDataConfig() MarketDataConfig {
	return MarketDataConfig{
		EnableTicker:    true,
		EnableCandles:   true,
		EnableOrderBook: true,
		EnableTrades:    false,
		EnableMarkPrice: false,
		EnableFunding:   false,

		CandleTimeframe: ws.Timeframe1m,
		OrderBookLevels: 5,

		TickerHandler: func(message string) {
			fmt.Printf("TICKER: %s\n", message)
		},
		CandleHandler: func(message string) {
			fmt.Printf("CANDLE: %s\n", message)
		},
		OrderBookHandler: func(message string) {
			fmt.Printf("ORDERBOOK: %s\n", message)
		},
	}
}

// DefaultTradingStreamConfig returns a default configuration for trading streams
func DefaultTradingStreamConfig() TradingStreamConfig {
	return TradingStreamConfig{
		EnableOrders:    true,
		EnableFills:     true,
		EnablePositions: true,
		EnableAccount:   false,

		OrderHandler: func(message string) {
			fmt.Printf("ORDER: %s\n", message)
		},
		FillHandler: func(message string) {
			fmt.Printf("FILL: %s\n", message)
		},
		PositionHandler: func(message string) {
			fmt.Printf("POSITION: %s\n", message)
		},
	}
}
