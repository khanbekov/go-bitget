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

func TestGetAccountBillService_FluentAPI(t *testing.T) {
	client := &MockClient{}
	service := &GetAccountBillService{c: client}

	// Test fluent API pattern
	result := service.
		Symbol("BTCUSDT").
		StartUnit("1640995200000").
		EndUnit("1641081600000")

	assert.Equal(t, "BTCUSDT", result.symbol)
	assert.Equal(t, "1640995200000", result.startUnit)
	assert.Equal(t, "1641081600000", result.endUnit)
	assert.Equal(t, service, result, "Should return the same service instance for chaining")
}

func TestGetAccountBillService_Do_Success(t *testing.T) {
	// Mock response data
	mockBillData := BillResponse{
		Symbol:    "BTCUSDT",
		StartUnit: "1640995200000",
		EndUnit:   "1641081600000",
	}

	// Wrap in the expected structure
	mockWrapper := struct {
		Data BillResponse `json:"data"`
	}{
		Data: mockBillData,
	}

	mockDataBytes, _ := json.Marshal(mockWrapper)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	// Set up service with all parameters
	service.Symbol("BTCUSDT").
		StartUnit("1640995200000").
		EndUnit("1641081600000")

	// Expected query parameters
	expectedParams := url.Values{}
	expectedParams.Add("symbol", "BTCUSDT")
	expectedParams.Add("startUnit", "1640995200000")
	expectedParams.Add("endUnit", "1641081600000")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	bill, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, bill)
	assert.Equal(t, "BTCUSDT", bill.Symbol)
	assert.Equal(t, "1640995200000", bill.StartUnit)
	assert.Equal(t, "1641081600000", bill.EndUnit)

	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Do_Success_WithoutOptionalParams(t *testing.T) {
	// Mock response data with only required symbol
	mockBillData := BillResponse{
		Symbol:    "ETHUSDT",
		StartUnit: "",
		EndUnit:   "",
	}

	mockWrapper := struct {
		Data BillResponse `json:"data"`
	}{
		Data: mockBillData,
	}

	mockDataBytes, _ := json.Marshal(mockWrapper)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	// Set up service with only required parameter
	service.Symbol("ETHUSDT")

	// Expected query parameters (only symbol)
	expectedParams := url.Values{}
	expectedParams.Add("symbol", "ETHUSDT")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	bill, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, bill)
	assert.Equal(t, "ETHUSDT", bill.Symbol)
	assert.Equal(t, "", bill.StartUnit)
	assert.Equal(t, "", bill.EndUnit)

	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Do_Success_WithOnlyStartUnit(t *testing.T) {
	// Mock response data with symbol and start unit
	mockBillData := BillResponse{
		Symbol:    "ADAUSDT",
		StartUnit: "1640995200000",
		EndUnit:   "",
	}

	mockWrapper := struct {
		Data BillResponse `json:"data"`
	}{
		Data: mockBillData,
	}

	mockDataBytes, _ := json.Marshal(mockWrapper)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	// Set up service with symbol and start unit
	service.Symbol("ADAUSDT").StartUnit("1640995200000")

	// Expected query parameters (symbol and startUnit)
	expectedParams := url.Values{}
	expectedParams.Add("symbol", "ADAUSDT")
	expectedParams.Add("startUnit", "1640995200000")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	bill, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, bill)
	assert.Equal(t, "ADAUSDT", bill.Symbol)
	assert.Equal(t, "1640995200000", bill.StartUnit)
	assert.Equal(t, "", bill.EndUnit)

	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Do_Success_WithOnlyEndUnit(t *testing.T) {
	// Mock response data with symbol and end unit
	mockBillData := BillResponse{
		Symbol:    "DOTUSDT",
		StartUnit: "",
		EndUnit:   "1641081600000",
	}

	mockWrapper := struct {
		Data BillResponse `json:"data"`
	}{
		Data: mockBillData,
	}

	mockDataBytes, _ := json.Marshal(mockWrapper)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	// Create mock client
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	// Set up service with symbol and end unit
	service.Symbol("DOTUSDT").EndUnit("1641081600000")

	// Expected query parameters (symbol and endUnit)
	expectedParams := url.Values{}
	expectedParams.Add("symbol", "DOTUSDT")
	expectedParams.Add("endUnit", "1641081600000")

	// Mock the API call
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	// Execute the test
	ctx := context.Background()
	bill, err := service.Do(ctx)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, bill)
	assert.Equal(t, "DOTUSDT", bill.Symbol)
	assert.Equal(t, "", bill.StartUnit)
	assert.Equal(t, "1641081600000", bill.EndUnit)

	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Do_RequiredParamValidation(t *testing.T) {
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	// Test with missing symbol (required parameter)
	ctx := context.Background()
	bill, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, bill)
	assert.Contains(t, err.Error(), "symbol is required")

	// No API call should be made
	mockClient.AssertNotCalled(t, "CallAPI")
}

