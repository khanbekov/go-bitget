package trading

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestCreateOrderService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &CreateOrderService{c: mockClient}

	// Test fluent API pattern
	result := service.
		ProductType(ProductTypeUSDTFutures).
		Symbol("BTCUSDT").
		MarginMode(MarginModeCrossed).
		MarginCoin("USDT").
		SideType(SideTypeBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000")

	assert.Equal(t, ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, MarginModeCrossed, result.marginMode)
	assert.Equal(t, "USDT", result.marginCoin)
	assert.Equal(t, SideTypeBuy, result.sideType)
	assert.Equal(t, OrderTypeLimit, result.orderType)
	assert.Equal(t, "0.001", result.size)
	assert.Equal(t, "50000", result.price)
	assert.Equal(t, service, result, "Should return the same service instance for chaining")
}

func TestCreateOrderService_Do_Success(t *testing.T) {
	// Mock response data
	mockOrderData := map[string]interface{}{
		"orderId":    "123456789",
		"clientOid":  "custom-123",
		"symbol":     "BTCUSDT",
		"size":       "0.001",
		"side":       "buy",
		"orderType":  "limit",
		"price":      "50000",
		"state":      "new",
		"marginCoin": "USDT",
	}

	mockDataBytes, _ := json.Marshal(mockOrderData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &CreateOrderService{c: mockClient}

	// Set up service parameters
	service.ProductType(ProductTypeUSDTFutures).
		Symbol("BTCUSDT").
		MarginMode(MarginModeCrossed).
		MarginCoin("USDT").
		SideType(SideTypeBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointPlaceOrder,
		mock.Anything, // query params (nil for POST)
		mock.Anything, // body bytes
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	orderInfo, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, orderInfo)

	mockClient.AssertExpectations(t)
}

func TestCreateOrderService_Do_MissingRequiredParam(t *testing.T) {
	mockClient := &MockClient{}
	service := &CreateOrderService{c: mockClient}

	// Set up service without required symbol parameter
	service.ProductType(ProductTypeUSDTFutures).
		MarginMode(MarginModeCrossed).
		MarginCoin("USDT").
		SideType(SideTypeBuy).
		OrderType(OrderTypeLimit).
		Size("0.001")
	// Missing symbol - should cause validation error

	ctx := context.Background()
	orderInfo, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, orderInfo)
	assert.Contains(t, err.Error(), "symbol is required")
}

func TestCreateOrderService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &CreateOrderService{c: mockClient}

	service.ProductType(ProductTypeUSDTFutures).
		Symbol("BTCUSDT").
		MarginMode(MarginModeCrossed).
		MarginCoin("USDT").
		SideType(SideTypeBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000")

	// Mock API error
	apiError := fmt.Errorf("API error: insufficient balance")
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointPlaceOrder,
		mock.Anything,
		mock.Anything,
		true).Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	orderInfo, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, orderInfo)
	assert.Contains(t, err.Error(), "insufficient balance")
	mockClient.AssertExpectations(t)
}

func TestCreateOrderService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	service := &CreateOrderService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		ProductType(ProductTypeUSDTFutures).
		Symbol("BTCUSDT").
		MarginMode(MarginModeCrossed).
		MarginCoin("USDT").
		SideType(SideTypeBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000")

	assert.Equal(t, service, result)
	assert.Equal(t, ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, MarginModeCrossed, service.marginMode)
	assert.Equal(t, "USDT", service.marginCoin)
	assert.Equal(t, SideTypeBuy, service.sideType)
	assert.Equal(t, OrderTypeLimit, service.orderType)
}
