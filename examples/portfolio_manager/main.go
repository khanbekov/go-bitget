package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/futures/market"
	"github.com/khanbekov/go-bitget/futures/position"
	"github.com/khanbekov/go-bitget/futures/trading"
)

// PortfolioManager demonstrates advanced portfolio management
type PortfolioManager struct {
	client      *futures.Client
	ctx         context.Context
	
	// Portfolio configuration
	symbols     []string
	productType string
	
	// Risk management
	maxPositions      int
	maxRiskPerTrade   float64 // % of portfolio
	maxTotalExposure  float64 // % of portfolio
	
	// Portfolio state
	totalBalance      float64
	currentExposure   float64
	positions         map[string]*PositionInfo
}

// PositionInfo holds detailed position information
type PositionInfo struct {
	Symbol        string
	Side          string
	Size          float64
	EntryPrice    float64
	CurrentPrice  float64
	UnrealizedPnL float64
	MarginUsed    float64
	RiskLevel     string
}

// PortfolioSummary contains portfolio analytics
type PortfolioSummary struct {
	TotalBalance     float64
	AvailableBalance float64
	TotalExposure    float64
	UnrealizedPnL    float64
	DailyPnL         float64
	PositionCount    int
	RiskLevel        string
	TopPerformers    []PositionInfo
	WorstPerformers  []PositionInfo
}

// NewPortfolioManager creates a new portfolio manager
func NewPortfolioManager() *PortfolioManager {
	client := futures.NewClient(
		os.Getenv("BITGET_API_KEY"),
		os.Getenv("BITGET_SECRET_KEY"),
		os.Getenv("BITGET_PASSPHRASE"),
	)

	return &PortfolioManager{
		client:      client,
		ctx:         context.Background(),
		symbols:     []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT", "DOGEUSDT"},
		productType: "USDT-FUTURES",
		
		// Risk management settings
		maxPositions:     3,    // Maximum 3 concurrent positions
		maxRiskPerTrade:  0.02, // Max 2% risk per trade
		maxTotalExposure: 0.5,  // Max 50% total exposure
		
		positions: make(map[string]*PositionInfo),
	}
}