func TestGetAccountBillService_Do_APIError(t *testing.T) {
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	service.Symbol("BTCUSDT").
		StartUnit("1640995200000").
		EndUnit("1641081600000")

	expectedParams := url.Values{}
	expectedParams.Add("symbol", "BTCUSDT")
	expectedParams.Add("startUnit", "1640995200000")
	expectedParams.Add("endUnit", "1641081600000")

	// Mock API error
	apiError := fmt.Errorf("API error: invalid time range")
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(nil, &fasthttp.ResponseHeader{}, apiError)

	ctx := context.Background()
	bill, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, bill)
	assert.Contains(t, err.Error(), "invalid time range")
	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Do_UnmarshalError(t *testing.T) {
	// Invalid JSON response
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        json.RawMessage(`invalid json`),
	}

	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	service.Symbol("BTCUSDT")

	expectedParams := url.Values{}
	expectedParams.Add("symbol", "BTCUSDT")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()
	bill, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, bill)
	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Do_ContextCancellation(t *testing.T) {
	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	service.Symbol("BTCUSDT")

	expectedParams := url.Values{}
	expectedParams.Add("symbol", "BTCUSDT")

	// Mock context cancellation error
	contextError := context.Canceled
	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(nil, &fasthttp.ResponseHeader{}, contextError)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	bill, err := service.Do(ctx)

	assert.Error(t, err)
	assert.Nil(t, bill)
	assert.Equal(t, context.Canceled, err)
	mockClient.AssertExpectations(t)
}

func TestGetAccountBillService_Integration(t *testing.T) {
	// Integration test style example
	client := &MockClient{}

	// Note: The client would need a method to create GetAccountBillService
	// For now, we'll test the service creation pattern
	service := &GetAccountBillService{c: client}

	assert.NotNil(t, service)
	assert.Equal(t, client, service.c)

	// Test chaining
	result := service.
		Symbol("BTCUSDT").
		StartUnit("1640995200000").
		EndUnit("1641081600000")

	assert.Equal(t, service, result)
	assert.Equal(t, "BTCUSDT", service.symbol)
	assert.Equal(t, "1640995200000", service.startUnit)
	assert.Equal(t, "1641081600000", service.endUnit)
}

func TestGetAccountBillService_checkRequiredParams(t *testing.T) {
	service := &GetAccountBillService{}

	// Test with missing symbol
	err := service.checkRequiredParams()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "symbol is required")

	// Test with symbol provided
	service.Symbol("BTCUSDT")
	err = service.checkRequiredParams()
	assert.NoError(t, err)
}

