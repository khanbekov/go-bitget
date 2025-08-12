package uta

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/common"
)

// CancelOrderService cancels an existing order
type CancelOrderService struct {
	c         ClientInterface
	symbol    *string
	category  *string
	orderId   *string
	clientOid *string
}

// Symbol sets the trading symbol (required)
func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = &symbol
	return s
}

// Category sets the product category (required)
func (s *CancelOrderService) Category(category string) *CancelOrderService {
	s.category = &category
	return s
}

// OrderId sets the order ID (either orderId or clientOid is required)
func (s *CancelOrderService) OrderId(orderId string) *CancelOrderService {
	s.orderId = &orderId
	return s
}

// ClientOid sets the client order ID (either orderId or clientOid is required)
func (s *CancelOrderService) ClientOid(clientOid string) *CancelOrderService {
	s.clientOid = &clientOid
	return s
}

// Do executes the cancel order request
func (s *CancelOrderService) Do(ctx context.Context) (*Order, error) {
	if s.symbol == nil {
		return nil, common.NewMissingParameterError("symbol")
	}
	if s.category == nil {
		return nil, common.NewMissingParameterError("category")
	}
	if s.orderId == nil && s.clientOid == nil {
		return nil, common.NewMissingParameterError("orderId or clientOid")
	}

	params := map[string]interface{}{
		"symbol":   *s.symbol,
		"category": *s.category,
	}

	if s.orderId != nil {
		params["orderId"] = *s.orderId
	}
	if s.clientOid != nil {
		params["clientOid"] = *s.clientOid
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", EndpointTradeCancelOrder, nil, body, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := common.UnmarshalJSON(res.Data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