// Start begins portfolio management
func (pm *PortfolioManager) Start() error {
	log.Println("🚀 Starting Portfolio Manager...")

	for {
		// Update portfolio data
		if err := pm.updatePortfolioData(); err != nil {
			log.Printf("❌ Failed to update portfolio: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Generate and display report
		summary := pm.generatePortfolioSummary()
		pm.displayPortfolioReport(summary)

		// Execute portfolio rebalancing if needed
		if err := pm.rebalancePortfolio(); err != nil {
			log.Printf("⚠️ Rebalancing error: %v", err)
		}

		// Sleep before next update
		time.Sleep(60 * time.Second)
	}
}

// updatePortfolioData refreshes all portfolio information
func (pm *PortfolioManager) updatePortfolioData() error {
	// 1. Update account balance
	if err := pm.updateAccountBalance(); err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	// 2. Update all positions
	if err := pm.updatePositions(); err != nil {
		return fmt.Errorf("failed to update positions: %w", err)
	}

	// 3. Update current prices
	if err := pm.updatePrices(); err != nil {
		return fmt.Errorf("failed to update prices: %w", err)
	}

	return nil
}

// updateAccountBalance gets the latest account balance
func (pm *PortfolioManager) updateAccountBalance() error {
	accountService := account.NewAccountListService(pm.client)
	accounts, err := accountService.
		ProductType(account.ProductType(pm.productType)).
		Do(pm.ctx)

	if err != nil {
		return err
	}

	// Find USDT account
	for _, acc := range accounts {
		if acc.MarginCoin == "USDT" {
			balance, err := strconv.ParseFloat(acc.Available, 64)
			if err != nil {
				return fmt.Errorf("invalid balance format: %w", err)
			}
			pm.totalBalance = balance
			break
		}
	}

	return nil
}

// updatePositions retrieves all current positions
func (pm *PortfolioManager) updatePositions() error {
	positionService := position.NewAllPositionsService(pm.client)
	positions, err := positionService.
		ProductType(position.ProductType(pm.productType)).
		Do(pm.ctx)

	if err != nil {
		return err
	}

	// Clear existing positions
	pm.positions = make(map[string]*PositionInfo)
	pm.currentExposure = 0

	// Process each position
	for _, pos := range positions {
		size, _ := strconv.ParseFloat(pos.Size, 64)
		if size == 0 {
			continue // Skip empty positions
		}

		entryPrice, _ := strconv.ParseFloat(pos.AverageOpenPrice, 64)
		unrealizedPnL, _ := strconv.ParseFloat(pos.UnrealizedPL, 64)
		marginUsed, _ := strconv.ParseFloat(pos.Margin, 64)

		posInfo := &PositionInfo{
			Symbol:        pos.Symbol,
			Side:          pos.HoldSide,
			Size:          size,
			EntryPrice:    entryPrice,
			UnrealizedPnL: unrealizedPnL,
			MarginUsed:    marginUsed,
		}

		// Calculate risk level
		riskPct := marginUsed / pm.totalBalance
		if riskPct > 0.1 {
			posInfo.RiskLevel = "HIGH"
		} else if riskPct > 0.05 {
			posInfo.RiskLevel = "MEDIUM"
		} else {
			posInfo.RiskLevel = "LOW"
		}

		pm.positions[pos.Symbol] = posInfo
		pm.currentExposure += marginUsed
	}

	return nil
}

// updatePrices gets current market prices for all positions
func (pm *PortfolioManager) updatePrices() error {
	for symbol, posInfo := range pm.positions {
		tickerService := market.NewTickerService(pm.client)
		ticker, err := tickerService.
			Symbol(symbol).
			ProductType(market.ProductType(pm.productType)).
			Do(pm.ctx)

		if err != nil {
			log.Printf("⚠️ Failed to get price for %s: %v", symbol, err)
			continue
		}

		currentPrice, err := strconv.ParseFloat(ticker.LastPr, 64)
		if err != nil {
			continue
		}

		posInfo.CurrentPrice = currentPrice
	}

	return nil
}

// generatePortfolioSummary creates comprehensive portfolio analytics
func (pm *PortfolioManager) generatePortfolioSummary() *PortfolioSummary {
	summary := &PortfolioSummary{
		TotalBalance:     pm.totalBalance,
		AvailableBalance: pm.totalBalance - pm.currentExposure,
		TotalExposure:    pm.currentExposure,
		PositionCount:    len(pm.positions),
	}

	// Calculate total unrealized PnL
	var allPositions []PositionInfo
	for _, pos := range pm.positions {
		summary.UnrealizedPnL += pos.UnrealizedPnL
		allPositions = append(allPositions, *pos)
	}

	// Determine overall risk level
	exposurePct := pm.currentExposure / pm.totalBalance
	if exposurePct > 0.7 {
		summary.RiskLevel = "HIGH"
	} else if exposurePct > 0.4 {
		summary.RiskLevel = "MEDIUM"  
	} else {
		summary.RiskLevel = "LOW"
	}

	// Sort positions by performance
	sort.Slice(allPositions, func(i, j int) bool {
		return allPositions[i].UnrealizedPnL > allPositions[j].UnrealizedPnL
	})

	// Get top and worst performers (max 3 each)
	topCount := min(3, len(allPositions))
	if topCount > 0 {
		summary.TopPerformers = allPositions[:topCount]
		
		// Reverse for worst performers
		worstStart := max(0, len(allPositions)-3)
		summary.WorstPerformers = allPositions[worstStart:]
	}

	return summary
}

// displayPortfolioReport shows a detailed portfolio report
func (pm *PortfolioManager) displayPortfolioReport(summary *PortfolioSummary) {
	fmt.Println("\n" + "="*60)
	fmt.Println("📊 PORTFOLIO SUMMARY", time.Now().Format("15:04:05"))
	fmt.Println("="*60)

	// Overall metrics
	fmt.Printf("💰 Total Balance:      $%.2f USDT\n", summary.TotalBalance)
	fmt.Printf("💸 Available:          $%.2f USDT (%.1f%%)\n", 
		summary.AvailableBalance, 
		summary.AvailableBalance/summary.TotalBalance*100)
	fmt.Printf("📈 Total Exposure:     $%.2f USDT (%.1f%%)\n", 
		summary.TotalExposure,
		summary.TotalExposure/summary.TotalBalance*100)
	fmt.Printf("📊 Unrealized PnL:     $%.2f USDT", summary.UnrealizedPnL)
	if summary.UnrealizedPnL >= 0 {
		fmt.Print(" 🟢\n")
	} else {
		fmt.Print(" 🔴\n")
	}
	fmt.Printf("🎯 Active Positions:   %d/%d\n", summary.PositionCount, pm.maxPositions)
	fmt.Printf("⚠️ Risk Level:         %s\n", summary.RiskLevel)

	// Position details
	if len(pm.positions) > 0 {
		fmt.Println("\n📋 ACTIVE POSITIONS:")
		fmt.Println("-" * 60)
		
		for symbol, pos := range pm.positions {
			pnlPct := (pos.UnrealizedPnL / pos.MarginUsed) * 100
			fmt.Printf("%s: %s %.4f @ $%.2f → $%.2f (%.2f%%) %s\n",
				symbol,
				pos.Side,
				pos.Size,
				pos.EntryPrice,
				pos.CurrentPrice,
				pnlPct,
				pos.RiskLevel)
		}
	}

	// Top performers
	if len(summary.TopPerformers) > 0 {
		fmt.Println("\n🏆 TOP PERFORMERS:")
		for _, pos := range summary.TopPerformers {
			fmt.Printf("  %s: +$%.2f\n", pos.Symbol, pos.UnrealizedPnL)
		}
	}

	// Worst performers  
	if len(summary.WorstPerformers) > 0 {
		fmt.Println("\n📉 NEEDS ATTENTION:")
		for _, pos := range summary.WorstPerformers {
			if pos.UnrealizedPnL < 0 {
				fmt.Printf("  %s: $%.2f\n", pos.Symbol, pos.UnrealizedPnL)
			}
		}
	}

	fmt.Println("="*60)
}

// rebalancePortfolio executes portfolio rebalancing logic
func (pm *PortfolioManager) rebalancePortfolio() error {
	// 1. Risk management - close high-risk positions
	for symbol, pos := range pm.positions {
		// Close position if loss exceeds 5%
		lossPct := pos.UnrealizedPnL / pos.MarginUsed
		if lossPct < -0.05 {
			log.Printf("🛑 Closing %s due to 5%% loss limit", symbol)
			if err := pm.closePosition(symbol); err != nil {
				return err
			}
		}

		// Close position if profit exceeds 10% (take profit)
		if lossPct > 0.1 {
			log.Printf("🎯 Taking profit on %s (10%% gain)", symbol)
			if err := pm.closePosition(symbol); err != nil {
				return err
			}
		}
	}

	// 2. Opportunity detection - look for new positions if under limits
	if len(pm.positions) < pm.maxPositions && pm.currentExposure < pm.totalBalance*pm.maxTotalExposure {
		if err := pm.scanForOpportunities(); err != nil {
			return fmt.Errorf("opportunity scanning failed: %w", err)
		}
	}

	return nil
}

// closePosition closes a specific position
func (pm *PortfolioManager) closePosition(symbol string) error {
	pos, exists := pm.positions[symbol]
	if !exists {
		return fmt.Errorf("position %s not found", symbol)
	}

	closeService := position.NewClosePositionService(pm.client)
	_, err := closeService.
		Symbol(symbol).
		ProductType(position.ProductType(pm.productType)).
		HoldSide(pos.Side).
		MarginCoin("USDT").
		Do(pm.ctx)

	if err != nil {
		return fmt.Errorf("failed to close position %s: %w", symbol, err)
	}

	log.Printf("✅ Closed position: %s", symbol)
	delete(pm.positions, symbol)
	return nil
}

// scanForOpportunities looks for new trading opportunities
func (pm *PortfolioManager) scanForOpportunities() error {
	log.Println("🔍 Scanning for new opportunities...")

	// Get market data for all symbols
	for _, symbol := range pm.symbols {
		// Skip if already have position
		if _, exists := pm.positions[symbol]; exists {
			continue
		}

		// Get recent price action
		candleService := market.NewCandlestickService(pm.client)
		candles, err := candleService.
			Symbol(symbol).
			ProductType(market.ProductType(pm.productType)).
			Granularity("15m").
			Limit("20").
			Do(pm.ctx)

		if err != nil {
			continue
		}

		// Simple momentum detection
		if len(candles) >= 3 {
			recent := candles[len(candles)-1]
			prev := candles[len(candles)-3]

			// Check for 2% momentum
			if recent.Close > prev.Close*1.02 {
				log.Printf("📈 Momentum detected in %s - potential BUY signal", symbol)
				// Could implement position opening logic here
			} else if recent.Close < prev.Close*0.98 {
				log.Printf("📉 Downward momentum in %s - potential SELL signal", symbol)
				// Could implement short position logic here
			}
		}
	}

	return nil
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	// Verify environment variables
	requiredEnvs := []string{"BITGET_API_KEY", "BITGET_SECRET_KEY", "BITGET_PASSPHRASE"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", env)
		}
	}

	// Create and start portfolio manager
	pm := NewPortfolioManager()
	
	log.Println("🎯 Portfolio Manager starting...")
	log.Println("📊 Monitoring:", pm.symbols)
	log.Printf("⚙️ Max Positions: %d", pm.maxPositions)
	log.Printf("⚙️ Max Risk per Trade: %.1f%%", pm.maxRiskPerTrade*100)
	log.Printf("⚙️ Max Total Exposure: %.1f%%", pm.maxTotalExposure*100)

	if err := pm.Start(); err != nil {
		log.Fatalf("❌ Portfolio manager failed: %v", err)
	}
}