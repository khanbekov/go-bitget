package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/market"
)

// MarketDataStreamer demonstrates real-time market data collection and analysis
type MarketDataStreamer struct {
	client      *futures.Client
	ctx         context.Context
	symbols     []string
	productType string

	// Data storage
	priceData map[string]*PriceHistory
	dataLock  sync.RWMutex

	// Analytics
	alertThresholds map[string]*AlertConfig
	updateInterval  time.Duration
}

// PriceHistory stores historical price data for a symbol
type PriceHistory struct {
	Symbol       string
	CurrentPrice float64
	PriceHistory []PricePoint
	Volume24h    float64
	Change24h    float64
	HighLow24h   [2]float64 // [low, high]
	LastUpdate   time.Time

	// Technical indicators
	SMA20          float64
	SMA50          float64
	RSI            float64
	BollingerBands [3]float64 // [lower, middle, upper]
}

// PricePoint represents a price at a specific time
type PricePoint struct {
	Timestamp time.Time
	Price     float64
	Volume    float64
}

// AlertConfig defines price alert conditions
type AlertConfig struct {
	Symbol          string
	PriceAlertHigh  float64
	PriceAlertLow   float64
	VolumeAlertHigh float64
	ChangeAlert24h  float64 // Alert if 24h change exceeds this %
	TechnicalAlerts bool
	LastAlertTime   time.Time
	AlertCooldown   time.Duration
}

// MarketAlert represents a triggered alert
type MarketAlert struct {
	Symbol         string
	Type           string
	Message        string
	Timestamp      time.Time
	CurrentValue   float64
	ThresholdValue float64
}

// NewMarketDataStreamer creates a new market data streamer
func NewMarketDataStreamer() *MarketDataStreamer {
	client := futures.NewClient(
		os.Getenv("BITGET_API_KEY"),
		os.Getenv("BITGET_SECRET_KEY"),
		os.Getenv("BITGET_PASSPHRASE"),
	)

	symbols := []string{
		"BTCUSDT", "ETHUSDT", "ADAUSDT", "SOLUSDT", "DOGEUSDT",
		"BNBUSDT", "XRPUSDT", "MATICUSDT", "DOTUSDT", "LTCUSDT",
	}

	streamer := &MarketDataStreamer{
		client:          client,
		ctx:             context.Background(),
		symbols:         symbols,
		productType:     "USDT-FUTURES",
		priceData:       make(map[string]*PriceHistory),
		alertThresholds: make(map[string]*AlertConfig),
		updateInterval:  10 * time.Second,
	}

	// Initialize price data storage
	for _, symbol := range symbols {
		streamer.priceData[symbol] = &PriceHistory{
			Symbol:       symbol,
			PriceHistory: make([]PricePoint, 0, 100), // Keep last 100 points
		}

		// Configure default alerts
		streamer.alertThresholds[symbol] = &AlertConfig{
			Symbol:          symbol,
			ChangeAlert24h:  5.0, // Alert on 5% change
			VolumeAlertHigh: 0,   // Will be set dynamically
			TechnicalAlerts: true,
			AlertCooldown:   5 * time.Minute,
		}
	}

	return streamer
}

// Start begins the market data streaming
func (mds *MarketDataStreamer) Start() error {
	log.Println("🚀 Starting Market Data Streamer...")
	log.Printf("📊 Monitoring %d symbols: %v", len(mds.symbols), mds.symbols)

	// Start data collection goroutine
	go mds.dataCollectionLoop()

	// Start analysis and alert goroutine
	go mds.analysisLoop()

	// Start display goroutine
	go mds.displayLoop()

	// Keep main thread alive
	select {}
}

// dataCollectionLoop continuously collects market data
func (mds *MarketDataStreamer) dataCollectionLoop() {
	ticker := time.NewTicker(mds.updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mds.collectMarketData()
		case <-mds.ctx.Done():
			return
		}
	}
}

