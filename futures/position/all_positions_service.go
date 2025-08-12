package position

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/common"
	"golang.org/x/net/context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// AllPositionsService retrieves all open positions for the account
type AllPositionsService struct {
	c           futures.ClientInterface
	productType futures.ProductType
	marginCoin  string
}

func (s *AllPositionsService) ProductType(productType futures.ProductType) *AllPositionsService {
	s.productType = productType
	return s
}

func (s *AllPositionsService) MarginCoin(marginCoin string) *AllPositionsService {
	s.marginCoin = marginCoin
	return s
}

func (s *AllPositionsService) Do(ctx context.Context) ([]*Position, error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("productType", string(s.productType))
	if s.marginCoin != "" {
		queryParams.Set("marginCoin", s.marginCoin)
	}

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointAllPositions, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var positions []*Position
	err = jsoniter.Unmarshal(res.Data, &positions)

	if err != nil {
		return nil, err
	}

	return positions, nil
}

type Position struct {
	// Margin coin used for margin trading
	MarginCoin string `json:"marginCoin"`

	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Position side: long/short
	HoldSide futures.HoldSideType `json:"holdSide"`

	// Position size (positive number)
	Size float64 `json:"size"`

	// Mark price
	MarkPrice float64 `json:"markPrice"`

	// Position value in margin coin
	PositionValue float64 `json:"positionValue"`

	// Average opening price
	AverageOpenPrice float64 `json:"averageOpenPrice"`

	// Unrealized profit and loss
	UnrealizedPL float64 `json:"unrealizedPL"`

	// Unrealized profit and loss rate
	UnrealizedPLR float64 `json:"unrealizedPLR"`

	// Position margin
	Margin float64 `json:"margin"`

	// Available size for closing position
	Available float64 `json:"available"`

	// Cross margin leverage
	CrossedLeverage float64 `json:"crossedLeverage"`

	// Isolated margin leverage
	IsolatedLeverage float64 `json:"isolatedLeverage"`

	// Margin mode: isolated/crossed
	MarginMode string `json:"marginMode"`

	// Position mode: one_way_mode/hedge_mode
	PosMode string `json:"posMode"`

	// Margin ratio
	MarginRatio float64 `json:"marginRatio"`

	// Maintenance margin ratio
	MaintenanceMarginRatio float64 `json:"maintenanceMarginRatio"`

	// Created timestamp
	Ctime int64 `json:"ctime"`

	// Updated timestamp
	Utime int64 `json:"utime"`

	// Break-even price
	BreakEvenPrice float64 `json:"breakEvenPrice"`

	// Total funding fee
	TotalFee float64 `json:"totalFee"`

	// Deducted funding fee
	DeductedFee float64 `json:"deductedFee"`

	// Auto margin increase flag
	AutoMargin string `json:"autoMargin"`

	// Asset mode
	AssetMode string `json:"assetMode"`
}

// UnmarshalJSON realization interface json.Unmarshaler for Position
func (p *Position) UnmarshalJSON(data []byte) error {
	tmp := map[string]interface{}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for key, value := range tmp {
		switch key {
		case "marginCoin":
			p.MarginCoin = value.(string)
		case "symbol":
			p.Symbol = value.(string)
		case "holdSide":
			p.HoldSide = futures.HoldSideType(value.(string))
		case "size":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.Size = v
		case "markPrice":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.MarkPrice = v
		case "positionValue":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.PositionValue = v
		case "averageOpenPrice":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.AverageOpenPrice = v
		case "unrealizedPL":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.UnrealizedPL = v
		case "unrealizedPLR":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.UnrealizedPLR = v
		case "margin":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.Margin = v
		case "available":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.Available = v
		case "crossedLeverage":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.CrossedLeverage = v
		case "isolatedLeverage":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.IsolatedLeverage = v
		case "marginMode":
			p.MarginMode = value.(string)
		case "posMode":
			p.PosMode = value.(string)
		case "marginRatio":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.MarginRatio = v
		case "maintenanceMarginRatio":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.MaintenanceMarginRatio = v
		case "ctime":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			p.Ctime = v
		case "utime":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			p.Utime = v
		case "breakEvenPrice":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.BreakEvenPrice = v
		case "totalFee":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.TotalFee = v
		case "deductedFee":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			p.DeductedFee = v
		case "autoMargin":
			p.AutoMargin = value.(string)
		case "assetMode":
			p.AssetMode = value.(string)
		default:
			// Unknown fields are ignored
		}
	}
	return nil
}
