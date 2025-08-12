package uta

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

func TestModifyOrderService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &ModifyOrderService{c: mockClient}

	// Test fluent API
	result := service.
		Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		OrderId("123456789").
		NewPrice("50000.00").
		NewSize("0.002")

	assert.Equal(t, service, result)
	assert.NotNil(t, service.symbol)
	assert.Equal(t, "BTCUSDT", *service.symbol)
	assert.NotNil(t, service.category)
	assert.Equal(t, CategoryUSDTFutures, *service.category)
	assert.NotNil(t, service.orderId)
	assert.Equal(t, "123456789", *service.orderId)
	assert.NotNil(t, service.newPrice)
	assert.Equal(t, "50000.00", *service.newPrice)
	assert.NotNil(t, service.newSize)
	assert.Equal(t, "0.002", *service.newSize)
}

func TestModifyOrderService_Do_Success(t *testing.T) {
	// Mock order data
	mockOrderData := map[string]interface{}{
		"orderId":     "123456789",
		"symbol":      "BTCUSDT",
		"category":    "USDT-FUTURES",
		"side":        "buy",
		"orderType":   "limit",
		"size":        "0.002",
		"price":       "50000.00",
		"timeInForce": "gtc",
		"status":      "new",
	}

	mockDataBytes, _ := json.Marshal(mockOrderData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := &ModifyOrderService{c: mockClient}

	// Configure service
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		OrderId("123456789").
		NewPrice("50000.00").
		NewSize("0.002")

	// Set up expected parameters
	expectedParams := map[string]interface{}{
		"symbol":   "BTCUSDT",
		"category": "USDT-FUTURES",
		"orderId":  "123456789",
		"price":    "50000.00",
		"qty":      "0.002",
	}
	expectedBody, _ := json.Marshal(expectedParams)

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradeModifyOrder,
		url.Values(nil),
		expectedBody,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123456789", result.OrderID)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, "0.002", result.Size)
	assert.Equal(t, "50000.00", result.Price)
	mockClient.AssertExpectations(t)
}

func TestModifyOrderService_Do_WithClientOid(t *testing.T) {
	// Mock order data
	mockOrderData := map[string]interface{}{
		"orderId":   "123456789",
		"clientOid": "client-123",
		"symbol":    "BTCUSDT",
		"category":  "USDT-FUTURES",
		"size":      "0.003",
		"price":     "49000.00",
		"status":    "new",
	}

	mockDataBytes, _ := json.Marshal(mockOrderData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := &ModifyOrderService{c: mockClient}

	// Configure service with clientOid instead of orderId
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		ClientOid("client-123").
		NewPrice("49000.00").
		NewSize("0.003").
		NewClientOid("new-client-123")

	// Set up expected parameters
	expectedParams := map[string]interface{}{
		"symbol":       "BTCUSDT",
		"category":     "USDT-FUTURES",
		"clientOid":    "client-123",
		"price":        "49000.00",
		"qty":          "0.003",
		"newClientOid": "new-client-123",
	}
	expectedBody, _ := json.Marshal(expectedParams)

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradeModifyOrder,
		url.Values(nil),
		expectedBody,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123456789", result.OrderID)
	mockClient.AssertExpectations(t)
}

func TestModifyOrderService_Do_MissingParameters(t *testing.T) {
	mockClient := &MockClient{}
	service := &ModifyOrderService{c: mockClient}

	// Test missing symbol
	service.Category(CategoryUSDTFutures).OrderId("123456789").NewPrice("50000.00")
	_, err := service.Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symbol")

	// Reset service
	service = &ModifyOrderService{c: mockClient}

	// Test missing category
	service.Symbol("BTCUSDT").OrderId("123456789").NewPrice("50000.00")
	_, err = service.Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "category")

	// Reset service
	service = &ModifyOrderService{c: mockClient}

	// Test missing orderId and clientOid
	service.Symbol("BTCUSDT").Category(CategoryUSDTFutures).NewPrice("50000.00")
	_, err = service.Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "orderId or clientOid")

	// Reset service
	service = &ModifyOrderService{c: mockClient}

	// Test missing modifications
	service.Symbol("BTCUSDT").Category(CategoryUSDTFutures).OrderId("123456789")
	_, err = service.Do(context.Background())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "newSize or newPrice")
}

func TestModifyOrderService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &ModifyOrderService{c: mockClient}

	// Configure service
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		OrderId("123456789").
		NewPrice("50000.00")

	// Mock API error
	apiError := fmt.Errorf("API error: order not found")
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradeModifyOrder,
		url.Values(nil),
		mock.Anything,
		true).Return((*ApiResponse)(nil), &fasthttp.ResponseHeader{}, apiError)

	// Execute
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "order not found")
	mockClient.AssertExpectations(t)
}

func TestModifyOrderService_Do_UnmarshalError(t *testing.T) {
	mockClient := &MockClient{}
	service := &ModifyOrderService{c: mockClient}

	// Configure service
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		OrderId("123456789").
		NewPrice("50000.00")

	// Mock response with invalid JSON
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradeModifyOrder,
		url.Values(nil),
		mock.Anything,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestModifyOrderService_Integration(t *testing.T) {
	// This test verifies the service can be created through the client
	client := NewClient("test", "test", "test")
	service := client.NewModifyOrderService()

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)
}
