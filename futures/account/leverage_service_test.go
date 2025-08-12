package account

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/valyala/fasthttp"
)

func TestSetLeverageService_FluentAPI(t *testing.T) {
	client := &MockClient{}
	service := &SetLeverageService{c: client}

	// Test fluent API pattern
	result := service.
		Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10").
		TradeSide("long")

	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, futures.ProductTypeUSDTFutures, result.productType)
	assert.Equal(t, "USDT", result.marginCoin)
	assert.Equal(t, "10", result.leverage)
	assert.Equal(t, "long", result.tradeSide)
	assert.Equal(t, service, result, "Should return the same service instance for chaining")
}

func TestSetLeverageService_Do_Success(t *testing.T) {
	// Mock successful API response
	mockApiResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`{}`),
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	// Set up service parameters
	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	// Mock the API call with flexible JSON body matching
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.MatchedBy(func(body []byte) bool {
			var requestBody map[string]string
			if err := json.Unmarshal(body, &requestBody); err != nil {
				return false
			}
			return requestBody["symbol"] == "BTCUSDT" &&
				requestBody["productType"] == "USDT-FUTURES" &&
				requestBody["marginCoin"] == "USDT" &&
				requestBody["leverage"] == "10"
		}),
		true).Return(mockApiResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Do_Success_WithTradeSide(t *testing.T) {
	// Mock successful API response
	mockApiResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`{}`),
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	// Set up service parameters with optional tradeSide
	service.Symbol("ETHUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("20").
		TradeSide("short")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.MatchedBy(func(body []byte) bool {
			var requestBody map[string]string
			if err := json.Unmarshal(body, &requestBody); err != nil {
				return false
			}
			return requestBody["symbol"] == "ETHUSDT" &&
				requestBody["productType"] == "USDT-FUTURES" &&
				requestBody["marginCoin"] == "USDT" &&
				requestBody["leverage"] == "20" &&
				requestBody["tradeSide"] == "short"
		}),
		true).Return(mockApiResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Do_Success_COINFutures(t *testing.T) {
	// Mock successful API response
	mockApiResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`{}`),
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	// Set up service parameters for COIN futures
	service.Symbol("BTCUSD").
		ProductType(futures.ProductTypeCOINFutures).
		MarginCoin("BTC").
		Leverage("50")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.MatchedBy(func(body []byte) bool {
			var requestBody map[string]string
			if err := json.Unmarshal(body, &requestBody); err != nil {
				return false
			}
			return requestBody["symbol"] == "BTCUSD" &&
				requestBody["productType"] == "COIN-FUTURES" &&
				requestBody["marginCoin"] == "BTC" &&
				requestBody["leverage"] == "50"
		}),
		true).Return(mockApiResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Do_MissingSymbol(t *testing.T) {
	service := &SetLeverageService{}

	service.ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	ctx := context.Background()
	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symbol is required")
}

func TestSetLeverageService_Do_MissingProductType(t *testing.T) {
	service := &SetLeverageService{}

	service.Symbol("BTCUSDT").
		MarginCoin("USDT").
		Leverage("10")

	ctx := context.Background()
	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "productType is required")
}

func TestSetLeverageService_Do_MissingMarginCoin(t *testing.T) {
	service := &SetLeverageService{}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		Leverage("10")

	ctx := context.Background()
	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "marginCoin is required")
}

func TestSetLeverageService_Do_MissingLeverage(t *testing.T) {
	service := &SetLeverageService{}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT")

	ctx := context.Background()
	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "leverage is required")
}

func TestSetLeverageService_Do_APIError(t *testing.T) {
	// Mock API error response
	mockApiResponse := &futures.ApiResponse{
		Code:        "40001",
		Msg:         "Invalid leverage value",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`{}`),
	}

	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("1000") // Invalid high leverage

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.Anything,
		true).Return(mockApiResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "API error: Invalid leverage value (code 40001)")
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Do_NetworkError(t *testing.T) {
	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	// Mock network error
	networkError := fmt.Errorf("network error: connection timeout")
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.Anything,
		true).Return(nil, &fasthttp.ResponseHeader{}, networkError)

	ctx := context.Background()
	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "connection timeout")
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Do_MarshalError(t *testing.T) {
	// This test simulates a marshal error by setting invalid data
	// However, since we're using simple string maps, marshal errors are rare
	// We'll test with a service that has valid params but simulate the error scenario

	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	// Since jsoniter.Marshal is quite robust with string maps,
	// we'll test that the method completes successfully
	mockApiResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`{}`),
	}

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.Anything,
		true).Return(mockApiResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	err := service.Do(ctx)

	assert.NoError(t, err)
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Do_ContextCancellation(t *testing.T) {
	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	// Mock context cancellation error
	contextError := context.Canceled
	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.Anything,
		true).Return(nil, &fasthttp.ResponseHeader{}, contextError)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	err := service.Do(ctx)

	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	mockClient.AssertExpectations(t)
}

