package trading

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestBatchCancelOrdersService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	// Test fluent API chaining
	result := service.
		Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT").
		AddOrderId("123456789").
		AddClientOid("client123").
		AddOrder("987654321", "client456")

	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "USDT", service.marginCoin)
	assert.Len(t, service.orderIdList, 3)
	assert.Equal(t, "123456789", service.orderIdList[0].OrderId)
	assert.Equal(t, "client123", service.orderIdList[1].ClientOid)
	assert.Equal(t, "987654321", service.orderIdList[2].OrderId)
	assert.Equal(t, "client456", service.orderIdList[2].ClientOid)
	assert.Same(t, service, result)
}

func TestBatchCancelOrdersService_Do_Success(t *testing.T) {
	// Mock response data
	mockResponseData := BatchCancelResponse{
		SuccessList: []OrderInfo{
			{OrderId: "123456789", ClientOrderId: "client123"},
			{OrderId: "987654321", ClientOrderId: "client456"},
		},
		FailureList: []OrderInfoFailed{},
	}

	mockDataBytes, _ := json.Marshal(mockResponseData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client and service
	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	// Configure service
	orderList := []BatchCancelOrderItem{
		{OrderId: "123456789"},
		{ClientOid: "client123"},
	}
	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT").
		OrderIdList(orderList)

	// Mock API call expectations with JSON matching
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointBatchCancelOrders,
		url.Values(nil),
		mock.MatchedBy(func(body []byte) bool {
			var actualBody map[string]interface{}
			if err := json.Unmarshal(body, &actualBody); err != nil {
				return false
			}
			return actualBody["symbol"] == "BTCUSDT" &&
				actualBody["productType"] == "USDT-FUTURES" &&
				actualBody["marginCoin"] == "USDT" &&
				len(actualBody["orderIdList"].([]interface{})) == 2
		}),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute test
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.SuccessList, 2)
	assert.Len(t, result.FailureList, 0)
	assert.Equal(t, "123456789", result.SuccessList[0].OrderId)
	assert.Equal(t, "client123", result.SuccessList[0].ClientOrderId)
	mockClient.AssertExpectations(t)
}

func TestBatchCancelOrdersService_Do_WithFailures(t *testing.T) {
	// Mock response data with some failures
	mockResponseData := BatchCancelResponse{
		SuccessList: []OrderInfo{
			{OrderId: "123456789", ClientOrderId: "client123"},
		},
		FailureList: []OrderInfoFailed{
			{OrderId: "invalid123", ClientOrderId: "invalid_client", ErrorMsg: "Order not found", ErrorCode: "40001"},
		},
	}

	mockDataBytes, _ := json.Marshal(mockResponseData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		AddOrderId("123456789").
		AddOrderId("invalid123")

	mockClient.On("CallAPI", mock.Anything, "POST", EndpointBatchCancelOrders, mock.Anything, mock.Anything, true).
		Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.SuccessList, 1)
	assert.Len(t, result.FailureList, 1)
	assert.Equal(t, "Order not found", result.FailureList[0].ErrorMsg)
	assert.Equal(t, "40001", result.FailureList[0].ErrorCode)
	mockClient.AssertExpectations(t)
}

