package uta

import (
	"encoding/json"
	"time"
)

// ApiResponse represents the standard UTA API response structure
type ApiResponse struct {
	Code        string          `json:"code"`
	Msg         string          `json:"msg"`
	RequestTime int64           `json:"requestTime"`
	Data        json.RawMessage `json:"data"`
}

// Account information structures

// AccountInfo represents account settings and configuration
type AccountInfo struct {
	AssetMode    string         `json:"assetMode"`
	HoldingMode  string         `json:"holdingMode"`
	STPMode      string         `json:"stpMode"`
	SymbolConfig []SymbolConfig `json:"symbolConfig"`
	CoinConfig   []CoinConfig   `json:"coinConfig"`
}

// SymbolConfig represents symbol-specific configuration
type SymbolConfig struct {
	Symbol     string `json:"symbol"`
	Category   string `json:"category"`
	Leverage   string `json:"leverage"`
	MarginMode string `json:"marginMode"`
}

// CoinConfig represents coin-specific configuration
type CoinConfig struct {
	Coin       string `json:"coin"`
	Category   string `json:"category"`
	Leverage   string `json:"leverage"`
	MarginMode string `json:"marginMode"`
}

// AccountAssets represents account asset information
type AccountAssets struct {
	AccountEquity         string      `json:"accountEquity"`
	UnrealizedPNL         string      `json:"unrealizedPNL"`
	EffectiveEquity       string      `json:"effectiveEquity"`
	CrossMarginRatio      string      `json:"crossMarginRatio"`
	IsolatedMarginRatio   string      `json:"isolatedMarginRatio"`
	CrossMaintenanceRatio string      `json:"crossMaintenanceRatio"`
	Assets                []AssetInfo `json:"assets"`
}

// AssetInfo represents individual asset information
type AssetInfo struct {
	Coin              string `json:"coin"`
	Available         string `json:"available"`
	Frozen            string `json:"frozen"`
	Balance           string `json:"balance"`
	UnrealizedPNL     string `json:"unrealizedPNL"`
	CrossMarginAssets string `json:"crossMarginAssets"`
	IsolatedBalance   string `json:"isolatedBalance"`
	BorrowAmount      string `json:"borrowAmount"`
	AccruedInterest   string `json:"accruedInterest"`
	NetAssets         string `json:"netAssets"`
	NetAssetsUSD      string `json:"netAssetsUSD"`
}

// FundingAssets represents funding account assets
type FundingAssets struct {
	Coin      string `json:"coin"`
	Available string `json:"available"`
	Frozen    string `json:"frozen"`
	Balance   string `json:"balance"`
}

// FeeRate represents trading fee rates
type FeeRate struct {
	Symbol    string `json:"symbol"`
	Category  string `json:"category"`
	MakerRate string `json:"makerRate"`
	TakerRate string `json:"takerRate"`
}

// Transfer structures

// TransferResult represents transfer operation result
type TransferResult struct {
	TransferID string `json:"transferId"`
	ClientOid  string `json:"clientOid,omitempty"`
}

// TransferRecord represents transfer history record
type TransferRecord struct {
	TransferID string `json:"transferId"`
	ClientOid  string `json:"clientOid"`
	FromType   string `json:"fromType"`
	ToType     string `json:"toType"`
	Amount     string `json:"amount"`
	Coin       string `json:"coin"`
	Symbol     string `json:"symbol,omitempty"`
	FromUserId string `json:"fromUserId,omitempty"`
	ToUserId   string `json:"toUserId,omitempty"`
	Status     string `json:"status"`
	Timestamp  string `json:"timestamp"`
}

// Deposit and withdrawal structures

// DepositAddress represents deposit address information
type DepositAddress struct {
	Address string `json:"address"`
	Chain   string `json:"chain"`
	Coin    string `json:"coin"`
	Tag     string `json:"tag,omitempty"`
	URL     string `json:"url,omitempty"`
}

// DepositRecord represents deposit history record
type DepositRecord struct {
	OrderID    string `json:"orderId"`
	TrxID      string `json:"trxId"`
	Coin       string `json:"coin"`
	Chain      string `json:"chain"`
	Amount     string `json:"amount"`
	Status     string `json:"status"`
	ToAddress  string `json:"toAddress"`
	ConfirmNum string `json:"confirmNum"`
	Timestamp  string `json:"timestamp"`
}

// WithdrawalResult represents withdrawal operation result
type WithdrawalResult struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid,omitempty"`
}