func TestBillResponse_JSONMarshaling(t *testing.T) {
	// Test JSON marshaling and unmarshaling
	original := BillResponse{
		Symbol:    "BTCUSDT",
		StartUnit: "1640995200000",
		EndUnit:   "1641081600000",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaled BillResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// Verify all fields are preserved
	assert.Equal(t, original.Symbol, unmarshaled.Symbol)
	assert.Equal(t, original.StartUnit, unmarshaled.StartUnit)
	assert.Equal(t, original.EndUnit, unmarshaled.EndUnit)
}

func TestBillResponse_JSONMarshaling_WithEmptyFields(t *testing.T) {
	// Test with empty optional fields
	original := BillResponse{
		Symbol:    "ETHUSDT",
		StartUnit: "",
		EndUnit:   "",
	}

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	assert.NoError(t, err)

	// Unmarshal back
	var unmarshaled BillResponse
	err = json.Unmarshal(jsonData, &unmarshaled)
	assert.NoError(t, err)

	// Verify all fields are preserved
	assert.Equal(t, original.Symbol, unmarshaled.Symbol)
	assert.Equal(t, original.StartUnit, unmarshaled.StartUnit)
	assert.Equal(t, original.EndUnit, unmarshaled.EndUnit)
}

func TestBillResponse_JSONUnmarshaling_InvalidJSON(t *testing.T) {
	jsonData := `invalid json`

	var bill BillResponse
	err := json.Unmarshal([]byte(jsonData), &bill)

	assert.Error(t, err)
}

func TestBillResponse_JSONUnmarshaling_MissingFields(t *testing.T) {
	jsonData := `{"symbol": "BTCUSDT"}`

	var bill BillResponse
	err := json.Unmarshal([]byte(jsonData), &bill)

	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", bill.Symbol)
	assert.Equal(t, "", bill.StartUnit) // Missing fields should be empty
	assert.Equal(t, "", bill.EndUnit)
}

func TestBillResponse_JSONUnmarshaling_ExtraFields(t *testing.T) {
	jsonData := `{
		"symbol": "BTCUSDT",
		"startUnit": "1640995200000",
		"endUnit": "1641081600000",
		"extraField": "should be ignored"
	}`

	var bill BillResponse
	err := json.Unmarshal([]byte(jsonData), &bill)

	assert.NoError(t, err)
	assert.Equal(t, "BTCUSDT", bill.Symbol)
	assert.Equal(t, "1640995200000", bill.StartUnit)
	assert.Equal(t, "1641081600000", bill.EndUnit)
	// Extra fields should be ignored
}

// Benchmark tests
func BenchmarkGetAccountBillService_Do(b *testing.B) {
	mockBillData := BillResponse{
		Symbol:    "BTCUSDT",
		StartUnit: "1640995200000",
		EndUnit:   "1641081600000",
	}

	mockWrapper := struct {
		Data BillResponse `json:"data"`
	}{
		Data: mockBillData,
	}

	mockDataBytes, _ := json.Marshal(mockWrapper)
	mockResponse := &futures.ApiResponse{
		Code:        "00000",
		Msg:         "success",
		RequestTime: 1640995200000,
		Data:        mockDataBytes,
	}

	mockClient := &MockClient{}
	service := &GetAccountBillService{c: mockClient}

	service.Symbol("BTCUSDT").
		StartUnit("1640995200000").
		EndUnit("1641081600000")

	expectedParams := url.Values{}
	expectedParams.Add("symbol", "BTCUSDT")
	expectedParams.Add("startUnit", "1640995200000")
	expectedParams.Add("endUnit", "1641081600000")

	mockClient.On("CallAPI",
		mock.Anything,
		"GET",
		futures.EndpointAccountBills,
		expectedParams,
		[]byte(nil),
		true).Return(mockResponse, &fasthttp.ResponseHeader{}, nil)

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Do(ctx)
	}
}

func BenchmarkBillResponse_JSONMarshaling(b *testing.B) {
	bill := BillResponse{
		Symbol:    "BTCUSDT",
		StartUnit: "1640995200000",
		EndUnit:   "1641081600000",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(bill)
	}
}

func BenchmarkBillResponse_JSONUnmarshaling(b *testing.B) {
	jsonData := `{
		"symbol": "BTCUSDT",
		"startUnit": "1640995200000",
		"endUnit": "1641081600000"
	}`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var bill BillResponse
		_ = json.Unmarshal([]byte(jsonData), &bill)
	}
}
