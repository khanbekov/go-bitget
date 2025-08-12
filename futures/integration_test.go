package futures_test

import (
	"context"
	"testing"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/futures/market"
	"github.com/khanbekov/go-bitget/futures/position"
	"github.com/khanbekov/go-bitget/futures/trading"
)

// TestServiceConstructors verifies that all service constructor functions work correctly
// and can be instantiated using the new package-based approach.
func TestServiceConstructors(t *testing.T) {
	// Create client
	client := futures.NewClient("test-key", "test-secret", "test-passphrase")

	t.Run("Account Services", func(t *testing.T) {
		// Test all account services can be instantiated
		services := []interface{}{
			account.NewAccountInfoService(client),
			account.NewAccountListService(client),
			account.NewSetLeverageService(client),
			account.NewAdjustMarginService(client),
			account.NewSetMarginModeService(client),
			account.NewSetPositionModeService(client),
			account.NewGetAccountBillService(client),
		}

		for i, service := range services {
			if service == nil {
				t.Errorf("Account service %d is nil", i)
			}
		}

		// Test fluent API works
		accountService := account.NewAccountInfoService(client)
		chained := accountService.Symbol("BTCUSDT").ProductType(account.ProductTypeUSDTFutures).MarginCoin("USDT")
		if chained == nil {
			t.Error("Account service fluent API returned nil")
		}
	})

	t.Run("Market Services", func(t *testing.T) {
		// Test all market services can be instantiated
		services := []interface{}{
			market.NewCandlestickService(client),
			market.NewAllTickersService(client),
			market.NewTickerService(client),
			market.NewOrderBookService(client),
			market.NewRecentTradesService(client),
			market.NewCurrentFundingRateService(client),
			market.NewHistoryFundingRateService(client),
			market.NewOpenInterestService(client),
			market.NewSymbolPriceService(client),
			market.NewContractsService(client),
		}

		for i, service := range services {
			if service == nil {
				t.Errorf("Market service %d is nil", i)
			}
		}

		// Test fluent API works
		candleService := market.NewCandlestickService(client)
		chained := candleService.Symbol("BTCUSDT").ProductType(market.ProductTypeUSDTFutures).Granularity("1m")
		if chained == nil {
			t.Error("Market service fluent API returned nil")
		}
	})

	t.Run("Position Services", func(t *testing.T) {
		// Test all position services can be instantiated
		services := []interface{}{
			position.NewAllPositionsService(client),
			position.NewSinglePositionService(client),
			position.NewHistoryPositionsService(client),
			position.NewClosePositionService(client),
		}

		for i, service := range services {
			if service == nil {
				t.Errorf("Position service %d is nil", i)
			}
		}

		// Test fluent API works
		positionService := position.NewAllPositionsService(client)
		chained := positionService.ProductType(futures.ProductTypeUSDTFutures)
		if chained == nil {
			t.Error("Position service fluent API returned nil")
		}
	})

	t.Run("Trading Services", func(t *testing.T) {
		// Test all trading services can be instantiated
		services := []interface{}{
			trading.NewCreateOrderService(client),
			trading.NewModifyOrderService(client),
			trading.NewCancelOrderService(client),
			trading.NewCancelAllOrdersService(client),
			trading.NewGetOrderDetailsService(client),
			trading.NewPendingOrdersService(client),
			trading.NewOrderHistoryService(client),
			trading.NewFillHistoryService(client),
			trading.NewCreatePlanOrderService(client),
			trading.NewModifyPlanOrderService(client),
			trading.NewCancelPlanOrderService(client),
			trading.NewPendingPlanOrdersService(client),
			trading.NewCreateBatchOrdersService(client),
		}

		for i, service := range services {
			if service == nil {
				t.Errorf("Trading service %d is nil", i)
			}
		}

		// Test fluent API works
		orderService := trading.NewCreateOrderService(client)
		chained := orderService.Symbol("BTCUSDT").ProductType(trading.ProductTypeUSDTFutures).SideType("buy")
		if chained == nil {
			t.Error("Trading service fluent API returned nil")
		}
	})

	// Test type constants are accessible
	t.Run("Type Constants", func(t *testing.T) {
		if account.ProductTypeUSDTFutures != "USDT-FUTURES" {
			t.Error("Account ProductType constant incorrect")
		}
		if market.ProductTypeUSDTFutures != "USDT-FUTURES" {
			t.Error("Market ProductType constant incorrect")
		}
		if futures.ProductTypeUSDTFutures != "USDT-FUTURES" {
			t.Error("Position ProductType constant incorrect")
		}
		if trading.ProductTypeUSDTFutures != "USDT-FUTURES" {
			t.Error("Trading ProductType constant incorrect")
		}
	})
}

// TestServiceIntegration tests that services properly integrate with the client
func TestServiceIntegration(t *testing.T) {
	client := futures.NewClient("test-key", "test-secret", "test-passphrase")

	// Test that services can use client methods
	accountService := account.NewAccountInfoService(client)
	if accountService == nil {
		t.Fatal("Account service is nil")
	}

	// Test that service has access to client interface methods
	// Note: This would normally call the API, so we just test the setup
	setupComplete := accountService.Symbol("BTCUSDT").ProductType(account.ProductTypeUSDTFutures) != nil
	if !setupComplete {
		t.Error("Service setup failed")
	}
}

// TestUsagePatterns demonstrates the correct usage patterns for the SDK
func TestUsagePatterns(t *testing.T) {
	client := futures.NewClient("test-key", "test-secret", "test-passphrase")
	ctx := context.Background()

	t.Run("Account Usage Pattern", func(t *testing.T) {
		// Demonstrate correct account service usage
		accountService := account.NewAccountInfoService(client)
		request := accountService.
			Symbol("BTCUSDT").
			ProductType(account.ProductTypeUSDTFutures).
			MarginCoin("USDT")
		
		// Would normally call Do(ctx) here, but that requires valid API credentials
		if request == nil {
			t.Error("Account service request setup failed")
		}
	})

	t.Run("Market Usage Pattern", func(t *testing.T) {
		// Demonstrate correct market service usage
		candleService := market.NewCandlestickService(client)
		request := candleService.
			Symbol("BTCUSDT").
			ProductType(market.ProductTypeUSDTFutures).
			Granularity("1m").
			Limit("100")
		
		if request == nil {
			t.Error("Market service request setup failed")
		}
	})

	t.Run("Trading Usage Pattern", func(t *testing.T) {
		// Demonstrate correct trading service usage  
		orderService := trading.NewCreateOrderService(client)
		request := orderService.
			Symbol("BTCUSDT").
			ProductType(trading.ProductTypeUSDTFutures).
			SideType(trading.SideBuy).
			OrderType(trading.OrderTypeLimit).
			Size("0.01").
			Price("45000")
			
		if request == nil {
			t.Error("Trading service request setup failed")
		}
	})

	t.Run("Position Usage Pattern", func(t *testing.T) {
		// Demonstrate correct position service usage
		positionService := position.NewAllPositionsService(client)
		request := positionService.ProductType(futures.ProductTypeUSDTFutures)
		
		if request == nil {
			t.Error("Position service request setup failed")
		}
	})

	// The ctx variable should be used (but we don't make actual API calls in unit tests)
	_ = ctx
}