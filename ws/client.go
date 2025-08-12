// Package ws provides WebSocket client functionality for real-time data streaming from Bitget.
// It supports both public market data and private authenticated channels with automatic
// reconnection, subscription management, and message handling.
//
// Example usage:
//
//	client := NewBitgetBaseWsClient(logger, "wss://ws.bitget.com/v2/ws/public", "")
//	client.SetListener(messageHandler, errorHandler)
//	client.Connect()
package ws

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/common/types"
	"github.com/rs/zerolog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robfig/cron/v3"
)

// OnReceive is a callback function type for handling incoming WebSocket messages.
// It receives the raw message string from the WebSocket connection.
type OnReceive func(message string)

// rateLimiter implements rate limiting to prevent exceeding message send limits
type rateLimiter struct {
	lastSend    time.Time     // Timestamp of last message sent
	minInterval time.Duration // Minimum interval between messages (100ms for 10 messages/second)
	mutex       sync.Mutex    // Mutex for thread-safe access
}

// BaseWsClient provides a WebSocket client for Bitget's real-time API.
// It handles connection management, authentication, subscription tracking,
// and automatic reconnection with configurable timeouts.
type BaseWsClient struct {
	needLogin             bool                           // Whether authentication is required
	connected             bool                           // Current connection status
	loginStatus           bool                           // Authentication status
	url                   string                         // WebSocket endpoint URL
	logger                zerolog.Logger                 // Logger for debugging and monitoring
	listener              OnReceive                      // Default message handler
	errorListener         OnReceive                      // Error message handler
	checkConnectionTicker *time.Ticker                   // Timer for connection health checks
	reconnectionTimeout   time.Duration                  // Timeout before attempting reconnection
	sendMutex             *sync.Mutex                    // Mutex for thread-safe message sending
	webSocketClient       *websocket.Conn                // Underlying WebSocket connection
	lastReceivedTime      time.Time                      // Timestamp of last received message
	connectionStartTime   time.Time                      // Timestamp when connection was established
	subscribeRequests     *types.Set                     // Set of active subscription requests
	signer                *common.Signer                 // Signer for authentication
	subscriptions         map[SubscriptionArgs]OnReceive // Map of subscriptions to their handlers
	rateLimiter           *rateLimiter                   // Rate limiter for message sending
}

// NewBitgetBaseWsClient creates a new WebSocket client for Bitget's real-time API.
//
// Parameters:
//   - logger: Logger instance for debugging and monitoring
//   - url: WebSocket endpoint URL (public or private)
//   - secretKey: Secret key for authentication (empty string for public channels)
//
// Returns a configured BaseWsClient ready for connection.
func NewBitgetBaseWsClient(logger zerolog.Logger, url, secretKey string) *BaseWsClient {
	return &BaseWsClient{
		logger:                logger,
		url:                   url,
		subscribeRequests:     types.NewSet(),
		signer:                common.NewSigner(secretKey),
		subscriptions:         make(map[SubscriptionArgs]OnReceive),
		sendMutex:             &sync.Mutex{},
		checkConnectionTicker: time.NewTicker(5 * time.Second),
		reconnectionTimeout:   120 * time.Second, // Increased from 60s to 120s for better stability
		lastReceivedTime:      time.Now(),
		connectionStartTime:   time.Now(),
		rateLimiter: &rateLimiter{
			minInterval: 100 * time.Millisecond, // 10 messages per second max
		},
	}
}

// SetCheckConnectionInterval configures how often the client checks connection health.
// Default is 5 seconds. Lower values provide faster reconnection but more overhead.
func (c *BaseWsClient) SetCheckConnectionInterval(interval time.Duration) {
	c.checkConnectionTicker = time.NewTicker(interval)
}

// SetReconnectionTimeout sets how long to wait without receiving messages before reconnecting.
// Default is 60 seconds. Shorter timeouts provide faster recovery but may cause unnecessary reconnections.
func (c *BaseWsClient) SetReconnectionTimeout(timeout time.Duration) {
	c.reconnectionTimeout = timeout
}

// SetListener sets the default message and error handlers for the WebSocket client.
//
// Parameters:
//   - msgListener: Function to handle incoming data messages
//   - errorListener: Function to handle error messages and API errors
func (c *BaseWsClient) SetListener(msgListener OnReceive, errorListener OnReceive) {
	c.listener = msgListener
	c.errorListener = errorListener
}

// Connect initiates the WebSocket connection and starts the monitoring loop.
// This method starts the connection health checker and ping mechanism.
func (c *BaseWsClient) Connect() {

	go c.tickerLoop() // Run ticker loop in background goroutine
	err := c.startPing()
	if err != nil {
		c.logger.Error().Err(err).Msg("fail to start ping")
		return
	}
}