// collectMarketData fetches latest market data for all symbols
func (mds *MarketDataStreamer) collectMarketData() {
	// Get all tickers in one API call for efficiency
	allTickersService := market.NewAllTickersService(mds.client)
	tickers, err := allTickersService.
		ProductType(futures.ProductType(mds.productType)).
		Do(mds.ctx)

	if err != nil {
		log.Printf("❌ Failed to get tickers: %v", err)
		return
	}

	mds.dataLock.Lock()
	defer mds.dataLock.Unlock()

	// Update price data for each monitored symbol
	for _, ticker := range tickers {
		if priceHistory, exists := mds.priceData[ticker.Symbol]; exists {
			currentPrice, _ := strconv.ParseFloat(ticker.LastPr, 64)
			volume24h, _ := strconv.ParseFloat(ticker.BaseVolume, 64)
			change24h, _ := strconv.ParseFloat(ticker.Change24h, 64)
			high24h, _ := strconv.ParseFloat(ticker.High24h, 64)
			low24h, _ := strconv.ParseFloat(ticker.Low24h, 64)

			// Update current data
			priceHistory.CurrentPrice = currentPrice
			priceHistory.Volume24h = volume24h
			priceHistory.Change24h = change24h
			priceHistory.HighLow24h = [2]float64{low24h, high24h}
			priceHistory.LastUpdate = time.Now()

			// Add to price history
			point := PricePoint{
				Timestamp: time.Now(),
				Price:     currentPrice,
				Volume:    volume24h,
			}

			priceHistory.PriceHistory = append(priceHistory.PriceHistory, point)

			// Keep only last 100 points
			if len(priceHistory.PriceHistory) > 100 {
				priceHistory.PriceHistory = priceHistory.PriceHistory[1:]
			}

			// Calculate technical indicators
			mds.calculateTechnicalIndicators(priceHistory)
		}
	}
}

// calculateTechnicalIndicators computes technical analysis indicators
func (mds *MarketDataStreamer) calculateTechnicalIndicators(ph *PriceHistory) {
	if len(ph.PriceHistory) < 20 {
		return // Need at least 20 data points
	}

	prices := make([]float64, len(ph.PriceHistory))
	for i, point := range ph.PriceHistory {
		prices[i] = point.Price
	}

	// Simple Moving Averages
	if len(prices) >= 20 {
		ph.SMA20 = calculateSMA(prices, 20)
	}
	if len(prices) >= 50 {
		ph.SMA50 = calculateSMA(prices, 50)
	}

	// RSI
	if len(prices) >= 14 {
		ph.RSI = calculateRSI(prices, 14)
	}

	// Bollinger Bands (20-period, 2 std dev)
	if len(prices) >= 20 {
		ph.BollingerBands = calculateBollingerBands(prices, 20, 2.0)
	}
}

// analysisLoop performs continuous market analysis and alert checking
func (mds *MarketDataStreamer) analysisLoop() {
	ticker := time.NewTicker(30 * time.Second) // Analyze every 30 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mds.checkAlerts()
			mds.performMarketAnalysis()
		case <-mds.ctx.Done():
			return
		}
	}
}

