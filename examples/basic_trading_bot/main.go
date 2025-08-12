package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/futures/market"
	"github.com/khanbekov/go-bitget/futures/position"
	"github.com/khanbekov/go-bitget/futures/trading"
)

// TradingBot demonstrates a complete trading bot implementation
type TradingBot struct {
	client      *futures.Client
	symbol      string
	productType string
	ctx         context.Context

	// Strategy parameters
	positionSize  string
	leverage      string
	takeProfitPct float64
	stopLossPct   float64

	// State tracking
	lastPrice float64
	isRunning bool
}

// NewTradingBot creates a new trading bot instance
func NewTradingBot() *TradingBot {
	// Initialize client with credentials from environment
	client := futures.NewClient(
		os.Getenv("BITGET_API_KEY"),
		os.Getenv("BITGET_SECRET_KEY"),
		os.Getenv("BITGET_PASSPHRASE"),
	)

	// Configure for testnet if needed
	if os.Getenv("BITGET_TESTNET") == "true" {
		client.SetApiEndpoint("https://testnet.bitget.com")
	}

	return &TradingBot{
		client:        client,
		symbol:        "BTCUSDT",
		productType:   "USDT-FUTURES",
		ctx:           context.Background(),
		positionSize:  "0.001", // Small size for demo
		leverage:      "2",     // Leverage for strategy
		takeProfitPct: 0.02,    // 2% take profit
		stopLossPct:   0.01,    // 1% stop loss
	}
}

// Start begins the trading bot operation
func (bot *TradingBot) Start() error {
	log.Println("🚀 Starting Trading Bot...")
	bot.isRunning = true

	// 1. Check account status
	if err := bot.checkAccountStatus(); err != nil {
		return fmt.Errorf("account check failed: %w", err)
	}

	// 2. Set leverage
	if err := bot.setLeverage(bot.leverage); err != nil {
		return fmt.Errorf("set leverage failed: %w", err)
	}

	// 3. Main trading loop
	ticker := time.NewTicker(5 * time.Second) // Check every 30 seconds
	defer ticker.Stop()

	for bot.isRunning {
		select {
		case <-ticker.C:
			if err := bot.executeTradingCycle(); err != nil {
				log.Printf("❌ Trading cycle error: %v", err)
				// Continue running despite errors
			}
		case <-bot.ctx.Done():
			log.Println("🛑 Bot stopping due to context cancellation")
			bot.isRunning = false
		}
	}

	return nil
}

// checkAccountStatus verifies account is ready for trading
func (bot *TradingBot) checkAccountStatus() error {
	log.Println("📊 Checking account status...")

	// Get account information
	accountService := account.NewAccountInfoService(bot.client)
	accountInfo, err := accountService.
		Symbol(bot.symbol).
		ProductType(account.ProductType(bot.productType)).
		MarginCoin("USDT").
		Do(bot.ctx)

	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}

	// Check available balance
	available := accountInfo.Available
	if err != nil {
		return fmt.Errorf("invalid available balance: %w", err)
	}

	if available < 10.0 { // Minimum 10 USDT
		return fmt.Errorf("insufficient balance: %.2f USDT (minimum 10 USDT required)", available)
	}

	log.Printf("✅ Account ready - Available: %.2f USDT", available)
	return nil
}

// executeTradingCycle runs one complete trading cycle
func (bot *TradingBot) executeTradingCycle() error {
	// 1. Get current market price
	currentPrice, err := bot.getCurrentPrice()
	if err != nil {
		return fmt.Errorf("failed to get current price: %w", err)
	}

	// 2. Check existing positions
	positions, err := bot.getCurrentPositions()
	if err != nil {
		return fmt.Errorf("failed to get positions: %w", err)
	}

	// 3. Execute strategy based on position state
	if len(positions) == 0 {
		// No position - look for entry signal
		return bot.checkEntrySignal(currentPrice)
	} else {
		// Has position - manage existing trades
		return bot.manageExistingPositions(positions, currentPrice)
	}
}

// getCurrentPrice fetches the latest market price
func (bot *TradingBot) getCurrentPrice() (float64, error) {
	tickerService := market.NewTickerService(bot.client)
	ticker, err := tickerService.
		Symbol(bot.symbol).
		ProductType(bot.productType).
		Do(bot.ctx)

	if err != nil {
		return 0, err
	}

	price, err := strconv.ParseFloat(ticker.LastPr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price format: %w", err)
	}

	return price, nil
}

// getCurrentPositions retrieves current open positions
func (bot *TradingBot) getCurrentPositions() ([]*position.Position, error) {
	positionService := position.NewAllPositionsService(bot.client)
	positions, err := positionService.
		ProductType(futures.ProductType(bot.productType)).
		Do(bot.ctx)

	if err != nil {
		return nil, err
	}

	// Filter positions for our symbol
	var filteredPositions []*position.Position
	for _, pos := range positions {
		if pos.Symbol == bot.symbol && pos.Size != 0 {
			filteredPositions = append(filteredPositions, pos)
		}
	}

	return filteredPositions, nil
}

