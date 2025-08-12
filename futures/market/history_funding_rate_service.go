package market

import (
	"context"
	"encoding/json"
	"net/url"
)

// HistoryFundingRateService handles retrieving historical funding rates for futures contracts.
type HistoryFundingRateService struct {
	c ClientInterface

	// Required parameters
	symbol      string
	productType ProductType

	// Optional parameters
	pageSize  *string // Number of results per page (default: 20, max: 100)
	pageNo    *string // Page number (default: 1)
	nextPage  *string // Next page token for pagination
	endTime   *string // End time timestamp (ms)
	startTime *string // Start time timestamp (ms)
}

// Symbol sets the trading symbol for which to get historical funding rates (e.g., "BTCUSDT").
func (s *HistoryFundingRateService) Symbol(symbol string) *HistoryFundingRateService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the funding rate request.
func (s *HistoryFundingRateService) ProductType(productType ProductType) *HistoryFundingRateService {
	s.productType = productType
	return s
}

// PageSize sets the number of results per page (default: 20, max: 100).
func (s *HistoryFundingRateService) PageSize(pageSize string) *HistoryFundingRateService {
	s.pageSize = &pageSize
	return s
}

// PageNo sets the page number (default: 1).
func (s *HistoryFundingRateService) PageNo(pageNo string) *HistoryFundingRateService {
	s.pageNo = &pageNo
	return s
}

// NextPage sets the next page token for pagination.
func (s *HistoryFundingRateService) NextPage(nextPage string) *HistoryFundingRateService {
	s.nextPage = &nextPage
	return s
}

// EndTime sets the end time timestamp (ms).
func (s *HistoryFundingRateService) EndTime(endTime string) *HistoryFundingRateService {
	s.endTime = &endTime
	return s
}

// StartTime sets the start time timestamp (ms).
func (s *HistoryFundingRateService) StartTime(startTime string) *HistoryFundingRateService {
	s.startTime = &startTime
	return s
}

// HistoryFundingRate represents a historical funding rate record.
type HistoryFundingRate struct {
	Symbol            string `json:"symbol"`            // Trading symbol
	FundingRate       string `json:"fundingRate"`       // Funding rate
	SettlementFunding string `json:"settlementFunding"` // Settlement funding
	FundingTime       string `json:"fundingTime"`       // Funding time (timestamp)
	SettlementTime    string `json:"settlementTime"`    // Settlement time (timestamp)
}

// HistoryFundingRateResponse represents the response from the history funding rate API.
type HistoryFundingRateResponse struct {
	FundingRates []HistoryFundingRate `json:"data"`
	NextPage     string               `json:"nextPage"`
}

// Do executes the history funding rate request.
func (s *HistoryFundingRateService) Do(ctx context.Context) (*HistoryFundingRateResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("symbol", s.symbol)
	params.Set("productType", string(s.productType))

	// Add optional parameters
	if s.pageSize != nil {
		params.Set("pageSize", *s.pageSize)
	}
	if s.pageNo != nil {
		params.Set("pageNo", *s.pageNo)
	}
	if s.nextPage != nil {
		params.Set("nextPage", *s.nextPage)
	}
	if s.endTime != nil {
		params.Set("endTime", *s.endTime)
	}
	if s.startTime != nil {
		params.Set("startTime", *s.startTime)
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointHistoryFundingRate, params, nil, false)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result HistoryFundingRateResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}