// Package ws provides WebSocket data type abstractions based on Bitget API documentation
package ws

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// =============================================================================
// BASE WEBSOCKET MESSAGE STRUCTURES
// =============================================================================

// WebSocketMessage represents the base structure for all WebSocket messages
type WebSocketMessage struct {
	Event  string           `json:"event,omitempty"`
	Action string           `json:"action,omitempty"`
	Arg    SubscriptionArgs `json:"arg"`
	Data   json.RawMessage  `json:"data,omitempty"`
	TS     int64            `json:"ts"`
	Code   string           `json:"code,omitempty"`
	Msg    string           `json:"msg,omitempty"`
}

// SubscriptionArgs represents the subscription argument structure used throughout the codebase
type SubscriptionArgs struct {
	ProductType string `json:"instType"`
	Channel     string `json:"channel"`
	Symbol      string `json:"instId,omitempty"`
	Coin        string `json:"coin,omitempty"`
}

// SubscriptionRequest represents a WebSocket subscription request
type SubscriptionRequest struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

// WsBaseReq represents a base WebSocket request structure
type WsBaseReq struct {
	Op   string        `json:"op"`
	Args []interface{} `json:"args"`
}

// WsLoginBaseReq represents a login request wrapper
type WsLoginBaseReq struct {
	Op   string       `json:"op"`
	Args []WsLoginReq `json:"args"`
}

// WsLoginReq represents login request parameters
type WsLoginReq struct {
	ApiKey     string `json:"apiKey"`
	Passphrase string `json:"passphrase"`
	Timestamp  string `json:"timestamp"`
	Sign       string `json:"sign"`
}

// =============================================================================
// TICKER DATA ABSTRACTION
// =============================================================================

// TickerData represents ticker market data based on official documentation
type TickerData struct {
	InstId          string `json:"instId"`          // Product ID, e.g., BTCUSDT
	LastPrice       string `json:"lastPr"`          // Latest price
	BidPrice        string `json:"bidPr"`           // Bid price
	AskPrice        string `json:"askPr"`           // Ask price
	BidSize         string `json:"bidSz"`           // Buying amount
	AskSize         string `json:"askSz"`           // Selling amount
	Open24h         string `json:"open24h"`         // Opening price 24h ago
	High24h         string `json:"high24h"`         // 24h high
	Low24h          string `json:"low24h"`          // 24h low
	Change24h       string `json:"change24h"`       // 24h change
	FundingRate     string `json:"fundingRate"`     // Funding rate
	NextFundingTime string `json:"nextFundingTime"` // Next funding time (timestamp)
	MarkPrice       string `json:"markPrice"`       // Mark price
	IndexPrice      string `json:"indexPrice"`      // Index price
	HoldingAmount   string `json:"holdingAmount"`   // Open interest
	BaseVolume      string `json:"baseVolume"`      // Trading volume of the coin
	QuoteVolume     string `json:"quoteVolume"`     // Trading volume of quote currency
	OpenUtc         string `json:"openUtc"`         // Price at 00:00 (UTC)
	SymbolType      int    `json:"symbolType"`      // 1=perpetual, 2=delivery
	Symbol          string `json:"symbol"`          // Trading pair
	DeliveryPrice   string `json:"deliveryPrice"`   // Delivery price
	Timestamp       string `json:"ts"`              // System timestamp

	// Parsed fields for convenience
	LastPriceFloat      float64   `json:"-"`
	BidPriceFloat       float64   `json:"-"`
	AskPriceFloat       float64   `json:"-"`
	Change24hFloat      float64   `json:"-"`
	NextFundingTimeDate time.Time `json:"-"`
	TimestampDate       time.Time `json:"-"`
}

