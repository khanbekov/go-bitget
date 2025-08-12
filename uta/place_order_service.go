package uta

import (
	"context"
	"encoding/json"

	"github.com/khanbekov/go-bitget/common"
)

// PlaceOrderService places a new order
type PlaceOrderService struct {
	c            ClientInterface
	symbol       *string
	category     *string
	side         *string
	orderType    *string
	size         *string
	price        *string
	clientOid    *string
	timeInForce  *string
	reduceOnly   *string
	positionSide *string
	stp          *string
}

// Symbol sets the trading symbol (required)
func (s *PlaceOrderService) Symbol(symbol string) *PlaceOrderService {
	s.symbol = &symbol
	return s
}

// Category sets the product category (required)
func (s *PlaceOrderService) Category(category string) *PlaceOrderService {
	s.category = &category
	return s
}

// Side sets the order side (required): "buy" or "sell"
func (s *PlaceOrderService) Side(side string) *PlaceOrderService {
	s.side = &side
	return s
}

// OrderType sets the order type (required): "limit" or "market"
func (s *PlaceOrderService) OrderType(orderType string) *PlaceOrderService {
	s.orderType = &orderType
	return s
}

// Size sets the order size (required)
func (s *PlaceOrderService) Size(size string) *PlaceOrderService {
	s.size = &size
	return s
}

// Price sets the order price (required for limit orders)
func (s *PlaceOrderService) Price(price string) *PlaceOrderService {
	s.price = &price
	return s
}

// ClientOid sets the client order ID (optional)
func (s *PlaceOrderService) ClientOid(clientOid string) *PlaceOrderService {
	s.clientOid = &clientOid
	return s
}

// TimeInForce sets the time in force (optional)
func (s *PlaceOrderService) TimeInForce(timeInForce string) *PlaceOrderService {
	s.timeInForce = &timeInForce
	return s
}

// ReduceOnly sets reduce only flag (optional)
func (s *PlaceOrderService) ReduceOnly(reduceOnly string) *PlaceOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// PositionSide sets the position side (optional, for futures in hedge mode)
func (s *PlaceOrderService) PositionSide(positionSide string) *PlaceOrderService {
	s.positionSide = &positionSide
	return s
}

// STP sets the self trade prevention mode (optional)
func (s *PlaceOrderService) STP(stp string) *PlaceOrderService {
	s.stp = &stp
	return s
}

// Do executes the place order request
func (s *PlaceOrderService) Do(ctx context.Context) (*Order, error) {
	if s.symbol == nil {
		return nil, common.NewMissingParameterError("symbol")
	}
	if s.category == nil {
		return nil, common.NewMissingParameterError("category")
	}
	if s.side == nil {
		return nil, common.NewMissingParameterError("side")
	}
	if s.orderType == nil {
		return nil, common.NewMissingParameterError("orderType")
	}
	if s.size == nil {
		return nil, common.NewMissingParameterError("size")
	}

	params := map[string]interface{}{
		"symbol":    *s.symbol,
		"category":  *s.category,
		"side":      *s.side,
		"orderType": *s.orderType,
		"qty":       *s.size, // UTA API uses 'qty' instead of 'size'
	}

	if s.price != nil {
		params["price"] = *s.price
	}
	if s.clientOid != nil {
		params["clientOid"] = *s.clientOid
	}
	if s.timeInForce != nil {
		params["timeInForce"] = *s.timeInForce
	}
	if s.reduceOnly != nil {
		params["reduceOnly"] = *s.reduceOnly
	}
	if s.positionSide != nil {
		params["posSide"] = *s.positionSide
	}
	if s.stp != nil {
		params["stp"] = *s.stp
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", EndpointTradePlaceOrder, nil, body, true)
	if err != nil {
		return nil, err
	}

	var order Order
	if err := common.UnmarshalJSON(res.Data, &order); err != nil {
		return nil, err
	}

	return &order, nil
}
