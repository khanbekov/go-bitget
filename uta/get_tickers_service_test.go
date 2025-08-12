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

func TestGetTickersService_FluentAPI(t *testing.T) {
	mockClient := &MockClient{}
	service := &GetTickersService{c: mockClient}

	// Test fluent API chaining
	result := service.
		Category(CategoryUSDTFutures).
		Symbol("BTCUSDT")

	assert.Equal(t, service, result)
	assert.Equal(t, CategoryUSDTFutures, *service.category)
	assert.Equal(t, "BTCUSDT", *service.symbol)
}

func TestGetTickersService_Do_Success_AllTickers(t *testing.T) {
	// Setup mock data for multiple tickers
	mockTickers := []Ticker{
		{
			Symbol:       "BTCUSDT",
			Category:     CategoryUSDTFutures,
			LastPrice:    "50000.0",
			OpenPrice24h: "49000.0",
			HighPrice24h: "51000.0",
			LowPrice24h:  "48000.0",
			Ask1Price:    "50001.0",
			Bid1Price:    "49999.0",
			Bid1Size:     "10.5",
			Ask1Size:     "8.2",
			Price24hPcnt: "0.02041",
			Volume24h:    "12345.67",
			Turnover24h:  "617283350.0",
			IndexPrice:   "50000.5",
			MarkPrice:    "50000.2",
			FundingRate:  "0.0001",
			OpenInterest: "98765.43",
			Timestamp:    "1640995200000",
		},
		{
			Symbol:       "ETHUSDT",
			Category:     CategoryUSDTFutures,
			LastPrice:    "3500.0",
			OpenPrice24h: "3450.0",
			HighPrice24h: "3600.0",
			LowPrice24h:  "3400.0",
			Ask1Price:    "3500.5",
			Bid1Price:    "3499.5",
			Bid1Size:     "15.0",
			Ask1Size:     "12.3",
			Price24hPcnt: "0.01449",
			Volume24h:    "8765.43",
			Turnover24h:  "30679005.0",
			IndexPrice:   "3500.1",
			MarkPrice:    "3500.0",
			FundingRate:  "0.0001",
			OpenInterest: "54321.0",
			Timestamp:    "1640995200000",
		},
	}

	mockDataBytes, _ := json.Marshal(mockTickers)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client and service
	mockClient := &MockClient{}
	service := &GetTickersService{c: mockClient}

	// Configure service parameters
	service.Category(CategoryUSDTFutures)

	// Set up expected parameters
	expectedParams := url.Values{}
	expectedParams.Set("category", CategoryUSDTFutures)

	// Set up expected API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointMarketTickers,
		expectedParams,
		[]byte(nil),
		false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute test
	ctx := context.Background()
	result, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 2)

	// Check first ticker
	btcTicker := result[0]
	assert.Equal(t, "BTCUSDT", btcTicker.Symbol)
	assert.Equal(t, CategoryUSDTFutures, btcTicker.Category)
	assert.Equal(t, "50000.0", btcTicker.LastPrice)
	assert.Equal(t, "49000.0", btcTicker.OpenPrice24h)
	assert.Equal(t, "51000.0", btcTicker.HighPrice24h)
	assert.Equal(t, "48000.0", btcTicker.LowPrice24h)
	assert.Equal(t, "50001.0", btcTicker.Ask1Price)
	assert.Equal(t, "49999.0", btcTicker.Bid1Price)
	assert.Equal(t, "0.02041", btcTicker.Price24hPcnt)
	assert.Equal(t, "12345.67", btcTicker.Volume24h)
	assert.Equal(t, "617283350.0", btcTicker.Turnover24h)
	assert.Equal(t, "50000.5", btcTicker.IndexPrice)
	assert.Equal(t, "50000.2", btcTicker.MarkPrice)
	assert.Equal(t, "0.0001", btcTicker.FundingRate)
	assert.Equal(t, "98765.43", btcTicker.OpenInterest)

	// Check second ticker
	ethTicker := result[1]
	assert.Equal(t, "ETHUSDT", ethTicker.Symbol)
	assert.Equal(t, "3500.0", ethTicker.LastPrice)

	mockClient.AssertExpectations(t)
}

