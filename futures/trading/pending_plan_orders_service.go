package trading

import (
	"context"
	"encoding/json"
	"net/url"
)

// PendingPlanOrdersService handles retrieving pending trigger/conditional orders (plan orders).
// This service allows you to query all pending plan orders for a specific symbol and product type.
type PendingPlanOrdersService struct {
	c ClientInterface

	// Required parameters
	symbol      string
	productType ProductType
	planType    PlanType

	// Optional parameters
	limit      *string
	idLessThan *string
}

// Symbol sets the trading symbol to filter plan orders (e.g., "BTCUSDT").
func (s *PendingPlanOrdersService) Symbol(symbol string) *PendingPlanOrdersService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type to filter plan orders.
func (s *PendingPlanOrdersService) ProductType(productType ProductType) *PendingPlanOrdersService {
	s.productType = productType
	return s
}

// PlanType sets the type of plan orders to retrieve (normal_plan, track_plan, stop_loss, take_profit, stop_surplus).
func (s *PendingPlanOrdersService) PlanType(planType PlanType) *PendingPlanOrdersService {
	s.planType = planType
	return s
}

// Limit sets the number of records to return (default 100, max 100).
func (s *PendingPlanOrdersService) Limit(limit string) *PendingPlanOrdersService {
	s.limit = &limit
	return s
}

// IdLessThan sets the cursor for pagination (return records with ID less than this value).
func (s *PendingPlanOrdersService) IdLessThan(idLessThan string) *PendingPlanOrdersService {
	s.idLessThan = &idLessThan
	return s
}

// PendingPlanOrder represents a pending plan order.
type PendingPlanOrder struct {
	OrderId      string `json:"orderId"`      // Plan order ID
	ClientOid    string `json:"clientOid"`    // Client order ID
	Symbol       string `json:"symbol"`       // Trading symbol
	ProductType  string `json:"productType"`  // Product type
	PlanType     string `json:"planType"`     // Plan type
	TriggerPrice string `json:"triggerPrice"` // Trigger price
	TriggerType  string `json:"triggerType"`  // Trigger type
	Side         string `json:"side"`         // Order side
	OrderType    string `json:"orderType"`    // Order type
	Size         string `json:"size"`         // Order size
	Price        string `json:"price"`        // Order price
	TimeInForce  string `json:"timeInForce"`  // Time in force
	Status       string `json:"status"`       // Order status
	CreateTime   string `json:"createTime"`   // Creation time
	UpdateTime   string `json:"updateTime"`   // Last update time
	ReduceOnly   string `json:"reduceOnly"`   // Reduce only flag
	MarginCoin   string `json:"marginCoin"`   // Margin coin
}

// PendingPlanOrdersResponse represents the response from getting pending plan orders.
type PendingPlanOrdersResponse struct {
	Orders []PendingPlanOrder `json:"orders"` // List of pending plan orders
}

// Do executes the get pending plan orders request.
func (s *PendingPlanOrdersService) Do(ctx context.Context) ([]*PendingPlanOrder, error) {
	// Build query parameters
	queryParams := url.Values{}
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))
	queryParams.Set("planType", string(s.planType))

	// Add optional parameters
	if s.limit != nil {
		queryParams.Set("limit", *s.limit)
	}
	if s.idLessThan != nil {
		queryParams.Set("idLessThan", *s.idLessThan)
	}

	// Make API call
	res, _, err := s.c.CallAPI(ctx, "GET", EndpointPendingPlanOrders, queryParams, nil, true)
	if err != nil {
		return nil, err
	}

	// Parse response
	var orders []*PendingPlanOrder
	if err := json.Unmarshal(res.Data, &orders); err != nil {
		return nil, err
	}

	return orders, nil
}
