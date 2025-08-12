package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/khanbekov/go-bitget/common"
	"github.com/khanbekov/go-bitget/uta"
	"github.com/khanbekov/go-bitget/ws"
	"github.com/rs/zerolog"
	"os"
	"strings"
	"time"
)

func main() {
	logger, apiKey, secretKey, passphrase := initialize()

	client := uta.NewClientWithLogger(apiKey, secretKey, passphrase, logger)

	client.SetDemoTrading(false)

	// Test account functionality
	testAccount(client, logger)

	// Test order functionality (careful with real orders)
	testOrderCreateModifyCancel(client, logger)

	// Test UTA-specific features
	testUTAFeatures(client, logger)
}

func testCandlestick(client *uta.Client, logger zerolog.Logger) {
	s := client.NewGetCandlesticksService().
		Symbol("BTCUSDT").
		Category(uta.CategoryUSDTFutures).
		Interval(uta.Interval1m).
		Limit("10")

	res, err := s.Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error from request")
	}
	logger.Info().Any("response", res).Msg("UTA Candlestick data")
}

func testAccount(client *uta.Client, logger zerolog.Logger) {
	// Test account info
	accountInfo, err := client.NewAccountInfoService().Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error getting account info")
	} else {
		logger.Info().Any("account_info", accountInfo).Msg("UTA Account Info")
	}

	// Test account assets
	assets, err := client.NewAccountAssetsService().Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error getting account assets")
	} else {
		logger.Info().Any("account_assets", assets).Msg("UTA Account Assets")
	}

	// Test fee rate (may not be available in demo mode)
	feeRate, err := client.NewAccountFeeRateService().
		Symbol("BTCUSDT").
		Category(uta.CategoryUSDTFutures).
		Do(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "40404") {
			logger.Warn().Msg("Fee rate endpoint not available in demo mode (this is normal)")
		} else {
			logger.Error().Err(err).Msg("Error getting fee rate")
		}
	} else {
		logger.Info().Any("fee_rate", feeRate).Msg("UTA Fee Rate")
	}
}

func testOrderCreateModifyCancel(client *uta.Client, logger zerolog.Logger) {
	// Note: UTA order operations require careful consideration in production
	// This example shows the complete create-modify-cancel workflow

	logger.Info().Msg("Starting UTA order management test (Create → Modify → Cancel)")

	// Step 1: Get current market price
	tickers, err := client.NewGetTickersService().
		Category(uta.CategoryUSDTFutures).
		Symbol("BTCUSDT").
		Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get ticker")
		return
	}

	var currentPrice float64 = 30000 // fallback price
	if len(tickers) > 0 {
		logger.Info().Str("current_price", tickers[0].LastPrice).Msg("Current BTC price")
		// Parse current price for calculations (simplified)
		currentPrice = 50000 // using safe fallback
	}

	// Step 2: Create initial order with very safe price
	safePrice := fmt.Sprintf("%.2f", currentPrice*0.3) // 30% below market price
	logger.Info().Str("safe_price", safePrice).Msg("Creating order with safe price")

	order, err := testCreateLimitOrder(client, logger, safePrice)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create order")
		return
	}

	if order == nil {
		logger.Error().Msg("Order creation returned nil")
		return
	}

	logger.Info().
		Str("order_id", order.OrderID).
		Str("initial_price", safePrice).
		Str("size", "0.001").
		Msg("✅ Step 1: Order created successfully")

	// Step 3: Modify the order (change price and size)
	modifiedPrice := fmt.Sprintf("%.2f", currentPrice*0.25) // Even lower price
	modifiedSize := "0.0015"                                // Slightly larger size

	logger.Info().
		Str("new_price", modifiedPrice).
		Str("new_size", modifiedSize).
		Msg("Modifying order with new parameters")

	_, err = testModifyOrder(client, logger, order.OrderID, modifiedPrice, modifiedSize)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to modify order")
		// Continue with cancellation even if modify failed
	} else {
		logger.Info().
			Str("order_id", order.OrderID).
			Str("modified_price", modifiedPrice).
			Str("modified_size", modifiedSize).
			Msg("✅ Step 2: Order modified successfully")
	}

	// Step 4: Cancel the order for safety
	logger.Info().Str("order_id", order.OrderID).Msg("Cancelling order for safety")

	_, err = testCancelOrder(client, logger, order.OrderID)
	if err != nil {
		logger.Error().Err(err).Msg("❌ Failed to cancel order - MANUAL INTERVENTION MAY BE REQUIRED")
		logger.Error().Str("order_id", order.OrderID).Msg("Please manually cancel this order in your Bitget account")
	} else {
		logger.Info().
			Str("order_id", order.OrderID).
			Msg("✅ Step 3: Order cancelled successfully")
	}

	logger.Info().Msg("🎉 Complete Create → Modify → Cancel workflow finished")
}

