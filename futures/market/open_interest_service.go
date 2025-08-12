package market

import (
	"context"
	"encoding/json"
	"net/url"
)

// OpenInterestService handles retrieving open interest data for futures contracts.
type OpenInterestService struct {
	c ClientInterface

	// Required parameters
	symbol      string
	productType ProductType
}

// Symbol sets the trading symbol for which to get open interest data (e.g., "BTCUSDT").
func (s *OpenInterestService) Symbol(symbol string) *OpenInterestService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the open interest request.
func (s *OpenInterestService) ProductType(productType ProductType) *OpenInterestService {
	s.productType = productType
	return s
}

// OpenInterest represents the open interest data for a symbol.
type OpenInterest struct {
	Symbol           string `json:"symbol"`           // Trading symbol
	Amount           string `json:"amount"`           // Open interest amount
	Timestamp        string `json:"timestamp"`        // Timestamp (ms)
	Size             string `json:"size"`             // Open interest size
	OpenInterestUSDT string `json:"openInterestUSDT"` // Open interest in USDT
}

// OpenInterestResponse represents the response from the open interest API.
type OpenInterestResponse struct {
	OpenInterests []OpenInterest
}

// Do executes the open interest request.
func (s *OpenInterestService) Do(ctx context.Context) (*OpenInterestResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("symbol", s.symbol)
	params.Set("productType", string(s.productType))

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointOpenInterest, params, nil, false)
	if err != nil {
		return nil, err
	}

	// Parse response - the API returns an array directly
	var openInterests []OpenInterest
	if err := json.Unmarshal(res.Data, &openInterests); err != nil {
		return nil, err
	}

	return &OpenInterestResponse{OpenInterests: openInterests}, nil
}