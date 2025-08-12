package uta

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// Example test to demonstrate UTA SDK usage
func TestExample_UTA_Usage(t *testing.T) {
	// Skip if no .env file (this is just an example)
	if err := godotenv.Load("../.env"); err != nil {
		t.Skip("No .env file found, skipping example test")
	}

	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		t.Skip("Missing API credentials in environment variables")
	}

	// Create UTA client
	client := NewClient(apiKey, secretKey, passphrase)
	ctx := context.Background()

	t.Run("Get USDT Futures Tickers", func(t *testing.T) {
		tickers, err := client.NewGetTickersService().
			Category(CategoryUSDTFutures).
			Do(ctx)

		require.NoError(t, err)
		require.NotEmpty(t, tickers)
		t.Logf("Found %d USDT futures tickers", len(tickers))
	})

	t.Run("Get SPOT Tickers", func(t *testing.T) {
		tickers, err := client.NewGetTickersService().
			Category(CategorySpot).
			Do(ctx)

		require.NoError(t, err)
		require.NotEmpty(t, tickers)
		t.Logf("Found %d SPOT tickers", len(tickers))
	})

	t.Run("Get BTC Candlesticks", func(t *testing.T) {
		candlesticks, err := client.NewGetCandlesticksService().
			Category(CategoryUSDTFutures).
			Symbol("BTCUSDT").
			Interval(Interval1m).
			Limit("5").
			Do(ctx)

		require.NoError(t, err)
		require.NotEmpty(t, candlesticks)
		require.LessOrEqual(t, len(candlesticks), 5)
		t.Logf("Found %d candlesticks for BTCUSDT", len(candlesticks))
	})

	t.Run("Get Account Info", func(t *testing.T) {
		accountInfo, err := client.NewAccountInfoService().Do(ctx)

		require.NoError(t, err)
		require.NotNil(t, accountInfo)
		t.Logf("Account - Asset Mode: %s, Holding Mode: %s",
			accountInfo.AssetMode, accountInfo.HoldingMode)
	})

	t.Run("Get Account Assets", func(t *testing.T) {
		accountAssets, err := client.NewAccountAssetsService().Do(ctx)

		require.NoError(t, err)
		require.NotNil(t, accountAssets)
		t.Logf("Account Equity: %s", accountAssets.AccountEquity)
	})

	t.Run("Get Fee Rate", func(t *testing.T) {
		feeRate, err := client.NewAccountFeeRateService().
			Symbol("BTCUSDT").
			Category(CategoryUSDTFutures).
			Do(ctx)

		require.NoError(t, err)
		require.NotNil(t, feeRate)
		t.Logf("BTCUSDT USDT-FUTURES fee rates - Maker: %s, Taker: %s",
			feeRate.MakerRate, feeRate.TakerRate)
	})
}

// Example of how to use UTA SDK in production code
func ExampleClient() {
	// Initialize the client
	client := NewClient("your_api_key", "your_secret_key", "your_passphrase")

	// Get all USDT futures tickers
	tickers, err := client.NewGetTickersService().
		Category(CategoryUSDTFutures).
		Do(context.Background())
	if err != nil {
		panic(err)
	}

	for _, ticker := range tickers {
		if ticker.Symbol == "BTCUSDT" {
			println("BTC Price:", ticker.LastPrice)
			break
		}
	}

	// Get account information
	accountInfo, err := client.NewAccountInfoService().Do(context.Background())
	if err != nil {
		panic(err)
	}

	println("Account Asset Mode:", accountInfo.AssetMode)

	// Place a limit order (commented out for safety)
	/*
		order, err := client.NewPlaceOrderService().
			Symbol("BTCUSDT").
			Category(CategoryUSDTFutures).
			Side(SideBuy).
			OrderType(OrderTypeLimit).
			Size("0.001").
			Price("30000"). // Set your desired price
			Do(context.Background())
		if err != nil {
			panic(err)
		}

		println("Order placed:", order.OrderID)
	*/
}