// ParseFloats parses string fields to float64 for easier calculations
func (t *TickerData) ParseFloats() error {
	if t.LastPrice != "" {
		if val, err := strconv.ParseFloat(t.LastPrice, 64); err != nil {
			return fmt.Errorf("failed to parse lastPrice: %w", err)
		} else {
			t.LastPriceFloat = val
		}
	}

	if t.BidPrice != "" {
		if val, err := strconv.ParseFloat(t.BidPrice, 64); err != nil {
			return fmt.Errorf("failed to parse bidPrice: %w", err)
		} else {
			t.BidPriceFloat = val
		}
	}

	if t.AskPrice != "" {
		if val, err := strconv.ParseFloat(t.AskPrice, 64); err != nil {
			return fmt.Errorf("failed to parse askPrice: %w", err)
		} else {
			t.AskPriceFloat = val
		}
	}

	if t.Change24h != "" {
		if val, err := strconv.ParseFloat(t.Change24h, 64); err != nil {
			return fmt.Errorf("failed to parse change24h: %w", err)
		} else {
			t.Change24hFloat = val
		}
	}

	return nil
}

// ParseTimestamps parses timestamp fields to time.Time
func (t *TickerData) ParseTimestamps() error {
	if t.NextFundingTime != "" {
		if timestamp, err := strconv.ParseInt(t.NextFundingTime, 10, 64); err != nil {
			return fmt.Errorf("failed to parse nextFundingTime: %w", err)
		} else {
			t.NextFundingTimeDate = time.UnixMilli(timestamp)
		}
	}

	if t.Timestamp != "" {
		if timestamp, err := strconv.ParseInt(t.Timestamp, 10, 64); err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		} else {
			t.TimestampDate = time.UnixMilli(timestamp)
		}
	}

	return nil
}

// Spread calculates the bid-ask spread
func (t *TickerData) Spread() float64 {
	return t.AskPriceFloat - t.BidPriceFloat
}

// SpreadPercentage calculates the spread as a percentage of mid price
func (t *TickerData) SpreadPercentage() float64 {
	mid := (t.BidPriceFloat + t.AskPriceFloat) / 2
	if mid == 0 {
		return 0
	}
	return (t.Spread() / mid) * 100
}

// =============================================================================
// CANDLESTICK DATA ABSTRACTION
// =============================================================================

