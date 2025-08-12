package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/account"
	"github.com/khanbekov/go-bitget/futures/market"
	"github.com/khanbekov/go-bitget/futures/position"
)

// RiskManager demonstrates comprehensive risk management for futures trading
type RiskManager struct {
	client      *futures.Client
	ctx         context.Context
	
	// Account configuration
	symbols        []string
	productType    string
	marginCoin     string
	
	// Risk limits
	maxPortfolioRisk   float64  // Maximum portfolio risk percentage
	maxSinglePosition  float64  // Maximum single position size (% of portfolio)
	maxDailyLoss       float64  // Maximum daily loss limit
	maxDrawdown        float64  // Maximum drawdown before stopping
	
	// Correlation limits
	maxCorrelatedExposure float64  // Maximum exposure to correlated assets
	correlationThreshold  float64  // Correlation threshold for risk grouping
	
	// Current state
	totalBalance       float64
	startDayBalance    float64
	maxBalanceToday    float64
	positions          map[string]*RiskPosition
	correlationMatrix  map[string]map[string]float64
}

// RiskPosition contains position information with risk metrics
type RiskPosition struct {
	Symbol        string
	Size          float64
	EntryPrice    float64
	CurrentPrice  float64
	UnrealizedPnL float64
	MarginUsed    float64
	
	// Risk metrics
	PositionRisk     float64  // % of portfolio at risk
	VaR95            float64  // 95% Value at Risk
	MaxLoss          float64  // Maximum possible loss
	RiskScore        float64  // Overall risk score (0-100)
	
	// Time-based metrics
	HoldingTime      time.Duration
	LastUpdate       time.Time
}

// RiskReport contains comprehensive risk assessment
type RiskReport struct {
	Timestamp            time.Time
	TotalPortfolioRisk   float64
	DailyPnL             float64
	MaxDrawdown          float64
	RiskScore            float64  // Overall portfolio risk score
	
	// Risk metrics
	TotalVaR             float64  // Portfolio VaR
	ConcentrationRisk    float64  // Position concentration risk
	CorrelationRisk      float64  // Cross-asset correlation risk
	LeverageRisk         float64  // Excessive leverage risk
	
	// Recommendations
	Warnings             []string
	Actions              []string
	RiskLevel           string   // LOW, MEDIUM, HIGH, CRITICAL
}

// NewRiskManager creates a new risk management system
func NewRiskManager() *RiskManager {
	client := futures.NewClient(
		os.Getenv("BITGET_API_KEY"),
		os.Getenv("BITGET_SECRET_KEY"),
		os.Getenv("BITGET_PASSPHRASE"),
	)

	return &RiskManager{
		client:      client,
		ctx:         context.Background(),
		symbols:     []string{"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT", "DOGEUSDT", "AVAXUSDT"},
		productType: "USDT-FUTURES",
		marginCoin:  "USDT",
		
		// Conservative risk limits
		maxPortfolioRisk:      0.10,  // 10% max portfolio risk
		maxSinglePosition:     0.05,  // 5% max single position
		maxDailyLoss:          0.03,  // 3% max daily loss
		maxDrawdown:           0.15,  // 15% max drawdown
		maxCorrelatedExposure: 0.15,  // 15% max correlated exposure
		correlationThreshold:  0.7,   // 70% correlation threshold
		
		positions:         make(map[string]*RiskPosition),
		correlationMatrix: make(map[string]map[string]float64),
	}
}

// Start begins the risk management monitoring
func (rm *RiskManager) Start() error {
	log.Println("🛡️ Starting Risk Management System...")
	
	// Initialize daily tracking
	if err := rm.initializeDailyTracking(); err != nil {
		return fmt.Errorf("failed to initialize daily tracking: %w", err)
	}
	
	// Calculate correlation matrix
	if err := rm.updateCorrelationMatrix(); err != nil {
		log.Printf("⚠️ Failed to update correlation matrix: %v", err)
	}

	// Main monitoring loop
	for {
		// Update portfolio state
		if err := rm.updatePortfolioState(); err != nil {
			log.Printf("❌ Failed to update portfolio: %v", err)
			time.Sleep(30 * time.Second)
			continue
		}

		// Generate risk report
		report := rm.generateRiskReport()
		rm.displayRiskReport(report)

		// Execute risk actions
		if err := rm.executeRiskActions(report); err != nil {
			log.Printf("⚠️ Risk action execution error: %v", err)
		}

		// Sleep between checks
		time.Sleep(60 * time.Second)
	}
}

