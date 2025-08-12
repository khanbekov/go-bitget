package trading

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
)

// BatchOrderInfo represents a single order in a batch
type BatchOrderInfo struct {
	Size                          string                  `json:"size"`
	Price                         string                  `json:"price,omitempty"`
	SideType                      SideType                `json:"side"`
	PositionSideType              PositionSideType        `json:"tradeSide,omitempty"`
	OrderType                     OrderType               `json:"orderType"`
	TimeInForceType               TimeInForceType         `json:"force,omitempty"`
	ClientOrderId                 string                  `json:"clientOid,omitempty"`
	ReduceOnlyType                ReduceOnlyType          `json:"reduceOnly,omitempty"`
	PresetStopSurplusPrice        string                  `json:"presetStopSurplusPrice,omitempty"`
	PresetStopLossPrice           string                  `json:"presetStopLossPrice,omitempty"`
	PresetStopSurplusExecutePrice string                  `json:"presetStopSurplusExecutePrice,omitempty"`
	PresetStopLossExecutePrice    string                  `json:"presetStopLossExecutePrice,omitempty"`
	SelfTradePreventionType       SelfTradePreventionType `json:"stpMode,omitempty"`
}

// CreateBatchOrdersResponse represents the response from batch order creation
type CreateBatchOrdersResponse struct {
	SuccessList []OrderInfo       `json:"successList"`
	FailureList []OrderInfoFailed `json:"failureList"`
}

// CreateBatchOrdersService provides methods to create multiple orders in a single request
type CreateBatchOrdersService struct {
	c           ClientInterface
	productType ProductType
	symbol      string
	marginMode  MarginModeType
	marginCoin  string
	orders      []BatchOrderInfo
}

// ProductType sets type of market on bitget (USDT-FUTURES, COIN-FUTURES etc.) REQUIRED
func (s *CreateBatchOrdersService) ProductType(productType ProductType) *CreateBatchOrdersService {
	s.productType = productType
	return s
}

// Symbol sets the trading pair. REQUIRED
func (s *CreateBatchOrdersService) Symbol(symbol string) *CreateBatchOrdersService {
	s.symbol = symbol
	return s
}

// MarginMode sets the margin mode for position (crossed/isolated). REQUIRED
func (s *CreateBatchOrdersService) MarginMode(marginMode MarginModeType) *CreateBatchOrdersService {
	s.marginMode = marginMode
	return s
}

// MarginCoin sets the margin coin for the order. REQUIRED
func (s *CreateBatchOrdersService) MarginCoin(marginCoin string) *CreateBatchOrdersService {
	s.marginCoin = marginCoin
	return s
}

// AddOrder adds a single order to the batch. Multiple calls build up the batch.
func (s *CreateBatchOrdersService) AddOrder(order BatchOrderInfo) *CreateBatchOrdersService {
	s.orders = append(s.orders, order)
	return s
}

// Orders sets all orders in the batch at once, replacing any previously added orders.
func (s *CreateBatchOrdersService) Orders(orders []BatchOrderInfo) *CreateBatchOrdersService {
	s.orders = orders
	return s
}

// ClearOrders removes all orders from the batch
func (s *CreateBatchOrdersService) ClearOrders() *CreateBatchOrdersService {
	s.orders = nil
	return s
}

// GetOrderCount returns the number of orders currently in the batch
func (s *CreateBatchOrdersService) GetOrderCount() int {
	return len(s.orders)
}

// validateOrder validates a single order in the batch
func (s *CreateBatchOrdersService) validateOrder(order BatchOrderInfo) error {
	if order.Size == "" {
		return fmt.Errorf("size is required for all orders")
	}
	if order.SideType == "" {
		return fmt.Errorf("sideType is required for all orders")
	}
	if order.OrderType == "" {
		return fmt.Errorf("orderType is required for all orders")
	}
	return nil
}

// checkRequiredParams validates required parameters for the batch order request
func (s *CreateBatchOrdersService) checkRequiredParams() error {
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if s.marginCoin == "" {
		return fmt.Errorf("marginCoin is required")
	}
	if s.marginMode == "" {
		return fmt.Errorf("marginMode is required")
	}
	if len(s.orders) == 0 {
		return fmt.Errorf("at least one order is required")
	}
	if len(s.orders) > 20 {
		return fmt.Errorf("maximum 20 orders allowed per batch")
	}

	// Validate each order in the batch
	for i, order := range s.orders {
		if err := s.validateOrder(order); err != nil {
			return fmt.Errorf("order %d: %w", i+1, err)
		}
	}

	return nil
}

// Do sends the batch order creation request
func (s *CreateBatchOrdersService) Do(ctx context.Context) (batchResponse *CreateBatchOrdersResponse, err error) {
	// Check required params before execution
	if err = s.checkRequiredParams(); err != nil {
		return nil, err
	}

	body := s.createBatchOrderRequestBody()

	// Marshal body to JSON
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Make request to API
	var res *ApiResponse

	res, _, err = s.c.CallAPI(ctx, "POST", EndpointBatchOrders, nil, bodyBytes, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	err = jsoniter.Unmarshal(res.Data, &batchResponse)

	if err != nil {
		return nil, err
	}

	return batchResponse, nil
}

// createBatchOrderRequestBody constructs the request payload for batch order creation
func (s *CreateBatchOrdersService) createBatchOrderRequestBody() map[string]interface{} {
	body := make(map[string]interface{})

	// Set product type
	body["productType"] = string(s.productType)
	body["symbol"] = s.symbol
	body["marginMode"] = string(s.marginMode)
	body["marginCoin"] = s.marginCoin

	// Convert orders to the format expected by the API
	orderList := make([]map[string]interface{}, len(s.orders))
	for i, order := range s.orders {
		orderMap := make(map[string]interface{})

		// Required fields
		orderMap["size"] = order.Size
		orderMap["side"] = string(order.SideType)
		orderMap["orderType"] = string(order.OrderType)

		// Optional fields
		if order.Price != "" {
			orderMap["price"] = order.Price
		}
		if order.PositionSideType != "" {
			orderMap["tradeSide"] = string(order.PositionSideType)
		}
		if order.TimeInForceType != "" {
			orderMap["force"] = string(order.TimeInForceType)
		}
		if order.ClientOrderId != "" {
			orderMap["clientOid"] = order.ClientOrderId
		}
		if order.ReduceOnlyType != "" {
			orderMap["reduceOnly"] = string(order.ReduceOnlyType)
		}
		if order.PresetStopSurplusPrice != "" {
			orderMap["presetStopSurplusPrice"] = order.PresetStopSurplusPrice
		}
		if order.PresetStopLossPrice != "" {
			orderMap["presetStopLossPrice"] = order.PresetStopLossPrice
		}
		if order.PresetStopSurplusExecutePrice != "" {
			orderMap["presetStopSurplusExecutePrice"] = order.PresetStopSurplusExecutePrice
		}
		if order.PresetStopLossExecutePrice != "" {
			orderMap["presetStopLossExecutePrice"] = order.PresetStopLossExecutePrice
		}
		if order.SelfTradePreventionType != "" {
			orderMap["stpMode"] = string(order.SelfTradePreventionType)
		}

		orderList[i] = orderMap
	}

	body["orderList"] = orderList

	return body
}
