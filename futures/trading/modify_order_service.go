package trading

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
)

// ModifyOrderService provides methods to modify an existing order.
type ModifyOrderService struct {
	c                         ClientInterface
	orderId                   string
	clientOrderId             string
	symbol                    string
	productType               ProductType
	marginCoin                string
	newClientOrderId          string
	newSize                   string
	newPrice                  string
	newPresetStopSurplusPrice string
	newPresetStopLossPrice    string
}

// OrderId sets the order ID to modify.
func (s *ModifyOrderService) OrderId(orderId string) *ModifyOrderService {
	s.orderId = orderId
	return s
}

// ClientOrderId sets the custom order ID to identify the order to modify.
func (s *ModifyOrderService) ClientOrderId(clientOrderId string) *ModifyOrderService {
	s.clientOrderId = clientOrderId
	return s
}

// Symbol sets the trading pair. (required)
func (s *ModifyOrderService) Symbol(symbol string) *ModifyOrderService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type (e.g., USDT-FUTURES). (required)
func (s *ModifyOrderService) ProductType(productType ProductType) *ModifyOrderService {
	s.productType = productType
	return s
}

// MarginCoin sets the margin coin for the order. (required)
func (s *ModifyOrderService) MarginCoin(marginCoin string) *ModifyOrderService {
	s.marginCoin = marginCoin
	return s
}

// NewClientOrderId sets the new custom order ID for the modified order. (required)
func (s *ModifyOrderService) NewClientOrderId(newClientOrderId string) *ModifyOrderService {
	s.newClientOrderId = newClientOrderId
	return s
}

// NewSize sets the new order size.
func (s *ModifyOrderService) NewSize(newSize string) *ModifyOrderService {
	s.newSize = newSize
	return s
}

// NewPrice sets the new order price.
func (s *ModifyOrderService) NewPrice(newPrice string) *ModifyOrderService {
	s.newPrice = newPrice
	return s
}

// NewPresetStopSurplusPrice sets the new take-profit price.
func (s *ModifyOrderService) NewPresetStopSurplusPrice(newPresetStopSurplusPrice string) *ModifyOrderService {
	s.newPresetStopSurplusPrice = newPresetStopSurplusPrice
	return s
}

// NewPresetStopLossPrice sets the new stop-loss price.
func (s *ModifyOrderService) NewPresetStopLossPrice(newPresetStopLossPrice string) *ModifyOrderService {
	s.newPresetStopLossPrice = newPresetStopLossPrice
	return s
}

// checkRequiredParams checks if all required parameters are set.
func (s *ModifyOrderService) checkRequiredParams() error {
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.marginCoin == "" {
		return fmt.Errorf("marginCoin is required")
	}
	if s.newClientOrderId == "" {
		return fmt.Errorf("newClientOrderId is required")
	}
	return nil
}

// Do sends the request to modify the order.
func (s *ModifyOrderService) Do(ctx context.Context) (modifyOrderResponse *OrderInfo, err error) {
	// Check required parameters before execution
	if err = s.checkRequiredParams(); err != nil {
		return nil, err
	}
	body := s.modifyOrderRequestBody()
	// Marshal body to JSON
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}
	// Make request to API
	var res *ApiResponse
	res, _, err = s.c.CallAPI(ctx, "POST", EndpointModifyOrder, nil, bodyBytes, true)
	if err != nil {
		return nil, err
	}
	// Unmarshal JSON from response
	err = jsoniter.Unmarshal(res.Data, &modifyOrderResponse)
	if err != nil {
		return nil, err
	}
	return modifyOrderResponse, nil
}

// modifyOrderRequestBody creates the request body for modifying an order.
func (s *ModifyOrderService) modifyOrderRequestBody() map[string]string {
	body := make(map[string]string)
	// Set required parameters
	body["symbol"] = s.symbol
	body["productType"] = string(s.productType)
	body["marginCoin"] = s.marginCoin
	body["newClientOid"] = s.newClientOrderId
	// Set optional parameters
	if s.orderId != "" {
		body["orderId"] = s.orderId
	}
	if s.clientOrderId != "" {
		body["clientOid"] = s.clientOrderId
	}
	if s.newSize != "" {
		body["newSize"] = s.newSize
	}
	if s.newPrice != "" {
		body["newPrice"] = s.newPrice
	}
	if s.newPresetStopSurplusPrice != "" {
		body["newPresetStopSurplusPrice"] = s.newPresetStopSurplusPrice
	}
	if s.newPresetStopLossPrice != "" {
		body["newPresetStopLossPrice"] = s.newPresetStopLossPrice
	}
	return body
}
