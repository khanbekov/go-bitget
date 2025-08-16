//go:build integration
// +build integration

package suites

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/tests/integration"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// AccountTestSuite manages integration tests for account endpoints
type AccountTestSuite struct {
	client    *futures.Client
	config    *integration.TestConfig
	logger    zerolog.Logger
	results   []integration.TestResult
	startTime time.Time
}

// NewAccountTestSuite creates a new account test suite
func NewAccountTestSuite(config *integration.TestConfig) *AccountTestSuite {
	client := futures.NewClient(config.APIKey, config.SecretKey, config.Passphrase)

	// Set the base URL from config if specified
	if config.BaseURL != "" {
		client.SetApiEndpoint(config.BaseURL)
	}

	return &AccountTestSuite{
		client:    client,
		config:    config,
		logger:    config.GetLogger(),
		results:   make([]integration.TestResult, 0),
		startTime: time.Now(),
	}
}

// TestAccountEndpoints runs all enabled account endpoint tests
func TestAccountEndpoints(t *testing.T) {
	config, err := integration.LoadConfig()
	require.NoError(t, err, "Failed to load test configuration")

	suite := NewAccountTestSuite(config)

	suite.logger.Info().
		Str("base_url", config.BaseURL).
		Bool("demo_trading", config.DemoTrading).
		Str("test_symbol", config.TestSymbol).
		Msg("Starting account integration tests")

	// Run tests based on configuration
	if config.IsEndpointEnabled("account_info") {
		suite.TestAccountInfo(t)
	}

	if config.IsEndpointEnabled("account_list") {
		suite.TestAccountList(t)
	}

	if config.IsEndpointEnabled("account_bills") {
		suite.TestAccountBills(t)
	}

	if config.IsEndpointEnabled("set_leverage") {
		suite.TestSetLeverage(t)
	}

	if config.IsEndpointEnabled("adjust_margin") {
		suite.TestAdjustMargin(t)
	}

	if config.IsEndpointEnabled("set_margin_mode") {
		suite.TestSetMarginMode(t)
	}

	if config.IsEndpointEnabled("set_position_mode") {
		suite.TestSetPositionMode(t)
	}

	// Generate test report
	if config.GenerateReport {
		err := suite.GenerateReport()
		if err != nil {
			suite.logger.Error().Err(err).Msg("Failed to generate test report")
		}
	}

	suite.PrintSummary()
}

// TestAccountInfo tests the account information endpoint
func (s *AccountTestSuite) TestAccountInfo(t *testing.T) {
	testName := "account_info"
	s.logger.Info().Str("test", testName).Msg("Testing account info endpoint")

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/account",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	// Test with USDT-FUTURES product type
	service := account.NewAccountInfoService(s.client).
		ProductType(account.ProductType(s.config.TestProductType)).
		MarginCoin(s.config.TestCoin)

	response, err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		t.Errorf("Account info test failed: %v", err)
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("MarginCoin: %s, USDT Equity: %.2f",
			response.MarginCoin, response.UsdtEquity)

		s.logger.Info().
			Str("test", testName).
			Str("margin_coin", response.MarginCoin).
			Float64("usdt_equity", response.UsdtEquity).
			Float64("btc_equity", response.BtcEquity).
			Msg("Test passed")

		// Validate response structure
		assert.NotEmpty(t, response.MarginCoin, "Margin coin should not be empty")
		assert.GreaterOrEqual(t, response.UsdtEquity, 0.0, "USDT equity should be non-negative")
	}

	s.results = append(s.results, result)
}

// TestAccountList tests the account list endpoint
func (s *AccountTestSuite) TestAccountList(t *testing.T) {
	testName := "account_list"
	s.logger.Info().Str("test", testName).Msg("Testing account list endpoint")

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/accounts",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	service := account.NewAccountListService(s.client).
		ProductType(futures.ProductType(s.config.TestProductType))

	response, err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		t.Errorf("Account list test failed: %v", err)
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("Found %d accounts", len(response.Accounts))

		s.logger.Info().
			Str("test", testName).
			Int("account_count", len(response.Accounts)).
			Msg("Test passed")

		// Validate response structure
		assert.NotNil(t, response, "Response should not be nil")
		if len(response.Accounts) > 0 {
			assert.NotEmpty(t, response.Accounts[0].MarginCoin, "First account margin coin should not be empty")
		}
	}

	s.results = append(s.results, result)
}

