package account

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


func TestAccountInfoService_FluentAPI(t *testing.T) {
	client := &MockClient{}
	service := &AccountInfoService{c: client}

	// Test fluent API pattern
	result := service.
		Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, "USDT", result.marginCoin)
	assert.Equal(t, service, result, "Should return the same service instance for chaining")
}

func TestAccountInfoService_Do_Success(t *testing.T) {
	// Mock response data
	mockAccountData := &Account{
		MarginCoin:            "USDT",
		Locked:                100.50,
		Available:             1000.25,
		CrossedMaxAvailable:   950.75,
		IsolatedMaxAvailable:  800.00,
		MaxTransferOut:        900.00,
		AccountEquity:         1100.75,
		UsdtEquity:            1100.75,
		BtcEquity:             0.025,
		CrossedRiskRate:       0.15,
		CrossedMarginLeverage: 10,
		IsolatedLongLever:     20,
		IsolatedShortLever:    15,
		MarginMode:            "isolated",
		PosMode:               "one_way_mode",
		UnrealizedPL:          25.50,
		Coupon:                0,
		CrossedUnrealizedPL:   0.00,
		IsolatedUnrealizedPL:  25.50,
		AssetMode:             "single",
		Grant:                 "",
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	// Set up service parameters
	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	// Expected query parameters
	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("marginCoin", "USDT")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountInfo,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	account, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, "USDT", account.MarginCoin)
	assert.Equal(t, 100.50, account.Locked)
	assert.Equal(t, 1000.25, account.Available)
	assert.Equal(t, 950.75, account.CrossedMaxAvailable)
	assert.Equal(t, 800.00, account.IsolatedMaxAvailable)
	assert.Equal(t, 900.00, account.MaxTransferOut)
	assert.Equal(t, 1100.75, account.AccountEquity)
	assert.Equal(t, 1100.75, account.UsdtEquity)
	assert.Equal(t, 0.025, account.BtcEquity)
	assert.Equal(t, 0.15, account.CrossedRiskRate)
	assert.Equal(t, int64(10), account.CrossedMarginLeverage)
	assert.Equal(t, int64(20), account.IsolatedLongLever)
	assert.Equal(t, int64(15), account.IsolatedShortLever)
	assert.Equal(t, "isolated", account.MarginMode)
	assert.Equal(t, "one_way_mode", account.PosMode)
	assert.Equal(t, 25.50, account.UnrealizedPL)
	assert.Equal(t, int64(0), account.Coupon)
	assert.Equal(t, 0.00, account.CrossedUnrealizedPL)
	assert.Equal(t, 25.50, account.IsolatedUnrealizedPL)
	assert.Equal(t, "single", account.AssetMode)
	assert.Equal(t, "", account.Grant)

	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Do_Success_WithDifferentProductType(t *testing.T) {
	// Mock response data for COIN futures
	mockAccountData := &Account{
		MarginCoin:            "BTC",
		Locked:                0.01,
		Available:             0.05,
		CrossedMaxAvailable:   0.045,
		IsolatedMaxAvailable:  0.04,
		MaxTransferOut:        0.045,
		AccountEquity:         0.06,
		UsdtEquity:            2500.00,
		BtcEquity:             0.06,
		CrossedRiskRate:       0.20,
		CrossedMarginLeverage: 5,
		IsolatedLongLever:     10,
		IsolatedShortLever:    10,
		MarginMode:            "crossed",
		PosMode:               "hedge_mode",
		UnrealizedPL:          0.001,
		Coupon:                100,
		CrossedUnrealizedPL:   0.001,
		IsolatedUnrealizedPL:  0.00,
		AssetMode:             "union",
		Grant:                 "premium",
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	// Set up service parameters for COIN futures
	service.Symbol("BTCUSD").
		ProductType(ProductTypeCoinFutures).
		MarginCoin("BTC")

	// Expected query parameters
	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSD")
	expectedParams.Set("productType", "COIN-FUTURES")
	expectedParams.Set("marginCoin", "BTC")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountInfo,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	account, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, account)
	assert.Equal(t, "BTC", account.MarginCoin)
	assert.Equal(t, 0.01, account.Locked)
	assert.Equal(t, 0.05, account.Available)
	assert.Equal(t, "crossed", account.MarginMode)
	assert.Equal(t, "hedge_mode", account.PosMode)
	assert.Equal(t, "union", account.AssetMode)
	assert.Equal(t, "premium", account.Grant)
	assert.Equal(t, int64(100), account.Coupon)

	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("marginCoin", "USDT")

	// Mock API error
	apiError := fmt.Errorf("API error: insufficient permissions")
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountInfo,
		expectedParams,
		[]byte(nil),
		true).Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	account, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Contains(t, err.Error(), "insufficient permissions")
	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Do_UnmarshalError(t *testing.T) {
	// Invalid JSON response
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("marginCoin", "USDT")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountInfo,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	account, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, account)
	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Do_ContextCancellation(t *testing.T) {
	// Test context cancellation during API call
	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("marginCoin", "USDT")

	// Mock context cancellation error
	contextError := context.Canceled
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountInfo,
		expectedParams,
		[]byte(nil),
		true).Return(nil, &fasthttp.ResponseHeader{}, contextError)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	account, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, account)
	assert.Equal(t, context.Canceled, err)
	mockClient.AssertExpectations(t)
}

