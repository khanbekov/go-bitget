package uta

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestIntegration contains integration tests that require real API credentials
// To run these tests, set up your .env file with:
// BITGET_API_KEY=your_api_key
// BITGET_SECRET_KEY=your_secret_key
// BITGET_PASSPHRASE=your_passphrase
//
// Run with: go test -v ./uta -run TestIntegration

func setupIntegrationClient(t *testing.T) *Client {
	// Load environment variables
	if err := godotenv.Load("../.env"); err != nil {
		t.Skip("No .env file found, skipping integration tests")
	}

	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		t.Skip("Missing API credentials in environment variables, skipping integration tests")
	}

	// Create logger for integration tests
	logger := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("service", "uta-integration-test").
		Logger().
		Level(zerolog.InfoLevel)

	client := NewClientWithLogger(apiKey, secretKey, passphrase, logger)
	return client
}

func TestIntegration_GetTickers_SPOT(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing GetTickers for SPOT category...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting all SPOT tickers
	tickers, err := client.NewGetTickersService().
		Category(CategorySpot).
		Do(ctx)

	require.NoError(t, err, "Failed to get SPOT tickers")
	require.NotEmpty(t, tickers, "No SPOT tickers returned")

	t.Logf("Retrieved %d SPOT tickers", len(tickers))

	// Validate ticker structure
	ticker := tickers[0]
	assert.NotEmpty(t, ticker.Symbol, "Ticker symbol should not be empty")
	assert.Equal(t, CategorySpot, ticker.Category, "Category should be SPOT")
	assert.NotEmpty(t, ticker.LastPrice, "Last price should not be empty")
	assert.NotEmpty(t, ticker.Timestamp, "Timestamp should not be empty")

	t.Logf("Sample ticker: %s - Price: %s, Volume: %s",
		ticker.Symbol, ticker.LastPrice, ticker.Volume24h)
}

func TestIntegration_GetTickers_USDT_FUTURES(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing GetTickers for USDT-FUTURES category...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting USDT-FUTURES tickers
	tickers, err := client.NewGetTickersService().
		Category(CategoryUSDTFutures).
		Do(ctx)

	require.NoError(t, err, "Failed to get USDT-FUTURES tickers")
	require.NotEmpty(t, tickers, "No USDT-FUTURES tickers returned")

	t.Logf("Retrieved %d USDT-FUTURES tickers", len(tickers))

	// Find BTC ticker for validation
	var btcTicker *Ticker
	for _, ticker := range tickers {
		if ticker.Symbol == "BTCUSDT" {
			btcTicker = &ticker
			break
		}
	}

	require.NotNil(t, btcTicker, "BTCUSDT ticker not found")
	assert.Equal(t, CategoryUSDTFutures, btcTicker.Category)
	assert.NotEmpty(t, btcTicker.LastPrice, "BTC last price should not be empty")
	assert.NotEmpty(t, btcTicker.IndexPrice, "BTC index price should not be empty")
	assert.NotEmpty(t, btcTicker.MarkPrice, "BTC mark price should not be empty")
	assert.NotEmpty(t, btcTicker.FundingRate, "BTC funding rate should not be empty")
	assert.NotEmpty(t, btcTicker.OpenInterest, "BTC open interest should not be empty")

	t.Logf("BTC Futures - Price: %s, Index: %s, Mark: %s, Funding: %s, OI: %s",
		btcTicker.LastPrice, btcTicker.IndexPrice, btcTicker.MarkPrice,
		btcTicker.FundingRate, btcTicker.OpenInterest)
}

func TestIntegration_GetTickers_SingleSymbol(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing GetTickers for single symbol...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting single ticker
	tickers, err := client.NewGetTickersService().
		Category(CategoryUSDTFutures).
		Symbol("BTCUSDT").
		Do(ctx)

	require.NoError(t, err, "Failed to get BTCUSDT ticker")
	require.Len(t, tickers, 1, "Should return exactly one ticker")

	ticker := tickers[0]
	assert.Equal(t, "BTCUSDT", ticker.Symbol)
	assert.Equal(t, CategoryUSDTFutures, ticker.Category)
	assert.NotEmpty(t, ticker.LastPrice)

	t.Logf("BTCUSDT ticker - Price: %s, 24h Change: %s%%",
		ticker.LastPrice, ticker.Price24hPcnt)
}

