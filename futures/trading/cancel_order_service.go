package trading

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// CancelOrderService provides methods to cancel an existing order.
type CancelOrderService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	marginCoin  string
	orderId     string
	clientOid   string
}

// Symbol sets the trading pair (required).
func (s *CancelOrderService) Symbol(symbol string) *CancelOrderService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type (required).
func (s *CancelOrderService) ProductType(productType ProductType) *CancelOrderService {
	s.productType = productType
	return s
}

// MarginCoin sets the margin coin (optional).
func (s *CancelOrderService) MarginCoin(marginCoin string) *CancelOrderService {
	s.marginCoin = marginCoin
	return s
}

// OrderId sets the order ID (either orderId or clientOid required).
func (s *CancelOrderService) OrderId(orderId string) *CancelOrderService {
	s.orderId = orderId
	return s
}

// ClientOid sets the custom order ID (either orderId or clientOid required).
func (s *CancelOrderService) ClientOid(clientOid string) *CancelOrderService {
	s.clientOid = clientOid
	return s
}

// checkRequiredParams validates required parameters.
func (s *CancelOrderService) checkRequiredParams() error {
	if s.symbol == "" {
		return fmt.Errorf("symbol is required")
	}
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.orderId == "" && s.clientOid == "" {
		return fmt.Errorf("either orderId or clientOid must be provided")
	}
	return nil
}

// Do sends the cancel order request.
func (s *CancelOrderService) Do(ctx context.Context) (orderInfo *OrderInfo, err error) {
	if err := s.checkRequiredParams(); err != nil {
		return nil, err
	}
	body := s.cancelOrderRequestBody()
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", EndpointCancelOrder, nil, bodyBytes, true)
	if err != nil {
		return nil, err
	}

	err = jsoniter.Unmarshal(res.Data, &orderInfo)
	if err != nil {
		return nil, err
	}
	return orderInfo, nil
}

// cancelOrderRequestBody constructs the request payload.
func (s *CancelOrderService) cancelOrderRequestBody() map[string]string {
	body := make(map[string]string)
	body["symbol"] = s.symbol
	body["productType"] = string(s.productType)

	if s.marginCoin != "" {
		body["marginCoin"] = s.marginCoin
	}

	if s.orderId != "" {
		body["orderId"] = s.orderId
	}

	if s.clientOid != "" {
		body["clientOid"] = s.clientOid
	}

	return body
}
