package market

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"
	"strconv"

	"github.com/khanbekov/go-bitget/futures"
)

// OrderBookService retrieves order book depth data
type OrderBookService struct {
	c           futures.ClientInterface
	symbol      string
	productType futures.ProductType
	precision   string
	limit       string
}

// Symbol sets the trading pair symbol for the order book request.
// Required parameter. Examples: "BTCUSDT", "ETHUSDT", "ADAUSDT".
func (s *OrderBookService) Symbol(symbol string) *OrderBookService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the request.
// Required parameter. Use ProductTypeUSDTFutures, ProductTypeUSDCFutures, or ProductTypeCOINFutures.
func (s *OrderBookService) ProductType(productType futures.ProductType) *OrderBookService {
	s.productType = productType
	return s
}

// Precision sets the price precision for the order book.
// Optional parameter. Examples: "0.1", "0.01", "0.001".
func (s *OrderBookService) Precision(precision string) *OrderBookService {
	s.precision = precision
	return s
}

// Limit sets the number of order book levels to return.
// Optional parameter. Valid values: "5", "15", "50", "100". Default is "100".
func (s *OrderBookService) Limit(limit string) *OrderBookService {
	s.limit = limit
	return s
}

// Do executes the order book request and returns the result.
// Returns an OrderBook object containing bid and ask levels.
//
// The context can be used for request cancellation and timeout control.
// Returns an error if the request fails or if required parameters are missing.
func (s *OrderBookService) Do(ctx context.Context) (*OrderBook, error) {
	queryParams := url.Values{}

	// Set required params
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))

	// Set optional params
	if s.precision != "" {
		queryParams.Set("precision", s.precision)
	}
	if s.limit != "" {
		queryParams.Set("limit", s.limit)
	}

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointMergeDepth, queryParams, nil, false)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var orderBook *OrderBook
	err = jsoniter.Unmarshal(res.Data, &orderBook)

	if err != nil {
		return nil, err
	}

	return orderBook, nil
}

// OrderBook represents order book depth data
type OrderBook struct {
	// Ask levels (sell orders) - sorted from lowest to highest price
	Asks []OrderBookLevel `json:"asks"`

	// Bid levels (buy orders) - sorted from highest to lowest price
	Bids []OrderBookLevel `json:"bids"`

	// Timestamp of the order book snapshot
	Ts string `json:"ts"`
}

// OrderBookLevel represents a single price level in the order book
type OrderBookLevel struct {
	// Price level
	Price float64 `json:"price"`

	// Total quantity at this price level
	Size float64 `json:"size"`
}

// UnmarshalJSON implements custom JSON unmarshaling for OrderBookLevel.
// The Bitget API returns order book levels as arrays of strings [price, size],
// so this method converts the array format to the struct fields.
func (o *OrderBookLevel) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 2 {
		return fmt.Errorf("expected 2 elements for order book level, received %d", len(arr))
	}

	// Parse price
	price, err := strconv.ParseFloat(arr[0], 64)
	if err != nil {
		return fmt.Errorf("failed parsing price: %v", err)
	}

	// Parse size
	size, err := strconv.ParseFloat(arr[1], 64)
	if err != nil {
		return fmt.Errorf("failed parsing size: %v", err)
	}

	o.Price = price
	o.Size = size

	return nil
}