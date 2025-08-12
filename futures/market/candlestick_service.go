package market

import (
	"encoding/json"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"
	"strconv"
	
	"github.com/khanbekov/go-bitget/futures"
)

// CandlestickService provides methods for retrieving OHLCV (candlestick) data from Bitget.
// It supports configurable time ranges, granularities, and result limits.
type CandlestickService struct {
	c           ClientInterface
	symbol      string      // Trading pair symbol (e.g., "BTCUSDT")
	productType ProductType // Product type (USDT-FUTURES, COIN-FUTURES, etc.)
	granularity string      // Time interval (1m, 5m, 15m, 1h, 4h, 1d, etc.)
	limit       string      // Maximum number of candlesticks to return
	startTime   string      // Start time in milliseconds
	endTime     string      // End time in milliseconds
}

// Symbol sets the trading pair symbol for the candlestick request.
// Required parameter. Examples: "BTCUSDT", "ETHUSDT", "ADAUSDT".
func (s *CandlestickService) Symbol(symbol string) *CandlestickService {
	s.symbol = symbol
	return s
}

// ProductType sets the product type for the request.
// Required parameter. Use ProductTypeUSDTFutures or ProductTypeCoinFutures.
func (s *CandlestickService) ProductType(productType ProductType) *CandlestickService {
	s.productType = productType
	return s
}

// Granularity sets the time interval for candlesticks.
// Required parameter. Supported values: "1m", "5m", "15m", "30m", "1h", "4h", "6h", "12h", "1d", "3d", "1w", "1M".
func (s *CandlestickService) Granularity(granularity string) *CandlestickService {
	s.granularity = granularity
	return s
}

// StartTime sets the start time for the candlestick data range.
// Optional parameter. Should be provided as milliseconds timestamp string.
func (s *CandlestickService) StartTime(startTime string) *CandlestickService {
	s.startTime = startTime
	return s
}

// EndTime sets the end time for the candlestick data range.
// Optional parameter. Should be provided as milliseconds timestamp string.
func (s *CandlestickService) EndTime(endTime string) *CandlestickService {
	s.endTime = endTime
	return s
}

// Limit sets the maximum number of candlesticks to return.
// Optional parameter. Default and maximum value is 1000.
func (s *CandlestickService) Limit(limit string) *CandlestickService {
	s.limit = limit
	return s
}

// Do executes the candlestick data request and returns the results.
// Returns a slice of Candlestick objects containing OHLCV data.
//
// The context can be used for request cancellation and timeout control.
// Returns an error if the request fails or if required parameters are missing.
func (s *CandlestickService) Do(ctx context.Context) (candles []Candlestick, err error) {
	queryParams := url.Values{}

	// Set params of request
	queryParams.Set("symbol", s.symbol)
	queryParams.Set("productType", string(s.productType))
	queryParams.Set("granularity", s.granularity)
	if s.limit != "" {
		queryParams.Set("limit", s.limit)
	}
	if s.startTime != "" {
		queryParams.Set("startTime", s.startTime)
	}
	if s.endTime != "" {
		queryParams.Set("endTime", s.endTime)
	}

	// Make request to API
	var res *futures.ApiResponse

	res, _, err = s.c.CallAPI(ctx, "GET", futures.EndpointCandlesticks, queryParams, nil, false)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	err = jsoniter.Unmarshal(res.Data, &candles)

	if err != nil {
		return nil, err
	}

	return candles, nil
}

// Candlestick represents OHLCV (Open, High, Low, Close, Volume) data for a specific time period.
// Contains all essential price and volume information for technical analysis.
type Candlestick struct {
	CloseTime        int64   `json:"time"`             // Timestamp when the candlestick period ended (milliseconds)
	Open             float64 `json:"entry"`            // Opening price at the start of the period
	High             float64 `json:"high"`             // Highest price during the period
	Low              float64 `json:"low"`              // Lowest price during the period
	Close            float64 `json:"exit"`             // Closing price at the end of the period
	Volume           float64 `json:"volume"`           // Base asset volume traded during the period
	QuoteAssetVolume float64 `json:"quoteAssetVolume"` // Quote asset volume traded during the period
}

// UnmarshalJSON implements custom JSON unmarshaling for Candlestick.
// The Bitget API returns candlestick data as arrays of strings rather than objects,
// so this method converts the array format to the struct fields.
//
// Expected array format: [timestamp, open, high, low, close, volume, quoteVolume]
func (c *Candlestick) UnmarshalJSON(data []byte) error {
	var arr []string
	if err := json.Unmarshal(data, &arr); err != nil {
		return err
	}
	if len(arr) != 7 {
		return fmt.Errorf("expected 7 elements, received %d", len(arr))
	}

	// Парсим CloseTime
	closeTime, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return fmt.Errorf("incorrect CloseTime: %v", err)
	}

	// Парсим числовые значения
	parsers := []func(string) (any, error){
		func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
		func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
		func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
		func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
		func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
		func(s string) (any, error) { return strconv.ParseFloat(s, 64) },
	}

	// Присваиваем значения полям
	c.CloseTime = closeTime
	for i := 1; i < 7; i++ {
		val, err := parsers[i-1](arr[i])
		if err != nil {
			return fmt.Errorf("failed parsing element %d: %v", i+1, err)
		}
		switch i {
		case 1:
			c.Open = val.(float64)
		case 2:
			c.High = val.(float64)
		case 3:
			c.Low = val.(float64)
		case 4:
			c.Close = val.(float64)
		case 5:
			c.Volume = val.(float64)
		case 6:
			c.QuoteAssetVolume = val.(float64)
		}
	}

	return nil
}