func TestAccountInfoService_Integration(t *testing.T) {
	// Integration test style example  
	client := &MockClient{}

	service := &AccountInfoService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "USDT", service.marginCoin)
}

func TestAccount_UnmarshalJSON_Success(t *testing.T) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "100.50",
		"available": "1000.25",
		"crossedMaxAvailable": "950.75",
		"isolatedMaxAvailable": "800.00",
		"maxTransferOut": "900.00",
		"accountEquity": "1100.75",
		"usdtEquity": "1100.75",
		"btcEquity": "0.025",
		"crossedRiskRate": "0.15",
		"crossedMarginLeverage": "10",
		"isolatedLongLever": "20",
		"isolatedShortLever": "15",
		"marginMode": "isolated",
		"posMode": "one_way_mode",
		"unrealizedPL": "25.50",
		"coupon": "0",
		"crossedUnrealizedPL": "0.00",
		"isolatedUnrealizedPL": "25.50",
		"assetMode": "single",
		"grant": ""
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.NoError(t, err)
	assert.Equal(t, "USDT", account.MarginCoin)
	assert.Equal(t, 100.50, account.Locked)
	assert.Equal(t, 1000.25, account.Available)
	assert.Equal(t, 950.75, account.CrossedMaxAvailable)
	assert.Equal(t, 800.00, account.IsolatedMaxAvailable)
	assert.Equal(t, 900.00, account.MaxTransferOut)
	assert.Equal(t, 1100.75, account.AccountEquity)
	assert.Equal(t, 1100.75, account.UsdtEquity)
	assert.Equal(t, 0.025, account.BtcEquity)
	assert.Equal(t, 0.15, account.CrossedRiskRate)
	assert.Equal(t, int64(10), account.CrossedMarginLeverage)
	assert.Equal(t, int64(20), account.IsolatedLongLever)
	assert.Equal(t, int64(15), account.IsolatedShortLever)
	assert.Equal(t, "isolated", account.MarginMode)
	assert.Equal(t, "one_way_mode", account.PosMode)
	assert.Equal(t, 25.50, account.UnrealizedPL)
	assert.Equal(t, int64(0), account.Coupon)
	assert.Equal(t, 0.00, account.CrossedUnrealizedPL)
	assert.Equal(t, 25.50, account.IsolatedUnrealizedPL)
	assert.Equal(t, "single", account.AssetMode)
	assert.Equal(t, "", account.Grant)
}

