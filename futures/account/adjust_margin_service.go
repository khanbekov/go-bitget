package account

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/futures"
)

// AdjustMarginService handles adjusting margin for a futures position.
// You can add or reduce margin for isolated positions.
type AdjustMarginService struct {
	c futures.ClientInterface

	// Required parameters
	symbol      string
	productType futures.ProductType
	marginCoin  string
	amount      string
	type_       string // "ADD" or "REDUCE"

	// Optional parameters
	holdSide *string // "LONG" or "SHORT" (for hedge mode)
}

// Symbol sets the trading symbol for margin adjustment (e.g., "BTCUSDT").
func (s *AdjustMarginService) Symbol(symbol string) *AdjustMarginService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the margin adjustment.
func (s *AdjustMarginService) ProductType(productType futures.ProductType) *AdjustMarginService {
	s.productType = productType
	return s
}

// MarginCoin sets the margin coin for the adjustment.
func (s *AdjustMarginService) MarginCoin(marginCoin string) *AdjustMarginService {
	s.marginCoin = marginCoin
	return s
}

// Amount sets the amount to adjust (positive value for both add and reduce).
func (s *AdjustMarginService) Amount(amount string) *AdjustMarginService {
	s.amount = amount
	return s
}

// Type sets the adjustment type ("ADD" or "REDUCE").
func (s *AdjustMarginService) Type(type_ string) *AdjustMarginService {
	s.type_ = type_
	return s
}

// AddMargin is a helper method to set type to "ADD".
func (s *AdjustMarginService) AddMargin() *AdjustMarginService {
	s.type_ = "ADD"
	return s
}

// ReduceMargin is a helper method to set type to "REDUCE".
func (s *AdjustMarginService) ReduceMargin() *AdjustMarginService {
	s.type_ = "REDUCE"
	return s
}

// HoldSide sets the hold side for hedge mode positions ("LONG" or "SHORT").
func (s *AdjustMarginService) HoldSide(holdSide string) *AdjustMarginService {
	s.holdSide = &holdSide
	return s
}

// AdjustMarginResponse represents the response from adjusting margin.
type AdjustMarginResponse struct {
	Symbol      string `json:"symbol"`      // Trading symbol
	ProductType string `json:"productType"` // Product type
	MarginCoin  string `json:"marginCoin"`  // Margin coin
	Amount      string `json:"amount"`      // Adjusted amount
	Type        string `json:"type"`        // Adjustment type
	HoldSide    string `json:"holdSide"`    // Hold side (for hedge mode)
	Success     bool   `json:"success"`     // Whether the adjustment was successful
}

// Do executes the adjust margin request.
func (s *AdjustMarginService) Do(ctx context.Context) (*AdjustMarginResponse, error) {
	// Build request body
	params := map[string]interface{}{
		"symbol":      s.symbol,
		"productType": string(s.productType),
		"marginCoin":  s.marginCoin,
		"amount":      s.amount,
		"type":        s.type_,
	}

	// Add optional parameters
	if s.holdSide != nil {
		params["holdSide"] = *s.holdSide
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "POST", futures.EndpointSetMargin, nil, body, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result AdjustMarginResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
