package account

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/futures"
)

// SetPositionModeService handles setting the position mode for futures trading.
// Position mode can be either one-way (single position) or hedge (dual position).
type SetPositionModeService struct {
	c futures.ClientInterface

	// Required parameters
	productType  futures.ProductType
	positionMode futures.PositionModeType
}

// ProductType sets the product type for the position mode setting.
func (s *SetPositionModeService) ProductType(productType futures.ProductType) *SetPositionModeService {
	s.productType = productType
	return s
}

// PositionMode sets the position mode (one-way or hedge).
func (s *SetPositionModeService) PositionMode(positionMode futures.PositionModeType) *SetPositionModeService {
	s.positionMode = positionMode
	return s
}

// SetPositionModeResponse represents the response from setting position mode.
type SetPositionModeResponse struct {
	ProductType  string `json:"productType"`  // Product type
	PositionMode string `json:"positionMode"` // New position mode
}

// Do executes the set position mode request.
func (s *SetPositionModeService) Do(ctx context.Context) (*SetPositionModeResponse, error) {
	// Build request body
	params := map[string]interface{}{
		"productType": string(s.productType),
		"posMode":     string(s.positionMode),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "POST", futures.EndpointSetPositionMode, nil, body, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result SetPositionModeResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
