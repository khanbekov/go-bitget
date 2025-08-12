package trading

import (
	"github.com/khanbekov/go-bitget/common/client"
)

// Re-export common types to avoid importing futures package
type (
	ClientInterface = client.ClientInterface
	ApiResponse     = client.ApiResponse
)

// Futures-specific enums (duplicated to avoid import cycle)
type ProductType string

const (
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	ProductTypeCoinFutures ProductType = "COIN-FUTURES"
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
)

// Order types and sides
type OrderType string

const (
	OrderTypeLimit  OrderType = "limit"
	OrderTypeMarket OrderType = "market"
)

type Side string
type SideType = Side // Alias for backward compatibility

const (
	SideBuy  Side = "buy"
	SideSell Side = "sell"

	// Alternative side constants for backward compatibility
	SideTypeBuy  SideType = "BUY"
	SideTypeSell SideType = "SELL"
)

// Trigger types for plan orders
type TriggerType string

const (
	TriggerTypeFillPrice TriggerType = "fill_price" // Trigger when fill price reaches trigger price
	TriggerTypeMarkPrice TriggerType = "mark_price" // Trigger when mark price reaches trigger price
)

type HoldSide string

const (
	HoldSideLong  HoldSide = "long"
	HoldSideShort HoldSide = "short"
)

// Plan order types
type PlanType string

const (
	PlanTypeNormalPlan  PlanType = "normal_plan"  // Normal trigger order
	PlanTypeTrackPlan   PlanType = "track_plan"   // Trailing stop order
	PlanTypeStopLoss    PlanType = "stop_loss"    // Stop loss order
	PlanTypeTakeProfit  PlanType = "take_profit"  // Take profit order
	PlanTypeStopSurplus PlanType = "stop_surplus" // Stop surplus order
)

// Time in force options
type TimeInForce string
type TimeInForceType = TimeInForce // Alias for backward compatibility

const (
	TimeInForceGTC      TimeInForce = "GTC"       // Good Till Cancel
	TimeInForceIOC      TimeInForce = "IOC"       // Immediate or Cancel
	TimeInForceFOK      TimeInForce = "FOK"       // Fill or Kill
	TimeInForcePostOnly TimeInForce = "post_only" // Post only order (prevent instant execution)
)

// Margin modes
type MarginMode string
type MarginModeType = MarginMode // Alias for backward compatibility

const (
	MarginModeCrossed  MarginMode = "crossed"
	MarginModeIsolated MarginMode = "isolated"
)

// Additional type aliases for batch orders and advanced features
type PositionSideType string

const (
	PositionSideOpen  PositionSideType = "open"
	PositionSideClose PositionSideType = "close"
)

type ReduceOnlyType string

const (
	ReduceOnlyTrue  ReduceOnlyType = "YES"
	ReduceOnlyFalse ReduceOnlyType = "NO"
)

type SelfTradePreventionType string

const (
	STPNone        SelfTradePreventionType = "none"
	STPCancelTaker SelfTradePreventionType = "cancel_taker"
	STPCancelMaker SelfTradePreventionType = "cancel_maker"
	STPCancelBoth  SelfTradePreventionType = "cancel_both"
)

// Order info models for responses
type OrderInfo struct {
	OrderId       string `json:"orderId"`
	ClientOrderId string `json:"clientOId"`
}

type OrderInfoFailed struct {
	OrderId       string `json:"orderId"`
	ClientOrderId string `json:"clientOId"`
	ErrorMsg      string `json:"errorMsg"`
	ErrorCode     string `json:"errorCode"`
}

// Batch cancel order types
type BatchCancelOrderItem struct {
	OrderId   string `json:"orderId,omitempty"`
	ClientOid string `json:"clientOid,omitempty"`
}

type BatchCancelResponse struct {
	SuccessList []OrderInfo       `json:"successList"`
	FailureList []OrderInfoFailed `json:"failureList"`
}