// initializeDailyTracking sets up daily risk tracking
func (rm *RiskManager) initializeDailyTracking() error {
	// Get current balance
	accountService := account.NewAccountListService(rm.client)
	accounts, err := accountService.
		ProductType(account.ProductType(rm.productType)).
		Do(rm.ctx)

	if err != nil {
		return err
	}

	// Find USDT account
	for _, acc := range accounts {
		if acc.MarginCoin == rm.marginCoin {
			balance, err := strconv.ParseFloat(acc.Equity, 64)
			if err != nil {
				return fmt.Errorf("invalid balance format: %w", err)
			}
			rm.totalBalance = balance
			rm.startDayBalance = balance
			rm.maxBalanceToday = balance
			break
		}
	}

	log.Printf("📊 Daily tracking initialized - Starting balance: $%.2f", rm.startDayBalance)
	return nil
}

// updateCorrelationMatrix calculates asset correlations
func (rm *RiskManager) updateCorrelationMatrix() error {
	log.Println("📈 Updating correlation matrix...")
	
	// Initialize matrix
	rm.correlationMatrix = make(map[string]map[string]float64)
	for _, symbol1 := range rm.symbols {
		rm.correlationMatrix[symbol1] = make(map[string]float64)
	}

	// Get price data for all symbols
	priceData := make(map[string][]float64)
	
	for _, symbol := range rm.symbols {
		candleService := market.NewCandlestickService(rm.client)
		candles, err := candleService.
			Symbol(symbol).
			ProductType(market.ProductType(rm.productType)).
			Granularity("1h").
			Limit("100").
			Do(rm.ctx)

		if err != nil {
			log.Printf("⚠️ Failed to get candles for %s: %v", symbol, err)
			continue
		}

		// Extract closing prices
		var prices []float64
		for _, candle := range candles {
			prices = append(prices, candle.Close)
		}
		priceData[symbol] = prices
	}

	// Calculate correlations
	for _, symbol1 := range rm.symbols {
		for _, symbol2 := range rm.symbols {
			if prices1, exists1 := priceData[symbol1]; exists1 {
				if prices2, exists2 := priceData[symbol2]; exists2 {
					correlation := rm.calculateCorrelation(prices1, prices2)
					rm.correlationMatrix[symbol1][symbol2] = correlation
				}
			}
		}
	}

	return nil
}

// calculateCorrelation computes Pearson correlation coefficient
func (rm *RiskManager) calculateCorrelation(x, y []float64) float64 {
	if len(x) != len(y) || len(x) == 0 {
		return 0.0
	}

	n := float64(len(x))
	
	// Calculate means
	var sumX, sumY float64
	for i := 0; i < len(x); i++ {
		sumX += x[i]
		sumY += y[i]
	}
	meanX := sumX / n
	meanY := sumY / n

	// Calculate correlation
	var numerator, sumXX, sumYY float64
	for i := 0; i < len(x); i++ {
		dx := x[i] - meanX
		dy := y[i] - meanY
		numerator += dx * dy
		sumXX += dx * dx
		sumYY += dy * dy
	}

	denominator := math.Sqrt(sumXX * sumYY)
	if denominator == 0 {
		return 0.0
	}

	return numerator / denominator
}

// updatePortfolioState refreshes all position and risk data
func (rm *RiskManager) updatePortfolioState() error {
	// Update account balance
	if err := rm.updateAccountBalance(); err != nil {
		return err
	}

	// Update positions
	if err := rm.updatePositions(); err != nil {
		return err
	}

	// Calculate risk metrics for each position
	for _, pos := range rm.positions {
		rm.calculatePositionRisk(pos)
	}

	return nil
}

