package uta

import (
	"context"

	"github.com/khanbekov/go-bitget/common"
)

// AccountAssetsService retrieves account asset information
type AccountAssetsService struct {
	c ClientInterface
}

// Do executes the account assets request
func (s *AccountAssetsService) Do(ctx context.Context) (*AccountAssets, error) {
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointAccountAssets, nil, nil, true)
	if err != nil {
		return nil, err
	}

	var accountAssets AccountAssets
	if err := common.UnmarshalJSON(res.Data, &accountAssets); err != nil {
		return nil, err
	}

	return &accountAssets, nil
}
