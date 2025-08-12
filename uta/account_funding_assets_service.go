package uta

import (
	"context"
	"net/url"

	"github.com/khanbekov/go-bitget/common"
)

// AccountFundingAssetsService retrieves funding account asset information
type AccountFundingAssetsService struct {
	c    ClientInterface
	coin *string
}

// Coin sets the coin parameter to filter by specific coin
func (s *AccountFundingAssetsService) Coin(coin string) *AccountFundingAssetsService {
	s.coin = &coin
	return s
}

// Do executes the funding assets request
func (s *AccountFundingAssetsService) Do(ctx context.Context) ([]FundingAssets, error) {
	params := url.Values{}
	if s.coin != nil {
		params.Set("coin", *s.coin)
	}

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointAccountFundingAssets, params, nil, true)
	if err != nil {
		return nil, err
	}

	var fundingAssets []FundingAssets
	if err := common.UnmarshalJSON(res.Data, &fundingAssets); err != nil {
		return nil, err
	}

	return fundingAssets, nil
}
