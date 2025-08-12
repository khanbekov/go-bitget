package trading

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
)

// CreateOrderService account info
type CreateOrderService struct {
	c                             ClientInterface
	productType                   ProductType
	symbol                        string
	marginMode                    MarginModeType
	marginCoin                    string
	size                          string
	price                         string
	sideType                      SideType
	positionSideType              PositionSideType
	orderType                     OrderType
	timeInForceType               TimeInForceType
	clientOrderId                 string
	reduceOnlyType                ReduceOnlyType
	presetStopSurplusPrice        string
	presetStopLossPrice           string
	presetStopSurplusExecutePrice string
	presetStopLossExecutePrice    string
	selfTradePreventionType       SelfTradePreventionType
}

// ProductType sets type of market on bitget (USDT-FUTURES, COIN-FUTURES etc.) REQUIRED
func (s *CreateOrderService) ProductType(productType ProductType) *CreateOrderService {
	s.productType = productType
	return s
}

// Symbol sets the trading pair. REQUIRED
func (s *CreateOrderService) Symbol(symbol string) *CreateOrderService {
	s.symbol = symbol
	return s
}

// MarginMode sets the margin mode for position (crossed/isolated). REQUIRED
func (s *CreateOrderService) MarginMode(marginMode MarginModeType) *CreateOrderService {
	s.marginMode = marginMode
	return s
}

// MarginCoin sets the margin coin for the order. REQUIRED
func (s *CreateOrderService) MarginCoin(marginCoin string) *CreateOrderService {
	s.marginCoin = marginCoin
	return s
}

// SideType sets the order side (buy/sell). REQUIRED
func (s *CreateOrderService) SideType(sideType SideType) *CreateOrderService {
	s.sideType = sideType
	return s
}

// OrderType sets the order type (limit/market/etc). REQUIRED
func (s *CreateOrderService) OrderType(orderType OrderType) *CreateOrderService {
	s.orderType = orderType
	return s
}

// Size sets the order size (quantity). REQUIRED
func (s *CreateOrderService) Size(size string) *CreateOrderService {
	s.size = size
	return s
}

// Price sets the order price
func (s *CreateOrderService) Price(price string) *CreateOrderService {
	s.price = price
	return s
}

// PositionSideType sets the position side (open/close). Only required in hedge-mode
func (s *CreateOrderService) PositionSideType(positionSideType PositionSideType) *CreateOrderService {
	s.positionSideType = positionSideType
	return s
}

// TimeInForceType sets the time in force (GTC/IOC/etc)
func (s *CreateOrderService) TimeInForceType(timeInForceType TimeInForceType) *CreateOrderService {
	s.timeInForceType = timeInForceType
	return s
}

// ClientOrderId sets custom order id
func (s *CreateOrderService) ClientOrderId(clientOrderId string) *CreateOrderService {
	s.clientOrderId = clientOrderId
	return s
}

// ReduceOnlyType sets whether the order is reduce-only
func (s *CreateOrderService) ReduceOnlyType(reduceOnlyType ReduceOnlyType) *CreateOrderService {
	s.reduceOnlyType = reduceOnlyType
	return s
}

// PresetStopSurplusPrice sets the preset stop surplus price
func (s *CreateOrderService) PresetStopSurplusPrice(presetStopSurplusPrice string) *CreateOrderService {
	s.presetStopSurplusPrice = presetStopSurplusPrice
	return s
}

// PresetStopLossPrice sets the preset stop loss price
func (s *CreateOrderService) PresetStopLossPrice(presetStopLossPrice string) *CreateOrderService {
	s.presetStopLossPrice = presetStopLossPrice
	return s
}

// PresetStopSurplusExecutePrice sets the surplus stop execute price
func (s *CreateOrderService) PresetStopSurplusExecutePrice(presetStopSurplusExecutePrice string) *CreateOrderService {
	s.presetStopSurplusExecutePrice = presetStopSurplusExecutePrice
	return s
}

// PresetStopLossExecutePrice sets the stop loss execute price
func (s *CreateOrderService) PresetStopLossExecutePrice(presetStopLossExecutePrice string) *CreateOrderService {
	s.presetStopLossExecutePrice = presetStopLossExecutePrice
	return s
}

// StpMode sets the self-trade prevention mode
func (s *CreateOrderService) StpMode(stpMode SelfTradePreventionType) *CreateOrderService {
	s.selfTradePreventionType = stpMode
	return s
}

func (s *CreateOrderService) checkRequiredParams() error {
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if s.marginMode == "" {
		return fmt.Errorf("marginMode is required")
	}
	if s.marginCoin == "" {
		return fmt.Errorf("marginCoin is required")
	}
	if s.size == "" {
		return fmt.Errorf("size is required")
	}
	if s.sideType == "" {
		return fmt.Errorf("sideType is required")
	}
	if s.orderType == "" {
		return fmt.Errorf("orderType is required")
	}
	return nil
}

func (s *CreateOrderService) Do(ctx context.Context) (createOrderResponse *OrderInfo, err error) {
	// check required params before execution
	if err = s.checkRequiredParams(); err != nil {
		return nil, err
	}

	body := s.createOrderRequrestBody()

	// Marshal body to JSON
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Make request to API
	var res *ApiResponse

	res, _, err = s.c.CallAPI(ctx, "POST", EndpointPlaceOrder, nil, bodyBytes, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	err = jsoniter.Unmarshal(res.Data, &createOrderResponse)

	if err != nil {
		return nil, err
	}

	return createOrderResponse, nil
}

func (s *CreateOrderService) createOrderRequrestBody() map[string]string {
	body := make(map[string]string)

	// Set params of request

	// required params
	body["productType"] = string(s.productType)
	body["symbol"] = s.symbol
	body["marginCoin"] = s.marginCoin
	body["marginMode"] = string(s.marginMode)
	body["size"] = s.size
	body["side"] = string(s.sideType)
	body["orderType"] = string(s.orderType)

	// optional params
	if s.price != "" {
		body["price"] = s.price
	}
	if s.positionSideType != "" {
		body["tradeSide"] = string(s.positionSideType)
	}
	if s.timeInForceType != "" {
		body["force"] = string(s.timeInForceType)
	}
	if s.clientOrderId != "" {
		body["clientOid"] = s.clientOrderId
	}
	if s.reduceOnlyType != "" {
		body["reduceOnly"] = string(s.reduceOnlyType)
	}
	if s.presetStopSurplusPrice != "" {
		body["presetStopSurplusPrice"] = s.presetStopSurplusPrice
	}
	if s.presetStopLossPrice != "" {
		body["presetStopLossPrice"] = s.presetStopLossPrice
	}
	if s.presetStopSurplusExecutePrice != "" {
		body["presetStopSurplusExecutePrice"] = s.presetStopSurplusExecutePrice
	}
	if s.presetStopLossExecutePrice != "" {
		body["presetStopLossExecutePrice"] = s.presetStopLossExecutePrice
	}
	if s.selfTradePreventionType != "" {
		body["stpMode"] = string(s.selfTradePreventionType)
	}
	return body
}