func TestIntegration_GetCandlesticks_USDT_FUTURES(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing GetCandlesticks for USDT-FUTURES...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting candlestick data
	candlesticks, err := client.NewGetCandlesticksService().
		Category(CategoryUSDTFutures).
		Symbol("BTCUSDT").
		Interval(Interval1m).
		Type(CandlestickTypeMarket).
		Limit("10").
		Do(ctx)

	require.NoError(t, err, "Failed to get candlesticks")
	require.NotEmpty(t, candlesticks, "No candlesticks returned")
	require.LessOrEqual(t, len(candlesticks), 10, "Should return at most 10 candlesticks")

	t.Logf("Retrieved %d candlesticks", len(candlesticks))

	// Validate candlestick structure
	candle := candlesticks[0]
	assert.NotEmpty(t, candle.Timestamp, "Timestamp should not be empty")
	assert.NotEmpty(t, candle.Open, "Open price should not be empty")
	assert.NotEmpty(t, candle.High, "High price should not be empty")
	assert.NotEmpty(t, candle.Low, "Low price should not be empty")
	assert.NotEmpty(t, candle.Close, "Close price should not be empty")
	assert.NotEmpty(t, candle.Volume, "Volume should not be empty")
	assert.NotEmpty(t, candle.Turnover, "Turnover should not be empty")

	t.Logf("Latest candle - OHLC: %s/%s/%s/%s, Volume: %s",
		candle.Open, candle.High, candle.Low, candle.Close, candle.Volume)
}

func TestIntegration_GetCandlesticks_SPOT(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing GetCandlesticks for SPOT...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting spot candlestick data
	candlesticks, err := client.NewGetCandlesticksService().
		Category(CategorySpot).
		Symbol("BTCUSDT").
		Interval(Interval5m).
		Limit("5").
		Do(ctx)

	require.NoError(t, err, "Failed to get SPOT candlesticks")
	require.NotEmpty(t, candlesticks, "No SPOT candlesticks returned")
	require.LessOrEqual(t, len(candlesticks), 5, "Should return at most 5 candlesticks")

	t.Logf("Retrieved %d SPOT candlesticks", len(candlesticks))

	// Validate first candlestick
	candle := candlesticks[0]
	assert.NotEmpty(t, candle.Timestamp)
	assert.NotEmpty(t, candle.Open)
	assert.NotEmpty(t, candle.Close)

	t.Logf("SPOT candle - Open: %s, Close: %s", candle.Open, candle.Close)
}

func TestIntegration_AccountInfo(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing AccountInfo...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting account information
	accountInfo, err := client.NewAccountInfoService().Do(ctx)

	require.NoError(t, err, "Failed to get account info")
	require.NotNil(t, accountInfo, "Account info should not be nil")

	assert.NotEmpty(t, accountInfo.AssetMode, "Asset mode should not be empty")
	assert.NotEmpty(t, accountInfo.HoldingMode, "Holding mode should not be empty")

	t.Logf("Account - Asset Mode: %s, Holding Mode: %s, STP Mode: %s",
		accountInfo.AssetMode, accountInfo.HoldingMode, accountInfo.STPMode)

	if len(accountInfo.SymbolConfig) > 0 {
		t.Logf("Symbol configs: %d", len(accountInfo.SymbolConfig))
		for i, config := range accountInfo.SymbolConfig {
			if i < 3 { // Show first 3
				t.Logf("  Symbol: %s, Category: %s, Leverage: %s",
					config.Symbol, config.Category, config.Leverage)
			}
		}
	}

	if len(accountInfo.CoinConfig) > 0 {
		t.Logf("Coin configs: %d", len(accountInfo.CoinConfig))
		for i, config := range accountInfo.CoinConfig {
			if i < 3 { // Show first 3
				t.Logf("  Coin: %s, Category: %s, Leverage: %s",
					config.Coin, config.Category, config.Leverage)
			}
		}
	}
}

