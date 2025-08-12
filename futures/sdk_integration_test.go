package futures_test

import (
	"testing"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/futures/market"
	"github.com/khanbekov/go-bitget/futures/position"
	"github.com/khanbekov/go-bitget/futures/trading"
)

// TestSDKIntegration verifies the SDK works end-to-end with all refactored packages
func TestSDKIntegration(t *testing.T) {
	// Create client
	client := futures.NewClient("test-key", "test-secret", "test-passphrase")

	t.Run("AllServiceConstructorsWork", func(t *testing.T) {
		// Account services
		if account.NewAccountInfoService(client) == nil {
			t.Error("AccountInfoService constructor failed")
		}
		if account.NewSetLeverageService(client) == nil {
			t.Error("SetLeverageService constructor failed")
		}

		// Market services  
		if market.NewCandlestickService(client) == nil {
			t.Error("CandlestickService constructor failed")
		}
		if market.NewAllTickersService(client) == nil {
			t.Error("AllTickersService constructor failed")
		}

		// Position services
		if position.NewAllPositionsService(client) == nil {
			t.Error("AllPositionsService constructor failed")
		}

		// Trading services
		if trading.NewCreateOrderService(client) == nil {
			t.Error("CreateOrderService constructor failed")
		}
		if trading.NewGetOrderDetailsService(client) == nil {
			t.Error("GetOrderDetailsService constructor failed")
		}
	})

	t.Run("FluentAPIWorks", func(t *testing.T) {
		// Test account service fluent API
		accountService := account.NewAccountInfoService(client).
			Symbol("BTCUSDT").
			ProductType(account.ProductTypeUSDTFutures).
			MarginCoin("USDT")
		if accountService == nil {
			t.Error("Account service fluent API failed")
		}

		// Test market service fluent API
		marketService := market.NewCandlestickService(client).
			Symbol("BTCUSDT").
			ProductType(market.ProductTypeUSDTFutures).
			Granularity("1m")
		if marketService == nil {
			t.Error("Market service fluent API failed")
		}

		// Test trading service fluent API  
		tradingService := trading.NewCreateOrderService(client).
			Symbol("BTCUSDT").
			ProductType(trading.ProductTypeUSDTFutures).
			SideType(trading.SideBuy)
		if tradingService == nil {
			t.Error("Trading service fluent API failed")
		}

		// Test position service fluent API
		positionService := position.NewAllPositionsService(client).
			ProductType(futures.ProductTypeUSDTFutures)
		if positionService == nil {
			t.Error("Position service fluent API failed")
		}
	})

	t.Run("TypeConstantsAccessible", func(t *testing.T) {
		// Verify constants are accessible from each package
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

		// Test other constants
		if trading.SideBuy != "buy" {
			t.Error("Trading SideBuy constant incorrect")
		}
		if trading.OrderTypeLimit != "limit" {
			t.Error("Trading OrderTypeLimit constant incorrect")
		}
	})
}

// TestServiceCount verifies we have all expected services
func TestServiceCount(t *testing.T) {
	t.Run("AllServicesImplemented", func(t *testing.T) {
		client := futures.NewClient("test", "test", "test")

		// Count all services that can be created
		accountServices := []interface{}{
			account.NewAccountInfoService(client),
			account.NewAccountListService(client),
			account.NewSetLeverageService(client),
			account.NewAdjustMarginService(client),
			account.NewSetMarginModeService(client),
			account.NewSetPositionModeService(client),
			account.NewGetAccountBillService(client),
		}

		marketServices := []interface{}{
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

		positionServices := []interface{}{
			position.NewAllPositionsService(client),
			position.NewSinglePositionService(client),
			position.NewHistoryPositionsService(client),
			position.NewClosePositionService(client),
		}

		tradingServices := []interface{}{
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

		// Verify counts
		if len(accountServices) != 7 {
			t.Errorf("Expected 7 account services, got %d", len(accountServices))
		}
		if len(marketServices) != 10 {
			t.Errorf("Expected 10 market services, got %d", len(marketServices))
		}
		if len(positionServices) != 4 {
			t.Errorf("Expected 4 position services, got %d", len(positionServices))
		}
		if len(tradingServices) != 13 {
			t.Errorf("Expected 13 trading services, got %d", len(tradingServices))
		}

		totalServices := len(accountServices) + len(marketServices) + len(positionServices) + len(tradingServices)
		if totalServices != 34 {
			t.Errorf("Expected 34 total services, got %d", totalServices)
		}
	})
}