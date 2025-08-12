package uta

import (
	"context"
	"net/url"

	"github.com/khanbekov/go-bitget/common"
)

// AccountFeeRateService retrieves trading fee rates for a symbol
type AccountFeeRateService struct {
	c        ClientInterface
	symbol   *string
	category *string
}

// Symbol sets the trading symbol (required)
func (s *AccountFeeRateService) Symbol(symbol string) *AccountFeeRateService {
	s.symbol = &symbol
	return s
}

// Category sets the product category (required)
func (s *AccountFeeRateService) Category(category string) *AccountFeeRateService {
	s.category = &category
	return s
}

// Do executes the fee rate request
func (s *AccountFeeRateService) Do(ctx context.Context) (*FeeRate, error) {
	if s.symbol == nil {
		return nil, common.NewMissingParameterError("symbol")
	}
	if s.category == nil {
		return nil, common.NewMissingParameterError("category")
	}

	params := url.Values{}
	params.Set("symbol", *s.symbol)
	params.Set("category", *s.category)

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointAccountFeeRate, params, nil, true)
	if err != nil {
		return nil, err
	}

	var feeRate FeeRate
	if err := common.UnmarshalJSON(res.Data, &feeRate); err != nil {
		return nil, err
	}

	return &feeRate, nil
}
