package uta

import (
	"context"
	"net/url"

	"github.com/khanbekov/go-bitget/common"
)

// GetCandlesticksService retrieves candlestick/OHLCV data
type GetCandlesticksService struct {
	c         ClientInterface
	category  *string
	symbol    *string
	interval  *string
	startTime *string
	endTime   *string
	dataType  *string
	limit     *string
}

// Category sets the product category (required)
func (s *GetCandlesticksService) Category(category string) *GetCandlesticksService {
	s.category = &category
	return s
}

// Symbol sets the trading symbol (required)
func (s *GetCandlesticksService) Symbol(symbol string) *GetCandlesticksService {
	s.symbol = &symbol
	return s
}

// Interval sets the time interval (required)
func (s *GetCandlesticksService) Interval(interval string) *GetCandlesticksService {
	s.interval = &interval
	return s
}

// StartTime sets the start timestamp (optional)
func (s *GetCandlesticksService) StartTime(startTime string) *GetCandlesticksService {
	s.startTime = &startTime
	return s
}

// EndTime sets the end timestamp (optional)
func (s *GetCandlesticksService) EndTime(endTime string) *GetCandlesticksService {
	s.endTime = &endTime
	return s
}

// Type sets the candlestick type (optional): "MARKET", "MARK", "INDEX"
func (s *GetCandlesticksService) Type(dataType string) *GetCandlesticksService {
	s.dataType = &dataType
	return s
}

// Limit sets the number of results to return (optional, max 100)
func (s *GetCandlesticksService) Limit(limit string) *GetCandlesticksService {
	s.limit = &limit
	return s
}

// Do executes the get candlesticks request
func (s *GetCandlesticksService) Do(ctx context.Context) ([]Candlestick, error) {
	if s.category == nil {
		return nil, common.NewMissingParameterError("category")
	}
	if s.symbol == nil {
		return nil, common.NewMissingParameterError("symbol")
	}
	if s.interval == nil {
		return nil, common.NewMissingParameterError("interval")
	}

	params := url.Values{}
	params.Set("category", *s.category)
	params.Set("symbol", *s.symbol)
	params.Set("interval", *s.interval)

	if s.startTime != nil {
		params.Set("startTime", *s.startTime)
	}
	if s.endTime != nil {
		params.Set("endTime", *s.endTime)
	}
	if s.dataType != nil {
		params.Set("type", *s.dataType)
	}
	if s.limit != nil {
		params.Set("limit", *s.limit)
	}

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointMarketCandles, params, nil, false)
	if err != nil {
		return nil, err
	}

	var candlesticks []Candlestick
	if err := common.UnmarshalJSON(res.Data, &candlesticks); err != nil {
		return nil, err
	}

	return candlesticks, nil
}
