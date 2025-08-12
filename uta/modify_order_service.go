package uta

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/common"
)

// ModifyOrderService modifies an existing order
type ModifyOrderService struct {
	c            ClientInterface
	symbol       *string
	category     *string
	orderId      *string
	clientOid    *string
	newSize      *string
	newPrice     *string
	newClientOid *string
}

// Symbol sets the trading symbol (required)
func (s *ModifyOrderService) Symbol(symbol string) *ModifyOrderService {
	s.symbol = &symbol
	return s
}

// Category sets the product category (required)
func (s *ModifyOrderService) Category(category string) *ModifyOrderService {
	s.category = &category
	return s
}

// OrderId sets the order ID (either orderId or clientOid is required)
func (s *ModifyOrderService) OrderId(orderId string) *ModifyOrderService {
	s.orderId = &orderId
	return s
}

// ClientOid sets the client order ID (either orderId or clientOid is required)
func (s *ModifyOrderService) ClientOid(clientOid string) *ModifyOrderService {
	s.clientOid = &clientOid
	return s
}

// NewSize sets the new order size (optional)
func (s *ModifyOrderService) NewSize(size string) *ModifyOrderService {
	s.newSize = &size
	return s
}

// NewPrice sets the new order price (optional)
func (s *ModifyOrderService) NewPrice(price string) *ModifyOrderService {
	s.newPrice = &price
	return s
}

// NewClientOid sets the new client order ID (optional)
func (s *ModifyOrderService) NewClientOid(clientOid string) *ModifyOrderService {
	s.newClientOid = &clientOid
	return s
}

// Do executes the modify order request
func (s *ModifyOrderService) Do(ctx context.Context) (*Order, error) {
	if s.symbol == nil {
		return nil, common.NewMissingParameterError("symbol")
	}
	if s.category == nil {
		return nil, common.NewMissingParameterError("category")
	}
	if s.orderId == nil && s.clientOid == nil {
		return nil, common.NewMissingParameterError("orderId or clientOid")
	}
	if s.newSize == nil && s.newPrice == nil {
		return nil, common.NewMissingParameterError("newSize or newPrice - at least one modification is required")
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
	if s.newSize != nil {
		params["qty"] = *s.newSize
	}
	if s.newPrice != nil {
		params["price"] = *s.newPrice
	}
	if s.newClientOid != nil {
		params["newClientOid"] = *s.newClientOid
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", EndpointTradeModifyOrder, nil, body, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := common.UnmarshalJSON(res.Data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
