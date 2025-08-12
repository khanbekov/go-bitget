package trading

import (
	"context"
	"encoding/json"
)

// CreatePlanOrderService handles placing trigger/conditional orders (plan orders).
// Plan orders are triggered when the market price reaches the specified trigger price.
type CreatePlanOrderService struct {
	c ClientInterface

	// Required parameters
	symbol       string
	productType  ProductType
	planType     PlanType
	triggerPrice string
	triggerType  TriggerType
	side         SideType
	orderType    OrderType
	size         string

	// Optional parameters
	price       *string
	timeInForce *TimeInForceType
	clientOid   *string
	reduceOnly  *bool
	marginCoin  *string
}

// Symbol sets the trading symbol (e.g., "BTCUSDT").
func (s *CreatePlanOrderService) Symbol(symbol string) *CreatePlanOrderService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the plan order.
func (s *CreatePlanOrderService) ProductType(productType ProductType) *CreatePlanOrderService {
	s.productType = productType
	return s
}

// PlanType sets the type of plan order (normal_plan, track_plan, stop_loss, take_profit, stop_surplus).
func (s *CreatePlanOrderService) PlanType(planType PlanType) *CreatePlanOrderService {
	s.planType = planType
	return s
}

// TriggerPrice sets the price that triggers the plan order execution.
func (s *CreatePlanOrderService) TriggerPrice(triggerPrice string) *CreatePlanOrderService {
	s.triggerPrice = triggerPrice
	return s
}

// TriggerType sets how the trigger price is compared (fill_price or mark_price).
func (s *CreatePlanOrderService) TriggerType(triggerType TriggerType) *CreatePlanOrderService {
	s.triggerType = triggerType
	return s
}

// Side sets the order side (buy or sell).
func (s *CreatePlanOrderService) Side(side SideType) *CreatePlanOrderService {
	s.side = side
	return s
}

// OrderType sets the order type (limit or market) for when the plan order is triggered.
func (s *CreatePlanOrderService) OrderType(orderType OrderType) *CreatePlanOrderService {
	s.orderType = orderType
	return s
}

// Size sets the order size/quantity.
func (s *CreatePlanOrderService) Size(size string) *CreatePlanOrderService {
	s.size = size
	return s
}

// Price sets the limit price (required for limit orders, optional for market orders).
func (s *CreatePlanOrderService) Price(price string) *CreatePlanOrderService {
	s.price = &price
	return s
}

// TimeInForce sets the time in force for the order.
func (s *CreatePlanOrderService) TimeInForce(timeInForce TimeInForceType) *CreatePlanOrderService {
	s.timeInForce = &timeInForce
	return s
}

// ClientOid sets the client order ID for tracking.
func (s *CreatePlanOrderService) ClientOid(clientOid string) *CreatePlanOrderService {
	s.clientOid = &clientOid
	return s
}

// ReduceOnly sets whether this is a reduce-only order.
func (s *CreatePlanOrderService) ReduceOnly(reduceOnly bool) *CreatePlanOrderService {
	s.reduceOnly = &reduceOnly
	return s
}

// MarginCoin sets the margin coin for the order.
func (s *CreatePlanOrderService) MarginCoin(marginCoin string) *CreatePlanOrderService {
	s.marginCoin = &marginCoin
	return s
}

// CreatePlanOrderResponse represents the response from placing a plan order.
type CreatePlanOrderResponse struct {
	OrderId   string `json:"orderId"`   // Plan order ID
	ClientOid string `json:"clientOid"` // Client order ID
}

// Do executes the create plan order request.
func (s *CreatePlanOrderService) Do(ctx context.Context) (*CreatePlanOrderResponse, error) {
	// Build request body
	params := map[string]interface{}{
		"symbol":       s.symbol,
		"productType":  string(s.productType),
		"planType":     string(s.planType),
		"triggerPrice": s.triggerPrice,
		"triggerType":  string(s.triggerType),
		"side":         string(s.side),
		"orderType":    string(s.orderType),
		"size":         s.size,
	}

	// Add optional parameters
	if s.price != nil {
		params["price"] = *s.price
	}
	if s.timeInForce != nil {
		params["timeInForce"] = string(*s.timeInForce)
	}
	if s.clientOid != nil {
		params["clientOid"] = *s.clientOid
	}
	if s.reduceOnly != nil {
		params["reduceOnly"] = *s.reduceOnly
	}
	if s.marginCoin != nil {
		params["marginCoin"] = *s.marginCoin
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "POST", EndpointCreatePlanOrder, nil, body, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result CreatePlanOrderResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