// updateAccountBalance gets latest account information
func (rm *RiskManager) updateAccountBalance() error {
	accountService := account.NewAccountListService(rm.client)
	accounts, err := accountService.
		ProductType(account.ProductType(rm.productType)).
		Do(rm.ctx)

	if err != nil {
		return err
	}

	for _, acc := range accounts {
		if acc.MarginCoin == rm.marginCoin {
			balance, err := strconv.ParseFloat(acc.Equity, 64)
			if err != nil {
				return fmt.Errorf("invalid balance format: %w", err)
			}
			rm.totalBalance = balance
			
			// Update daily max
			if balance > rm.maxBalanceToday {
				rm.maxBalanceToday = balance
			}
			break
		}
	}

	return nil
}

// updatePositions retrieves and processes all positions
func (rm *RiskManager) updatePositions() error {
	positionService := position.NewAllPositionsService(rm.client)
	positions, err := positionService.
		ProductType(position.ProductType(rm.productType)).
		Do(rm.ctx)

	if err != nil {
		return err
	}

	// Clear existing positions
	rm.positions = make(map[string]*RiskPosition)

	// Process each position
	for _, pos := range positions {
		size, _ := strconv.ParseFloat(pos.Size, 64)
		if size == 0 {
			continue // Skip empty positions
		}

		entryPrice, _ := strconv.ParseFloat(pos.AverageOpenPrice, 64)
		unrealizedPnL, _ := strconv.ParseFloat(pos.UnrealizedPL, 64)
		marginUsed, _ := strconv.ParseFloat(pos.Margin, 64)

		// Get current price
		currentPrice, err := rm.getCurrentPrice(pos.Symbol)
		if err != nil {
			log.Printf("⚠️ Failed to get price for %s: %v", pos.Symbol, err)
			currentPrice = entryPrice // Fallback to entry price
		}

		riskPos := &RiskPosition{
			Symbol:        pos.Symbol,
			Size:          size,
			EntryPrice:    entryPrice,
			CurrentPrice:  currentPrice,
			UnrealizedPnL: unrealizedPnL,
			MarginUsed:    marginUsed,
			LastUpdate:    time.Now(),
		}

		rm.positions[pos.Symbol] = riskPos
	}

	return nil
}

// getCurrentPrice fetches current market price for a symbol
func (rm *RiskManager) getCurrentPrice(symbol string) (float64, error) {
	tickerService := market.NewTickerService(rm.client)
	ticker, err := tickerService.
		Symbol(symbol).
		ProductType(market.ProductType(rm.productType)).
		Do(rm.ctx)

	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(ticker.LastPr, 64)
}

// calculatePositionRisk computes risk metrics for a position
func (rm *RiskManager) calculatePositionRisk(pos *RiskPosition) {
	// Position risk as % of portfolio
	pos.PositionRisk = pos.MarginUsed / rm.totalBalance

	// Simple VaR calculation (95% confidence, assuming 2% daily volatility)
	dailyVolatility := 0.02  // 2% daily volatility assumption
	pos.VaR95 = pos.MarginUsed * dailyVolatility * 1.96  // 95% confidence

	// Maximum possible loss (100% of margin)
	pos.MaxLoss = pos.MarginUsed

	// Risk score calculation (0-100)
	riskScore := 0.0
	
	// Position size component (0-40 points)
	sizeRisk := (pos.PositionRisk / rm.maxSinglePosition) * 40
	if sizeRisk > 40 {
		sizeRisk = 40
	}
	riskScore += sizeRisk

	// P&L component (0-30 points)
	pnlRisk := 0.0
	if pos.UnrealizedPnL < 0 {
		lossRatio := math.Abs(pos.UnrealizedPnL) / pos.MarginUsed
		pnlRisk = lossRatio * 30
		if pnlRisk > 30 {
			pnlRisk = 30
		}
	}
	riskScore += pnlRisk

	// Holding time component (0-30 points)
	holdingHours := time.Since(pos.LastUpdate).Hours()
	if holdingHours > 168 { // More than a week
		riskScore += 30
	} else if holdingHours > 72 { // More than 3 days
		riskScore += 20
	} else if holdingHours > 24 { // More than a day
		riskScore += 10
	}

	pos.RiskScore = riskScore
}

