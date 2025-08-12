package trading

import (
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"
)

// FillHistoryService retrieves order execution/fill history
type FillHistoryService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	orderId     string
	startTime   string
	endTime     string
	pageSize    string
	lastEndId   string
}

func (s *FillHistoryService) Symbol(symbol string) *FillHistoryService {
	s.symbol = symbol
	return s
}

func (s *FillHistoryService) ProductType(productType ProductType) *FillHistoryService {
	s.productType = productType
	return s
}

func (s *FillHistoryService) OrderId(orderId string) *FillHistoryService {
	s.orderId = orderId
	return s
}

func (s *FillHistoryService) StartTime(startTime string) *FillHistoryService {
	s.startTime = startTime
	return s
}

func (s *FillHistoryService) EndTime(endTime string) *FillHistoryService {
	s.endTime = endTime
	return s
}

func (s *FillHistoryService) PageSize(pageSize string) *FillHistoryService {
	s.pageSize = pageSize
	return s
}

func (s *FillHistoryService) LastEndId(lastEndId string) *FillHistoryService {
	s.lastEndId = lastEndId
	return s
}

func (s *FillHistoryService) Do(ctx context.Context) (*FillHistoryResponse, error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("productType", string(s.productType))
	if s.symbol != "" {
		queryParams.Set("symbol", s.symbol)
	}
	if s.orderId != "" {
		queryParams.Set("orderId", s.orderId)
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

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointFillHistory, queryParams, nil, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var response *FillHistoryResponse
	err = jsoniter.Unmarshal(res.Data, &response)

	if err != nil {
		return nil, err
	}

	return response, nil
}

type FillHistoryResponse struct {
	// List of fill records
	List []*FillRecord `json:"list"`

	// End ID for pagination
	EndId string `json:"endId"`
}

type FillRecord struct {
	// Trade ID
	TradeId string `json:"tradeId"`

	// Order ID
	OrderId string `json:"orderId"`

	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Fill size/quantity
	Size string `json:"size"`

	// Fill price
	Price string `json:"price"`

	// Order side (buy/sell)
	Side string `json:"side"`

	// Fill amount (size * price)
	Amount string `json:"amount"`

	// Trading fee paid
	Fee string `json:"fee"`

	// Fee coin
	FeeCcy string `json:"feeCcy"`

	// Profits from this fill
	Profit string `json:"profit"`

	// Position side
	PosSide string `json:"posSide"`

	// Margin coin
	MarginCoin string `json:"marginCoin"`

	// Order type
	OrderType string `json:"orderType"`

	// Margin mode
	MarginMode string `json:"marginMode"`

	// Trade side
	TradeSide string `json:"tradeSide"`

	// Position mode
	PosMode string `json:"posMode"`

	// Liquidity role (maker/taker)
	Role string `json:"role"`

	// Fill time
	CTime string `json:"cTime"`
}