func TestBatchCancelOrdersService_Do_MinimalParams(t *testing.T) {
	// Test with only productType (no orderIdList)
	mockResponseData := BatchCancelResponse{
		SuccessList: []OrderInfo{},
		FailureList: []OrderInfoFailed{},
	}

	mockDataBytes, _ := json.Marshal(mockResponseData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	service.ProductType(ProductTypeUSDTFutures)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointBatchCancelOrders,
		url.Values(nil),
		mock.MatchedBy(func(body []byte) bool {
			var actualBody map[string]interface{}
			if err := json.Unmarshal(body, &actualBody); err != nil {
				return false
			}
			return actualBody["productType"] == "USDT-FUTURES"
		}),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	mockClient.AssertExpectations(t)
}

func TestBatchCancelOrdersService_ValidationErrors(t *testing.T) {
	mockClient := &MockClient{}

	tests := []struct {
		name          string
		setupService  func() *BatchCancelOrdersService
		expectedError string
	}{
		{
			name: "Missing ProductType",
			setupService: func() *BatchCancelOrdersService {
				return NewBatchCancelOrdersService(mockClient).
					Symbol("BTCUSDT")
			},
			expectedError: "productType is required",
		},
		{
			name: "Missing Symbol with OrderIdList",
			setupService: func() *BatchCancelOrdersService {
				return NewBatchCancelOrdersService(mockClient).
					ProductType(ProductTypeUSDTFutures).
					AddOrderId("123456")
			},
			expectedError: "symbol is required when orderIdList is provided",
		},
		{
			name: "Too Many Orders",
			setupService: func() *BatchCancelOrdersService {
				service := NewBatchCancelOrdersService(mockClient).
					Symbol("BTCUSDT").
					ProductType(ProductTypeUSDTFutures)

				// Add 51 orders (exceeds limit of 50)
				for i := 0; i < 51; i++ {
					service.AddOrderId(fmt.Sprintf("order_%d", i))
				}
				return service
			},
			expectedError: "orderIdList can contain maximum 50 orders",
		},
		{
			name: "Empty Order Item",
			setupService: func() *BatchCancelOrdersService {
				orderList := []BatchCancelOrderItem{
					{OrderId: "123456"},
					{}, // Empty order item
				}
				return NewBatchCancelOrdersService(mockClient).
					Symbol("BTCUSDT").
					ProductType(ProductTypeUSDTFutures).
					OrderIdList(orderList)
			},
			expectedError: "order at index 1 must have either orderId or clientOid",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := test.setupService()
			ctx := context.Background()
			result, err := service.Do(ctx)

			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), test.expectedError)
		})
	}
}

func TestBatchCancelOrdersService_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	service.ProductType(ProductTypeUSDTFutures)

	// Mock API error
	apiError := fmt.Errorf("API error: rate limit exceeded")
	mockClient.On("CallAPI", mock.Anything, "POST", EndpointBatchCancelOrders, mock.Anything, mock.Anything, true).
		Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	mockClient.AssertExpectations(t)
}

func TestBatchCancelOrdersService_UnmarshalError(t *testing.T) {
	// Mock response with invalid JSON data
	mockResponse := &ApiResponse{
		Code: "00000",
		Msg:  "success",
		Data: json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	service.ProductType(ProductTypeUSDTFutures)

	mockClient.On("CallAPI", mock.Anything, "POST", EndpointBatchCancelOrders, mock.Anything, mock.Anything, true).
		Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestBatchCancelOrdersService_OrderIdListMethods(t *testing.T) {
	mockClient := &MockClient{}
	service := NewBatchCancelOrdersService(mockClient)

	// Test individual add methods
	service.AddOrderId("order1")
	assert.Len(t, service.orderIdList, 1)
	assert.Equal(t, "order1", service.orderIdList[0].OrderId)
	assert.Empty(t, service.orderIdList[0].ClientOid)

	service.AddClientOid("client1")
	assert.Len(t, service.orderIdList, 2)
	assert.Empty(t, service.orderIdList[1].OrderId)
	assert.Equal(t, "client1", service.orderIdList[1].ClientOid)

	service.AddOrder("order2", "client2")
	assert.Len(t, service.orderIdList, 3)
	assert.Equal(t, "order2", service.orderIdList[2].OrderId)
	assert.Equal(t, "client2", service.orderIdList[2].ClientOid)

	// Test OrderIdList method (should replace existing list)
	newList := []BatchCancelOrderItem{
		{OrderId: "new_order"},
	}
	service.OrderIdList(newList)
	assert.Len(t, service.orderIdList, 1)
	assert.Equal(t, "new_order", service.orderIdList[0].OrderId)
}

func TestBatchCancelOrdersService_ConstructorIntegration(t *testing.T) {
	mockClient := &MockClient{}

	// Test that the constructor creates a service that works with the fluent API
	service := NewBatchCancelOrdersService(mockClient)

	assert.NotNil(t, service)
	assert.Equal(t, mockClient, service.c)
	assert.Empty(t, service.symbol)
	assert.Empty(t, service.productType)
	assert.Empty(t, service.marginCoin)
	assert.Nil(t, service.orderIdList)
}
