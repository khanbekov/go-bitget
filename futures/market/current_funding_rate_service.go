package market

import (
	"context"
	"encoding/json"
	"net/url"
)

// CurrentFundingRateService handles retrieving current funding rates for futures contracts.
type CurrentFundingRateService struct {
	c ClientInterface

	// Required parameters
	symbol      string
	productType ProductType
}

// Symbol sets the trading symbol for which to get the current funding rate (e.g., "BTCUSDT").
func (s *CurrentFundingRateService) Symbol(symbol string) *CurrentFundingRateService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the funding rate request.
func (s *CurrentFundingRateService) ProductType(productType ProductType) *CurrentFundingRateService {
	s.productType = productType
	return s
}

// CurrentFundingRate represents the current funding rate information.
type CurrentFundingRate struct {
	Symbol            string `json:"symbol"`            // Trading symbol
	FundingRate       string `json:"fundingRate"`       // Current funding rate
	SettlementFunding string `json:"settlementFunding"` // Settlement funding
	FundingTime       string `json:"fundingTime"`       // Next funding time (timestamp)
	SettlementTime    string `json:"settlementTime"`    // Settlement time (timestamp)
}

// CurrentFundingRateResponse represents the response from the current funding rate API.
type CurrentFundingRateResponse struct {
	FundingRates []CurrentFundingRate
}

// Do executes the current funding rate request.
func (s *CurrentFundingRateService) Do(ctx context.Context) (*CurrentFundingRateResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("symbol", s.symbol)
	params.Set("productType", string(s.productType))

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointCurrentFundingRate, params, nil, false)
	if err != nil {
		return nil, err
	}

	// Parse response - the API returns an array directly
	var fundingRates []CurrentFundingRate
	if err := json.Unmarshal(res.Data, &fundingRates); err != nil {
		return nil, err
	}

	return &CurrentFundingRateResponse{FundingRates: fundingRates}, nil
}