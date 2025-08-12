package market

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// TickerService retrieves 24hr ticker statistics for a specific symbol
type TickerService struct {
	c           futures.ClientInterface
	symbol      string
	productType futures.ProductType
}

// Symbol sets the trading pair symbol for the ticker request.
// Required parameter. Examples: "BTCUSDT", "ETHUSDT", "ADAUSDT".
func (s *TickerService) Symbol(symbol string) *TickerService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the request.
// Required parameter. Use "USDT-FUTURES", "USDC-FUTURES", or "COIN-FUTURES".
func (s *TickerService) ProductType(productType string) *TickerService {
	s.productType = futures.ProductType(productType)
	return s
}

// Do executes the ticker request and returns the result.
// Returns a Ticker object containing 24hr statistics for the specified symbol.
//
// The context can be used for request cancellation and timeout control.
// Returns an error if the request fails or if required parameters are missing.
func (s *TickerService) Do(ctx context.Context) (*Ticker, error) {
	queryParams := url.Values{}

	// Set required params
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointTicker, queryParams, nil, false)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response - API returns array of tickers
	var tickers []Ticker
	err = jsoniter.Unmarshal(res.Data, &tickers)

	if err != nil {
		return nil, err
	}

	// Return first ticker if available
	if len(tickers) == 0 {
		return nil, fmt.Errorf("no ticker data returned for symbol %s", s.symbol)
	}

	return &tickers[0], nil
}