func TestAccount_UnmarshalJSON_WithNumericValues(t *testing.T) {
	jsonData := `{
		"marginCoin": "BTC",
		"locked": 0.01,
		"available": 0.05,
		"crossedMaxAvailable": 0.045,
		"isolatedMaxAvailable": 0.04,
		"maxTransferOut": 0.045,
		"accountEquity": 0.06,
		"usdtEquity": 2500.00,
		"btcEquity": 0.06,
		"crossedRiskRate": 0.20,
		"crossedMarginLeverage": 5,
		"isolatedLongLever": 10,
		"isolatedShortLever": 10,
		"marginMode": "crossed",
		"posMode": "hedge_mode",
		"unrealizedPL": 0.001,
		"coupon": 100,
		"crossedUnrealizedPL": 0.001,
		"isolatedUnrealizedPL": 0.00,
		"assetMode": "union",
		"grant": "premium"
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.NoError(t, err)
	assert.Equal(t, "BTC", account.MarginCoin)
	assert.Equal(t, 0.01, account.Locked)
	assert.Equal(t, 0.05, account.Available)
	assert.Equal(t, 0.045, account.CrossedMaxAvailable)
	assert.Equal(t, 0.04, account.IsolatedMaxAvailable)
	assert.Equal(t, 0.045, account.MaxTransferOut)
	assert.Equal(t, 0.06, account.AccountEquity)
	assert.Equal(t, 2500.00, account.UsdtEquity)
	assert.Equal(t, 0.06, account.BtcEquity)
	assert.Equal(t, 0.20, account.CrossedRiskRate)
	assert.Equal(t, int64(5), account.CrossedMarginLeverage)
	assert.Equal(t, int64(10), account.IsolatedLongLever)
	assert.Equal(t, int64(10), account.IsolatedShortLever)
	assert.Equal(t, "crossed", account.MarginMode)
	assert.Equal(t, "hedge_mode", account.PosMode)
	assert.Equal(t, 0.001, account.UnrealizedPL)
	assert.Equal(t, int64(100), account.Coupon)
	assert.Equal(t, 0.001, account.CrossedUnrealizedPL)
	assert.Equal(t, 0.00, account.IsolatedUnrealizedPL)
	assert.Equal(t, "union", account.AssetMode)
	assert.Equal(t, "premium", account.Grant)
}

func TestAccount_UnmarshalJSON_WithEmptyValues(t *testing.T) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "",
		"available": "",
		"crossedMaxAvailable": null,
		"isolatedMaxAvailable": null,
		"maxTransferOut": "",
		"accountEquity": "",
		"usdtEquity": "",
		"btcEquity": "",
		"crossedRiskRate": "",
		"crossedMarginLeverage": "",
		"isolatedLongLever": "",
		"isolatedShortLever": "",
		"marginMode": "isolated",
		"posMode": "one_way_mode",
		"unrealizedPL": "",
		"coupon": "",
		"crossedUnrealizedPL": "",
		"isolatedUnrealizedPL": "",
		"assetMode": "single"
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.NoError(t, err)
	assert.Equal(t, "USDT", account.MarginCoin)
	assert.Equal(t, 0.0, account.Locked)
	assert.Equal(t, 0.0, account.Available)
	assert.Equal(t, 0.0, account.CrossedMaxAvailable)
	assert.Equal(t, 0.0, account.IsolatedMaxAvailable)
	assert.Equal(t, 0.0, account.MaxTransferOut)
	assert.Equal(t, 0.0, account.AccountEquity)
	assert.Equal(t, 0.0, account.UsdtEquity)
	assert.Equal(t, 0.0, account.BtcEquity)
	assert.Equal(t, 0.0, account.CrossedRiskRate)
	assert.Equal(t, int64(0), account.CrossedMarginLeverage)
	assert.Equal(t, int64(0), account.IsolatedLongLever)
	assert.Equal(t, int64(0), account.IsolatedShortLever)
	assert.Equal(t, "isolated", account.MarginMode)
	assert.Equal(t, "one_way_mode", account.PosMode)
	assert.Equal(t, 0.0, account.UnrealizedPL)
	assert.Equal(t, int64(0), account.Coupon)
	assert.Equal(t, 0.0, account.CrossedUnrealizedPL)
	assert.Equal(t, 0.0, account.IsolatedUnrealizedPL)
	assert.Equal(t, "single", account.AssetMode)
	assert.Equal(t, "", account.Grant)
}

func TestAccount_UnmarshalJSON_WithMissingFields(t *testing.T) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "100.50",
		"available": "1000.25",
		"marginMode": "isolated",
		"posMode": "one_way_mode",
		"assetMode": "single"
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.NoError(t, err)
	assert.Equal(t, "USDT", account.MarginCoin)
	assert.Equal(t, 100.50, account.Locked)
	assert.Equal(t, 1000.25, account.Available)
	assert.Equal(t, "isolated", account.MarginMode)
	assert.Equal(t, "one_way_mode", account.PosMode)
	assert.Equal(t, "single", account.AssetMode)
	// Missing fields should have default values
	assert.Equal(t, 0.0, account.CrossedMaxAvailable)
	assert.Equal(t, 0.0, account.IsolatedMaxAvailable)
	assert.Equal(t, int64(0), account.CrossedMarginLeverage)
	assert.Equal(t, "", account.Grant)
}

func TestAccount_UnmarshalJSON_WithUnknownFields(t *testing.T) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "100.50",
		"available": "1000.25",
		"marginMode": "isolated",
		"posMode": "one_way_mode",
		"assetMode": "single",
		"unknownField": "should be ignored",
		"anotherUnknownField": 12345
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.NoError(t, err)
	assert.Equal(t, "USDT", account.MarginCoin)
	assert.Equal(t, 100.50, account.Locked)
	assert.Equal(t, 1000.25, account.Available)
	assert.Equal(t, "isolated", account.MarginMode)
	assert.Equal(t, "one_way_mode", account.PosMode)
	assert.Equal(t, "single", account.AssetMode)
	// Unknown fields should be ignored
}

func TestAccount_UnmarshalJSON_InvalidJSON(t *testing.T) {
	jsonData := `invalid json`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.Error(t, err)
}

func TestAccount_UnmarshalJSON_InvalidFloatConversion(t *testing.T) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "invalid_float",
		"available": "1000.25",
		"marginMode": "isolated"
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unable to convert string")
}

func TestAccount_UnmarshalJSON_InvalidIntConversion(t *testing.T) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "100.50",
		"available": "1000.25",
		"crossedMarginLeverage": "invalid_int",
		"marginMode": "isolated"
	}`

	var account Account
	err := json.Unmarshal([]byte(jsonData), &account)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Unable to convert string")
}