// generateRiskReport creates comprehensive risk assessment
func (rm *RiskManager) generateRiskReport() *RiskReport {
	report := &RiskReport{
		Timestamp: time.Now(),
		Warnings:  make([]string, 0),
		Actions:   make([]string, 0),
	}

	// Calculate daily P&L
	report.DailyPnL = rm.totalBalance - rm.startDayBalance

	// Calculate maximum drawdown
	report.MaxDrawdown = (rm.maxBalanceToday - rm.totalBalance) / rm.maxBalanceToday

	// Calculate total portfolio risk
	totalMarginUsed := 0.0
	totalVaR := 0.0
	highRiskPositions := 0

	for _, pos := range rm.positions {
		totalMarginUsed += pos.MarginUsed
		totalVaR += pos.VaR95
		if pos.RiskScore > 60 {
			highRiskPositions++
		}
	}

	report.TotalPortfolioRisk = totalMarginUsed / rm.totalBalance
	report.TotalVaR = totalVaR

	// Calculate concentration risk
	if len(rm.positions) > 0 {
		var positions []float64
		for _, pos := range rm.positions {
			positions = append(positions, pos.PositionRisk)
		}
		sort.Float64s(positions)
		
		// Concentration risk = largest position / average position
		if len(positions) > 0 {
			largestPosition := positions[len(positions)-1]
			avgPosition := totalMarginUsed / float64(len(positions)) / rm.totalBalance
			if avgPosition > 0 {
				report.ConcentrationRisk = largestPosition / avgPosition
			}
		}
	}

	// Calculate correlation risk
	report.CorrelationRisk = rm.calculateCorrelationRisk()

	// Calculate leverage risk
	report.LeverageRisk = totalMarginUsed / rm.totalBalance

	// Overall risk score
	riskComponents := []float64{
		report.TotalPortfolioRisk * 100,
		report.ConcentrationRisk * 20,
		report.CorrelationRisk * 30,
		math.Abs(report.MaxDrawdown) * 100,
		float64(highRiskPositions) * 10,
	}
	
	totalRiskScore := 0.0
	for _, component := range riskComponents {
		totalRiskScore += component
	}
	report.RiskScore = totalRiskScore / float64(len(riskComponents))

	// Determine risk level
	if report.RiskScore > 80 {
		report.RiskLevel = "CRITICAL"
	} else if report.RiskScore > 60 {
		report.RiskLevel = "HIGH"
	} else if report.RiskScore > 40 {
		report.RiskLevel = "MEDIUM"
	} else {
		report.RiskLevel = "LOW"
	}

	// Generate warnings and actions
	rm.generateRiskWarnings(report)
	rm.generateRiskActions(report)

	return report
}

// calculateCorrelationRisk assesses cross-asset correlation risk
func (rm *RiskManager) calculateCorrelationRisk() float64 {
	if len(rm.positions) < 2 {
		return 0.0
	}

	totalCorrelationRisk := 0.0
	comparisons := 0

	// Check all position pairs for high correlation
	symbols := make([]string, 0, len(rm.positions))
	for symbol := range rm.positions {
		symbols = append(symbols, symbol)
	}

	for i := 0; i < len(symbols); i++ {
		for j := i + 1; j < len(symbols); j++ {
			symbol1, symbol2 := symbols[i], symbols[j]
			
			if correlationMap, exists := rm.correlationMatrix[symbol1]; exists {
				if correlation, exists := correlationMap[symbol2]; exists {
					if math.Abs(correlation) > rm.correlationThreshold {
						// High correlation found - calculate combined risk
						pos1Risk := rm.positions[symbol1].PositionRisk
						pos2Risk := rm.positions[symbol2].PositionRisk
						combinedRisk := pos1Risk + pos2Risk
						
						// Risk increases with correlation strength
						correlationMultiplier := math.Abs(correlation)
						totalCorrelationRisk += combinedRisk * correlationMultiplier
					}
				}
			}
			comparisons++
		}
	}

	if comparisons > 0 {
		return totalCorrelationRisk / float64(comparisons)
	}
	
	return 0.0
}