// checkAlerts monitors for price alerts and technical signals
func (mds *MarketDataStreamer) checkAlerts() {
	mds.dataLock.RLock()
	defer mds.dataLock.RUnlock()

	now := time.Now()

	for symbol, priceHistory := range mds.priceData {
		alertConfig := mds.alertThresholds[symbol]

		// Skip if in cooldown period
		if now.Sub(alertConfig.LastAlertTime) < alertConfig.AlertCooldown {
			continue
		}

		var alerts []MarketAlert

		// Price-based alerts
		if alertConfig.PriceAlertHigh > 0 && priceHistory.CurrentPrice >= alertConfig.PriceAlertHigh {
			alerts = append(alerts, MarketAlert{
				Symbol:         symbol,
				Type:           "PRICE_HIGH",
				Message:        fmt.Sprintf("Price above $%.2f", alertConfig.PriceAlertHigh),
				Timestamp:      now,
				CurrentValue:   priceHistory.CurrentPrice,
				ThresholdValue: alertConfig.PriceAlertHigh,
			})
		}

		if alertConfig.PriceAlertLow > 0 && priceHistory.CurrentPrice <= alertConfig.PriceAlertLow {
			alerts = append(alerts, MarketAlert{
				Symbol:         symbol,
				Type:           "PRICE_LOW",
				Message:        fmt.Sprintf("Price below $%.2f", alertConfig.PriceAlertLow),
				Timestamp:      now,
				CurrentValue:   priceHistory.CurrentPrice,
				ThresholdValue: alertConfig.PriceAlertLow,
			})
		}

		// 24h change alert
		if abs(priceHistory.Change24h) >= alertConfig.ChangeAlert24h {
			direction := "UP"
			if priceHistory.Change24h < 0 {
				direction = "DOWN"
			}
			alerts = append(alerts, MarketAlert{
				Symbol:         symbol,
				Type:           "VOLUME_SPIKE",
				Message:        fmt.Sprintf("24h change: %s %.2f%%", direction, abs(priceHistory.Change24h)),
				Timestamp:      now,
				CurrentValue:   abs(priceHistory.Change24h),
				ThresholdValue: alertConfig.ChangeAlert24h,
			})
		}

		// Technical alerts
		if alertConfig.TechnicalAlerts {
			// RSI alerts
			if priceHistory.RSI > 70 {
				alerts = append(alerts, MarketAlert{
					Symbol:       symbol,
					Type:         "TECHNICAL",
					Message:      fmt.Sprintf("RSI Overbought: %.1f", priceHistory.RSI),
					Timestamp:    now,
					CurrentValue: priceHistory.RSI,
				})
			} else if priceHistory.RSI < 30 {
				alerts = append(alerts, MarketAlert{
					Symbol:       symbol,
					Type:         "TECHNICAL",
					Message:      fmt.Sprintf("RSI Oversold: %.1f", priceHistory.RSI),
					Timestamp:    now,
					CurrentValue: priceHistory.RSI,
				})
			}

			// Bollinger Band breakouts
			if priceHistory.CurrentPrice > priceHistory.BollingerBands[2] {
				alerts = append(alerts, MarketAlert{
					Symbol:    symbol,
					Type:      "TECHNICAL",
					Message:   "Price above upper Bollinger Band",
					Timestamp: now,
				})
			} else if priceHistory.CurrentPrice < priceHistory.BollingerBands[0] {
				alerts = append(alerts, MarketAlert{
					Symbol:    symbol,
					Type:      "TECHNICAL",
					Message:   "Price below lower Bollinger Band",
					Timestamp: now,
				})
			}
		}

		// Display alerts
		for _, alert := range alerts {
			mds.displayAlert(alert)
			alertConfig.LastAlertTime = now
		}
	}
}

// performMarketAnalysis generates market insights
func (mds *MarketDataStreamer) performMarketAnalysis() {
	mds.dataLock.RLock()
	defer mds.dataLock.RUnlock()

	// Find top movers
	var symbols []string
	for symbol := range mds.priceData {
		symbols = append(symbols, symbol)
	}

	// Sort by 24h change
	sort.Slice(symbols, func(i, j int) bool {
		return mds.priceData[symbols[i]].Change24h > mds.priceData[symbols[j]].Change24h
	})

	log.Printf("📊 Top gainers: %s (+%.2f%%), %s (+%.2f%%), %s (+%.2f%%)",
		symbols[0], mds.priceData[symbols[0]].Change24h,
		symbols[1], mds.priceData[symbols[1]].Change24h,
		symbols[2], mds.priceData[symbols[2]].Change24h,
	)

	// Bottom movers
	n := len(symbols)
	log.Printf("📉 Top losers: %s (%.2f%%), %s (%.2f%%), %s (%.2f%%)",
		symbols[n-1], mds.priceData[symbols[n-1]].Change24h,
		symbols[n-2], mds.priceData[symbols[n-2]].Change24h,
		symbols[n-3], mds.priceData[symbols[n-3]].Change24h,
	)
}

// displayLoop shows real-time market data dashboard
func (mds *MarketDataStreamer) displayLoop() {
	ticker := time.NewTicker(15 * time.Second) // Update display every 15 seconds
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mds.displayMarketDashboard()
		case <-mds.ctx.Done():
			return
		}
	}
}

