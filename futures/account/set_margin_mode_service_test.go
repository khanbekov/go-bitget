package account

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

func TestSetMarginModeService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	// Test fluent API
	result := service.
		Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginMode(futures.MarginModeIsolated)

	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, futures.ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, futures.MarginModeIsolated, result.marginMode)
}

func TestSetMarginModeService_Do_Success_Isolated(t *testing.T) {
	mockMarginModeData := map[string]interface{}{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginMode":  "ISOLATED",
	}

	mockDataBytes, _ := json.Marshal(mockMarginModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginMode(futures.MarginModeIsolated)

	expectedBody := map[string]interface{}{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginMode":  "ISOLATED",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMarginMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "ISOLATED", result.MarginMode)
	mockClient.AssertExpectations(t)
}

func TestSetMarginModeService_Do_Success_Cross(t *testing.T) {
	mockMarginModeData := map[string]interface{}{
		"symbol":      "ETHUSDT",
		"productType": "USDT-FUTURES",
		"marginMode":  "CROSSED",
	}

	mockDataBytes, _ := json.Marshal(mockMarginModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	service.Symbol("ETHUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginMode(futures.MarginModeCrossed)

	expectedBody := map[string]interface{}{
		"symbol":      "ETHUSDT",
		"productType": "USDT-FUTURES",
		"marginMode":  "CROSSED",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMarginMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ETHUSDT", result.Symbol)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "CROSSED", result.MarginMode)
	mockClient.AssertExpectations(t)
}

func TestSetMarginModeService_Do_CoinFutures(t *testing.T) {
	mockMarginModeData := map[string]interface{}{
		"symbol":      "BTCUSD",
		"productType": "COIN-FUTURES",
		"marginMode":  "ISOLATED",
	}

	mockDataBytes, _ := json.Marshal(mockMarginModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	service.Symbol("BTCUSD").
		ProductType(futures.ProductTypeCOINFutures).
		MarginMode(futures.MarginModeIsolated)

	expectedBody := map[string]interface{}{
		"symbol":      "BTCUSD",
		"productType": "COIN-FUTURES",
		"marginMode":  "ISOLATED",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMarginMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "BTCUSD", result.Symbol)
	assert.Equal(t, "COIN-FUTURES", result.ProductType)
	assert.Equal(t, "ISOLATED", result.MarginMode)
	mockClient.AssertExpectations(t)
}

func TestSetMarginModeService_Do_USDCFutures(t *testing.T) {
	mockMarginModeData := map[string]interface{}{
		"symbol":      "ADAUSDT",
		"productType": "USDC-FUTURES",
		"marginMode":  "CROSSED",
	}

	mockDataBytes, _ := json.Marshal(mockMarginModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	service.Symbol("ADAUSDT").
		ProductType(futures.ProductTypeUSDCFutures).
		MarginMode(futures.MarginModeCrossed)

	expectedBody := map[string]interface{}{
		"symbol":      "ADAUSDT",
		"productType": "USDC-FUTURES",
		"marginMode":  "CROSSED",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMarginMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ADAUSDT", result.Symbol)
	assert.Equal(t, "USDC-FUTURES", result.ProductType)
	assert.Equal(t, "CROSSED", result.MarginMode)
	mockClient.AssertExpectations(t)
}

func TestSetMarginModeService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginMode(futures.MarginModeIsolated)

	apiError := fmt.Errorf("API error: margin mode not supported")
	mockClient.On("CallAPI", mock.Anything, "POST", futures.EndpointSetMarginMode, mock.Anything, mock.Anything, true).
		Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "margin mode not supported")
	mockClient.AssertExpectations(t)
}

func TestSetMarginModeService_Do_UnmarshalError(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &SetMarginModeService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginMode(futures.MarginModeIsolated)

	mockClient.On("CallAPI", mock.Anything, "POST", futures.EndpointSetMarginMode, mock.Anything, mock.Anything, true).
		Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestSetMarginModeService_MarginModeConstants(t *testing.T) {
	assert.Equal(t, futures.MarginModeType("ISOLATED"), futures.MarginModeIsolated)
	assert.Equal(t, futures.MarginModeType("CROSSED"), futures.MarginModeCrossed)
}

func TestSetMarginModeService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	service := &SetMarginModeService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginMode(futures.MarginModeIsolated)

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, futures.ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, futures.MarginModeIsolated, service.marginMode)
}
