package trading

import (
	"context"
	"encoding/json"
)

// ModifyPlanOrderService handles modifying existing trigger/conditional orders (plan orders).
// This service allows you to modify trigger price, order type, size, and price of existing plan orders.
type ModifyPlanOrderService struct {
	c ClientInterface

	// Required parameters
	orderId      string
	orderType    OrderType
	triggerPrice string

	// Optional parameters
	triggerType *TriggerType
	size        *string
	price       *string
}

// OrderId sets the plan order ID to modify.
func (s *ModifyPlanOrderService) OrderId(orderId string) *ModifyPlanOrderService {
	s.orderId = orderId
	return s
}

// OrderType sets the order type (limit or market) for when the plan order is triggered.
func (s *ModifyPlanOrderService) OrderType(orderType OrderType) *ModifyPlanOrderService {
	s.orderType = orderType
	return s
}

// TriggerPrice sets the new trigger price for the plan order.
func (s *ModifyPlanOrderService) TriggerPrice(triggerPrice string) *ModifyPlanOrderService {
	s.triggerPrice = triggerPrice
	return s
}

// TriggerType sets how the trigger price is compared (fill_price or mark_price).
func (s *ModifyPlanOrderService) TriggerType(triggerType TriggerType) *ModifyPlanOrderService {
	s.triggerType = &triggerType
	return s
}

// Size sets the new order size/quantity.
func (s *ModifyPlanOrderService) Size(size string) *ModifyPlanOrderService {
	s.size = &size
	return s
}

// Price sets the new limit price (required for limit orders).
func (s *ModifyPlanOrderService) Price(price string) *ModifyPlanOrderService {
	s.price = &price
	return s
}

// ModifyPlanOrderResponse represents the response from modifying a plan order.
type ModifyPlanOrderResponse struct {
	OrderId   string `json:"orderId"`   // Modified plan order ID
	ClientOid string `json:"clientOid"` // Client order ID
}

// Do executes the modify plan order request.
func (s *ModifyPlanOrderService) Do(ctx context.Context) (*ModifyPlanOrderResponse, error) {
	// Build request body
	params := map[string]interface{}{
		"orderId":      s.orderId,
		"orderType":    string(s.orderType),
		"triggerPrice": s.triggerPrice,
	}

	// Add optional parameters
	if s.triggerType != nil {
		params["triggerType"] = string(*s.triggerType)
	}
	if s.size != nil {
		params["size"] = *s.size
	}
	if s.price != nil {
		params["price"] = *s.price
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "POST", EndpointModifyPlanOrder, nil, body, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result ModifyPlanOrderResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