func TestAccount_UnmarshalJSON_InvalidStringTypeAssertion(t *testing.T) {
	jsonData := `{
		"marginCoin": 123,
		"locked": "100.50",
		"available": "1000.25",
		"marginMode": "isolated"
	}`

	var account Account

	// This should panic due to type assertion, so we catch it
	defer func() {
		if r := recover(); r != nil {
			// Expected panic due to type assertion
			assert.Contains(t, fmt.Sprintf("%v", r), "interface conversion")
		}
	}()

	err := json.Unmarshal([]byte(jsonData), &account)

	// If we reach here without panic, the test should fail
	if err == nil {
		t.Error("Expected panic or error, but got nil")
	}
}

// Benchmark tests
func BenchmarkAccountInfoService_Do(b *testing.B) {
	mockAccountData := &Account{
		MarginCoin:            "USDT",
		Locked:                100.50,
		Available:             1000.25,
		CrossedMaxAvailable:   950.75,
		IsolatedMaxAvailable:  800.00,
		MaxTransferOut:        900.00,
		AccountEquity:         1100.75,
		UsdtEquity:            1100.75,
		BtcEquity:             0.025,
		CrossedRiskRate:       0.15,
		CrossedMarginLeverage: 10,
		IsolatedLongLever:     20,
		IsolatedShortLever:    15,
		MarginMode:            "isolated",
		PosMode:               "one_way_mode",
		UnrealizedPL:          25.50,
		Coupon:                0,
		CrossedUnrealizedPL:   0.00,
		IsolatedUnrealizedPL:  25.50,
		AssetMode:             "single",
		Grant:                 "",
	}

	mockDataBytes, _ := json.Marshal(mockAccountData)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := &AccountInfoService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		MarginCoin("USDT")

	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("marginCoin", "USDT")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointAccountInfo,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Do(ctx)
	}
}

func BenchmarkAccount_UnmarshalJSON(b *testing.B) {
	jsonData := `{
		"marginCoin": "USDT",
		"locked": "100.50",
		"available": "1000.25",
		"crossedMaxAvailable": "950.75",
		"isolatedMaxAvailable": "800.00",
		"maxTransferOut": "900.00",
		"accountEquity": "1100.75",
		"usdtEquity": "1100.75",
		"btcEquity": "0.025",
		"crossedRiskRate": "0.15",
		"crossedMarginLeverage": "10",
		"isolatedLongLever": "20",
		"isolatedShortLever": "15",
		"marginMode": "isolated",
		"posMode": "one_way_mode",
		"unrealizedPL": "25.50",
		"coupon": "0",
		"crossedUnrealizedPL": "0.00",
		"isolatedUnrealizedPL": "25.50",
		"assetMode": "single",
		"grant": ""
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var account Account
		_ = json.Unmarshal([]byte(jsonData), &account)
	}
}
