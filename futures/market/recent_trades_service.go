package market

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"
	"strconv"
)

// RecentTradesService retrieves recent trade executions
type RecentTradesService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	limit       string
}

// Symbol sets the trading pair symbol for the recent trades request.
// Required parameter. Examples: "BTCUSDT", "ETHUSDT", "ADAUSDT".
func (s *RecentTradesService) Symbol(symbol string) *RecentTradesService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the request.
// Required parameter. Use ProductTypeUSDTFutures, ProductTypeUSDCFutures, or ProductTypeCoinFutures.
func (s *RecentTradesService) ProductType(productType ProductType) *RecentTradesService {
	s.productType = productType
	return s
}

// Limit sets the number of recent trades to return.
// Optional parameter. Default and maximum value is 100.
func (s *RecentTradesService) Limit(limit string) *RecentTradesService {
	s.limit = limit
	return s
}

// Do executes the recent trades request and returns the results.
// Returns a slice of RecentTrade objects containing recent trade executions.
//
// The context can be used for request cancellation and timeout control.
// Returns an error if the request fails or if required parameters are missing.
func (s *RecentTradesService) Do(ctx context.Context) ([]*RecentTrade, error) {
	queryParams := url.Values{}

	// Set required params
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))

	// Set optional params
	if s.limit != "" {
		queryParams.Set("limit", s.limit)
	}

	// Make request to API
	var res *ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointRecentTrades, queryParams, nil, false)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var trades []*RecentTrade
	err = jsoniter.Unmarshal(res.Data, &trades)

	if err != nil {
		return nil, err
	}

	return trades, nil
}

// RecentTrade represents a recent trade execution
type RecentTrade struct {
	// Trade ID
	TradeId string `json:"tradeId"`

	// Trade price
	Price float64 `json:"price"`

	// Trade size/quantity
	Size float64 `json:"size"`

	// Trade side (buy/sell from taker perspective)
	Side string `json:"side"`

	// Trade timestamp
	Ts int64 `json:"ts"`
}

// UnmarshalJSON implements custom JSON unmarshaling for RecentTrade.
// The Bitget API returns trade data as arrays of strings rather than objects,
// so this method converts the array format to the struct fields.
//
// Expected array format: [tradeId, price, size, side, timestamp]
func (r *RecentTrade) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 5 {
		return fmt.Errorf("expected 5 elements for recent trade, received %d", len(arr))
	}

	// Parse tradeId
	r.TradeId = arr[0]

	// Parse price
	price, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		return fmt.Errorf("failed parsing price: %v", err)
	}
	r.Price = price

	// Parse size
	size, err := strconv.ParseFloat(arr[2], 64)
	if err != nil {
		return fmt.Errorf("failed parsing size: %v", err)
	}
	r.Size = size

	// Parse side
	r.Side = arr[3]

	// Parse timestamp
	ts, err := strconv.ParseInt(arr[4], 10, 64)
	if err != nil {
		return fmt.Errorf("failed parsing timestamp: %v", err)
	}
	r.Ts = ts

	return nil
}