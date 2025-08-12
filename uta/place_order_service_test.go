package uta

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/khanbekov/go-bitget/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestPlaceOrderService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &PlaceOrderService{c: mockClient}

	// Test fluent API chaining
	result := service.
		Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		Side(SideBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000").
		ClientOid("test_client_oid").
		TimeInForce(TimeInForceGTC).
		ReduceOnly("false").
		PositionSide(PositionSideLong).
		STP(STPNone)

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", *service.symbol)
	assert.Equal(t, CategoryUSDTFutures, *service.category)
	assert.Equal(t, SideBuy, *service.side)
	assert.Equal(t, OrderTypeLimit, *service.orderType)
	assert.Equal(t, "0.001", *service.size)
	assert.Equal(t, "50000", *service.price)
	assert.Equal(t, "test_client_oid", *service.clientOid)
	assert.Equal(t, TimeInForceGTC, *service.timeInForce)
	assert.Equal(t, "false", *service.reduceOnly)
	assert.Equal(t, PositionSideLong, *service.positionSide)
	assert.Equal(t, STPNone, *service.stp)
}

func TestPlaceOrderService_Do_Success(t *testing.T) {
	// Setup mock data
	mockOrder := Order{
		OrderID:      "123456789",
		ClientOid:    "test_client_oid",
		Symbol:       "BTCUSDT",
		Category:     CategoryUSDTFutures,
		Side:         SideBuy,
		OrderType:    OrderTypeLimit,
		Price:        "50000",
		Size:         "0.001",
		FilledSize:   "0",
		FilledAmount: "0",
		AvgPrice:     "0",
		Status:       OrderStatusLive,
		TimeInForce:  TimeInForceGTC,
		CreatedTime:  "1640995200000",
		UpdatedTime:  "1640995200000",
	}

	mockDataBytes, _ := json.Marshal(mockOrder)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client and service
	mockClient := &MockClient{}
	service := &PlaceOrderService{c: mockClient}

	// Configure service parameters
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		Side(SideBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000").
		ClientOid("test_client_oid")

	// Create expected request body
	expectedBody := map[string]interface{}{
		"symbol":    "BTCUSDT",
		"category":  CategoryUSDTFutures,
		"side":      SideBuy,
		"orderType": OrderTypeLimit,
		"size":      "0.001",
		"price":     "50000",
		"clientOid": "test_client_oid",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	// Set up expected API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradePlaceOrder,
		(*url.Values)(nil),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute test
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "123456789", result.OrderID)
	assert.Equal(t, "test_client_oid", result.ClientOid)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, CategoryUSDTFutures, result.Category)
	assert.Equal(t, SideBuy, result.Side)
	assert.Equal(t, OrderTypeLimit, result.OrderType)
	assert.Equal(t, "50000", result.Price)
	assert.Equal(t, "0.001", result.Size)
	assert.Equal(t, OrderStatusLive, result.Status)

	mockClient.AssertExpectations(t)
}

func TestPlaceOrderService_Do_MissingRequiredParameters(t *testing.T) {
	mockClient := &MockClient{}
	service := &PlaceOrderService{c: mockClient}
	ctx := context.Background()

	// Test missing symbol
	_, err := service.Do(ctx)
	assert.Error(t, err)
	assert.IsType(t, &common.MissingParameterError{}, err)

	// Test missing category
	service.Symbol("BTCUSDT")
	_, err = service.Do(ctx)
	assert.Error(t, err)
	assert.IsType(t, &common.MissingParameterError{}, err)

	// Test missing side
	service.Category(CategoryUSDTFutures)
	_, err = service.Do(ctx)
	assert.Error(t, err)
	assert.IsType(t, &common.MissingParameterError{}, err)

	// Test missing orderType
	service.Side(SideBuy)
	_, err = service.Do(ctx)
	assert.Error(t, err)
	assert.IsType(t, &common.MissingParameterError{}, err)

	// Test missing size
	service.OrderType(OrderTypeLimit)
	_, err = service.Do(ctx)
	assert.Error(t, err)
	assert.IsType(t, &common.MissingParameterError{}, err)
}

