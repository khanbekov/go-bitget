package account

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// AccountListService handles retrieving the list of all futures accounts.
type AccountListService struct {
	c futures.ClientInterface

	// Optional parameters
	productType *futures.ProductType
}

// ProductType sets the product type to filter accounts by.
func (s *AccountListService) ProductType(productType futures.ProductType) *AccountListService {
	s.productType = &productType
	return s
}

// AccountListItem represents a single account in the account list.
type AccountListItem struct {
	MarginCoin          string `json:"marginCoin"`          // Margin coin
	Locked              string `json:"locked"`              // Locked amount
	Available           string `json:"available"`           // Available amount
	CrossMaxSize        string `json:"crossMaxSize"`        // Maximum cross margin size
	IsolatedMaxSize     string `json:"isolatedMaxSize"`     // Maximum isolated margin size
	MaxTransferOut      string `json:"maxTransferOut"`      // Maximum transfer out amount
	AccountEquity       string `json:"accountEquity"`       // Account equity
	UsdtEquity          string `json:"usdtEquity"`          // USDT equity
	BtcEquity           string `json:"btcEquity"`           // BTC equity
	CrossRiskRate       string `json:"crossRiskRate"`       // Cross margin risk rate
	CrossMarginLeverage string `json:"crossMarginLeverage"` // Cross margin leverage
	FixedLongLeverage   string `json:"fixedLongLeverage"`   // Fixed long leverage
	FixedShortLeverage  string `json:"fixedShortLeverage"`  // Fixed short leverage
	MarginMode          string `json:"marginMode"`          // Margin mode
	PositionMode        string `json:"positionMode"`        // Position mode
	UnrealizedPL        string `json:"unrealizedPL"`        // Unrealized profit/loss
	CouponAmount        string `json:"couponAmount"`        // Coupon amount
}

// AccountListResponse represents the response from the account list API.
type AccountListResponse struct {
	Accounts []AccountListItem
}

// Do executes the account list request.
func (s *AccountListService) Do(ctx context.Context) (*AccountListResponse, error) {
	// Build query parameters
	params := url.Values{}
	if s.productType != nil {
		params.Set("productType", string(*s.productType))
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointAccountList, params, nil, true)
	if err != nil {
		return nil, err
	}

	// Parse response - the API returns an array directly
	var accounts []AccountListItem
	if err := json.Unmarshal(res.Data, &accounts); err != nil {
		return nil, err
	}

	return &AccountListResponse{Accounts: accounts}, nil
}
