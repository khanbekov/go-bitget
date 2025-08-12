package trading

import (
	"context"
	"fmt"
	"net/url"

	jsoniter "github.com/json-iterator/go"
)

// OrderDetail represents detailed order information
type OrderDetail struct {
	Symbol                 string `json:"symbol"`
	Size                   string `json:"size"`
	OrderId                string `json:"orderId"`
	ClientOid              string `json:"clientOid"`
	BaseVolume             string `json:"baseVolume"`
	PriceAvg               string `json:"priceAvg"`
	Fee                    string `json:"fee"`
	Price                  string `json:"price"`
	State                  string `json:"state"`
	Side                   string `json:"side"`
	Force                  string `json:"force"`
	TotalProfits           string `json:"totalProfits"`
	PosSide                string `json:"posSide"`
	MarginCoin             string `json:"marginCoin"`
	PresetStopSurplusPrice string `json:"presetStopSurplusPrice"`
	PresetStopLossPrice    string `json:"presetStopLossPrice"`
	QuoteVolume            string `json:"quoteVolume"`
	OrderType              string `json:"orderType"`
	Leverage               string `json:"leverage"`
	MarginMode             string `json:"marginMode"`
	ReduceOnly             string `json:"reduceOnly"`
	EnterPointSource       string `json:"enterPointSource"`
	TradeSide              string `json:"tradeSide"`
	PosMode                string `json:"posMode"`
	OrderSource            string `json:"orderSource"`
	CancelReason           string `json:"cancelReason"`
	CTime                  string `json:"cTime"`
	UTime                  string `json:"uTime"`
}

// GetOrderDetailsService provides methods to retrieve order details
type GetOrderDetailsService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	orderId     string
	clientOid   string
}

// Symbol sets the trading pair (required)
func (s *GetOrderDetailsService) Symbol(symbol string) *GetOrderDetailsService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type (required)
func (s *GetOrderDetailsService) ProductType(productType ProductType) *GetOrderDetailsService {
	s.productType = productType
	return s
}

// OrderId sets the order ID (either orderId or clientOid required)
func (s *GetOrderDetailsService) OrderId(orderId string) *GetOrderDetailsService {
	s.orderId = orderId
	return s
}

// ClientOid sets the custom order ID (either orderId or clientOid required)
func (s *GetOrderDetailsService) ClientOid(clientOid string) *GetOrderDetailsService {
	s.clientOid = clientOid
	return s
}

// checkRequiredParams validates required parameters
func (s *GetOrderDetailsService) checkRequiredParams() error {
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.orderId == "" && s.clientOid == "" {
		return fmt.Errorf("either orderId or clientOid must be provided")
	}
	return nil
}

// Do sends the order detail request
func (s *GetOrderDetailsService) Do(ctx context.Context) (orderDetail *OrderDetail, err error) {
	if err := s.checkRequiredParams(); err != nil {
		return nil, err
	}

	queryParams := url.Values{}
	queryParams.Add("symbol", s.symbol)
	queryParams.Add("productType", string(s.productType))
	if s.orderId != "" {
		queryParams.Add("orderId", s.orderId)
	}
	if s.clientOid != "" {
		queryParams.Add("clientOid", s.clientOid)
	}

	res, _, err := s.c.CallAPI(ctx, "GET", EndpointOrderDetails, queryParams, nil, true)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data OrderDetail `json:"data"`
	}
	if err := jsoniter.Unmarshal(res.Data, &wrapper); err != nil {
		return nil, err
	}
	return &wrapper.Data, nil
}