func TestIntegration_AccountAssets(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing AccountAssets...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting account assets
	accountAssets, err := client.NewAccountAssetsService().Do(ctx)

	require.NoError(t, err, "Failed to get account assets")
	require.NotNil(t, accountAssets, "Account assets should not be nil")

	t.Logf("Account Equity: %s, Unrealized PNL: %s, Effective Equity: %s",
		accountAssets.AccountEquity, accountAssets.UnrealizedPNL, accountAssets.EffectiveEquity)

	if len(accountAssets.Assets) > 0 {
		t.Logf("Found %d assets", len(accountAssets.Assets))
		for i, asset := range accountAssets.Assets {
			if i < 5 { // Show first 5 non-zero assets
				if asset.Balance != "0" && asset.Balance != "" {
					t.Logf("  %s: Balance=%s, Available=%s, Frozen=%s",
						asset.Coin, asset.Balance, asset.Available, asset.Frozen)
				}
			}
		}
	}
}

func TestIntegration_FundingAssets(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing FundingAssets...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting funding assets
	fundingAssets, err := client.NewAccountFundingAssetsService().Do(ctx)

	require.NoError(t, err, "Failed to get funding assets")
	require.NotNil(t, fundingAssets, "Funding assets should not be nil")

	t.Logf("Found %d funding assets", len(fundingAssets))

	// Show non-zero funding assets
	for i, asset := range fundingAssets {
		if i < 5 && asset.Balance != "0" && asset.Balance != "" {
			t.Logf("  %s: Balance=%s, Available=%s, Frozen=%s",
				asset.Coin, asset.Balance, asset.Available, asset.Frozen)
		}
	}
}

func TestIntegration_FeeRate(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing FeeRate...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test getting fee rate for USDT futures
	feeRate, err := client.NewAccountFeeRateService().
		Symbol("BTCUSDT").
		Category(CategoryUSDTFutures).
		Do(ctx)

	require.NoError(t, err, "Failed to get fee rate")
	require.NotNil(t, feeRate, "Fee rate should not be nil")

	assert.Equal(t, "BTCUSDT", feeRate.Symbol)
	assert.Equal(t, CategoryUSDTFutures, feeRate.Category)
	assert.NotEmpty(t, feeRate.MakerRate, "Maker rate should not be empty")
	assert.NotEmpty(t, feeRate.TakerRate, "Taker rate should not be empty")

	t.Logf("BTCUSDT fee rates - Maker: %s, Taker: %s",
		feeRate.MakerRate, feeRate.TakerRate)

	// Test getting fee rate for SPOT
	spotFeeRate, err := client.NewAccountFeeRateService().
		Symbol("BTCUSDT").
		Category(CategorySpot).
		Do(ctx)

	require.NoError(t, err, "Failed to get SPOT fee rate")
	require.NotNil(t, spotFeeRate, "SPOT fee rate should not be nil")

	t.Logf("BTCUSDT SPOT fee rates - Maker: %s, Taker: %s",
		spotFeeRate.MakerRate, spotFeeRate.TakerRate)
}

func TestIntegration_ErrorHandling(t *testing.T) {
	client := setupIntegrationClient(t)

	t.Log("Testing error handling with invalid requests...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Test with invalid symbol
	_, err := client.NewGetTickersService().
		Category(CategoryUSDTFutures).
		Symbol("INVALIDSYMBOL").
		Do(ctx)

	// This might or might not error depending on the API behavior
	// Some APIs return empty results, others return errors
	if err != nil {
		t.Logf("Expected error for invalid symbol: %v", err)
	} else {
		t.Log("API returned empty result for invalid symbol (this is also valid)")
	}

	// Test with invalid category in fee rate
	_, err = client.NewAccountFeeRateService().
		Symbol("BTCUSDT").
		Category("INVALID_CATEGORY").
		Do(ctx)

	// This should typically return an error
	if err != nil {
		t.Logf("Expected error for invalid category: %v", err)
	}
}

// Benchmark tests for performance measurement
func BenchmarkIntegration_GetTickers(b *testing.B) {
	if os.Getenv("BITGET_API_KEY") == "" {
		b.Skip("No API credentials, skipping benchmark")
	}

	client := setupIntegrationClient(&testing.T{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.NewGetTickersService().
			Category(CategoryUSDTFutures).
			Do(ctx)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}

func BenchmarkIntegration_GetSingleTicker(b *testing.B) {
	if os.Getenv("BITGET_API_KEY") == "" {
		b.Skip("No API credentials, skipping benchmark")
	}

	client := setupIntegrationClient(&testing.T{})
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.NewGetTickersService().
			Category(CategoryUSDTFutures).
			Symbol("BTCUSDT").
			Do(ctx)
		if err != nil {
			b.Fatalf("Benchmark failed: %v", err)
		}
	}
}
