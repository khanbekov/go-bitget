package account

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/common"
	"golang.org/x/net/context"
	"net/url"
)

// AccountInfoService account info
type AccountInfoService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	marginCoin  string
}

func (s *AccountInfoService) Symbol(symbol string) *AccountInfoService {
	s.symbol = symbol
	return s
}

func (s *AccountInfoService) ProductType(productType ProductType) *AccountInfoService {
	s.productType = productType
	return s
}

func (s *AccountInfoService) MarginCoin(marginCoin string) *AccountInfoService {
	s.marginCoin = marginCoin
	return s
}

func (s *AccountInfoService) Do(ctx context.Context) (acc *Account, err error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))
	queryParams.Set("marginCoin", s.marginCoin)

	// Make request to API
	var res *ApiResponse

	res, _, err = s.c.CallAPI(ctx, "GET", EndpointAccountInfo, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	err = jsoniter.Unmarshal(res.Data, &acc)

	if err != nil {
		return nil, err
	}

	return acc, nil
}

type Account struct {
	// Margin coin used for margin trading
	MarginCoin string `json:"marginCoin"`

	// Locked quantity (margin coin). Lockup will be triggered when there is a position to be closed.
	Locked float64 `json:"locked"`

	// Available quantity in the account
	Available float64 `json:"available"`

	// Maximum available balance to open positions under the cross margin mode (margin coin)
	CrossedMaxAvailable float64 `json:"crossedMaxAvailable"`

	// Maximum available balance to open positions under the isolated margin mode (margin coin)
	IsolatedMaxAvailable float64 `json:"isolatedMaxAvailable"`

	// Maximum transferable amount
	MaxTransferOut float64 `json:"maxTransferOut"`

	// Account equity (margin coin), including unrealized PnL (based on mark price)
	AccountEquity float64 `json:"accountEquity"`

	// Account equity in USDT
	UsdtEquity float64 `json:"usdtEquity"`

	// Account equity in BTC
	BtcEquity float64 `json:"btcEquity"`

	// Risk ratio in cross margin mode
	CrossedRiskRate float64 `json:"crossedRiskRate"`

	// Leverage in cross margin mode
	CrossedMarginLeverage int64 `json:"crossedMarginLeverage"`

	// Leverage of long positions in isolated margin mode
	IsolatedLongLever int64 `json:"isolatedLongLever"`

	// Leverage of short positions in isolated margin mode
	IsolatedShortLever int64 `json:"isolatedShortLever"`

	// Margin mode. isolated – isolated margin mode; crossed – cross margin mode
	MarginMode string `json:"marginMode"`

	// Position mode: one_way_mode (one-way mode) / hedge_mode (hedge mode)
	PosMode string `json:"posMode"`

	// Unrealized PnL
	UnrealizedPL float64 `json:"unrealizedPL"`

	// Trading bonus
	Coupon int64 `json:"coupon"`

	// Unrealized PnL for cross margin mode
	CrossedUnrealizedPL float64 `json:"crossedUnrealizedPL"`

	// Unrealized PnL for isolated margin mode
	IsolatedUnrealizedPL float64 `json:"isolatedUnrealizedPL"`

	// Assets mode: "union" (Multi-assets mode) / "single" (Single-assets mode)
	AssetMode string `json:"assetMode"`

	// Trading grant (omitted if empty)
	Grant string `json:"grant,omitempty"`
}

// UnmarshalJSON realization interface json.Unmarshaler for Account
func (d *Account) UnmarshalJSON(data []byte) error {
	type Alias Account // Чтобы избежать рекурсии
	tmp := map[string]interface{}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for key, value := range tmp {
		switch key {
		case "marginCoin":
			d.MarginCoin = common.SafeStringCast(value)
		case "locked":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.Locked = v
		case "available":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.Available = v
		case "crossedMaxAvailable":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.CrossedMaxAvailable = v
		case "isolatedMaxAvailable":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.IsolatedMaxAvailable = v
		case "maxTransferOut":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.MaxTransferOut = v
		case "accountEquity":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.AccountEquity = v
		case "usdtEquity":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.UsdtEquity = v
		case "btcEquity":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.BtcEquity = v
		case "crossedRiskRate":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.CrossedRiskRate = v
		case "crossedMarginLeverage":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			d.CrossedMarginLeverage = v
		case "isolatedLongLever":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			d.IsolatedLongLever = v
		case "isolatedShortLever":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			d.IsolatedShortLever = v
		case "marginMode":
			d.MarginMode = common.SafeStringCast(value)
		case "posMode":
			d.PosMode = common.SafeStringCast(value)
		case "unrealizedPL":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.UnrealizedPL = v
		case "coupon":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			d.Coupon = v
		case "crossedUnrealizedPL":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.CrossedUnrealizedPL = v
		case "isolatedUnrealizedPL":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			d.IsolatedUnrealizedPL = v
		case "assetMode":
			d.AssetMode = common.SafeStringCast(value)
		case "grant":
			d.Grant = common.SafeStringCast(value)
		default:
			// Неизвестные поля игнорируются
		}
	}
	return nil
}