// generateRiskWarnings creates warning messages based on risk levels
func (rm *RiskManager) generateRiskWarnings(report *RiskReport) {
	// Portfolio-level warnings
	if report.TotalPortfolioRisk > rm.maxPortfolioRisk {
		report.Warnings = append(report.Warnings, 
			fmt.Sprintf("Portfolio risk (%.1f%%) exceeds limit (%.1f%%)", 
				report.TotalPortfolioRisk*100, rm.maxPortfolioRisk*100))
	}

	// Daily loss warning
	dailyLossPct := report.DailyPnL / rm.startDayBalance
	if dailyLossPct < -rm.maxDailyLoss {
		report.Warnings = append(report.Warnings, 
			fmt.Sprintf("Daily loss (%.1f%%) exceeds limit (%.1f%%)", 
				dailyLossPct*100, rm.maxDailyLoss*100))
	}

	// Drawdown warning
	if report.MaxDrawdown > rm.maxDrawdown {
		report.Warnings = append(report.Warnings, 
			fmt.Sprintf("Drawdown (%.1f%%) exceeds limit (%.1f%%)", 
				report.MaxDrawdown*100, rm.maxDrawdown*100))
	}

	// Position-specific warnings
	for symbol, pos := range rm.positions {
		if pos.PositionRisk > rm.maxSinglePosition {
			report.Warnings = append(report.Warnings, 
				fmt.Sprintf("%s position size (%.1f%%) exceeds limit (%.1f%%)", 
					symbol, pos.PositionRisk*100, rm.maxSinglePosition*100))
		}

		if pos.RiskScore > 70 {
			report.Warnings = append(report.Warnings, 
				fmt.Sprintf("%s has high risk score: %.0f/100", symbol, pos.RiskScore))
		}
	}

	// Correlation warnings
	if report.CorrelationRisk > 0.1 {
		report.Warnings = append(report.Warnings, 
			fmt.Sprintf("High correlation risk detected: %.1f%%", report.CorrelationRisk*100))
	}
}

// generateRiskActions creates recommended actions based on risk assessment
func (rm *RiskManager) generateRiskActions(report *RiskReport) {
	// Critical risk level actions
	if report.RiskLevel == "CRITICAL" {
		report.Actions = append(report.Actions, "IMMEDIATE: Consider closing high-risk positions")
		report.Actions = append(report.Actions, "IMMEDIATE: Stop opening new positions")
	}

	// High portfolio risk
	if report.TotalPortfolioRisk > rm.maxPortfolioRisk*1.2 {
		report.Actions = append(report.Actions, "Reduce overall position sizes by 20%")
	}

	// Daily loss limit approaching
	dailyLossPct := math.Abs(report.DailyPnL / rm.startDayBalance)
	if dailyLossPct > rm.maxDailyLoss*0.8 {
		report.Actions = append(report.Actions, "Approaching daily loss limit - consider position review")
	}

	// High-risk individual positions
	for symbol, pos := range rm.positions {
		if pos.RiskScore > 80 {
			report.Actions = append(report.Actions, fmt.Sprintf("Consider closing %s (risk: %.0f)", symbol, pos.RiskScore))
		} else if pos.RiskScore > 60 {
			report.Actions = append(report.Actions, fmt.Sprintf("Monitor %s closely (risk: %.0f)", symbol, pos.RiskScore))
		}
	}

	// Correlation risk actions
	if report.CorrelationRisk > 0.15 {
		report.Actions = append(report.Actions, "Reduce positions in highly correlated assets")
	}
}

