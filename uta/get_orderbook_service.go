package uta

import (
	"context"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"net/url"
)

// GetOrderBookService retrieves order book data from UTA API
type GetOrderBookService struct {
	c        ClientInterface
	symbol   string
	category string
	limit    *int
}

// NewGetOrderBookService creates a new GetOrderBookService instance
func NewGetOrderBookService(c ClientInterface) *GetOrderBookService {
	return &GetOrderBookService{
		c: c,
	}
}

// Symbol sets the trading pair symbol for the order book request
func (s *GetOrderBookService) Symbol(symbol string) *GetOrderBookService {
	s.symbol = symbol
	return s
}

// Category sets the product category (SPOT, USDT-FUTURES, etc.)
func (s *GetOrderBookService) Category(category string) *GetOrderBookService {
	s.category = category
	return s
}

// Limit sets the order book depth limit (optional, default based on category)
func (s *GetOrderBookService) Limit(limit int) *GetOrderBookService {
	s.limit = &limit
	return s
}

// Do executes the order book request and returns the result
func (s *GetOrderBookService) Do(ctx context.Context) (*OrderBook, error) {
	// Validate required parameters
	if s.symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}
	if s.category == "" {
		return nil, fmt.Errorf("category is required")
	}

	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("category", s.category)

	if s.limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *s.limit))
	}

	// Make API request
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointMarketOrderbook, queryParams, nil, false)
	if err != nil {
		return nil, err
	}

	// Parse response
	var orderBook OrderBook
	if err := jsoniter.Unmarshal(res.Data, &orderBook); err != nil {
		return nil, fmt.Errorf("failed to unmarshal order book data: %w", err)
	}

	return &orderBook, nil
}