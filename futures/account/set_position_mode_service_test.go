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

func TestSetPositionModeService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	// Test fluent API
	result := service.
		ProductType(futures.ProductTypeUSDTFutures).
		PositionMode(futures.PositionModeOneWay)

	assert.Equal(t, futures.ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, futures.PositionModeOneWay, result.positionMode)
}

func TestSetPositionModeService_Do_Success_OneWay(t *testing.T) {
	mockPositionModeData := map[string]interface{}{
		"productType":  "USDT-FUTURES",
		"positionMode": "one_way",
	}

	mockDataBytes, _ := json.Marshal(mockPositionModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures).
		PositionMode(futures.PositionModeOneWay)

	expectedBody := map[string]interface{}{
		"productType":  "USDT-FUTURES",
		"positionMode": "one_way",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetPositionMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "one_way", result.PositionMode)
	mockClient.AssertExpectations(t)
}

func TestSetPositionModeService_Do_Success_Hedge(t *testing.T) {
	mockPositionModeData := map[string]interface{}{
		"productType":  "USDT-FUTURES",
		"positionMode": "hedge",
	}

	mockDataBytes, _ := json.Marshal(mockPositionModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures).
		PositionMode(futures.PositionModeHedge)

	expectedBody := map[string]interface{}{
		"productType":  "USDT-FUTURES",
		"positionMode": "hedge",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetPositionMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "hedge", result.PositionMode)
	mockClient.AssertExpectations(t)
}

func TestSetPositionModeService_Do_CoinFutures(t *testing.T) {
	mockPositionModeData := map[string]interface{}{
		"productType":  "COIN-FUTURES",
		"positionMode": "one_way",
	}

	mockDataBytes, _ := json.Marshal(mockPositionModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	service.ProductType(futures.ProductTypeCOINFutures).
		PositionMode(futures.PositionModeOneWay)

	expectedBody := map[string]interface{}{
		"productType":  "COIN-FUTURES",
		"positionMode": "one_way",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetPositionMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "COIN-FUTURES", result.ProductType)
	assert.Equal(t, "one_way", result.PositionMode)
	mockClient.AssertExpectations(t)
}

func TestSetPositionModeService_Do_USDCFutures(t *testing.T) {
	mockPositionModeData := map[string]interface{}{
		"productType":  "USDC-FUTURES",
		"positionMode": "hedge",
	}

	mockDataBytes, _ := json.Marshal(mockPositionModeData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDCFutures).
		PositionMode(futures.PositionModeHedge)

	expectedBody := map[string]interface{}{
		"productType":  "USDC-FUTURES",
		"positionMode": "hedge",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetPositionMode,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "USDC-FUTURES", result.ProductType)
	assert.Equal(t, "hedge", result.PositionMode)
	mockClient.AssertExpectations(t)
}

func TestSetPositionModeService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures).
		PositionMode(futures.PositionModeOneWay)

	apiError := fmt.Errorf("API error: position mode change not allowed")
	mockClient.On("CallAPI", mock.Anything, "POST", futures.EndpointSetPositionMode, mock.Anything, mock.Anything, true).
		Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "position mode change not allowed")
	mockClient.AssertExpectations(t)
}

func TestSetPositionModeService_Do_UnmarshalError(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &SetPositionModeService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures).
		PositionMode(futures.PositionModeOneWay)

	mockClient.On("CallAPI", mock.Anything, "POST", futures.EndpointSetPositionMode, mock.Anything, mock.Anything, true).
		Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestSetPositionModeService_PositionModeConstants(t *testing.T) {
	assert.Equal(t, futures.PositionModeType("one_way"), futures.PositionModeOneWay)
	assert.Equal(t, futures.PositionModeType("hedge"), futures.PositionModeHedge)
}

func TestSetPositionModeService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	service := &SetPositionModeService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		ProductType(futures.ProductTypeUSDTFutures).
		PositionMode(futures.PositionModeOneWay)

	assert.Equal(t, service, result)
	assert.Equal(t, futures.ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, futures.PositionModeOneWay, service.positionMode)
}