// displayRiskReport shows formatted risk assessment
func (rm *RiskManager) displayRiskReport(report *RiskReport) {
	fmt.Println("\n" + "="*70)
	fmt.Printf("🛡️ RISK MANAGEMENT REPORT - %s\n", report.Timestamp.Format("15:04:05"))
	fmt.Println("="*70)

	// Risk level indicator
	riskEmoji := map[string]string{
		"LOW":      "🟢",
		"MEDIUM":   "🟡", 
		"HIGH":     "🟠",
		"CRITICAL": "🔴",
	}
	fmt.Printf("⚠️ Risk Level: %s %s (Score: %.1f/100)\n", 
		riskEmoji[report.RiskLevel], report.RiskLevel, report.RiskScore)

	// Key metrics
	fmt.Printf("💰 Portfolio Value: $%.2f USDT\n", rm.totalBalance)
	fmt.Printf("📊 Daily P&L: $%.2f (%.2f%%)\n", 
		report.DailyPnL, (report.DailyPnL/rm.startDayBalance)*100)
	fmt.Printf("📉 Max Drawdown: %.2f%%\n", report.MaxDrawdown*100)
	fmt.Printf("⚡ Portfolio Risk: %.1f%% (Limit: %.1f%%)\n", 
		report.TotalPortfolioRisk*100, rm.maxPortfolioRisk*100)
	fmt.Printf("🎯 VaR (95%%): $%.2f\n", report.TotalVaR)

	// Position details
	if len(rm.positions) > 0 {
		fmt.Println("\n📋 POSITION RISK ANALYSIS:")
		fmt.Println("-"*70)
		fmt.Printf("%-10s %-8s %-10s %-8s %-6s %-8s\n", 
			"SYMBOL", "SIZE%", "P&L", "MARGIN", "RISK", "VaR")
		fmt.Println("-"*70)

		// Sort positions by risk score
		var sortedPositions []*RiskPosition
		for _, pos := range rm.positions {
			sortedPositions = append(sortedPositions, pos)
		}
		sort.Slice(sortedPositions, func(i, j int) bool {
			return sortedPositions[i].RiskScore > sortedPositions[j].RiskScore
		})

		for _, pos := range sortedPositions {
			riskIcon := "🟢"
			if pos.RiskScore > 70 {
				riskIcon = "🔴"
			} else if pos.RiskScore > 50 {
				riskIcon = "🟡"
			}

			fmt.Printf("%-10s %-7.1f%% $%-8.2f $%-6.0f %s%-5.0f $%-7.0f\n",
				pos.Symbol,
				pos.PositionRisk*100,
				pos.UnrealizedPnL,
				pos.MarginUsed,
				riskIcon,
				pos.RiskScore,
				pos.VaR95)
		}
	}

	// Warnings
	if len(report.Warnings) > 0 {
		fmt.Println("\n⚠️ RISK WARNINGS:")
		for _, warning := range report.Warnings {
			fmt.Printf("  • %s\n", warning)
		}
	}

	// Recommended actions
	if len(report.Actions) > 0 {
		fmt.Println("\n🎯 RECOMMENDED ACTIONS:")
		for _, action := range report.Actions {
			fmt.Printf("  • %s\n", action)
		}
	}

	fmt.Println("="*70)
}

// executeRiskActions performs automatic risk management actions
func (rm *RiskManager) executeRiskActions(report *RiskReport) error {
	// Only execute critical actions automatically
	if report.RiskLevel == "CRITICAL" {
		log.Println("🚨 CRITICAL RISK - Executing automatic risk management")
		
		// Example: Close positions with risk score > 90
		for symbol, pos := range rm.positions {
			if pos.RiskScore > 90 {
				log.Printf("🛑 Auto-closing high risk position: %s", symbol)
				// In a real implementation, you would call the close position API here
				// For demo purposes, we just log the action
			}
		}
	}

	// Log daily loss limit breach
	dailyLossPct := report.DailyPnL / rm.startDayBalance
	if dailyLossPct < -rm.maxDailyLoss {
		log.Printf("🚨 Daily loss limit breached: %.2f%% (limit: %.2f%%)", 
			dailyLossPct*100, rm.maxDailyLoss*100)
		// In production: halt all trading, send alerts, etc.
	}

	return nil
}

func main() {
	// Verify environment variables
	requiredEnvs := []string{"BITGET_API_KEY", "BITGET_SECRET_KEY", "BITGET_PASSPHRASE"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", env)
		}
	}

	// Create and start risk manager
	riskManager := NewRiskManager()

	log.Println("🛡️ Risk Management System Configuration:")
	log.Printf("📊 Monitoring symbols: %v", riskManager.symbols)
	log.Printf("⚙️ Max Portfolio Risk: %.1f%%", riskManager.maxPortfolioRisk*100)
	log.Printf("⚙️ Max Single Position: %.1f%%", riskManager.maxSinglePosition*100)
	log.Printf("⚙️ Max Daily Loss: %.1f%%", riskManager.maxDailyLoss*100)
	log.Printf("⚙️ Max Drawdown: %.1f%%", riskManager.maxDrawdown*100)

	if err := riskManager.Start(); err != nil {
		log.Fatalf("❌ Risk manager failed: %v", err)
	}
}