func TestPlaceOrderService_Do_WithOptionalParameters(t *testing.T) {
	// Setup mock data
	mockOrder := Order{
		OrderID:      "123456789",
		ClientOid:    "test_client_oid",
		Symbol:       "BTCUSDT",
		Category:     CategoryUSDTFutures,
		Side:         SideBuy,
		OrderType:    OrderTypeLimit,
		Price:        "50000",
		Size:         "0.001",
		Status:       OrderStatusLive,
		TimeInForce:  TimeInForceGTC,
		ReduceOnly:   "false",
		PositionSide: PositionSideLong,
		STP:          STPNone,
	}

	mockDataBytes, _ := json.Marshal(mockOrder)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client and service
	mockClient := &MockClient{}
	service := &PlaceOrderService{c: mockClient}

	// Configure service with all parameters
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		Side(SideBuy).
		OrderType(OrderTypeLimit).
		Size("0.001").
		Price("50000").
		ClientOid("test_client_oid").
		TimeInForce(TimeInForceGTC).
		ReduceOnly("false").
		PositionSide(PositionSideLong).
		STP(STPNone)

	// Create expected request body with all parameters
	expectedBody := map[string]interface{}{
		"symbol":       "BTCUSDT",
		"category":     CategoryUSDTFutures,
		"side":         SideBuy,
		"orderType":    OrderTypeLimit,
		"size":         "0.001",
		"price":        "50000",
		"clientOid":    "test_client_oid",
		"timeInForce":  TimeInForceGTC,
		"reduceOnly":   "false",
		"positionSide": PositionSideLong,
		"stp":          STPNone,
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	// Set up expected API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradePlaceOrder,
		(*url.Values)(nil),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute test
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, TimeInForceGTC, result.TimeInForce)
	assert.Equal(t, "false", result.ReduceOnly)
	assert.Equal(t, PositionSideLong, result.PositionSide)
	assert.Equal(t, STPNone, result.STP)

	mockClient.AssertExpectations(t)
}

func TestPlaceOrderService_Do_MarketOrder(t *testing.T) {
	// Test market order without price
	mockOrder := Order{
		OrderID:   "123456789",
		Symbol:    "BTCUSDT",
		Category:  CategoryUSDTFutures,
		Side:      SideBuy,
		OrderType: OrderTypeMarket,
		Size:      "0.001",
		Status:    OrderStatusLive,
	}

	mockDataBytes, _ := json.Marshal(mockOrder)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := &PlaceOrderService{c: mockClient}

	// Configure service for market order (no price)
	service.Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		Side(SideBuy).
		OrderType(OrderTypeMarket).
		Size("0.001")

	expectedBody := map[string]interface{}{
		"symbol":    "BTCUSDT",
		"category":  CategoryUSDTFutures,
		"side":      SideBuy,
		"orderType": OrderTypeMarket,
		"size":      "0.001",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		EndpointTradePlaceOrder,
		(*url.Values)(nil),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, OrderTypeMarket, result.OrderType)
	assert.Empty(t, result.Price) // Market orders don't have price

	mockClient.AssertExpectations(t)
}

func TestPlaceOrderService_Integration(t *testing.T) {
	// Integration-style test using the real client structure
	client := NewClient("test_api_key", "test_secret_key", "test_passphrase")
	service := client.NewPlaceOrderService()

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test fluent API works with real service
	service.Symbol("BTCUSDT").Category(CategoryUSDTFutures).Side(SideBuy)
	assert.Equal(t, "BTCUSDT", *service.symbol)
	assert.Equal(t, CategoryUSDTFutures, *service.category)
	assert.Equal(t, SideBuy, *service.side)
}
