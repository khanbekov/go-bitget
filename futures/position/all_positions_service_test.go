package position

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"testing"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestAllPositionsService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &AllPositionsService{c: mockClient}

	// Test fluent API pattern
	result := service.
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT")

	assert.Equal(t, futures.ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, "USDT", result.marginCoin)
	assert.Equal(t, service, result, "Should return the same service instance for chaining")
}

func TestAllPositionsService_Do_Success(t *testing.T) {
	// Mock response data
	mockPositionData := []map[string]interface{}{
		{
			"marginCoin":             "USDT",
			"symbol":                 "BTCUSDT",
			"holdSide":               "long",
			"size":                   "0.001",
			"markPrice":              "50000.00",
			"positionValue":          "50.00",
			"averageOpenPrice":       "49000.00",
			"unrealizedPL":           "1.00",
			"unrealizedPLR":          "0.02",
			"margin":                 "25.00",
			"available":              "0.001",
			"crossedLeverage":        "2.0",
			"isolatedLeverage":       "0.0",
			"marginMode":             "cross",
			"posMode":                "one_way_mode",
			"marginRatio":            "0.05",
			"maintenanceMarginRatio": "0.005",
			"ctime":                  "1640995200000",
			"utime":                  "1640995200000",
			"breakEvenPrice":         "49500.00",
			"totalFee":               "0.05",
			"deductedFee":            "0.02",
			"autoMargin":             "on",
			"assetMode":              "single_asset",
		},
	}

	mockDataBytes, _ := json.Marshal(mockPositionData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &AllPositionsService{c: mockClient}

	// Set up service parameters
	service.ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT")

	// Expected query parameters
	expectedParams := url.Values{}
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("marginCoin", "USDT")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAllPositions,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	positions, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, positions, 1)

	// Check position data
	pos := positions[0]
	assert.Equal(t, "USDT", pos.MarginCoin)
	assert.Equal(t, "BTCUSDT", pos.Symbol)
	assert.Equal(t, futures.HoldSideLong, pos.HoldSide)
	assert.Equal(t, 0.001, pos.Size)
	assert.Equal(t, 50000.0, pos.MarkPrice)

	mockClient.AssertExpectations(t)
}

func TestAllPositionsService_Do_WithoutMarginCoin(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`[]`),
	}

	mockClient := &MockClient{}
	service := &AllPositionsService{c: mockClient}

	// Set up service without margin coin
	service.ProductType(futures.ProductTypeUSDTFutures)

	// Expected query parameters without margin coin
	expectedParams := url.Values{}
	expectedParams.Set("productType", "USDT-FUTURES")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAllPositions,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	positions, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.Empty(t, positions)
	mockClient.AssertExpectations(t)
}

func TestAllPositionsService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &AllPositionsService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures)

	// Mock API error
	apiError := fmt.Errorf("API error: insufficient permissions")
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAllPositions,
		mock.Anything,
		[]byte(nil),
		true).Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	positions, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, positions)
	assert.Contains(t, err.Error(), "insufficient permissions")
	mockClient.AssertExpectations(t)
}

func TestAllPositionsService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	service := &AllPositionsService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT")

	assert.Equal(t, service, result)
	assert.Equal(t, futures.ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "USDT", service.marginCoin)
}