// API Endpoints for trading operations
const (
	EndpointPlaceOrder        = "/api/v2/mix/order/place-order"
	EndpointModifyOrder       = "/api/v2/mix/order/modify-order"
	EndpointCancelOrder       = "/api/v2/mix/order/cancel-order"
	EndpointCancelAllOrders   = "/api/v2/mix/order/cancel-all-orders"
	EndpointBatchCancelOrders = "/api/v2/mix/order/batch-cancel-orders"
	EndpointOrderDetails      = "/api/v2/mix/order/detail"
	EndpointPendingOrders     = "/api/v2/mix/order/orders-pending"
	EndpointOrderHistory      = "/api/v2/mix/order/history"
	EndpointFillHistory       = "/api/v2/mix/order/fills"
	EndpointBatchOrders       = "/api/v2/mix/order/batch-place-order"
	EndpointCancelPlanOrder   = "/api/v2/mix/order/cancel-plan-order"
	EndpointCreatePlanOrder   = "/api/v2/mix/order/place-plan-order"
	EndpointModifyPlanOrder   = "/api/v2/mix/order/modify-plan-order"
	EndpointPendingPlanOrders = "/api/v2/mix/order/plan-current"
)

// Service Constructor Functions

// NewCreateOrderService creates a new create order service.
func NewCreateOrderService(client ClientInterface) *CreateOrderService {
	return &CreateOrderService{c: client}
}

// NewModifyOrderService creates a new modify order service.
func NewModifyOrderService(client ClientInterface) *ModifyOrderService {
	return &ModifyOrderService{c: client}
}

// NewCancelOrderService creates a new cancel order service.
func NewCancelOrderService(client ClientInterface) *CancelOrderService {
	return &CancelOrderService{c: client}
}

// NewCancelAllOrdersService creates a new cancel all orders service.
func NewCancelAllOrdersService(client ClientInterface) *CancelAllOrdersService {
	return &CancelAllOrdersService{c: client}
}

// NewGetOrderDetailsService creates a new get order details service.
func NewGetOrderDetailsService(client ClientInterface) *GetOrderDetailsService {
	return &GetOrderDetailsService{c: client}
}

// NewPendingOrdersService creates a new pending orders service.
func NewPendingOrdersService(client ClientInterface) *PendingOrdersService {
	return &PendingOrdersService{c: client}
}

// NewOrderHistoryService creates a new order history service.
func NewOrderHistoryService(client ClientInterface) *OrderHistoryService {
	return &OrderHistoryService{c: client}
}

// NewFillHistoryService creates a new fill history service.
func NewFillHistoryService(client ClientInterface) *FillHistoryService {
	return &FillHistoryService{c: client}
}

// NewCreatePlanOrderService creates a new create plan order service.
func NewCreatePlanOrderService(client ClientInterface) *CreatePlanOrderService {
	return &CreatePlanOrderService{c: client}
}

// NewModifyPlanOrderService creates a new modify plan order service.
func NewModifyPlanOrderService(client ClientInterface) *ModifyPlanOrderService {
	return &ModifyPlanOrderService{c: client}
}

// NewCancelPlanOrderService creates a new cancel plan order service.
func NewCancelPlanOrderService(client ClientInterface) *CancelPlanOrderService {
	return &CancelPlanOrderService{c: client}
}

// NewPendingPlanOrdersService creates a new pending plan orders service.
func NewPendingPlanOrdersService(client ClientInterface) *PendingPlanOrdersService {
	return &PendingPlanOrdersService{c: client}
}

// NewCreateBatchOrdersService creates a new create batch orders service.
func NewCreateBatchOrdersService(client ClientInterface) *CreateBatchOrdersService {
	return &CreateBatchOrdersService{c: client}
}

// NewBatchCancelOrdersService creates a new batch cancel orders service.
func NewBatchCancelOrdersService(client ClientInterface) *BatchCancelOrdersService {
	return &BatchCancelOrdersService{c: client}
}
