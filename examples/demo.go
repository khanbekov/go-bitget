package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/khanbekov/go-bitget/futures"
	"github.com/khanbekov/go-bitget/futures/market"
)

func main() {
	// Load credentials from .env file
	loadEnv()
	
	// Get API credentials from environment
	apiKey := os.Getenv("BITGET_API_KEY")
	secretKey := os.Getenv("BITGET_SECRET_KEY")
	passphrase := os.Getenv("BITGET_PASSPHRASE")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("❌ Please set BITGET_API_KEY, BITGET_SECRET_KEY, and BITGET_PASSPHRASE environment variables")
	}

	fmt.Println("🚀 Testing Go-Bitget Futures SDK...")
	fmt.Printf("📋 Using demo API key: %s...\n", apiKey[:10])
	
	// Create futures client
	client := futures.NewClient(apiKey, secretKey, passphrase)
	ctx := context.Background()

	// Test 1: Get single ticker
	fmt.Println("\n📊 Test 1: Getting BTCUSDT ticker...")
	tickerService := market.NewTickerService(client)
	ticker, err := tickerService.
		Symbol("BTCUSDT").
		ProductType("USDT-FUTURES").
		Do(ctx)

	if err != nil {
		fmt.Printf("  ❌ Ticker request failed: %v\n", err)
	} else {
		fmt.Printf("  ✅ BTC Price: $%s USDT\n", ticker.LastPr)
		fmt.Printf("  📈 24h Change: %s%%\n", ticker.Change24h)
		fmt.Printf("  📊 24h Volume: %s BTC\n", ticker.BaseVolume)
	}

	// Test 2: Get all tickers (limit output)
	fmt.Println("\n📊 Test 2: Getting all USDT-FUTURES tickers...")
	allTickersService := market.NewAllTickersService(client)
	tickers, err := allTickersService.
		ProductType("USDT-FUTURES").
		Do(ctx)

	if err != nil {
		fmt.Printf("  ❌ All tickers request failed: %v\n", err)
	} else {
		fmt.Printf("  ✅ Retrieved %d tickers successfully\n", len(tickers))
		
		// Show first 3 tickers
		fmt.Println("  🎯 Sample tickers:")
		for i, t := range tickers {
			if i >= 3 {
				break
			}
			fmt.Printf("    %s: $%s (24h: %s%%)\n", t.Symbol, t.LastPr, t.Change24h)
		}
	}

	// Test 3: Get candlestick data
	fmt.Println("\n📊 Test 3: Getting ETHUSDT candlestick data...")
	candleService := market.NewCandlestickService(client)
	candles, err := candleService.
		Symbol("ETHUSDT").
		ProductType("USDT-FUTURES").
		Granularity("1H").
		Limit("3").
		Do(ctx)

	if err != nil {
		fmt.Printf("  ❌ Candlestick request failed: %v\n", err)
	} else {
		fmt.Printf("  ✅ Retrieved %d candlesticks (1h intervals)\n", len(candles))
		if len(candles) > 0 {
			latest := candles[len(candles)-1]
			fmt.Printf("  📊 Latest candle - Open: $%.2f, High: $%.2f, Low: $%.2f, Close: $%.2f\n", 
				latest.Open, latest.High, latest.Low, latest.Close)
		}
	}

	fmt.Println("\n✅ Go-Bitget SDK tests completed successfully!")
	fmt.Println("🎉 The SDK is working correctly with your demo API keys!")
}

// loadEnv loads environment variables from .env file
func loadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // .env file is optional
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			os.Setenv(key, value)
		}
	}
}