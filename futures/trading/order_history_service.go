package trading

import (
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"
)

// OrderHistoryService retrieves historical orders (filled, cancelled, rejected)
type OrderHistoryService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	startTime   string
	endTime     string
	pageSize    string
	lastEndId   string
}

func (s *OrderHistoryService) Symbol(symbol string) *OrderHistoryService {
	s.symbol = symbol
	return s
}

func (s *OrderHistoryService) ProductType(productType ProductType) *OrderHistoryService {
	s.productType = productType
	return s
}

func (s *OrderHistoryService) StartTime(startTime string) *OrderHistoryService {
	s.startTime = startTime
	return s
}

func (s *OrderHistoryService) EndTime(endTime string) *OrderHistoryService {
	s.endTime = endTime
	return s
}

func (s *OrderHistoryService) PageSize(pageSize string) *OrderHistoryService {
	s.pageSize = pageSize
	return s
}

func (s *OrderHistoryService) LastEndId(lastEndId string) *OrderHistoryService {
	s.lastEndId = lastEndId
	return s
}

func (s *OrderHistoryService) Do(ctx context.Context) (*OrderHistoryResponse, error) {
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
	var res *ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointOrderHistory, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var response *OrderHistoryResponse
	err = jsoniter.Unmarshal(res.Data, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type OrderHistoryResponse struct {
	// List of historical orders
	List []*HistoricalOrder `json:"list"`

	// End ID for pagination
	EndId string `json:"endId"`
}

type HistoricalOrder struct {
	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Order size
	Size string `json:"size"`

	// Order ID
	OrderId string `json:"orderId"`

	// Client order ID
	ClientOid string `json:"clientOid"`

	// Filled quantity
	FilledQty string `json:"filledQty"`

	// Filled amount/volume
	FilledAmount string `json:"filledAmount"`

	// Average fill price
	PriceAvg string `json:"priceAvg"`

	// Order price
	Price string `json:"price"`

	// Order state (filled, cancelled, rejected, etc.)
	State string `json:"state"`

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

	// Trading fee paid
	Fee string `json:"fee"`

	// Order type
	OrderType string `json:"orderType"`

	// Leverage
	Leverage string `json:"leverage"`

	// Margin mode
	MarginMode string `json:"marginMode"`

	// Reduce only flag
	ReduceOnly string `json:"reduceOnly"`

	// Enter point source
	EnterPointSource string `json:"enterPointSource"`

	// Trade side
	TradeSide string `json:"tradeSide"`

	// Position mode
	PosMode string `json:"posMode"`

	// Order source
	OrderSource string `json:"orderSource"`

	// Cancel reason (if cancelled)
	CancelReason string `json:"cancelReason"`

	// Create time
	CTime string `json:"cTime"`

	// Update time
	UTime string `json:"uTime"`

	// Preset stop surplus price
	PresetStopSurplusPrice string `json:"presetStopSurplusPrice"`

	// Preset stop loss price
	PresetStopLossPrice string `json:"presetStopLossPrice"`
}