func TestGetTickersService_Do_Success_SingleTicker(t *testing.T) {
	// Setup mock data for single ticker
	mockTickers := []Ticker{
		{
			Symbol:       "BTCUSDT",
			Category:     CategorySpot,
			LastPrice:    "50000.0",
			OpenPrice24h: "49000.0",
			Ask1Price:    "50001.0",
			Bid1Price:    "49999.0",
			Timestamp:    "1640995200000",
		},
	}

	mockDataBytes, _ := json.Marshal(mockTickers)
	mockResponse := &ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := &GetTickersService{c: mockClient}

	// Configure service for single symbol
	service.Category(CategorySpot).Symbol("BTCUSDT")

	// Set up expected parameters
	expectedParams := url.Values{}
	expectedParams.Set("category", CategorySpot)
	expectedParams.Set("symbol", "BTCUSDT")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointMarketTickers,
		expectedParams,
		[]byte(nil),
		false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result, 1)
	assert.Equal(t, "BTCUSDT", result[0].Symbol)
	assert.Equal(t, CategorySpot, result[0].Category)

	mockClient.AssertExpectations(t)
}

func TestGetTickersService_Do_MissingCategory(t *testing.T) {
	mockClient := &MockClient{}
	service := &GetTickersService{c: mockClient}
	ctx := context.Background()

	// Test missing category
	_, err := service.Do(ctx)
	assert.Error(t, err)
	assert.IsType(t, &common.MissingParameterError{}, err)
}

func TestGetTickersService_Do_DifferentCategories(t *testing.T) {
	categories := []string{CategorySpot, CategoryMargin, CategoryUSDTFutures, CategoryCoinFutures, CategoryUSDCFutures}

	for _, category := range categories {
		t.Run("Category_"+category, func(t *testing.T) {
			mockTicker := []Ticker{
				{
					Symbol:    "BTCUSDT",
					Category:  category,
					LastPrice: "50000.0",
					Timestamp: "1640995200000",
				},
			}

			mockDataBytes, _ := json.Marshal(mockTicker)
			mockResponse := &ApiResponse{
				Code:        "00000",
				Msg:         "success",
				RequestTime: 1640995200000,
				Data:        mockDataBytes,
			}

			mockClient := &MockClient{}
			service := &GetTickersService{c: mockClient}
			service.Category(category)

			expectedParams := url.Values{}
			expectedParams.Set("category", category)

			mockClient.On("CallAPI",
				mock.Anything,
				"GET",
				EndpointMarketTickers,
				expectedParams,
				[]byte(nil),
				false).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

			ctx := context.Background()
			result, err := service.Do(ctx)

			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.Len(t, result, 1)
			assert.Equal(t, category, result[0].Category)

			mockClient.AssertExpectations(t)
		})
	}
}

func TestGetTickersService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &GetTickersService{c: mockClient}
	service.Category(CategoryUSDTFutures)

	expectedParams := url.Values{}
	expectedParams.Set("category", CategoryUSDTFutures)

	// Set up expected API call to return error
	expectedError := assert.AnError
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		EndpointMarketTickers,
		expectedParams,
		[]byte(nil),
		false).Return(nil, &fasthttp.ResponseHeader{}, expectedError)

	ctx := context.Background()
	result, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)

	mockClient.AssertExpectations(t)
}

func TestGetTickersService_Integration(t *testing.T) {
	// Integration-style test using the real client structure
	client := NewClient("test_api_key", "test_secret_key", "test_passphrase")
	service := client.NewGetTickersService()

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test fluent API works with real service
	service.Category(CategoryUSDTFutures).Symbol("BTCUSDT")
	assert.Equal(t, CategoryUSDTFutures, *service.category)
	assert.Equal(t, "BTCUSDT", *service.symbol)
}
