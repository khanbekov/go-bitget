package account

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/futures"
)

// SetMarginModeService handles setting the margin mode for futures trading positions.
// Margin mode can be either isolated (each position has its own margin) or cross (all positions share margin).
type SetMarginModeService struct {
	c futures.ClientInterface

	// Required parameters
	symbol      string
	productType futures.ProductType
	marginMode  futures.MarginModeType
	marginCoin  string
}

// Symbol sets the trading symbol for which to set the margin mode (e.g., "BTCUSDT").
func (s *SetMarginModeService) Symbol(symbol string) *SetMarginModeService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the margin mode setting.
func (s *SetMarginModeService) ProductType(productType futures.ProductType) *SetMarginModeService {
	s.productType = productType
	return s
}

// MarginCoin sets the margin coin for the adjustment.
func (s *SetMarginModeService) MarginCoin(marginCoin string) *SetMarginModeService {
	s.marginCoin = marginCoin
	return s
}

// MarginMode sets the margin mode (isolated or cross).
func (s *SetMarginModeService) MarginMode(marginMode futures.MarginModeType) *SetMarginModeService {
	s.marginMode = marginMode
	return s
}

// SetMarginModeResponse represents the response from setting margin mode.
type SetMarginModeResponse struct {
	Symbol      string `json:"symbol"`      // Trading symbol
	ProductType string `json:"productType"` // Product type
	MarginMode  string `json:"marginMode"`  // New margin mode
}

// Do executes the set margin mode request.
func (s *SetMarginModeService) Do(ctx context.Context) (*SetMarginModeResponse, error) {
	// Build request body
	params := map[string]interface{}{
		"symbol":      s.symbol,
		"productType": string(s.productType),
		"marginMode":  string(s.marginMode),
		"marginCoin":  s.marginCoin,
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "POST", futures.EndpointSetMarginMode, nil, body, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result SetMarginModeResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