// ConnectWebSocket establishes the actual WebSocket connection to the Bitget server.
// This method is called internally by Connect() and during reconnection attempts.
func (c *BaseWsClient) ConnectWebSocket() {
	var err error
	c.logger.Info().Msg("WebSocket connecting...")
	c.webSocketClient, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		fmt.Printf("WebSocket connected error: %s\n", err)
		return
	}
	c.logger.Info().Msg("WebSocket connected")
	c.connected = true
	c.connectionStartTime = time.Now() // Reset connection start time
	c.lastReceivedTime = time.Now()    // Reset last received time to prevent immediate timeout

	// Restore subscriptions after reconnection
	if len(c.subscriptions) > 0 {
		c.logger.Info().Int("subscription_count", len(c.subscriptions)).Msg("Restoring subscriptions after reconnection")
		c.restoreSubscriptions()
	}
}

// Login performs authentication for private WebSocket channels.
// Required for accessing account-specific data like positions and orders.
//
// Parameters:
//   - apiKey: Your Bitget API key
//   - passphrase: Your API passphrase
//   - signType: Signature algorithm (SHA256 or RSA)
func (c *BaseWsClient) Login(apiKey, passphrase string, signType common.SignType) {
	timesStamp := common.TimestampSec()

	var sign string
	if signType == common.SHA256 {
		sign = c.signer.Sign(common.WsAuthMethod, common.WsAuthPath, "", timesStamp)
	} else {
		sign = c.signer.SignByRSA(common.WsAuthMethod, common.WsAuthPath, "", timesStamp)
	}

	loginReq := WsLoginReq{
		ApiKey:     apiKey,
		Passphrase: passphrase,
		Timestamp:  timesStamp,
		Sign:       sign,
	}
	var args []interface{}
	args = append(args, loginReq)

	baseReq := WsBaseReq{
		Op:   common.WsOpLogin,
		Args: args,
	}
	c.SendByType(baseReq)
}

func (c *BaseWsClient) StartReadLoop() {
	go c.ReadLoop()
}

func (c *BaseWsClient) startPing() error {
	cr := cron.New(cron.WithSeconds()) // Enable seconds field
	_, err := cr.AddFunc("*/15 * * * * *", c.ping)
	if err != nil {
		return err
	}
	cr.Start()
	return nil
}
func (c *BaseWsClient) ping() {
	c.Send("ping")
}

func (c *BaseWsClient) SendByType(req WsBaseReq) {
	json, _ := jsoniter.MarshalToString(req)
	c.Send(json)
}

func (c *BaseWsClient) Send(data string) {
	if c.webSocketClient == nil {
		c.logger.Error().Msg("WebSocket sent error: no connection available")
		return
	}

	// Apply rate limiting (max 10 messages per second)
	c.rateLimiter.mutex.Lock()
	timeSinceLastSend := time.Since(c.rateLimiter.lastSend)
	if timeSinceLastSend < c.rateLimiter.minInterval {
		sleepDuration := c.rateLimiter.minInterval - timeSinceLastSend
		c.rateLimiter.mutex.Unlock()
		c.logger.Debug().Dur("sleep", sleepDuration).Msg("Rate limiting: sleeping before send")
		time.Sleep(sleepDuration)
		c.rateLimiter.mutex.Lock()
	}
	c.rateLimiter.lastSend = time.Now()
	c.rateLimiter.mutex.Unlock()

	c.logger.Debug().Str("message", data).Msg("send message")
	c.sendMutex.Lock()
	err := c.webSocketClient.WriteMessage(websocket.TextMessage, []byte(data))
	c.sendMutex.Unlock()
	if err != nil {
		c.logger.Error().Err(err).Str("message", data).Msg("failed to send message to websocket")
	}
}

func (c *BaseWsClient) tickerLoop() {
	c.logger.Info().Msg("tickerLoop started")
	for {
		select {
		case <-c.checkConnectionTicker.C:
			elapsedSecond := time.Now().Sub(c.lastReceivedTime)
			connectionAge := time.Now().Sub(c.connectionStartTime)

			// Check for 24-hour force disconnect (as per WebSocket spec)
			if connectionAge > 24*time.Hour {
				c.logger.Info().Msg("24-hour limit reached, forcing WebSocket reconnection")
				c.disconnectWebSocket()
				c.ConnectWebSocket()
				continue
			}

			// Check for message timeout
			if elapsedSecond > c.reconnectionTimeout {
				c.logger.Warn().Dur("elapsed", elapsedSecond).Msg("WebSocket reconnect due to timeout...")
				c.disconnectWebSocket()
				c.ConnectWebSocket()
			}
		}
	}
}

