package trading

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

// BatchCancelOrdersService provides methods to cancel multiple orders in a single request.
type BatchCancelOrdersService struct {
	c           ClientInterface
	symbol      string
	productType ProductType
	marginCoin  string
	orderIdList []BatchCancelOrderItem
}

// Symbol sets the trading pair (required when orderIdList is provided).
func (s *BatchCancelOrdersService) Symbol(symbol string) *BatchCancelOrdersService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type (required).
func (s *BatchCancelOrdersService) ProductType(productType ProductType) *BatchCancelOrdersService {
	s.productType = productType
	return s
}

// MarginCoin sets the margin coin (optional).
func (s *BatchCancelOrdersService) MarginCoin(marginCoin string) *BatchCancelOrdersService {
	s.marginCoin = marginCoin
	return s
}

// OrderIdList sets the list of order IDs to cancel (maximum 50 orders).
func (s *BatchCancelOrdersService) OrderIdList(orderIdList []BatchCancelOrderItem) *BatchCancelOrdersService {
	s.orderIdList = orderIdList
	return s
}

// AddOrderId adds a single order ID to the cancel list.
func (s *BatchCancelOrdersService) AddOrderId(orderId string) *BatchCancelOrdersService {
	if s.orderIdList == nil {
		s.orderIdList = make([]BatchCancelOrderItem, 0)
	}
	s.orderIdList = append(s.orderIdList, BatchCancelOrderItem{
		OrderId: orderId,
	})
	return s
}

// AddClientOid adds a single client order ID to the cancel list.
func (s *BatchCancelOrdersService) AddClientOid(clientOid string) *BatchCancelOrdersService {
	if s.orderIdList == nil {
		s.orderIdList = make([]BatchCancelOrderItem, 0)
	}
	s.orderIdList = append(s.orderIdList, BatchCancelOrderItem{
		ClientOid: clientOid,
	})
	return s
}

// AddOrder adds a single order item to the cancel list.
func (s *BatchCancelOrdersService) AddOrder(orderId, clientOid string) *BatchCancelOrdersService {
	if s.orderIdList == nil {
		s.orderIdList = make([]BatchCancelOrderItem, 0)
	}
	s.orderIdList = append(s.orderIdList, BatchCancelOrderItem{
		OrderId:   orderId,
		ClientOid: clientOid,
	})
	return s
}

// checkRequiredParams validates required parameters.
func (s *BatchCancelOrdersService) checkRequiredParams() error {
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}

	if len(s.orderIdList) > 0 {
		if s.symbol == "" {
			return fmt.Errorf("symbol is required when orderIdList is provided")
		}
		if len(s.orderIdList) > 50 {
			return fmt.Errorf("orderIdList can contain maximum 50 orders")
		}
		// Validate each order item has either orderId or clientOid
		for i, order := range s.orderIdList {
			if order.OrderId == "" && order.ClientOid == "" {
				return fmt.Errorf("order at index %d must have either orderId or clientOid", i)
			}
		}
	}

	return nil
}

// Do sends the batch cancel orders request.
func (s *BatchCancelOrdersService) Do(ctx context.Context) (*BatchCancelResponse, error) {
	if err := s.checkRequiredParams(); err != nil {
		return nil, err
	}

	body := s.batchCancelOrdersRequestBody()
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}

	res, _, err := s.c.CallAPI(ctx, "POST", EndpointBatchCancelOrders, nil, bodyBytes, true)
	if err != nil {
		return nil, err
	}

	var response BatchCancelResponse
	err = jsoniter.Unmarshal(res.Data, &response)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

// batchCancelOrdersRequestBody constructs the request payload.
func (s *BatchCancelOrdersService) batchCancelOrdersRequestBody() map[string]interface{} {
	body := make(map[string]interface{})
	body["productType"] = string(s.productType)

	if s.symbol != "" {
		body["symbol"] = s.symbol
	}

	if s.marginCoin != "" {
		body["marginCoin"] = s.marginCoin
	}

	if len(s.orderIdList) > 0 {
		body["orderIdList"] = s.orderIdList
	}

	return body
}