// TestAccountBills tests the account bills/transaction history endpoint
func (s *AccountTestSuite) TestAccountBills(t *testing.T) {
	testName := "account_bills"
	s.logger.Info().Str("test", testName).Msg("Testing account bills endpoint")

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/bill",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	service := account.NewGetAccountBillService(s.client).
		Symbol(s.config.TestSymbol)

	response, err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		t.Errorf("Account bills test failed: %v", err)
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("Bill info retrieved for %s", response.Symbol)

		s.logger.Info().
			Str("test", testName).
			Str("symbol", response.Symbol).
			Msg("Test passed")

		// Validate response structure
		assert.NotNil(t, response, "Response should not be nil")
		assert.NotEmpty(t, response.Symbol, "Symbol should not be empty")
	}

	s.results = append(s.results, result)
}

// TestSetLeverage tests the set leverage endpoint (with caution - has side effects)
func (s *AccountTestSuite) TestSetLeverage(t *testing.T) {
	testName := "set_leverage"
	s.logger.Warn().Str("test", testName).Msg("Testing set leverage endpoint - HAS SIDE EFFECTS")

	if !s.config.DemoTrading {
		s.logger.Error().Str("test", testName).Msg("Skipping leverage test - not in demo mode")
		t.Skip("Leverage test skipped - not in demo trading mode")
		return
	}

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/set-leverage",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	// Use safe leverage value (20x is typically safe for testing)
	service := account.NewSetLeverageService(s.client).
		Symbol(s.config.TestSymbol).
		ProductType(futures.ProductType(s.config.TestProductType)).
		MarginCoin(s.config.TestCoin).
		Leverage("20").
		HoldSide(string(account.HoldSideLong))

	err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		// Don't fail the test suite for leverage issues in demo mode
		s.logger.Warn().Msg("Leverage test failed - this may be expected in demo mode")
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("Set leverage for %s completed", s.config.TestSymbol)

		s.logger.Info().
			Str("test", testName).
			Str("symbol", s.config.TestSymbol).
			Msg("Test passed")
	}

	s.results = append(s.results, result)
}

// TestSetMarginMode tests the set margin mode endpoint (with caution - has side effects)
func (s *AccountTestSuite) TestSetMarginMode(t *testing.T) {
	testName := "set_margin_mode"
	s.logger.Warn().Str("test", testName).Msg("Testing set margin mode endpoint - HAS SIDE EFFECTS")

	if !s.config.DemoTrading {
		s.logger.Error().Str("test", testName).Msg("Skipping margin mode test - not in demo mode")
		t.Skip("Margin mode test skipped - not in demo trading mode")
		return
	}

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/set-margin-mode",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	service := account.NewSetMarginModeService(s.client).
		Symbol(s.config.TestSymbol).
		ProductType(futures.ProductType(s.config.TestProductType)).
		MarginMode(futures.MarginModeIsolated)

	response, err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		s.logger.Warn().Msg("Margin mode test failed - this may be expected in demo mode")
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("Set margin mode for %s: %s",
			response.Symbol, response.MarginMode)

		s.logger.Info().
			Str("test", testName).
			Str("symbol", response.Symbol).
			Str("margin_mode", response.MarginMode).
			Msg("Test passed")
	}

	s.results = append(s.results, result)
}

// TestSetPositionMode tests the set position mode endpoint (with caution - has side effects)
func (s *AccountTestSuite) TestSetPositionMode(t *testing.T) {
	testName := "set_position_mode"
	s.logger.Warn().Str("test", testName).Msg("Testing set position mode endpoint - HAS SIDE EFFECTS")

	if !s.config.DemoTrading {
		s.logger.Error().Str("test", testName).Msg("Skipping position mode test - not in demo mode")
		t.Skip("Position mode test skipped - not in demo trading mode")
		return
	}

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/set-position-mode",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	service := account.NewSetPositionModeService(s.client).
		ProductType(futures.ProductType(s.config.TestProductType)).
		PositionMode(futures.PositionModeOneWay)

	response, err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		s.logger.Warn().Msg("Position mode test failed - this may be expected in demo mode")
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("Set position mode: %s", response.PositionMode)

		s.logger.Info().
			Str("test", testName).
			Str("product_type", response.ProductType).
			Str("position_mode", response.PositionMode).
			Msg("Test passed")
	}

	s.results = append(s.results, result)
}