// checkEntrySignal implements a simple momentum strategy
func (bot *TradingBot) checkEntrySignal(currentPrice float64) error {
	log.Printf("📈 Checking entry signal - Price: %.2f", currentPrice)

	increase := 1.0
	decrease := 1.0

	// Simple strategy: Buy if price increased by "increase" in last check
	if bot.lastPrice > 0 && currentPrice > bot.lastPrice*increase {
		log.Println("🟢 Entry signal detected - Placing BUY order")
		return bot.placeBuyOrder(currentPrice)
	}

	// Sell if price decreased by "decrease"
	if bot.lastPrice > 0 && currentPrice < bot.lastPrice*decrease {
		log.Println("🔴 Entry signal detected - Placing SELL order")
		return bot.placeSellOrder(currentPrice)
	}
	bot.lastPrice = currentPrice

	log.Println("⏸️ No entry signal")
	return nil
}

// setLeverage set leverage on account for both sides
func (bot *TradingBot) setLeverage(leverage string) error {
	setLeverageService := account.NewSetLeverageService(bot.client)
	err := setLeverageService.
		Leverage(leverage).
		Symbol(bot.symbol).
		ProductType(futures.ProductType(bot.productType)).
		MarginCoin("USDT").Do(bot.ctx)

	if err != nil {
		log.Printf("❌ Failed to set leverage: %v err: %v", leverage, err)
	} else {
		log.Printf("✅ Set leverage to %v", leverage)
	}
	return err
}

// placeBuyOrder places a long position order
func (bot *TradingBot) placeBuyOrder(currentPrice float64) error {
	orderService := trading.NewCreateOrderService(bot.client)

	order, err := orderService.
		Symbol(bot.symbol).
		ProductType(trading.ProductType(bot.productType)).
		MarginMode(trading.MarginModeCrossed).
		MarginCoin("USDT").
		SideType(trading.SideBuy).
		PositionSideType(trading.PositionSideOpen).
		OrderType(trading.OrderTypeMarket).
		Size(bot.positionSize).
		Do(bot.ctx)

	if err != nil {
		return fmt.Errorf("failed to place buy order: %w", err)
	}

	log.Printf("✅ BUY order placed - ID: %s", order.OrderId)
	return nil
}

// placeSellOrder places a short position order
func (bot *TradingBot) placeSellOrder(currentPrice float64) error {
	orderService := trading.NewCreateOrderService(bot.client)

	order, err := orderService.
		Symbol(bot.symbol).
		ProductType(trading.ProductType(bot.productType)).
		MarginMode(trading.MarginModeCrossed).
		MarginCoin("USDT").
		SideType(trading.SideSell).
		PositionSideType(trading.PositionSideOpen).
		OrderType(trading.OrderTypeMarket).
		Size(bot.positionSize).
		Do(bot.ctx)

	if err != nil {
		return fmt.Errorf("failed to place sell order: %w", err)
	}

	log.Printf("✅ SELL order placed - ID: %s", order.OrderId)
	return nil
}

// manageExistingPositions handles open positions with risk management
func (bot *TradingBot) manageExistingPositions(positions []*position.Position, currentPrice float64) error {
	for _, pos := range positions {
		log.Printf("📊 Managing position: %s %s @ %s", pos.HoldSide, pos.Size, pos.AverageOpenPrice)

		// Calculate P&L percentage
		entryPrice := pos.AverageOpenPrice
		var pnlPct float64

		if pos.HoldSide == "long" {
			pnlPct = (currentPrice - entryPrice) / entryPrice
		} else {
			pnlPct = (entryPrice - currentPrice) / entryPrice
		}

		log.Printf("💰 P&L: %.2f%%", pnlPct*100)

		// Check take profit
		if pnlPct >= bot.takeProfitPct {
			log.Printf("🎯 Take profit triggered - Closing position")
			return bot.closePosition(pos)
		}

		// Check stop loss
		if pnlPct <= -bot.stopLossPct {
			log.Printf("🛑 Stop loss triggered - Closing position")
			return bot.closePosition(pos)
		}
	}

	return nil
}

// closePosition closes an existing position
func (bot *TradingBot) closePosition(pos *position.Position) error {
	closeService := position.NewClosePositionService(bot.client)

	_, err := closeService.
		Symbol(pos.Symbol).
		ProductType(futures.ProductType(bot.productType)).
		HoldSide(pos.HoldSide).
		Do(bot.ctx)

	if err != nil {
		return fmt.Errorf("failed to close position: %w", err)
	}

	log.Printf("✅ Position closed successfully")
	return nil
}

// Stop gracefully stops the trading bot
func (bot *TradingBot) Stop() {
	log.Println("🛑 Stopping Trading Bot...")
	bot.isRunning = false
}

func main() {
	// Check required environment variables
	requiredEnvs := []string{"BITGET_API_KEY", "BITGET_SECRET_KEY", "BITGET_PASSPHRASE"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", env)
		}
	}

	// Create and start bot
	bot := NewTradingBot()

	// Handle graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bot.ctx = ctx

	// Start the bot
	if err := bot.Start(); err != nil {
		log.Fatalf("❌ Bot failed to start: %v", err)
	}
}
