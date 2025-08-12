package uta

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestAccountInfoService_Do_Success(t *testing.T) {
	// Setup mock data
	mockAccountInfo := AccountInfo{
		AssetMode:   "unified",
		HoldingMode: "one_way_mode",
		STPMode:     "none",
		SymbolConfig: []SymbolConfig{
			{
				Symbol:     "BTCUSDT",
				Category:   "USDT-FUTURES",
				Leverage:   "10",
				MarginMode: "cross",
			},
		},
		CoinConfig: []CoinConfig{
			{
				Coin:       "USDT",
				Category:   "USDT-FUTURES",
				Leverage:   "10",
				MarginMode: "cross",
			},
		},
	}

	mockDataBytes, _ := json.Marshal(mockAccountInfo)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client and service
	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	// Set up expected API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountSettings,
		mock.MatchedBy(func(params url.Values) bool { return params == nil }),
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute test
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "unified", result.AssetMode)
	assert.Equal(t, "one_way_mode", result.HoldingMode)
	assert.Equal(t, "none", result.STPMode)
	assert.Len(t, result.SymbolConfig, 1)
	assert.Equal(t, "BTCUSDT", result.SymbolConfig[0].Symbol)
	assert.Len(t, result.CoinConfig, 1)
	assert.Equal(t, "USDT", result.CoinConfig[0].Coin)

	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Do_APIError(t *testing.T) {
	// Create mock client and service
	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	// Set up expected API call to return error
	expectedError := assert.AnError
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountSettings,
		mock.MatchedBy(func(params url.Values) bool { return params == nil }),
		[]byte(nil),
		true).Return(nil, &fasthttp.ResponseHeader{}, expectedError)

	// Execute test
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Integration(t *testing.T) {
	// Integration-style test using the real client structure
	client := NewClient("test_api_key", "test_secret_key", "test_passphrase")
	service := client.NewAccountInfoService()

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)
}
