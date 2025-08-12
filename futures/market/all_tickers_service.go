package market

import (
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// AllTickersService retrieves 24hr ticker statistics for all symbols
type AllTickersService struct {
	c           futures.ClientInterface
	productType futures.ProductType
}

// ProductType sets the product type for the request.
// Required parameter. Use ProductTypeUSDTFutures, ProductTypeUSDCFutures, or ProductTypeCOINFutures.
func (s *AllTickersService) ProductType(productType futures.ProductType) *AllTickersService {
	s.productType = productType
	return s
}

// Do executes the all tickers request and returns the results.
// Returns a slice of Ticker objects containing 24hr statistics for all symbols.
//
// The context can be used for request cancellation and timeout control.
// Returns an error if the request fails or if required parameters are missing.
func (s *AllTickersService) Do(ctx context.Context) ([]*Ticker, error) {
	queryParams := url.Values{}

	// Set required params
	queryParams.Set("productType", string(s.productType))

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointAllTickers, queryParams, nil, false)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var tickers []*Ticker
	err = jsoniter.Unmarshal(res.Data, &tickers)

	if err != nil {
		return nil, err
	}

	return tickers, nil
}

// Ticker represents 24hr ticker statistics for a symbol
type Ticker struct {
	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Last traded price
	LastPr string `json:"lastPr"`

	// Best bid price
	BidPr string `json:"bidPr"`

	// Best bid size
	BidSz string `json:"bidSz"`

	// Best ask price
	AskPr string `json:"askPr"`

	// Best ask size
	AskSz string `json:"askSz"`

	// 24hr price change
	Change24h string `json:"change24h"`

	// Open price (24hr ago)
	Open24h string `json:"open24h"`

	// Highest price in last 24hr
	High24h string `json:"high24h"`

	// Lowest price in last 24hr
	Low24h string `json:"low24h"`

	// 24hr trading volume in base asset
	BaseVolume string `json:"baseVolume"`

	// 24hr trading volume in quote asset
	QuoteVolume string `json:"quoteVolume"`

	// 24hr trading volume in USD
	UsdtVolume string `json:"usdtVolume"`

	// Timestamp of ticker data
	Ts string `json:"ts"`

	// Open interest
	OpenI string `json:"openI"`
}