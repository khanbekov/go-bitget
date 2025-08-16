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

// loginCredentials stores authentication details for re-authentication after reconnection
type loginCredentials struct {
	apiKey     string
	passphrase string
	signType   common.SignType
}

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
	reconnecting          bool                           // Flag to prevent multiple concurrent reconnection attempts
	reconnectMutex        sync.Mutex                     // Mutex for thread-safe reconnection
	maxReconnectAttempts  int                            // Maximum number of reconnection attempts
	reconnectAttempts     int                            // Current number of reconnection attempts
	storedLoginCreds      *loginCredentials              // Stored login credentials for re-authentication
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
		maxReconnectAttempts:  5, // Maximum 5 reconnection attempts before giving up
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
// Default is 120 seconds. Shorter timeouts provide faster recovery but may cause unnecessary reconnections.
func (c *BaseWsClient) SetReconnectionTimeout(timeout time.Duration) {
	c.reconnectionTimeout = timeout
}

// SetMaxReconnectAttempts sets the maximum number of reconnection attempts before giving up.
// Default is 5 attempts. Set to 0 for unlimited attempts.
func (c *BaseWsClient) SetMaxReconnectAttempts(maxAttempts int) {
	c.maxReconnectAttempts = maxAttempts
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
// Stores credentials for automatic re-authentication after reconnection.
//
// Parameters:
//   - apiKey: Your Bitget API key
//   - passphrase: Your API passphrase
//   - signType: Signature algorithm (SHA256 or RSA)
func (c *BaseWsClient) Login(apiKey, passphrase string, signType common.SignType) {
	// Store credentials for re-authentication after reconnection
	c.storedLoginCreds = &loginCredentials{
		apiKey:     apiKey,
		passphrase: passphrase,
		signType:   signType,
	}
	c.needLogin = true

	c.performLogin()
}

// performLogin executes the actual login process
func (c *BaseWsClient) performLogin() {
	if c.storedLoginCreds == nil {
		c.logger.Error().Msg("No stored login credentials available")
		return
	}

	timesStamp := common.TimestampSec()

	var sign string
	if c.storedLoginCreds.signType == common.SHA256 {
		sign = c.signer.Sign(common.WsAuthMethod, common.WsAuthPath, "", timesStamp)
	} else {
		sign = c.signer.SignByRSA(common.WsAuthMethod, common.WsAuthPath, "", timesStamp)
	}

	loginReq := WsLoginReq{
		ApiKey:     c.storedLoginCreds.apiKey,
		Passphrase: c.storedLoginCreds.passphrase,
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
	_, err := cr.AddFunc("*/30 * * * * *", c.ping)
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
			// Skip checks if already reconnecting
			if c.reconnecting {
				continue
			}

			elapsedSecond := time.Now().Sub(c.lastReceivedTime)
			connectionAge := time.Now().Sub(c.connectionStartTime)

			// Check for 24-hour force disconnect (as per WebSocket spec)
			if connectionAge > 24*time.Hour {
				c.logger.Info().Msg("24-hour limit reached, forcing WebSocket reconnection")
				go func() {
					if err := c.performReconnection(); err != nil {
						c.logger.Error().Err(err).Msg("Failed to perform 24-hour reconnection")
					}
				}()
				continue
			}

			// Check for message timeout
			if elapsedSecond > c.reconnectionTimeout {
				c.logger.Warn().Dur("elapsed", elapsedSecond).Msg("WebSocket reconnect due to timeout...")
				go func() {
					if err := c.performReconnection(); err != nil {
						c.logger.Error().Err(err).Msg("Failed to perform timeout reconnection")
					}
				}()
			}
		}
	}
}

// Reconnect manually triggers a WebSocket reconnection.
// This method can be called to force a reconnection in case of network issues.
// It includes exponential backoff and retry logic.
func (c *BaseWsClient) Reconnect() error {
	c.reconnectMutex.Lock()
	defer c.reconnectMutex.Unlock()

	if c.reconnecting {
		c.logger.Debug().Msg("Reconnection already in progress, skipping")
		return nil
	}

	c.logger.Info().Msg("Manual reconnection triggered")
	return c.performReconnection()
}

