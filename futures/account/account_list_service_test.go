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


func TestAccountListService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	// Test fluent API
	result := service.ProductType(futures.ProductTypeUSDTFutures)

	assert.Equal(t, futures.ProductTypeUSDTFutures, *result.productType)
}

func TestAccountListService_Do_Success_WithProductType(t *testing.T) {
	mockAccountData := []map[string]interface{}{
		{
			"marginCoin":          "USDT",
			"locked":              "0.00000000",
			"available":           "999.12345678",
			"crossMaxSize":        "5000.00000000",
			"isolatedMaxSize":     "1000.00000000",
			"maxTransferOut":      "999.12345678",
			"accountEquity":       "1000.00000000",
			"usdtEquity":          "1000.00000000",
			"btcEquity":           "0.02500000",
			"crossRiskRate":       "0.0000",
			"crossMarginLeverage": "20",
			"fixedLongLeverage":   "10",
			"fixedShortLeverage":  "10",
			"marginMode":          "CROSSED",
			"positionMode":        "one_way",
			"unrealizedPL":        "0.00000000",
			"couponAmount":        "0.00000000",
		},
		{
			"marginCoin":          "BTC",
			"locked":              "0.00000000",
			"available":           "0.05000000",
			"crossMaxSize":        "0.50000000",
			"isolatedMaxSize":     "0.10000000",
			"maxTransferOut":      "0.05000000",
			"accountEquity":       "0.05000000",
			"usdtEquity":          "2000.00000000",
			"btcEquity":           "0.05000000",
			"crossRiskRate":       "0.0000",
			"crossMarginLeverage": "15",
			"fixedLongLeverage":   "5",
			"fixedShortLeverage":  "5",
			"marginMode":          "ISOLATED",
			"positionMode":        "hedge",
			"unrealizedPL":        "0.00000000",
			"couponAmount":        "0.00000000",
		},
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures)

	expectedParams := url.Values{}
	expectedParams.Set("productType", "USDT-FUTURES")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountList,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Accounts, 2)

	// Check first account
	assert.Equal(t, "USDT", result.Accounts[0].MarginCoin)
	assert.Equal(t, "0.00000000", result.Accounts[0].Locked)
	assert.Equal(t, "999.12345678", result.Accounts[0].Available)
	assert.Equal(t, "5000.00000000", result.Accounts[0].CrossMaxSize)
	assert.Equal(t, "1000.00000000", result.Accounts[0].IsolatedMaxSize)
	assert.Equal(t, "CROSSED", result.Accounts[0].MarginMode)
	assert.Equal(t, "one_way", result.Accounts[0].PositionMode)

	// Check second account
	assert.Equal(t, "BTC", result.Accounts[1].MarginCoin)
	assert.Equal(t, "0.05000000", result.Accounts[1].Available)
	assert.Equal(t, "ISOLATED", result.Accounts[1].MarginMode)
	assert.Equal(t, "hedge", result.Accounts[1].PositionMode)
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Do_Success_WithoutProductType(t *testing.T) {
	mockAccountData := []map[string]interface{}{
		{
			"marginCoin":          "USDT",
			"locked":              "0.00000000",
			"available":           "500.00000000",
			"crossMaxSize":        "2500.00000000",
			"isolatedMaxSize":     "500.00000000",
			"maxTransferOut":      "500.00000000",
			"accountEquity":       "500.00000000",
			"usdtEquity":          "500.00000000",
			"btcEquity":           "0.01250000",
			"crossRiskRate":       "0.0000",
			"crossMarginLeverage": "10",
			"fixedLongLeverage":   "5",
			"fixedShortLeverage":  "5",
			"marginMode":          "CROSSED",
			"positionMode":        "one_way",
			"unrealizedPL":        "0.00000000",
			"couponAmount":        "0.00000000",
		},
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	// No product type set
	expectedParams := url.Values{}

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountList,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Accounts, 1)
	assert.Equal(t, "USDT", result.Accounts[0].MarginCoin)
	assert.Equal(t, "500.00000000", result.Accounts[0].Available)
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Do_CoinFutures(t *testing.T) {
	mockAccountData := []map[string]interface{}{
		{
			"marginCoin":          "BTC",
			"locked":              "0.00000000",
			"available":           "0.10000000",
			"crossMaxSize":        "1.00000000",
			"isolatedMaxSize":     "0.20000000",
			"maxTransferOut":      "0.10000000",
			"accountEquity":       "0.10000000",
			"usdtEquity":          "4000.00000000",
			"btcEquity":           "0.10000000",
			"crossRiskRate":       "0.0000",
			"crossMarginLeverage": "20",
			"fixedLongLeverage":   "10",
			"fixedShortLeverage":  "10",
			"marginMode":          "ISOLATED",
			"positionMode":        "hedge",
			"unrealizedPL":        "0.00000000",
			"couponAmount":        "0.00000000",
		},
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	service.ProductType(futures.ProductTypeCOINFutures)

	expectedParams := url.Values{}
	expectedParams.Set("productType", "COIN-FUTURES")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountList,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Accounts, 1)
	assert.Equal(t, "BTC", result.Accounts[0].MarginCoin)
	assert.Equal(t, "0.10000000", result.Accounts[0].Available)
	assert.Equal(t, "ISOLATED", result.Accounts[0].MarginMode)
	assert.Equal(t, "hedge", result.Accounts[0].PositionMode)
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Do_USDCFutures(t *testing.T) {
	mockAccountData := []map[string]interface{}{
		{
			"marginCoin":          "USDC",
			"locked":              "0.00000000",
			"available":           "750.00000000",
			"crossMaxSize":        "3750.00000000",
			"isolatedMaxSize":     "750.00000000",
			"maxTransferOut":      "750.00000000",
			"accountEquity":       "750.00000000",
			"usdtEquity":          "750.00000000",
			"btcEquity":           "0.01875000",
			"crossRiskRate":       "0.0000",
			"crossMarginLeverage": "15",
			"fixedLongLeverage":   "7",
			"fixedShortLeverage":  "7",
			"marginMode":          "CROSSED",
			"positionMode":        "one_way",
			"unrealizedPL":        "0.00000000",
			"couponAmount":        "0.00000000",
		},
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDCFutures)

	expectedParams := url.Values{}
	expectedParams.Set("productType", "USDC-FUTURES")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountList,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Accounts, 1)
	assert.Equal(t, "USDC", result.Accounts[0].MarginCoin)
	assert.Equal(t, "750.00000000", result.Accounts[0].Available)
	assert.Equal(t, "CROSSED", result.Accounts[0].MarginMode)
	assert.Equal(t, "one_way", result.Accounts[0].PositionMode)
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Do_EmptyResponse(t *testing.T) {
	mockAccountData := []map[string]interface{}{}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(mockDataBytes),
	}

	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures)

	expectedParams := url.Values{}
	expectedParams.Set("productType", "USDT-FUTURES")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountList,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Accounts, 0)
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures)

	apiError := fmt.Errorf("API error: unauthorized access")
	mockClient.On("CallAPI", mock.Anything, "GET", futures.EndpointAccountList, mock.Anything, []byte(nil), true).
		Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unauthorized access")
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Do_UnmarshalError(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &AccountListService{c: mockClient}

	service.ProductType(futures.ProductTypeUSDTFutures)

	mockClient.On("CallAPI", mock.Anything, "GET", futures.EndpointAccountList, mock.Anything, []byte(nil), true).
		Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mockClient.AssertExpectations(t)
}

func TestAccountListService_Integration(t *testing.T) {
	// Integration test style example  
	client := &MockClient{}

	service := &AccountListService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.ProductType(futures.ProductTypeUSDTFutures)

	assert.Equal(t, service, result)
	assert.Equal(t, futures.ProductTypeUSDTFutures, *service.productType)
}
