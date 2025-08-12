package uta

import (
	"context"
	"net/url"

	"github.com/khanbekov/go-bitget/common"
)

// GetTickersService retrieves ticker information
type GetTickersService struct {
	c        ClientInterface
	category *string
	symbol   *string
}

// Category sets the product category (required)
func (s *GetTickersService) Category(category string) *GetTickersService {
	s.category = &category
	return s
}

// Symbol sets the trading symbol (optional, if not set returns all symbols)
func (s *GetTickersService) Symbol(symbol string) *GetTickersService {
	s.symbol = &symbol
	return s
}

// Do executes the get tickers request
func (s *GetTickersService) Do(ctx context.Context) ([]Ticker, error) {
	if s.category == nil {
		return nil, common.NewMissingParameterError("category")
	}

	params := url.Values{}
	params.Set("category", *s.category)

	if s.symbol != nil {
		params.Set("symbol", *s.symbol)
	}

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointMarketTickers, params, nil, false)
	if err != nil {
		return nil, err
	}

	var tickers []Ticker
	if err := common.UnmarshalJSON(res.Data, &tickers); err != nil {
		return nil, err
	}

	return tickers, nil
}