// performReconnection handles the actual reconnection logic with exponential backoff
func (c *BaseWsClient) performReconnection() error {
	c.reconnecting = true
	defer func() {
		c.reconnecting = false
	}()

	// Disconnect current connection
	c.disconnectWebSocket()
	c.connected = false
	c.loginStatus = false

	// Reset connection state
	c.reconnectAttempts = 0

	for {
		// Check if we've exceeded max attempts (0 means unlimited)
		if c.maxReconnectAttempts > 0 && c.reconnectAttempts >= c.maxReconnectAttempts {
			c.logger.Error().
				Int("max_attempts", c.maxReconnectAttempts).
				Msg("Maximum reconnection attempts reached, giving up")
			return fmt.Errorf("maximum reconnection attempts (%d) exceeded", c.maxReconnectAttempts)
		}

		c.reconnectAttempts++
		c.logger.Info().
			Int("attempt", c.reconnectAttempts).
			Int("max_attempts", c.maxReconnectAttempts).
			Msg("Attempting to reconnect WebSocket")

		// Try to reconnect
		err := c.attemptConnection()
		if err == nil {
			c.logger.Info().
				Int("attempts_used", c.reconnectAttempts).
				Msg("WebSocket reconnection successful")

			// Reset attempts counter on success
			c.reconnectAttempts = 0
			return nil
		}

		c.logger.Warn().
			Err(err).
			Int("attempt", c.reconnectAttempts).
			Msg("Reconnection attempt failed")

		// Exponential backoff: wait 2^attempt seconds (max 30 seconds)
		backoffDuration := time.Duration(1<<uint(c.reconnectAttempts)) * time.Second
		if backoffDuration > 30*time.Second {
			backoffDuration = 30 * time.Second
		}

		c.logger.Debug().
			Dur("backoff", backoffDuration).
			Msg("Waiting before next reconnection attempt")

		time.Sleep(backoffDuration)
	}
}

// attemptConnection tries to establish a new WebSocket connection
func (c *BaseWsClient) attemptConnection() error {
	c.logger.Debug().Str("url", c.url).Msg("Attempting WebSocket connection")

	var err error
	c.webSocketClient, _, err = websocket.DefaultDialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to dial WebSocket: %w", err)
	}

	c.connected = true
	c.connectionStartTime = time.Now()
	c.lastReceivedTime = time.Now()

	// Re-authenticate if needed
	if c.needLogin && c.storedLoginCreds != nil {
		c.logger.Info().Msg("Re-authenticating after reconnection")
		c.performLogin()

		// Wait a bit for authentication to complete
		time.Sleep(1 * time.Second)
	}

	// Restore subscriptions
	if len(c.subscriptions) > 0 {
		c.logger.Info().Int("subscription_count", len(c.subscriptions)).Msg("Restoring subscriptions after reconnection")
		c.restoreSubscriptions()
	}

	return nil
}

func (c *BaseWsClient) disconnectWebSocket() {
	defer func() {
		if r := recover(); r != nil {
			c.logger.Error().Interface("panic", r).Msg("Panic recovered during WebSocket disconnection")
		}
		// Always ensure these are set regardless of panic
		c.connected = false
		c.webSocketClient = nil
	}()

	if c.webSocketClient == nil {
		return
	}

	c.logger.Debug().Msg("WebSocket disconnecting...")
	c.connected = false

	err := c.webSocketClient.Close()
	if err != nil {
		c.logger.Warn().Err(err).Msg("WebSocket disconnect error")
	} else {
		c.logger.Debug().Msg("WebSocket disconnected successfully")
	}

	c.webSocketClient = nil
}

func (c *BaseWsClient) ReadLoop() {
	for {

		if c.webSocketClient == nil {
			c.logger.Error().Msg("error on message read: no connection available")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		_, buf, err := c.webSocketClient.ReadMessage()
		if err != nil {
			c.logger.Warn().Err(err).Str("msg", string(buf)).Msg("error on message read")

			// Handle different types of connection errors
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseAbnormalClosure,
				websocket.CloseGoingAway,
				websocket.CloseAbnormalClosure,
				websocket.CloseNoStatusReceived,
				websocket.CloseProtocolError,
				websocket.CloseInternalServerErr,
				websocket.CloseServiceRestart,
			) {
				c.logger.Info().Msg("WebSocket closed, attempting reconnection")

				// Use improved reconnection logic
				if err := c.performReconnection(); err != nil {
					c.logger.Error().Err(err).Msg("Failed to reconnect after close error")
				}
			}
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
	if c.needLogin && c.storedLoginCreds != nil {
		c.logger.Info().Msg("Re-authenticating private WebSocket after reconnection")
		c.performLogin()

		// Wait for authentication to complete
		time.Sleep(500 * time.Millisecond)
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
