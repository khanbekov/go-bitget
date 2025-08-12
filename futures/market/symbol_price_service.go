package market

import (
	"context"
	"encoding/json"
	"net/url"
)

// SymbolPriceService handles retrieving symbol prices (mark/index/market) for futures contracts.
type SymbolPriceService struct {
	c ClientInterface

	// Required parameters
	symbol      string
	productType ProductType
}

// Symbol sets the trading symbol for which to get price data (e.g., "BTCUSDT").
func (s *SymbolPriceService) Symbol(symbol string) *SymbolPriceService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the symbol price request.
func (s *SymbolPriceService) ProductType(productType ProductType) *SymbolPriceService {
	s.productType = productType
	return s
}

// SymbolPrice represents the price information for a symbol.
type SymbolPrice struct {
	Symbol     string `json:"symbol"`     // Trading symbol
	MarkPrice  string `json:"markPrice"`  // Mark price
	IndexPrice string `json:"indexPrice"` // Index price
	LastPrice  string `json:"lastPrice"`  // Last traded price
	Timestamp  string `json:"timestamp"`  // Timestamp (ms)
}

// SymbolPriceResponse represents the response from the symbol price API.
type SymbolPriceResponse struct {
	SymbolPrices []SymbolPrice
}

// Do executes the symbol price request.
func (s *SymbolPriceService) Do(ctx context.Context) (*SymbolPriceResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("symbol", s.symbol)
	params.Set("productType", string(s.productType))

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointSymbolPrice, params, nil, false)
	if err != nil {
		return nil, err
	}

	// Parse response - the API returns an array directly
	var symbolPrices []SymbolPrice
	if err := json.Unmarshal(res.Data, &symbolPrices); err != nil {
		return nil, err
	}

	return &SymbolPriceResponse{SymbolPrices: symbolPrices}, nil
}