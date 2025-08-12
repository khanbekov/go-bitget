package uta

import (
	"context"

	"github.com/khanbekov/go-bitget/common"
)

// AccountInfoService retrieves account settings and configuration
type AccountInfoService struct {
	c ClientInterface
}

// Do executes the account info request
func (s *AccountInfoService) Do(ctx context.Context) (*AccountInfo, error) {
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointAccountSettings, nil, nil, true)
	if err != nil {
		return nil, err
	}

	var accountInfo AccountInfo
	if err := common.UnmarshalJSON(res.Data, &accountInfo); err != nil {
		return nil, err
	}

	return &accountInfo, nil
}