func (c *BaseWsClient) disconnectWebSocket() {
	if c.webSocketClient == nil {
		return
	}

	fmt.Println("WebSocket disconnecting...")
	err := c.webSocketClient.Close()
	if err != nil {
		c.logger.Warn().Err(err).Msg("WebSocket disconnect error")
		return
	}

	c.logger.Info().Msg("WebSocket disconnected")
}

func (c *BaseWsClient) ReadLoop() {
	for {

		if c.webSocketClient == nil {
			c.logger.Error().Msg("error on message read: no connection available")
			//time.Sleep(TimerIntervalSecond * time.Second)
			continue
		}

		_, buf, err := c.webSocketClient.ReadMessage()
		if err != nil {
			c.logger.Warn().Err(err).Msg("error on message read")
			continue
		}
		c.lastReceivedTime = time.Now()
		message := string(buf)

		if message == "pong" {
			c.logger.Debug().Str("message", message).Msg("keep connected")
			continue
		}
		c.logger.Debug().Str("message", message).Msg("read message from websocket")

		jsonMap := make(map[string]interface{})
		err = jsoniter.Unmarshal(buf, &jsonMap)
		if err != nil {
			c.logger.Warn().Err(err).Msg("error on umarshalling message")
			continue
		}

		v, e := jsonMap["code"]

		if e {
			code, ok := v.(float64)
			if !ok || code != 0 {
				c.errorListener(message)
				continue
			}
		}

		v, e = jsonMap["event"]
		if e && v == "login" {
			c.logger.Debug().Str("message", message).Msg("login")
			c.loginStatus = true
			continue
		}

		v, e = jsonMap["data"]
		if e {
			listener := c.GetListener(jsonMap["arg"])
			listener(message)
			continue
		}
	}

}

func (c *BaseWsClient) GetListener(argJson interface{}) OnReceive {

	mapData := argJson.(map[string]interface{})

	subscribeReq := SubscriptionArgs{
		ProductType: fmt.Sprintf("%v", mapData["instType"]),
		Channel:     fmt.Sprintf("%v", mapData["channel"]),
		Symbol:      fmt.Sprintf("%v", mapData["instId"]),
	}

	v, e := c.subscriptions[subscribeReq]

	if !e {
		return c.listener
	}
	return v
}

// IsConnected returns true if the WebSocket connection is established and active
func (c *BaseWsClient) IsConnected() bool {
	return c.connected && c.webSocketClient != nil
}

// IsLoggedIn returns true if the WebSocket is authenticated (for private channels)
func (c *BaseWsClient) IsLoggedIn() bool {
	return c.loginStatus
}

// GetSubscriptionCount returns the number of active subscriptions
func (c *BaseWsClient) GetSubscriptionCount() int {
	return len(c.subscriptions)
}

// restoreSubscriptions resubscribes to all previously active subscriptions after reconnection
func (c *BaseWsClient) restoreSubscriptions() {
	c.logger.Info().Msg("Starting subscription restoration after reconnection")

	// Create a copy of current subscriptions to avoid map iteration issues
	var subscriptionsToRestore []SubscriptionArgs
	for args := range c.subscriptions {
		subscriptionsToRestore = append(subscriptionsToRestore, args)
	}

	// Wait a bit for the connection to stabilize
	time.Sleep(500 * time.Millisecond)

	// Re-authenticate if this is a private WebSocket
	if c.needLogin {
		c.logger.Info().Msg("Re-authenticating private WebSocket after reconnection")
		// Note: Login credentials should be stored if we need to re-authenticate
		// For now, we'll assume the parent connector will handle re-authentication
	}

	// Restore each subscription
	restoredCount := 0
	for _, args := range subscriptionsToRestore {
		c.logger.Debug().
			Str("channel", args.Channel).
			Str("symbol", args.Symbol).
			Str("productType", args.ProductType).
			Msg("Restoring subscription")

		// Use subscribe method to restore the subscription
		c.subscribe(args)
		restoredCount++

		// Small delay between subscriptions to avoid rate limiting
		time.Sleep(100 * time.Millisecond)
	}

	c.logger.Info().
		Int("restored_subscriptions", restoredCount).
		Int("total_subscriptions", len(c.subscriptions)).
		Msg("Subscription restoration completed")
}

func (c *BaseWsClient) Close() {
	if c.connected {
		cm := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close")

		if err := c.webSocketClient.WriteMessage(websocket.CloseMessage, cm); err != nil {
			c.logger.Error().Err(err).Msg("WebSocket disconnection error")
		}
		c.disconnectWebSocket()
	}
}