// CandlestickData represents candlestick data based on official documentation
// Data is received as array: [timestamp, open, high, low, close, baseVolume, quoteVolume, usdtVolume]
type CandlestickData struct {
	Timestamp   string `json:"-"` // Start time (milliseconds)
	Open        string `json:"-"` // Opening price
	High        string `json:"-"` // Highest price
	Low         string `json:"-"` // Lowest price
	Close       string `json:"-"` // Closing price
	BaseVolume  string `json:"-"` // Trading volume of left coin
	QuoteVolume string `json:"-"` // Trading volume of quote currency
	UsdtVolume  string `json:"-"` // Trading volume of USDT

	// Parsed fields for convenience
	TimestampDate    time.Time `json:"-"`
	OpenFloat        float64   `json:"-"`
	HighFloat        float64   `json:"-"`
	LowFloat         float64   `json:"-"`
	CloseFloat       float64   `json:"-"`
	BaseVolumeFloat  float64   `json:"-"`
	QuoteVolumeFloat float64   `json:"-"`
	UsdtVolumeFloat  float64   `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling for candlestick array data
func (c *CandlestickData) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) < 8 {
		return fmt.Errorf("candlestick data array must have at least 8 elements, got %d", len(arr))
	}

	c.Timestamp = arr[0]
	c.Open = arr[1]
	c.High = arr[2]
	c.Low = arr[3]
	c.Close = arr[4]
	c.BaseVolume = arr[5]
	c.QuoteVolume = arr[6]
	c.UsdtVolume = arr[7]

	return c.ParseAll()
}

// ParseAll parses all string fields to appropriate types
func (c *CandlestickData) ParseAll() error {
	var err error

	// Parse timestamp
	if c.Timestamp != "" {
		timestamp, err := strconv.ParseInt(c.Timestamp, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
		c.TimestampDate = time.UnixMilli(timestamp)
	}

	// Parse OHLC values
	if c.Open != "" {
		c.OpenFloat, err = strconv.ParseFloat(c.Open, 64)
		if err != nil {
			return fmt.Errorf("failed to parse open: %w", err)
		}
	}

	if c.High != "" {
		c.HighFloat, err = strconv.ParseFloat(c.High, 64)
		if err != nil {
			return fmt.Errorf("failed to parse high: %w", err)
		}
	}

	if c.Low != "" {
		c.LowFloat, err = strconv.ParseFloat(c.Low, 64)
		if err != nil {
			return fmt.Errorf("failed to parse low: %w", err)
		}
	}

	if c.Close != "" {
		c.CloseFloat, err = strconv.ParseFloat(c.Close, 64)
		if err != nil {
			return fmt.Errorf("failed to parse close: %w", err)
		}
	}

	// Parse volume values
	if c.BaseVolume != "" {
		c.BaseVolumeFloat, err = strconv.ParseFloat(c.BaseVolume, 64)
		if err != nil {
			return fmt.Errorf("failed to parse baseVolume: %w", err)
		}
	}

	if c.QuoteVolume != "" {
		c.QuoteVolumeFloat, err = strconv.ParseFloat(c.QuoteVolume, 64)
		if err != nil {
			return fmt.Errorf("failed to parse quoteVolume: %w", err)
		}
	}

	if c.UsdtVolume != "" {
		c.UsdtVolumeFloat, err = strconv.ParseFloat(c.UsdtVolume, 64)
		if err != nil {
			return fmt.Errorf("failed to parse usdtVolume: %w", err)
		}
	}

	return nil
}

// BodyRange calculates the body range (close - open)
func (c *CandlestickData) BodyRange() float64 {
	return c.CloseFloat - c.OpenFloat
}

// BodyRangePercentage calculates the body range as percentage of open price
func (c *CandlestickData) BodyRangePercentage() float64 {
	if c.OpenFloat == 0 {
		return 0
	}
	return (c.BodyRange() / c.OpenFloat) * 100
}

// WickRange calculates the total wick range (high - low - body)
func (c *CandlestickData) WickRange() float64 {
	return (c.HighFloat - c.LowFloat) - abs(c.BodyRange())
}

// IsBullish returns true if the candlestick is bullish (close > open)
func (c *CandlestickData) IsBullish() bool {
	return c.CloseFloat > c.OpenFloat
}

// IsBearish returns true if the candlestick is bearish (close < open)
func (c *CandlestickData) IsBearish() bool {
	return c.CloseFloat < c.OpenFloat
}

// =============================================================================
// ORDER BOOK DATA ABSTRACTION
// =============================================================================

// OrderBookLevel represents a single price level in the order book
type OrderBookLevel struct {
	Price  string `json:"-"`
	Amount string `json:"-"`

	// Parsed fields
	PriceFloat  float64 `json:"-"`
	AmountFloat float64 `json:"-"`
}

// UnmarshalJSON implements custom JSON unmarshaling for order book level array
func (l *OrderBookLevel) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) < 2 {
		return fmt.Errorf("order book level must have price and amount, got %d elements", len(arr))
	}

	l.Price = arr[0]
	l.Amount = arr[1]

	return l.ParseFloats()
}

// ParseFloats parses price and amount to float64
func (l *OrderBookLevel) ParseFloats() error {
	var err error

	if l.Price != "" {
		l.PriceFloat, err = strconv.ParseFloat(l.Price, 64)
		if err != nil {
			return fmt.Errorf("failed to parse price: %w", err)
		}
	}

	if l.Amount != "" {
		l.AmountFloat, err = strconv.ParseFloat(l.Amount, 64)
		if err != nil {
			return fmt.Errorf("failed to parse amount: %w", err)
		}
	}

	return nil
}

// OrderBookData represents order book data based on official documentation
type OrderBookData struct {
	Asks     []OrderBookLevel `json:"asks"`     // Seller depth
	Bids     []OrderBookLevel `json:"bids"`     // Buyer depth
	Checksum int64            `json:"checksum"` // Checksum for data integrity
	Seq      int64            `json:"seq"`      // Sequence number
	TS       string           `json:"ts"`       // Match engine timestamp

	// Parsed fields
	TimestampDate time.Time `json:"-"`
}

// ParseTimestamp parses the timestamp field
func (o *OrderBookData) ParseTimestamp() error {
	if o.TS != "" {
		timestamp, err := strconv.ParseInt(o.TS, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
		o.TimestampDate = time.UnixMilli(timestamp)
	}
	return nil
}

// BestBid returns the best bid price and amount
func (o *OrderBookData) BestBid() (float64, float64) {
	if len(o.Bids) == 0 {
		return 0, 0
	}
	return o.Bids[0].PriceFloat, o.Bids[0].AmountFloat
}

// BestAsk returns the best ask price and amount
func (o *OrderBookData) BestAsk() (float64, float64) {
	if len(o.Asks) == 0 {
		return 0, 0
	}
	return o.Asks[0].PriceFloat, o.Asks[0].AmountFloat
}

// Spread returns the bid-ask spread
func (o *OrderBookData) Spread() float64 {
	bestBidPrice, _ := o.BestBid()
	bestAskPrice, _ := o.BestAsk()
	return bestAskPrice - bestBidPrice
}

// MidPrice returns the mid price between best bid and ask
func (o *OrderBookData) MidPrice() float64 {
	bestBidPrice, _ := o.BestBid()
	bestAskPrice, _ := o.BestAsk()
	return (bestBidPrice + bestAskPrice) / 2
}

// TotalBidVolume calculates the total volume on the bid side
func (o *OrderBookData) TotalBidVolume() float64 {
	total := 0.0
	for _, bid := range o.Bids {
		total += bid.AmountFloat
	}
	return total
}

// TotalAskVolume calculates the total volume on the ask side
func (o *OrderBookData) TotalAskVolume() float64 {
	total := 0.0
	for _, ask := range o.Asks {
		total += ask.AmountFloat
	}
	return total
}

// =============================================================================
// TRADE DATA ABSTRACTION
// =============================================================================

// TradeData represents public trade data based on official documentation
type TradeData struct {
	TS      string `json:"ts"`      // Fill time
	Price   string `json:"price"`   // Filled price
	Size    string `json:"size"`    // Filled amount
	Side    string `json:"side"`    // Filled side: sell/buy
	TradeId string `json:"tradeId"` // Trade ID

	// Parsed fields
	TimestampDate time.Time `json:"-"`
	PriceFloat    float64   `json:"-"`
	SizeFloat     float64   `json:"-"`
	IsBuy         bool      `json:"-"`
	IsSell        bool      `json:"-"`
}

// ParseAll parses all string fields to appropriate types
func (t *TradeData) ParseAll() error {
	var err error

	// Parse timestamp
	if t.TS != "" {
		timestamp, err := strconv.ParseInt(t.TS, 10, 64)
		if err != nil {
			return fmt.Errorf("failed to parse timestamp: %w", err)
		}
		t.TimestampDate = time.UnixMilli(timestamp)
	}

	// Parse price
	if t.Price != "" {
		t.PriceFloat, err = strconv.ParseFloat(t.Price, 64)
		if err != nil {
			return fmt.Errorf("failed to parse price: %w", err)
		}
	}

	// Parse size
	if t.Size != "" {
		t.SizeFloat, err = strconv.ParseFloat(t.Size, 64)
		if err != nil {
			return fmt.Errorf("failed to parse size: %w", err)
		}
	}

	// Parse side
	t.IsBuy = t.Side == "buy"
	t.IsSell = t.Side == "sell"

	return nil
}

// Value calculates the trade value (price * size)
func (t *TradeData) Value() float64 {
	return t.PriceFloat * t.SizeFloat
}

// =============================================================================
// HELPER FUNCTIONS
// =============================================================================

// abs returns the absolute value of x
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// =============================================================================
// NEW CHANNEL CONSTANTS BASED ON DOCUMENTATION
// =============================================================================

// Product types based on documentation
const (
	ProductTypeUSDTFutures = "USDT-FUTURES"
	ProductTypeCOINFutures = "COIN-FUTURES"
	ProductTypeUSDCFutures = "USDC-FUTURES"
)

// WebSocket operations
const (
	OpSubscribe   = "subscribe"
	OpUnsubscribe = "unsubscribe"
)

// Action types for push data
const (
	ActionSnapshot = "snapshot"
	ActionUpdate   = "update"
)