// WithdrawalRecord represents withdrawal history record
type WithdrawalRecord struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	Coin      string `json:"coin"`
	Chain     string `json:"chain"`
	Amount    string `json:"amount"`
	Fee       string `json:"fee"`
	Status    string `json:"status"`
	ToAddress string `json:"toAddress"`
	TrxID     string `json:"trxId,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Sub-account structures

// SubAccount represents sub-account information
type SubAccount struct {
	Username    string `json:"username"`
	SubUid      string `json:"subUid"`
	Status      string `json:"status"`
	Note        string `json:"note"`
	CreatedTime string `json:"createdTime"`
	UpdatedTime string `json:"updatedTime"`
}

// SubAccountAssets represents sub-account asset information
type SubAccountAssets struct {
	SubUid          string      `json:"subUid"`
	AccountEquity   string      `json:"accountEquity"`
	UnrealizedPNL   string      `json:"unrealizedPNL"`
	EffectiveEquity string      `json:"effectiveEquity"`
	Assets          []AssetInfo `json:"assets"`
}

// SubAccountAPIKey represents sub-account API key information
type SubAccountAPIKey struct {
	APIKey      string   `json:"apiKey"`
	Note        string   `json:"note"`
	Type        string   `json:"type"`
	Permissions []string `json:"permissions"`
	IPs         []string `json:"ips"`
	Status      string   `json:"status"`
	CreatedTime string   `json:"createdTime"`
	UpdatedTime string   `json:"updatedTime"`
}

// Trading structures

// Order represents order information
type Order struct {
	OrderID      string `json:"orderId"`
	ClientOid    string `json:"clientOid"`
	Symbol       string `json:"symbol"`
	Category     string `json:"category"`
	Side         string `json:"side"`
	OrderType    string `json:"orderType"`
	Price        string `json:"price,omitempty"`
	Size         string `json:"size"`
	FilledSize   string `json:"filledSize"`
	FilledAmount string `json:"filledAmount"`
	AvgPrice     string `json:"avgPrice"`
	Status       string `json:"status"`
	TimeInForce  string `json:"timeInForce,omitempty"`
	ReduceOnly   string `json:"reduceOnly,omitempty"`
	PositionSide string `json:"positionSide,omitempty"`
	STP          string `json:"stp,omitempty"`
	CreatedTime  string `json:"createdTime"`
	UpdatedTime  string `json:"updatedTime"`
}

// BatchOrderResult represents batch order operation result
type BatchOrderResult struct {
	OrderID   string `json:"orderId"`
	ClientOid string `json:"clientOid"`
	Code      string `json:"code"`
	Msg       string `json:"msg"`
}

// Fill represents trade fill information
type Fill struct {
	FillID     string `json:"fillId"`
	OrderID    string `json:"orderId"`
	ClientOid  string `json:"clientOid"`
	Symbol     string `json:"symbol"`
	Category   string `json:"category"`
	Side       string `json:"side"`
	FillPrice  string `json:"fillPrice"`
	FillSize   string `json:"fillSize"`
	FillAmount string `json:"fillAmount"`
	Fee        string `json:"fee"`
	FeeCoin    string `json:"feeCoin"`
	TradeRole  string `json:"tradeRole"`
	Timestamp  string `json:"timestamp"`
}

// Position represents position information
type Position struct {
	Symbol           string `json:"symbol"`
	Category         string `json:"category"`
	Side             string `json:"side"`
	Size             string `json:"size"`
	AvgPrice         string `json:"avgPrice"`
	MarkPrice        string `json:"markPrice"`
	UnrealizedPNL    string `json:"unrealizedPNL"`
	Leverage         string `json:"leverage"`
	MarginMode       string `json:"marginMode"`
	PositionMargin   string `json:"positionMargin"`
	LiquidationPrice string `json:"liquidationPrice"`
	CreatedTime      string `json:"createdTime"`
	UpdatedTime      string `json:"updatedTime"`
}

// StrategyOrder represents strategy order information
type StrategyOrder struct {
	OrderID      string `json:"orderId"`
	ClientOid    string `json:"clientOid"`
	Symbol       string `json:"symbol"`
	Category     string `json:"category"`
	StrategyType string `json:"strategyType"`
	TriggerPrice string `json:"triggerPrice"`
	TriggerType  string `json:"triggerType"`
	OrderType    string `json:"orderType"`
	Price        string `json:"price,omitempty"`
	Size         string `json:"size"`
	Side         string `json:"side"`
	PositionSide string `json:"positionSide,omitempty"`
	Status       string `json:"status"`
	CreatedTime  string `json:"createdTime"`
	UpdatedTime  string `json:"updatedTime"`
}

// Market data structures

// Ticker represents ticker information
type Ticker struct {
	Symbol            string `json:"symbol"`
	Category          string `json:"category"`
	LastPrice         string `json:"lastPrice"`
	OpenPrice24h      string `json:"openPrice24h"`
	HighPrice24h      string `json:"highPrice24h"`
	LowPrice24h       string `json:"lowPrice24h"`
	Ask1Price         string `json:"ask1Price"`
	Bid1Price         string `json:"bid1Price"`
	Bid1Size          string `json:"bid1Size"`
	Ask1Size          string `json:"ask1Size"`
	Price24hPcnt      string `json:"price24hPcnt"`
	Volume24h         string `json:"volume24h"`
	Turnover24h       string `json:"turnover24h"`
	IndexPrice        string `json:"indexPrice,omitempty"`
	MarkPrice         string `json:"markPrice,omitempty"`
	FundingRate       string `json:"fundingRate,omitempty"`
	OpenInterest      string `json:"openInterest,omitempty"`
	DeliveryStartTime string `json:"deliveryStartTime,omitempty"`
	DeliveryTime      string `json:"deliveryTime,omitempty"`
	DeliveryStatus    string `json:"deliveryStatus,omitempty"`
	Timestamp         string `json:"ts"`
}

// Candlestick represents OHLCV data
type Candlestick struct {
	Timestamp string
	Open      string
	High      string
	Low       string
	Close     string
	Volume    string
	Turnover  string
}

// UnmarshalJSON implements custom JSON unmarshaling for Candlestick
func (c *Candlestick) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}

	if len(arr) < 7 {
		return nil
	}

	c.Timestamp = arr[0]
	c.Open = arr[1]
	c.High = arr[2]
	c.Low = arr[3]
	c.Close = arr[4]
	c.Volume = arr[5]
	c.Turnover = arr[6]

	return nil
}

// OrderBook represents order book data
type OrderBook struct {
	Symbol    string     `json:"symbol"`
	Category  string     `json:"category"`
	Bids      [][]string `json:"bids"`
	Asks      [][]string `json:"asks"`
	Timestamp string     `json:"ts"`
}

// Financial record structures

// FinancialRecord represents financial transaction record
type FinancialRecord struct {
	RecordID  string `json:"recordId"`
	Coin      string `json:"coin"`
	Type      string `json:"type"`
	Amount    string `json:"amount"`
	Balance   string `json:"balance"`
	Fee       string `json:"fee,omitempty"`
	FeeCoin   string `json:"feeCoin,omitempty"`
	Timestamp string `json:"timestamp"`
}

// ConvertRecord represents coin conversion record
type ConvertRecord struct {
	FromCoin     string `json:"fromCoin"`
	FromCoinSize string `json:"fromCoinSize"`
	ToCoin       string `json:"toCoin"`
	ToCoinSize   string `json:"toCoinSize"`
	Price        string `json:"price"`
	Timestamp    string `json:"timestamp"`
}

// PaymentCoin represents payment coin information
type PaymentCoin struct {
	Coin          string `json:"coin"`
	Available     string `json:"available"`
	AvailableUSD  string `json:"availableUSD"`
	MaxSelectable string `json:"maxSelectable"`
}

// RepayableCoin represents repayable coin information
type RepayableCoin struct {
	Coin          string `json:"coin"`
	RepayableSize string `json:"repayableSize"`
	MaxSelectable string `json:"maxSelectable"`
}

// RepayResult represents repayment operation result
type RepayResult struct {
	Result      string `json:"result"`
	RepayAmount string `json:"repayAmount"`
}

// Pagination structures

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	NextFlag bool   `json:"nextFlag"`
	Cursor   string `json:"cursor,omitempty"`
}

// General utility structures

// SwitchStatus represents account switch status
type SwitchStatus struct {
	Status string `json:"status"` // process, success, fail
}

// DeductInfo represents BGB deduction information
type DeductInfo struct {
	Deduct string `json:"deduct"` // on or off
}

// TransferableCoin represents transferable coin information
type TransferableCoin struct {
	Coin string `json:"coin"`
}

// MaxOpenAvailable represents maximum open available information
type MaxOpenAvailable struct {
	Symbol      string `json:"symbol"`
	Category    string `json:"category"`
	MaxBuySize  string `json:"maxBuySize"`
	MaxSellSize string `json:"maxSellSize"`
}

// Institutional loan structures

// LoanOrder represents institutional loan order
type LoanOrder struct {
	OrderID   string `json:"orderId"`
	Coin      string `json:"coin"`
	Amount    string `json:"amount"`
	Rate      string `json:"rate"`
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

// Helper functions for time conversion

// ParseTimestamp converts Unix millisecond timestamp string to time.Time
func ParseTimestamp(ts string) (time.Time, error) {
	if ts == "" {
		return time.Time{}, nil
	}

	var timestamp int64
	if err := json.Unmarshal([]byte(ts), &timestamp); err != nil {
		return time.Time{}, err
	}

	return time.Unix(timestamp/1000, (timestamp%1000)*1000000), nil
}

// FormatTimestamp converts time.Time to Unix millisecond timestamp string
func FormatTimestamp(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return string(rune(t.UnixNano() / 1000000))
}
