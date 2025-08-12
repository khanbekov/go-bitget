package market

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

func TestCandlestickService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &CandlestickService{c: mockClient}

	// Test fluent API pattern
	result := service.
		Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		Granularity("1m").
		Limit("100").
		StartTime("1640995200000").
		EndTime("1641081600000")

	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, "1m", result.granularity)
	assert.Equal(t, "100", result.limit)
	assert.Equal(t, "1640995200000", result.startTime)
	assert.Equal(t, "1641081600000", result.endTime)
	assert.Equal(t, service, result, "Should return the same service instance for chaining")
}

func TestCandlestickService_Do_Success(t *testing.T) {
	// Mock response data
	mockCandlestickData := [][]string{
		{"1640995200000", "47000.50", "47500.00", "46500.00", "47200.00", "1.5", "70800.00"},
		{"1640995260000", "47200.00", "47300.00", "46800.00", "47000.00", "2.0", "94000.00"},
	}

	mockDataBytes, _ := json.Marshal(mockCandlestickData)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &CandlestickService{c: mockClient}

	// Set up service parameters
	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		Granularity("1m").
		Limit("2")

	// Expected query parameters
	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("granularity", "1m")
	expectedParams.Set("limit", "2")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointCandlesticks,
		expectedParams,
		[]byte(nil),
		false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	candles, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, candles, 2)

	// Check first candlestick
	assert.Equal(t, int64(1640995200000), candles[0].CloseTime)
	assert.Equal(t, 47000.50, candles[0].Open)
	assert.Equal(t, 47500.00, candles[0].High)
	assert.Equal(t, 46500.00, candles[0].Low)
	assert.Equal(t, 47200.00, candles[0].Close)
	assert.Equal(t, 1.5, candles[0].Volume)
	assert.Equal(t, 70800.00, candles[0].QuoteAssetVolume)

	// Check second candlestick
	assert.Equal(t, int64(1640995260000), candles[1].CloseTime)
	assert.Equal(t, 47200.00, candles[1].Open)
	assert.Equal(t, 47300.00, candles[1].High)
	assert.Equal(t, 46800.00, candles[1].Low)
	assert.Equal(t, 47000.00, candles[1].Close)
	assert.Equal(t, 2.0, candles[1].Volume)
	assert.Equal(t, 94000.00, candles[1].QuoteAssetVolume)

	mockClient.AssertExpectations(t)
}

func TestCandlestickService_Do_WithOptionalParams(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`[]`),
	}

	mockClient := &MockClient{}
	service := &CandlestickService{c: mockClient}

	// Set up service with optional parameters
	service.Symbol("ETHUSDT").
		ProductType(ProductTypeUSDCFutures).
		Granularity("5m").
		Limit("50").
		StartTime("1640995200000").
		EndTime("1641081600000")

	// Expected query parameters including optional ones
	expectedParams := url.Values{}
	expectedParams.Set("symbol", "ETHUSDT")
	expectedParams.Set("productType", "USDC-FUTURES")
	expectedParams.Set("granularity", "5m")
	expectedParams.Set("limit", "50")
	expectedParams.Set("startTime", "1640995200000")
	expectedParams.Set("endTime", "1641081600000")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointCandlesticks,
		expectedParams,
		[]byte(nil),
		false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	candles, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.Empty(t, candles)
	mockClient.AssertExpectations(t)
}

func TestCandlestickService_Do_WithoutOptionalParams(t *testing.T) {
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`[]`),
	}

	mockClient := &MockClient{}
	service := &CandlestickService{c: mockClient}

	// Set up service without optional parameters
	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		Granularity("1h")

	// Expected query parameters without optional ones
	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("granularity", "1h")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointCandlesticks,
		expectedParams,
		[]byte(nil),
		false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	candles, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.Empty(t, candles)
	mockClient.AssertExpectations(t)
}

func TestCandlestickService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &CandlestickService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		Granularity("1m")

	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("granularity", "1m")

	// Mock API error
	apiError := fmt.Errorf("API error: rate limit exceeded")
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointCandlesticks,
		expectedParams,
		[]byte(nil),
		false).Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	candles, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, candles)
	assert.Contains(t, err.Error(), "rate limit exceeded")
	mockClient.AssertExpectations(t)
}

func TestCandlestickService_Do_UnmarshalError(t *testing.T) {
	// Invalid JSON response
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &CandlestickService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		Granularity("1m")

	expectedParams := url.Values{}
	expectedParams.Set("symbol", "BTCUSDT")
	expectedParams.Set("productType", "USDT-FUTURES")
	expectedParams.Set("granularity", "1m")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointCandlesticks,
		expectedParams,
		[]byte(nil),
		false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	candles, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, candles)
	mockClient.AssertExpectations(t)
}

func TestCandlestick_UnmarshalJSON_Success(t *testing.T) {
	jsonData := `["1640995200000", "47000.50", "47500.00", "46500.00", "47200.00", "1.5", "70800.00"]`

	var candlestick Candlestick
	err := json.Unmarshal([]byte(jsonData), &candlestick)

	assert.NoError(t, err)
	assert.Equal(t, int64(1640995200000), candlestick.CloseTime)
	assert.Equal(t, 47000.50, candlestick.Open)
	assert.Equal(t, 47500.00, candlestick.High)
	assert.Equal(t, 46500.00, candlestick.Low)
	assert.Equal(t, 47200.00, candlestick.Close)
	assert.Equal(t, 1.5, candlestick.Volume)
	assert.Equal(t, 70800.00, candlestick.QuoteAssetVolume)
}

func TestCandlestick_UnmarshalJSON_InvalidLength(t *testing.T) {
	jsonData := `["1640995200000", "47000.50", "47500.00"]` // Only 3 elements instead of 7

	var candlestick Candlestick
	err := json.Unmarshal([]byte(jsonData), &candlestick)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected 7 elements, received 3")
}

func TestCandlestick_UnmarshalJSON_InvalidTimestamp(t *testing.T) {
	jsonData := `["invalid_timestamp", "47000.50", "47500.00", "46500.00", "47200.00", "1.5", "70800.00"]`

	var candlestick Candlestick
	err := json.Unmarshal([]byte(jsonData), &candlestick)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "incorrect CloseTime")
}

func TestCandlestick_UnmarshalJSON_InvalidFloatValues(t *testing.T) {
	jsonData := `["1640995200000", "invalid_float", "47500.00", "46500.00", "47200.00", "1.5", "70800.00"]`

	var candlestick Candlestick
	err := json.Unmarshal([]byte(jsonData), &candlestick)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed parsing element")
}

func TestCandlestick_UnmarshalJSON_InvalidJSONStructure(t *testing.T) {
	jsonData := `{"not": "an_array"}`

	var candlestick Candlestick
	err := json.Unmarshal([]byte(jsonData), &candlestick)

	assert.Error(t, err)
}

// Integration test style example
func TestCandlestickService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	service := &CandlestickService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		Symbol("BTCUSDT").
		ProductType(ProductTypeUSDTFutures).
		Granularity("1m").
		Limit("100")

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "1m", service.granularity)
	assert.Equal(t, "100", service.limit)
}

// Benchmark test for UnmarshalJSON performance
func BenchmarkCandlestick_UnmarshalJSON(b *testing.B) {
	jsonData := `["1640995200000", "47000.50", "47500.00", "46500.00", "47200.00", "1.5", "70800.00"]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var candlestick Candlestick
		_ = json.Unmarshal([]byte(jsonData), &candlestick)
	}
}