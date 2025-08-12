package trading

import (
	"context"
	"encoding/json"
)

// CancelPlanOrderService handles canceling trigger/conditional orders (plan orders).
// This service allows you to cancel existing plan orders by order ID and plan type.
type CancelPlanOrderService struct {
	c ClientInterface

	// Required parameters
	orderId  string
	planType PlanType
}

// OrderId sets the plan order ID to cancel.
func (s *CancelPlanOrderService) OrderId(orderId string) *CancelPlanOrderService {
	s.orderId = orderId
	return s
}

// PlanType sets the type of plan order to cancel (normal_plan, track_plan, stop_loss, take_profit, stop_surplus).
func (s *CancelPlanOrderService) PlanType(planType PlanType) *CancelPlanOrderService {
	s.planType = planType
	return s
}

// CancelPlanOrderResponse represents the response from canceling a plan order.
type CancelPlanOrderResponse struct {
	OrderId   string `json:"orderId"`   // Canceled plan order ID
	ClientOid string `json:"clientOid"` // Client order ID
}

// Do executes the cancel plan order request.
func (s *CancelPlanOrderService) Do(ctx context.Context) (*CancelPlanOrderResponse, error) {
	// Build request body
	params := map[string]interface{}{
		"orderId":  s.orderId,
		"planType": string(s.planType),
	}

	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "POST", EndpointCancelPlanOrder, nil, body, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var result CancelPlanOrderResponse
	if err := json.Unmarshal(res.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
