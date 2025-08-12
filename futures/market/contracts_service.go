package market

import (
	jsoniter "github.com/json-iterator/go"
	"golang.org/x/net/context"
	"net/url"

	"github.com/khanbekov/go-bitget/futures"
)

// ContractsService retrieves contract configuration and trading rules
type ContractsService struct {
	c           futures.ClientInterface
	productType futures.ProductType
	symbol      string
}

// ProductType sets the product type for the request.
// Required parameter. Use ProductTypeUSDTFutures, ProductTypeUSDCFutures, or ProductTypeCOINFutures.
func (s *ContractsService) ProductType(productType futures.ProductType) *ContractsService {
	s.productType = productType
	return s
}

// Symbol sets the trading pair symbol for the contracts request.
// Optional parameter. If not provided, returns all contracts for the product type.
func (s *ContractsService) Symbol(symbol string) *ContractsService {
	s.symbol = symbol
	return s
}

// Do executes the contracts request and returns the results.
// Returns a slice of Contract objects containing contract specifications.
//
// The context can be used for request cancellation and timeout control.
// Returns an error if the request fails or if required parameters are missing.
func (s *ContractsService) Do(ctx context.Context) ([]*Contract, error) {
	queryParams := url.Values{}

	// Set required params
	queryParams.Set("productType", string(s.productType))

	// Set optional params
	if s.symbol != "" {
		queryParams.Set("symbol", s.symbol)
	}

	// Make request to API
	var res *futures.ApiResponse

	res, _, err := s.c.CallAPI(ctx, "GET", futures.EndpointContracts, queryParams, nil, false)

	if err != nil {
		return nil, err
	}

	// Unmarshal json from response
	var contracts []*Contract
	err = jsoniter.Unmarshal(res.Data, &contracts)

	if err != nil {
		return nil, err
	}

	return contracts, nil
}

// Contract represents contract configuration and trading rules
type Contract struct {
	// Symbol identifier (e.g., BTCUSDT)
	Symbol string `json:"symbol"`

	// Base coin (e.g., BTC)
	BaseCoin string `json:"baseCoin"`

	// Quote coin (e.g., USDT)
	QuoteCoin string `json:"quoteCoin"`

	// Buy limit price factor
	BuyLimitPriceRatio string `json:"buyLimitPriceRatio"`

	// Sell limit price factor
	SellLimitPriceRatio string `json:"sellLimitPriceRatio"`

	// Fee rate for taker orders
	FeeRateUpRatio string `json:"feeRateUpRatio"`

	// Maker fee rate
	MakerFeeRate string `json:"makerFeeRate"`

	// Taker fee rate
	TakerFeeRate string `json:"takerFeeRate"`

	// Open cost for maintaining positions
	OpenCostUpRatio string `json:"openCostUpRatio"`

	// Support margin coin
	SupportMarginCoins []string `json:"supportMarginCoins"`

	// Minimum trade amount
	MinTradeNum string `json:"minTradeNum"`

	// Price precision (decimal places)
	PriceEndStep string `json:"priceEndStep"`

	// Price scale (minimum price increment)
	PricePlace string `json:"pricePlace"`

	// Volume scale (minimum quantity increment)
	VolumePlace string `json:"volumePlace"`

	// Symbol type (perpetual for example)
	SymbolType string `json:"symbolType"`

	// Minimum order size
	MinTradeUSDT string `json:"minTradeUSDT"`

	// Quantity multiplier, the quantity of the order must be greater than minTradeNum and is a multiple of sizeMulti.
	SizeMultiplier string `json:"sizeMultiplier"`

	// MaxSymbolOrderNum
	MaxSymbolOrderNum string `json:"maxSymbolOrderNum"`

	// MaxProductOrderNum
	MaxProductOrderNum string `json:"maxProductOrderNum"`

	// MaxPositionNum
	MaxPositionNum string `json:"maxPositionNum"`

	// Symbol status (online/offline)
	SymbolStatus string `json:"symbolStatus"`

	// Offline time if symbol is offline
	OffTime string `json:"offTime"`

	// Limit open time
	LimitOpenTime string `json:"limitOpenTime"`

	// Delivery time
	DeliveryTime string `json:"deliveryTime"`

	// Delivery start time
	DeliveryStartTime string `json:"deliveryStartTime"`

	// Delivery start time
	DeliveryPeriod string `json:"deliveryPeriod"`

	// Launch time
	LaunchTime string `json:"launchTime"`

	// Fund interval (for funding rate calculation)
	FundInterval string `json:"fundInterval"`

	// Minimum leverage
	MinLever string `json:"minLever"`

	// Maximum leverage
	MaxLever string `json:"maxLever"`

	// Position limit
	PosLimit string `json:"posLimit"`

	// Maintenance margin rate
	MaintainTime string `json:"maintainTime"`

	// Open time
	OpenTime string `json:"openTime"`

	// Maximum market order quantity
	MaxMarketOrderQty string `json:"maxMarketOrderQty"`

	// Maximum order quantity
	MaxOrderQty string `json:"maxOrderQty"`
}
