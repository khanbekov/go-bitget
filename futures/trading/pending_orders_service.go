package trading

import (
	"context"
	"encoding/json"
	jsoniter "github.com/json-iterator/go"
	"net/url"
)

// PendingOrdersService retrieves all open/pending orders
type PendingOrdersService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	marginCoin  string
}

func (s *PendingOrdersService) Symbol(symbol string) *PendingOrdersService {
	s.symbol = symbol
	return s
}

func (s *PendingOrdersService) ProductType(productType ProductType) *PendingOrdersService {
	s.productType = productType
	return s
}

func (s *PendingOrdersService) MarginCoin(marginCoin string) *PendingOrdersService {
	s.marginCoin = marginCoin
	return s
}

func (s *PendingOrdersService) Do(ctx context.Context) ([]*PendingOrder, error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("productType", string(s.productType))
	if s.symbol != "" {
		queryParams.Set("symbol", s.symbol)
	}
	if s.marginCoin != "" {
		queryParams.Set("marginCoin", s.marginCoin)
	}

	// Make request to API
	var res *ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointPendingOrders, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var response PendingOrdersResponse
	err = jsoniter.Unmarshal(res.Data, &response)

	if err != nil {
		return nil, err
	}

	// Return empty slice if entrustedList is null or empty
	if response.EntrustedList == nil {
		return []*PendingOrder{}, nil
	}

	return response.EntrustedList, nil
}

// PendingOrdersResponse represents the API response structure
type PendingOrdersResponse struct {
	EntrustedList []*PendingOrder `json:"entrustedList"`
	EndId         *string         `json:"endId"` // Pointer to handle null values
}

type PendingOrder struct {
	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Order size
	Size string `json:"size"`

	// Order ID
	OrderId string `json:"orderId"`

	// Client order ID
	ClientOid string `json:"clientOid"`

	// Base volume
	BaseVolume string `json:"baseVolume"`

	// Fee
	Fee string `json:"fee"`

	// Order price
	Price string `json:"price"`

	// Average filled price
	PriceAvg string `json:"priceAvg"`

	// Order status
	Status string `json:"status"`

	// Order side (buy/sell)
	Side string `json:"side"`

	// Time in force
	Force string `json:"force"`

	// Total profits
	TotalProfits string `json:"totalProfits"`

	// Position side
	PosSide string `json:"posSide"`

	// Margin coin
	MarginCoin string `json:"marginCoin"`

	// Quote volume
	QuoteVolume string `json:"quoteVolume"`

	// Leverage
	Leverage string `json:"leverage"`

	// Margin mode
	MarginMode string `json:"marginMode"`

	// Enter point source
	EnterPointSource string `json:"enterPointSource"`

	// Trade side
	TradeSide string `json:"tradeSide"`

	// Position mode
	PosMode string `json:"posMode"`

	// Order type
	OrderType string `json:"orderType"`

	// Order source
	OrderSource string `json:"orderSource"`

	// Created time
	CTime string `json:"cTime"`

	// Updated time
	UTime string `json:"uTime"`

	// Preset stop surplus price
	PresetStopSurplusPrice string `json:"presetStopSurplusPrice"`

	// Preset stop surplus type
	PresetStopSurplusType string `json:"presetStopSurplusType"`

	// Preset stop surplus execute price
	PresetStopSurplusExecutePrice string `json:"presetStopSurplusExecutePrice"`

	// Preset stop loss price
	PresetStopLossPrice string `json:"presetStopLossPrice"`

	// Preset stop loss type
	PresetStopLossType string `json:"presetStopLossType"`

	// Preset stop loss execute price
	PresetStopLossExecutePrice string `json:"presetStopLossExecutePrice"`
}

// UnmarshalJSON implements custom JSON unmarshaling to handle null values
func (r *PendingOrdersResponse) UnmarshalJSON(data []byte) error {
	type Alias PendingOrdersResponse
	aux := &struct {
		EntrustedList json.RawMessage `json:"entrustedList"`
		*Alias
	}{
		Alias: (*Alias)(r),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Handle null entrustedList
	if string(aux.EntrustedList) == "null" || len(aux.EntrustedList) == 0 {
		r.EntrustedList = []*PendingOrder{}
		return nil
	}

	// Parse entrustedList if not null
	return json.Unmarshal(aux.EntrustedList, &r.EntrustedList)
}
