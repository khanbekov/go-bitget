package position

import (
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"github.com/khanbekov/go-bitget/common"
	"golang.org/x/net/context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// HistoryPositionsService retrieves closed/historical positions
type HistoryPositionsService struct {
	c           futures.ClientInterface
	symbol      string
	productType futures.ProductType
	startTime   string
	endTime     string
	pageSize    string
	lastEndId   string
}

func (s *HistoryPositionsService) Symbol(symbol string) *HistoryPositionsService {
	s.symbol = symbol
	return s
}

func (s *HistoryPositionsService) ProductType(productType futures.ProductType) *HistoryPositionsService {
	s.productType = productType
	return s
}

func (s *HistoryPositionsService) StartTime(startTime string) *HistoryPositionsService {
	s.startTime = startTime
	return s
}

func (s *HistoryPositionsService) EndTime(endTime string) *HistoryPositionsService {
	s.endTime = endTime
	return s
}

func (s *HistoryPositionsService) PageSize(pageSize string) *HistoryPositionsService {
	s.pageSize = pageSize
	return s
}

func (s *HistoryPositionsService) LastEndId(lastEndId string) *HistoryPositionsService {
	s.lastEndId = lastEndId
	return s
}

func (s *HistoryPositionsService) Do(ctx context.Context) (*HistoryPositionsResponse, error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("productType", string(s.productType))
	if s.symbol != "" {
		queryParams.Set("symbol", s.symbol)
	}
	if s.startTime != "" {
		queryParams.Set("startTime", s.startTime)
	}
	if s.endTime != "" {
		queryParams.Set("endTime", s.endTime)
	}
	if s.pageSize != "" {
		queryParams.Set("pageSize", s.pageSize)
	}
	if s.lastEndId != "" {
		queryParams.Set("lastEndId", s.lastEndId)
	}

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointHistoryPositions, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var response *HistoryPositionsResponse
	err = jsoniter.Unmarshal(res.Data, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type HistoryPositionsResponse struct {
	// List of historical positions
	List []*HistoryPosition `json:"list"`

	// End ID for pagination
	EndId string `json:"endId"`
}

type HistoryPosition struct {
	// Position ID
	PositionId string `json:"positionId"`

	// Margin coin used for margin trading
	MarginCoin string `json:"marginCoin"`

	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Position side: long/short
	HoldSide string `json:"holdSide"`

	// Position size (positive number)
	Size float64 `json:"size"`

	// Mark price at close
	MarkPrice float64 `json:"markPrice"`

	// Average opening price
	AverageOpenPrice float64 `json:"averageOpenPrice"`

	// Average closing price
	AverageClosePrice float64 `json:"averageClosePrice"`

	// Realized profit and loss
	RealizedPL float64 `json:"realizedPL"`

	// Total funding fee
	TotalFee float64 `json:"totalFee"`

	// Margin mode: isolated/crossed
	MarginMode string `json:"marginMode"`

	// Position mode: one_way_mode/hedge_mode
	PosMode string `json:"posMode"`

	// Position opening time
	OpenTime int64 `json:"openTime"`

	// Position closing time
	CloseTime int64 `json:"closeTime"`

	// Asset mode
	AssetMode string `json:"assetMode"`
}

// UnmarshalJSON realization interface json.Unmarshaler for HistoryPosition
func (hp *HistoryPosition) UnmarshalJSON(data []byte) error {
	tmp := map[string]interface{}{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	for key, value := range tmp {
		switch key {
		case "positionId":
			hp.PositionId = common.SafeStringCast(value)
		case "marginCoin":
			hp.MarginCoin = common.SafeStringCast(value)
		case "symbol":
			hp.Symbol = common.SafeStringCast(value)
		case "holdSide":
			hp.HoldSide = common.SafeStringCast(value)
		case "size":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			hp.Size = v
		case "markPrice":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			hp.MarkPrice = v
		case "averageOpenPrice":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			hp.AverageOpenPrice = v
		case "averageClosePrice":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			hp.AverageClosePrice = v
		case "realizedPL":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			hp.RealizedPL = v
		case "totalFee":
			v, err := common.ConvertToFloat64(value)
			if err != nil {
				return err
			}
			hp.TotalFee = v
		case "marginMode":
			hp.MarginMode = common.SafeStringCast(value)
		case "posMode":
			hp.PosMode = common.SafeStringCast(value)
		case "openTime":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			hp.OpenTime = v
		case "closeTime":
			v, err := common.ConvertToInt64(value)
			if err != nil {
				return err
			}
			hp.CloseTime = v
		case "assetMode":
			hp.AssetMode = common.SafeStringCast(value)
		default:
			// Unknown fields are ignored
		}
	}
	return nil
}
