package position

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"

	"github.com/khanbekov/go-bitget/futures"
)

// ClosePositionService closes all or part of a position
type ClosePositionService struct {
	c           futures.ClientInterface
	symbol      string
	productType futures.ProductType
	holdSide    futures.HoldSideType
}

func (s *ClosePositionService) Symbol(symbol string) *ClosePositionService {
	s.symbol = symbol
	return s
}

func (s *ClosePositionService) ProductType(productType futures.ProductType) *ClosePositionService {
	s.productType = productType
	return s
}

func (s *ClosePositionService) HoldSide(holdSide futures.HoldSideType) *ClosePositionService {
	s.holdSide = holdSide
	return s
}

func (s *ClosePositionService) checkRequiredParams() error {
	if s.productType == "" {
		return fmt.Errorf("productType is required")
	}
	if s.holdSide == "" {
		return fmt.Errorf("holdSide is required")
	}
	return nil
}

func (s *ClosePositionService) Do(ctx context.Context) (*ClosePositionResponse, error) {
	// check required params before execution
	if err := s.checkRequiredParams(); err != nil {
		return nil, err
	}

	body := s.createClosePositionRequestBody()

	// Marshal body to JSON
	bodyBytes, err := jsoniter.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Make request to API
	var res *futures.ApiResponse

	res, _, err = s.c.CallAPI(ctx, "POST", futures.EndpointClosePosition, nil, bodyBytes, true)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var closePositionResponse *ClosePositionResponse
	err = jsoniter.Unmarshal(res.Data, &closePositionResponse)

	if err != nil {
		return nil, err
	}

	return closePositionResponse, nil
}

func (s *ClosePositionService) createClosePositionRequestBody() map[string]string {
	body := make(map[string]string)

	// Set required params
	body["symbol"] = s.symbol
	body["productType"] = string(s.productType)
	body["holdSide"] = string(s.holdSide)

	return body
}

type ClosePositionResponse struct {
	// Order ID of the closing order
	OrderId string `json:"orderId"`

	// Client order ID if provided
	ClientOrderId string `json:"clientOrderId"`

	// Symbol of the position being closed
	Symbol string `json:"symbol"`

	// Size of the position being closed
	Size string `json:"size"`

	// Side of the closing order (buy/sell)
	Side string `json:"side"`

	// Execution status
	Status string `json:"status"`
}
