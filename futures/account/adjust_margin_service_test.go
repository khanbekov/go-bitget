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

func TestAdjustMarginService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	// Test fluent API
	result := service.
		Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("100").
		Type("ADD").
		HoldSide("LONG")

	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, futures.ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, "USDT", result.marginCoin)
	assert.Equal(t, "100", result.amount)
	assert.Equal(t, "ADD", result.type_)
	assert.Equal(t, "LONG", *result.holdSide)
}

func TestAdjustMarginService_AddMargin_Helper(t *testing.T) {
	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	// Test AddMargin helper method
	result := service.AddMargin()

	assert.Equal(t, "ADD", result.type_)
}

func TestAdjustMarginService_ReduceMargin_Helper(t *testing.T) {
	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	// Test ReduceMargin helper method
	result := service.ReduceMargin()

	assert.Equal(t, "REDUCE", result.type_)
}

func TestAdjustMarginService_Do_Success_AddMargin(t *testing.T) {
	mockAdjustMarginData := map[string]interface{}{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "100.00000000",
		"type":        "ADD",
		"holdSide":    "LONG",
		"success":     true,
	}

	mockDataBytes, _ := json.Marshal(mockAdjustMarginData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("100.00000000").
		AddMargin().
		HoldSide("LONG")

	expectedBody := map[string]interface{}{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "100.00000000",
		"type":        "ADD",
		"holdSide":    "LONG",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMargin,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "USDT", result.MarginCoin)
	assert.Equal(t, "100.00000000", result.Amount)
	assert.Equal(t, "ADD", result.Type)
	assert.Equal(t, "LONG", result.HoldSide)
	assert.Equal(t, true, result.Success)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_Success_ReduceMargin(t *testing.T) {
	mockAdjustMarginData := map[string]interface{}{
		"symbol":      "ETHUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "50.00000000",
		"type":        "REDUCE",
		"holdSide":    "SHORT",
		"success":     true,
	}

	mockDataBytes, _ := json.Marshal(mockAdjustMarginData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("ETHUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("50.00000000").
		ReduceMargin().
		HoldSide("SHORT")

	expectedBody := map[string]interface{}{
		"symbol":      "ETHUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "50.00000000",
		"type":        "REDUCE",
		"holdSide":    "SHORT",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMargin,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ETHUSDT", result.Symbol)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "USDT", result.MarginCoin)
	assert.Equal(t, "50.00000000", result.Amount)
	assert.Equal(t, "REDUCE", result.Type)
	assert.Equal(t, "SHORT", result.HoldSide)
	assert.Equal(t, true, result.Success)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_WithoutHoldSide(t *testing.T) {
	mockAdjustMarginData := map[string]interface{}{
		"symbol":      "ADAUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "25.00000000",
		"type":        "ADD",
		"holdSide":    "",
		"success":     true,
	}

	mockDataBytes, _ := json.Marshal(mockAdjustMarginData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("ADAUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("25.00000000").
		AddMargin()

	expectedBody := map[string]interface{}{
		"symbol":      "ADAUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "25.00000000",
		"type":        "ADD",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMargin,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ADAUSDT", result.Symbol)
	assert.Equal(t, "USDT-FUTURES", result.ProductType)
	assert.Equal(t, "USDT", result.MarginCoin)
	assert.Equal(t, "25.00000000", result.Amount)
	assert.Equal(t, "ADD", result.Type)
	assert.Equal(t, "", result.HoldSide)
	assert.Equal(t, true, result.Success)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_CoinFutures(t *testing.T) {
	mockAdjustMarginData := map[string]interface{}{
		"symbol":      "BTCUSD",
		"productType": "COIN-FUTURES",
		"marginCoin":  "BTC",
		"amount":      "0.01000000",
		"type":        "ADD",
		"holdSide":    "LONG",
		"success":     true,
	}

	mockDataBytes, _ := json.Marshal(mockAdjustMarginData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("BTCUSD").
		ProductType(futures.ProductTypeCOINFutures).
		MarginCoin("BTC").
		Amount("0.01000000").
		AddMargin().
		HoldSide("LONG")

	expectedBody := map[string]interface{}{
		"symbol":      "BTCUSD",
		"productType": "COIN-FUTURES",
		"marginCoin":  "BTC",
		"amount":      "0.01000000",
		"type":        "ADD",
		"holdSide":    "LONG",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMargin,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "BTCUSD", result.Symbol)
	assert.Equal(t, "COIN-FUTURES", result.ProductType)
	assert.Equal(t, "BTC", result.MarginCoin)
	assert.Equal(t, "0.01000000", result.Amount)
	assert.Equal(t, "ADD", result.Type)
	assert.Equal(t, "LONG", result.HoldSide)
	assert.Equal(t, true, result.Success)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_USDCFutures(t *testing.T) {
	mockAdjustMarginData := map[string]interface{}{
		"symbol":      "ETHUSDC",
		"productType": "USDC-FUTURES",
		"marginCoin":  "USDC",
		"amount":      "75.00000000",
		"type":        "REDUCE",
		"holdSide":    "SHORT",
		"success":     true,
	}

	mockDataBytes, _ := json.Marshal(mockAdjustMarginData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("ETHUSDC").
		ProductType(futures.ProductTypeUSDCFutures).
		MarginCoin("USDC").
		Amount("75.00000000").
		ReduceMargin().
		HoldSide("SHORT")

	expectedBody := map[string]interface{}{
		"symbol":      "ETHUSDC",
		"productType": "USDC-FUTURES",
		"marginCoin":  "USDC",
		"amount":      "75.00000000",
		"type":        "REDUCE",
		"holdSide":    "SHORT",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMargin,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "ETHUSDC", result.Symbol)
	assert.Equal(t, "USDC-FUTURES", result.ProductType)
	assert.Equal(t, "USDC", result.MarginCoin)
	assert.Equal(t, "75.00000000", result.Amount)
	assert.Equal(t, "REDUCE", result.Type)
	assert.Equal(t, "SHORT", result.HoldSide)
	assert.Equal(t, true, result.Success)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_Failed(t *testing.T) {
	mockAdjustMarginData := map[string]interface{}{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "1000.00000000",
		"type":        "REDUCE",
		"holdSide":    "LONG",
		"success":     false,
	}

	mockDataBytes, _ := json.Marshal(mockAdjustMarginData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("1000.00000000").
		ReduceMargin().
		HoldSide("LONG")

	expectedBody := map[string]interface{}{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"amount":      "1000.00000000",
		"type":        "REDUCE",
		"holdSide":    "LONG",
	}
	expectedBodyBytes, _ := json.Marshal(expectedBody)

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetMargin,
		mock.IsType(url.Values(nil)),
		expectedBodyBytes,
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "BTCUSDT", result.Symbol)
	assert.Equal(t, "REDUCE", result.Type)
	assert.Equal(t, false, result.Success)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("100").
		AddMargin()

	apiError := fmt.Errorf("API error: insufficient balance")
	mockClient.On("CallAPI", mock.Anything, "POST", futures.EndpointSetMargin, mock.Anything, mock.Anything, true).
		Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "insufficient balance")
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Do_UnmarshalError(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &AdjustMarginService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("100").
		AddMargin()

	mockClient.On("CallAPI", mock.Anything, "POST", futures.EndpointSetMargin, mock.Anything, mock.Anything, true).
		Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAdjustMarginService_Integration(t *testing.T) {
	// Integration test style example  
	client := &MockClient{}

	service := &AdjustMarginService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Amount("100").
		AddMargin().
		HoldSide("LONG")

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, futures.ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "USDT", service.marginCoin)
	assert.Equal(t, "100", service.amount)
	assert.Equal(t, "ADD", service.type_)
	assert.Equal(t, "LONG", *service.holdSide)
}