func testCancelAllOrders(client *uta.Client, logger zerolog.Logger) {
	// Note: This is a stub implementation in UTA SDK
	logger.Info().Msg("Cancel all orders functionality available in UTA SDK (stub implementation)")

	// In a real implementation, this would be:
	// err := client.NewCancelAllOrdersService().
	//     Category(uta.CategoryUSDTFutures).
	//     Symbol("BTCUSDT").
	//     Do(context.Background())
}

func initialize() (zerolog.Logger, string, string, string) {
	// Человекочитаемые консольные логи с цветами
	var writer zerolog.ConsoleWriter = zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.StampMicro,
		NoColor:    false,
	}

	logger := zerolog.New(writer).
		With().
		Timestamp().
		Logger()

	err := godotenv.Load() // Загружает .env файл в текущей директории
	if err != nil {
		logger.Warn().Msg(".env file not found")
	}

	// Получаем переменные окружения
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	// Проверяем наличие обязательных переменных
	if apiKey == "" || secretKey == "" || passphrase == "" {
		logger.Error().Msg("API keys or passphrase not set in environment variables")
		os.Exit(1)
	}
	return logger, apiKey, secretKey, passphrase
}

func testCreateLimitOrder(client *uta.Client, logger zerolog.Logger, price string) (*uta.Order, error) {
	clientOid := uuid.New().String()

	s := client.NewPlaceOrderService().
		Category(uta.CategoryUSDTFutures).
		Symbol("BTCUSDT").
		Side(uta.SideBuy).
		OrderType(uta.OrderTypeLimit).
		Size("0.001").
		Price(price). // Using provided safe price
		ClientOid(clientOid).
		TimeInForce(uta.TimeInForceGTC).
		PositionSide(uta.PositionSideLong) // Required for futures trading

	res, err := s.Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error creating UTA order")

		// Provide specific help for common errors
		if strings.Contains(err.Error(), "40099") {
			logger.Error().Msg("🔧 Solution: Check if your API keys match the environment (testnet vs production)")
			logger.Error().Msg("   - Production keys: https://www.bitget.com/")
			logger.Error().Msg("   - Testnet keys: https://testnet.bitget.com/")
			logger.Error().Msg("   - Add BITGET_TESTNET=true to .env if using testnet")
		} else if strings.Contains(err.Error(), "40084") {
			logger.Error().Msg("🔧 Solution: Your account is in Classic mode, switch to UTA mode")
			logger.Error().Msg("   - Visit Bitget website → Account Settings → Switch to UTA mode")
		} else if strings.Contains(err.Error(), "25236") {
			logger.Error().Msg("🔧 Solution: Missing required parameters for futures trading")
			logger.Error().Msg("   - positionSide (long/short) is required for futures orders")
			logger.Error().Msg("   - This has been fixed in the code - the error should not occur again")
		}

		return nil, err
	}
	logger.Info().Any("response", res).Msg("UTA order created")
	return res, nil
}

func testModifyOrder(client *uta.Client, logger zerolog.Logger, orderID string, newPrice string, newSize string) (*uta.Order, error) {
	logger.Info().
		Str("order_id", orderID).
		Str("new_price", newPrice).
		Str("new_size", newSize).
		Msg("Modifying order with new parameters")

	s := client.NewModifyOrderService().
		Category(uta.CategoryUSDTFutures).
		Symbol("BTCUSDT").
		OrderId(orderID).
		NewPrice(newPrice).
		NewSize(newSize)

	res, err := s.Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error modifying UTA order")

		// Provide specific help for common errors
		if strings.Contains(err.Error(), "40007") {
			logger.Error().Msg("🔧 Solution: Order not found or already filled/cancelled")
		} else if strings.Contains(err.Error(), "40008") {
			logger.Error().Msg("🔧 Solution: Order cannot be modified (wrong status)")
		} else if strings.Contains(err.Error(), "40084") {
			logger.Error().Msg("🔧 Solution: Your account is in Classic mode, switch to UTA mode")
		}
		return nil, err
	}

	logger.Info().Any("response", res).Msg("UTA order modified")
	return res, nil
}

func testCancelOrder(client *uta.Client, logger zerolog.Logger, orderId string) (*uta.Order, error) {
	s := client.NewCancelOrderService().
		Category(uta.CategoryUSDTFutures).
		Symbol("BTCUSDT").
		OrderId(orderId)

	res, err := s.Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error cancelling UTA order")
		return nil, err
	}
	logger.Info().Any("response", res).Msg("UTA order cancelled")
	return res, nil
}

