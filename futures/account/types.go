package account

import (
	"github.com/khanbekov/go-bitget/common/client"
)

// Re-export common types to avoid importing futures package
type (
	ClientInterface = client.ClientInterface
	ApiResponse     = client.ApiResponse
)

// Futures-specific enums (duplicated to avoid import cycle)
type ProductType string

const (
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	ProductTypeCoinFutures ProductType = "COIN-FUTURES"
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
)

type MarginMode string

const (
	MarginModeCrossed  MarginMode = "crossed"
	MarginModeIsolated MarginMode = "isolated"
)

type PositionMode string

const (
	PositionModeOneWay PositionMode = "one_way_mode"
	PositionModeHedge  PositionMode = "hedge_mode"
)

type HoldSide string

const (
	HoldSideLong  HoldSide = "long"
	HoldSideShort HoldSide = "short"
)

// API Endpoints for account operations
const (
	EndpointAccountInfo      = "/api/v2/mix/account/account"
	EndpointAccountList      = "/api/v2/mix/account/accounts"
	EndpointAccountBills     = "/api/v2/mix/account/bill"
	EndpointSetLeverage      = "/api/v2/mix/account/set-leverage"
	EndpointAdjustMargin     = "/api/v2/mix/account/set-margin"
	EndpointSetMarginMode    = "/api/v2/mix/account/set-margin-mode"
	EndpointSetPositionMode  = "/api/v2/mix/account/set-position-mode"
)

// Service Constructor Functions

// NewAccountInfoService creates a new account information service.
func NewAccountInfoService(client ClientInterface) *AccountInfoService {
	return &AccountInfoService{c: client}
}

// NewAccountListService creates a new account list service.
func NewAccountListService(client ClientInterface) *AccountListService {
	return &AccountListService{c: client}
}

// NewSetLeverageService creates a new leverage setting service.
func NewSetLeverageService(client ClientInterface) *SetLeverageService {
	return &SetLeverageService{c: client}
}

// NewAdjustMarginService creates a new margin adjustment service.
func NewAdjustMarginService(client ClientInterface) *AdjustMarginService {
	return &AdjustMarginService{c: client}
}

// NewSetMarginModeService creates a new margin mode setting service.
func NewSetMarginModeService(client ClientInterface) *SetMarginModeService {
	return &SetMarginModeService{c: client}
}

// NewSetPositionModeService creates a new position mode setting service.
func NewSetPositionModeService(client ClientInterface) *SetPositionModeService {
	return &SetPositionModeService{c: client}
}

// NewGetAccountBillService creates a new account bill service.
func NewGetAccountBillService(client ClientInterface) *GetAccountBillService {
	return &GetAccountBillService{c: client}
}