// TestAdjustMargin tests the adjust margin endpoint (with caution - has side effects)
func (s *AccountTestSuite) TestAdjustMargin(t *testing.T) {
	testName := "adjust_margin"
	s.logger.Warn().Str("test", testName).Msg("Testing adjust margin endpoint - HAS SIDE EFFECTS")

	if !s.config.DemoTrading {
		s.logger.Error().Str("test", testName).Msg("Skipping adjust margin test - not in demo mode")
		t.Skip("Adjust margin test skipped - not in demo trading mode")
		return
	}

	startTime := time.Now()
	result := integration.TestResult{
		TestName:  testName,
		Endpoint:  "/api/v2/mix/account/set-margin",
		Timestamp: startTime,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(s.config.TestTimeout)*time.Second)
	defer cancel()

	// Use minimal amount for testing
	service := account.NewAdjustMarginService(s.client).
		Symbol(s.config.TestSymbol).
		ProductType(futures.ProductType(s.config.TestProductType)).
		MarginCoin(s.config.TestCoin).
		Amount("1").
		HoldSide(string(account.HoldSideLong)).
		Type("ADD")

	response, err := service.Do(ctx)
	result.Duration = time.Since(startTime)

	if err != nil {
		result.Success = false
		result.Error = err.Error()
		s.logger.Error().Err(err).Str("test", testName).Msg("Test failed")
		s.logger.Warn().Msg("Adjust margin test failed - this may be expected without open positions")
	} else {
		result.Success = true
		result.Response = fmt.Sprintf("Adjusted margin for %s", response.Symbol)

		s.logger.Info().
			Str("test", testName).
			Str("symbol", response.Symbol).
			Str("margin_coin", response.MarginCoin).
			Msg("Test passed")
	}

	s.results = append(s.results, result)
}

// GenerateReport creates a comprehensive test report
func (s *AccountTestSuite) GenerateReport() error {
	endTime := time.Now()

	passedTests := 0
	failedTests := 0
	for _, result := range s.results {
		if result.Success {
			passedTests++
		} else {
			failedTests++
		}
	}

	report := integration.TestReport{
		Config:      *s.config,
		StartTime:   s.startTime,
		EndTime:     endTime,
		Duration:    endTime.Sub(s.startTime),
		TotalTests:  len(s.results),
		PassedTests: passedTests,
		FailedTests: failedTests,
		TestResults: s.results,
		Summary:     fmt.Sprintf("Account Integration Tests: %d passed, %d failed out of %d total", passedTests, failedTests, len(s.results)),
		EnvironmentInfo: map[string]string{
			"demo_trading": strconv.FormatBool(s.config.DemoTrading),
			"base_url":     s.config.BaseURL,
			"test_symbol":  s.config.TestSymbol,
			"test_coin":    s.config.TestCoin,
		},
	}

	return report.SaveToFile(s.config.ReportPath)
}

// PrintSummary prints a test execution summary
func (s *AccountTestSuite) PrintSummary() {
	passedTests := 0
	failedTests := 0
	totalDuration := time.Duration(0)

	for _, result := range s.results {
		if result.Success {
			passedTests++
		} else {
			failedTests++
		}
		totalDuration += result.Duration
	}

	s.logger.Info().
		Int("total_tests", len(s.results)).
		Int("passed_tests", passedTests).
		Int("failed_tests", failedTests).
		Dur("total_duration", totalDuration).
		Msg("Account integration test summary")

	if failedTests > 0 {
		s.logger.Warn().Msg("Some tests failed - check logs and configuration")
	} else {
		s.logger.Info().Msg("All enabled tests passed successfully!")
	}
}