func testTickers(client *uta.Client, logger zerolog.Logger) {
	logger.Info().Msg("Testing UTA Tickers")

	// Test USDT futures tickers
	futuresTickers, err := client.NewGetTickersService().
		Category(uta.CategoryUSDTFutures).
		Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error getting USDT futures tickers")
	} else {
		logger.Info().Int("count", len(futuresTickers)).Msg("USDT Futures tickers retrieved")
		if len(futuresTickers) > 0 {
			// Show BTC ticker
			for _, ticker := range futuresTickers {
				if ticker.Symbol == "BTCUSDT" {
					logger.Info().
						Str("symbol", ticker.Symbol).
						Str("last_price", ticker.LastPrice).
						Str("24h_change", ticker.Price24hPcnt).
						Str("volume", ticker.Volume24h).
						Str("funding_rate", ticker.FundingRate).
						Str("open_interest", ticker.OpenInterest).
						Msg("BTC Futures Ticker")
					break
				}
			}
		}
	}

	// Test spot tickers
	spotTickers, err := client.NewGetTickersService().
		Category(uta.CategorySpot).
		Symbol("BTCUSDT").
		Do(context.Background())
	if err != nil {
		logger.Error().Err(err).Msg("Error getting SPOT ticker")
	} else if len(spotTickers) > 0 {
		ticker := spotTickers[0]
		logger.Info().
			Str("symbol", ticker.Symbol).
			Str("last_price", ticker.LastPrice).
			Str("24h_change", ticker.Price24hPcnt).
			Str("volume", ticker.Volume24h).
			Msg("BTC Spot Ticker")
	}
}

func testUTAFeatures(client *uta.Client, logger zerolog.Logger) {
	logger.Info().Msg("Testing UTA-specific features")

	// Test funding assets (may not be available in demo mode)
	fundingAssets, err := client.NewAccountFundingAssetsService().
		Coin("USDT").
		Do(context.Background())
	if err != nil {
		if strings.Contains(err.Error(), "40404") {
			logger.Warn().Msg("Funding assets endpoint not available in demo mode (this is normal)")
		} else {
			logger.Error().Err(err).Msg("Error getting funding assets")
		}
	} else {
		logger.Info().Any("funding_assets", fundingAssets).Msg("USDT Funding Assets")
	}

	// Test multiple category tickers
	categories := []string{
		uta.CategorySpot,
		uta.CategoryUSDTFutures,
		uta.CategoryMargin,
	}

	for _, category := range categories {
		tickers, err := client.NewGetTickersService().
			Category(category).
			Symbol("BTCUSDT").
			Do(context.Background())
		if err != nil {
			logger.Error().Err(err).Str("category", category).Msg("Error getting ticker for category")
		} else if len(tickers) > 0 {
			logger.Info().
				Str("category", category).
				Str("price", tickers[0].LastPrice).
				Msg("BTC price in category")
		}
	}

	// Test candlesticks for different intervals
	intervals := []string{uta.Interval1m, uta.Interval5m, uta.Interval1H}
	for _, interval := range intervals {
		candles, err := client.NewGetCandlesticksService().
			Category(uta.CategoryUSDTFutures).
			Symbol("BTCUSDT").
			Interval(interval).
			Limit("3").
			Do(context.Background())
		if err != nil {
			logger.Error().Err(err).Str("interval", interval).Msg("Error getting candlesticks")
		} else {
			logger.Info().
				Str("interval", interval).
				Int("count", len(candles)).
				Msg("Candlesticks retrieved")
		}
	}
}

func testPublicWS(secret string, logger zerolog.Logger) {
	public_ws_url := "wss://ws.bitget.com/v2/ws/public"
	client := ws.NewBitgetBaseWsClient(logger, public_ws_url, secret)
	client.SetListener(func(message string) { logger.Info().Msg(message) }, func(message string) { logger.Error().Msg(message) })

	client.ConnectWebSocket()
	go client.Connect()
	client.SendByType(ws.WsBaseReq{
		Op: common.WsOpSubscribe,
		Args: []interface{}{
			ws.SubscriptionArgs{
				ProductType: "USDT-FUTURES",
				Channel:     "ticker",
				Symbol:      "BTCUSDT",
			},
		},
	})
	client.StartReadLoop()
	time.Sleep(100 * time.Second)
}

func testPrivateWS(apiKey, secret, passphrase string, logger zerolog.Logger) {
	public_ws_url := "wss://ws.bitget.com/v2/ws/private"
	client := ws.NewBitgetBaseWsClient(logger, public_ws_url, secret)
	client.SetListener(func(message string) { logger.Info().Msg(message) }, func(message string) { logger.Error().Msg(message) })

	client.ConnectWebSocket()
	go client.Connect()
	defer client.Close()
	client.Login(apiKey, passphrase, common.SHA256)

	time.Sleep(3 * time.Second)

	client.SendByType(ws.WsBaseReq{
		Op: common.WsOpSubscribe,
		Args: []interface{}{
			ws.SubscriptionArgs{
				ProductType: "USDT-FUTURES",
				Channel:     "orders",
				Symbol:      "default",
			},
		},
	})
	client.StartReadLoop()
	time.Sleep(100 * time.Second)
}
