package market

import (
	"github.com/khanbekov/go-bitget/common/client"
)

// Re-export common types to avoid importing futures package
type (
	ClientInterface = client.ClientInterface
	ApiResponse     = client.ApiResponse
)

// Futures-specific enums (duplicated to avoid import cycle)
type ProductType string

const (
	ProductTypeUSDTFutures ProductType = "USDT-FUTURES"
	ProductTypeCoinFutures ProductType = "COIN-FUTURES"
	ProductTypeUSDCFutures ProductType = "USDC-FUTURES"
)

// API Endpoints for market data operations
const (
	EndpointAllTickers          = "/api/v2/mix/market/tickers"
	EndpointCandlesticks        = "/api/v2/mix/market/candles"
	EndpointTicker              = "/api/v2/mix/market/ticker"
	EndpointMergeDepth          = "/api/v2/mix/market/merge-depth"
	EndpointContractConfig      = "/api/v2/mix/market/contracts"
	EndpointRecentTrades        = "/api/v2/mix/market/fills"
	EndpointCurrentFundingRate  = "/api/v2/mix/market/current-funding-rate"
	EndpointHistoryFundingRate  = "/api/v2/mix/market/history-funding-rate"
	EndpointOpenInterest        = "/api/v2/mix/market/open-interest"
	EndpointSymbolPrice         = "/api/v2/mix/market/symbol-price"
)

// Service Constructor Functions

// NewCandlestickService creates a new candlestick service.
func NewCandlestickService(client ClientInterface) *CandlestickService {
	return &CandlestickService{c: client}
}

// NewAllTickersService creates a new all tickers service.
func NewAllTickersService(client ClientInterface) *AllTickersService {
	return &AllTickersService{c: client}
}

// NewTickerService creates a new ticker service.
func NewTickerService(client ClientInterface) *TickerService {
	return &TickerService{c: client}
}

// NewOrderBookService creates a new order book service.
func NewOrderBookService(client ClientInterface) *OrderBookService {
	return &OrderBookService{c: client}
}

// NewRecentTradesService creates a new recent trades service.
func NewRecentTradesService(client ClientInterface) *RecentTradesService {
	return &RecentTradesService{c: client}
}

// NewCurrentFundingRateService creates a new current funding rate service.
func NewCurrentFundingRateService(client ClientInterface) *CurrentFundingRateService {
	return &CurrentFundingRateService{c: client}
}

// NewHistoryFundingRateService creates a new history funding rate service.
func NewHistoryFundingRateService(client ClientInterface) *HistoryFundingRateService {
	return &HistoryFundingRateService{c: client}
}

// NewOpenInterestService creates a new open interest service.
func NewOpenInterestService(client ClientInterface) *OpenInterestService {
	return &OpenInterestService{c: client}
}

// NewSymbolPriceService creates a new symbol price service.
func NewSymbolPriceService(client ClientInterface) *SymbolPriceService {
	return &SymbolPriceService{c: client}
}

// NewContractsService creates a new contracts service.
func NewContractsService(client ClientInterface) *ContractsService {
	return &ContractsService{c: client}
}