// displayMarketDashboard shows formatted market data
func (mds *MarketDataStreamer) displayMarketDashboard() {
	mds.dataLock.RLock()
	defer mds.dataLock.RUnlock()

	fmt.Print("\033[2J\033[H") // Clear screen
	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("📊 REAL-TIME MARKET DASHBOARD", time.Now().Format("15:04:05"))
	fmt.Println(strings.Repeat("=", 80))

	// Create sorted slice for consistent display
	var symbols []string
	for symbol := range mds.priceData {
		symbols = append(symbols, symbol)
	}
	sort.Strings(symbols)

	fmt.Printf("%-10s %-10s %-8s %-8s %-8s %-6s %-6s\n",
		"SYMBOL", "PRICE", "24H%", "RSI", "VOL(M)", "SMA20", "BB")
	fmt.Println(strings.Repeat("-", 80))

	for _, symbol := range symbols {
		ph := mds.priceData[symbol]

		// Format change with color indicator
		changeIcon := "🟢"
		if ph.Change24h < 0 {
			changeIcon = "🔴"
		}

		// Format volume in millions
		volumeM := ph.Volume24h / 1000000

		// RSI color indicator
		rsiIcon := "⚪"
		if ph.RSI > 70 {
			rsiIcon = "🔴" // Overbought
		} else if ph.RSI < 30 {
			rsiIcon = "🟢" // Oversold
		}

		fmt.Printf("%-10s $%-8.2f %s%-6.2f%% %s%-5.1f $%-8.0f $%-8.2f $%-6.2f\n",
			symbol,
			ph.CurrentPrice,
			changeIcon, ph.Change24h,
			rsiIcon, ph.RSI,
			volumeM,
			ph.SMA20,
			ph.BollingerBands[1], // Middle band (SMA)
		)
	}

	fmt.Println(strings.Repeat("=", 80))
	fmt.Println("Legend: 🟢=Up/Oversold 🔴=Down/Overbought ⚪=Neutral | RSI: <30=Oversold >70=Overbought")
}

// displayAlert shows formatted alerts
func (mds *MarketDataStreamer) displayAlert(alert MarketAlert) {
	icon := "⚠️"
	switch alert.Type {
	case "PRICE_HIGH":
		icon = "🚀"
	case "PRICE_LOW":
		icon = "🔻"
	case "VOLUME_SPIKE":
		icon = "📈"
	case "TECHNICAL":
		icon = "🔍"
	}

	log.Printf("%s ALERT [%s]: %s", icon, alert.Symbol, alert.Message)
}

// Technical indicator calculation functions

func calculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		return 0
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}
	return sum / float64(period)
}

func calculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		return 0
	}

	gains := 0.0
	losses := 0.0

	// Calculate initial average gain/loss
	for i := len(prices) - period; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains += change
		} else {
			losses -= change
		}
	}

	avgGain := gains / float64(period)
	avgLoss := losses / float64(period)

	if avgLoss == 0 {
		return 100
	}

	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

func calculateBollingerBands(prices []float64, period int, stdDev float64) [3]float64 {
	if len(prices) < period {
		return [3]float64{0, 0, 0}
	}

	// Calculate SMA
	sma := calculateSMA(prices, period)

	// Calculate standard deviation
	variance := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		diff := prices[i] - sma
		variance += diff * diff
	}
	variance /= float64(period)
	std := math.Sqrt(variance)

	return [3]float64{
		sma - (stdDev * std), // Lower band
		sma,                  // Middle band (SMA)
		sma + (stdDev * std), // Upper band
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	// Verify environment variables
	requiredEnvs := []string{"BITGET_API_KEY", "BITGET_SECRET_KEY", "BITGET_PASSPHRASE"}
	for _, env := range requiredEnvs {
		if os.Getenv(env) == "" {
			log.Fatalf("❌ Required environment variable %s is not set", env)
		}
	}

	// Create and start market data streamer
	streamer := NewMarketDataStreamer()

	if err := streamer.Start(); err != nil {
		log.Fatalf("❌ Market data streamer failed: %v", err)
	}
}