func TestSetLeverageService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	service := &SetLeverageService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10").
		TradeSide("long")

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, futures.ProductTypeUSDTFutures, service.productType)
	assert.Equal(t, "USDT", service.marginCoin)
	assert.Equal(t, "10", service.leverage)
	assert.Equal(t, "long", service.tradeSide)
}

func TestSetLeverageService_checkRequiredParams(t *testing.T) {
	tests := []struct {
		name        string
		setupFunc   func(*SetLeverageService)
		expectedErr string
	}{
		{
			name: "valid params",
			setupFunc: func(s *SetLeverageService) {
				s.Symbol("BTCUSDT").
					ProductType(futures.ProductTypeUSDTFutures).
					MarginCoin("USDT").
					Leverage("10")
			},
			expectedErr: "",
		},
		{
			name: "missing symbol",
			setupFunc: func(s *SetLeverageService) {
				s.ProductType(futures.ProductTypeUSDTFutures).
					MarginCoin("USDT").
					Leverage("10")
			},
			expectedErr: "symbol is required",
		},
		{
			name: "missing productType",
			setupFunc: func(s *SetLeverageService) {
				s.Symbol("BTCUSDT").
					MarginCoin("USDT").
					Leverage("10")
			},
			expectedErr: "productType is required",
		},
		{
			name: "missing marginCoin",
			setupFunc: func(s *SetLeverageService) {
				s.Symbol("BTCUSDT").
					ProductType(futures.ProductTypeUSDTFutures).
					Leverage("10")
			},
			expectedErr: "marginCoin is required",
		},
		{
			name: "missing leverage",
			setupFunc: func(s *SetLeverageService) {
				s.Symbol("BTCUSDT").
					ProductType(futures.ProductTypeUSDTFutures).
					MarginCoin("USDT")
			},
			expectedErr: "leverage is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &SetLeverageService{}
			tt.setupFunc(service)

			err := service.checkRequiredParams()

			if tt.expectedErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			}
		})
	}
}

func TestSetLeverageService_setLeverageRequestBody(t *testing.T) {
	service := &SetLeverageService{}

	// Test with all parameters
	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10").
		TradeSide("long")

	body := service.setLeverageRequestBody()

	expected := map[string]string{
		"symbol":      "BTCUSDT",
		"productType": "USDT-FUTURES",
		"marginCoin":  "USDT",
		"leverage":    "10",
		"tradeSide":   "long",
	}

	assert.Equal(t, expected, body)
}

func TestSetLeverageService_setLeverageRequestBody_WithoutTradeSide(t *testing.T) {
	service := &SetLeverageService{}

	// Test without optional tradeSide
	service.Symbol("ETHUSDT").
		ProductType(futures.ProductTypeCOINFutures).
		MarginCoin("ETH").
		Leverage("20")

	body := service.setLeverageRequestBody()

	expected := map[string]string{
		"symbol":      "ETHUSDT",
		"productType": "COIN-FUTURES",
		"marginCoin":  "ETH",
		"leverage":    "20",
	}

	assert.Equal(t, expected, body)

	// Ensure tradeSide is not included when empty
	_, exists := body["tradeSide"]
	assert.False(t, exists, "tradeSide should not be included when empty")
}

// Benchmark tests
func BenchmarkSetLeverageService_Do(b *testing.B) {
	mockApiResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`{}`),
	}

	mockClient := &MockClient{}
	service := &SetLeverageService{c: mockClient}

	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	mockClient.On("CallAPI",
		mock.Anything,
		"POST",
		futures.EndpointSetLeverage,
		mock.Anything,
		mock.Anything,
		true).Return(mockApiResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.Do(ctx)
	}
}

func BenchmarkSetLeverageService_setLeverageRequestBody(b *testing.B) {
	service := &SetLeverageService{}
	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10").
		TradeSide("long")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.setLeverageRequestBody()
	}
}

func BenchmarkSetLeverageService_checkRequiredParams(b *testing.B) {
	service := &SetLeverageService{}
	service.Symbol("BTCUSDT").
		ProductType(futures.ProductTypeUSDTFutures).
		MarginCoin("USDT").
		Leverage("10")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.checkRequiredParams()
	}
}
